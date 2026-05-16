package v1

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for sites.
type Handler struct {
	service *Service
	logger  *slog.Logger
}

// NewHandler creates a new site handler.
func NewHandler(service *Service, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}

	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Create godoc
// @Summary Create a new site
// @Description Create a new site with auto-generated security keys
// @Tags sites
// @Accept json
// @Produce json
// @Param site body CreateRequest true "Site creation request"
// @Security BasicAuth
// @Success 201 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Domain already exists"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	site, err := h.service.CreateSite(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, site)
}

// List godoc
// @Summary List all sites
// @Description Get a list of all sites
// @Tags sites
// @Produce json
// @Security BasicAuth
// @Success 200 {object} ListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites [get]
func (h *Handler) List(c *gin.Context) {
	sites, err := h.service.ListSites(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, sites)
}

// Get godoc
// @Summary Get a site by UUID
// @Description Get detailed information about a specific site
// @Tags sites
// @Produce json
// @Param uuid path string true "Site UUID"
// @Security BasicAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites/{uuid} [get]
func (h *Handler) Get(c *gin.Context) {
	uuid := c.Param("uuid")

	site, err := h.service.GetSite(c.Request.Context(), uuid)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, site)
}

// Update godoc
// @Summary Update a site
// @Description Update site information
// @Tags sites
// @Accept json
// @Produce json
// @Param uuid path string true "Site UUID"
// @Param site body UpdateRequest true "Site update request"
// @Security BasicAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Domain already in use"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites/{uuid} [patch]
func (h *Handler) Update(c *gin.Context) {
	uuid := c.Param("uuid")

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	site, err := h.service.UpdateSite(c.Request.Context(), uuid, req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, site)
}

// Delete godoc
// @Summary Delete a site
// @Description Permanently delete a site
// @Tags sites
// @Param uuid path string true "Site UUID"
// @Security BasicAuth
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites/{uuid} [delete]
func (h *Handler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")

	if err := h.service.DeleteSite(c.Request.Context(), uuid); err != nil {
		h.handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// RegenerateKeys godoc
// @Summary Regenerate site security keys
// @Description Generate new public and secret keys for a site
// @Tags sites
// @Produce json
// @Param uuid path string true "Site UUID"
// @Security BasicAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sites/{uuid}/regenerate-keys [post]
func (h *Handler) RegenerateKeys(c *gin.Context) {
	uuid := c.Param("uuid")

	site, err := h.service.RegenerateKeys(c.Request.Context(), uuid)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, site)
}

// handleError maps service errors to appropriate HTTP responses.
func (h *Handler) handleError(c *gin.Context, err error) {
	switch {
	case err == ErrSiteNotFound:
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Site not found",
		})
	case err == ErrDomainExists:
		c.JSON(http.StatusConflict, ErrorResponse{
			Error: "Domain already exists",
		})
	case err == ErrInvalidDomain:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid domain",
			Details: err.Error(),
		})
	case err == ErrInvalidName:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid name",
			Details: err.Error(),
		})
	case err == ErrEmptyUpdate:
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "At least one field must be updated",
		})
	case err == ErrKeyGeneration:
		h.logger.Error("key generation failed", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to generate security keys",
		})
	default:
		h.logger.Error("unexpected error", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
	}
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
