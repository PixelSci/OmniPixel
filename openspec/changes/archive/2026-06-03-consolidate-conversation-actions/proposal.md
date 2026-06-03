## Why

当前 `DELETE /conversations/:id` 只能删除，没有修改标题的接口。用一个 `POST /conversations/:id` 统一处理两种操作，通过 body 中的 `action` 区分。

## What Changes

- **路由** — 删除 `DELETE /conversations/:conversation_id`，新增 `POST /conversations/:conversation_id`
- **Controller** — `ConversationAction` handler：解析 `action`，分发到 update title / delete
- **UseCase** — `UpdateTitle(conversationID, userID, title)` 方法
- **Repository** — `Update` 方法（更新 title）
- **Domain** — `ConversationActionRequest` struct
- **Frontend** — `deleteConversation` 改为 POST 请求

## Impact

- `apps/server/api/routes/conversation_routes.go`
- `apps/server/api/controller/conversation_controller.go`
- `apps/server/usecase/conversation_usecase.go`
- `apps/server/repository/conversation_repository.go`
- `apps/server/domain/conversation.go`
- `apps/omni-pixel/src/lib/conversation.ts`
