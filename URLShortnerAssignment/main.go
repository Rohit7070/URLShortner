package main

import (
	"URLShortnerAssignment/internal/controller"
	"URLShortnerAssignment/internal/middleware"
	"URLShortnerAssignment/internal/models"
	"URLShortnerAssignment/internal/repository"
	"URLShortnerAssignment/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func main() {
	// Gin mode setup
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// DB setup
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "urls.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// Auto-migrate model
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// Initialize layers
	repo := repository.NewRepository(db)
	baseDomain := os.Getenv("BASE_DOMAIN")
	if baseDomain == "" {
		baseDomain = "http://localhost:8080"
	}
	svc := service.NewURLService(repo, baseDomain)
	ctrl := controller.NewURLController(svc)

	// Rate limiter setup (5 req/min)
	limiter := middleware.NewLimiterStore(rate.Every(time.Minute/5), 5)

	// Gin router
	r := gin.Default()
	ctrl.RegisterRoutes(r, limiter)

	// Start server
	addr := ":8090"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	log.Printf("🚀 URL Shortener running on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
