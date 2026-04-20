package db

import (
	"embed"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	// Blank import registers the "pgx5://" scheme with golang-migrate.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// Migrate brings the database up to the latest schema. Safe to call on every
// startup — golang-migrate guards concurrent runners with a Postgres advisory
// lock, and no-op migrations are returned as migrate.ErrNoChange and swallowed.
func Migrate(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("db: migrate DSN is required")
	}
	src, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("db: build migration source: %w", err)
	}

	// migrate's pgx/v5 driver registers the "pgx5" scheme; rewrite the user-
	// facing postgres:// URL so they share the same DATABASE_URL env var.
	migrateDSN := dsn
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if strings.HasPrefix(migrateDSN, prefix) {
			migrateDSN = "pgx5://" + strings.TrimPrefix(migrateDSN, prefix)
			break
		}
	}

	m, err := migrate.NewWithSourceInstance("iofs", src, migrateDSN)
	if err != nil {
		return fmt.Errorf("db: init migrate: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db: run migrations: %w", err)
	}
	return nil
}
