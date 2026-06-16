现在我已全面了解项目架构，以下是完整的开发文档。

---

# GoRDS — 数据库工单平台开发文档

> 项目定位：一个面向 DBA 和开发者的 **数据库工单审批 + SQL 查询 + SQL 审核** 平台。支持 MySQL / TiDB / ClickHouse，完全中文化。

---

## 目录

1. [项目概览](#1-项目概览)
2. [技术栈](#2-技术栈)
3. [项目结构](#3-项目结构)
4. [启动与配置](#4-启动与配置)
5. [后端架构详解](#5-后端架构详解)
6. [前端架构详解](#6-前端架构详解)
7. [核心业务流程](#7-核心业务流程)
8. [开发指南](#8-开发指南)

---

## 1. 项目概览

**GoRDS** 解决数据库变更和查询的管控问题：

| 功能模块 | 描述 |
|---------|------|
| **SQL 查询 (DAS)** | 在线查询 MySQL/TiDB/ClickHouse，支持表名补全、收藏、历史 |
| **数据库工单** | DDL/DML 变更提交→SQL 审核→审批流→自动执行，全生命周期管理 |
| **SQL 审核 (Inspect)** | 基于 TiDB Parser 的语法树遍历，80+ 审核规则（表名长度、索引规范、DDL 安全等） |
| **用户权限** | 角色/组织/审批流管理，粒度到 schema 和 table 级别 |
| **消息通知** | 企业微信/钉钉/邮件多渠道推送 |
| **OTP 双因素认证** | 登录可选绑定 TOTP |

---

## 2. 技术栈

### 后端 (`backend/`)

| 组件 | 选型 |
|------|------|
| 语言 | **Go 1.22+** |
| HTTP 框架 | **Gin** + gin-jwt(v2) |
| ORM | **GORM** (MySQL 驱动) |
| 数据库 | **MySQL 8.0** (业务数据) / DAS 查询直连用户数据库 |
| 缓存 | **Redis** (可选) |
| SQL Parser | **pingcap/tidb/parser** — 解析 SQL 为 AST，用于审核 |
| 配置管理 | **Viper** (yaml) |
| 定时任务 | **robfig/cron** |
| 日志 | **Logrus** |
| 前端嵌入 | `embed` 打包 dist 到 Go 二进制 |

### 前端 (`www/`)

| 组件 | 选型 |
|------|------|
| 框架 | **Vue 3** (Composition API + `<script setup>`) |
| 构建工具 | **Vite 7** |
| UI 组件库 | **Ant Design Vue 4** |
| 状态管理 | **Pinia** (持久化插件) |
| 路由 | **Vue Router 4** |
| 代码编辑器 | **CodeMirror 6** (SQL + JSON 语法高亮) |
| HTTP 请求 | **Axios** (带统一错误处理) |

---

## 3. 项目结构

```
GoRDS/
├── backend/                          # Go 后端
│   ├── cmd/main.go                   # 入口：解析 flags → app.Run()
│   ├── config.yaml / .template       # 应用配置
│   ├── go.mod / go.sum
│   ├── api/
│   │   └── api.go                    # Gin 引擎初始化（CORS/日志/路由注册）
│   ├── middleware/
│   │   ├── jwt.go                    # JWT 认证中间件（gin-jwt v2）
│   │   ├── otp.go                    # OTP 双因素中间件
│   │   ├── permissions.go            # 管理员权限检查
│   │   ├── claims.go                 # JWT Claims 工具函数
│   │   └── log.go                    # 请求日志
│   ├── internal/
│   │   ├── global/global.go          # 全局单例：Config/DB/Redis/JWT/Log/Cron
│   │   ├── config/config.go          # Configuration 结构体（YAML → struct）
│   │   ├── app/app.go                # Run() + RunServer() 启动流程
│   │   ├── bootstrap/
│   │   │   ├── config.go             # Viper 加载配置
│   │   │   ├── db.go                 # GORM + Redis 初始化 + 表 AutoMigrate + 种子数据
│   │   │   ├── cron.go               # 定时任务注册
│   │   │   └── log.go                # Logger 初始化
│   │   ├── common/                   # 公共模块
│   │   │   ├── models/               # 基础模型（Model/EnumType）+ 实例/环境/通知模型
│   │   │   ├── forms/                # 请求表单结构体
│   │   │   ├── views/                # 管理员 API 处理器（环境/实例/通知 CRUD）
│   │   │   ├── services/             # 业务逻辑
│   │   │   ├── routers/              # 路由注册
│   │   │   └── tasks/                # 定时任务：同步 DB 元数据
│   │   ├── users/                    # 用户模块
│   │   │   ├── models/               # 用户/角色/组织模型
│   │   │   ├── forms/views/services/ # 登录/OTP/个人信息/组织/角色 CRUD
│   │   │   └── routers/              # /api/v1/user/login, /api/v1/profile, /api/v1/admin
│   │   ├── das/                      # 数据查询服务（Data Analytics Service）
│   │   │   ├── dao/                  # MySQL/ClickHouse 原生查询
│   │   │   ├── parser/               # TiDB Parser 封装
│   │   │   ├── models/forms/         # 数据模型 + 请求表单
│   │   │   ├── services/             # 查询执行、权限、收藏、历史、字典、表元数据
│   │   │   ├── views/                # /api/v1/das/* 处理器
│   │   │   └── routers/              # 路由：查询/管理两组
│   │   ├── orders/                   # 工单模块
│   │   │   ├── models/               # 工单记录/审批流/任务/日志模型
│   │   │   ├── forms/                # 建单/列表/操作/语法检查/任务
│   │   │   ├── services/             # 提交/审批/认领/执行/生任务/导出
│   │   │   ├── views/                # API 处理器 + WebSocket
│   │   │   └── routers/              # 路由
│   │   └── inspect/                  # SQL 审核模块
│   │       ├── config/               # 审核参数配置
│   │       ├── controllers/
│   │       │   ├── rules/            # 80+ 审核规则（alter.go / create.go / dml.go / drop.go / ...）
│   │       │   ├── traverses/        # AST 遍历器（每个规则一个 Visitor）
│   │       │   ├── logics/           # 审核逻辑（检查参数 → 生成审核结果）
│   │       │   ├── process/          # 审核流程编排
│   │       │   └── hint.go           # 审核结果结构
│   │       ├── models/forms/views/services/routers/  # 标准 4 层
│   │   ├── pkg/                      # 工具包
│   │       ├── kv/cache.go           # 内存 KV 缓存
│   │       ├── notifier/             # 通知推送（企业微信/钉钉/邮件）
│   │       ├── pagination/           # 分页工具
│   │       ├── query/finger.go       # SQL Fingerprint
│   │       ├── response/             # 统一 JSON 响应
│   │       └── utils/                # 加解密/AES-256-GCM/工具函数
│   └── web/web.go                    # embed: dist/
│
├── www/                              # Vue 3 前端
│   ├── index.html
│   ├── vite.config.js                # 代理 /api→localhost:8083
│   ├── package.json
│   ├── src/
│   │   ├── main.js                   # 入口：Pinia/Router/Antd/CodeMirror
│   │   ├── App.vue                   # 根组件：antd 主题 + 中文
│   │   ├── permission.js             # 路由守卫（登录/动态路由/权限）
│   │   ├── router/
│   │   │   ├── index.js              # 路由创建
│   │   │   └── ark.js                # 静态路由 + 异步路由定义（Layout）
│   │   ├── store/
│   │   │   ├── user.js               # 用户状态（Pinia + persist）
│   │   │   ├── permission.js         # 路由生成
│   │   │   └── prefill.js            # 表单预填充
│   │   ├── api/                      # 所有 API 调用
│   │   │   ├── login.js / das.js / order.js / admin.js / profile.js
│   │   ├── utils/
│   │   │   ├── request.js            # Axios 实例 + 拦截器
│   │   │   ├── permission.js         # filterAsyncRoutes()
│   │   │   └── time.js / validate.js
│   │   ├── views/
│   │   │   ├── login/                # 登录 + OTP
│   │   │   ├── das/                  # SQL 查询控制台
│   │   │   ├── orders/               # 工单（创建/列表/详情/任务）
│   │   │   ├── admin/                # 后台管理
│   │   │   │   ├── perms/            # 用户/角色/组织/审批流
│   │   │   │   └── system/           # 环境/实例/审核参数/通知/DAS权限
│   │   │   └── account/              # 个人中心
│   │   └── components/
│   │       ├── layout/               # 侧边栏 + 顶栏布局
│   │       ├── edit/Codemirror.vue   # SQL 编辑器封装
│   │       └── patterns/             # 通用页面组件
│   └── public/
```

---

## 4. 启动与配置

### 4.1 环境要求

| 组件 | 版本 |
|------|------|
| Go | ≥1.22 |
| Node.js | ≥20.19 |
| MySQL | 8.0+ |
| Redis | 可选 |

### 4.2 配置 `backend/config.yaml`

```yaml
app:
    listen_address: "localhost:8083"
    environment: "dev"             # dev / prod
    secret_key: "A9$k!pZ3@rT7&xQ1#Lf8^Vm2*Ws6%Hd0"  # 必须32位，用于 AES/JWT

database:
    driver: mysql
    host: 127.0.0.1
    port: 3306
    database: GoRDS
    username: GoRDS_rw
    password: "1234.Com!"

redis:
    host: 127.0.0.1
    port: 6379
    db: 0
```

> **重要**：首次启动会自动创建所有表（AutoMigrate）并写入种子数据（管理员账号 `admin` / 密码 `admin123`、默认审核参数、允许的操作列表）。

### 4.3 启动命令

```bash
# 启动后端（会自动嵌入前端 dist）
cd backend && go run cmd/main.go

# 单独启动前端 dev server（热更新）
cd www && npm run dev          # 代理到 localhost:8083
```

前端构建后嵌入 Go 二进制，生产部署只需 `go build -o GoRDS cmd/main.go` 即可单文件运行。

---

## 5. 后端架构详解

### 5.1 三层架构

后端严格遵循 **View → Service → Model/DAO** 的分层模式：

```
Gin Router → middleware → View (处理 HTTP) → Service (业务逻辑) → Model/DAO (数据访问)
```

每一层在一个独立的 `.go` 文件中，放在对应模块的 `views/` / `services/` / `models/` / `forms/` 目录下。

### 5.2 启动流程

`cmd/main.go` → `app.Run(configFile)`:

```
bootstrap.InitializeConfig()    # Viper 读 YAML → global.App.Config
bootstrap.InitializeLog()       # Logrus 初始化
bootstrap.InitializeDB()        # GORM 连接 + AutoMigrate + 种子数据
bootstrap.InitializeRedis()     # Redis 连接（可选）
bootstrap.InitializeCron()      # 定时任务
app.RunServer()                 # Gin 引擎 → 路由注册 → 静态文件 → 启动
```

### 5.3 全局对象 `global.App`

在 `internal/global/global.go` 中定义的全局单例，贯穿整个应用：

```go
type Application struct {
    ConfigViper *viper.Viper
    Config      config.Configuration  // 所有 YAML 配置
    JWT         *jwt.GinJWTMiddleware
    Log         *logrus.Logger
    DB          *gorm.DB
    Redis       *redis.Client
    Cron        *cron.Cron
}
```

在任意 service/view 中通过 `global.App.DB` / `global.App.Config` 使用。

### 5.4 路由注册机制

提供了一个插件式的路由注册模式 [`api/api.go`](backend/api/api.go:16)：

```go
// 各模块在 init 或显式调用时注册自己的路由函数
func Include(opts ...Option)      // 收集所有路由函数
func Init() *gin.Engine           // 创建引擎，遍历 options 调用

// 在 app/app.go 中:
api.Include(
    userRouter.Routers,      // → 注册 /api/v1/user/* 等
    commonRouter.Routers,    // → 注册 /api/v1/admin/*（环境/实例/通知）
    inspectRouter.Routers,   // → 注册 /api/v1/admin/inspect/*
    dasRouter.Routers,       // → 注册 /api/v1/das/*, /api/v1/admin/das/*
    ordersRouter.Routers,    // → 注册 /ws/*, /api/v1/orders/*
)
```

每个模块的 `Routers` 方法接收 `*gin.Engine`，在其内部注册自己的路由组。

### 5.5 认证与权限

```
Login: POST /api/v1/user/login
  ├── OTPMiddleware()      → 验证用户密码 → 检查是否需要 OTP
  ├── JWT.LoginHandler()   → 验证 OTP（如有）→ 签发 JWT Token
  └── 返回 {"token": "JWT xxx", "expire": "..."}

后续请求:
  Header: Authorization: JWT <token>
  └── JWT.MiddlewareFunc()      → 解析 JWT，注入 username
  └── HasAdminPermission()      → 检查 insight_users.is_superuser
```

### 5.6 五大核心模块

#### 5.6.1 用户管理 — `internal/users/`

| 路由 | 功能 |
|------|------|
| `POST /api/v1/user/login` | 登录（密码 + OTP） |
| `POST /api/v1/user/logout` | 登出 |
| `POST /api/v1/user/otp-auth-url` | 获取 OTP 绑定二维码 |
| `POST /api/v1/user/otp-auth-callback` | OTP 绑定回调 |
| `GET /api/v1/profile` | 个人信息 |
| `PUT /api/v1/profile` | 修改个人信息 |
| `GET/POST/PUT/DELETE /api/v1/admin/users` | 用户管理 CRUD |
| `GET/POST/PUT/DELETE /api/v1/admin/roles` | 角色管理 |
| `GET/POST/PUT/DELETE /api/v1/admin/organizations` | 组织管理 |

模型关系：
- `InsightUsers` ↔ `InsightOrgs` (多对多，通过 `InsightOrgUsers`)
- `InsightUsers` → `InsightRoles`

#### 5.6.2 DB 实例与环境 — `internal/common/`

管理和发现数据库实例：

- **环境**：`InsightInstanceEnvironments`（如"生产"、"测试"）
- **实例**：`InsightInstances`（主机/端口/密码/DB类型/用途）
  - `use_type`: "查询" 或 "工单"
  - `db_type`: MySQL / TiDB / ClickHouse
  - 密码用 AES-256-GCM 加密存储（`pkg/utils/crypto.go`）
- **Schema**：`InsightInstanceSchemas` — 定时任务自动采集实例下的库列表

#### 5.6.3 数据查询 (DAS) — `internal/das/`

在线执行 SELECT 查询，不可修改数据：

```
/api/v1/das/query/mysql          POST  执行 MySQL/TiDB 查询
/api/v1/das/query/clickhouse     POST  执行 ClickHouse 查询
/api/v1/das/schemas              GET   获取用户的可用库
/api/v1/das/schema/tables        GET   获取库下表
/api/v1/das/schema/grants        GET   获取表级别权限
/api/v1/das/history              GET   查询历史
/api/v1/das/favorites            CRUD  收藏 SQL
/api/v1/das/table-info           GET   表元数据
/api/v1/das/dbdict               GET   数据字典

/api/v1/admin/das/schemas        GET   管理员管理用户Schema权限
/api/v1/admin/das/tables         GET   管理员管理用户表权限
```

核心架构：

```
View（HTTP 处理）
  → Service（ExecuteMySQLQueryService / ExecuteClickHouseQueryService）
    → DAO（dao.DB / dao.ClickhouseDB — 原生 SQL 查询）
    → 结果统一为 []map[string]interface{}
```

- 查询使用 `database/sql` 原生驱动（非 GORM）
- MySQL 驱动：`github.com/go-sql-driver/mysql`
- ClickHouse 驱动：`github.com/clickhouse/clickhouse-go/v2`
- 所有查询记录审计日志到 `InsightDASRecords`

#### 5.6.4 工单系统 — `internal/orders/`

完整的工单生命周期：

```
PENDING → APPROVED → CLAIMED → EXECUTING → COMPLETED
                      ↓                    ↓ (失败)
                    REJECTED             FAILED
                                           ↓
                                         REVIEWED
  (任意状态可撤销) → REVOKED
```

`/api/v1/orders/` 路由概览：

| 路由 | 方法 | 功能 |
|------|------|------|
| `/orders` | POST | 创建工单 |
| `/orders` | GET | 工单列表（分页+搜索） |
| `/orders/:order_id` | GET | 工单详情 |
| `/orders/approvals/:order_id` | GET | 审批状态 |
| `/orders/logs/:order_id` | GET | 操作日志 |
| `/orders/actions/approval` | PUT | 审批 |
| `/orders/actions/claim` | PUT | 认领 |
| `/orders/actions/revoke` | PUT | 撤销 |
| `/orders/actions/transfer` | PUT | 转交 |
| `/orders/actions/complete` | PUT | 完成 |
| `/orders/actions/fail` | PUT | 标记失败 |
| `/orders/actions/review` | PUT | 复核 |
| `/orders/inspect-syntax` | POST | SQL 语法审核 |
| `/orders/tasks` | POST | 生成执行任务 |
| `/orders/tasks/:order_id` | GET | 获取任务列表 |
| `/orders/tasks/execute` | POST | 单条任务执行 |
| `/orders/tasks/execute-batch` | POST | 批量执行 |
| `/orders/tasks/exports/:filename` | GET | 下载导出文件 |

关键模型关系：

```
InsightOrderRecords            工单主表 (1)
  └── InsightOrderTasks        工单任务 (N) — 每一条 SQL 一个 Task
  └── InsightApprovalRecords   审批记录 (N) — 多级审批
  └── InsightOrderLogs         操作日志 (N)
  └── InsightOrderMessages     消息推送记录 (N)

InsightApprovalFlows           审批流定义
  └── InsightApprovalFlowUsers 审批流用户映射
```

审批流定义示例（JSON）：
```json
[
  {"stage":1, "stage_name":"部门审批", "approvers":["zhangsan","lisi"], "type":"OR"},
  {"stage":2, "stage_name":"DBA审批",   "approvers":["dba1","dba2"],    "type":"AND"}
]
```

`OR` 类型：任一审批人通过即可；`AND` 类型：所有人都需要批准。

WebSocket 端点 `/ws/:channel` 用于实时推送工单状态变更。

#### 5.6.5 SQL 审核 — `internal/inspect/`

使用 TiDB Parser 将 SQL 解析为 AST，然后遍历 AST 节点检查各项规范。

处理流程：
```
SQL text
  → TiDB Parser (ast.StmtNode)
    → Process (流程编排，遍历所有规则)
      → Rule[1..N] (检查函数)
        → Traverse (AST Visitor 遍历)
          → Logic (检查参数 → 生成 Hint)

返回审核结果: []InspectHint{Level, Message}
```

80+ 审核规则分布在 `controllers/rules/` 下：

| 文件 | 规则类别 | 示例 |
|------|---------|------|
| `create.go` | 建表/建索引 | 表名长度≤32、必须有主键、检查字符集 |
| `alter.go` | 表变更 | 是否允许 DROP 列、是否合并 ALTER |
| `dml.go` | 数据操作 | 必须有 WHERE、禁止子查询、限制影响行数 |
| `drop.go` | 删除 | 是否允许 DROP TABLE/TRUNCATE |
| `database.go` | 库操作 | 建库规范 |
| `rename.go` | 重命名 | 是否允许 RENAME TABLE |
| `view.go` | 视图 | 是否允许 CREATE VIEW |
| `analyze.go` | 分析表 | ANALYZE TABLE 规范 |
| `rule.go` | 规则结构定义 | Rule struct + CheckFunc |

审核参数两层覆盖：
1. **全局参数** (`InsightInspectGlobalParams`) — 系统级默认值
2. **实例参数** (`InsightInspectInstanceParams`) — 覆盖特定实例的规则

### 5.7 工具包

- **加解密** `pkg/utils/crypto.go` — AES-256-GCM，密钥来自 `config.app.secret_key`
- **通知推送** `pkg/notifier/` — 企业微信机器人 / 钉钉机器人 / SMTP 邮件
- **分页** `pkg/pagination/` — 基于 GORM 的通用分页
- **响应格式** `pkg/response/` — 统一 `{code, data, message, request_id}` 格式
- **SQL 指纹** `pkg/query/finger.go` — 去除参数值的 SQL 归一化

---

## 6. 前端架构详解

### 6.1 路由系统

```
/login          → Login.vue          （无需登录）
/403            → 403.vue

/               → Layout.vue（侧边栏+顶栏布局，需要登录）
  ├─ /das       → SQL查询（DAS控制台）
  ├─ /orders    → 工单系统
  │   ├─ /orders/create            新建工单
  │   ├─ /orders                   工单列表
  │   ├─ /orders/:order_id         工单详情
  │   ├─ /orders/tasks/:order_id   工单任务
  │   └─ /orders/tasks/exports/:filename  下载页面
  ├─ /account   → 个人中心
  └─ /admin     → 后台管理（需要 is_superuser）
      ├─ /admin/users     用户管理
      ├─ /admin/roles     角色管理
      ├─ /admin/orgs      组织管理
      ├─ /admin/flows     审批流
      ├─ /admin/environment   环境管理
      ├─ /admin/instance      实例配置
      ├─ /admin/inspect       审核参数
      ├─ /admin/notify        消息通知
      └─ /admin/das           DAS 权限
```

动态路由由 [`permission.js`](www/src/permission.js:1) 生成：

1. 首次访问时加载用户信息
2. 调用 `permissionStore.GenerateRoutes(is_superuser)` 
3. `filterAsyncRoutes()` 过滤出需要管理员权限的路由
4. 通过 `router.addRoute()` 动态添加

### 6.2 状态管理（Pinia）

- **userStore**: 用户信息 + Token + 登录状态判断
- **permissionStore**: 已生成路由列表 + 菜单路由
- **prefillStore**: 工单表单预填充数据

Token 在 localStorage 中持久化，通过 Axios 请求拦截器自动注入 `Authorization: JWT <token>`。

### 6.3 统一请求层

`utils/request.js` 封装 Axios：
- 请求拦截器：注入 JWT Token、启动 NProgress
- 响应拦截器：统一解析 `code !== '0000'` 为错误、401 自动跳转登录页

### 6.4 UI 组件

- **Layout.vue**: 侧边栏菜单 + 顶栏（用户信息/全屏/退出）
- **Codemirror.vue**: SQL 编辑器封装（语法高亮、自动补全、格式化）
- **PageCardShell/PageTableSection/PageToolbar**: 通用列表页面模式

---

## 7. 核心业务流程

### 7.1 工单提交流程

```
开发者提交工单
  → 前端语法审核 (inspect-syntax)
    → TiDB Parser 解析 SQL → 遍历 80+ 规则 → 返回审核结果
      ↓ 审核不通过 → 提示修改
      ↓ 审核通过
  → 创建工单 (POST /orders)
    → 写入 InsightOrderRecords
    → 创建审批记录 InsightApprovalRecords（根据审批流定义）
    → 通知审批人
  → 审批人审批
    → AND/OR 多级审批 → 全部通过后工单状态变为 APPROVED
  → 执行人认领
    → 生成执行任务 InsightOrderTasks
    → 逐条执行 SQL（DDL 使用 gh-ost 在线变更 / DML 直接执行）
    → 执行完成后自动复核
  → 工单完成
```

### 7.2 DDL 安全执行

DDL 变更通过 **gh-ost**（GitHub Online Schema Migration）执行，配置在 `config.yaml` 的 `ghost` 段：

```yaml
ghost:
    path: "/usr/local/bin/gh-ost"
    args: ["--allow-on-master", "--assume-rbr", ...]
```

### 7.3 定时任务

| 任务 | 周期 | 功能 |
|------|------|------|
| `SyncDBMeta` | 默认每5分钟 | 自动采集所有实例的库表结构到 `InsightInstanceSchemas` |
| `KillTiDBQuery` | 默认每5分钟 | 清理 TiDB 超时查询 |

---

## 8. 开发指南

### 8.1 新增一个 API 端点

以在 `orders` 模块添加一个"工单标签"功能为例：

**1. Forms — 请求结构体** (`internal/orders/forms/`)
```go
type AdminTagForm struct {
    ID   uint64 `form:"id"`
    Name string `form:"name" binding:"required"`
}
```

**2. Models — 数据模型** (`internal/orders/models/`)
```go
type InsightOrderTags struct {
    *models.Model
    OrderID uuid.UUID `gorm:"type:char(36);index"`
    Tag     string    `gorm:"type:varchar(32)"`
}
```
然后在 `bootstrap/db.go` 的 `initializeTables` 中添加 `&ordersModels.InsightOrderTags{}`。

**3. Services — 业务逻辑** (`internal/orders/services/`)
```go
type AdminTagService struct {
    *forms.AdminTagForm
    C *gin.Context
}

func (s *AdminTagService) Run() error {
    return global.App.DB.Create(&ordersModels.InsightOrderTags{...}).Error
}
```

**4. Views — HTTP 处理器** (`internal/orders/views/`)
```go
func AdminCreateTagView(c *gin.Context) {
    var form forms.AdminTagForm
    if err := c.ShouldBind(&form); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }
    service := services.AdminTagService{AdminTagForm: &form, C: c}
    if err := service.Run(); err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, nil, "创建成功")
}
```

**5. Routers — 路由注册** (`internal/orders/routers/`)
```go
admin.POST("/tags", views.AdminCreateTagView)
```

**6. 前端 API** (`www/src/api/order.js`)
```js
export const createOrderTagApi = (data) => post('/api/v1/admin/approval-flows/tags', data)
```

### 8.2 新增一个审核规则

**1. Traverser — AST 遍历器** (`internal/inspect/controllers/traverses/`)
```go
type TraverseCheckNoExistEnum struct {
    Exist bool  // 是否检查通过
}

func (v *TraverseCheckNoExistEnum) Enter(in ast.Node) (ast.Node, bool) {
    // 根据 in.(type) 判断节点类型
    return in, false
}

func (v *TraverseCheckNoExistEnum) Leave(in ast.Node) (ast.Node, bool) {
    return in, true
}
```

**2. Logic — 审核逻辑** (`internal/inspect/controllers/logics/`)
```go
func LogicCheckNoExistEnum(v *traverses.TraverseCheckNoExistEnum, hint *controllers.RuleHint) {
    if !v.Exist {
        hint.HandleError("不允许使用xxx语法")
    }
}
```

**3. Rule — 注册规则** (`internal/inspect/controllers/rules/`)
```go
func DDLRules() []Rule {
    return append(DDLRules(), Rule{
        Hint:      "DDL#不允许使用xxx",
        CheckFunc: (*Rule).RuleCheckNoExistEnum,
    })
}

func (r *Rule) RuleCheckNoExistEnum(tistmt *ast.StmtNode) {
    v := &traverses.TraverseCheckNoExistEnum{}
    (*tistmt).Accept(v)
    logics.LogicCheckNoExistEnum(v, r.RuleHint)
}
```

**4. 审核参数** — 如需可配置，在 `bootstrap/db.go` 的 `initializeGlobalInspectParams` 中添加。

### 8.3 编码约定

- **View 只做 HTTP 处理**：解析参数、调用 Service、返回响应。不写业务逻辑。
- **Service 只做业务逻辑**：调用 DB、调用外部服务。不处理 HTTP 细节。
- **Model 只定义数据结构**：TableName、字段标签、钩子（BeforeCreate）。不写业务方法。
- **所有 API 响应使用 `response.Success/Fail`**：统一 `{code, data, message, request_id}`。
- **请求参数使用独立的 Forms 结构体**：不使用 Model 直接绑定。
- **日志使用 `global.App.Log`**：区分 Info/Error/Fatal，关联 request_id。
- **密码和敏感信息使用 AES-256-GCM 加密**：`utils.Encrypt()` / `utils.Decrypt()`。

### 8.4 常见开发场景

| 场景 | 操作 |
|------|------|
| 新建数据库表 | Models 定义结构 → 添加到 `initializeTables()` → `go run` 自动迁移 |
| 添加配置项 | `config/config.go` 添加字段 → `config.yaml` 添加值 |
| 增加前端页面 | `src/views/` 下新建 → `route.js` 导出 → `ark.js` 引用 |
| 修改前端 API | `src/api/` 下修改 → 调用 `get/post/put/del` 工具函数 |
| 新增通知渠道 | `pkg/notifier/` 新文件 → `common/services/notify.go` 添加测试发送 |
| 构建部署 | `cd www && npm run build` → `cd backend && go build -o GoRDS cmd/main.go` |