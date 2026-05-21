## Context

当前 Controller 和 Middleware 中的错误响应使用硬编码字符串，每处都手动写 `c.Status(code).JSON(fiber.Map{"message": msg})`。domain 层定义的领域错误（`ErrConversationNotFound` 等）在 controller 中通过 `errors.Is` 手动映射到 HTTP 状态码和消息。

需要建立集中的错误码体系、构建方法、映射表和输出函数。

## Goals / Non-Goals

**Goals:**
- 按模块分段定义数值型错误码（10xxx Account、20xxx Conversation、90xxx 系统通用）
- `APIError` 结构体携带 `code` + `message`，通过 `New()` 构建
- `HTTPStatus(code)` 将业务错误码映射为标准 HTTP 状态码
- `DomainError(c, err)` 将 domain error 自动映射为 HTTP 响应
- `Write(c, apiErr)` 通用写入方法
- 迁移所有现有硬编码到新体系
- 响应 JSON body 格式统一为 `{"code": 10001, "message": "邮箱或密码错误"}`

**Non-Goals:**
- 不改变 domain 层错误定义（`ErrInvalidCredentials` 等保留不变）
- 不在 domain 层引入 HTTP 概念
- 不处理 validation errors（字段级校验）——后续单独做

## Decisions

### 错误码编码规则

```
模块号 * 1000 + 具体序号

10xxx  Account
20xxx  Conversation
30xxx  Message（预留）
90xxx  系统/通用
```

每个业务错误码唯一标识一类错误场景，前端通过 code 判断分支。

### HTTP 状态码从 ErrorCode 推导

`code / 1000 * 1000` 映射到 HTTP 状态码区间，由 `HTTPStatus()` 函数集中维护映射。不在 `APIError` 结构体中存储 HTTP 状态码，保持其与传输层解耦。

```go
10001 → 401
20001 → 400
20002 → 404
90000 → 500
```

### APIError 实现 error 接口

`APIError` 实现 `Error() string`，返回 message。这样在 usecase 层也可以使用 `APIError` 作为通用错误返回，但 domain 层仍使用纯领域错误。

### DomainError() 不替代所有错误码

`DomainError(c, err)` 仅处理 domain error 映射。对于非 domain 场景（uuid 解析失败、请求体绑定失败），使用 `Write(c, apiErr)` + 预构建常量（如 `ErrInvalidConvID`）。

### 文件拆分

拆为 4 个文件而非单文件，职责清晰：

| 文件 | 职责 |
|---|---|
| `code.go` | ErrorCode 常量 + `HTTPStatus()` |
| `error.go` | APIError + `New()` |
| `mapping.go` | domain error → *APIError 映射表 |
| `write.go` | Write() / DomainError() Fiber 输出 |

## Risks / Trade-offs

- `HTTPStatus()` 需要为每个 ErrorCode 维护映射 → 新增错误码时容易遗漏，但编译期不会报错。后续可通过测试覆盖验证
- 响应格式从 `{"message": "..."}` 变为 `{"code": N, "message": "..."}` → 前端需要适配，当前 MVP 阶段无前端依赖，可直接改
