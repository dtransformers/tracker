package v1

import (
	"context"
	"log/slog"
	"time"
	"tracking/internal/ent"
	"tracking/internal/middlewares"
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

func NewModule(client *ent.Client, sitesService *sitesv1.Service, logger *slog.Logger) *Module {
	repo := NewRepository(client)
	deviceRepo := devicesv1.NewRepository(client)
	sessionRepo := visitor_sessionsv1.NewRepository(client)
	svc := NewService(repo, sitesService, deviceRepo, sessionRepo, logger)
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
	visitors := router.Group("/visitors")
	visitors.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: func(publicKey string) (string, error) {
			site, err := m.sitesService.GetSiteKey(context.Background(), publicKey)
			if err != nil {
				return "", err
			}
			return site.SecretKey, nil
		},
		NonceStore: state.NewNonceStore(state.Redis),
	}))

	visitors.Use(VerifyKeys(m.sitesService))
	visitors.Use(middlewares.RateLimit(1*time.Minute, 600))
	visitors.Use(middlewares.Validation())
	visitors.Use(middlewares.Cors())
	visitors.Use(middlewares.Metrics(m.logger))
	visitors.Use(middlewares.RequestTracking())

	visitors.POST("/", m.handler.CreateVisitor)
}
