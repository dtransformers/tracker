package v1

import (
	"context"
	"log/slog"
	"tracking/internal/ent"
	"tracking/internal/middlewares"
	sitesv1 "tracking/internal/modules/sites/v1"

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
	svc := NewService(repo, sitesService, logger)
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

	// Apply HMAC middleware
	visitors.Use(middlewares.HMAC(func(publicKey string) (string, error) {
		site, err := m.sitesService.GetSiteKey(context.Background(), publicKey)
		if err != nil {
			return "", err
		}
		// Assuming site struct has SecretKey field accessible
		// We need to make sure it's exported or available.
		// In ent, it generates a field in the struct.
		return site.SecretKey, nil
	}))

	visitors.POST("/track", m.handler.Track)
}
