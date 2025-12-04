package billing

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetPlans(ctx context.Context) ([]Plan, error)
	CreateSubscription(ctx context.Context, sub Subscription) error
	GetSubscriptionByUser(ctx context.Context, userID string) (*Subscription, error)
	UpdateUsage(ctx context.Context, userID string, scans, qrCreates int) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) GetPlans(ctx context.Context) ([]Plan, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, name, price_monthly, price_yearly, scan_limit, qr_limit 
		 FROM billing_plans ORDER BY price_monthly ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []Plan
	for rows.Next() {
		var p Plan
		if err := rows.Scan(&p.ID, &p.Name, &p.PriceMonthly, &p.PriceYearly, &p.ScanLimit, &p.QRLimit); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, nil
}

func (r *repository) CreateSubscription(ctx context.Context, s Subscription) error {
	_, err := r.pg.Exec(ctx,
		`INSERT INTO billing_subscriptions (id, user_id, plan_id, status, renew_at, created_at, stripe_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		s.ID, s.UserID, s.PlanID, s.Status, s.RenewAt, s.CreatedAt, s.StripeID)
	return err
}

func (r *repository) GetSubscriptionByUser(ctx context.Context, userID string) (*Subscription, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, plan_id, status, renew_at, created_at, stripe_id
		 FROM billing_subscriptions WHERE user_id=$1`, userID)

	var s Subscription
	err := row.Scan(&s.ID, &s.UserID, &s.PlanID, &s.Status, &s.RenewAt, &s.CreatedAt, &s.StripeID)

	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) UpdateUsage(ctx context.Context, userID string, scans, qrCreates int) error {
	_, err := r.pg.Exec(ctx,
		`INSERT INTO billing_usage (user_id, date, scans_used, qr_creates)
         VALUES ($1, CURRENT_DATE, $2, $3)
         ON CONFLICT (user_id, date)
         DO UPDATE SET scans_used = billing_usage.scans_used + $2,
                       qr_creates = billing_usage.qr_creates + $3`,
		userID, scans, qrCreates)
	return err
}
