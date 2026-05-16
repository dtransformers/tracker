package v1

import (
	"tracking/internal/ent"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, client *ent.Client) {
	repo := NewRepository(client)
	svc := NewService(repo)
	handler := NewHandler(svc)

	sites := router.Group("/sites")
	{
		sites.POST("", handler.Create)
		sites.GET("", handler.List)
		sites.GET("/:uuid", handler.Get)
		sites.PATCH("/:uuid", handler.Update)
		sites.DELETE("/:uuid", handler.Delete)
	}
}
