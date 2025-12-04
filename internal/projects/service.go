package projects

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	CreateProject(ctx context.Context, userID, name, color string) (*Project, error)
	ListProjects(ctx context.Context, userID string) ([]Project, error)
	GetProject(ctx context.Context, userID, id string) (*Project, error)
	UpdateProject(ctx context.Context, userID, id string, req UpdateProjectRequest) (*Project, error)
	DeleteProject(ctx context.Context, userID, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateProject(ctx context.Context, userID, name, color string) (*Project, error) {
	if color == "" {
		color = "#3B82F6" // default Tailwind blue
	}
	p := &Project{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      name,
		Color:     color,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) ListProjects(ctx context.Context, userID string) ([]Project, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *service) GetProject(ctx context.Context, userID, id string) (*Project, error) {
	return s.repo.GetByID(ctx, userID, id)
}

func (s *service) UpdateProject(ctx context.Context, userID, id string, req UpdateProjectRequest) (*Project, error) {
	p, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Color != nil {
		p.Color = *req.Color
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) DeleteProject(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, userID, id)
}
