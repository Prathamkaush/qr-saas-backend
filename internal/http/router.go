package http

import (
	"os"
	"time"

	"github.com/gin-contrib/cors" 
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"qr-saas/internal/config" 
	"qr-saas/internal/http/middleware"
)

// NewRouter accepts both Redis client and the global config struct.
func NewRouter(redis *redis.Client, cfg config.Config) *gin.Engine {
	r := gin.New()

	// ----------------------------------------------
	// FIX: Dynamically configure CORS for Render/Vercel
	// ----------------------------------------------
	
	// Check the FRONTEND_URL environment variable
	frontendURL := os.Getenv("FRONTEND_URL")
	
	// Allow multiple domains for Vercel/Render and local dev
	allowedOrigins := []string{"http://localhost:3000"} 

	if frontendURL != "" {
		// Add the deployed HTTPS domain to the allowed list
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins, 
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())

	// Keep your Rate Limiter
	r.Use(middleware.RateLimit(redis, 200, time.Minute))

	return r
}
