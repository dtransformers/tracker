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
	consents := router.Group("/consents")

	consents.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: middlewares.HMACSecretKeyFunc(m.siteSecretFunc()),
		NonceStore:   state.NewNonceStore(state.Redis),
	}))

	consents.Use(middlewares.VerifySite(m.siteSecretFunc()))
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
