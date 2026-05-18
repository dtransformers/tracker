package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func RateLimit(RateLimitWindow time.Duration, RateLimitRequests int) gin.HandlerFunc {
	// Create a rate limiter store (in-memory for single instance)
	// For distributed systems, use Redis store
	rate := limiter.Rate{
		Period: RateLimitWindow,
		Limit:  int64(RateLimitRequests),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// metricsMiddleware tracks request metrics.
func Metrics(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// Log request
		logger.Info("request completed",
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// You can also emit metrics to your monitoring system here
		// metrics.RecordRequest(path, c.Request.Method, status, duration)
	}
}

// validationMiddleware validates common request parameters.
func Validation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate content type for POST/PATCH requests
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "" {
				c.JSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// corsMiddleware configures CORS settings.
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
