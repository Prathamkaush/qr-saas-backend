package settings

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	GetByUserID(ctx context.Context, userID string) (*Settings, error)
	Upsert(ctx context.Context, s *Settings) error
}

type repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) Repository {
	return &repository{pg}
}
func (r *repository) GetByUserID(ctx context.Context, userID string) (*Settings, error) {
	row := r.pg.QueryRow(ctx,
		`SELECT user_id, theme, language, timezone, email_notifications,
                brand_name, custom_domain, logo_url, updated_at
         FROM user_settings WHERE user_id = $1`, userID)

	var s Settings
	err := row.Scan(
		&s.UserID,
		&s.Theme,
		&s.Language,
		&s.Timezone,
		&s.EmailNotifications,
		&s.BrandName,
		&s.CustomDomain,
		&s.LogoURL,
		&s.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) Upsert(ctx context.Context, s *Settings) error {
	_, err := r.pg.Exec(ctx, `
		INSERT INTO user_settings (
			user_id, theme, language, timezone, email_notifications,
			brand_name, custom_domain, logo_url, updated_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (user_id)
		DO UPDATE SET
			theme = EXCLUDED.theme,
			language = EXCLUDED.language,
			timezone = EXCLUDED.timezone,
			email_notifications = EXCLUDED.email_notifications,
			brand_name = EXCLUDED.brand_name,
			custom_domain = EXCLUDED.custom_domain,
			logo_url = EXCLUDED.logo_url,
			updated_at = EXCLUDED.updated_at
	`,
		s.UserID,
		s.Theme,
		s.Language,
		s.Timezone,
		s.EmailNotifications,
		s.BrandName,
		s.CustomDomain,
		s.LogoURL,
		s.UpdatedAt,
	)
	return err
}
