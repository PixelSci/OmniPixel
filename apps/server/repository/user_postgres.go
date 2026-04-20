package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/model"
)

// Postgres error code for a unique-constraint violation.
const pgUniqueViolation = "23505"

// userSelectColumns uses COALESCE so nullable password_hash scans into a
// plain string ("" means no password / OAuth-only user).
const userSelectColumns = "id, name, email, COALESCE(password_hash, ''), created_at, updated_at"

func NewPostgresUserRepo(pool *pgxpool.Pool) UserRepository {
	return &postgresUserRepo{pool: pool}
}

type postgresUserRepo struct {
	pool *pgxpool.Pool
}

func (r *postgresUserRepo) Create(ctx context.Context, u *model.User) error {
	const q = `
		INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
	`
	if _, err := r.pool.Exec(ctx, q, u.ID, u.Name, u.Email, u.PasswordHash, u.CreatedAt, u.UpdatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return ErrDuplicate
		}
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *postgresUserRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	return r.findOne(ctx, `SELECT `+userSelectColumns+` FROM users WHERE id = $1`, id)
}

func (r *postgresUserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return r.findOne(ctx, `SELECT `+userSelectColumns+` FROM users WHERE email = $1`, email)
}

func (r *postgresUserRepo) findOne(ctx context.Context, query string, args ...any) (*model.User, error) {
	var u model.User
	err := r.pool.QueryRow(ctx, query, args...).
		Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query user: %w", err)
	}
	return &u, nil
}

func (r *postgresUserRepo) List(ctx context.Context, offset, limit int) ([]*model.User, int, error) {
	var total int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	if total == 0 {
		return []*model.User{}, 0, nil
	}

	const q = `
		SELECT ` + userSelectColumns + `
		FROM users
		ORDER BY created_at ASC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.pool.Query(ctx, q, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]*model.User, 0, limit)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, &u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate users: %w", err)
	}
	return users, total, nil
}
