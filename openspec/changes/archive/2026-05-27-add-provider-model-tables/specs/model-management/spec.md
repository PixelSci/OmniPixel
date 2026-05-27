## ADDED Requirements

### Requirement: Model domain model
The system SHALL define a `Model` struct with GORM tags mapping to the `models` table, containing fields: `ID` (UUID PK), `ProviderID` (UUID FK), `ModelName` (string), `IsEnabled` (bool, default true), `ExpireTime` (*time.Time, nullable).

#### Scenario: Model struct compiles
- **WHEN** the `Model` struct is compiled
- **THEN** it includes `gorm:"primaryKey"` on ID, `gorm:"column:..."` tags, and a foreign key association to `providers`

### Requirement: Model repository
The system SHALL provide a `ModelRepository` with:
- `FindByModelName(modelName string) (*Model, error)` — find an enabled, non-expired model by name
- `ListEnabled() ([]Model, error)` — list all enabled, non-expired models

#### Scenario: Find model by name
- **WHEN** `FindByModelName("deepseek-v4-flash")` is called and the model exists, is enabled, and not expired
- **THEN** the repository returns the `Model` struct with its `ProviderID`

#### Scenario: Model not found or disabled
- **WHEN** the model does not exist, is disabled, or is expired
- **THEN** the repository returns a not-found error

#### Scenario: List enabled models
- **WHEN** `ListEnabled` is called
- **THEN** the repository returns all models where `is_enabled = true` and `expire_time IS NULL OR expire_time > now()`

### Requirement: Model seed data
The system SHALL seed the `models` table with initial data at startup using `FirstOrCreate`, including common DeepSeek models (deepseek-chat, deepseek-v4-flash) linked to the DeepSeek provider.

#### Scenario: First startup
- **WHEN** the application starts and the `models` table is empty
- **THEN** seed records are inserted for the DeepSeek provider

### Requirement: Chat uses models table for routing
The system SHALL look up the model by name from the `models` table, then find the corresponding provider to obtain the API key and base URL for the AI call.

#### Scenario: Chat request with valid model
- **WHEN** a chat request arrives with `model_id = "deepseek-v4-flash"`
- **THEN** the usecase looks up the model by name, retrieves its provider, and uses that provider's API key and base URL for the AI call

#### Scenario: Chat request with unknown model
- **WHEN** a chat request arrives with a `model_id` that does not match any enabled model
- **THEN** the usecase returns an error indicating the model is not available
