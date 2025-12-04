package analytics

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // <--- Switched to pgx
)

// ScanEvent matches your Postgres table
type ScanEvent struct {
	EventID    string
	QRID       string
	UserID     string
	ScannedAt  time.Time
	IP         string
	Country    string
	City       string
	UserAgent  string
	DeviceType string
	OS         string
	Browser    string
	Referer    string
}

// Summary Struct (Same as before)
type Summary struct {
	TotalScans int64          `json:"total_scans"`
	UniqueIPs  int64          `json:"unique_ips"`
	Countries  map[string]int `json:"countries"`
	Devices    map[string]int `json:"devices"`
	Browsers   map[string]int `json:"browsers"`
}

type TimePoint struct {
	Timestamp time.Time `json:"timestamp"`
	Count     int64     `json:"count"`
}

type Repository interface {
	InsertScanEvent(ctx context.Context, ev ScanEvent) error
	GetSummary(ctx context.Context, userID, qrID string, from, to time.Time) (*Summary, error)
	GetGlobalStats(ctx context.Context, userID string) (*Summary, error)
	GetTimeSeries(ctx context.Context, userID, qrID string, from, to time.Time, granularity string) ([]TimePoint, error)
	GetGlobalTimeSeries(ctx context.Context, userID string, from, to time.Time, granularity string) ([]TimePoint, error)
}

type repository struct {
	db *pgxpool.Pool // <--- Changed from clickhouse.Conn
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

// ---------------------------------------------------------
// 1. INSERT (Postgres)
// ---------------------------------------------------------
func (r *repository) InsertScanEvent(ctx context.Context, ev ScanEvent) error {
	query := `
		INSERT INTO scan_events (
			id, qr_id, user_id, scanned_at, 
			ip, country, city, user_agent, 
			device_type, os, browser, referer
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		ev.EventID, ev.QRID, ev.UserID, ev.ScannedAt,
		ev.IP, ev.Country, ev.City, ev.UserAgent,
		ev.DeviceType, ev.OS, ev.Browser, ev.Referer,
	)
	return err
}

// ---------------------------------------------------------
// 2. GET SUMMARY (Single QR)
// ---------------------------------------------------------
func (r *repository) GetSummary(ctx context.Context, userID, qrID string, from, to time.Time) (*Summary, error) {
	return r.fetchStats(ctx, userID, qrID, from, to)
}

// ---------------------------------------------------------
// 3. GET GLOBAL STATS (Dashboard)
// ---------------------------------------------------------
func (r *repository) GetGlobalStats(ctx context.Context, userID string) (*Summary, error) {
	// For global dashboard, we usually want all-time or last 30 days
	// Let's default to all-time for the counters
	from := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Now().Add(24 * time.Hour)
	return r.fetchStats(ctx, userID, "", from, to)
}

// --- Shared Helper for Stats ---
func (r *repository) fetchStats(ctx context.Context, userID, qrID string, from, to time.Time) (*Summary, error) {
	summary := &Summary{
		Countries: make(map[string]int),
		Devices:   make(map[string]int),
		Browsers:  make(map[string]int),
	}

	// 1. Totals
	countQuery := `
		SELECT count(*), count(DISTINCT ip) 
		FROM scan_events 
		WHERE user_id = $1 AND scanned_at BETWEEN $2 AND $3
	`
	args := []interface{}{userID, from, to}
	if qrID != "" {
		countQuery += " AND qr_id = $4"
		args = append(args, qrID)
	}

	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&summary.TotalScans, &summary.UniqueIPs)
	if err != nil {
		return nil, err
	}

	if summary.TotalScans == 0 {
		return summary, nil
	}

	// 2. Breakdowns (Helper function usage)
	summary.Countries, _ = r.getBreakdown(ctx, "country", userID, qrID, from, to)
	summary.Devices, _ = r.getBreakdown(ctx, "device_type", userID, qrID, from, to)
	summary.Browsers, _ = r.getBreakdown(ctx, "browser", userID, qrID, from, to)

	return summary, nil
}

func (r *repository) getBreakdown(ctx context.Context, column, userID, qrID string, from, to time.Time) (map[string]int, error) {
	// Postgres SQL uses $1, $2 syntax
	query := fmt.Sprintf(`
		SELECT %s, count(*) as c
		FROM scan_events
		WHERE user_id = $1 AND scanned_at BETWEEN $2 AND $3
	`, column)

	args := []interface{}{userID, from, to}
	if qrID != "" {
		query += " AND qr_id = $4"
		args = append(args, qrID)
	}

	query += fmt.Sprintf(" GROUP BY %s ORDER BY c DESC LIMIT 5", column)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]int)
	for rows.Next() {
		var name string
		var count int64
		// Handle NULLs in Postgres
		if err := rows.Scan(&name, &count); err == nil {
			if name == "" {
				name = "Unknown"
			}
			result[name] = int(count)
		}
	}
	return result, nil
}

// ---------------------------------------------------------
// 4. TIME SERIES (Graph)
// ---------------------------------------------------------
func (r *repository) GetTimeSeries(ctx context.Context, userID, qrID string, from, to time.Time, granularity string) ([]TimePoint, error) {
	return r.fetchTimeSeries(ctx, userID, qrID, from, to)
}

func (r *repository) GetGlobalTimeSeries(ctx context.Context, userID string, from, to time.Time, granularity string) ([]TimePoint, error) {
	return r.fetchTimeSeries(ctx, userID, "", from, to)
}

func (r *repository) fetchTimeSeries(ctx context.Context, userID, qrID string, from, to time.Time) ([]TimePoint, error) {
	// Postgres date_trunc function is great for this
	query := `
		SELECT date_trunc('day', scanned_at) as ts, count(*) 
		FROM scan_events 
		WHERE user_id = $1 AND scanned_at BETWEEN $2 AND $3
	`
	args := []interface{}{userID, from, to}

	if qrID != "" {
		query += " AND qr_id = $4"
		args = append(args, qrID)
	}

	query += " GROUP BY ts ORDER BY ts ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []TimePoint
	for rows.Next() {
		var p TimePoint
		if err := rows.Scan(&p.Timestamp, &p.Count); err == nil {
			points = append(points, p)
		}
	}
	return points, nil
}
