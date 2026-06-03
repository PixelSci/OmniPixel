## 1. Domain

- [x] 1.1 `domain/conversation.go` — 新增 `ConversationActionRequest` struct

## 2. Repository

- [x] 2.1 `repository/conversation_repository.go` — 新增 `UpdateTitle(conversationID, userID, title)` 方法

## 3. UseCase

- [x] 3.1 `usecase/conversation_usecase.go` — 新增 `UpdateConversationTitle` 方法

## 4. Controller & Routes

- [x] 4.1 `controller/conversation_controller.go` — 新增 `ConversationAction` handler，按 action 分发
- [x] 4.2 `routes/conversation_routes.go` — 删除 `DELETE`，加 `POST /conversations/:conversation_id`

## 5. Frontend API

- [x] 5.1 `lib/conversation.ts` — `deleteConversation` 改为 `http.post`；新增 `updateConversationTitle(id, title)` 函数

## 6. Frontend UI

- [x] 6.1 `components/HigChatSessionList.vue` — `Trash2` → `Ellipsis`；`...` 按钮 + `UiDropdownMenu` 菜单（修改标题 / 删除对话）；新增 `rename` emit；移除 `confirm()`
- [x] 6.2 `pages/index.vue` — 处理 `@rename` 事件（弹出 prompt 输入新标题 → 调 `updateConversationTitle`）

## 7. Verify

- [x] 7.1 `go build ./...` + `pnpm build` 编译通过
