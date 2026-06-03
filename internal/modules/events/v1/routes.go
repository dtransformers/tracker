package v1

import (
	"context"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	sitesv1 "tracking/internal/modules/sites/v1"
	"tracking/internal/modules/state"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler      *Handler
	service      *Service
	repo         Repository
	sitesService     *sitesv1.Service
	sessionValidator middlewares.SessionValidator
	logger           *slog.Logger
}

func NewModule(client *ent.Client, sitesService *sitesv1.Service, logger *slog.Logger) *Module {
	if logger == nil {
		logger = slog.Default()
	}
	repo := NewRepository(client)
	svc := NewService(repo, logger)
	handler := NewHandler(svc, logger)

	return &Module{
		handler:      handler,
		service:      svc,
		repo:         repo,
		sitesService: sitesService,
		logger:       logger,
	}
}

func (m *Module) Service() *Service {
	return m.service
}

// siteSecretFunc returns a SiteSecretFunc backed by the sites service.
func (m *Module) siteSecretFunc() middlewares.SiteSecretFunc {
	return func(ctx context.Context, publicKey string) (string, *ent.Site, error) {
		site, err := m.sitesService.GetSiteKey(ctx, publicKey)
		if err != nil {
			return "", nil, err
		}
		return site.SecretKey, site, nil
	}
}

func (m *Module) Register(router *gin.RouterGroup) {
	events := router.Group("/events")

	events.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: middlewares.HMACSecretKeyFunc(m.siteSecretFunc()),
		NonceStore:   state.NewNonceStore(state.Redis),
	}))

	events.Use(middlewares.VerifySite(m.siteSecretFunc()))
	events.Use(middlewares.RateLimit(1*time.Minute, 600))
	events.Use(middlewares.Validation())
	events.Use(middlewares.Cors())
	events.Use(middlewares.Metrics(m.logger))
	events.Use(middlewares.RequestTracking())
	events.Use(func(c *gin.Context) {
		if m.sessionValidator == nil {
			c.JSON(500, gin.H{"error": "internal error: session validator not configured"})
			c.Abort()
			return
		}
		middlewares.RequireVisitorAuthentication(m.sessionValidator)(c)
	})

	events.POST("/custom", m.handler.CreateEvent)
	events.POST("/pageviews", m.handler.CreatePageView)
}

func (m *Module) SetSessionValidator(validator middlewares.SessionValidator) {
	m.sessionValidator = validator
}
