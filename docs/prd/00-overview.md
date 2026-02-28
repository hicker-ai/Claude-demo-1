# 用户管理系统 — 子任务总览

## 技术栈

| 层 | 技术 |
|---|---|
| 后端 | Go 1.22+ / Gin / Ent / Zap / Cobra |
| 前端 | React 18 / Ant Design 5 / Vite / TypeScript (web/) |
| 数据库 | PostgreSQL |
| LDAP | gldap (server) / go-ldap/ldap/v3 (filter) |

## 架构

单二进制、双端口：HTTP(:8080) + LDAP(:10389)，共享 Service→DAO→PostgreSQL。

## 子任务清单

| # | PRD 文件 | 模块 | 依赖 |
|---|---------|------|------|
| 01 | `01-project-scaffold.md` | 项目脚手架、Go Module、Cobra CLI、配置 | 无 |
| 02 | `02-ent-schema.md` | Ent Schema 定义与代码生成 | 01 |
| 03 | `03-domain-dao.md` | 领域模型与数据访问层 | 02 |
| 04 | `04-service-layer.md` | 业务逻辑层 (User/Group/Auth Service) | 03 |
| 05 | `05-ldap-filter-parser.md` | LDAP 过滤条件解析器 (RFC 4515) | 01 |
| 06 | `06-ldap-dn-attrs.md` | LDAP DN 构建/解析 + 属性映射 (OpenLDAP/AD) | 01 |
| 07 | `07-ldap-handler.md` | LDAP Handler (Bind/Search) | 04, 05, 06 |
| 08 | `08-http-handler.md` | HTTP REST API Handler + 中间件 | 04 |
| 09 | `09-server-startup.md` | 双端口服务启动与优雅关闭 | 07, 08 |
| 10 | `10-frontend-setup-login.md` | 前端项目初始化 + 登录页 | 无 |
| 11 | `11-frontend-user-pages.md` | 前端用户管理页面 | 10 |
| 12 | `12-frontend-group-pages.md` | 前端用户组管理页面 | 10 |
| 13 | `13-frontend-ldap-config.md` | 前端 LDAP 配置页面 | 10 |
| 14 | `14-integration-tests.md` | 集成测试 | 09, 13 |

## 执行原则

- 每个子任务独立完成后必须自测通过
- 遵循 TDD：先写测试 → 验证失败 → 实现 → 验证通过
- 每个子任务完成后提交 git commit
- 遵循 `.ai-context/rules/` 中的所有规范
