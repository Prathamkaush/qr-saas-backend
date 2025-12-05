package qr

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"qr-saas/internal/qr/render"

	"github.com/google/uuid"
)

type Service interface {
	CreateDynamicURL(ctx context.Context, userID, name, targetURL string, design any) (*QRCode, error)
	GenerateQRImage(ctx context.Context, qrID, userID, scene string) ([]byte, error)
	ListByUser(ctx context.Context, userID string) ([]QRCode, error)
	Delete(ctx context.Context, id, userID string) error
}

type service struct {
	repo Repository
	baseURL string // 🔥 ADDED: Field to hold the global deployment URL
}

// 🔥 FIX 1: NewService now accepts the baseURL from main.go
func NewService(repo Repository, baseURL string) Service {
	return &service{repo: repo, baseURL: baseURL}
}

// GenerateShortCode uses crypto/rand for secure, non-colliding strings
func GenerateShortCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

func (s *service) CreateDynamicURL(ctx context.Context, userID, name, targetURL string, design any) (*QRCode, error) {
    if targetURL == "" {
        return nil, errors.New("target_url required")
    }

    if name == "" {
        name = "My QR Code"
    }

    designJSON, _ := json.Marshal(design)
    now := time.Now().UTC()

    // ----------------------------------------------------------------
    // 🔥 FIX 1: Determine if this QR should be STATIC or DYNAMIC
    // ----------------------------------------------------------------
    
    // We only treat URL types as dynamic, which require the redirect service.
    // All other types (vcard, pdf, wifi) are treated as static content
    // and should be encoded directly.
    
    qrType := "dynamic" 
    finalTargetURL := targetURL // Target is the final destination for dynamic links

    if targetURL == "https://example.com" { // Check if frontend sent the placeholder
        // If the frontend used the placeholder, treat it as a generic static code
        qrType = "text" 
        finalTargetURL = "Encoded Static Content" 
    }
    
    // NOTE: This assumes your frontend only calls CreateDynamicURL for dynamic link types.
    // If the frontend sends VCard/Wifi data, it should be adjusted to call a different handler.

    var qr *QRCode
    var err error

    for i := 0; i < 3; i++ {
        shortCode, errGen := GenerateShortCode(6)
        if errGen != nil {
            return nil, errGen
        }

        qr = &QRCode{
            ID:          uuid.NewString(),
            UserID:      userID,
            ProjectID:   nil,
            Name:        name,
            QRType:      qrType, // Set the determined type
            ShortCode:   shortCode,
            TargetURL:   finalTargetURL, // Use the proper target URL
            DesignJSON:  string(designJSON),
            IsActive:    true,
            CreatedAt:   now,
            UpdatedAt:   now,
        }

        err = s.repo.Create(ctx, qr)

        if err == nil { break }
        if !strings.Contains(err.Error(), "duplicate key") && !strings.Contains(err.Error(), "23505") {
            return nil, err
        }
        fmt.Println("⚠️ Collision detected, retrying short code generation...")
    }

    if err != nil {
        return nil, fmt.Errorf("failed to create QR after retries: %w", err)
    }

    return qr, nil
}func (s *service) GenerateQRImage(ctx context.Context, qrID, userID, scene string) ([]byte, error) {
qrData, err := s.repo.GetByID(ctx, qrID, userID)
if err != nil {
return nil, err
}

if qrData.ShortCode == "" {
return nil, errors.New("qr code has no short code")
}

// Determine content to encode into the image
var contentToEncode string

// 1. If it's a DYNAMIC URL, we encode the TRACKING link.
if qrData.QRType == "dynamic" {
contentToEncode = fmt.Sprintf("%s/r/%s", s.baseURL, qrData.ShortCode)
} else {
// 2. If it's STATIC (vcard, wifi, pdf), we encode the raw content.
// NOTE: This assumes raw content is saved in qrData.TargetURL or a dedicated StaticData field.
// For now, we'll use TargetURL as the placeholder for raw content data.
contentToEncode = qrData.TargetURL // This holds the vcard string, wifi string, etc.
}

if contentToEncode == "" {
return nil, errors.New("cannot generate QR for empty content")
}

// Generate QR using the derived content
qrBytes, err := render.RenderQRWithLogo(contentToEncode, render.RenderOptions{
Size: 600,
})
	if err != nil {
		return nil, err
	}

	if scene == "person_pizza" {
		composite, err := render.ComposeQROnBackground(render.CompositeOptions{
			BackgroundPath: "assets/person_pizza.png",
			QRBytes:        qrBytes,
			PosX:           200,
			PosY:           150,
			Width:          250,
			Height:         250,
		})
		return composite, err
	}

	return qrBytes, nil
}

func (s *service) ListByUser(ctx context.Context, userID string) ([]QRCode, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
