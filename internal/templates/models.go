package templates

import "time"

type Template struct {
	ID         string                 `json:"id"`
	UserID     *string                `json:"user_id"` // null = global template
	Category   string                 `json:"category"`
	Name       string                 `json:"name"`
	Thumbnail  string                 `json:"thumbnail"`   // image URL
	DesignJSON map[string]interface{} `json:"design_json"` // frame, colors, shapes, logo position
	CreatedAt  time.Time              `json:"created_at"`
}

type CreateTemplateRequest struct {
	Category   string                 `json:"category" binding:"required"`
	Name       string                 `json:"name" binding:"required"`
	Thumbnail  string                 `json:"thumbnail"`
	DesignJSON map[string]interface{} `json:"design_json" binding:"required"`
}

type UpdateTemplateRequest struct {
	Category   *string                `json:"category"`
	Name       *string                `json:"name"`
	Thumbnail  *string                `json:"thumbnail"`
	DesignJSON map[string]interface{} `json:"design_json"`
}
type TemplateInstance struct {
	ID         string                 `json:"id"`
	UserID     string                 `json:"user_id"`
	TemplateID string                 `json:"template_id"`
	URLID      string                 `json:"url_id"`
	Data       map[string]interface{} `json:"data"`
	CreatedAt  time.Time              `json:"created_at"`
}
