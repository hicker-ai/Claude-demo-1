# PRD-14: 集成测试

## 目标

编写端到端集成测试，验证 HTTP API 和 LDAP 协议的完整工作流程。

## 依赖

- PRD-09 (服务启动)
- PRD-13 (前端 LDAP 配置)

## 交付物

- `tests/integration/ldap_test.go` — LDAP 集成测试
- `tests/integration/api_test.go` — HTTP API 集成测试
- `tests/integration/setup_test.go` — 测试环境搭建

## 详细要求

### 测试环境

`tests/integration/setup_test.go`:
- 启动测试服务器（HTTP + LDAP），使用 SQLite in-memory
- 提供 helper 函数：创建测试用户/组、获取 auth token
- TestMain 中启动/关闭服务器

### LDAP 集成测试

使用 `go-ldap/ldap/v3` 客户端：

1. **完整 Bind 流程：**
   - HTTP API 创建用户 → LDAP Bind 验证该用户 → 成功
   - LDAP Bind 错误密码 → 失败
   - LDAP Bind 禁用用户 → 失败

2. **Search 过滤测试（重点）：**
   - `(objectClass=inetOrgPerson)` → 返回所有用户
   - `(uid=testuser)` → 返回特定用户
   - `(&(objectClass=inetOrgPerson)(mail=*@test.com))` → AND + Substring
   - `(|(uid=user1)(uid=user2))` → OR
   - `(!(status=disabled))` → NOT
   - `(cn=Test*)` → Substring prefix
   - `(cn=*User)` → Substring final
   - `(&(objectClass=groupOfNames)(cn=admin*))` → 组搜索 + Substring
   - 嵌套：`(&(|(cn=A)(cn=B))(!(status=disabled)))` → 复杂嵌套
   - `(cn>=A)` → GreaterOrEqual
   - `(cn<=Z)` → LessOrEqual
   - 切换 AD 模式后：`(objectClass=user)` 和 `sAMAccountName` 正常工作

3. **用户组搜索：**
   - 搜索组返回 member 属性（DN 列表）
   - 搜索用户返回 memberOf 属性（AD 模式）

### HTTP API 集成测试

1. **用户生命周期：** 创建 → 获取 → 列表（搜索/分页） → 更新 → 修改密码 → 禁用 → 删除
2. **组生命周期：** 创建 → 创建子组 → 获取（含层级） → 添加成员 → 获取成员 → 移除成员 → 删除
3. **认证流程：** 登录 → 带 token 访问 → 无 token 访问被拒 → 错误 token 被拒
4. **校验错误：** 缺少必填字段、重复用户名、不存在的 ID

## 验收标准

- `go test ./tests/integration/ -v` 全部通过
- LDAP 过滤测试覆盖所有操作符类型
- HTTP API 测试覆盖全生命周期
- 两种 LDAP 模式均测试
