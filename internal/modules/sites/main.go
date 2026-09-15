package sites

import (
	"fmt"
	"tracking/internal/ent"

	v1 "tracking/internal/modules/sites/v1"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine, client *ent.Client, rootUser, rootPass string) {
	v1Group := router.Group("/events/api/v1")
	v1.RegisterWithCredentials(v1Group, client, rootUser, rootPass)
	fmt.Printf("Sites module endpoints under /events/api/v1/sites Are Ready To use \n")

}
