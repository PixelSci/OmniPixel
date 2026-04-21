package repository

import (
	"context"
	"database/sql"

	"omni-pixel/domain"
)

// NewHealth wires PostgreSQL for liveness; when db is nil, Ping is a no-op (dev without DB).
func NewHealth(db *sql.DB) domain.HealthRepository {
	if db == nil {
		return noopHealth{}
	}
	return &sqlHealth{db: db}
}

type noopHealth struct{}

func (noopHealth) Ping(context.Context) error { return nil }

type sqlHealth struct {
	db *sql.DB
}

func (r *sqlHealth) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
