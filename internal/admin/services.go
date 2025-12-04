package admin

import "context"

type Service interface {
	ListUsers(ctx context.Context) ([]UserListItem, error)
	UpdateUserRole(ctx context.Context, userID, role string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) ListUsers(ctx context.Context) ([]UserListItem, error) {
	return s.repo.ListUsers(ctx)
}

func (s *service) UpdateUserRole(ctx context.Context, userID, role string) error {
	return s.repo.UpdateUserRole(ctx, userID, role)
}
