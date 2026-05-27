## 1. Domain Models

- [ ] 1.1 Create `domain/provider.go` — `Provider` struct with GORM tags (id, name, base_url, api_key)
- [ ] 1.2 Create `domain/model.go` — `Model` struct with GORM tags (id, provider_id FK, model_name, is_enabled, expire_time)

## 2. Seed Data

- [ ] 2.1 Add `Provider` and `Model` to `AutoMigrate` in `bootstrap/db.go`
- [ ] 2.2 Seed default DeepSeek provider from `AI_BASE_URL` / `AI_API_KEY` env values using `FirstOrCreate`

## 3. Repositories

- [ ] 3.1 Create `repository/provider_repository.go` — `FindByID(id uuid.UUID) (*domain.Provider, error)`
- [ ] 3.2 Create `repository/model_repository.go` — `FindByModelName(name string) (*domain.Model, error)` and `ListEnabled() ([]domain.Model, error)`

## 4. AI Provider Routing

- [ ] 4.1 Modify `usecase/conversation_usecase.go` Chat method — look up model by name, then provider, use provider's key and base_url
- [ ] 4.2 Update `AIProvider` interface or factory to accept per-request api_key/base_url instead of hardcoded env values
- [ ] 4.3 Fix `Message.ModelID` field in `domain/conversation.go` — change from `uuid.UUID` to `string` to store model_name
