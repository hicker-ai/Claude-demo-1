# PRD-08: HTTP REST API Handler + 中间件

## 目标

使用 Gin 实现 RESTful HTTP API，提供用户/组 CRUD、认证、LDAP 配置管理接口，以及 JWT 认证中间件和日志中间件。

## 依赖

- PRD-04 (Service Layer)

## 交付物

- `internal/handler/http/response.go` — 统一响应格式
- `internal/handler/http/router.go` — 路由注册
- `internal/handler/http/user.go` — 用户 API handler
- `internal/handler/http/group.go` — 用户组 API handler
- `internal/handler/http/auth.go` — 认证 API handler
- `internal/handler/http/ldap_config.go` — LDAP 配置 API handler
- `internal/middleware/auth.go` — JWT 认证中间件
- `internal/middleware/logger.go` — 请求日志中间件
- `internal/handler/http/user_test.go`
- `internal/handler/http/group_test.go`
- `internal/handler/http/auth_test.go`

## 详细要求

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误响应：`code` 非零，`message` 描述错误原因，`data` 可选。

### 请求 DTO（遵循代码风格规则）

所有请求 DTO 必须包含 `json` tag 和 `binding` tag：

```go
type CreateUserReq struct {
    Username    string `json:"username" binding:"required,max=64"`
    DisplayName string `json:"display_name" binding:"required,max=128"`
    Email       string `json:"email" binding:"required,email,max=255"`
    Password    string `json:"password" binding:"required,min=8"`
    Phone       string `json:"phone,omitempty" binding:"omitempty,max=32"`
}
```

### API 路由

```
POST   /api/v1/auth/login
POST   /api/v1/auth/logout

# 以下需要 JWT 认证
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
PUT    /api/v1/users/:id/password
PUT    /api/v1/users/:id/status
GET    /api/v1/users/:id/groups

POST   /api/v1/groups
GET    /api/v1/groups
GET    /api/v1/groups/:id
PUT    /api/v1/groups/:id
DELETE /api/v1/groups/:id
POST   /api/v1/groups/:id/members
DELETE /api/v1/groups/:id/members/:uid
GET    /api/v1/groups/:id/members

GET    /api/v1/ldap/config
PUT    /api/v1/ldap/config
GET    /api/v1/ldap/status
```

### 中间件

**JWT 认证中间件：**
- 从 `Authorization: Bearer <token>` 头提取 token
- 调用 AuthService.ValidateToken 验证
- 失败返回 401
- 成功将用户信息注入 gin.Context

**日志中间件：**
- 使用 Zap 结构化日志
- 记录：method, path, status, latency, request_id
- request_id 由中间件生成并注入 context

### 分页

列表接口支持分页参数：`?page=1&page_size=20&search=keyword`

响应格式：
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}
```

### 测试要求

- 使用 `httptest` + Gin test mode
- 表驱动测试覆盖每个接口：
  - 正常操作
  - 校验错误（缺少必填字段、格式错误）
  - 资源不存在 (404)
  - 重复键 (409)
  - 未认证 (401)

## 验收标准

- `go test ./internal/handler/http/ -v` 全部通过
- 所有 API 遵循 RESTful 规范
- 请求 DTO 包含正确的 json + binding tag
- 统一响应格式
