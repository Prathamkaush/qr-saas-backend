package projects

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, p *Project) error
	ListByUser(ctx context.Context, userID string) ([]Project, error)
	GetByID(ctx context.Context, userID, id string) (*Project, error)
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, userID, id string) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) Create(ctx context.Context, p *Project) error {
	_, err := r.pg.Exec(ctx,
		`INSERT INTO projects (id, user_id, name, color, created_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		p.ID, p.UserID, p.Name, p.Color, p.CreatedAt)
	return err
}

func (r *repository) ListByUser(ctx context.Context, userID string) ([]Project, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, user_id, name, color, created_at
		 FROM projects WHERE user_id=$1 ORDER BY created_at DESC`, userID)
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

func (r *repository) GetByID(ctx context.Context, userID, id string) (*Project, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, name, color, created_at
         FROM projects WHERE id=$1 AND user_id=$2 LIMIT 1`,
		id, userID)

	var p Project
	if err := row.Scan(
		&p.ID, &p.UserID, &p.Name, &p.Color, &p.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *repository) Update(ctx context.Context, p *Project) error {
	_, err := r.pg.Exec(ctx,
		`UPDATE projects SET name=$1, color=$2 WHERE id=$3 AND user_id=$4`,
		p.Name, p.Color, p.ID, p.UserID)
	return err
}

func (r *repository) Delete(ctx context.Context, userID, id string) error {
	_, err := r.pg.Exec(ctx,
		`DELETE FROM projects WHERE id=$1 AND user_id=$2`, id, userID)
	return err
}
