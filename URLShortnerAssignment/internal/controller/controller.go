package controller

import (
	"URLShortnerAssignment/internal/interfaces"
	"URLShortnerAssignment/internal/middleware"
	"URLShortnerAssignment/internal/service"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type URLController struct {
	svc interfaces.URLService
}

func NewURLController(svc interfaces.URLService) *URLController {
	return &URLController{svc: svc}
}

type ShortenRequest struct {
	LongURL string `json:"long_url" binding:"required,url"`
	Custom  string `json:"custom" binding:"omitempty,alphanum"`
}

func (c *URLController) RegisterRoutes(r *gin.Engine, limiter *middleware.LimiterStore) {
	r.POST("/shorten", middleware.RateLimitMiddleware(limiter), c.Shorten)
	r.GET("/stats/:code", c.Stats)
	r.GET("/:code", c.Redirect)
}

func (c *URLController) Shorten(ctx *gin.Context) {
	var req ShortenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := c.svc.Shorten(req.LongURL, req.Custom)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAlreadyExists):
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrNotFound):
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	baseURL := getBaseURL(ctx)
	ctx.JSON(http.StatusCreated, gin.H{
		"short_url":  fmt.Sprintf("%s/%s", baseURL, u.ShortCode),
		"short_code": u.ShortCode,
		"long_url":   u.LongURL,
	})
}

func (c *URLController) Stats(ctx *gin.Context) {
	code := ctx.Param("code")
	u, err := c.svc.Stats(code)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.JSON(http.StatusOK, u)
}

func (c *URLController) Redirect(ctx *gin.Context) {
	code := ctx.Param("code")
	u, err := c.svc.Resolve(code)
	if err != nil {
		ctx.String(http.StatusNotFound, "not found")
		return
	}
	ctx.Redirect(http.StatusFound, u.LongURL)
}

func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}
