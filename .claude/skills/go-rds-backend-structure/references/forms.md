# 表单层（forms/）

定义请求参数结构体，使用 `form`/`json` tag 和 `binding` 验证。

## 模板

```go
package forms

import "github.com/keepsty/go_rds/pkg/pagination"

// 列表查询：嵌入 PaginationQ 实现分页
type GetXxxForm struct {
    PaginationQ pagination.Pagination
    Search      string `form:"search"`
    Status      string `form:"status"`
}

// 创建：使用 binding validate tag
type CreateXxxForm struct {
    Name  string `form:"name"  json:"name" binding:"required,min=2,max=32"`
    Email string `form:"email" json:"email" binding:"required,email"`
    Type  string `form:"type"  json:"type"  binding:"required,oneof=typeA typeB"`
}

// 更新：同创建，可选字段用 omitempty
type UpdateXxxForm struct {
    Name   string `form:"name"  json:"name" binding:"required,min=2,max=32"`
    Status string `form:"status" json:"status" binding:"required,oneof=active inactive"`
}
```

## 常用 binding 验证 tag

| Tag | 含义 |
|-----|------|
| `required` | 必填 |
| `min=N` | 最小长度 |
| `max=N` | 最大长度 |
| `oneof=a b c` | 枚举值 |
| `email` | 邮箱格式 |
| `omitempty` | 可空（用于更新表单中允许不传的字段） |
| `len=6` | 固定长度 |
| `numeric` | 纯数字 |
| `uuid` | UUID 格式 |
| `boolean` | 布尔值 |
