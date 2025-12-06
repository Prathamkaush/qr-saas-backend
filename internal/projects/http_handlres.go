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

	// ------------------------
	// MUST COME FIRST
	// ------------------------
	r.GET("/:id/qr", h.ListQRs)
	r.PUT("/:id/add/:qrID", h.AddQR)
	r.PUT("/:id/remove/:qrID", h.RemoveQR)

	// ------------------------
	// CRUD ROUTES (simple ones last)
	// ------------------------
	r.POST("/", h.Create)
	r.GET("/", h.List)
	r.GET("/:id", h.GetOne)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

//
// ------------------------
// CREATE PROJECT
// ------------------------
//

// Create godoc
// @Summary Create project (folder)
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

//
// ------------------------
// LIST PROJECTS
// ------------------------
//

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

//
// ------------------------
// GET ONE PROJECT
// ------------------------
//

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

//
// ------------------------
// UPDATE PROJECT
// ------------------------
//

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

//
// ------------------------
// DELETE PROJECT
// ------------------------
//

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

//
// ------------------------
// LIST QRs INSIDE PROJECT
// ------------------------
//

// @Summary List all QRs inside a project
// @Tags Projects
// @Security BearerAuth
// @Produce json
// @Success 200 {array} qr.QRCode
// @Router /api/projects/{id}/qr [get]
func (h *Handler) ListQRs(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	qrs, err := h.svc.ListProjectQRs(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to load qrs"})
		return
	}

	c.JSON(200, qrs)
}

//
// ------------------------
// MOVE QR INTO PROJECT
// ------------------------
//

// @Summary Add a QR to a project
// @Tags Projects
// @Security BearerAuth
// @Router /api/projects/{id}/add/{qrID} [put]
func (h *Handler) AddQR(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")
	qrID := c.Param("qrID")

	if err := h.svc.AssignQR(c.Request.Context(), userID, qrID, projectID); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}

//
// ------------------------
// REMOVE QR FROM PROJECT
// ------------------------
//

// @Summary Remove QR from project
// @Tags Projects
// @Security BearerAuth
// @Router /api/projects/{id}/remove/{qrID} [put]
func (h *Handler) RemoveQR(c *gin.Context) {
	userID := c.GetString("user_id")
	projectID := c.Param("id")
	qrID := c.Param("qrID")

	// optional: ensure project exists
	if _, err := h.svc.GetProject(c.Request.Context(), userID, projectID); err != nil {
		c.JSON(404, gin.H{"error": "project not found"})
		return
	}

	if err := h.svc.AssignQR(c.Request.Context(), userID, qrID, ""); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "ok"})
}
