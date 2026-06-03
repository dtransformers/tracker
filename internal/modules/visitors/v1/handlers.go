package v1

import (
	"log/slog"
	"net/http"

	"tracking/internal/config"
	"tracking/internal/ent"
	"tracking/pkg/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
	appCfg  config.AppConfig
}

func NewHandler(service *Service, logger *slog.Logger, appCfg config.AppConfig) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
		appCfg:  appCfg,
	}
}

// CreateSession godoc
// @Summary Create a visitor
// @Description Create a visitor on a site. Creates a new visitor if not exists, or updates last seen time. Requires anonymous_id cookie.
// @Tags visitors
// @Accept json
// @Produce json
// @Param request body CreateVisitorRequest true "Track visitor request"
// @Param X-API-Key header string true "Public API Key"
// @Param X-Signature header string true "HMAC Signature"
// @Param X-Timestamp header string true "Request Timestamp"
// @Param X-Nonce header string true "Request Nonce"
// @Param Cookie header string true "Cookie must contain anonymous_id"
// @Success 200 {object} VisitorResponse
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /api/v1/visitors [post]
func (h *Handler) CreateVisitor(c *gin.Context) {
	var req CreateVisitorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, httputil.ErrorResponse{
			Error:   "Invalid request body",
			Details: httputil.FormatBindError(err),
		})
		return
	}

	siteObj, exists := c.Get("site")
	if !exists {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Error: "Unauthorized"})
		return
	}

	site, ok := siteObj.(*ent.Site)
	if !ok || site == nil || !site.IsActive || site.IsDeleted {
		c.JSON(http.StatusUnauthorized, httputil.ErrorResponse{Error: "Unauthorized"})
		return
	}
	response, err := h.service.CreateSession(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create visitor session", "error", err)
		httputil.RespondError(c, err)
		return
	}
	sessionExpiration := h.appCfg.SessionExpiration
	visitorExpiration := h.appCfg.VisitorExpiration
	c.SetCookie("session_id", response.SessionID, sessionExpiration, "/", "", true, true)
	c.SetCookie("visitor_uuid", response.UUID, visitorExpiration, "/", "", true, true)

	c.JSON(http.StatusOK, response)
}
