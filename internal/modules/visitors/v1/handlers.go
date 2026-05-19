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

// CreateVisitor godoc
// @Summary Create a visitor
// @Description Create a visitor on a site. Creates a new visitor if not exists, or updates last seen time.
// @Tags visitors
// @Accept json
// @Produce json
// @Param request body TrackVisitorRequest true "Track visitor request"
// @Param X-API-Key header string true "Public API Key"
// @Param X-Signature header string true "HMAC Signature"
// @Success 200 {object} VisitorResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/visitors [post]
func (h *Handler) CreateVisitor(c *gin.Context) {
	var req TrackVisitorRequest

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
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Site info not found in context"})
		return
	}

	// Use the ent.Site retrieved from the middleware context
	site, ok := siteObj.(*ent.Site)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid site info in context"})
		return
	}

	res, err := h.service.TrackVisitor(c.Request.Context(), req, site)
	if err != nil {
		h.logger.Error("failed to track visitor", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
