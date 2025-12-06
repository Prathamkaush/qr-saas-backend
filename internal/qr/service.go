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
	CreateDynamicURL(ctx context.Context, userID, name, targetURL string, qrType string, design any) (*QRCode, error)
	GenerateQRImage(ctx context.Context, qrID, userID, scene string) ([]byte, error)
	ListByUser(ctx context.Context, userID string) ([]QRCode, error)
	Delete(ctx context.Context, id, userID string) error
}

type service struct {
	repo    Repository
	baseURL string
}

func NewService(repo Repository, baseURL string) Service {
	return &service{repo: repo, baseURL: baseURL}
}

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

func (s *service) CreateDynamicURL(ctx context.Context, userID, name, targetURL string, qrType string, design any) (*QRCode, error) {
	if targetURL == "" {
		return nil, errors.New("target_url required")
	}
	if name == "" {
		name = "My QR Code"
	}

	designJSON, _ := json.Marshal(design)
	now := time.Now().UTC()

	// 1. Determine Logic based on QR Type
	finalQRType := qrType
	finalTargetURL := targetURL

	// If it is a URL, we treat it as dynamic (redirects via server)
	if qrType == "url" {
		finalQRType = "dynamic"
	}
	// For "wifi", "vcard", "text", etc., we keep the type as-is so 
	// GenerateQRImage knows to encode the raw content (TargetURL) directly.

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
			QRType:      finalQRType,
			ShortCode:   shortCode,
			TargetURL:   finalTargetURL,
			DesignJSON:  string(designJSON),
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		err = s.repo.Create(ctx, qr)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "duplicate key") && !strings.Contains(err.Error(), "23505") {
			return nil, err
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create QR after retries: %w", err)
	}

	return qr, nil
}

type DesignConfig struct {
	Color   string `json:"color"`
	BgColor string `json:"bgColor"`
}

func (s *service) GenerateQRImage(ctx context.Context, qrID, userID, scene string) ([]byte, error) {
	qrData, err := s.repo.GetByID(ctx, qrID, userID)
	if err != nil {
		return nil, err
	}

	if qrData.ShortCode == "" {
		return nil, errors.New("qr code has no short code")
	}

	// 1. Determine Content to Encode
	var contentToEncode string
	if qrData.QRType == "dynamic" || qrData.QRType == "url" {
		contentToEncode = fmt.Sprintf("%s/r/%s", s.baseURL, qrData.ShortCode)
	} else {
		contentToEncode = qrData.TargetURL
	}

	if contentToEncode == "" {
		return nil, errors.New("cannot generate QR for empty content")
	}

	// 2. Extract Colors
	var design DesignConfig
	fgColor := "#000000"
	bgColor := "#ffffff"

	if qrData.DesignJSON != "" {
		if err := json.Unmarshal([]byte(qrData.DesignJSON), &design); err == nil {
			if design.Color != "" {
				fgColor = design.Color
			}
			if design.BgColor != "" {
				bgColor = design.BgColor
			}
		}
	}

	// 3. Render QR with Options
	qrBytes, err := render.RenderQRWithLogo(contentToEncode, render.RenderOptions{
		Size:            600,
		Color:           fgColor,
		BackgroundColor: bgColor,
	})
	if err != nil {
		return nil, err
	}

	if scene == "person_pizza" {
		return render.ComposeQROnBackground(render.CompositeOptions{
			BackgroundPath: "assets/person_pizza.png",
			QRBytes:        qrBytes,
			PosX:           200,
			PosY:           150,
			Width:          250,
			Height:         250,
		})
	}

	return qrBytes, nil
}

func (s *service) ListByUser(ctx context.Context, userID string) ([]QRCode, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) Delete(ctx context.Context, id, userID string) error {
	return s.repo.Delete(ctx, id, userID)
}
