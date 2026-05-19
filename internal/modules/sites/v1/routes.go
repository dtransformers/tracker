package v1

import (
	"log/slog"
	"tracking/internal/ent"
	"tracking/internal/middlewares"

	"github.com/gin-gonic/gin"
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

	group.Use(middlewares.Metrics(m.logger))

	if m.config.EnableRateLimit {
		group.Use(middlewares.RateLimit(m.config.RateLimitWindow, m.config.RateLimitRequests))
	}

	group.Use(middlewares.Basic(m.config.RootUsername, m.config.RootPassword))
	group.Use(middlewares.Validation())
}

func (m *Module) registerRoutes(group *gin.RouterGroup) {
	group.POST("", m.handler.Create)
	group.GET("", m.handler.List)

	group.GET("/:uuid", m.handler.Get)
	group.PATCH("/:uuid", m.handler.Update)
	group.DELETE("/:uuid", m.handler.Delete)

	group.POST("/:uuid/regenerate-keys", m.handler.RegenerateKeys)
}

func (m *Module) Health() map[string]interface{} {
	return map[string]interface{}{
		"status":  "healthy",
		"service": "sites-v1",
	}
}

func (m *Module) Shutdown() error {
	m.logger.Info("shutting down sites module")
	return nil
}

func Register(router *gin.RouterGroup, client *ent.Client) {
	RegisterWithLogger(router, client, slog.Default())
}

func RegisterWithLogger(router *gin.RouterGroup, client *ent.Client, logger *slog.Logger) {
	module := NewModule(client, logger, DefaultRouterConfig())
	module.Register(router)
}

func RegisterWithCredentials(router *gin.RouterGroup, client *ent.Client, rootUser, rootPass string) {
	RegisterWithLoggerAndCredentials(router, client, slog.Default(), rootUser, rootPass)
}

func RegisterWithLoggerAndCredentials(router *gin.RouterGroup, client *ent.Client, logger *slog.Logger, rootUser, rootPass string) {
	cfg := DefaultRouterConfig()
	cfg.RootUsername = rootUser
	cfg.RootPassword = rootPass
	module := NewModule(client, logger, cfg)
	module.Register(router)
}
