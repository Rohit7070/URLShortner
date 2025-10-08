package test

import (
	"URLShortnerAssignment/internal/controller"
	"URLShortnerAssignment/internal/middleware"
	"URLShortnerAssignment/internal/models"
	"URLShortnerAssignment/internal/repository"
	"URLShortnerAssignment/internal/service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	// In-memory DB
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.URL{})

	// Init repo → service → controller
	repo := repository.NewRepository(db)
	svc := service.NewURLService(repo, "http://localhost:8080")
	ctrl := controller.NewURLController(svc)

	// Setup rate limiter
	limiter := middleware.NewLimiterStore(rate.Every(time.Second), 5)

	// Register routes
	r := gin.Default()
	ctrl.RegisterRoutes(r, limiter)

	return r
}

func TestURLShortener_AllEndpoints(t *testing.T) {
	r := setupRouter()

	// 1️⃣ POST /shorten — create short URL
	reqBody, _ := json.Marshal(map[string]string{"long_url": "https://golang.org"})
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code, "Expected 201 Created")

	var createResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &createResp)
	assert.NotEmpty(t, createResp["short_code"], "short_code must be returned")
	assert.NotEmpty(t, createResp["short_url"], "short_url must be returned")

	shortCode := createResp["short_code"]

	// 2️⃣ GET /stats/:code — fetch stats
	statsReq := httptest.NewRequest(http.MethodGet, "/stats/"+shortCode, nil)
	statsW := httptest.NewRecorder()
	r.ServeHTTP(statsW, statsReq)

	assert.Equal(t, http.StatusOK, statsW.Code, "Expected 200 OK for stats")
	assert.Contains(t, statsW.Body.String(), `"long_url":"https://golang.org"`)

	// 3️⃣ GET /:code — redirect
	redirectReq := httptest.NewRequest(http.MethodGet, "/"+shortCode, nil)
	redirectW := httptest.NewRecorder()
	r.ServeHTTP(redirectW, redirectReq)

	assert.Equal(t, http.StatusFound, redirectW.Code, "Expected 302 Redirect")
	assert.Equal(t, "https://golang.org", redirectW.Header().Get("Location"))

	// 4️⃣ Check stats again → hits incremented
	statsReq2 := httptest.NewRequest(http.MethodGet, "/stats/"+shortCode, nil)
	statsW2 := httptest.NewRecorder()
	r.ServeHTTP(statsW2, statsReq2)
	assert.Contains(t, statsW2.Body.String(), `"hits":1`)
}

func TestRateLimiting(t *testing.T) {
	r := setupRouter()

	for i := 0; i < 5; i++ {
		body, _ := json.Marshal(map[string]string{"long_url": "https://example.com"})
		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusTooManyRequests)
	}

	// Exceed limit
	body, _ := json.Marshal(map[string]string{"long_url": "https://blocked.com"})
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
}
