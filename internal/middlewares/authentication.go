package middlewares

import (
	"context"
	"net/http"

	"tracking/pkg/constants"

	"github.com/gin-gonic/gin"
)

func Basic(username, password string) gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{
		username: password,
	})
}

type SessionValidator func(ctx context.Context, visitorUUID, sessionUUID string) (*SessionDetails, error)

// RequireVisitorAuthentication ensures that the request has valid Visitor and Session UUIDs via cookies.
// It also validates them against the database using the provided SessionValidator.
// This middleware must be used after RequestTracking middleware.
func RequireVisitorAuthentication(validator SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracking, exists := GetTrackingInfo(c.Request.Context())
		if !exists || tracking == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
			c.Abort()
			return
		}

		if tracking.VisitorUUID == "" || tracking.SessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request: missing visitor session cookies"})
			c.Abort()
			return
		}

		details, err := validator(c.Request.Context(), tracking.VisitorUUID, tracking.SessionToken)
		if err != nil || details == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request: invalid visitor session"})
			c.Abort()
			return
		}

		c.Set(SessionDetailsKey, details)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), SessionDetailsKey, details))

		c.Next()
	}
}

func GetSessionDetails(ctx context.Context) (*SessionDetails, bool) {
	if ctx == nil {
		return nil, false
	}
	if info := ctx.Value(SessionDetailsKey); info != nil {
		if details, ok := info.(*SessionDetails); ok {
			return details, true
		}
	}
	return nil, false
}

func SetTrackingCookies(c *gin.Context, sessionToken, visitorUUID string, sessionExp, visitorExp int) {
	c.SetCookieData(&http.Cookie{
		Name:     constants.CookieSessionToken,
		Value:    sessionToken,
		MaxAge:   sessionExp,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
	c.SetCookieData(&http.Cookie{
		Name:     constants.CookieVisitorUUID,
		Value:    visitorUUID,
		MaxAge:   visitorExp,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})
}
