package audit

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc}

	r.GET("/my", h.GetMyEvents)
}

// GetMyEvents godoc
// @Summary Get audit events for logged-in user
// @Tags Audit
// @Security BearerAuth
// @Produce json
// @Success 200 {array} AuditEvent
// @Router /api/audit/my [get]
func (h *Handler) GetMyEvents(c *gin.Context) {
	userID := c.GetString("user_id")

	limitStr := c.DefaultQuery("limit", "20")
	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 {
		limit = 20
	}

	events, err := h.svc.GetUserEvents(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, events)
}
