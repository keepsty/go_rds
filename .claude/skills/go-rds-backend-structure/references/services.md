# 服务层（services/）

服务层负责**业务逻辑**，与 `forms`、`models`、`global.App.DB` 交互。

## 核心模式

每个业务操作定义一个结构体，嵌入对应 form 指针 + Gin Context，通过 `Run()` 方法执行。

## CRUD 模板

```go
package services

import (
    "errors"
    "fmt"
    "github.com/keepsty/go_rds/internal/global"
    "github.com/keepsty/go_rds/pkg/pagination"
    "github.com/keepsty/go_rds/internal/{模块}/forms"
    "github.com/keepsty/go_rds/internal/{模块}/models"
    "github.com/gin-gonic/gin"
    "github.com/go-sql-driver/mysql"
)

// ── 列表 ──
type GetXxxService struct {
    *forms.GetXxxForm
    C *gin.Context
}

func (s *GetXxxService) Run() (responseData interface{}, total int64, err error) {
    var list []models.XxxModel
    tx := global.App.DB.Model(&models.XxxModel{})
    if s.Search != "" {
        tx = tx.Where("`name` like ?", "%"+s.Search+"%")
    }
    total = pagination.Pager(&s.PaginationQ, tx, &list)
    return &list, total, nil
}

// ── 创建 ──
type CreateXxxService struct {
    *forms.CreateXxxForm
    C *gin.Context
}

func (s *CreateXxxService) Run() error {
    record := models.XxxModel{Name: s.Name}
    result := global.App.DB.Model(&models.XxxModel{}).Create(&record)
    if result.Error != nil {
        var mysqlErr *mysql.MySQLError
        if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
            return fmt.Errorf("记录`%s`已存在", s.Name)
        }
        return result.Error
    }
    return nil
}

// ── 更新 ──
type UpdateXxxService struct {
    *forms.UpdateXxxForm
    C  *gin.Context
    ID uint64
}

func (s *UpdateXxxService) Run() error {
    result := global.App.DB.Model(&models.XxxModel{}).Where("id=?", s.ID).
        Updates(map[string]interface{}{"name": s.Name})
    if result.Error != nil { return result.Error }
    return nil
}

// ── 删除 ──
type DeleteXxxService struct {
    C  *gin.Context
    ID uint64
}

func (s *DeleteXxxService) Run() error {
    tx := global.App.DB.Where("id=?", s.ID).Delete(&models.XxxModel{})
    if tx.Error != nil { return tx.Error }
    return nil
}
```

## Run() 方法签名

| 场景 | 签名 |
|------|------|
| 分页列表 | `(responseData interface{}, total int64, err error)` |
| 单条数据 | `(responseData interface{}, err error)` |
| 写操作（创建/更新/删除） | `(err error)` |

## 常用模式

### 数据库事务

```go
return global.App.DB.Transaction(func(tx *gorm.DB) error {
    if err := tx.Model(&models.Xxx{}).Create(&record).Error; err != nil {
        return err
    }
    return nil
})
```

### MySQL 1062 重复键处理

```go
var mysqlErr *mysql.MySQLError
if errors.As(result.Error, &mysqlErr) && mysqlErr.Number == 1062 {
    return fmt.Errorf("记录`%s`已存在", s.Name)
}
```

### 组织范围查询

```go
// 放在 services/org_scope.go 中
func applyOrgDescendantScope(tx *gorm.DB, column, key string) *gorm.DB {
    return tx.Where(fmt.Sprintf("(%s = ? OR %s LIKE ?)", column, column), key, key+"-%")
}
```
