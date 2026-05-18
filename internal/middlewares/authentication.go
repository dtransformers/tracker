package middlewares

import "github.com/gin-gonic/gin"

func Basic(username, password string) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}
