## Context

`providers` 和 `models` 表、domain struct、AutoMigrate 已在上一个 change 中完成。需要新建完整应用层来暴露查询 API。

## Goals / Non-Goals

**Goals:**
- 新增 `GET /api/v1/providers` 返回所有厂商（不含 api_key）
- 新增 `GET /api/v1/models` 返回可用模型列表，支持按厂商过滤
- 遵循现有 Clean Architecture 分层：domain → repository → usecase → controller → route

**Non-Goals:**
- 不提供创建/更新/删除接口（当前为只读查询）
- 不做分页（表数据量小，几十条以内）
- 不做缓存（作为后续优化）

## Decisions

### API 设计

```
GET /api/v1/providers           →  { providers: [...] }
GET /api/v1/models              →  { models: [...] }
GET /api/v1/models?provider_id=xxx&include_disabled=true
```

这两个接口均为 **JWT 保护**（与 conversation 接口一致），挂载在 auth group 下。

### Response 形状

**Provider** — 排除 `api_key`，通过独立的 response DTO 实现安全：

```json
{
  "providers": [
    { "id": "uuid", "name": "DeepSeek", "base_url": "https://api.deepseek.com" }
  ]
}
```

**Model** — 包含 provider 名称，方便前端直接展示：

```json
{
  "models": [
    {
      "id": "uuid",
      "provider_id": "uuid",
      "provider_name": "DeepSeek",
      "model_name": "deepseek-v4-flash",
      "is_enabled": true,
      "expire_time": null,
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

`provider_name` 通过 GORM Preload 一次查出，避免 N+1。

### Repository 方法

```go
// domain/repository interfaces
type ProviderRepository interface {
    FindAll() ([]Provider, error)
}

type ModelRepository interface {
    FindAll(filter ModelFilter) ([]Model, error)
}

type ModelFilter struct {
    ProviderID     *uuid.UUID
    IncludeDisabled bool
}
```

- `ModelFilter` 是值类型（非指针），空值表示"不过滤"
- 默认行为：只返回 `is_enabled = true` 且未过期的模型

### 错误码分配

| 范围 | 用途 |
|------|------|
| 10xxx | Account |
| 20xxx | Conversation |
| 30xxx | Provider |
| 40xxx | Model |
| 90xxx | System |

本次新增：`30001`（provider not found），`40001`（model not found）

### 文件结构

```
                                   domain/
                                   ├── provider.go         ← 已有 GORM struct
                                   ├── model.go            ← 已有 GORM struct
                                   ├── provider_repository.go  ← NEW: interface + response DTO
                                   └── model_repository.go     ← NEW: interface + filter + response DTO

                                   repository/
                                   ├── provider_repository.go  ← NEW: GORM impl
                                   └── model_repository.go     ← NEW: GORM impl

                                   usecase/
                                   ├── provider_usecase.go     ← NEW
                                   └── model_usecase.go        ← NEW

                                   api/controller/
                                   ├── provider_controller.go  ← NEW
                                   └── model_controller.go     ← NEW

                                   api/routes/
                                   ├── provider_routes.go      ← NEW
                                   └── model_routes.go         ← NEW
```

## Risks / Trade-offs

- **API Key 泄露风险**：通过 response DTO 而非 `json:"-"` tag 解决，保留 API Key 在其他场景（AI 调用）中的序列化能力
- **无缓存**：每次请求查 DB，厂商/模型数据变更频率极低 → 后续可在 usecase 层加 in-memory cache
- **model 查询走 Preload**：Provider 信息通过 GORM Preload 加载，数据量小，性能无问题
