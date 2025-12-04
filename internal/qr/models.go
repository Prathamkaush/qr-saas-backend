package qr

import "time"

type QRCode struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	ProjectID  *string   `json:"project_id,omitempty"`
	Name       string    `json:"name"`
	QRType     string    `json:"qr_type"`
	ShortCode  string    `json:"short_code"`
	TargetURL  string    `json:"target_url"`
	DesignJSON string    `json:"design_json"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
