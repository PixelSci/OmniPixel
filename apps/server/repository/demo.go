package repository

import (
	"context"
	"database/sql"

	"omni-pixel/domain"
)

// NewDemo returns a demo repository: PostgreSQL when db is set, otherwise a fixed in-memory reply.
func NewDemo(db *sql.DB) domain.DemoRepository {
	if db == nil {
		return noopDemo{}
	}
	return &sqlDemo{db: db}
}

type noopDemo struct{}

func (noopDemo) Greeting(context.Context) (string, error) {
	return "hello (no database configured)", nil
}

type sqlDemo struct {
	db *sql.DB
}

func (r *sqlDemo) Greeting(ctx context.Context) (string, error) {
	var msg string
	err := r.db.QueryRowContext(ctx, `SELECT 'hello from PostgreSQL'`).Scan(&msg)
	return msg, err
}
