package audit

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	LogEvent(ctx context.Context, userID, action, entity, entityID, metadata string) error
	GetUserEvents(ctx context.Context, userID string, limit int) ([]AuditEvent, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) LogEvent(ctx context.Context, userID, action, entity, entityID, metadata string) error {
	a := AuditEvent{
		ID:        uuid.New().String(),
		UserID:    userID,
		Action:    action,
		Entity:    entity,
		EntityID:  entityID,
		Metadata:  metadata,
		CreatedAt: time.Now().UTC(),
	}
	return s.repo.CreateEvent(ctx, a)
}

func (s *service) GetUserEvents(ctx context.Context, userID string, limit int) ([]AuditEvent, error) {
	return s.repo.ListUserEvents(ctx, userID, limit)
}
