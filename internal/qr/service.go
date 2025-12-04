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
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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

	// basic safety: give default name if empty
	if name == "" {
		name = "My QR Code"
	}

	designJSON, _ := json.Marshal(design)
	now := time.Now().UTC()

	var qr *QRCode
	var err error

	// ---------------------------------------------------------
	// 🔥 RETRY LOOP: Try 3 times to generate unique code
	// ---------------------------------------------------------
	for i := 0; i < 3; i++ {
		// 1. Generate Secure Code
		shortCode, errGen := GenerateShortCode(6)
		if errGen != nil {
			return nil, errGen
		}

		qr = &QRCode{
			ID:         uuid.NewString(),
			UserID:     userID,
			ProjectID:  nil,
			Name:       name,
			QRType:     "dynamic",
			ShortCode:  shortCode, // Uses the generated code
			TargetURL:  targetURL,
			DesignJSON: string(designJSON),
			IsActive:   true,
			CreatedAt:  now,
			UpdatedAt:  now,
		}

		// 2. Try to insert into Database
		err = s.repo.Create(ctx, qr)

		// If success (err == nil), break the loop!
		if err == nil {
			break
		}

		// If error is NOT a duplicate key, fail immediately
		// (Postgres error 23505 is unique_violation)
		if !strings.Contains(err.Error(), "duplicate key") && !strings.Contains(err.Error(), "23505") {
			return nil, err
		}

		// If it WAS a duplicate, the loop continues and tries a new code...
		fmt.Println("⚠️ Collision detected, retrying short code generation...")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create QR after retries: %w", err)
	}

	return qr, nil
}

func (s *service) GenerateQRImage(ctx context.Context, qrID, userID, scene string) ([]byte, error) {
	qrData, err := s.repo.GetByID(ctx, qrID, userID)
	if err != nil {
		return nil, err
	}

	// ---------------------------------------------------------
	// 🔥 CRITICAL FIX: Encode the Redirect URL, NOT the Target URL
	// ---------------------------------------------------------
	// If you encode TargetURL directly, the phone goes straight to the website.
	// The backend NEVER sees the scan, so NO analytics are logged.

	// Check if ShortCode exists
	if qrData.ShortCode == "" {
		return nil, errors.New("qr code has no short code")
	}

	// 1. Construct the Tracking URL
	// This hits your Go server (/r/:code) -> Logs to ClickHouse -> Redirects to Target
	// IMPORTANT: For local testing with a PHONE, 'localhost' won't work.
	// You need your computer's IP (e.g., http://192.168.1.5:8080/r/...) or ngrok.
	url := fmt.Sprintf("http://localhost:8080/r/%s", qrData.ShortCode)

	// 2. Generate QR using the Tracking URL
	qrBytes, err := render.RenderQRWithLogo(url, render.RenderOptions{
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

	// default: return plain QR
	return qrBytes, nil
}

func (s *service) ListByUser(ctx context.Context, userID string) ([]QRCode, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
