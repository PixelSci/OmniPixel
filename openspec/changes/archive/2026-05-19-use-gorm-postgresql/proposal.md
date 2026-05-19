## Why

The server currently uses `pgx/v5` directly with raw SQL queries in repositories. GORM provides a higher-level ORM that reduces boilerplate, improves type safety, and enables auto-migration for schema management. GORM is already pulled in transitively — making it the primary data-access layer simplifies the stack and accelerates feature development.

## What Changes

- Replace `pgxpool.Pool` with `*gorm.DB` as the database handle throughout the application
- Add GORM model tags to domain structs (User, Conversation, Message)
- Rewrite repository implementations to use GORM query methods instead of raw SQL
- Replace hand-written connection setup in `bootstrap/db.go` with GORM's `gorm.Open`
- Enable GORM auto-migration on startup for schema management
- Remove the `db/postgres.go` type alias and the direct `pgx/v5` dependency
- **BREAKING**: The `db *pgxpool.Pool` field in repositories changes to `db *gorm.DB`

## Capabilities

### New Capabilities

- `gorm-database`: Database connection management and auto-migration via GORM, replacing the pgxpool-based connection layer

### Modified Capabilities

<!-- No existing specs to modify -->

## Impact

- `bootstrap/db.go` — rewrite to create `*gorm.DB` instead of `*pgxpool.Pool`
- `bootstrap/wire.go` — update DI graph type
- `db/postgres.go` — remove or repurpose type alias
- `domain/*.go` — add GORM struct tags
- `repository/*.go` — rewrite query methods with GORM
- `go.mod` — promote gorm.io/gorm and gorm.io/driver/postgres to direct dependencies, remove pgx/v5
