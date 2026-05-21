## Why

当前 `GET /conversations/:conversation_id` 路由注册了一个空函数 `func() {}`，前端无法获取单条会话详情及关联消息。需要打通从路由到数据库的完整链路。

## What Changes

- 新增 `GET /conversations/:conversation_id` 端点，返回会话元数据及其消息列表
- `ConversationRepository` 接口新增 `FindByID`、`ListMessagesByConversationID` 方法
- 新增 `ConversationDetailResponse` DTO，不暴露 `UserID`
- 取消注释领域错误 `ErrConversationNotFound`、`ErrInvalidConversationID`

## Capabilities

### New Capabilities
- `conversation-detail`: 通过 conversation_id 获取会话详情及其消息列表，按 user_id 做权限隔离

### Modified Capabilities
- `gorm-database`: 新增 `FindByID`、`ListMessagesByConversationID` 两个 GORM 仓储方法

## Impact

- `domain/conversation.go` — 新增 DTO + 错误 + 接口方法签名
- `repository/conversation_repository.go` — 新增两个 GORM 查询实现
- `usecase/conversation_usecase.go` — 新增 `GetConversation` 方法
- `api/controller/conversation_controller.go` — 新增 `GetConversation` handler
- `api/routes/conversation_routes.go` — 替换空函数为实际 handler
