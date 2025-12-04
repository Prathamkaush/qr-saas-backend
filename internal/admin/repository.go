package admin

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	ListUsers(ctx context.Context) ([]UserListItem, error)
	UpdateUserRole(ctx context.Context, userID, role string) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}

func (r *repository) ListUsers(ctx context.Context) ([]UserListItem, error) {
	rows, err := r.pg.Query(ctx,
		`SELECT id, email, role, created_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []UserListItem
	for rows.Next() {
		var u UserListItem
		err := rows.Scan(&u.ID, &u.Email, &u.Role, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}

func (r *repository) UpdateUserRole(ctx context.Context, userID, role string) error {
	_, err := r.pg.Exec(ctx,
		`UPDATE users SET role=$1 WHERE id=$2`, role, userID)
	return err
}
