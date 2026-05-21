# OmniPixel Server 架构

## 概述

`apps/server` 是 OmniPixel 的后端 API 服务，采用 **Clean Architecture**（整洁架构）分层设计，以 Fiber v3 作为 HTTP 框架，GORM v2 作为 ORM，PostgreSQL 作为数据库。

技术栈：Go 1.26 + Fiber v3 + GORM v2 + PostgreSQL + JWT

## 目录结构与职责

| 文件夹 | 层级 | 工作范围 | 关键文件 |
|---|---|---|---|
| `cmd/` | 入口 | 应用启动入口，仅调用 bootstrap 层。不含业务逻辑 | `main.go` — 调用 `bootstrap.NewApp()` 并启动 |
| `bootstrap/` | 启动引导 | 应用组装中枢：环境配置加载、数据库连接、手动依赖注入、Fiber 实例创建与路由挂载 | `env.go` — Viper 读取 `.env` → `Env` 结构体<br>`db.go` — GORM PostgreSQL 连接 + `AutoMigrate`<br>`wire.go` — 按序构建依赖图，产出 `Providers`<br>`app.go` — Fiber 实例、CORS、路由分组、`Listen()` |
| `domain/` | 领域层 | 核心业务模型（GORM 标签）、请求/响应 DTO、Repository 接口、领域错误。**不依赖任何外部框架** | `account.go` — `User` 模型、`SigninRequest/Response`、`UserRepository` 接口、`ErrInvalidCredentials` / `ErrUserInactive`<br>`conversation.go` — `Conversation` / `Message` 模型、`ConversationRepository` 接口<br>`jwt.go` — `JwtCustomClaims`（嵌入 `RegisteredClaims`） |
| `usecase/` | 用例层 | 业务逻辑编排。持有 repository 接口引用，不感知 HTTP 与数据库细节 | `health_usecase.go` — 健康检查（无外部依赖）<br>`account_usecase.go` — 登录/注册逻辑（依赖 `UserRepository` + JWT 密钥）<br>`conversation_usecase.go` — 会话列表查询（依赖 `ConversationRepository`） |
| `repository/` | 仓储层 | 实现 `domain` 中定义的 Repository 接口，封装 GORM 查询 | `user_repository.go` — 实现 `UserRepository`（`FindByEmail` 大小写不敏感、`UpdateLastLogin`）<br>`conversation_repository.go` — 实现 `ConversationRepository`（`ListByUserID`） |
| `api/controller/` | 控制器 | HTTP 请求处理：绑定请求体 → 调用 usecase → 领域错误到 HTTP 状态码映射 → JSON 响应 | `health_controller.go` — `GET /health`<br>`account_controller.go` — `POST /signin`、`POST /signup`（`ErrInvalidCredentials → 401`，`ErrUserInactive → 403`）<br>`conversation_controller.go` — `GET /conversations`（从 JWT Locals 取 `user_id`） |
| `api/middleware/` | 中间件 | 全局 CORS 配置、JWT Bearer Token 校验 | `cors_middleware.go` — 允许 `localhost:5173` / `localhost:3000` 跨域<br>`jwt_middleware.go` — 解析 Bearer Token → 验证 HS256 → `fiber.Locals` 注入 `user_id` / `email` |
| `api/routes/` | 路由注册 | 接收 `fiber.Router` + controller，注册具体路径。区分公开路由与受保护路由 | `health_route.go` — `GET /health`（公开）<br>`account_routes.go` — `POST /signin`、`POST /signup`（公开）<br>`conversation_routes.go` — `GET /conversations`、`POST /conversation`、`GET /conversations/:id`（需 JWT） |
| `internal/` | 内部工具 | 应用级共享工具函数 | `token.go` — `CreateAccessToken()`（HS256 生成 JWT） |

**依赖构建顺序**：`Env → DB → Repository → UseCase → Controller → Route`

## 系统流程

### 请求生命周期

```
HTTP Request
  │
  ▼
Fiber App (bootstrap/app.go)
  │
  ├── CORS Middleware (全局)
  ├── JWT Middleware (受保护路由组)
  │     └── 解析 Bearer Token → 验证签名 → Locals 注入 user_id / email
  ▼
Controller (api/controller/)
  │
  ├── 绑定请求体 (c.Bind().Body)
  ├── 调用 UseCase
  ├── 错误 → HTTP 状态码映射
  └── 成功 → JSON 响应
  ▼
UseCase (usecase/)
  │
  ├── 业务逻辑编排
  ├── 调用 Repository 接口
  └── 返回领域对象或错误
  ▼
Repository (repository/)
  │
  ├── GORM 查询
  └── PostgreSQL
```

### 启动流程

```
main.go
  → bootstrap.NewApp()
      → fiber.New()                  // 创建 Fiber 实例
      → CORSMiddleware()             // 全局 CORS
      → NewProviders()               // 构建依赖图
          → NewEnv()                 // .env → Env 结构体
          → NewPostgresDB(env)       // GORM 连接 + AutoMigrate
          → NewUserRepository(db)
          → NewConversationRepository(db)
          → NewHealthUseCase()
          → NewAccountUseCase(repo, secret, expiry)
          → NewConversationUseCase(repo)
          → New*Controller(usecase)
      → 注册路由
          → /api/v1 (公开组)
          → /api/v1 + JWT Middleware (受保护组)
  → b.Start()
      → app.Listen(SERVER_ADDRESS)
```
