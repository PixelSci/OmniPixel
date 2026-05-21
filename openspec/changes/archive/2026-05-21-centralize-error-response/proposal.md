## Why

当前所有错误信息以字符串硬编码在各 controller 和 middleware 中，共 11+ 处 `"message"` 散落 3 个文件。缺少统一的错误码体系、构建方法和 domain error 到 HTTP 响应的映射机制，导致重复代码和新增错误时的一致性风险。

## What Changes

- 新增 `internal/response/` 包，集中管理错误码定义、APIError 结构、预构建错误常量、domain error 映射表及 Fiber 写入函数
- 按模块分错误码：10xxx Account、20xxx Conversation、90xxx 系统通用
- `HTTPStatus(code)` 方法将业务错误码映射为标准 HTTP 状态码
- `New(code, message)` 构建 `*APIError`
- `DomainError(c, err)` 自动查表将 domain error 转为 HTTP 响应
- `Write(c, apiErr)` 通用错误写入
- 迁移所有 controller 和 middleware 中的硬编码错误

## Capabilities

### New Capabilities
- `error-response-system`: 集中的错误码体系、APIError 构建、domain error 映射及 HTTP 响应输出

## Impact

- `internal/response/code.go` — ErrorCode 常量定义
- `internal/response/error.go` — APIError 结构体 + New() 构建
- `internal/response/mapping.go` — HTTPStatus() 映射 + domain error 映射表
- `internal/response/write.go` — Write() / DomainError()
- `api/controller/account_controller.go` — 替换硬编码
- `api/controller/conversation_controller.go` — 替换硬编码
- `api/middleware/jwt_middleware.go` — 替换硬编码
