# PRD-03: 领域模型与数据访问层

## 目标

定义领域模型（与 Ent 解耦），实现 DAO 层封装 Ent 查询，提供 User/Group 的完整 CRUD 和关联操作。

## 依赖

- PRD-02 完成

## 交付物

- `internal/domain/user.go` — User 领域模型与输入结构体
- `internal/domain/group.go` — Group 领域模型与输入结构体
- `internal/dao/dao.go` — DAO 初始化与迁移
- `internal/dao/user.go` — User DAO 实现
- `internal/dao/group.go` — Group DAO 实现
- `internal/dao/user_test.go` — User DAO 测试
- `internal/dao/group_test.go` — Group DAO 测试

## 详细要求

### 领域模型

- `domain.User` — 包含所有字段和 Groups 切片
- `domain.Group` — 包含所有字段、Children 和 Users 切片
- `domain.CreateUserInput` / `domain.UpdateUserInput` — 输入 DTO
- `domain.CreateGroupInput` / `domain.UpdateGroupInput` — 输入 DTO
- `domain.ListResult[T]` — 分页结果泛型

### DAO 接口

**UserDAO:**
- `Create(ctx, *ent.UserCreate) → *domain.User`
- `GetByID(ctx, uuid) → *domain.User`
- `GetByUsername(ctx, string) → *domain.User`
- `List(ctx, page, pageSize, search) → *domain.ListResult[domain.User]`
- `Update(ctx, id, input) → *domain.User`
- `Delete(ctx, id) → error`
- `GetUserGroups(ctx, userID) → []*domain.Group`

**GroupDAO:**
- `Create(ctx, input) → *domain.Group`
- `GetByID(ctx, uuid) → *domain.Group`
- `List(ctx) → []*domain.Group` (含层级)
- `Update(ctx, id, input) → *domain.Group`
- `Delete(ctx, id) → error`
- `AddMembers(ctx, groupID, userIDs) → error`
- `RemoveMember(ctx, groupID, userID) → error`
- `GetMembers(ctx, groupID) → []*domain.User`

### 测试要求

- 使用 SQLite in-memory 进行测试（Ent 支持方言切换）
- 表驱动测试，覆盖正常路径和边界条件
- 覆盖：创建、查询、更新、删除、分页、关联操作

## 验收标准

- `go test ./internal/dao/ -v` 全部通过
- DAO 层与领域模型完全解耦于 Ent 生成代码
