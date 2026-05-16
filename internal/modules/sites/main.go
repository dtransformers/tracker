package sites

import (
	"log/slog"
	"tracking/internal/ent"

	v1 "tracking/internal/modules/sites/v1"

	"github.com/gin-gonic/gin"
)

// RegisterV1 registers v1 site routes.
// This is the simple version for backward compatibility.
func RegisterV1(router *gin.Engine, client *ent.Client) {
	v1Group := router.Group("/api/v1")
	v1.Register(v1Group, client)
}

// RegisterV1WithLogger registers v1 site routes with a custom logger.
func RegisterV1WithLogger(router *gin.Engine, client *ent.Client, logger *slog.Logger) {
	v1Group := router.Group("/api/v1")
	v1.RegisterWithLogger(v1Group, client, logger)
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
