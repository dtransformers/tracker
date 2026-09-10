package v1

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"tracking/internal/ent"
	"tracking/pkg/constants"
	custom_errors "tracking/pkg/errors"
	"tracking/pkg/httputil"
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

// Create godoc
// @Summary Create a new site
// @Description Create a new site with auto-generated security keys and assigned tenant ID
// @Tags sites
// @Accept json
// @Produce json
// @Param site body CreateRequest true "Site creation request"
// @Security BasicAuth
// @Success 201 {object} Response
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 409 {object} httputil.ErrorResponse "Domain already exists"
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "Invalid request body", err))
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
// @Param tenant_id query int64 false "Filter by Tenant ID"
// @Security BasicAuth
// @Success 200 {object} ListResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites [get]
func (h *Handler) List(c *gin.Context) {
	var tenantID *int64
	if tidStr := c.Query("tenant_id"); tidStr != "" {
		tid, err := strconv.ParseInt(tidStr, 10, 64)
		if err != nil {
			h.logger.Warn("invalid tenant_id", "error", err)
			httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "Invalid tenant_id format", err))
			return
		}
		tenantID = &tid
	}

	sites, err := h.service.ListSites(c.Request.Context(), tenantID)
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
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites/{uuid} [get]
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
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Failure 409 {object} httputil.ErrorResponse "Domain already in use"
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites/{uuid} [patch]
func (h *Handler) Update(c *gin.Context) {
	uuid := c.Param("uuid")

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", "error", err)
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "Invalid request parameters", err))
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
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites/{uuid} [delete]
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
// @Failure 400 {object} httputil.ErrorResponse
// @Failure 404 {object} httputil.ErrorResponse
// @Failure 500 {object} httputil.ErrorResponse
// @Router /events/api/v1/sites/{uuid}/regenerate-keys [post]
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
	if err == nil {
		return
	}

	if ent.IsNotFound(err) || errors.Is(err, ErrSiteNotFound) {
		httputil.RespondError(c, custom_errors.NewNotFoundError(constants.Version, "Site not found", err))
		return
	}

	if errors.Is(err, ErrInvalidDomain) {
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, err.Error(), err))
		return
	}
	if errors.Is(err, ErrInvalidName) {
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, err.Error(), err))
		return
	}
	if errors.Is(err, ErrEmptyUpdate) {
		httputil.RespondError(c, custom_errors.NewBadRequestError(constants.Version, "At least one field must be updated", err))
		return
	}
	if errors.Is(err, ErrDomainExists) {
		httputil.RespondError(c, custom_errors.NewConflictError(constants.Version, "Domain already exists", err))
		return
	}
	if errors.Is(err, ErrKeyGeneration) {
		h.logger.Error("key generation failed", "error", err)
		httputil.RespondError(c, custom_errors.NewInternalError(constants.Version, "Failed to generate security keys", err))
		return
	}

	// Default fallback
	h.logger.Error("unexpected error", "error", err)
	httputil.RespondError(c, custom_errors.NewInternalError(constants.Version, "Internal server error", err))
}
