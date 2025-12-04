package qrtypes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, svc Service) {
	r.POST("/:id/wifi", func(c *gin.Context) { /* ... */ })
	r.POST("/:id/vcard", func(c *gin.Context) { /* ... */ })
	r.GET("/:id", func(c *gin.Context) { /* ... */ })
}
