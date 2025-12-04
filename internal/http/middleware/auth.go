// internal/http/middleware/auth.go
package middleware

import (
	"net/http"
	"strings"

	"qr-saas/internal/auth"

	"github.com/gin-gonic/gin"
)

func JWTAuth(authSvc auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		// ⭐ FIX: Let OPTIONS go through for CORS preflight
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		_, userID, err := authSvc.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
