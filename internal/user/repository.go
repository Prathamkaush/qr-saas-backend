package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, email, name, avatar_url, created_at
		 FROM users WHERE id=$1 LIMIT 1`, id)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, email, name, avatar_url, created_at
		 FROM users WHERE email=$1 LIMIT 1`, email)

	var u User
	if err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.CreatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) Create(ctx context.Context, u *User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	_, err := r.pg.Exec(ctx,
		`INSERT INTO users (id, email, name, avatar_url, created_at)
		 VALUES ($1,$2,$3,$4,$5)`,
		u.ID, u.Email, u.Name, u.AvatarURL, u.CreatedAt)
	return err
}
