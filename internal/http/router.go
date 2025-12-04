package http

import (
	"time"

	"github.com/gin-contrib/cors" // <--- Import the official library
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"qr-saas/internal/config"
	"qr-saas/internal/http/middleware"
)

func NewRouter(redis *redis.Client, cfg config.Config) *gin.Engine {
	r := gin.New()

	// ----------------------------------------------
	// FIX: Use Official CORS Config
	// ----------------------------------------------
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // <--- Specific Origin (No *)
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // <--- This works now because Origin is not *
		MaxAge:           12 * time.Hour,
	}))

	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())

	// Keep your Rate Limiter
	r.Use(middleware.RateLimit(redis, 200, time.Minute))

	return r
}
