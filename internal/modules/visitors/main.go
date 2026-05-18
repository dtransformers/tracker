package visitors

import (
	"log/slog"
	"tracking/internal/ent"
	sitesv1 "tracking/internal/modules/sites/v1"
	v1 "tracking/internal/modules/visitors/v1"

	"github.com/gin-gonic/gin"
)

// RegisterV1 registers v1 visitor routes.
func RegisterV1(router *gin.Engine, client *ent.Client, logger *slog.Logger) {
	v1Group := router.Group("/api/v1")
	
	// Create sites service to pass to visitors
	sitesRepo := sitesv1.NewRepository(client)
	sitesService := sitesv1.NewService(sitesRepo, logger)
	
	module := v1.NewModule(client, sitesService, logger)
	module.Register(v1Group)
}
