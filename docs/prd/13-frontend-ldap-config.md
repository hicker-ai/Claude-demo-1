# PRD-13: 前端 LDAP 配置页面

## 目标

实现 LDAP 服务配置管理页面，展示服务状态，支持修改 LDAP 配置参数。

## 依赖

- PRD-10 (前端项目初始化)

## 交付物

- `web/src/pages/LDAPConfig.tsx` — LDAP 配置页
- `web/src/api/ldap.ts` — LDAP 配置 API

## 详细要求

### LDAP API

`web/src/api/ldap.ts`:
```typescript
getConfig()     // GET /api/v1/ldap/config
updateConfig(data: LDAPConfigReq)  // PUT /api/v1/ldap/config
getStatus()     // GET /api/v1/ldap/status
```

### LDAP 配置页面

`web/src/pages/LDAPConfig.tsx`:

**服务状态卡片：**
- 运行状态指示灯（绿色=运行中，红色=停止）
- LDAP 端口
- 当前连接数
- 自动刷新（轮询 /ldap/status）

**配置表单：**
- Base DN：文本输入框，如 `dc=example,dc=com`
- 模式选择：Radio 组件，`OpenLDAP` / `Active Directory`
- LDAP 端口：数字输入框
- 保存按钮

**使用说明卡片：**
- 根据当前模式展示示例 ldapsearch 命令
- OpenLDAP 示例：`ldapsearch -x -H ldap://localhost:10389 -b "dc=example,dc=com" "(objectClass=inetOrgPerson)"`
- AD 示例：`ldapsearch -x -H ldap://localhost:10389 -b "dc=example,dc=com" "(objectClass=user)"`

## 验收标准

- 状态卡片正确展示 LDAP 服务状态
- 配置修改后保存成功
- 模式切换后示例命令自动更新
