package v1

import (
	"context"
	"log/slog"
	"time"
	"tracking/internal/config"
	"tracking/internal/ent"
	"tracking/internal/middlewares"
	consentsv1 "tracking/internal/modules/consents/v1"
	devicesv1 "tracking/internal/modules/devices/v1"
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

func NewModule(client *ent.Client, sitesService *sitesv1.Service, sessionSvc *visitor_sessionsv1.Service, consentsSvc *consentsv1.Service, logger *slog.Logger, appCfg config.AppConfig) *Module {
	if logger == nil {
		logger = slog.Default()
	}
	repo := NewRepository(client)
	deviceRepo := devicesv1.NewRepository(client)
	svc := NewService(repo, sitesService, deviceRepo, sessionSvc, consentsSvc, logger)
	handler := NewHandler(svc, logger, appCfg)

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
// Using the request context ensures timeouts and cancellations propagate correctly.
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
	visitors := router.Group("/visitors")

	visitors.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: middlewares.HMACSecretKeyFunc(m.siteSecretFunc()),
		NonceStore:   state.NewNonceStore(state.Redis),
	}))

	visitors.Use(middlewares.VerifySite(m.siteSecretFunc()))
	visitors.Use(middlewares.RateLimit(1*time.Minute, 600))
	visitors.Use(middlewares.Validation())
	visitors.Use(middlewares.Cors())
	visitors.Use(middlewares.Metrics(m.logger))
	visitors.Use(middlewares.RequestTracking())

	visitors.POST("/", m.handler.CreateVisitor)
}
