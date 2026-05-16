package sites

import (
	"tracking/internal/ent"
	v1 "tracking/internal/modules/sites/v1"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine, client *ent.Client) {
	group := router.Group("/api/v1")
	v1.Register(group, client)
}

// Future versions:
// func RegisterV2(router *gin.Engine, client *ent.Client) { ... }
