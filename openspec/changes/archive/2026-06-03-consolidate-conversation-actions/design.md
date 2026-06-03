## Context

用 `POST /conversations/:conversation_id` 统一处理修改标题和删除操作，通过 `action` 字段区分。

## Goals / Non-Goals

**Goals:**
- `POST /conversations/:conversation_id` 支持 `update` 和 `delete` 两种 action
- 删除原有 `DELETE` 路由
- 前端 `deleteConversation` 改为 POST 调用

**Non-Goals:**
- 不在同一个请求里同时做多个操作

## Decisions

### 请求格式

```json
// 修改标题
POST /conversations/:conversation_id
{ "action": "update", "title": "新标题" }

// 删除
POST /conversations/:conversation_id
{ "action": "delete" }
```

### Domain

```go
type ConversationActionRequest struct {
    Action string `json:"action"` // "update" | "delete"
    Title  string `json:"title"`
}
```

### Controller 分发

```
ConversationAction(c)
    ├─ action == "update" → useCase.UpdateTitle(id, userID, title)
    ├─ action == "delete"       → useCase.DeleteConversation(id, userID)
    └─ 其他                     → 400
```

### Repository

新增 `Update(conversationID, userID, title)` — 只更新 title 字段，WHERE 带上 user_id 校验归属。

### 前端

```typescript
// 改标题
http.post(`/conversations/${id}`, { action: 'update', title: '新标题' })

// 删除
http.post(`/conversations/${id}`, { action: 'delete' })
```

### 前端 UI

`HigChatSessionList.vue` 改造：

```
hover 时显示 ... 按钮 → 点击弹出下拉菜单
                          ├─ ✏️ 修改标题
                          └─ 🗑️ 删除对话
```

- `Trash2` icon → `Ellipsis` icon
- 用 `UiDropdownMenu` 包裹菜单
- 删除 `confirm()` 弹窗
- 新增 `emit('rename', id)` — "修改标题"触发，父组件弹出输入框
- "删除对话" → `emit('delete', id)`

## Risks

- 前端 `deleteConversation` 函数签名不变，只是内部实现从 `http.delete` 改为 `http.post`
