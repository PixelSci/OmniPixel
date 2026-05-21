## Context

当前 `GET /conversations/:conversation_id` 路由处理器为空函数。需要从 domain → repository → usecase → controller → route 五层完整实现。

该路由在 JWT 保护组内，`user_id` 已通过 middleware 注入 `fiber.Locals`。Conversation 和 Message 均为已存在的 GORM 模型，数据库表结构无需变更。

## Goals / Non-Goals

**Goals:**
- 实现 `GET /conversations/:conversation_id`，返回会话元数据 + 消息列表
- 按 `user_id` 隔离，用户只能访问自己的会话
- 过滤软删除的会话（`is_visible = false`）

**Non-Goals:**
- 不包含消息分页（后续独立实现）
- 不涉及 Conversation 的 PATCH/DELETE
- 不新建 MessageRepository（消息查询放在 ConversationRepository 内）

## Decisions

### 权限控制：统一返回 404

用户请求他人会话或不存在会话时均返回 404，不泄露会话存在性。

查询条件 `WHERE id = ? AND user_id = ? AND is_visible = true` 保证无法区分"不存在"和"无权访问"。

### Conversation.UserID 不返回

`ConversationDetailResponse` 从 `Conversation` 结构体中显式排除 `UserID`，该字段仅用于外键查询。Message.UserId 保留，前端需据此区分用户消息和 AI 消息。

### 不在 Controller 中做 DTO 转换

UseCase 直接返回 `domain.ConversationDetailResponse`，Controller 仅做 HTTP 层面的绑定和状态码映射，保持 Controller 薄层职责。

### ListMessagesByConversationID 放在 ConversationRepository

不新建 MessageRepository。当前消息查询场景仅此一处，足够简单。未来消息逻辑膨胀时再拆。

## Risks / Trade-offs

- `ListMessagesByConversationID` 无分页 → 长对话可能返回大量 Message。当前 MVP 阶段可接受，后续加分页参数和 LIMIT/OFFSET
- `is_visible` 只标记查询，无 `deleted_at` 时间戳 → 缺少审计追溯能力，但不影响当前功能
