package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method
		reqID, _ := c.Get("request_id")

		log.Printf("[%v] %s %s %d %s", reqID, method, path, status, latency)
	}
}
