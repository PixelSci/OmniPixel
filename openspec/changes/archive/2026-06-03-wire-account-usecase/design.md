## Context

所有基础设施已就绪，只需在 usecase 层串联调用。Signin 流程：`FindByEmail → CheckPassword → CreateAccessToken → UpdateLastLogin`。Signup 流程：`ExistsByEmail → HashPassword → Create → CreateAccessToken`。

## Goals / Non-Goals

**Goals:**
- 实现 usecase.Signin() 和 usecase.Signup()，串联已有接口
- 修复 controller.Signup 绑定错误的 request 类型
- SignupRequest/SignupResponse 补充必要字段
- UserRepository 新增 ExistsByEmail、Create 方法
- 前端 http.ts 请求拦截器附加 JWT

**Non-Goals:**
- 不新增 API 端点（只实现已有 stub）
- 不添加邮箱验证、OAuth、密码重置等新功能
- 不新增中间件或路由

## Decisions

### Signin 使用已有依赖

```
controller.Signin()
  → usecase.Signin(request)
      → userRepo.FindByEmail(email)      // 已存在
      → internal.CheckPassword(pw, hash)  // 已存在 (apps/server/internal/password.go)
      → internal.CreateAccessToken(...)   // 已存在 (apps/server/internal/token.go)
      → userRepo.UpdateLastLogin(userID)  // 已存在
```

`FindByEmail` 在找不到用户时返回 `ErrInvalidCredentials`，usecase 透传该错误。密码不匹配时 usecase 返回 `ErrInvalidCredentials`。

### Signup 补充最小字段

`SignupRequest` 已有 username/email/password 字段（前端 auth.ts 已定义），domain 层 struct 当前为空。补充为：

```go
type SignupRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

### Repository 新增两个方法

`ExistsByEmail` 复用已有 `FindByEmail` 的 `lower(email)` 查询模式，`Create` 使用 GORM `db.Create`。

### 前端 JWT 拦截器

`useAuth.ts` 已通过 `useLocalStorage` 管理 token。拦截器从同一 key 读取：

```ts
// http.ts request interceptor
const token = localStorage.getItem('omni-pixel:access-token')
if (token) config.headers.Authorization = `Bearer ${JSON.parse(token)}`
```

## Risks / Trade-offs

- Signin/Signup 当前在公开路由组（无 JWT 保护）— 正确行为
- `password.go` 使用 bcrypt cost 12 — 登录时可能较慢（~300ms），暂不调整
