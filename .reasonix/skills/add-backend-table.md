---
name: add-backend-table
description: Add a new database table with model struct + AutoMigrate DDL + seed data to a sqlx-based backend module
runAs: subagent
allowed-tools: read_file, write_file, edit_file, multi_edit, search_content
---
You are a DB migration generator for GoRDS. Add a new table with its Go model struct, DDL, and optional seed data.

## Input

- `module`: module name (e.g. `cluster`)
- `table_name`: MySQL table name (e.g. `my_new_table`)
- `comment`: table comment
- `columns`: list of `{name, type, go_type, db_tag, nullable, default, comment}`
  - type: `BIGINT`, `VARCHAR(N)`, `INT`, `DECIMAL(20,4)`, `TEXT`, `TINYINT`, `DATETIME`
  - go_type: `int64`, `string`, `float64`, `time.Time`
  - db_tag: column name for sqlx scanning
- `indexes`: optional list of `{name, columns}` e.g. `idx_user_id (user_id)`
- `seed_data`: optional list of INSERT rows

## Steps

### Step 0: Create Go model struct

Create or update `models/{feature}.go`:

```go
package models

import "time"

type TableName struct {
    ID         int64     `json:"id" db:"id"`
    Name       string    `json:"name" db:"name"`
    Status     int64     `json:"status" db:"status"`
    CreateTime time.Time `json:"create_time" db:"create_time"`
    UpdateTime time.Time `json:"update_time" db:"update_time"`
}
```

Naming:
- Go struct: PascalCase singular of table_name (e.g. `my_new_table` → `MyNewTable`)
- JSON tags: snake_case
- db tags: same as MySQL column name

### Step 1: Add AutoMigrate DDL

Read `dao/mysql/migrate.go`, find the `tables` slice in `AutoMigrate()`, append:

```go
{
    name: "{table_name}",
    sql: `CREATE TABLE IF NOT EXISTS {table_name} (
        id BIGINT AUTO_INCREMENT PRIMARY KEY COMMENT '主键ID',
        name VARCHAR(128) NOT NULL DEFAULT '' COMMENT '名称',
        status INT NOT NULL DEFAULT 0 COMMENT '状态 0:off 1:on',
        create_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
        update_time DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
        INDEX idx_status (status)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='{comment}'`,
},
```

**Column type conventions:**

| MySQL type | Go type | Go zero-safe default |
|---|---|---|
| `BIGINT` | `int64` | `0` |
| `INT` | `int` / `int64` | `0` |
| `VARCHAR(N)` | `string` | `''` |
| `DECIMAL(20,4)` | `float64` | `0.0` |
| `TEXT` | `string` | `''` |
| `TINYINT` | `int8` / `bool` | `0` / `false` |
| `DATETIME` | `time.Time` | — |
| `TIMESTAMP` (自动更新) | `time.Time` | `CURRENT_TIMESTAMP` |

- 业务时间用 `DATETIME`，系统时间戳用 `TIMESTAMP`
- 所有字符串字段设置 `NOT NULL DEFAULT ''`
- 所有数字字段设置 `NOT NULL DEFAULT 0`
- `id` 总是 `BIGINT AUTO_INCREMENT PRIMARY KEY`
- 时间戳字段总是 `create_time` + `update_time`

### Step 2: Add seed data (optional)

Read `dao/mysql/seed.go`, find `InitializeSeedData()`, append:

```go
{
    name: "{table_name}",
    sql: `INSERT IGNORE INTO {table_name} (name, status) VALUES
        ('示例A', 1),
        ('示例B', 0)`,
},
```

### Step 3: Verify

```bash
go build ./internal/{module}/...
```
