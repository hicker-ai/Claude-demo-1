# PRD-06: LDAP DN 构建/解析 + 属性映射

## 目标

实现 LDAP Distinguished Name (DN) 的构建和解析模块，以及支持 OpenLDAP 和 Microsoft AD 两种规范的属性映射器。

## 依赖

- PRD-01 完成

## 交付物

- `internal/ldap/dn/dn.go` — DN 构建与解析
- `internal/ldap/dn/dn_test.go`
- `internal/ldap/attrs/mapper.go` — 属性映射器
- `internal/ldap/attrs/mapper_test.go`

## 详细要求

### DN 模块

**功能：**

- `BuildUserDN(username, displayName, baseDN, mode)` — 构建用户 DN
  - OpenLDAP: `uid=<username>,ou=users,<baseDN>`
  - AD: `cn=<displayName>,cn=Users,<baseDN>`
- `BuildGroupDN(groupName, baseDN, mode)` — 构建组 DN
  - OpenLDAP: `cn=<groupName>,ou=groups,<baseDN>`
  - AD: `cn=<groupName>,cn=Groups,<baseDN>`
- `ParseDN(dn)` — 解析 DN 为 RDN 列表 `[{Type:"uid", Value:"john"}, {Type:"ou", Value:"users"}, ...]`
- `ExtractUsername(dn, mode)` — 从 DN 提取用户名
  - OpenLDAP: 取 uid= 的值
  - AD: 取 cn= 的值 或 通过反查
- `BaseDNFromDN(dn)` — 提取 base DN 部分
- `IsUserDN(dn, baseDN, mode)` — 判断是否为用户 DN
- `IsGroupDN(dn, baseDN, mode)` — 判断是否为组 DN

### 属性映射器

**Mode 类型：** `openldap` | `activedirectory`

**User 属性映射：**

| DB 字段 | OpenLDAP 属性 | AD 属性 |
|---------|--------------|---------|
| username | uid | sAMAccountName |
| display_name | cn, displayName | cn, displayName |
| email | mail | mail |
| phone | telephoneNumber | telephoneNumber |
| status | (custom) status | userAccountControl |

**Group 属性映射：**

| DB 字段 | OpenLDAP 属性 | AD 属性 |
|---------|--------------|---------|
| name | cn | cn |
| description | description | description |
| members | member (DN list) | member (DN list) |

**接口：**

```go
type Mapper struct { mode Mode }

// MapAttribute 将 LDAP 属性名映射到数据库列名
func (m *Mapper) MapAttribute(ldapAttr string) (dbColumn string, ok bool)

// UserObjectClasses 返回当前模式下用户的 objectClass 列表
func (m *Mapper) UserObjectClasses() []string

// GroupObjectClasses 返回当前模式下组的 objectClass 列表
func (m *Mapper) GroupObjectClasses() []string

// UserToLDAPAttrs 将 domain.User 转换为 LDAP 属性 map
func (m *Mapper) UserToLDAPAttrs(user *domain.User) map[string][]string

// GroupToLDAPAttrs 将 domain.Group 转换为 LDAP 属性 map（含 member DN 列表）
func (m *Mapper) GroupToLDAPAttrs(group *domain.Group, memberDNs []string) map[string][]string
```

### 测试要求

**DN 测试：**
- OpenLDAP/AD 两种模式的 DN 构建
- DN 解析和组件提取
- 用户名提取
- 用户/组 DN 判断
- 特殊字符 DN 处理（含逗号、等号的值）

**属性映射测试：**
- 两种模式的属性名映射
- 未知属性返回 false
- User/Group → LDAP 属性转换完整性
- objectClass 返回值正确

## 验收标准

- `go test ./internal/ldap/dn/ -v` 全部通过
- `go test ./internal/ldap/attrs/ -v` 全部通过
- 两种 LDAP 规范均正确处理
