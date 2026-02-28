# PRD-04: 业务逻辑层

## 目标

实现 UserService、GroupService、AuthService，封装业务逻辑，包括密码哈希、JWT 认证、输入校验。

## 依赖

- PRD-03 完成

## 交付物

- `internal/service/user.go` — UserService
- `internal/service/group.go` — GroupService
- `internal/service/auth.go` — AuthService
- `internal/service/user_test.go`
- `internal/service/group_test.go`
- `internal/service/auth_test.go`

## 详细要求

### UserService

- `CreateUser(ctx, CreateUserInput)` — 校验输入，bcrypt 哈希密码，调用 DAO
- `GetUser(ctx, id)` — 获取用户（含 groups）
- `GetUserByUsername(ctx, username)` — 按用户名查询
- `ListUsers(ctx, ListUsersInput)` — 分页列表
- `UpdateUser(ctx, id, UpdateUserInput)` — 部分更新
- `DeleteUser(ctx, id)` — 删除用户
- `ChangePassword(ctx, id, oldPw, newPw)` — 验证旧密码后更新
- `SetUserStatus(ctx, id, status)` — 启用/禁用
- `Authenticate(ctx, username, password)` — 验证凭据，返回 User（供 LDAP Bind 和 HTTP Login 共用）
- `SearchUsers(ctx, predicates...)` — 供 LDAP Search 调用

### GroupService

- `CreateGroup(ctx, CreateGroupInput)` — 校验父组存在性
- `GetGroup(ctx, id)` — 获取组（含 users 和 children）
- `ListGroups(ctx)` — 树形结构列表
- `UpdateGroup(ctx, id, UpdateGroupInput)` — 校验无环引用
- `DeleteGroup(ctx, id)` — 校验无子组和成员后删除
- `AddMembers(ctx, groupID, userIDs)`
- `RemoveMember(ctx, groupID, userID)`
- `GetGroupMembers(ctx, groupID)`
- `GetUserGroups(ctx, userID)`
- `SearchGroups(ctx, predicates...)` — 供 LDAP Search 调用

### AuthService

- `Login(ctx, username, password)` — 认证 + 生成 JWT
- `ValidateToken(ctx, token)` — 验证 JWT，返回用户信息
- JWT secret 从配置读取，token 有效期可配置
- 使用 `github.com/golang-jwt/jwt/v5`

### 密码策略

- bcrypt，cost=10
- `golang.org/x/crypto/bcrypt`

## 验收标准

- `go test ./internal/service/ -v` 全部通过
- 认证逻辑同时服务 HTTP Login 和 LDAP Bind
- 密码存储安全（bcrypt 哈希）
