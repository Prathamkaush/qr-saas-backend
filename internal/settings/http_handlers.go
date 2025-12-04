package settings

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	r.GET("/", h.GetMySettings)
	r.PUT("/", h.UpdateMySettings)
}

// GetMySettings godoc
// @Summary Get settings for logged-in user
// @Tags Settings
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Settings
// @Router /api/settings/ [get]
func (h *Handler) GetMySettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	sett, err := h.svc.GetSettings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load settings"})
		return
	}

	c.JSON(http.StatusOK, sett)
}

// UpdateMySettings godoc
// @Summary Update settings for logged-in user
// @Tags Settings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body UpdateSettingsRequest true "Settings payload"
// @Success 200 {object} Settings
// @Router /api/settings/ [put]
func (h *Handler) UpdateMySettings(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	sett, err := h.svc.UpdateSettings(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, sett)
}
