# PRD-10: 前端项目初始化 + 登录页

## 目标

使用 Vite + React 18 + TypeScript 初始化前端项目，配置 Ant Design 5，实现 API 客户端封装和登录页面。

## 依赖

- 无（前端可与后端并行开发）

## 交付物

- `web/` — Vite React 项目（含 package.json、tsconfig 等）
- `web/src/api/client.ts` — Axios 封装
- `web/src/api/auth.ts` — 认证 API
- `web/src/pages/Login.tsx` — 登录页
- `web/src/App.tsx` — 路由配置
- `web/src/store/auth.ts` — 认证状态管理

## 详细要求

### 项目初始化

```bash
npm create vite@latest web -- --template react-ts
cd web
npm install antd @ant-design/icons axios react-router-dom
```

### API Client

`web/src/api/client.ts`:
- 基于 Axios 创建实例，baseURL: `/api/v1`
- Request 拦截器：从 localStorage 读取 token，添加 `Authorization: Bearer <token>`
- Response 拦截器：401 → 清除 token，跳转登录页；其他错误 → Ant Design message 提示

### 认证状态

`web/src/store/auth.ts`:
- token 存储在 localStorage
- 提供 `isAuthenticated()`, `getToken()`, `setToken()`, `clearToken()` 方法

### 登录页

`web/src/pages/Login.tsx`:
- 居中卡片布局
- Ant Design Form：用户名（required）、密码（required）
- 登录按钮 + loading 状态
- 成功 → 存储 token → 跳转 /users
- 失败 → 错误提示

### 路由

`web/src/App.tsx`:
- `/login` → Login 页面
- `/users`, `/groups`, `/ldap-config` → 需认证（AuthGuard）
- 未认证访问受保护路由 → 重定向到 /login
- 默认路由 `/` → 重定向到 /users

### Vite 代理配置

`web/vite.config.ts`:
```typescript
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080'
    }
  }
})
```

## 验收标准

- `cd web && npm run dev` 启动成功
- 登录页正常渲染
- API 客户端拦截器工作正常
- 路由守卫工作正常
