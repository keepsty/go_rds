---
name: go-rds-backend-structure
description: >-
  GoRDS 后端代码结构与模块模式参考。当用户需要新增后端模块、API 接口，或询问项目架构问题时，**必须**使用此技能。涵盖 routers → views → services → models 的完整 CRUD 流程、路由注册、中间件、全局对象、响应封装等规范。适用于：
  (1) 创建新的业务模块（如 reports、audit 等），
  (2) 为已有模块增加 CRUD 接口，
  (3) 理解请求生命周期和代码组织方式，
  (4) 遵循项目约定写代码而非猜测。
  本技能描述的是 gin-gonic + GORM 的项目模式，推送到 AI coding agent 时优先参考。
compatibility:
  - go 1.21+
  - gin-gonic
  - gorm
  - gin-jwt/v2
---

# GoRDS 后端代码结构指南

该 Go 项目采用 **Gin (Go)** + **GORM** 框架，遵循固定的模块化分层。每个业务模块（users、common、das、orders、inspect、cluster）独立为 `internal/{模块}/` 目录，依赖流向：`routers/ → views/ → services/ → models/`。

```
backend/
├── cmd/main.go            # 入口
├── api/api.go             # Gin 引擎工厂（CORS、request-id、日志）
├── middleware/             # JWT、权限、OTP、请求日志
├── pkg/                   # response、pagination、utils
├── internal/
│   ├── app/app.go         # 应用组装
│   ├── bootstrap/         # config → log → DB → Redis → Cron
│   ├── config/config.go   # Viper 配置映射
│   ├── global/global.go   # Application 全局单例
│   ├── common/            # 公共模块（实例、环境、通知、基础模型）
│   ├── users/das/orders/inspect/cluster/  # 业务模块
│   └── {新模块}/          # 新增模块遵循相同结构
```

---

## 按需加载子 skill

根据当前任务，只加载需要的部分：

| 任务场景 | 需加载的文件 |
|----------|------------|
| **新增完整模块** | 主 SKILL.md + `references/startup.md` + `references/router.md` + `references/views.md` + `references/forms.md` + `references/services.md` + `references/models.md` |
| **加一个 CRUD 接口** | `references/views.md` + `references/forms.md` + `references/services.md` |
| **加一个新数据库表** | `references/models.md` + `references/startup.md`（仅 auto-migrate 部分） |
| **修改响应格式** | `references/response.md` |
| **了解请求流程/架构** | `references/startup.md` + `references/middleware.md` |
| **配置中间件** | `references/middleware.md` |
| **路由注册** | `references/router.md` |

> 提示：如果不确定加载哪个，直接加载 `references/router.md` + `references/views.md` + `references/forms.md` + `references/services.md` + `references/models.md` 这五个核心文件即可覆盖 90% 场景。

---

## 新增模块速查（10 步）

```bash
mkdir -p internal/{模块}/{routers,views,services,forms,models,tasks}
```

1. **models** — 定义 GORM 模型 + `TableName()`，继承 `*commonModels.Model`
2. **bootstrap/db.go** — 在 `initializeTables()` 注册模型
3. **forms** — 定义请求参数结构体 + binding tag
4. **services** — Service 结构体 + `Run()` 方法
5. **views** — handler 函数（绑定 form → 调用 Service → 返回 response）
6. **views/helpers.go** — `parseUint64Param` / `getUsername` 辅助函数
7. **routers** — `Routers()` + 子路由注册 + 中间件配置
8. **app.go** — `api.Include(yourModuleRouter.Routers)` 注册
9. （可选）种子数据到 `bootstrap/db.go`
10. **更新 CLAUDE.md**

---

## 黄金规则

- 表名统一 `insight_` 前缀
- 全部响应走 `response.Success/Fail/ValidateFail/PaginationSuccess`，禁止直接 `c.JSON()`
- `app.secret_key` 必须恰好 32 字符
- 模型必须嵌入 `*commonModels.Model` 获取 `ID/CreatedAt/UpdatedAt`
