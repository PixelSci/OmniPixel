## 1. Domain 层

- [x] 1.1 `domain/account.go` — SignupRequest 补充 Username/Email/Password 字段，SignupResponse 补充 AccessToken/User 字段，UserRepository 接口新增 ExistsByEmail、Create 方法签名

## 2. Repository 层

- [x] 2.1 `repository/user_repository.go` — 实现 ExistsByEmail(email) 和 Create(user) 方法

## 3. UseCase 层

- [x] 3.1 `usecase/account_usecase.go` — 实现 Signin：FindByEmail → CheckPassword → CreateAccessToken → UpdateLastLogin
- [x] 3.2 `usecase/account_usecase.go` — 实现 Signup：ExistsByEmail → HashPassword → Create → CreateAccessToken

## 4. Controller 层

- [x] 4.1 `api/controller/account_controller.go` — Signup handler 绑定 SignupRequest 并调用 usecase.Signup()

## 5. 前端

- [x] 5.1 `apps/omni-pixel/src/lib/http.ts` — 请求拦截器从 localStorage 读取 JWT，附加 Bearer Authorization header

## 6. 验证

- [x] 6.1 `go build ./...` 编译通过
- [x] 6.2 前端 `pnpm build` 或 `npm run build` 编译通过
