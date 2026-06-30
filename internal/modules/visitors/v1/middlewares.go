package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
)

// AutoResolveSessionMiddleware ensures that a request has a valid visitor and session UUID.
// If the session is missing or expired, it will automatically try to create a new one using the tracking info.
// This prevents valid pageviews and events from being blocked just because their session expired.
func (s *Service) AutoResolveSessionMiddleware(sessionExp, visitorExp int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tracking, exists := middlewares.GetTrackingInfo(c.Request.Context())
		if !exists || tracking == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request"})
			c.Abort()
			return
		}

		if s.tryValidateExistingSession(c, tracking, sessionExp, visitorExp) {
			c.Next()
			return
		}

		site, valid := getSiteFromContext(c)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request: missing or invalid site"})
			c.Abort()
			return
		}

		if !s.autoResolveNewSession(c, tracking, site, sessionExp, visitorExp) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized request: failed to auto-resolve session"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (s *Service) tryValidateExistingSession(c *gin.Context, tracking *middlewares.TrackingInfo, sessionExp, visitorExp int) bool {
	if tracking.VisitorUUID == "" || tracking.SessionToken == "" {
		return false
	}
	details, err := s.sessionSvc.ValidateSessionDetails(c.Request.Context(), tracking.VisitorUUID, tracking.SessionToken)
	if err != nil || details == nil {
		return false
	}

	middlewares.SetTrackingCookies(c, tracking.SessionToken, tracking.VisitorUUID, sessionExp, visitorExp, s.appCfg.CookiesDomain)
	c.Set(middlewares.SessionDetailsKey, details)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), middlewares.SessionDetailsKey, details))
	return true
}

func getSiteFromContext(c *gin.Context) (*ent.Site, bool) {
	siteObj, exists := c.Get("site")
	if !exists {
		return nil, false
	}
	site, ok := siteObj.(*ent.Site)
	if !ok || site == nil || !site.IsActive || site.IsDeleted {
		return nil, false
	}
	return site, true
}

func (s *Service) autoResolveNewSession(c *gin.Context, tracking *middlewares.TrackingInfo, site *ent.Site, sessionExp, visitorExp int) bool {
	req := CreateVisitorRequest{
		VisitorUUID:     tracking.VisitorUUID,
		SessionToken:    tracking.SessionToken,
		FingerprintHash: c.GetHeader("X-Fingerprint"),
	}

	res, err := s.CreateSession(c.Request.Context(), req, site)
	if err != nil {
		return false
	}

	details, err := s.sessionSvc.ValidateSessionDetails(c.Request.Context(), res.UUID, res.SessionToken)
	if err != nil || details == nil {
		return false
	}

	tracking.VisitorUUID = res.UUID
	tracking.SessionToken = res.SessionToken

	middlewares.SetTrackingCookies(c, res.SessionToken, res.UUID, sessionExp, visitorExp, s.appCfg.CookiesDomain)
	c.Set(middlewares.SessionDetailsKey, details)
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), middlewares.SessionDetailsKey, details))
	return true
}
