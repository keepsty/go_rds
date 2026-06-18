# 统一响应规范（pkg/response/）

所有 API 响应必须通过 `response` 包返回，禁止直接 `c.JSON()`。

## 响应方法

| 方法 | Code | 说明 |
|------|------|------|
| `response.Success(c, data, msg)` | `0000` | 操作成功 |
| `response.Fail(c, msg)` | `0001` | 业务失败 |
| `response.ValidateFail(c, msg)` | `0001` | 参数验证失败 |
| `response.PaginationSuccess(c, total, data)` | `0000` | 分页成功 |

## JSON 格式

```json
// 普通响应
{"request_id": "xxx", "code": "0000", "data": {...}, "message": "success"}

// 分页响应
{"request_id": "xxx", "code": "0000", "data": [...], "message": "success", "total": 100}
```

内部自动注入 `request_id` 并记录日志。

## 分页工具（pkg/pagination/）

```go
type Pagination struct {
    PageSize int  `form:"page_size" json:"page_size"`
    Page     int  `form:"page" json:"page"`
    IsPage   bool `form:"is_page" json:"is_page"`
}

// 用法：返回 total 同时自动填充数据到 list
total = pagination.Pager(&s.PaginationQ, tx, &list)
```
