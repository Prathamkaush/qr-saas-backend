package billing

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	GetPlans(ctx context.Context) ([]Plan, error)
	Subscribe(ctx context.Context, userID string, planID string, stripeToken string) error
	GetActiveSubscription(ctx context.Context, userID string) (*Subscription, error)
	AddUsage(ctx context.Context, userID string, scans, qrCreates int) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) GetPlans(ctx context.Context) ([]Plan, error) {
	return s.repo.GetPlans(ctx)
}

func (s *service) Subscribe(ctx context.Context, userID string, planID string, stripeToken string) error {
	// Normally you'd process stripe token here.
	// For now, we will simulate.
	sub := Subscription{
		ID:        uuid.New().String(),
		UserID:    userID,
		PlanID:    planID,
		Status:    "active",
		RenewAt:   time.Now().AddDate(0, 1, 0),
		CreatedAt: time.Now(),
		StripeID:  "fake_stripe_id",
	}
	return s.repo.CreateSubscription(ctx, sub)
}

func (s *service) GetActiveSubscription(ctx context.Context, userID string) (*Subscription, error) {
	return s.repo.GetSubscriptionByUser(ctx, userID)
}

func (s *service) AddUsage(ctx context.Context, userID string, scans, qrCreates int) error {
	return s.repo.UpdateUsage(ctx, userID, scans, qrCreates)
}
