## 1. Error response package

- [x] 1.1 创建 `code.go` — ErrorCode 常量定义（10xxx Account、20xxx Conversation、90xxx 系统）+ `HTTPStatus(code)` 映射函数
- [x] 1.2 创建 `error.go` — `APIError` 结构体（`code` + `message`）+ `New(code, msg)` 构建函数
- [x] 1.3 创建 `mapping.go` — 预构建错误常量（`ErrInvalidCredentials`、`ErrConversationNotFound` 等）+ domain error → *APIError 映射表
- [x] 1.4 创建 `write.go` — `Write(c, apiErr)` 通用写入 + `DomainError(c, err)` 查表写入

## 2. Controller 迁移

- [x] 2.1 `account_controller.go` — 将 `errors.Is` 分支替换为 `response.DomainError(c, err)`，将通用错误替换为 `response.Write(c, response.ErrXxx)`
- [x] 2.2 `conversation_controller.go` — 将 `errors.Is` 分支替换为 `response.DomainError(c, err)`，将通用错误替换为 `response.Write(c, response.ErrXxx)`

## 3. Middleware 迁移

- [x] 3.1 `jwt_middleware.go` — 将 `fiber.Map{"message": "..."}` 替换为 `response.Write(c, response.ErrXxx)`

## 4. 验证

- [x] 4.1 `go build ./...` 编译通过
