package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter is a middleware that limits the rate of requests.
func RateLimiter(limit rate.Limit, burst int) gin.HandlerFunc {
	var (
		visitors     = make(map[string]*rate.Limiter)
		mu           sync.Mutex
		cleanupInterval = 1 * time.Minute // Clean up visitors every minute
	)

	// Start a goroutine to clean up visitors periodically
	go func() {
		for {
			time.Sleep(cleanupInterval)
			mu.Lock()
			for ip, lastSeen := range visitors {
				if time.Since(lastSeen.lastEvent) > cleanupInterval {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		limiter, exists := visitors[ip]
		if !exists {
			limiter = &visitor{
				limiter:   rate.NewLimiter(limit, burst),
				lastEvent: time.Now(),
			}
			visitors[ip] = limiter
		}
		limiter.lastEvent = time.Now()
		mu.Unlock()

		if !limiter.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}

type visitor struct {
	limiter   *rate.Limiter
	lastEvent time.Time
}
