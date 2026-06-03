## Why

`POST /signin` 和 `POST /signup` 的 usecase 层当前为空函数 stub。所有依赖（FindByEmail、CheckPassword、CreateAccessToken、UpdateLastLogin）已存在，只需要在 usecase 中串联调用。前端 signin.vue 和 useAuth composable 已就绪，http.ts 缺少 JWT 拦截器。

## What Changes

- **usecase/account_usecase.go** — 实现 Signin（查用户 → 验密码 → 生成 JWT → 更新最后登录时间）和 Signup（查重 → 哈希 → 创建用户 → 生成 JWT）
- **controller/account_controller.go** — Signup 绑定正确的 `SignupRequest`（当前绑了 SigninRequest）
- **domain/account.go** — SignupRequest/SignupResponse 补充字段，UserRepository 接口新增 ExistsByEmail 和 Create 方法签名
- **repository/user_repository.go** — 实现 ExistsByEmail 和 Create
- **lib/http.ts** — 请求拦截器从 localStorage 读取 JWT 并附加 Authorization header

## Impact

- `apps/server/usecase/account_usecase.go`
- `apps/server/api/controller/account_controller.go`
- `apps/server/domain/account.go`
- `apps/server/repository/user_repository.go`
- `apps/omni-pixel/src/lib/http.ts`
