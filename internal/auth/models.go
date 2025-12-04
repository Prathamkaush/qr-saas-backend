package auth

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // "-" means never send this to frontend
	Name         string    `json:"name"`
	AvatarURL    string    `json:"avatar_url"`
	Provider     string    `json:"provider"` // "google", "local"
	ProviderID   string    `json:"provider_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// Request Models
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse (Keep this one)
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
