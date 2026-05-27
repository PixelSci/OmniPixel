## 1. Domain layer

- [x] 1.1 `domain/conversation.go` — 新增 `AIChatMessage`、`AIStreamChunk`、`ChatRequest` 结构体，`StreamWriter` 接口，`AIProvider` 接口
- [x] 1.2 `ConversationRepository` 接口新增 `Insert(*Conversation) error`、`InsertMessage(*Message) error`

## 2. AI Provider

- [x] 2.1 `internal/ai/openai.go` — OpenAI 兼容实现：HTTP POST + SSE reader → channel

## 3. Repository layer

- [x] 3.1 `repository/conversation_repository.go` — 实现 `Insert`（`db.Create`）、`InsertMessage`（`db.Create`）

## 4. UseCase layer

- [x] 4.1 `usecase/conversation_usecase.go` — 新增 `Chat` 方法：新建/校验会话 → 存用户消息 → 查历史 → 调 AI → SSE → 存 AI 回复

## 5. Controller layer

- [x] 5.1 `api/controller/conversation_controller.go` — 新增 `sseAdapter` 实现 `domain.StreamWriter`，新增 `Chat` handler 绑定请求并调用 usecase

## 6. Route + Bootstrap

- [x] 6.1 `api/routes/conversation_routes.go` — 替换 `POST /conversation` stub 为 `Chat` handler
- [x] 6.2 `bootstrap/wire.go` — 注入 OpenAIProvider；`env.go` 新增 `AI_BASE_URL` / `AI_API_KEY`

## 7. 验证

- [x] 7.1 `go build ./...` 编译通过
