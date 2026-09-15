package events

import (
	"fmt"
	"log/slog"
	"tracking/internal/ent"
	v1 "tracking/internal/modules/events/v1"
	sitesv1 "tracking/internal/modules/sites/v1"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine, client *ent.Client, sessionMiddleware gin.HandlerFunc, logger *slog.Logger) *v1.Service {
	if logger == nil {
		logger = slog.Default()
	}
	v1Group := router.Group("/events/api/v1")

	sitesRepo := sitesv1.NewRepository(client)
	sitesService := sitesv1.NewService(sitesRepo, logger)

	module := v1.NewModule(client, sitesService, logger)
	module.SetSessionMiddleware(sessionMiddleware)
	module.Register(v1Group)
	fmt.Printf("Events module endpoints under /events/api/v1/events Are Ready To use \n")
	return module.Service()
}
