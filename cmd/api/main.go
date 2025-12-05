// @title QR SaaS API
// @version 2.0
// @description Scalable QR Code SaaS backend with dynamic QR, analytics, billing, templates, folders, OAuth, and more.
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"log"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "qr-saas/docs" // <-- REQUIRED FOR SWAGGER !!!

	"qr-saas/internal/admin"
	"qr-saas/internal/analytics"
	"qr-saas/internal/audit"
	"qr-saas/internal/auth"
	"qr-saas/internal/billing"
	"qr-saas/internal/config"
	"qr-saas/internal/db"
	internalhttp "qr-saas/internal/http"
	"qr-saas/internal/http/middleware"
	"qr-saas/internal/projects"
	"qr-saas/internal/qr"
	"qr-saas/internal/redirect"
	"qr-saas/internal/settings"
	"qr-saas/internal/templates"
)

func main() {
	cfg := config.Load()

	fmt.Println("FRONTEND_URL READ:", os.Getenv("FRONTEND_URL")) // Log the raw ENV var
    fmt.Println("FRONTEND_URL USED:", cfg.FrontendURL) // Log the value passed to the service

	
	// --------------------------
	// DATABASE CONNECTIONS
	// --------------------------

	// ClickHouse (analytics)

	// PostgreSQL (core DB)
	pgDB := db.NewPostgresPool(cfg)

	// Redis (rate limit, cache)
	redisClient := db.NewRedis(cfg.RedisURL)

	// Router
	r := internalhttp.NewRouter(redisClient, cfg)

	// --------------------------
	// SERVICES + REPOSITORIES
	// --------------------------

	// Auth
	authRepo := auth.NewRepository(pgDB)
	googleOAuth := auth.NewGoogleOAuth(cfg)
	authSvc := auth.NewService(authRepo, googleOAuth, cfg.JWTSecret)

	// QR
	qrRepo := qr.NewRepository(pgDB)
	qrSvc := qr.NewService(qrRepo)

	// Analytics
	analyticsRepo := analytics.NewRepository(pgDB)
	analyticsSvc := analytics.NewService(analyticsRepo)

	// Redirect
	redirectSvc := redirect.NewService(qrRepo, analyticsSvc)

	// Projects
	projectsRepo := projects.NewRepository(pgDB)
	projectsSvc := projects.NewService(projectsRepo)

	// Settings
	settingsRepo := settings.NewRepository(pgDB)
	settingsSvc := settings.NewService(settingsRepo)

	// Templates
	templatesRepo := templates.NewRepository(pgDB)
	templatesSvc := templates.NewService(templatesRepo)

	// Billing
	billingRepo := billing.NewRepository(pgDB)
	billingSvc := billing.NewService(billingRepo)

	// Admin
	adminRepo := admin.NewRepository(pgDB)
	adminSvc := admin.NewService(adminRepo)

	// Audit
	auditRepo := audit.NewRepository(pgDB)
	auditSvc := audit.NewService(auditRepo)

	fmt.Println("JWT Secret Loaded:", cfg.JWTSecret)

	// --------------------------
	// ROUTES
	// --------------------------

	// AUTH
	authGroup := r.Group("/api/auth")
	auth.RegisterRoutes(authGroup, authSvc)

	// QR
	apiQR := r.Group("/api/qr")
	apiQR.Use(middleware.JWTAuth(authSvc))
	qr.RegisterRoutes(apiQR, qrSvc)

	// ANALYTICS
	apiAnalytics := r.Group("/api/analytics")
	apiAnalytics.Use(middleware.JWTAuth(authSvc))
	analytics.RegisterRoutes(apiAnalytics, analyticsSvc)

	// PROJECTS
	apiProjects := r.Group("/api/projects")
	apiProjects.Use(middleware.JWTAuth(authSvc))
	projects.RegisterRoutes(apiProjects, projectsSvc)

	// SETTINGS
	apiSettings := r.Group("/api/settings")
	apiSettings.Use(middleware.JWTAuth(authSvc))
	settings.RegisterRoutes(apiSettings, settingsSvc)

	// TEMPLATES
	apiTemplates := r.Group("/api/templates")
	apiTemplates.Use(middleware.JWTAuth(authSvc))
	templates.RegisterRoutes(apiTemplates, templatesSvc)

	// BILLING
	apiBilling := r.Group("/api/billing")
	apiBilling.Use(middleware.JWTAuth(authSvc))
	billing.RegisterRoutes(apiBilling, billingSvc)

	// ADMIN
	apiAdmin := r.Group("/api/admin")
	apiAdmin.Use(middleware.JWTAuth(authSvc))
	admin.RegisterRoutes(apiAdmin, adminSvc)

	// AUDIT
	apiAudit := r.Group("/api/audit")
	apiAudit.Use(middleware.JWTAuth(authSvc))
	audit.RegisterRoutes(apiAudit, auditSvc)

	// REDIRECT (public)
	redirect.RegisterRoutes(r, redirectSvc)

	// --------------------------
	// SWAGGER DOCS
	// --------------------------
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// --------------------------
	// START SERVER
	// --------------------------
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Server error:", err)
	}
}
