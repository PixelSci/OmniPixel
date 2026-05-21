## 1. Domain layer

- [x] 1.1 取消注释 `ErrConversationNotFound` 和 `ErrInvalidConversationID` 错误定义
- [x] 1.2 新增 `ConversationDetailResponse` 结构体（Conversation 字段排除 UserID，附带 Messages 切片）
- [x] 1.3 在 `ConversationRepository` 接口中新增 `FindByID(conversationID, userID uuid.UUID) (*Conversation, error)` 和 `ListMessagesByConversationID(conversationID uuid.UUID) ([]Message, error)` 方法签名

## 2. Repository layer

- [x] 2.1 实现 `FindByID` — GORM 查询 `WHERE id = ? AND user_id = ? AND is_visible = ?`，未找到返回 `ErrConversationNotFound`
- [x] 2.2 实现 `ListMessagesByConversationID` — GORM 查询 `WHERE conversation_id = ? ORDER BY created_at ASC`

## 3. UseCase layer

- [x] 3.1 新增 `GetConversation(userID, conversationID uuid.UUID) (*domain.ConversationDetailResponse, error)` 方法，依次调用 `FindByID` 和 `ListMessagesByConversationID`，组装返回 DTO

## 4. Controller layer

- [x] 4.1 新增 `GetConversation(c fiber.Ctx) error` handler：解析 `conversation_id` 参数为 UUID（失败 → 400），获取 `user_id` 从 Locals，调用 usecase，映射领域错误到 HTTP 状态码（`ErrInvalidConversationID → 400`，`ErrConversationNotFound → 404`，其他 → 500）

## 5. Route layer

- [x] 5.1 将 `GET /conversations/:conversation_id` 的空函数替换为 `conversationController.GetConversation`
