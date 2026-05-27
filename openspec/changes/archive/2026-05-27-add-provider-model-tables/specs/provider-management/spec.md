## ADDED Requirements

### Requirement: Provider domain model
The system SHALL define a `Provider` struct with GORM tags mapping to the `providers` table, containing fields: `ID` (UUID PK), `Name` (string), `BaseURL` (string), `APIKey` (string).

#### Scenario: Provider struct compiles
- **WHEN** the `Provider` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on ID and `gorm:"column:..."` tags on each database-mapped field

### Requirement: Provider repository
The system SHALL provide a `ProviderRepository` with a method `FindByID(id uuid.UUID) (*Provider, error)` that looks up a provider by primary key.

#### Scenario: Find provider by ID
- **WHEN** `FindByID` is called with an existing provider UUID
- **THEN** the repository returns the `Provider` struct

#### Scenario: Provider not found
- **WHEN** `FindByID` is called with a non-existent UUID
- **THEN** the repository returns a not-found error

### Requirement: Provider seed data
The system SHALL seed the `providers` table with initial data at startup using `FirstOrCreate`, including at minimum a DeepSeek provider entry from existing `AI_BASE_URL` and `AI_API_KEY` env values.

#### Scenario: First startup
- **WHEN** the application starts and the `providers` table is empty
- **THEN** the DeepSeek seed record is inserted

#### Scenario: Subsequent startup
- **WHEN** the application starts and the seed record already exists
- **THEN** no duplicate records are created
