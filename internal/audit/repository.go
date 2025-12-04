package audit

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateEvent(ctx context.Context, event AuditEvent) error
	ListUserEvents(ctx context.Context, userID string, limit int) ([]AuditEvent, error)
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) CreateEvent(ctx context.Context, e AuditEvent) error {
	_, err := r.pg.Exec(ctx,
		`INSERT INTO audit_logs (id, user_id, action, entity, entity_id, metadata, created_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		e.ID, e.UserID, e.Action, e.Entity, e.EntityID, e.Metadata, e.CreatedAt)
	return err
}

func (r *repository) ListUserEvents(ctx context.Context, userID string, limit int) ([]AuditEvent, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, user_id, action, entity, entity_id, metadata, created_at
		 FROM audit_logs WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2`,
		userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditEvent
	for rows.Next() {
		var a AuditEvent
		err := rows.Scan(&a.ID, &a.UserID, &a.Action, &a.Entity, &a.EntityID, &a.Metadata, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, nil
}
