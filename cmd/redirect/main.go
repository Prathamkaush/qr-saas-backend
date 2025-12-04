package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"qr-saas/internal/analytics"
	"qr-saas/internal/config"
	"qr-saas/internal/db"
	"qr-saas/internal/qr"
	"qr-saas/internal/redirect"
)

func main() {
	cfg := config.Load()

	// ClickHouse for analytics (OK)
	chConn, err := db.NewClickHouse(cfg)
	if err != nil {
		log.Fatal("❌ ClickHouse:", err)
	}

	// Postgres for QR metadata (IMPORTANT)
	pgDB := db.NewPostgresPool(cfg)

	// Repositories
	qrRepo := qr.NewRepository(pgDB)
	analyticsRepo := analytics.NewRepository(chConn)

	// Services
	analyticsSvc := analytics.NewService(analyticsRepo)
	redirectSvc := redirect.NewService(qrRepo, analyticsSvc)

	// Router
	r := gin.Default()
	redirect.RegisterRoutes(r, redirectSvc)

	log.Println("🚀 Redirect service running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
