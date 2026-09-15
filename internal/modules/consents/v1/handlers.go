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

// @Summary Update visitor consent
// @Description Updates the consent status for a specific site and visitor
// @Tags Consents
// @Accept json
// @Produce json
// @Param X-Api-Key header string true "Site API Key"
// @Param X-Fingerprint header string true "Client Fingerprint"
// @Param X-Nonce header string true "Request Nonce"
// @Param request body UpdateConsentRequest true "Consent Request Payload"
// @Success 201 "Created"
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/consents [post]
func (h *Handler) UpdateConsent(c *gin.Context) {
	var req UpdateConsentRequest

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

	_, err := h.service.UpdateVisitorConsent(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to update consent", "error", err)
		httputil.RespondError(c, err)
		return
	}

	c.Status(http.StatusCreated)
}
