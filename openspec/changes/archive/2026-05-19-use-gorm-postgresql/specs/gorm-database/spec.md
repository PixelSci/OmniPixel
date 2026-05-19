## ADDED Requirements

### Requirement: GORM database connection

The system SHALL establish a PostgreSQL database connection using `gorm.io/driver/postgres` with the existing environment configuration (`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `DB_NAME`).

#### Scenario: Successful connection
- **WHEN** the application starts with valid database credentials
- **THEN** a `*gorm.DB` instance is created and the connection pool is ready

#### Scenario: Invalid credentials
- **WHEN** the application starts with invalid database credentials
- **THEN** the application logs a fatal error and exits

### Requirement: Auto-migration on startup

The system SHALL run GORM's `AutoMigrate` on startup for all known domain models to ensure the database schema matches the struct definitions.

#### Scenario: Schema in sync
- **WHEN** the application starts and all tables match their GORM model definitions
- **THEN** no DDL changes are executed

#### Scenario: Schema drift detected
- **WHEN** the application starts and a model has a new column or missing table
- **THEN** GORM executes the necessary DDL to align the schema

### Requirement: Repository methods use GORM query API

The system SHALL use GORM's query builder methods instead of raw SQL for all repository operations.

#### Scenario: Find user by email
- **WHEN** `UserRepository.FindByEmail` is called with an email address
- **THEN** the repository uses `db.Where("lower(email) = lower(?)", email).First(&user)` and returns the user or a not-found error

#### Scenario: Update user's last login
- **WHEN** `UserRepository.UpdateLastLogin` is called with a user ID
- **THEN** the repository uses `db.Model(&domain.User{}).Where("id = ?", userID).Updates(...)` to set `last_login_at` and `updated_at`

#### Scenario: List conversations by user ID
- **WHEN** `ConversationRepository.ListByUserID` is called with a user UUID
- **THEN** the repository uses `db.Where("user_id = ?", userID).Order("...").Find(&conversations)` and returns the slice

### Requirement: Domain models include GORM struct tags

The system SHALL annotate domain structs (`User`, `Conversation`, `Message`) with GORM struct tags specifying table names, column mappings, primary keys, and default values.

#### Scenario: User model has GORM tags
- **WHEN** the `User` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on the ID field and `gorm:"column:..."` on each database-mapped field

#### Scenario: Conversation model has GORM tags
- **WHEN** the `Conversation` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` and `gorm:"column:..."` tags matching the `conversations` table schema

### Requirement: Remove pgxpool dependency from application code

The system SHALL no longer reference `github.com/jackc/pgx/v5/pgxpool` in application code. The `db/postgres.go` type alias SHALL be removed or repurposed.

#### Scenario: Application compiles without pgxpool imports
- **WHEN** the application compiles
- **THEN** no source file imports `pgxpool` directly — only `gorm.io/gorm` is referenced for database access
