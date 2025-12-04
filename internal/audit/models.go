package audit

import "time"

// AuditEvent describes an action performed in the system.
type AuditEvent struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	EntityID  string    `json:"entity_id"` // qrID, userID, projectID, etc.
	Entity    string    `json:"entity"`    // "qr", "user", "project", "billing"
	Metadata  string    `json:"metadata"`  // JSON string
	CreatedAt time.Time `json:"created_at"`
}

type CreateEventRequest struct {
	Action   string `json:"action" binding:"required"`
	Entity   string `json:"entity" binding:"required"`
	EntityID string `json:"entity_id" binding:"required"`
	Metadata string `json:"metadata"`
}
