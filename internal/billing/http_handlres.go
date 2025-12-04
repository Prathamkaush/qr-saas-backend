package billing

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	h := &Handler{svc}

	r.GET("/plans", h.GetPlans)
	r.POST("/subscribe", h.Subscribe)
	r.GET("/subscription", h.GetSubscription)
}

// GetPlans godoc
// @Summary List billing plans
// @Tags Billing
// @Produce json
// @Success 200 {array} Plan
// @Router /api/billing/plans [get]
func (h *Handler) GetPlans(c *gin.Context) {
	plans, err := h.svc.GetPlans(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

// Subscribe godoc
// @Summary Subscribe to a plan
// @Tags Billing
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param data body SubscribeRequest true "subscription"
// @Success 200 {string} string "ok"
// @Router /api/billing/subscribe [post]
func (h *Handler) Subscribe(c *gin.Context) {
	userID := c.GetString("user_id")

	var req SubscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}

	err := h.svc.Subscribe(c.Request.Context(), userID, req.PlanID, req.StripeToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// GetSubscription godoc
// @Summary Get active subscription
// @Tags Billing
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Subscription
// @Router /api/billing/subscription [get]
func (h *Handler) GetSubscription(c *gin.Context) {
	userID := c.GetString("user_id")

	sub, err := h.svc.GetActiveSubscription(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}
	c.JSON(http.StatusOK, sub)
}
