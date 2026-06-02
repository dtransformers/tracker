package v1

import (
	"context"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	sitesv1 "tracking/internal/modules/sites/v1"
	"tracking/internal/modules/state"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *Handler
	service      *Service
	repo         Repository
	sitesService *sitesv1.Service
	logger       *slog.Logger
}

func NewModule(client *ent.Client, sessionSvc *visitor_sessionsv1.Service, sitesService *sitesv1.Service, logger *slog.Logger) *Module {
	if logger == nil {
		logger = slog.Default()
	}
	repo := NewRepository(client)
	svc := NewService(repo, sessionSvc, logger)
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
	consents := router.Group("/consents")
	consents.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: func(publicKey string) (string, error) {
			site, err := m.sitesService.GetSiteKey(context.Background(), publicKey)
			if err != nil {
				return "", err
			}
			return site.SecretKey, nil
		},
		NonceStore: state.NewNonceStore(state.Redis),
	}))

	// Reusing verify keys from visitors if exported, or duplicating here.
	// We'll assume the middlewares are shared or we use standard middleware.
	// Since VerifyKeys was in visitors, we need to extract it or move it to middlewares.
	// Actually we can just copy VerifyKeys logic or create it in middlewares later.
	// For now, I will omit VerifyKeys if it's not exported, or just use it from visitors.
	// Let's use it from tracking/internal/modules/visitors/v1 (aliased) or create it here.
	// It's cleaner to duplicate or use middlewares. For now:
	consents.Use(func(c *gin.Context) {
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

	consents.Use(middlewares.RateLimit(1*time.Minute, 600))
	consents.Use(middlewares.Validation())
	consents.Use(middlewares.Cors())
	consents.Use(middlewares.Metrics(m.logger))
	consents.Use(middlewares.RequestTracking())

	consents.POST("/", m.handler.UpdateConsent)
}

func (m *Module) Service() *Service {
	return m.service
}
