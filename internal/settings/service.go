package settings

import (
	"context"
	"time"
)

type Service interface {
	GetSettings(ctx context.Context, userID string) (*Settings, error)
	UpdateSettings(ctx context.Context, userID string, req UpdateSettingsRequest) (*Settings, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func defaultSettings(userID string) *Settings {
	return &Settings{
		UserID:             userID,
		Theme:              "light",
		Language:           "en",
		Timezone:           "Asia/Kolkata",
		EmailNotifications: true,
		BrandName:          "",
		CustomDomain:       "",
		LogoURL:            "",
		UpdatedAt:          time.Now().UTC(),
	}
}

func (s *service) GetSettings(ctx context.Context, userID string) (*Settings, error) {
	sett, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if sett == nil {
		// create default
		def := defaultSettings(userID)
		if err := s.repo.Upsert(ctx, def); err != nil {
			return nil, err
		}
		return def, nil
	}
	return sett, nil
}

func (s *service) UpdateSettings(ctx context.Context, userID string, req UpdateSettingsRequest) (*Settings, error) {
	sett, err := s.GetSettings(ctx, userID)
	if err != nil {
		return nil, err
	}

	if req.Theme != nil {
		sett.Theme = *req.Theme
	}
	if req.Language != nil {
		sett.Language = *req.Language
	}
	if req.Timezone != nil {
		sett.Timezone = *req.Timezone
	}
	if req.EmailNotifications != nil {
		sett.EmailNotifications = *req.EmailNotifications
	}
	if req.BrandName != nil {
		sett.BrandName = *req.BrandName
	}
	if req.CustomDomain != nil {
		sett.CustomDomain = *req.CustomDomain
	}
	if req.LogoURL != nil {
		sett.LogoURL = *req.LogoURL
	}

	sett.UpdatedAt = time.Now().UTC()

	if err := s.repo.Upsert(ctx, sett); err != nil {
		return nil, err
	}
	return sett, nil
}
