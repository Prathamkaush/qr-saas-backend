package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc}

	r.GET("/global", h.ListGlobal)
	r.GET("/mine", h.ListMine)
	r.POST("/", h.Create)
	r.GET("/:id", h.GetOne)
	r.PUT("/:id", h.Update)
	r.DELETE("/:id", h.Delete)
}

// @Summary List global templates
// @Tags Templates
// @Produce json
// @Success 200 {array} Template
// @Router /api/templates/global [get]
func (h *Handler) ListGlobal(c *gin.Context) {
	out, err := h.svc.ListGlobal(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// @Summary List user templates
// @Tags Templates
// @Security BearerAuth
// @Produce json
// @Success 200 {array} Template
// @Router /api/templates/mine [get]
func (h *Handler) ListMine(c *gin.Context) {
	userID := c.GetString("user_id")
	out, err := h.svc.ListMine(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, out)
}

// @Summary Create a template
// @Tags Templates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body CreateTemplateRequest true "template"
// @Success 201 {object} Template
// @Router /api/templates/ [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")
	var req CreateTemplateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	tpl, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create failed"})
		return
	}

	c.JSON(http.StatusCreated, tpl)
}

// @Summary Get template
// @Tags Templates
// @Produce json
// @Security BearerAuth
// @Param id path string true "Template ID"
// @Success 200 {object} Template
// @Router /api/templates/{id} [get]
func (h *Handler) GetOne(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	t, err := h.svc.Get(c.Request.Context(), userID, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if t == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// @Summary Update template
// @Tags Templates
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Template ID"
// @Param data body UpdateTemplateRequest true "update"
// @Success 200 {object} Template
// @Router /api/templates/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	t, err := h.svc.Update(c.Request.Context(), userID, id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	if t == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, t)
}

// @Summary Delete template
// @Tags Templates
// @Security BearerAuth
// @Param id path string true "Template ID"
// @Success 200 {string} string "ok"
// @Router /api/templates/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	if err := h.svc.Delete(c.Request.Context(), userID, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
