# PRD-05: LDAP 过滤条件解析器 (RFC 4515)

## 目标

实现完整的 LDAP 过滤条件解析器，将 RFC 4515 过滤字符串解析为 AST，再将 AST 转换为 Ent ORM 的 SQL predicate，支持所有标准过滤操作符。

## 依赖

- PRD-01 完成（仅需 go module）

## 交付物

- `internal/ldap/filter/ast.go` — 过滤条件 AST 类型定义
- `internal/ldap/filter/parser.go` — 过滤字符串解析器
- `internal/ldap/filter/evaluator.go` — AST → Ent predicate 转换器
- `internal/ldap/filter/parser_test.go` — 解析器测试
- `internal/ldap/filter/evaluator_test.go` — 转换器测试

## 详细要求

### AST 类型

```
FilterType: And, Or, Not, Equal, Substring, GreaterOrEqual, LessOrEqual, Present, ApproxMatch

Filter {
    Type     FilterType
    Attr     string           // 属性名（叶子节点）
    Value    string           // 值（equal/gte/lte/approx）
    Children []*Filter        // 子过滤（AND/OR/NOT）
    Substr   *SubstringFilter // 子串匹配
}

SubstringFilter {
    Initial string   // 前缀（第一个 * 之前）
    Any     []string // 中间部分（* 之间）
    Final   string   // 后缀（最后一个 * 之后）
}
```

### 解析器

使用 `github.com/go-ldap/ldap/v3` 的 `CompileFilter` 将过滤字符串解析为 BER packet，然后遍历 BER 树构建自定义 AST。

**必须支持的过滤语法：**

| 语法 | 类型 | 示例 |
|------|------|------|
| `(attr=value)` | Equal | `(cn=John)` |
| `(attr=*)` | Presence | `(email=*)` |
| `(attr>=value)` | GreaterOrEqual | `(createTimestamp>=20240101)` |
| `(attr<=value)` | LessOrEqual | `(createTimestamp<=20241231)` |
| `(attr~=value)` | ApproxMatch | `(cn~=Jon)` |
| `(attr=init*any*final)` | Substring | `(cn=J*o*hn)`, `(cn=*ohn)`, `(cn=Jo*)` |
| `(&(f1)(f2)...)` | AND | `(&(cn=John)(mail=*@e.com))` |
| `(\|(f1)(f2)...)` | OR | `(\|(cn=John)(cn=Jane))` |
| `(!(filter))` | NOT | `(!(status=disabled))` |

**转义字符处理：**
- `\2a` → `*`
- `\28` → `(`
- `\29` → `)`
- `\5c` → `\`
- `\00` → NUL

**嵌套支持：** 任意深度嵌套，如 `(&(|(cn=A)(cn=B))(!(status=disabled))(mail=*@example.com))`

### Evaluator（AST → Ent Predicate）

接收属性映射器（AttrMapper 接口），将 LDAP 属性名转换为数据库列名。

**转换规则：**

| Filter Type | Ent Predicate |
|-------------|---------------|
| Equal | `sql.P().EQ(col, val)` |
| Present | `sql.P().Not().IsNull(col)` |
| Substring (initial) | `sql.P().HasPrefix(col, val)` |
| Substring (final) | `sql.P().HasSuffix(col, val)` |
| Substring (any) | `sql.P().Contains(col, val)` |
| Substring (complex) | 组合 LIKE pattern |
| GreaterOrEqual | `sql.P().GTE(col, val)` |
| LessOrEqual | `sql.P().LTE(col, val)` |
| ApproxMatch | `sql.P().EqualFold(col, val)` |
| AND | `sql.And(children...)` |
| OR | `sql.Or(children...)` |
| NOT | `sql.Not(child)` |

### 测试要求（重点模块，测试必须全面）

**解析器测试用例（至少覆盖以下）：**
1. 每种过滤类型的基本用例
2. 嵌套组合：`(&(|(a=1)(b=2))(!(c=3)))`
3. Substring 各种变体：仅前缀、仅后缀、仅中间、前+中+后
4. 转义字符
5. objectClass 过滤：`(objectClass=inetOrgPerson)`
6. 错误输入：空字符串、未闭合括号、无效语法
7. 复杂真实世界过滤：`(&(objectClass=inetOrgPerson)(|(uid=admin)(mail=*@company.com))(!(userAccountControl=disabled)))`

**Evaluator 测试用例：**
1. 每种过滤类型转换为正确的 predicate
2. 未知属性返回错误
3. 嵌套过滤生成正确的 AND/OR/NOT predicate 树

## 验收标准

- `go test ./internal/ldap/filter/ -v` 全部通过
- 所有 RFC 4515 标准过滤操作符均已实现
- 解析器测试用例 ≥ 15 个
- Evaluator 测试用例 ≥ 10 个
