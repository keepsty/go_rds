# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

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


## 项目 Skills
以下 Skill 定义了标准化的开发流程（每个 Skill 是一个目录，核心指令在 SKILL.md 中）：
- `.Codex/skills/ant-design-vue/` - vue前端组件生成规范
- `.Codex/skills/go-rds-backend-structure/` - Gin后端API生成规范   
- `.Codex/skills/git-commit/` - Git 提交规范
- `.Codex/skills/security-audit/` - 代码安全审计

执行相关任务时，请先阅读对应 Skill 目录下的 SKILL.md 并严格遵循。
如 Skill 中包含 scripts/、resources/ 或 references/，请一并参考。

<!-- superpowers-zh:begin (do not edit between these markers) -->
# Superpowers-ZH 中文增强版

本项目已安装 superpowers-zh 技能框架（20 个 skills）。

## 核心规则

1. **收到任务时，先检查是否有匹配的 skill** — 哪怕只有 1% 的可能性也要检查
2. **设计先于编码** — 收到功能需求时，先用 brainstorming skill 做需求分析
3. **测试先于实现** — 写代码前先写测试（TDD）
4. **验证先于完成** — 声称完成前必须运行验证命令

## 可用 Skills

Skills 位于 `.Codex/skills/` 目录，每个 skill 有独立的 `SKILL.md` 文件。

- **brainstorming**: 在任何创造性工作之前必须使用此技能——创建功能、构建组件、添加功能或修改行为。在实现之前先探索用户意图、需求和设计。
- **chinese-code-review**: 中文 review 沟通参考——话术模板、分级标注（必须修复/建议修改/仅供参考）、国内团队常见反模式应对。仅在用户显式 /chinese-code-review 时调用，不要根据上下文自动触发。
- **chinese-commit-conventions**: 中文 commit 与 changelog 配置参考——Conventional Commits 中文适配、commitlint/husky/commitizen 中文模板、conventional-changelog 中文配置。仅在用户显式 /chinese-commit-conventions 时调用，不要根据上下文自动触发。
- **chinese-documentation**: 中文文档排版参考——中英文空格、全半角标点、术语保留、链接格式、中文文案排版指北约定。仅在用户显式 /chinese-documentation 时调用，不要根据上下文自动触发。
- **chinese-git-workflow**: 国内 Git 平台配置参考——Gitee、Coding.net、极狐 GitLab、CNB 的 SSH/HTTPS/凭据/CI 接入差异与镜像同步配置。仅在用户显式 /chinese-git-workflow 时调用，不要根据上下文自动触发。
- **dispatching-parallel-agents**: 当面对 2 个以上可以独立进行、无共享状态或顺序依赖的任务时使用
- **executing-plans**: 当你有一份书面实现计划需要在单独的会话中执行，并设有审查检查点时使用
- **finishing-a-development-branch**: 当实现完成、所有测试通过、需要决定如何集成工作时使用——通过提供合并、PR 或清理等结构化选项来引导开发工作的收尾
- **mcp-builder**: MCP 服务器构建方法论 — 系统化构建生产级 MCP 工具，让 AI 助手连接外部能力
- **receiving-code-review**: 收到代码审查反馈后、实施建议之前使用，尤其当反馈不明确或技术上有疑问时——需要技术严谨性和验证，而非敷衍附和或盲目执行
- **requesting-code-review**: 完成任务、实现重要功能或合并前使用，用于验证工作成果是否符合要求
- **subagent-driven-development**: 当在当前会话中执行包含独立任务的实现计划时使用
- **systematic-debugging**: 遇到任何 bug、测试失败或异常行为时使用，在提出修复方案之前执行
- **test-driven-development**: 在实现任何功能或修复 bug 时使用，在编写实现代码之前
- **using-git-worktrees**: 当需要开始与当前工作区隔离的功能开发，或在执行实现计划之前使用——通过原生工具或 git worktree 回退机制确保隔离工作区存在
- **using-superpowers**: 在开始任何对话时使用——确立如何查找和使用技能，要求在任何响应（包括澄清性问题）之前调用 Skill 工具
- **verification-before-completion**: 在宣称工作完成、已修复或测试通过之前使用，在提交或创建 PR 之前——必须运行验证命令并确认输出后才能声称成功；始终用证据支撑断言
- **workflow-runner**: 在 Codex / OpenClaw / Cursor 中直接运行 agency-orchestrator YAML 工作流——无需 API key，使用当前会话的 LLM 作为执行引擎。当用户提供 .yaml 工作流文件或要求多角色协作完成任务时触发。
- **writing-plans**: 当你有规格说明或需求用于多步骤任务时使用，在动手写代码之前
- **writing-skills**: 当创建新技能、编辑现有技能或在部署前验证技能是否有效时使用

## 如何使用

当任务匹配某个 skill 时，使用 `Skill` 工具加载对应 skill 并严格遵循其流程。绝不要用 Read 工具读取 SKILL.md 文件。

如果你认为哪怕只有 1% 的可能性某个 skill 适用于你正在做的事情，你必须调用该 skill 检查。
<!-- superpowers-zh:end -->
