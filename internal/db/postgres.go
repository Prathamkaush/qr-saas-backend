package db

import (
	"context"
	"log"
	"time"

	"qr-saas/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(cfg config.Config) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use PGURL from config
	pool, err := pgxpool.New(ctx, cfg.PGURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ Postgres ping failed: %v", err)
	}

	log.Println("✅ Connected to Postgres")
	return pool
}
