package qrtypes

import "context"

type Repository interface {
	SaveTypeData(ctx context.Context, qrID string, data interface{}) error
	GetTypeData(ctx context.Context, qrID string) (map[string]interface{}, error)
}
