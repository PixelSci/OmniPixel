## Why

当前 AI 厂商和模型信息全部硬编码在前端 `useModelSettings.ts` 中，API Key 存在前端 localStorage。后端只有一个全局 `AI_API_KEY`，无法支持多厂商、多模型的管理和灵活切换。增加厂商和模型的两张数据库表，实现集中管理和动态配置。

## What Changes

- 新增 `providers` 表，存储 AI 厂商信息（名称、API Key、base_url）
- 新增 `models` 表，存储模型信息（归属厂商、模型名、启用状态、过期时间）
- 新增对应的 Go domain model 和 GORM AutoMigrate
- `ChatRequest.ModelID` 改为参考 models 表中的模型名来调用对应厂商
- 移除前端 localStorage 中的 API Key 管理逻辑，改为从后端获取模型列表

## Capabilities

### New Capabilities
- `provider-management`: 厂商信息的增删改查，存储 API Key 和接入地址
- `model-management`: 模型信息管理，归属厂商，控制启用状态和有效期

### Modified Capabilities
- `gorm-database`: AutoMigrate 注册新增两张表

## Impact

- `apps/server/domain/` — 新增 `provider.go`，修改 `conversation.go`（Message.ModelID 字段调整）
- `apps/server/bootstrap/db.go` — AutoMigrate 注册
- `apps/server/bootstrap/wire.go` — AIProvider 工厂模式改造
- `apps/omni-pixel/` — 模型列表改为从后端 API 获取，移除 localStorage key 管理
