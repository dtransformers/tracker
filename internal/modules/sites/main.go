package sites

import (
	"log/slog"
	"tracking/internal/ent"

	v1 "tracking/internal/modules/sites/v1"

	"github.com/gin-gonic/gin"
)

// RegisterV1 registers v1 site routes.
func RegisterV1(router *gin.Engine, client *ent.Client, rootUser, rootPass string) {
	v1Group := router.Group("/api/v1")
	v1.RegisterWithCredentials(v1Group, client, rootUser, rootPass)
}

// RegisterV1WithLogger registers v1 site routes with a custom logger and credentials.
func RegisterV1WithLogger(router *gin.Engine, client *ent.Client, logger *slog.Logger, rootUser, rootPass string) {
	v1Group := router.Group("/api/v1")
	v1.RegisterWithLoggerAndCredentials(v1Group, client, logger, rootUser, rootPass)
}

func RegisterV1WithConfig(
	router *gin.Engine,
	client *ent.Client,
	logger *slog.Logger,
	cfg v1.RouterConfig,
) *v1.Module {
	v1Group := router.Group("/api/v1")

	module := v1.NewModule(client, logger, cfg)
	module.Register(v1Group)

	return module
}
