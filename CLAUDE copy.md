# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

GoRDS 是一个企业级数据库工单与数据查询平台，基于 **Gin (Go)** + **Vue 3 + Ant Design Vue** 构建。支持 SQL 审核、DDL/DML 工单管理与审批流、数据导出、MySQL/TiDB 实例访问控制。

## 常用命令

### 后端 (Go)

```bash
# 启动服务
cd backend && go run cmd/main.go -config config.yaml

# 运行所有测试
cd backend && go test ./...

# 运行单个包的测试
cd backend && go test ./middleware/

# 运行单个测试用例
cd backend && go test ./middleware/ -run TestOtpValidate

# 生产构建
cd backend && go build -ldflags "-s -w" -o GoRDS cmd/main.go
```

### 前端 (Vue 3 + Vite)

```bash
cd www
npm install
npm run dev       # 开发服务器，代理 /api 到 localhost:8083
npm run build     # 构建到 www/dist/，由后端嵌入
npm run lint      # ESLint 自动修复
npm run format    # Prettier 格式化
```

## 架构

### 后端启动流程

`cmd/main.go` → `app.Run(configFile)` → `bootstrap.InitializeConfig`（加载配置）→ `bootstrap.InitializeLog`（初始化日志）→ `bootstrap.InitializeDB`（自动迁移所有模型 + 种子数据）→ `bootstrap.InitializeRedis` → `bootstrap.InitializeCron`（定时任务）→ `app.RunServer()`（注册路由、挂载嵌入式前端 SPA、启动 Gin）。

### 请求生命周期

请求 → `api.Init()`（CORS、request-id、请求日志中间件）→ 各模块路由注册在 `/api/v1/` 或 `/api/v1/admin/` 分组下 → JWT 认证通过 `global.App.JWT.MiddlewareFunc()` 保护路由，管理员路由额外使用 `middleware.HasAdminPermission()`。

### 后端目录结构 (`backend/`)

| 路径 | 说明 |
|---|---|
| `cmd/main.go` | 入口，解析命令行参数，输出版本信息 |
| `internal/app/app.go` | 应用组装：注册路由、静态文件、启动 Gin |
| `internal/bootstrap/` | 配置加载、DB/Redis 初始化、定时任务、自动迁移与种子数据（管理员用户、DAS 允许操作、审核参数、通知配置） |
| `internal/global/global.go` | 全局单例 `Application` 结构体，持有 Config、JWT、Log、DB、Redis、Cron |
| `internal/config/` | Viper 映射的配置结构体 |
| `api/api.go` | Gin 引擎工厂：CORS、request-id、请求日志，收集路由选项 |
| `middleware/` | JWT 认证 (`jwt.go`)、管理员权限校验 (`permissions.go`)、TOTP 双因素认证 (`otp.go`)、请求日志 (`log.go`)、JWT claims 工具 (`claims.go`) |
| `pkg/` | 公共工具：`response/`（标准 JSON 响应信封）、`utils/`（加密、UUID、WebSocket、数据库工具）、`notifier/`（钉钉/微信/邮件通知）、`parser/`、`query/`（SQL 指纹）、`pagination/`、`kv/`（缓存） |
| `web/` | 通过 `//go:embed` 嵌入构建好的 Vue 前端 `dist/` |
| `internal/{模块}/` | 每个业务模块，遵循统一规范：`models/`、`routers/`、`services/`、`forms/`（请求绑定）、`views/`（Gin handler） |

### `internal/` 下的业务模块

- **users** — 用户、角色、组织架构、组织用户关联
- **common** — 数据库实例、实例环境、实例库表元数据（定时同步）、通知设置
- **inspect** — SQL 审核引擎：68 项审核规则，支持全局/实例级参数配置，使用 TiDB parser 做 AST 语法分析
- **das** — 数据查询：库表级访问控制、SQL 操作白名单、查询执行、收藏、Kill 任务
- **orders** — DDL/DML/数据导出工单：会签/或签审批流、工单任务、执行、回滚、消息通知

### 前端 (`www/`)

- Vue 3 + Pinia + Vue Router + Ant Design Vue 4
- Vite 开发服务器将 `/api`、`/media`、`/ws` 代理到 `http://localhost:8083`
- 路由守卫在 `src/permission.js`，静态路由在 `src/router/ark.js`
- Store：`user.js`（认证状态）、`permission.js`、`prefill.js`
- `src/api/` 下的 API 模块与后端一一对应：`login.js`、`das.js`、`order.js`、`admin.js`、`profile.js`

### 关键约定

- 后端模块遵循 **routers → views (handler) → services → models** 模式，`forms/` 存放请求结构体
- 所有 API 响应统一使用 `response.Response`（Code、Data、Message、RequestID）
- 配置一次性加载到 `global.App.Config`（`config.Configuration` 结构体）
- `app.secret_key` **必须恰好 32 个字符**——用于 JWT 签名和密码加密
- 前端 SPA 通过 `embed.FS` 嵌入；`npm run build` 后 Go 二进制文件即可作为单一部署产物
- gh-ost 用于 MySQL 的 Online DDL；TiDB 使用原生 Online DDL
