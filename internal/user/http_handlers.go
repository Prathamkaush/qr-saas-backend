package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc: svc}

	r.GET("/me", h.Me)
}

// Me godoc
// @Summary Get logged-in user
// @Description Returns user profile for current JWT token
// @Tags User
// @Produce json
// @Success 200 {object} User
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/user/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	u, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, u)
}
