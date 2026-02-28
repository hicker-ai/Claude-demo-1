# PRD-11: 前端用户管理页面

## 目标

实现用户列表和用户详情/编辑页面，包括用户 CRUD、状态管理、密码修改、用户组关联展示。

## 依赖

- PRD-10 (前端项目初始化)

## 交付物

- `web/src/components/Layout.tsx` — 管理后台布局
- `web/src/pages/UserList.tsx` — 用户列表页
- `web/src/pages/UserDetail.tsx` — 用户详情/编辑页
- `web/src/api/users.ts` — 用户 API

## 详细要求

### 管理后台布局

`web/src/components/Layout.tsx`:
- Ant Design Layout + Sider + Content
- 左侧菜单：用户管理、用户组管理、LDAP 配置
- 顶部：应用标题、当前用户信息、退出登录按钮
- 响应式侧边栏（可折叠）

### 用户 API

`web/src/api/users.ts`:
```typescript
listUsers(params: { page: number; page_size: number; search?: string })
getUser(id: string)
createUser(data: CreateUserReq)
updateUser(id: string, data: UpdateUserReq)
deleteUser(id: string)
changePassword(id: string, data: { old_password: string; new_password: string })
setUserStatus(id: string, status: 'enabled' | 'disabled')
getUserGroups(id: string)
```

### 用户列表页

`web/src/pages/UserList.tsx`:
- Ant Design Table，列：用户名、显示名、邮箱、状态、创建时间、操作
- 状态列：Tag 组件（enabled=绿色，disabled=红色）
- 搜索框：搜索用户名、显示名、邮箱
- 分页组件
- "创建用户" 按钮 → 弹出 Modal 表单
- 行操作：编辑（跳转详情页）、启用/禁用（开关）、删除（Popconfirm 确认）
- 批量操作可选

### 用户详情/编辑页

`web/src/pages/UserDetail.tsx`:
- URL: `/users/:id`
- 用户信息表单：显示名、邮箱、手机（可编辑）、用户名（只读）
- 保存按钮
- 修改密码区域：旧密码 + 新密码 + 确认密码
- 所属用户组列表：Table 展示，可跳转到组详情
- 账号状态开关

## 验收标准

- 用户列表正确展示、分页、搜索
- 创建/编辑/删除用户流程顺畅
- 状态切换实时更新
- 密码修改功能正常
- 用户所属组列表正确展示
