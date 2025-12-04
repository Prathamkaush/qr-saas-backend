package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	FindByProvider(ctx context.Context, provider, providerID string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) FindByProvider(ctx context.Context, provider, providerID string) (*User, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, email, name, avatar_url, provider, provider_id, created_at
         FROM users 
         WHERE provider=$1 AND provider_id=$2
         LIMIT 1`,
		provider, providerID,
	)

	var u User
	err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Name,
		&u.AvatarURL,
		&u.Provider,
		&u.ProviderID,
		&u.CreatedAt,
	)

	if err != nil {
		// not found is OK → return nil, nil so service can create user
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		// real DB error → bubble up
		return nil, err
	}

	return &u, nil
}

// Update your Create function to include password_hash
func (r *repository) Create(ctx context.Context, u *User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	_, err := r.pg.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name, avatar_url, provider, provider_id, created_at)
         VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		u.ID, u.Email, u.PasswordHash, u.Name, u.AvatarURL, u.Provider, u.ProviderID, u.CreatedAt,
	)
	return err
}

// Add this function to your repository implementation
func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT id, email, password_hash, name, avatar_url, provider, provider_id, created_at
         FROM users 
         WHERE email=$1 LIMIT 1`,
		email,
	)

	var u User
	err := row.Scan(
		&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.AvatarURL,
		&u.Provider, &u.ProviderID, &u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
