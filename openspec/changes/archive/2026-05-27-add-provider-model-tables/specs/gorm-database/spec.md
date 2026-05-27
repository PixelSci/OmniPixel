## MODIFIED Requirements

### Requirement: Auto-migration on startup

The system SHALL run GORM's `AutoMigrate` on startup for all known domain models, including `Provider` and `Model`, to ensure the database schema matches the struct definitions.

#### Scenario: Schema in sync
- **WHEN** the application starts and all tables match their GORM model definitions
- **THEN** no DDL changes are executed

#### Scenario: Schema drift detected
- **WHEN** the application starts and a model has a new column or missing table
- **THEN** GORM executes the necessary DDL to align the schema

### Requirement: Domain models include GORM struct tags

The system SHALL annotate domain structs (`User`, `Conversation`, `Message`, `Provider`, `Model`) with GORM struct tags specifying table names, column mappings, primary keys, and default values.

#### Scenario: User model has GORM tags
- **WHEN** the `User` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on the ID field and `gorm:"column:..."` on each database-mapped field

#### Scenario: Conversation model has GORM tags
- **WHEN** the `Conversation` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` and `gorm:"column:..."` tags matching the `conversations` table schema

#### Scenario: Provider model has GORM tags
- **WHEN** the `Provider` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on ID and `gorm:"column:..."` tags matching the `providers` table schema

#### Scenario: Model model has GORM tags
- **WHEN** the `Model` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on ID, `gorm:"column:..."` tags, and a foreign key association to `providers`
