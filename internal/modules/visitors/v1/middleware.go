package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	sitesv1 "tracking/internal/modules/sites/v1"
)

func VerifyKeys(s *sitesv1.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		pubKey := c.GetString("publicKey")
		if pubKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing public key"})
			return
		}

		site, err := s.GetSiteKey(c.Request.Context(), pubKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid site key"})
			return
		}

		if !site.IsActive || site.IsDeleted {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Site is inactive or deleted"})
			return
		}

		c.Set("site", site)
		c.Next()
	}
}
