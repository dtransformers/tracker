package v1

import (
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	sitesv1 "tracking/internal/modules/sites/v1"
	"tracking/internal/modules/state"

	"github.com/gin-gonic/gin"
)

type Module struct {
	handler          *Handler
	service          *Service
	repo             Repository
	sitesService      *sitesv1.Service
	sessionMiddleware gin.HandlerFunc
	logger            *slog.Logger
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

func (m *Module) Register(router *gin.RouterGroup) {
	consents := router.Group("/consents")

	consents.Use(middlewares.HMAC(middlewares.HMACConfig{
		GetSecretKey: middlewares.HMACSecretKeyFunc(middlewares.DefaultSiteSecretFunc(m.sitesService)),
		NonceStore:   state.NewNonceStore(state.Redis),
	}))

	consents.Use(middlewares.VerifySite(middlewares.DefaultSiteSecretFunc(m.sitesService)))
	consents.Use(middlewares.RateLimit(1*time.Minute, 600))
	consents.Use(middlewares.Validation())
	consents.Use(middlewares.Cors())
	consents.Use(middlewares.Metrics(m.logger))
	consents.Use(middlewares.RequestTracking())
	consents.Use(func(c *gin.Context) {
		if m.sessionMiddleware == nil {
			c.JSON(500, gin.H{"error": "internal error: session middleware not configured"})
			c.Abort()
			return
		}
		m.sessionMiddleware(c)
	})

	consents.POST("/", m.handler.UpdateConsent)
}

func (m *Module) Service() *Service {
	return m.service
}

func (m *Module) SetSessionMiddleware(middleware gin.HandlerFunc) {
	m.sessionMiddleware = middleware
}
