package v1

import (
	"log/slog"
	"net/http"
	"time"
	"tracking/internal/ent"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

type Module struct {
	handler *Handler
	service *Service
	repo    Repository
	logger  *slog.Logger
	config  RouterConfig
}

func NewModule(client *ent.Client, logger *slog.Logger, cfg RouterConfig) *Module {
	repo := NewRepository(client)
	svc := NewServiceWithConfig(repo, logger, cfg.ServiceConfig)
	handler := NewHandler(svc, logger)

	return &Module{
		handler: handler,
		service: svc,
		repo:    repo,
		logger:  logger,
		config:  cfg,
	}
}

func (m *Module) Register(router *gin.RouterGroup) {
	sites := router.Group("/sites")

	m.applyMiddlewares(sites)

	m.registerRoutes(sites)
}

func (m *Module) applyMiddlewares(group *gin.RouterGroup) {

	group.Use(m.metricsMiddleware())

	if m.config.EnableRateLimit {
		group.Use(m.rateLimitMiddleware())
	}

	group.Use(validationMiddleware())
}

func (m *Module) registerRoutes(group *gin.RouterGroup) {
	group.POST("", m.handler.Create)
	group.GET("", m.handler.List)

	group.GET("/:uuid", m.handler.Get)
	group.PATCH("/:uuid", m.handler.Update)
	group.DELETE("/:uuid", m.handler.Delete)

	// Special actions
	group.POST("/:uuid/regenerate-keys", m.handler.RegenerateKeys)
}

// rateLimitMiddleware creates a rate limiting middleware using ulule/limiter.
func (m *Module) rateLimitMiddleware() gin.HandlerFunc {
	// Create a rate limiter store (in-memory for single instance)
	// For distributed systems, use Redis store
	rate := limiter.Rate{
		Period: m.config.RateLimitWindow,
		Limit:  int64(m.config.RateLimitRequests),
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// metricsMiddleware tracks request metrics.
func (m *Module) metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Calculate metrics
		duration := time.Since(start)
		status := c.Writer.Status()

		// Log request
		m.logger.Info("request completed",
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
func validationMiddleware() gin.HandlerFunc {
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
func corsMiddleware() gin.HandlerFunc {
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

// Health returns the health status of the module.
func (m *Module) Health() map[string]interface{} {
	return map[string]interface{}{
		"status":  "healthy",
		"service": "sites-v1",
	}
}

// Shutdown performs graceful shutdown of the module.
func (m *Module) Shutdown() error {
	m.logger.Info("shutting down sites module")
	// Perform any necessary cleanup
	return nil
}

// Register is a convenience function to register the sites module with default config.
func Register(router *gin.RouterGroup, client *ent.Client) {
	RegisterWithLogger(router, client, slog.Default())
}

// RegisterWithLogger is a convenience function to register the sites module with a logger and default config.
func RegisterWithLogger(router *gin.RouterGroup, client *ent.Client, logger *slog.Logger) {
	module := NewModule(client, logger, DefaultRouterConfig())
	module.Register(router)
}
