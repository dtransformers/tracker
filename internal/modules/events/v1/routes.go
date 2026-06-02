package v1

import (
	"context"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	sitesv1 "tracking/internal/modules/sites/v1"
	"tracking/internal/modules/state"
	visitorsv1 "tracking/internal/modules/visitors/v1"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *Handler
	service      *Service
	repo         Repository
	sitesService *sitesv1.Service
	logger       *slog.Logger
}

func NewModule(client *ent.Client, visitorsSvc *visitorsv1.Service, sitesService *sitesv1.Service, logger *slog.Logger) *Module {
	if logger == nil {
		logger = slog.Default()
	}
	repo := NewRepository(client)
	svc := NewService(repo, visitorsSvc, logger)
	handler := NewHandler(svc, logger)

	return &Module{
		handler:      handler,
		service:      svc,
		repo:         repo,
		sitesService: sitesService,
		logger:       logger,
	}
}

func (m *Module) Register(router *gin.RouterGroup) {
	events := router.Group("/events")
	events.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: func(publicKey string) (string, error) {
			site, err := m.sitesService.GetSiteKey(context.Background(), publicKey)
			if err != nil {
				return "", err
			}
			return site.SecretKey, nil
		},
		NonceStore: state.NewNonceStore(state.Redis),
	}))

	events.Use(func(c *gin.Context) {
		siteKey := c.GetHeader("X-Api-Key")
		if siteKey == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		site, err := m.sitesService.GetSiteKey(c.Request.Context(), siteKey)
		if err != nil || site == nil {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("site", site)
		c.Next()
	})

	events.Use(middlewares.RateLimit(1*time.Minute, 600))
	events.Use(middlewares.Validation())
	events.Use(middlewares.Cors())
	events.Use(middlewares.Metrics(m.logger))
	events.Use(middlewares.RequestTracking())

	events.POST("/custom", m.handler.CreateEvent)
	events.POST("/pageviews", m.handler.CreatePageView)
}

func (m *Module) Service() *Service {
	return m.service
}
