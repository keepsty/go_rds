# 视图层（views/）

视图层职责：**参数绑定 → 构造 Service → 调用 Run() → 统一响应**

## 风格一：早期返回（推荐，代码简洁）

```go
func GetResourceView(c *gin.Context) {
    var form forms.GetResourceForm
    if err := c.ShouldBind(&form); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }
    service := services.GetResourceService{
        GetResourceForm: &form,
        C:               c,
    }
    returnData, total, err := service.Run()
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.PaginationSuccess(c, total, returnData)
}

func CreateResourceView(c *gin.Context) {
    var form forms.CreateResourceForm
    if err := c.ShouldBind(&form); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }
    service := services.CreateResourceService{
        CreateResourceForm: &form,
        C:                  c,
    }
    if err := service.Run(); err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, nil, "success")
}

func UpdateResourceView(c *gin.Context) {
    id, ok := parseUint64Param(c, "id")
    if !ok { return }
    var form forms.UpdateResourceForm
    if err := c.ShouldBind(&form); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }
    service := services.UpdateResourceService{
        UpdateResourceForm: &form,
        C:  c,
        ID: id,
    }
    if err := service.Run(); err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, nil, "success")
}

func DeleteResourceView(c *gin.Context) {
    id, ok := parseUint64Param(c, "id")
    if !ok { return }
    service := services.DeleteResourceService{C: c, ID: id}
    if err := service.Run(); err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, nil, "success")
}
```

## 风格二：嵌套 if-else（项目中也存在）

```go
func GetRolesView(c *gin.Context) {
    var form *forms.GetRolesForm = &forms.GetRolesForm{}
    if err := c.ShouldBind(&form); err == nil {
        service := services.GetRolesServices{GetRolesForm: form, C: c}
        returnData, total, err := service.Run()
        if err != nil {
            response.Fail(c, err.Error())
        } else {
            response.PaginationSuccess(c, total, returnData)
        }
    } else {
        response.ValidateFail(c, err.Error())
    }
}
```

> 同一模块内请统一风格。

## URL 参数 + JWT 提取（helpers.go）

每个模块的 `views/helpers.go`：

```go
package views

import (
    "strconv"
    "github.com/keepsty/go_rds/middleware"
    "github.com/keepsty/go_rds/pkg/response"
    "github.com/gin-gonic/gin"
)

func getUsername(c *gin.Context) (string, bool) {
    username, ok := middleware.GetUserNameFromJWT(c)
    if !ok { response.Fail(c, "认证信息无效"); return "", false }
    return username, true
}

func parseUint64Param(c *gin.Context, name string) (uint64, bool) {
    raw := c.Param(name)
    id, err := strconv.ParseUint(raw, 10, 64)
    if err != nil { response.ValidateFail(c, "非法参数: "+name); return 0, false }
    return id, true
}
```

## CRUD View 速查表

| 操作 | View 逻辑 | Service.Return | Response |
|------|----------|----------------|----------|
| 列表 | 绑定 form → Service | `(data, total, err)` | `PaginationSuccess(c, total, data)` |
| 详情 | 解析 id → Service | `(data, err)` | `Success(c, data, "success")` |
| 创建 | 绑定 form → Service | `err` | `Success(c, nil, "success")` |
| 更新 | 解析 id + 绑定 form → Service | `err` | `Success(c, nil, "success")` |
| 删除 | 解析 id → Service | `err` | `Success(c, nil, "success")` |
