package redirect

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, svc *Service) {
	h := &Handler{svc: svc}
	r.GET("/r/:code", h.RedirectQR)
}

type Handler struct {
	svc *Service
}

// RedirectQR godoc
// @Summary Redirect QR Scan
// @Description Redirect user to target URL and log scan analytics
// @Tags Redirect
// @Param code path string true "QR Code Short ID"
// @Success 302 {string} string "redirect"
// @Failure 404 {object} map[string]string
// @Router /r/{code} [get]
func (h *Handler) RedirectQR(c *gin.Context) {
	code := c.Param("code")
	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	ref := c.Request.Referer()

	target, err := h.svc.ResolveAndLog(c.Request.Context(), code, ip, ua, ref)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "QR code not found"})
		return
	}

	c.Redirect(http.StatusFound, target)
}
