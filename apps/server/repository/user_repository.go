package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/domain"
)

type UserRepository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewUserRepository(db *pgxpool.Pool, timeout time.Duration) *UserRepository {
	return &UserRepository{db: db, timeout: timeout}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		SELECT id, username, email, password_hash, device_fingerprint, display_name,
			avatar_url, status, email_verified_at, last_login_at, created_at, updated_at
		FROM users
		WHERE lower(email) = lower($1)
		LIMIT 1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.DeviceFingerprint,
		&user.DisplayName,
		&user.AvatarURL,
		&user.Status,
		&user.EmailVerifiedAt,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	_, err := r.db.Exec(ctx, `UPDATE users SET last_login_at = now(), updated_at = now() WHERE id = $1`, userID)
	return err
}
