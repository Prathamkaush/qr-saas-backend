// shared/types.go
package shared

import "time"

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type PageMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
