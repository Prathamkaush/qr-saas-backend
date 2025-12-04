// user/service.go
package user

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetOrCreateFromEmail(ctx context.Context, email, name, avatarURL string) (*User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetOrCreateFromEmail(ctx context.Context, email, name, avatarURL string) (*User, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err == nil && u != nil {
		return u, nil
	}

	u = &User{
		ID:        uuid.NewString(),
		Email:     email,
		Name:      name,
		AvatarURL: avatarURL,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
