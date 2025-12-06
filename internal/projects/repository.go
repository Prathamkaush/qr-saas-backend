package projects

import (
	"context"

	"qr-saas/internal/qr"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, p *Project) error
	ListByUser(ctx context.Context, userID string) ([]Project, error)
	GetByID(ctx context.Context, userID, id string) (*Project, error)
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, userID, id string) error

	// NEW FEATURES
	ListProjectQRs(ctx context.Context, userID, projectID string) ([]qr.QRCode, error)
	AssignQR(ctx context.Context, userID, qrID, projectID string) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

// ------------------------
// CREATE
// ------------------------
func (r *repository) Create(ctx context.Context, p *Project) error {
	_, err := r.pg.Exec(ctx,
		`INSERT INTO projects (id, user_id, name, color, created_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.UserID, p.Name, p.Color, p.CreatedAt,
	)
	return err
}

// ------------------------
// LIST PROJECTS
// ------------------------
func (r *repository) ListByUser(ctx context.Context, userID string) ([]Project, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, user_id, name, color, created_at
		 FROM projects 
		 WHERE user_id=$1 
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Name, &p.Color, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, nil
}

// ------------------------
// GET ONE
// ------------------------
func (r *repository) GetByID(ctx context.Context, userID, id string) (*Project, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, name, color, created_at
		 FROM projects 
		 WHERE id=$1 AND user_id=$2 
		 LIMIT 1`,
		id, userID,
	)

	var p Project
	if err := row.Scan(
		&p.ID, &p.UserID, &p.Name, &p.Color, &p.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

// ------------------------
// UPDATE
// ------------------------
func (r *repository) Update(ctx context.Context, p *Project) error {
	_, err := r.pg.Exec(ctx,
		`UPDATE projects 
		 SET name=$1, color=$2 
		 WHERE id=$3 AND user_id=$4`,
		p.Name, p.Color, p.ID, p.UserID,
	)
	return err
}

// ------------------------
// DELETE PROJECT
// ------------------------
func (r *repository) Delete(ctx context.Context, userID, id string) error {
	_, err := r.pg.Exec(ctx,
		`DELETE FROM projects 
		 WHERE id=$1 AND user_id=$2`,
		id, userID,
	)
	return err
}

// ====================================================================
// NEW: LIST ALL QRs INSIDE A PROJECT
// ====================================================================
func (r *repository) ListProjectQRs(ctx context.Context, userID, projectID string) ([]qr.QRCode, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT 
            id, user_id, project_id, name, qr_type, short_code,
            target_url, design_json, is_active, created_at, updated_at
         FROM qr_codes
         WHERE user_id=$1 AND project_id=$2
         ORDER BY created_at DESC`,
		userID, projectID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []qr.QRCode

	for rows.Next() {
		var q qr.QRCode
		err := rows.Scan(
			&q.ID,
			&q.UserID,
			&q.ProjectID,
			&q.Name,
			&q.QRType,
			&q.ShortCode,
			&q.TargetURL,
			&q.DesignJSON,
			&q.IsActive,
			&q.CreatedAt,
			&q.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}

	return out, nil
}

// ====================================================================
// NEW: ASSIGN QR → PROJECT (MOVE QR)
// ====================================================================
func (r *repository) AssignQR(ctx context.Context, userID, qrID, projectID string) error {
	// If projectID == "", remove QR from project
	var err error
	if projectID == "" {
		_, err = r.pg.Exec(ctx,
			`UPDATE qr_codes 
             SET project_id=NULL 
             WHERE id=$1 AND user_id=$2`,
			qrID, userID,
		)
	} else {
		_, err = r.pg.Exec(ctx,
			`UPDATE qr_codes 
             SET project_id=$1 
             WHERE id=$2 AND user_id=$3`,
			projectID, qrID, userID,
		)
	}

	if err != nil {
		return err
	}

	return nil
}
