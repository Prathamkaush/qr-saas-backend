// internal/analytics/service.go
package analytics

import (
	"context"
	"time"
)

type Service interface {
	InsertScanEvent(ctx context.Context, ev ScanEvent) error
	GetSummary(ctx context.Context, userID, qrID string, from, to time.Time) (*Summary, error)
	GetTimeSeries(ctx context.Context, userID, qrID string, from, to time.Time, granularity string) ([]TimePoint, error)
	GetGlobalStats(ctx context.Context, userID string) (*Summary, error)
	GetGlobalTimeSeries(ctx context.Context, userID string, from, to time.Time, granularity string) ([]TimePoint, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) InsertScanEvent(ctx context.Context, ev ScanEvent) error {
	return s.repo.InsertScanEvent(ctx, ev)
}

func (s *service) GetSummary(ctx context.Context, userID, qrID string, from, to time.Time) (*Summary, error) {
	return s.repo.GetSummary(ctx, userID, qrID, from, to)
}

func (s *service) GetTimeSeries(ctx context.Context, userID, qrID string, from, to time.Time, granularity string) ([]TimePoint, error) {
	if granularity == "" {
		granularity = "day"
	}
	return s.repo.GetTimeSeries(ctx, userID, qrID, from, to, granularity)
}

func (s *service) GetGlobalStats(ctx context.Context, userID string) (*Summary, error) {

	return s.repo.GetGlobalStats(ctx, userID)

}

func (s *service) GetGlobalTimeSeries(ctx context.Context, userID string, from, to time.Time, granularity string) ([]TimePoint, error) {
	if granularity == "" {
		granularity = "day"
	}
	return s.repo.GetGlobalTimeSeries(ctx, userID, from, to, granularity)
}
