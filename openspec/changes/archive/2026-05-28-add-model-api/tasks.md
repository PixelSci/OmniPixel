## 1. Domain — Repository Interfaces + DTOs

- [x] 1.1 Add `ProviderRepository` interface + `ProviderResponse` DTO to `domain/provider.go`
- [x] 1.2 Add `ModelRepository` interface + `ModelFilter` + `ModelResponse` DTO to `domain/model.go`

## 2. Repository — GORM Implementations

- [x] 2.1 Create `repository/provider_repository.go` — implement `FindAll()`
- [x] 2.2 Create `repository/model_repository.go` — implement `FindAll(filter ModelFilter)` with Preload("Provider") and dynamic conditions

## 3. UseCase

- [x] 3.1 Create `usecase/provider_usecase.go` — `ListProviders() ([]domain.ProviderResponse, error)`
- [x] 3.2 Create `usecase/model_usecase.go` — `ListModels(filter domain.ModelFilter) ([]domain.ModelResponse, error)`

## 4. Controller

- [x] 4.1 Create `api/controller/provider_controller.go` — `ListProviders(c fiber.Ctx) error`
- [x] 4.2 Create `api/controller/model_controller.go` — `ListModels(c fiber.Ctx) error`，parse optional query params

## 5. Routes

- [x] 5.1 Create `api/routes/provider_routes.go` — `GET /providers` on auth group
- [x] 5.2 Create `api/routes/model_routes.go` — `GET /models` on auth group

## 6. Error Codes

- [x] 6.1 Add provider/model error codes to `internal/response/code.go` (30001, 40001)
- [x] 6.2 Add domain-error-to-APIError mappings to `internal/response/mapping.go`

## 7. Wiring

- [x] 7.1 Register new repos, usecases, controllers in `bootstrap/wire.go` (`Providers` struct + `NewProviders`)
- [x] 7.2 Register new routes in `bootstrap/app.go`
