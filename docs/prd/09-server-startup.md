# PRD-09: 双端口服务启动与优雅关闭

## 目标

实现 Cobra `serve` 子命令，同时启动 HTTP Server 和 LDAP Server，监听不同端口，共享 Service 层，支持优雅关闭。

## 依赖

- PRD-07 (LDAP Handler)
- PRD-08 (HTTP Handler)

## 交付物

- `cmd/server/serve.go` — serve 子命令
- 修改 `cmd/server/root.go` — 注册 serve 命令

## 详细要求

### 启动流程

1. 加载配置（Viper）
2. 初始化 Zap logger
3. 初始化 Ent Client（PostgreSQL 连接）
4. 执行数据库 auto-migrate
5. 初始化各层：DAO → Service → Handler
6. 启动 HTTP Server（goroutine）
7. 启动 LDAP Server（goroutine）
8. 等待 SIGINT/SIGTERM 信号
9. 优雅关闭：HTTP Server → LDAP Server → DB 连接 → Logger flush

### 命令行用法

```bash
./usermanager serve --config configs/config.yaml
```

### 日志输出

启动时输出：
```
{"level":"info","msg":"HTTP server started","port":8080}
{"level":"info","msg":"LDAP server started","port":10389}
```

关闭时输出：
```
{"level":"info","msg":"Shutting down servers..."}
{"level":"info","msg":"Servers stopped"}
```

### 错误处理

- 配置加载失败 → 立即退出，输出错误
- 数据库连接失败 → 立即退出，输出错误
- HTTP/LDAP 端口被占用 → 立即退出，输出错误

## 验收标准

- `go build -o bin/usermanager ./cmd/server/` 编译成功
- `./bin/usermanager serve` 同时启动 HTTP 和 LDAP 服务
- Ctrl+C 优雅关闭，无 goroutine 泄漏
- 所有资源（DB 连接、logger）正确清理
