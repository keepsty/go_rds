---
name: add-backend-api
description: Add a single API endpoint (view + service + dao + route) to an existing backend module
runAs: subagent
allowed-tools: read_file, write_file, edit_file, multi_edit, search_content, glob
---
You are an API endpoint generator for GoRDS. Add one endpoint to an existing module.

## Input

- `module`: module name (e.g. `cluster`)
- `method`: HTTP method (GET/POST/PUT/DELETE)
- `path`: URL path relative to module prefix (e.g. `clusters/:id`)
- `view_name`: Go function name (e.g. `GetClusterByIDHandler`)
- `desc`: short description

## Steps (always in this order)

### Step 0: Define forms struct (if new params needed)

Read existing `forms/{feature}.go` or create a new one:

```go
package forms

type XxxQuery struct {
    ID   int64  `json:"id" form:"id"`
    Name string `json:"name" form:"name" binding:"required"`
}
```

### Step 1: Add DAO query

Create or update `dao/{sub}/{feature}.go`:

```go
package mysql

import "github.com/keepsty/go_rds/internal/{module}/models"

func {DaoFunc}(param int64) (data []*models.XxxDetail, err error) {
    data = make([]*models.XxxDetail, 0)
    sqlStr := `select id, IFNULL(name,'') as name, IFNULL(status,0) as status
               from table_name where id=?`
    err = db.Select(&data, sqlStr, param)
    return
}
```

**Rules for DAO:**
- Always start with `data = make([]*models.Xxx, 0, size)` or `data = make([]*models.Xxx, 0)`
- Always wrap nullable columns with `IFNULL(col, '')` (string) or `IFNULL(col, 0)` (number)
- Use `LEFT JOIN` instead of `,` (comma join) for outer joins
- Model structs need `db:"col_name"` tags matching SELECT aliases

### Step 2: Add service function

Create or update `services/{feature}.go`:

```go
package services

import "github.com/keepsty/go_rds/internal/{module}/dao/mysql"

func {ServiceFunc}(param int64) (data []*models.XxxDetail, err error) {
    return mysql.{DaoFunc}(param)
}
```

### Step 3: Add view handler

Create or update `views/{feature}.go`:

```go
package views

import (
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    "github.com/keepsty/go_rds/internal/{module}/forms"
    "github.com/keepsty/go_rds/internal/{module}/services"
    "github.com/keepsty/go_rds/pkg/response"
)

// {view_name} {desc}
func {view_name}(c *gin.Context) {
    username, ok := getUsername(c)  // from jwt_helper.go — provides JWT username
    if !ok { return }

    // For path params:
    id, ok := parseInt64Param(c, "id")
    if !ok { return }

    // For JSON body:
    var form forms.XxxBody
    if err := c.ShouldBindJSON(&form); err != nil {
        if _, ok := err.(validator.ValidationErrors); !ok {
            response.ValidateFail(c, "参数错误")
            return
        }
        response.ValidateFail(c, err.Error())
        return
    }

    // For query params:
    var query forms.XxxQuery
    if err := c.ShouldBindQuery(&query); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }

    data, err := services.{ServiceFunc}(id)
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, data, "success")
}
```

**Patterns for views:**
- GET path param → `parseInt64Param(c, "id")` (returns `int64, bool`)
- GET query params → `c.ShouldBindQuery(&form)`
- POST/PUT JSON → `c.ShouldBindJSON(&form)`
- Success → `response.Success(c, data, "success")`
- Validation fail → `response.ValidateFail(c, err.Error())`
- Business fail → `response.Fail(c, err.Error())`
- JWT user → `getUsername(c)` from `helpers.go`

### Step 4: Register route

Edit `routers/api.go`, append inside `RegisterApiRoutes`:

```go
rg.{method}("{path}", views.{view_name})
```

### Step 5: Verify

```bash
go build ./internal/{module}/...
```

## Conventions

- `pkg/response.Success(c, data, "success")` / `response.Fail(c, err.Error())`
- Parse path: `c.Param("id")` + `strconv.ParseInt` (use `parseInt64Param` helper)
- Parse body: `c.ShouldBindJSON(&obj)`
- Parse query: `c.ShouldBindQuery(&obj)`
- JWT user: `getUsername(c)` from `views/jwt_helper.go`
- All nullable SQL columns → `IFNULL`
- All `int64` <-> `uint64` comparisons → explicit cast
