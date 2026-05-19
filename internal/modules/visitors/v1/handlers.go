package v1

import (
	"log/slog"
	"net/http"

	"tracking/internal/ent"

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

// CreateSession godoc
// @Summary Create a visitor
// @Description Create a visitor on a site. Creates a new visitor if not exists, or updates last seen time. Requires anonymous_id cookie.
// @Tags visitors
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Track visitor request"
// @Param X-API-Key header string true "Public API Key"
// @Param X-Signature header string true "HMAC Signature"
// @Param Cookie header string true "Cookie must contain anonymous_id"
// @Success 200 {object} VisitorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/visitors [post]
func (h *Handler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
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

	res, err := h.service.CreateSession(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to create visitor session", "error", err)
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
