# PRD-12: 前端用户组管理页面

## 目标

实现用户组列表（树形结构）和用户组详情页面，支持组的 CRUD、成员管理。

## 依赖

- PRD-10 (前端项目初始化)

## 交付物

- `web/src/pages/GroupList.tsx` — 用户组列表页（树形）
- `web/src/pages/GroupDetail.tsx` — 用户组详情页
- `web/src/api/groups.ts` — 用户组 API

## 详细要求

### 用户组 API

`web/src/api/groups.ts`:
```typescript
listGroups()
getGroup(id: string)
createGroup(data: { name: string; description?: string; parent_id?: string })
updateGroup(id: string, data: { name?: string; description?: string; parent_id?: string })
deleteGroup(id: string)
getGroupMembers(id: string)
addMembers(id: string, userIds: string[])
removeMember(id: string, userId: string)
```

### 用户组列表页

`web/src/pages/GroupList.tsx`:
- Ant Design Tree 组件展示组的层级结构
- 树节点显示：组名、成员数
- 点击节点 → 右侧展示组详情（或跳转到详情页）
- "创建用户组" 按钮 → Modal 表单：
  - 组名（required）
  - 描述
  - 父组选择（TreeSelect，可选）
- 树节点右键菜单或操作按钮：编辑、删除（需确认）

### 用户组详情页

`web/src/pages/GroupDetail.tsx`:
- URL: `/groups/:id`
- 组信息表单：组名、描述、父组（可编辑）
- 保存按钮
- 成员列表：Ant Design Table（用户名、显示名、邮箱、操作-移除）
- "添加成员" 按钮 → Modal：
  - Ant Design Select/Transfer 组件
  - 搜索可用用户（排除已有成员）
  - 支持多选
- 移除成员：Popconfirm 确认后调用 API

## 验收标准

- 树形结构正确展示组层级
- 创建子组功能正常，树自动更新
- 成员添加/移除流程顺畅
- 删除含子组或成员的组给出提示
