## Context

`POST /conversation` 当前为 `func() {}` 空函数。需要实现为 SSE 流式聊天接口：`conversation_id` 为 null 时新建 Conversation（title 取首条消息前 60 字），否则校验归属后继续对话。用户消息存 DB，调 AI stream，AI 回复流完后存 DB。

## Goals / Non-Goals

**Goals:**
- `POST /conversation` 实现 SSE 流式聊天
- `StreamWriter` 接口隔离 UseCase 与 HTTP 层
- `AIProvider` 接口抽象 AI 调用（兼容 OpenAI）
- `ConversationRepository` 新增 `Insert`、`InsertMessage`

**Non-Goals:**
- 不新建 Model 表（`model_id` 字段暂不校验）
- 不实现历史 session 管理（上下文窗口由调用方控制）
- 不在 domain 层引入 Fiber 依赖

## Decisions

### StreamWriter 放在 domain 层

遵循与 Repository 接口相同的依赖反转模式：

```
domain/StreamWriter        ← UseCase 依赖
api/controller/sseAdapter  ← Fiber 适配实现
```

UseCase 完全不感知 Fiber 和 SSE 协议细节。`sseAdapter` 封装 `fiber.Ctx` 实现 `StreamWriter`。

### AIProvider 放在 domain 层

```go
type AIChatMessage struct {
    Role    string // "system" | "user" | "assistant"
    Content string
}

type AIStreamChunk struct {
    Token string
    Done  bool
}

type AIProvider interface {
    ChatStream(messages []AIChatMessage, modelID string) (<-chan AIStreamChunk, error)
}
```

`modelID` 传 string（UUID 的 string 形式）。内部实现去查配置或直接传给 OpenAI API。

### 标题生成

新建 Conversation 时，title 取首条消息前 60 字：`cutString(message, 60)`。后面不改 title。

### 流完成后存储

AI 回复所有 token 拼接成完整内容后一次写入：`InsertMessage(aiMessage)`，不在流过程中逐条写 DB。

### SSE 格式

```
data: {"token":"你"}
data: {"token":"好"}
data: {"done":true,"conversation_id":"...","message_id":"..."}
```

### Controller 通过 goroutine 桥接 StreamWriter

Controller 创建 `sseAdapter{fiber.Ctx}`，UseCase 接收 `StreamWriter`，UseCase 内部阻塞式调用 AI stream，每收到一个 chunk 就 `WriteToken`。Controller 启动 goroutine 执行 UseCase，主 goroutine 等待完成。

实际上 Fiber 3 已经支持 streaming，Controller 直接传 writer 给 UseCase 同步调用即可，因为 Fiber 的 handler 本身就是阻塞式的，SSE 逐条写入不需要额外 goroutine。

## Risks / Trade-offs

- AI API 调用失败时流已开启 → 通过 `data: {"error": "..."}` 写入流结束
- OpenAI 实现硬编码 API key 和 endpoint → 后续通过 env 配置注入
- `modelID` 不作校验 → 无效 model 由 AI API 报错，前端传正确的 UUID 字符串
