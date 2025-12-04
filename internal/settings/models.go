package settings

import "time"

type Settings struct {
	UserID             string `json:"user_id"`
	Theme              string `json:"theme"`               // light, dark
	Language           string `json:"language"`            // en, hi, etc
	Timezone           string `json:"timezone"`            // Asia/Kolkata, etc
	EmailNotifications bool   `json:"email_notifications"` // true/false

	// White-label / branding
	BrandName    string `json:"brand_name"`
	CustomDomain string `json:"custom_domain"`
	LogoURL      string `json:"logo_url"`

	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateSettingsRequest struct {
	Theme              *string `json:"theme"`
	Language           *string `json:"language"`
	Timezone           *string `json:"timezone"`
	EmailNotifications *bool   `json:"email_notifications"`
	BrandName          *string `json:"brand_name"`
	CustomDomain       *string `json:"custom_domain"`
	LogoURL            *string `json:"logo_url"`
}
