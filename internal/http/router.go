package http

import (
	"os" // Needed for reading environment variables
	"time"

	"github.com/gin-contrib/cors"package http

import (
    "os"
    "time"

    "github.com/gin-contrib/cors" 
    "github.com/gin-gonic/gin"
    "github.com/redis/go-redis/v9"

    "qr-saas/internal/config" // <-- Need this import
    "qr-saas/internal/http/middleware"
)

// 🔥 FIX 1: Update the function signature to accept cfg
func NewRouter(redis *redis.Client, cfg config.Config) *gin.Engine {
    r := gin.New()

    // ----------------------------------------------
    // FIX: Read FRONTEND_URL from environment for CORS
    // ----------------------------------------------
    frontendURL := os.Getenv("FRONTEND_URL")
    allowedOrigins := []string{"http://localhost:3000"} 

    // Add the deployed Vercel/Render URL if it exists
    if frontendURL != "" {
        allowedOrigins = append(allowedOrigins, frontendURL)
    }

    r.Use(cors.New(cors.Config{
        AllowOrigins:     allowedOrigins, 
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
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
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	// Assuming this exists
)

func NewRouter(redis *redis.Client) *gin.Engine {
	r := gin.New()

	// ----------------------------------------------
	// 🔥 FIX: Read FRONTEND_URL from environment for CORS
	// ----------------------------------------------
	frontendURL := os.Getenv("FRONTEND_URL")
	allowedOrigins := []string{"http://localhost:3000"} // Always allow local dev

	// Add the deployed Vercel/Render URL if it exists
	if frontendURL != "" {
		allowedOrigins = append(allowedOrigins, frontendURL)
	}

	r.Use(cors.New(cors.Config{
		// Allow the dynamic list of origins
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ... rest of r.Use statements ...

	return r
}
