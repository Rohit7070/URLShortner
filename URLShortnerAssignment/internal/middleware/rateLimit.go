package middleware

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type LimiterStore struct {
	limiters map[string]*rate.Limiter
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func NewLimiterStore(r rate.Limit, b int) *LimiterStore {
	return &LimiterStore{limiters: make(map[string]*rate.Limiter), r: r, b: b}
}

func (s *LimiterStore) getLimiter(ip string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()
	if l, ok := s.limiters[ip]; ok {
		return l
	}
	lim := rate.NewLimiter(s.r, s.b)
	s.limiters[ip] = lim
	return lim
}

func RateLimitMiddleware(store *LimiterStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		lim := store.getLimiter(ip)
		if !lim.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}
