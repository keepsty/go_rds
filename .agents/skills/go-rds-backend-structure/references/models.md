# 模型层（models/）

定义 GORM 模型，继承基础 `commonModels.Model`。

## 结构体模板

```go
package models

import (
    "github.com/keepsty/go_rds/internal/common/models"
    "gorm.io/datatypes"
)

type XxxModel struct {
    *models.Model
    Name   string         `gorm:"type:varchar(32);not null;uniqueIndex:uniq_name;comment:名称" json:"name"`
    Status models.EnumType `gorm:"type:ENUM('active','inactive');default:'active';comment:状态" json:"status"`
}

func (XxxModel) TableName() string {
    return "insight_xxx"     // 统一使用 insight_ 前缀
}
```

## 基础模型（common/models/base.go）

```go
type Model struct {
    ID        uint64    `gorm:"primaryKey" json:"id"`
    CreatedAt LocalTime `gorm:"index:idx_created_at;autoCreateTime;comment:创建时间" json:"created_at"`
    UpdatedAt LocalTime `gorm:"index:idx_updated_at;autoUpdateTime;comment:更新时间" json:"updated_at"`
}

type EnumType string   // 枚举类型，实现 Scan/Value 接口
```

## LocalTime（common/models/time.go）

自定义时间类型，JSON 序列化为 `"2006-01-02 15:04:05"` 格式，零值输出 `null`。

## GORM 注解速查

```go
`gorm:"primaryKey"`                            // 主键
`gorm:"autoIncrement"`                          // 自增
`gorm:"type:varchar(32)"`                       // 列类型
`gorm:"not null"`                               // 非空
`gorm:"default:''"`                             // 默认值
`gorm:"uniqueIndex:uniq_name"`                  // 唯一索引
`gorm:"index:idx_name"`                         // 普通索引
`gorm:"comment:描述"`                            // 注释
`gorm:"autoCreateTime"`                         // 自动创建时间
`gorm:"autoUpdateTime"`                         // 自动更新时间
`gorm:"type:ENUM('a','b');default:'a'"`         // 枚举类型
```

## 模型注册自动迁移

在 `bootstrap/db.go` 的 `initializeTables()` 注册：

```go
func initializeTables(db *gorm.DB) {
    err := db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
        &usersModels.InsightUsers{},
        &commonModels.InsightInstances{},
        // ...
        &{模块}Models.XxxModel{},   // ← 新增
    )
}
```
