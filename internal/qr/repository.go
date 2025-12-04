package qr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, qr *QRCode) error
	GetByID(ctx context.Context, id, userID string) (*QRCode, error)
	GetByShortCode(ctx context.Context, shortCode string) (*QRCode, error)
	ListByUser(ctx context.Context, userID string) ([]QRCode, error)
	Delete(ctx context.Context, id, userID string) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg: pg}
}

func (r *repository) Create(ctx context.Context, qr *QRCode) error {
	if qr.CreatedAt.IsZero() {
		qr.CreatedAt = time.Now().UTC()
	}
	qr.UpdatedAt = qr.CreatedAt

	_, err := r.pg.Exec(ctx, `
		INSERT INTO qr_codes (
			id,
			user_id,
			project_id,
			name,
			qr_type,
			short_code,
			target_url,
			design_json,
			is_active,
			created_at,
			updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
	`,
		qr.ID,
		qr.UserID,
		qr.ProjectID,
		qr.Name,
		qr.QRType,
		qr.ShortCode,
		qr.TargetURL,
		qr.DesignJSON,
		qr.IsActive,
		qr.CreatedAt,
		qr.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert qr_codes failed: %w", err)
	}

	// TODO: later also write into ClickHouse qr_codes for analytics
	return nil
}

func (r *repository) GetByID(ctx context.Context, id, userID string) (*QRCode, error) {
	row := r.pg.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			project_id,
			name,
			qr_type,
			short_code,
			target_url,
			design_json,
			is_active,
			created_at,
			updated_at
		FROM qr_codes
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`, id, userID)

	var qr QRCode
	if err := row.Scan(
		&qr.ID,
		&qr.UserID,
		&qr.ProjectID,
		&qr.Name,
		&qr.QRType,
		&qr.ShortCode,
		&qr.TargetURL,
		&qr.DesignJSON,
		&qr.IsActive,
		&qr.CreatedAt,
		&qr.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &qr, nil
}

func (r *repository) GetByShortCode(ctx context.Context, code string) (*QRCode, error) {
	row := r.pg.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			project_id,
			name,
			qr_type,
			short_code,
			target_url,
			design_json,
			is_active,
			created_at,
			updated_at
		FROM qr_codes
		WHERE short_code = $1
		LIMIT 1
	`, code)

	var qr QRCode
	if err := row.Scan(
		&qr.ID,
		&qr.UserID,
		&qr.ProjectID,
		&qr.Name,
		&qr.QRType,
		&qr.ShortCode,
		&qr.TargetURL,
		&qr.DesignJSON,
		&qr.IsActive,
		&qr.CreatedAt,
		&qr.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &qr, nil
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]QRCode, error) {
	rows, err := r.pg.Query(ctx, `
        SELECT
            id,
            user_id,
            project_id,
            name,
            qr_type,
            short_code,
            target_url,
            design_json,
            is_active,
            created_at,
            updated_at
        FROM qr_codes
        WHERE user_id = $1
        ORDER BY created_at DESC
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []QRCode
	for rows.Next() {
		var qr QRCode
		if err := rows.Scan(
			&qr.ID,
			&qr.UserID,
			&qr.ProjectID,
			&qr.Name,
			&qr.QRType,
			&qr.ShortCode,
			&qr.TargetURL,
			&qr.DesignJSON,
			&qr.IsActive,
			&qr.CreatedAt,
			&qr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, qr)
	}
	return out, nil
}

func (r *repository) Delete(ctx context.Context, id, userID string) error {
	// We include userID in the WHERE clause for security.
	// This prevents User A from deleting User B's QR code.
	query := `DELETE FROM qr_codes WHERE id=$1 AND user_id=$2`

	tag, err := r.pg.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	if tag.RowsAffected() == 0 {
		return errors.New("qr code not found or access denied")
	}

	return nil
}
