package templates

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, t *Template) error
	ListGlobal(ctx context.Context) ([]Template, error)
	ListUserTemplates(ctx context.Context, userID string) ([]Template, error)
	GetByID(ctx context.Context, id string, userID string) (*Template, error)
	Update(ctx context.Context, t *Template) error
	Delete(ctx context.Context, id string, userID string) error

	GetInstanceByURL(ctx context.Context, urlID string) (*TemplateInstance, error)
	GetTemplateMeta(ctx context.Context, id string) (*Template, error)
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func scanTemplate(row pgx.Row) (*Template, error) {
	var t Template
	var designBytes []byte
	var userID *string

	err := row.Scan(
		&t.ID,
		&userID,
		&t.Category,
		&t.Name,
		&t.Thumbnail,
		&designBytes,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	t.UserID = userID
	_ = json.Unmarshal(designBytes, &t.DesignJSON)
	return &t, nil
}

func (r *repository) Create(ctx context.Context, t *Template) error {
	designBytes, _ := json.Marshal(t.DesignJSON)

	_, err := r.pg.Exec(ctx,
		`INSERT INTO templates (id, user_id, category, name, thumbnail, design_json, created_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		t.ID, t.UserID, t.Category, t.Name, t.Thumbnail, designBytes, t.CreatedAt,
	)
	return err
}

func (r *repository) ListGlobal(ctx context.Context) ([]Template, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, user_id, category, name, thumbnail, design_json, created_at
         FROM templates WHERE user_id IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Template
	for rows.Next() {
		var t Template
		var designBytes []byte
		var uID *string

		if err := rows.Scan(
			&t.ID, &uID, &t.Category, &t.Name,
			&t.Thumbnail, &designBytes, &t.CreatedAt,
		); err != nil {
			return nil, err
		}

		t.UserID = uID
		_ = json.Unmarshal(designBytes, &t.DesignJSON)
		out = append(out, t)
	}

	return out, nil
}

func (r *repository) ListUserTemplates(ctx context.Context, userID string) ([]Template, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, user_id, category, name, thumbnail, design_json, created_at
         FROM templates WHERE user_id=$1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Template
	for rows.Next() {
		var t Template
		var designBytes []byte
		var uID *string

		if err := rows.Scan(
			&t.ID, &uID, &t.Category, &t.Name,
			&t.Thumbnail, &designBytes, &t.CreatedAt,
		); err != nil {
			return nil, err
		}

		t.UserID = uID
		_ = json.Unmarshal(designBytes, &t.DesignJSON)
		out = append(out, t)
	}

	return out, nil
}

func (r *repository) GetByID(ctx context.Context, id string, userID string) (*Template, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, category, name, thumbnail, design_json, created_at
         FROM templates WHERE id=$1 AND (user_id=$2 OR user_id IS NULL)`,
		id, userID)

	return scanTemplate(row)
}

func (r *repository) Update(ctx context.Context, t *Template) error {
	designBytes, _ := json.Marshal(t.DesignJSON)

	_, err := r.pg.Exec(ctx,
		`UPDATE templates
         SET category=$1, name=$2, thumbnail=$3, design_json=$4
         WHERE id=$5 AND user_id=$6`,
		t.Category, t.Name, t.Thumbnail, designBytes, t.ID, t.UserID,
	)
	return err
}

func (r *repository) Delete(ctx context.Context, id string, userID string) error {
	_, err := r.pg.Exec(ctx,
		`DELETE FROM templates WHERE id=$1 AND user_id=$2`,
		id, userID)
	return err
}

func (r *repository) GetInstanceByURL(ctx context.Context, urlID string) (*TemplateInstance, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, template_id, data, url_id, created_at
         FROM template_data WHERE url_id=$1`,
		urlID)

	var inst TemplateInstance
	var dataBytes []byte

	err := row.Scan(&inst.ID, &inst.UserID, &inst.TemplateID, &dataBytes, &inst.URLID, &inst.CreatedAt)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(dataBytes, &inst.Data)
	return &inst, nil
}

func (r *repository) GetTemplateMeta(ctx context.Context, id string) (*Template, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, user_id, category, name, thumbnail, design_json, created_at 
		 FROM templates WHERE id=$1`,
		id)

	return scanTemplate(row)
}
