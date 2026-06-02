package v1

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"tracking/internal/ent"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/visitors [post]
func (h *Handler) CreateVisitor(c *gin.Context) {
	var req CreateVisitorRequest

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
	response, err := h.service.CreateSession(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create visitor session", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	sessionExpiration := 30 * 60         // 30 Minutes In Sec
	visitorExpiration := 3600 * 24 * 365 // 1 Year In Sec
	c.SetCookie("session_id", response.SessionID, sessionExpiration, "/", "", true, true)
	c.SetCookie("visitor_uuid", response.UUID, visitorExpiration, "/", "", true, true)

	c.JSON(http.StatusOK, response)
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
