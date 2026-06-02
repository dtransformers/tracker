package v1

import (
	"fmt"
	"log/slog"
	"net/http"

	"tracking/internal/ent"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
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
// @Success 200 {object} UpdateConsentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/consents [post]
func (h *Handler) UpdateConsent(c *gin.Context) {
	var req UpdateConsentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: formatBindError(err),
		})
		return
	}

	siteObj, exists := c.Get("site")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	site, ok := siteObj.(*ent.Site)
	if !ok || site == nil || !site.IsActive || site.IsDeleted {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	res, err := h.service.UpdateVisitorConsent(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to update consent", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	TraceId string `json:"trace_id"`
	Version string `json:"version"`
}

func formatBindError(err error) string {
	if err == nil {
		return "invalid request body"
	}

	if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
		return fmt.Sprintf("Invalid data type for field: %s", typeErr.Field)
	}

	if _, ok := err.(*json.SyntaxError); ok {
		return "Malformed JSON payload"
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		if len(validationErrs) > 0 {
			firstErr := validationErrs[0]
			return fmt.Sprintf("Validation failed on field: %s", firstErr.Field())
		}
	}

	return "The provided payload format is invalid or contains mismatched data types."
}
