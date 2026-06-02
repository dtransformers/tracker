package visitors

import (
	"fmt"
	"log/slog"
	"tracking/internal/config"
	"tracking/internal/ent"
	consentsv1 "tracking/internal/modules/consents/v1"
	sitesv1 "tracking/internal/modules/sites/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
	v1 "tracking/internal/modules/visitors/v1"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine, client *ent.Client, sessionSvc *visitor_sessionsv1.Service, consentsSvc *consentsv1.Service, logger *slog.Logger, appCfg config.AppConfig) *v1.Service {
	if logger == nil {
		logger = slog.Default()
	}
	v1Group := router.Group("/api/v1")

	sitesRepo := sitesv1.NewRepository(client)
	sitesService := sitesv1.NewService(sitesRepo, logger)

	module := v1.NewModule(client, sitesService, sessionSvc, consentsSvc, logger, appCfg)
	module.Register(v1Group)
	fmt.Printf("Visitors module endpoints under /api/v1/visitor Are Ready To use \n")
	return module.Service()
}
