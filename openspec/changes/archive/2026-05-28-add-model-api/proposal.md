## Why

`providers` 和 `models` 两张表已经建好并通过 AutoMigrate 注册，但目前没有任何 API 可以查询这些数据。前端无法获取可用的厂商和模型列表（例如设置页面需要展示可用模型供用户选择），必须先 hardcode 或绕过。

## What Changes

- 新增 `GET /api/v1/providers` — 查询所有 AI 厂商列表（不含 api_key）
- 新增 `GET /api/v1/models` — 查询可用模型列表，支持按 `provider_id` 过滤
- 新建完整的应用层：repository 接口与实现、usecase、controller、route、wire
- 为 provider 和 model 添加独立的错误码（40xxx / 50xxx）

## Capabilities

### Modified Capabilities
- `provider-management`: 新增查询接口和 response DTO，扩展 repository 方法
- `model-management`: 新增查询接口和 response DTO，扩展 repository 方法

## Impact

- `apps/server/domain/` — 新增 repository 接口 + response DTO
- `apps/server/repository/` — 新增 `provider_repository.go`、`model_repository.go`
- `apps/server/usecase/` — 新增 `provider_usecase.go`、`model_usecase.go`
- `apps/server/api/controller/` — 新增 `provider_controller.go`、`model_controller.go`
- `apps/server/api/routes/` — 新增 `provider_routes.go`、`model_routes.go`
- `apps/server/bootstrap/` — wire.go 和 app.go 注册新组件
- `apps/server/internal/response/` — 新增错误码
