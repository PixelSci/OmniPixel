## Context

The server currently uses `pgx/v5` (`pgxpool.Pool`) for PostgreSQL connectivity. Repositories hand-write SQL queries with manual scanning. GORM (`gorm.io/gorm` v1.31.1 and `gorm.io/driver/postgres` v1.6.0) already exists in the dependency tree as indirect dependencies. The domain models (User, Conversation, Message) need GORM struct tags to participate in the ORM. The `bootstrap` package constructs the database handle and injects it into repositories via a providers struct.

## Goals / Non-Goals

**Goals:**
- Replace `*pgxpool.Pool` with `*gorm.DB` across all layers
- Add GORM model annotations to domain structs
- Rewrite repository methods to use GORM query builders
- Enable GORM auto-migration on startup for schema management
- Remove `pgx/v5` from direct dependencies

**Non-Goals:**
- Changing the database schema (tables, columns, indexes)
- Adding new repository methods or domain logic
- Introducing GORM hooks, callbacks, or relationships beyond basic queries
- Changing the application's startup flow or configuration format

## Decisions

### Decision 1: Use `gorm.io/driver/postgres` with the existing `pgx` driver underneath

GORM's postgres driver (`gorm.io/driver/postgres` v1.6.0) already uses `pgx/v5` internally as the SQL driver. This means we keep the same wire protocol and connection semantics — only the query-building layer changes.

### Decision 2: Keep domain structs in `domain/` package, add GORM tags directly

The `domain.User`, `domain.Conversation`, and `domain.Message` structs get GORM struct tags (`gorm:"column:..."`) alongside existing `json` tags. This avoids creating separate ORM model types while keeping the domain layer aware of persistence details — acceptable for a single-service app with no multi-DB abstraction.

### Decision 3: Auto-migrate in `bootstrap/db.go` after connecting

GORM's `AutoMigrate` runs on startup to sync struct definitions with the database schema. This is safe for development and simple deployments. The list of models passed to `AutoMigrate` lives in the db bootstrap function.

### Decision 4: Leverage GORM query methods in repositories

Replace raw SQL + manual `Scan` with GORM's chainable query API (`db.Where().First()`, `db.Find()`, `db.Model().Updates()`). This reduces boilerplate, eliminates manual column listing, and improves readability.

### Decision 5: Remove `db/postgres.go` type alias

The `type Postgres = pgxpool.Pool` type alias is no longer needed. Replace with a `db` package reference or remove the indirection entirely.

## Risks / Trade-offs

- **ORM magic hides query behavior** → Mitigation: GORM logging can be enabled in development for query visibility. The repository API remains thin so switching back to raw SQL is mechanical.
- **AutoMigrate on every startup** → Mitigation: `AutoMigrate` is idempotent and only runs DDL when schema differs. For production, migrations can be extracted later.
- **GORM tags couple domain to persistence** → Mitigation: This is a single-service app; the coupling is acceptable. If multi-DB support is needed later, separate model types can be introduced.
