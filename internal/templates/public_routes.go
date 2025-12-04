package templates

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(r *gin.Engine, svc Service) {
	r.GET("/t/:url_id", func(c *gin.Context) {
		urlID := c.Param("url_id")

		html, err := svc.RenderPublicPage(c.Request.Context(), urlID)
		if err != nil {
			c.String(http.StatusNotFound, "Template not found")
			return
		}

		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, html)
	})
}
