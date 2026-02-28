# User Management System with LDAP Server

用户管理系统，提供 HTTP REST API 和 LDAP 协议双端口服务。支持用户/用户组 CRUD、LDAP Bind 认证和 Search 查询，兼容 OpenLDAP 和 Active Directory 两种模式。

## 技术栈

- **后端**: Go 1.24 + Gin + Ent ORM + gldap + Zap
- **前端**: React 19 + Ant Design + TypeScript + Vite
- **数据库**: PostgreSQL
- **认证**: JWT (HTTP) / LDAP Bind

## 前置要求

- Go 1.24+
- Node.js 18+
- PostgreSQL 14+

## 快速开始

### 1. 准备数据库

```bash
# 创建数据库
createdb usermanager

# 或通过 psql
psql -U postgres -c "CREATE DATABASE usermanager;"
```

### 2. 修改配置

编辑 `configs/config.yaml`，根据实际环境修改数据库连接信息和 JWT 密钥：

```yaml
server:
  http_port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  dbname: usermanager
  sslmode: disable

ldap:
  port: 10389
  base_dn: "dc=example,dc=com"
  mode: "openldap"  # 或 "activedirectory"

jwt:
  secret: "your-secret-key"  # 生产环境务必修改
  expire_hours: 24
```

### 3. 启动后端

```bash
# 编译
go build -o bin/usermanager .

# 启动（自动执行数据库迁移）
./bin/usermanager serve --config configs/config.yaml
```

启动后会同时监听两个端口：
- **HTTP API**: `http://localhost:8080`
- **LDAP Server**: `ldap://localhost:10389`

### 4. 启动前端（开发模式）

```bash
cd web
npm install
npm run dev
```

浏览器访问 `http://localhost:5173`，前端开发服务器会自动代理 `/api` 请求到后端 `localhost:8080`。

### 5. 前端生产构建

```bash
cd web
npm run build
```

构建产物在 `web/dist/`，可部署到 Nginx 等静态文件服务器，并将 `/api` 反向代理到后端。

## 创建初始用户

系统启动后没有默认用户，需要通过 API 创建第一个用户（登录接口不需要 token）：

```bash
# 1. 先通过 API 创建用户（此时无需认证）
#    注意：如果 /api/v1/users 需要认证，可以临时修改代码或直接操作数据库
#    推荐做法：通过 LDAP 或直接插入数据库创建管理员

# 2. 登录获取 token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'

# 返回示例：
# {"code":0,"message":"success","data":{"token":"eyJhbG...","user":{"id":"...","username":"admin","display_name":"Admin"}}}
```

## HTTP API

所有 API 路径前缀为 `/api/v1`，除登录/登出外均需 `Authorization: Bearer <token>` 头。

### 认证

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/auth/login` | 登录，返回 JWT token |
| POST | `/auth/logout` | 登出 |

### 用户管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/users` | 创建用户 |
| GET | `/users` | 用户列表（支持 `?search=&page=&page_size=`） |
| GET | `/users/:id` | 获取用户详情 |
| PUT | `/users/:id` | 更新用户信息 |
| DELETE | `/users/:id` | 删除用户 |
| PUT | `/users/:id/password` | 修改密码 |
| PUT | `/users/:id/status` | 启用/禁用用户 |
| GET | `/users/:id/groups` | 获取用户所属组 |

### 用户组管理

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/groups` | 创建用户组 |
| GET | `/groups` | 用户组列表 |
| GET | `/groups/:id` | 获取用户组详情 |
| PUT | `/groups/:id` | 更新用户组 |
| DELETE | `/groups/:id` | 删除用户组 |
| POST | `/groups/:id/members` | 添加成员 |
| DELETE | `/groups/:id/members/:uid` | 移除成员 |
| GET | `/groups/:id/members` | 获取成员列表 |

### LDAP 配置

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/ldap/config` | 获取 LDAP 配置 |
| PUT | `/ldap/config` | 更新 LDAP 配置 |
| GET | `/ldap/status` | 获取 LDAP 服务状态 |

## LDAP 使用

### Bind 认证

```bash
# OpenLDAP 模式
ldapwhoami -H ldap://localhost:10389 \
  -D "uid=admin,ou=users,dc=example,dc=com" \
  -w password123

# Active Directory 模式
ldapwhoami -H ldap://localhost:10389 \
  -D "CN=Admin,CN=Users,dc=example,dc=com" \
  -w password123
```

### Search 查询

```bash
# 搜索所有用户
ldapsearch -H ldap://localhost:10389 -x \
  -b "dc=example,dc=com" \
  "(objectClass=inetOrgPerson)"

# 按用户名搜索
ldapsearch -H ldap://localhost:10389 -x \
  -b "dc=example,dc=com" \
  "(uid=admin)"

# 组合过滤条件
ldapsearch -H ldap://localhost:10389 -x \
  -b "dc=example,dc=com" \
  "(&(objectClass=inetOrgPerson)(mail=*@example.com))"

# 搜索用户组
ldapsearch -H ldap://localhost:10389 -x \
  -b "dc=example,dc=com" \
  "(objectClass=groupOfNames)" cn member

# 复杂嵌套过滤
ldapsearch -H ldap://localhost:10389 -x \
  -b "dc=example,dc=com" \
  "(&(|(cn=Admin*)(cn=Dev*))(!(status=disabled)))"
```

### 支持的 LDAP 过滤类型

- 等于: `(uid=admin)`
- 存在: `(mail=*)`
- 子串: `(cn=Admin*)`, `(cn=*test*)`, `(cn=*User)`
- 大于等于: `(uid>=b)`
- 小于等于: `(uid<=m)`
- 近似匹配: `(cn~=admin)`
- 与: `(&(a=1)(b=2))`
- 或: `(|(a=1)(a=2))`
- 非: `(!(status=disabled))`
- 嵌套组合: `(&(|(cn=A)(cn=B))(!(status=disabled)))`

## 运行测试

```bash
# 全部测试
go test ./...

# 集成测试（含 HTTP API + LDAP 协议）
go test ./tests/integration/ -v

# 前端类型检查
cd web && npx tsc --noEmit
```

## 项目结构

```
├── main.go              # 程序入口
├── cmd/                 # Cobra 命令（root、serve）
├── configs/             # 配置文件
├── internal/
│   ├── config/          # 配置加载
│   ├── dao/             # 数据访问层
│   ├── domain/          # 领域模型
│   ├── ent/             # Ent 生成代码
│   ├── handler/
│   │   ├── http/        # HTTP API 处理器
│   │   └── ldap/        # LDAP 协议处理器
│   ├── ldap/
│   │   ├── attrs/       # LDAP 属性映射
│   │   ├── dn/          # DN 构建与解析
│   │   └── filter/      # RFC 4515 过滤条件解析
│   ├── middleware/       # HTTP 中间件
│   ├── schema/          # Ent Schema 定义
│   └── service/         # 业务逻辑层
├── tests/integration/   # 集成测试
└── web/                 # React 前端
```
