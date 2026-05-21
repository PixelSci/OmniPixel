## Why

`POST /conversation` 当前为空函数 stub，需要实现为 SSE 流式聊天接口，支持 conversation 自动创建或继续对话，并通过 AI Provider 抽象支持扩展。

## What Changes

- `POST /conversation` 实现 SSE 流式聊天（`conversation_id` 为 null 时自动新建 Conversation）
- `domain` 新增 `ChatRequest`、`StreamWriter` 接口、`AIProvider` 接口
- `ConversationRepository` 新增 `Insert`、`InsertMessage` 方法
- `internal/ai/` 新增 AI Provider 抽象接口 + OpenAI 兼容实现
- UseCase 通过 `StreamWriter` 接口写流，不依赖 Fiber

## Capabilities

### New Capabilities
- `conversation-chat`: SSE 流式聊天端点，自动创建 Conversation + 流式返回 AI 回复

### Modified Capabilities
- `gorm-database`: `ConversationRepository` 新增 `Insert`、`InsertMessage` 方法
- `conversation-detail`: `ConversationRepository` 接口新增方法签名

## Impact

- `domain/conversation.go` — ChatRequest、StreamWriter、AIProvider、Repository 接口新增
- `internal/ai/provider.go` — AI Provider 抽象接口
- `internal/ai/openai.go` — OpenAI 兼容实现
- `repository/conversation_repository.go` — Insert、InsertMessage
- `usecase/conversation_usecase.go` — Chat 方法
- `api/controller/conversation_controller.go` — Chat handler + sseAdapter
- `api/routes/conversation_routes.go` — 替换 stub
