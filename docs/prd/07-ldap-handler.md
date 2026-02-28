# PRD-07: LDAP Handler (Bind/Search)

## 目标

基于 gldap 框架实现 LDAP Bind（认证）和 Search（搜索）handler，将 LDAP 协议请求转换为 Service 层调用。

## 依赖

- PRD-04 (Service Layer)
- PRD-05 (Filter Parser)
- PRD-06 (DN & Attrs)

## 交付物

- `internal/handler/ldap/handler.go` — Handler 结构体与路由注册
- `internal/handler/ldap/bind.go` — Bind 处理
- `internal/handler/ldap/search.go` — Search 处理
- `internal/handler/ldap/handler_test.go` — 集成测试

## 详细要求

### Handler 结构

```go
type Handler struct {
    userService  UserService   // interface, 定义在此层
    groupService GroupService  // interface, 定义在此层
    mapper       *attrs.Mapper
    baseDN       string
    logger       *zap.Logger
}
```

**接口定义（在消费层定义，遵循项目架构规则）：**

```go
type UserService interface {
    Authenticate(ctx context.Context, username, password string) (*domain.User, error)
    GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
    SearchUsers(ctx context.Context, ...) ([]*domain.User, error)
}

type GroupService interface {
    SearchGroups(ctx context.Context, ...) ([]*domain.Group, error)
    GetGroupMembers(ctx context.Context, groupID uuid.UUID) ([]*domain.User, error)
}
```

### Bind Handler

1. 接收 `gldap.Request`，获取 SimpleBindMessage
2. 从 DN 中解析用户名（使用 dn.ExtractUsername）
3. 调用 `UserService.Authenticate(username, password)`
4. 成功 → `ResultSuccess`，失败 → `ResultInvalidCredentials`
5. 日志记录绑定尝试（成功/失败，不记录密码）

### Search Handler

1. 接收 `gldap.Request`，获取 SearchMessage
2. 获取 BaseDN、Scope、Filter、Attributes、SizeLimit、TimeLimit
3. 解析 Filter 字符串 → AST（`filter.Parse`）
4. 分析 Filter 判断搜索目标：
   - 包含 `objectClass=inetOrgPerson` / `objectClass=user` → 搜索用户
   - 包含 `objectClass=groupOfNames` / `objectClass=group` → 搜索组
   - 包含 `objectClass=*` 或无 objectClass → 搜索用户和组
5. 将 Filter AST 转换为 Ent predicate（`evaluator.ToPredicate`）
6. 调用 Service 层查询
7. 将结果转换为 LDAP Entry（使用 attrs.Mapper）
8. 处理 Scope:
   - BaseObject: 仅匹配 BaseDN 指定的对象
   - SingleLevel: BaseDN 直接子节点
   - WholeSubtree: BaseDN 及所有后代
9. 处理 Attributes 请求：仅返回客户端请求的属性（空列表 = 全部）
10. 处理 SizeLimit：限制返回条目数
11. 写入 SearchResponseEntry 和 SearchDoneResponse

### 错误处理

| 场景 | LDAP ResultCode |
|------|----------------|
| 认证失败 | ResultInvalidCredentials (49) |
| 无权限 | ResultInsufficientAccessRights (50) |
| DN 不存在 | ResultNoSuchObject (32) |
| 过滤条件无效 | ResultProtocolError (2) |
| 服务器错误 | ResultOther (80) |

### 测试要求

使用 `go-ldap/ldap/v3` 作为 LDAP 客户端进行测试：

1. **Bind 测试：**
   - 正确凭据绑定成功
   - 错误密码绑定失败
   - 不存在用户绑定失败
   - 禁用用户绑定失败

2. **Search 测试：**
   - `(objectClass=inetOrgPerson)` — 返回所有用户
   - `(uid=john)` — 返回特定用户
   - `(&(objectClass=inetOrgPerson)(mail=*@example.com))` — 复合过滤
   - `(objectClass=groupOfNames)` — 返回所有组
   - `(|(cn=admin)(cn=users))` — OR 过滤组
   - `(!(status=disabled))` — NOT 过滤
   - `(cn=J*)` — 子串过滤
   - Scope 处理正确性
   - Attributes 选择正确性
   - SizeLimit 限制正确性
   - AD 模式下的 `(objectClass=user)` 和 `sAMAccountName`

## 验收标准

- `go test ./internal/handler/ldap/ -v` 全部通过
- 能通过标准 LDAP 客户端（ldapsearch 命令行工具）连接和查询
- OpenLDAP 和 AD 模式均正常工作
