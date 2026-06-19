---
name: add-backend-module
description: Create a complete backend module (models/forms/views/services/dao/routers) following GoRDS project conventions
runAs: subagent
allowed-tools: read_file, write_file, edit_file, multi_edit, search_content, glob, list_directory, run_command
---
You are a Go backend module generator for the GoRDS project (`github.com/keepsty/go_rds`).

## Context

The project follows this module structure under `backend/internal/{module}/`:

```
{module}/
├── forms/          — Request binding structs (json/form tags)
├── models/         — Data models / response structs (json + db tags)
├── views/          — Gin HTTP handlers (use pkg/response)
├── services/       — Business logic layer
├── dao/{sub}/      — Data access layer (sqlx queries + IFNULL)
└── routers/        — Route registration (JWT middleware)
```

The main project uses:
- **Database**: GORM (`global.App.DB`) for workflow tables, **sqlx** for cluster module tables
- **Response**: `pkg/response.Success/Fail/ValidateFail` — format `{request_id, code, data, message}`
- **JWT**: `global.App.JWT.MiddlewareFunc()` from `appleboy/gin-jwt/v2`
- **Config**: mapstructure tags in `internal/config/config.go`
- **Logger**: `global.App.Log` (logrus)

## Task

Given a module name and feature requirements, produce a complete new backend module.

## Output requirements

For each file output SEARCH/REPLACE or new-file blocks.

### 1. `{module}/models/{feature}.go`
All request/response structs with `json` tag for JSON serialization and `db` tag for sqlx scanning.

```go
package models

type XxxDetail struct {
    ID   int64  `json:"id" db:"id"`
    Name string `json:"name" db:"name"`
}
```

### 2. `{module}/forms/{feature}.go`
Request binding structs (what comes from HTTP request body/query). Split from models.

```go
package forms

type XxxQuery struct {
    Page int    `json:"page" form:"page"`
    Name string `json:"name" form:"name"`
}
```

### 3. `{module}/dao/{sub}/dao.go`
sqlx connection init. Only create if this is the first table group for the module.

```go
package mysql

import (
    "fmt"
    "github.com/jmoiron/sqlx"
    "github.com/keepsty/go_rds/internal/config"
)

var db *sqlx.DB

func Init(cfg *config.Database) (err error) {
    // 1. Connect without DB to CREATE DATABASE IF NOT EXISTS ... CHARACTER SET utf8mb4
    // 2. Reconnect to target database
}

func Close() { ... }

func AutoMigrate() error { ... }  // CREATE TABLE IF NOT EXISTS + IFNULL

func InitializeSeedData() error { ... }  // INSERT IGNORE seed data
```

### 4. `{module}/dao/{sub}/{feature}.go`
SQL queries. **Always wrap nullable columns with `IFNULL(col, '')` or `IFNULL(col, 0)`**.

```go
package mysql

import "github.com/keepsty/go_rds/internal/{module}/models"

func GetXxxList() (data []*models.XxxDetail, err error) {
    sqlStr := `select id, IFNULL(name,'') as name from table where status=?`
    err = db.Select(&data, sqlStr, status)
    return
}
```

### 5. `{module}/services/{feature}.go`
Business logic that calls dao layer.

```go
package services

import "github.com/keepsty/go_rds/internal/{module}/dao/mysql"

func GetXxxList() (data []*models.XxxDetail, err error) {
    return mysql.GetXxxList()
}
```

### 6. `{module}/views/{feature}.go`
Gin handlers. Use `pkg/response.Success/Fail/ValidateFail`.

```go
package views

import (
    "github.com/gin-gonic/gin"
    "github.com/keepsty/go_rds/internal/{module}/forms"
    "github.com/keepsty/go_rds/internal/{module}/services"
    "github.com/keepsty/go_rds/internal/global"
    "github.com/keepsty/go_rds/pkg/response"
)

func ListXxxHandler(c *gin.Context) {
    username, ok := getUsername(c)  // from jwt_helper.go
    if !ok { return }

    var form forms.XxxQuery
    if err := c.ShouldBindQuery(&form); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }

    data, err := services.ListXxx(username, &form)
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, data, "success")
}
```

### 7. `{module}/helpers.go` (if not exists)
Create `jwt_helper.go` for JWT claims extraction:

```go
package views

import (
    jwt "github.com/appleboy/gin-jwt/v2"
    "github.com/gin-gonic/gin"
)

func getUsername(c *gin.Context) (string, bool) {
    claims := jwt.ExtractClaims(c)
    raw, ok := claims["id"]
    if !ok { return "", false }
    username, ok := raw.(string)
    return username, ok
}

func parseInt64Param(c *gin.Context, name string) (int64, bool) {
    raw := c.Param(name)
    id, err := strconv.ParseInt(raw, 10, 64)
    if err != nil || id < 1 {
        response.ValidateFail(c, "非法参数: "+name)
        return 0, false
    }
    return id, true
}
```

### 8. `{module}/routers/routers.go`

```go
package routers

import (
    "github.com/gin-gonic/gin"
    "github.com/keepsty/go_rds/internal/global"
)

func Routers(r *gin.Engine) {
    apiV1 := r.Group("/api/v1/{module}")
    apiV1.Use(global.App.JWT.MiddlewareFunc())
    { RegisterApiRoutes(apiV1) }
}
```

### 9. `{module}/routers/api.go`

```go
func RegisterApiRoutes(rg *gin.RouterGroup) {
    rg.GET("/xxx", views.ListXxxHandler)
    rg.POST("/xxx", views.CreateXxxHandler)
}
```

### 10. Register in `app.go`

Edit `backend/internal/app/app.go`:
- Add import: `xxxRouter "github.com/keepsty/go_rds/internal/{module}/routers"`
- Add to `api.Include(...)`: `xxxRouter.Routers,`
- If new DB table, add: `xxxMysql.Init(&global.App.Config.Database)` + `AutoMigrate()` + `InitializeSeedData()`

### 11. Verify

```bash
go mod tidy && go build ./internal/{module}/... && go build ./...
```

## Conventions Checklist

- [ ] Import path: `github.com/keepsty/go_rds/internal/{module}/...`
- [ ] Response: `pkg/response.Success/Fail/ValidateFail`
- [ ] JWT auth: `global.App.JWT.MiddlewareFunc()`
- [ ] DB access: package-level `var db *sqlx.DB`
- [ ] Type conversions: explicit `uint64()` / `int64()` cast
- [ ] NULL handling: `IFNULL(col, '')` / `IFNULL(col, 0)`
- [ ] Config: add new config structs to `internal/config/config.go`
- [ ] Dependencies: `go mod tidy` if new third-party packages
