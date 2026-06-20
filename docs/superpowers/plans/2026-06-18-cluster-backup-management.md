# 集群备份管理实现计划

> **面向 AI 代理的工作者：** 使用 subagent-driven-development 逐任务实现。

**目标：** 在集群管理中增加 MySQL/Redis/TiDB 备份管理页面，后端存储配置、任务、记录到独立备份表。

**架构：** 后端沿用 cluster 模块的 sqlx + DAO 模式，新增 `backup_config / backup_task / backup_record` 三张表；前端作为 cluster 子页面，按数据库类型分 tab 展示。

**技术栈：** Go (sqlx) + Vue 3 + Ant Design Vue 4.x

---

## 后端设计

### 新增数据库表

```sql
-- 备份配置
CREATE TABLE database_backup_config (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(128) NOT NULL COMMENT '配置名称',
  db_type VARCHAR(32) NOT NULL COMMENT '数据库类型: MySQL/Redis/TiDB',
  instance_id VARCHAR(128) NOT NULL COMMENT '关联实例标识',
  backup_type VARCHAR(32) NOT NULL DEFAULT 'full' COMMENT '备份类型: full/incremental',
  schedule_cron VARCHAR(64) DEFAULT '' COMMENT '定时策略(cron表达式)',
  retention_days INT NOT NULL DEFAULT 7 COMMENT '保留天数',
  storage_path VARCHAR(256) DEFAULT '' COMMENT '存储路径',
  status VARCHAR(32) NOT NULL DEFAULT 'enabled' COMMENT '状态: enabled/disabled',
  remark TEXT COMMENT '备注',
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP,
  update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 备份任务
CREATE TABLE database_backup_task (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  config_id BIGINT NOT NULL COMMENT '关联配置ID',
  name VARCHAR(128) NOT NULL COMMENT '任务名称',
  db_type VARCHAR(32) NOT NULL COMMENT '数据库类型',
  backup_type VARCHAR(32) NOT NULL DEFAULT 'full' COMMENT '备份类型',
  instance_id VARCHAR(128) NOT NULL COMMENT '实例标识',
  status VARCHAR(32) NOT NULL DEFAULT 'pending' COMMENT '状态: pending/running/success/failed',
  started_at DATETIME COMMENT '开始时间',
  finished_at DATETIME COMMENT '完成时间',
  created_by VARCHAR(64) NOT NULL COMMENT '创建人',
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 备份记录
CREATE TABLE database_backup_record (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  task_id BIGINT NOT NULL COMMENT '关联任务ID',
  file_name VARCHAR(256) DEFAULT '' COMMENT '备份文件名',
  file_size BIGINT DEFAULT 0 COMMENT '文件大小(字节)',
  file_path VARCHAR(512) DEFAULT '' COMMENT '文件路径',
  status VARCHAR(32) NOT NULL DEFAULT 'running' COMMENT '状态: running/success/failed',
  output TEXT COMMENT '执行输出',
  started_at DATETIME COMMENT '开始时间',
  finished_at DATETIME COMMENT '完成时间',
  create_time DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 新增后端文件

| 文件 | 说明 |
|------|------|
| `internal/cluster/models/backup.go` | 3 个模型结构体 |
| `internal/cluster/forms/backup.go` | 请求参数表单 |
| `internal/cluster/dao/mysql/backup.go` | DAO 查询方法 |
| `internal/cluster/services/backup.go` | 业务逻辑层 |
| `internal/cluster/views/backup.go` | HTTP handler |
| `internal/cluster/routers/api.go` | 追加备份路由 |

### 修改的后端文件

| 文件 | 变更 |
|------|------|
| `internal/cluster/dao/mysql/migrate.go` | 新增 3 张表的 auto-migrate |

### API 路由

```
# 新增到 cluster 模块的 router/api.go
GET    /clusters/backup/configs          → 备份配置列表
POST   /clusters/backup/configs          → 创建备份配置
PUT    /clusters/backup/configs/:id      → 更新备份配置
DELETE /clusters/backup/configs/:id      → 删除备份配置

GET    /clusters/backup/tasks            → 备份任务列表
POST   /clusters/backup/tasks            → 创建备份任务
PUT    /clusters/backup/tasks/:id/status → 更新任务状态

GET    /clusters/backup/records          → 备份记录列表
```

### 后端代码模板

```go
// models/backup.go
type BackupConfig struct {
    ID            int64  `db:"id" json:"id"`
    Name          string `db:"name" json:"name"`
    DbType        string `db:"db_type" json:"db_type"`
    InstanceID    string `db:"instance_id" json:"instance_id"`
    BackupType    string `db:"backup_type" json:"backup_type"`
    ScheduleCron  string `db:"schedule_cron" json:"schedule_cron"`
    RetentionDays int    `db:"retention_days" json:"retention_days"`
    StoragePath   string `db:"storage_path" json:"storage_path"`
    Status        string `db:"status" json:"status"`
    Remark        string `db:"remark" json:"remark"`
    CreateTime    string `db:"create_time" json:"create_time"`
    UpdateTime    string `db:"update_time" json:"update_time"`
}

type BackupTask struct {
    ID         int64  `db:"id" json:"id"`
    ConfigID   int64  `db:"config_id" json:"config_id"`
    Name       string `db:"name" json:"name"`
    DbType     string `db:"db_type" json:"db_type"`
    BackupType string `db:"backup_type" json:"backup_type"`
    InstanceID string `db:"instance_id" json:"instance_id"`
    Status     string `db:"status" json:"status"`
    StartedAt  string `db:"started_at" json:"started_at"`
    FinishedAt string `db:"finished_at" json:"finished_at"`
    CreatedBy  string `db:"created_by" json:"created_by"`
    CreateTime string `db:"create_time" json:"create_time"`
}

type BackupRecord struct {
    ID         int64  `db:"id" json:"id"`
    TaskID     int64  `db:"task_id" json:"task_id"`
    FileName   string `db:"file_name" json:"file_name"`
    FileSize   int64  `db:"file_size" json:"file_size"`
    FilePath   string `db:"file_path" json:"file_path"`
    Status     string `db:"status" json:"status"`
    Output     string `db:"output" json:"output"`
    StartedAt  string `db:"started_at" json:"started_at"`
    FinishedAt string `db:"finished_at" json:"finished_at"`
    CreateTime string `db:"create_time" json:"create_time"`
}
```

---

## 前端设计

### 路由结构

在 `www/src/views/cluster/route.js` 中新增子路由：

```javascript
{
  name: 'view.cluster.backup',
  path: '/cluster/backup',
  component: () => import('./backup/index.vue'),
  meta: { title: '备份管理', keepAlive: true },
}
```

### 页面结构

```
/cluster/backup/
├── index.vue                    ← 容器页，Tabs 切换数据库类型
│   ├── config-tab.vue           ← 备份配置列表（按 db_type 筛选）
│   ├── task-tab.vue             ← 备份任务列表
│   └── record-tab.vue           ← 备份记录列表
```

### 页面功能

**备份配置 Tab：**
- 表格展示：名称、数据库类型、实例、备份类型、保留天数、状态
- 操作：新增/编辑/删除 配置
- 数据库类型标签：MySQL(蓝)、Redis(绿)、TiDB(紫)
- 新增/编辑弹窗：名称、选择类型、实例ID、备份类型、定时cron、保留天数、存储路径

**备份任务 Tab：**
- 表格展示：任务名称、数据库类型、备份类型、状态、创建人、时间
- 关联查看（可跳到对应配置或记录）

**备份记录 Tab：**
- 表格展示：备份文件名、大小、状态、开始/完成时间
- 可点击查看执行输出

---

## 要创建/修改的文件清单

### 后端

| 操作 | 文件 |
|------|------|
| 创建 | `internal/cluster/models/backup.go` |
| 创建 | `internal/cluster/forms/backup.go` |
| 创建 | `internal/cluster/dao/mysql/backup.go` |
| 创建 | `internal/cluster/services/backup.go` |
| 创建 | `internal/cluster/views/backup.go` |
| 修改 | `internal/cluster/dao/mysql/migrate.go` |
| 修改 | `internal/cluster/routers/api.go` |

### 前端

| 操作 | 文件 |
|------|------|
| 修改 | `www/src/api/cluster.js` |
| 修改 | `www/src/views/cluster/route.js` |
| 创建 | `www/src/views/cluster/backup/index.vue` |

---

## 验证方法

1. `go build ./internal/cluster/...`
2. `npm run build`
3. 启动后访问 `/cluster/backup` → 看到三个 Tab 页（配置/任务/记录）
4. 测试增删改查流程