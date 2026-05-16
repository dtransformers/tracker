package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Create handles site creation
// @Summary Create a new site
// @Description Create a new site with name and domain
// @Tags sites
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Site creation request"
// @Success 201 {object} ent.Site
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/sites [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.CreateSite(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Get handles fetching a single site
// @Summary Get a site
// @Description Get a site by its UUID
// @Tags sites
// @Produce json
// @Param uuid path string true "Site UUID"
// @Success 200 {object} ent.Site
// @Failure 404 {object} map[string]string
// @Router /api/v1/sites/{uuid} [get]
func (h *Handler) Get(c *gin.Context) {
	uuid := c.Param("uuid")
	res, err := h.service.GetSite(c.Request.Context(), uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// List handles listing all sites
// @Summary List all sites
// @Description Get a list of all active sites
// @Tags sites
// @Produce json
// @Success 200 {array} ent.Site
// @Failure 500 {object} map[string]string
// @Router /api/v1/sites [get]
func (h *Handler) List(c *gin.Context) {
	res, err := h.service.ListSites(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Update handles site updates
// @Summary Update a site
// @Description Update site name, domain, or status
// @Tags sites
// @Accept json
// @Produce json
// @Param uuid path string true "Site UUID"
// @Param request body UpdateRequest true "Site update request"
// @Success 200 {object} ent.Site
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/sites/{uuid} [patch]
func (h *Handler) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.service.UpdateSite(c.Request.Context(), uuid, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Delete handles soft-deleting a site
// @Summary Delete a site
// @Description Soft-delete a site by its UUID
// @Tags sites
// @Param uuid path string true "Site UUID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /api/v1/sites/{uuid} [delete]
func (h *Handler) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := h.service.DeleteSite(c.Request.Context(), uuid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
