package qrtypes

import "context"

type Service interface {
	CreateWiFi(ctx context.Context, qrID string, payload WiFiData) error
	CreateVCard(ctx context.Context, qrID string, payload VCardData) error
	GetQRTypeData(ctx context.Context, qrID string) (interface{}, error)
}
