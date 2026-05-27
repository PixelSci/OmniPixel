## Context

目前 AI 调用链路中，模型列表硬编码在前端，API Key 存在 localStorage 或单个环境变量中。后端只有一个 `OpenAIProvider`，无法按厂商/模型路由。

## Goals / Non-Goals

**Goals:**
- `providers` 表存储厂商信息（名称、API Key、base_url），一个厂商一把 key
- `models` 表存储模型信息，归属厂商，控制启用/过期状态
- 后端根据模型名查找厂商，用对应 key 和 base_url 调用

**Non-Goals:**
- 不实现多协议支持（Anthropic 原生、Gemini 原生），当前仅 OpenAI 兼容协议
- 不做用户级别的 API Key（当前为系统级 key）
- 不提供厂商/模型的管理后台 UI

## Decisions

### 两张表，职责分离

```
providers                      models
┌────────────────────┐        ┌────────────────────┐
│ id        UUID PK   │──┐     │ id        UUID PK   │
│ name      string    │  │     │ provider_id UUID FK │
│ base_url  string    │  │     │ model_name string   │
│ api_key   string    │◄─┘     │ is_enabled bool     │
└────────────────────┘        │ expire_time timestamp?│
                              └────────────────────┘
```

**理由**：厂商和模型是一对多关系，分开存避免数据冗余。后续扩展（按模型定价、限流）可直接在 models 表加字段。

### API Key 存在 provider 级别

同一厂商的所有模型共享同一把 key。不设置用户级 key，保持当前系统级 key 的模式。**理由**：当前阶段不需要区分用户配额，简化设计。

### `model_name` 即 API 调用时传入的字符串

ChatRequest 中的 model_id 直接对应 `models.model_name`。不再要求 UUID 格式。

### AIProvider 初始化为工厂模式

启动时从 `providers` 表加载所有厂商，为每个厂商创建一个 `OpenAIProvider` 实例。根据请求的 model_name 查 `models` 表找到 provider_id，再路由到对应实例。

### 种子数据

providers 和 models 的初始数据用 GORM 的 `FirstOrCreate` 方式插入，不写 SQL 迁移文件。

## Risks / Trade-offs

- **数据库依赖**：每次请求需查 models 表和 providers 表 → 模型列表可缓存到内存，减少 DB 查询
- **表结构变更**：未来支持多协议时，需要给 providers 加 `protocol` 字段 → 加字段成本低，不阻塞当前设计
