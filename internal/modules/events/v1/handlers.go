package v1

import (
	"log/slog"
	"net/http"

	"tracking/internal/ent"
	"tracking/pkg/httputil"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
	logger  *slog.Logger
}

func NewHandler(service *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// @Summary Create custom event
// @Description Records a custom event for a visitor
// @Tags Events
// @Accept json
// @Produce json
// @Param X-Api-Key header string true "Site API Key"
// @Param X-Fingerprint header string true "Client Fingerprint"
// @Param X-Nonce header string true "Request Nonce"
// @Param request body CreateEventRequest true "Event Request Payload"
// @Success 200 {object} ent.VisitorEvent
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /api/v1/events/custom [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest

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

	res, err := h.service.CreateEvent(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create event", "error", err)
		httputil.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// @Summary Create page view
// @Description Records a page view for a visitor
// @Tags Events
// @Accept json
// @Produce json
// @Param X-Api-Key header string true "Site API Key"
// @Param X-Fingerprint header string true "Client Fingerprint"
// @Param X-Nonce header string true "Request Nonce"
// @Param request body CreatePageViewRequest true "Page View Request Payload"
// @Success 200 {object} ent.PageView
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /api/v1/events/pageviews [post]
func (h *Handler) CreatePageView(c *gin.Context) {
	var req CreatePageViewRequest

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

	res, err := h.service.CreatePageView(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create page view", "error", err)
		httputil.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
