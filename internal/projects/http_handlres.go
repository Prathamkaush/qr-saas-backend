package projects

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	r.POST("/", h.Create)
	r.GET("/", h.List)
	r.GET("/:id", h.GetOne)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary Create project (folder / campaign)
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body CreateProjectRequest true "Project data"
// @Success 201 {object} Project
// @Failure 400 {object} map[string]string
// @Router /api/projects/ [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	p, err := h.svc.CreateProject(c.Request.Context(), userID, req.Name, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// List godoc
// @Summary List projects for logged-in user
// @Tags Projects
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Project
// @Router /api/projects/ [get]
func (h *Handler) List(c *gin.Context) {
	userID := c.GetString("user_id")

	projects, err := h.svc.ListProjects(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list projects"})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetOne godoc
// @Summary Get single project
// @Tags Projects
// @Security BearerAuth
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} Project
// @Failure 404 {object} map[string]string
// @Router /api/projects/{id} [get]
func (h *Handler) GetOne(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	p, err := h.svc.GetProject(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// Update godoc
// @Summary Update project
// @Tags Projects
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param data body UpdateProjectRequest true "Update payload"
// @Success 200 {object} Project
// @Failure 404 {object} map[string]string
// @Router /api/projects/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	p, err := h.svc.UpdateProject(c.Request.Context(), userID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// Delete godoc
// @Summary Delete project
// @Tags Projects
// @Security BearerAuth
// @Param id path string true "Project ID"
// @Success 200 {string} string "ok"
// @Router /api/projects/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.DeleteProject(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
