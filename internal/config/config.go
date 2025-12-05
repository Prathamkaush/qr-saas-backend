package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort string

	// Deleted U+00A0 characters from fields below
	ClickHouseDSN      string
	ClickHouseUser     string
	ClickHousePassword string
	ClickHouseDatabase string

	RedisURL  string
	JWTSecret string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	PGURL string
	
	// FIX 1: Add Frontend URL for dynamic routing and CORS
	FrontendURL string 
    
    // FIX 2: Add BaseURL (The public host of this API)
	BaseURL string
}

func Load() Config {
	// Get project root reliably
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(filepath.Dir(filepath.Dir(b)))

	envPath := filepath.Join(basepath, ".env")
	_ = godotenv.Load(envPath)

	fmt.Println("Loaded .env from:", envPath)
    
    // FIX 3: Prioritize standard cloud 'PORT' environment variable
    port := getEnv("PORT", getEnv("HTTP_PORT", "8080")) 

	cfg := Config{
		HTTPPort: port, // Use the priority port

		// Deleted U+00A0 characters from keys below
		ClickHouseDSN:      getEnv("CLICKHOUSE_DSN", ""),
		ClickHouseUser:     getEnv("CLICKHOUSE_USER", "default"),
		ClickHousePassword: getEnv("CLICKHOUSE_PASSWORD", ""),
		ClickHouseDatabase: getEnv("CLICKHOUSE_DATABASE", "qr_saas"),

		RedisURL:  getEnv("REDIS_URL", ""),
		JWTSecret: getEnv("JWT_SECRET", ""),

		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),

		PGURL: getEnv("PG_URL", ""),
		
		// Read Frontend URL (used by auth redirect)
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
        
        // Read Base URL (used by QR image generation for tracking)
		BaseURL: getEnv("BASE_URL", "http://localhost:8080"), 
	}

	fmt.Println("CLICKHOUSE_HOST LOADED =>", cfg.ClickHouseDSN)
	fmt.Println("JWT_SECRET LOADED =>", cfg.JWTSecret)
	fmt.Println("PG_URL LOADED =>", cfg.PGURL)

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
