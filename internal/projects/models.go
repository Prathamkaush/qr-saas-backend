package projects

import "time"

type Project struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	Count     int       `json:"count"`
}

type CreateProjectRequest struct {
	Name  string `json:"name" binding:"required"`
	Color string `json:"color"`
}

type UpdateProjectRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}
