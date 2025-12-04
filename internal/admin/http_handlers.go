package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc}

	r.GET("/users", h.ListUsers)
	r.PUT("/user/role", h.UpdateUserRole)
}

// ListUsers godoc
// @Summary Admin: List all users
// @Tags Admin
// @Security BearerAuth
// @Success 200 {array} UserListItem
// @Router /api/admin/users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.svc.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// UpdateUserRole godoc
// @Summary Admin: Update user role
// @Tags Admin
// @Security BearerAuth
// @Param data body UpdateUserRoleRequest true "role update"
// @Success 200 {string} string "ok"
// @Router /api/admin/user/role [put]
func (h *Handler) UpdateUserRole(c *gin.Context) {
	var req UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	err := h.svc.UpdateUserRole(c.Request.Context(), req.UserID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
