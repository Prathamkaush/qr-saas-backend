package templates

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID string, req CreateTemplateRequest) (*Template, error)
	ListGlobal(ctx context.Context) ([]Template, error)
	ListMine(ctx context.Context, userID string) ([]Template, error)
	Get(ctx context.Context, userID, id string) (*Template, error)
	Update(ctx context.Context, userID, id string, req UpdateTemplateRequest) (*Template, error)
	Delete(ctx context.Context, userID, id string) error

	// Add this ↓↓↓
	RenderPublicPage(ctx context.Context, urlID string) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Create(ctx context.Context, userID string, req CreateTemplateRequest) (*Template, error) {
	t := &Template{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Category:   req.Category,
		Name:       req.Name,
		Thumbnail:  req.Thumbnail,
		DesignJSON: req.DesignJSON,
		CreatedAt:  time.Now().UTC(),
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *service) ListGlobal(ctx context.Context) ([]Template, error) {
	return s.repo.ListGlobal(ctx)
}

func (s *service) ListMine(ctx context.Context, userID string) ([]Template, error) {
	return s.repo.ListUserTemplates(ctx, userID)
}

func (s *service) Get(ctx context.Context, userID, id string) (*Template, error) {
	return s.repo.GetByID(ctx, id, userID)
}

func (s *service) Update(ctx context.Context, userID, id string, req UpdateTemplateRequest) (*Template, error) {
	t, err := s.repo.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if t == nil || t.UserID == nil {
		return nil, nil
	}

	if req.Category != nil {
		t.Category = *req.Category
	}
	if req.Name != nil {
		t.Name = *req.Name
	}
	if req.Thumbnail != nil {
		t.Thumbnail = *req.Thumbnail
	}
	if req.DesignJSON != nil {
		t.DesignJSON = req.DesignJSON
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *service) Delete(ctx context.Context, userID, id string) error {
	return s.repo.Delete(ctx, id, userID)
}
func (s *service) RenderPublicPage(ctx context.Context, urlID string) (string, error) {
	instance, err := s.repo.GetInstanceByURL(ctx, urlID)
	if err != nil {
		return "", err
	}

	// THIS ONE fetches template without user permission checks
	template, err := s.repo.GetTemplateMeta(ctx, instance.TemplateID)
	if err != nil || template == nil {
		return "", err
	}

	switch template.Category {
	case "vcard":
		return RenderVCard(template, instance)
	case "social":
		return RenderSocial(template, instance)
	case "event":
		return RenderEvent(template, instance)
	default:
		return RenderGeneric(template, instance)
	}
}
