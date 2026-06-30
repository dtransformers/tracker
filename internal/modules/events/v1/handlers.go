package v1

import (
	"log/slog"
	"net/http"

	"tracking/internal/ent"
	"tracking/pkg/constants"
	custom_errors "tracking/pkg/errors"
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
// @Success 201 "Created"
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /api/v1/events/custom [post]
func (h *Handler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "Invalid request body", err))
		return
	}

	siteObj, exists := c.Get("site")
	if !exists {
		httputil.RespondError(c, custom_errors.NewUnauthorizedError(constants.Version, "Unauthorized"))
		return
	}

	site, ok := siteObj.(*ent.Site)
	if !ok || site == nil || !site.IsActive || site.IsDeleted {
		httputil.RespondError(c, custom_errors.NewUnauthorizedError(constants.Version, "Unauthorized"))
		return
	}

	_, err := h.service.CreateEvent(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create event", "error", err)
		httputil.RespondError(c, err)
		return
	}

	c.Status(http.StatusCreated)
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
// @Success 201 "Created"
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /api/v1/events/pageviews [post]
func (h *Handler) CreatePageView(c *gin.Context) {
	var req CreatePageViewRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "Invalid request body", err))
		return
	}

	siteObj, exists := c.Get("site")
	if !exists {
		httputil.RespondError(c, custom_errors.NewUnauthorizedError(constants.Version, "Unauthorized"))
		return
	}

	site, ok := siteObj.(*ent.Site)
	if !ok || site == nil || !site.IsActive || site.IsDeleted {
		httputil.RespondError(c, custom_errors.NewUnauthorizedError(constants.Version, "Unauthorized"))
		return
	}

	_, err := h.service.CreatePageView(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create page view", "error", err)
		httputil.RespondError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
