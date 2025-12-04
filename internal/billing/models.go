package billing

import "time"

type Plan struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	PriceMonthly int    `json:"price_monthly"`
	PriceYearly  int    `json:"price_yearly"`
	ScanLimit    int    `json:"scan_limit"`
	QRLimit      int    `json:"qr_limit"`
}

type Subscription struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	PlanID    string    `json:"plan_id"`
	Status    string    `json:"status"` // active, canceled, expired
	RenewAt   time.Time `json:"renew_at"`
	CreatedAt time.Time `json:"created_at"`
	StripeID  string    `json:"stripe_id"`
}

type Usage struct {
	UserID    string `json:"user_id"`
	Date      string `json:"date"`
	ScansUsed int    `json:"scans_used"`
	QRCreates int    `json:"qr_creates"`
}

type SubscribeRequest struct {
	PlanID      string `json:"plan_id" binding:"required"`
	StripeToken string `json:"stripe_token" binding:"required"`
}
