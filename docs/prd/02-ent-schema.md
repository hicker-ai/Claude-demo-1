# PRD-02: Ent Schema 定义

## 目标

使用 Ent ORM 定义 User 和 Group 的数据库 Schema，包含多对多关系和层级分组，生成 Ent 代码。

## 依赖

- PRD-01 完成

## 交付物

- `internal/schema/user.go` — User schema
- `internal/schema/group.go` — Group schema
- `internal/ent/` — Ent 生成代码
- `internal/ent/generate.go` — go:generate 指令

## 详细要求

### User Schema

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK, auto-gen, immutable |
| username | string | unique, not empty, max 64 |
| display_name | string | not empty, max 128 |
| email | string | unique, not empty, max 255 |
| password_hash | string | sensitive |
| phone | string | optional, max 32 |
| status | enum(enabled, disabled) | default: enabled |
| created_at | time | immutable, auto |
| updated_at | time | auto-update |

**Edge:** `groups` — many-to-many with Group (反向引用 Group.users)

### Group Schema

| 字段 | 类型 | 约束 |
|------|------|------|
| id | UUID | PK, auto-gen, immutable |
| name | string | unique, not empty, max 64 |
| description | string | optional, max 255 |
| parent_id | UUID | optional, nullable, FK → self |
| created_at | time | immutable, auto |
| updated_at | time | auto-update |

**Edges:**
- `users` — many-to-many with User
- `children` / `parent` — self-referencing (parent_id 字段)

## 验收标准

- `go generate ./internal/ent/` 成功生成代码
- `go build ./...` 编译通过
- 生成的代码包含 User、Group 的 CRUD 方法和 Edge 查询
