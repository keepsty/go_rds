# Salt 部署管理实现计划（v2）

> **面向 AI 代理的工作者：** 使用 subagent-driven-development 逐任务实现。

**目标：** 实现 Salt 部署管理功能。模板从数据库读取，部署同步等待结果。

**架构：** 后端新增模型 + 种子数据 + API → 前端新建页面模块

**技术栈：** Go (GORM) + Vue 3 + Ant Design Vue 4.x

---

## 需要创建/修改的文件

### 后端

| 文件 | 操作 | 说明 |
|------|------|------|
| `internal/salt/models/salt.go` | 修改 | 新增 `InsightSaltTemplates` GORM 模型 + `SaltDeployRecord` |
| `internal/bootstrap/db.go` | 修改 | 注册新模型 + 初始化种子模板数据 |
| `internal/salt/services/salt.go` | 修改 | 新增 `GetTemplates` 从 DB 读取 + `DeployFromTemplate` 同步执行 |
| `internal/salt/views/salt.go` | 修改 | 新增 `GetTemplatesView` + `DeployTemplateView` + `ListMinionsView` |
| `internal/salt/forms/salt.go` | 修改 | 新增 `DeployRequest` form |
| `internal/salt/routers/routers.go` | 修改 | 注册新路由 |

### 前端

| 文件 | 操作 | 说明 |
|------|------|------|
| `src/api/salt.js` | 新建 | API 调用模块 |
| `src/views/salt/route.js` | 新建 | 路由定义 |
| `src/views/salt/index.vue` | 新建 | 容器页 |
| `src/views/salt/deploy/index.vue` | 新建 | 部署管理页面 |
| `src/router/ark.js` | 修改 | 注册 salt 路由 |

---

## 1. 模型定义

### `internal/salt/models/salt.go` — 新增 GORM 模型

```go
package models

import (
    "github.com/keepsty/go_rds/internal/common/models"
    "gorm.io/datatypes"
)

// InsightSaltTemplates Salt 部署任务模板
type InsightSaltTemplates struct {
    *models.Model
    Name        string         `gorm:"type:varchar(64);not null;uniqueIndex:uniq_name;comment:模板标识" json:"name"`
    Title       string         `gorm:"type:varchar(128);not null;comment:模板标题" json:"title"`
    Description string         `gorm:"type:varchar(512);not null;default:'';comment:模板描述" json:"description"`
    FieldsSchema datatypes.JSON `gorm:"type:json;comment:字段定义JSON" json:"fields_schema"`
    Defaults    datatypes.JSON `gorm:"type:json;comment:默认值JSON" json:"defaults"`
}

func (InsightSaltTemplates) TableName() string {
    return "insight_salt_templates"
}
```

`FieldsSchema` JSON 结构示例：
```json
[
  {"key": "port", "label": "端口", "type": "number", "required": true, "default": 3306, "description": "MySQL 监听端口"},
  {"key": "version", "label": "版本", "type": "string", "required": true, "default": "8.0"},
  {"key": "command", "label": "命令", "type": "text", "required": true, "description": "要执行的 shell 命令"}
]
```

## 2. 种子数据

### `internal/bootstrap/db.go` — 新增 `initializeSaltTemplates`

```go
func initializeSaltTemplates(db *gorm.DB) {
    templates := []saltModels.InsightSaltTemplates{
        {
            Name: "cmd-run", Title: "远程执行命令",
            Description: "在目标主机上同步执行 shell 命令并返回结果",
            FieldsSchema: datatypes.JSON(`[{"key":"command","label":"命令","type":"text","required":true,"description":"要执行的 shell 命令"}]`),
            Defaults: datatypes.JSON(`{}`),
        },
        {
            Name: "state-apply", Title: "Salt State 部署",
            Description: "在目标主机上应用 Salt state.sls 文件",
            FieldsSchema: datatypes.JSON(`[{"key":"state_file","label":"State文件","type":"string","required":true,"description":"不含 .sls 后缀"},{"key":"saltenv","label":"环境","type":"string","default":"base"}]`),
            Defaults: datatypes.JSON(`{"saltenv":"base"}`),
        },
        {
            Name: "mysql-deploy", Title: "MySQL 实例部署",
            Description: "在目标主机上部署 MySQL 实例（需先上传配置到 Salt Master）",
            FieldsSchema: datatypes.JSON(`[{"key":"port","label":"端口","type":"number","required":true,"default":3306},{"key":"version","label":"版本","type":"string","required":true,"default":"8.0"},{"key":"datadir","label":"数据目录","type":"string","default":"/data/mysql"}]`),
            Defaults: datatypes.JSON(`{"port":3306,"version":"8.0","datadir":"/data/mysql"}`),
        },
    }
    for _, t := range templates {
        var existing saltModels.InsightSaltTemplates
        err := db.Where("name = ?", t.Name).First(&existing).Error
        if errors.Is(err, gorm.ErrRecordNotFound) {
            db.Create(&t)
        }
    }
}
```

在 `initializeMySQLGorm()` 中调用：
```go
initializeTables(db)
initializeAllowedOperations(db)
initializeGlobalInspectParams(db)
initializeAdminUser(db)
initializeNotifySettings(db)
initializeSaltTemplates(db)  // ← 新增
```

## 3. 服务层

### `internal/salt/services/salt.go` — 新增方法

```go
type DeployRequest struct {
    Template string              `json:"template" binding:"required"`
    Targets  []string            `json:"targets" binding:"required,min=1"`
    Config   map[string]interface{} `json:"config"`
}

// GetTemplates 从数据库读取模板列表
func GetTemplates() ([]models.InsightSaltTemplates, error) {
    var templates []models.InsightSaltTemplates
    result := global.App.DB.Model(&models.InsightSaltTemplates{}).Find(&templates)
    return templates, result.Error
}

// DeployFromTemplate 根据模板名和配置执行部署，同步等待结果
func DeployFromTemplate(name string, config map[string]interface{}, targets []string) (interface{}, error) {
    // 1. 查找模板
    var tmpl models.InsightSaltTemplates
    if err := global.App.DB.Where("name = ?", name).First(&tmpl).Error; err != nil {
        return nil, fmt.Errorf("模板 '%s' 不存在", name)
    }

    // 2. 获取 Salt 服务
    svc := NewSaltService()

    // 3. 根据模板类型执行
    switch name {
    case "cmd-run":
        cmd, _ := config["command"].(string)
        if cmd == "" {
            return nil, fmt.Errorf("命令不能为空")
        }
        result, err := svc.RunCommand(targets[0], cmd, false)
        if err != nil {
            return nil, err
        }
        return map[string]interface{}{
            "success": result.Success,
            "message": result.Error,
            "detail":  result.Detail,
        }, nil

    case "state-apply":
        stateFile, _ := config["state_file"].(string)
        if stateFile == "" {
            return nil, fmt.Errorf("state 文件不能为空")
        }
        result, err := svc.RunState(targets[0], stateFile, false)
        if err != nil {
            return nil, err
        }
        return map[string]interface{}{
            "success": result.Success,
            "message": result.Error,
            "detail":  result.Detail,
        }, nil

    case "mysql-deploy":
        // MySQL 部署使用现有的 InstallMySQLHandler
        hp := []*models.SaltMysqlHostPost{
            {
                Port:    int64(config["port"].(float64)),
                Host:    targets[0],
                Version: config["version"].(string),
            },
        }
        si := &models.SaltMysqlServerInfo{
            HostPort: hp,
        }
        data, err := InstallMySQLHandler("prod", si, &global.App.Config.Salt)
        if err != nil {
            return nil, err
        }
        return map[string]interface{}{
            "success": true,
            "detail":  data,
        }, nil

    default:
        return nil, fmt.Errorf("不支持模板类型: %s", name)
    }
}
```

## 4. 视图层

### `internal/salt/views/salt.go` — 新增 handler

```go
func GetTemplatesView(c *gin.Context) {
    templates, err := services.GetTemplates()
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, templates, "success")
}

func DeployTemplateView(c *gin.Context) {
    var req forms.DeployRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidateFail(c, err.Error())
        return
    }
    result, err := services.DeployFromTemplate(req.Template, req.Config, req.Targets)
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, result, "部署完成")
}

func ListMinionsView(c *gin.Context) {
    svc := services.NewSaltService()
    list, err := svc.ListMinions()
    if err != nil {
        response.Fail(c, err.Error())
        return
    }
    response.Success(c, list, "success")
}
```

### `internal/salt/forms/salt.go` — 新增

```go
type DeployRequest struct {
    Template string                 `json:"template" binding:"required"`
    Targets  []string               `json:"targets" binding:"required,min=1"`
    Config   map[string]interface{} `json:"config"`
}
```

## 5. 路由

### `internal/salt/routers/routers.go`

```go
func RegisterApiRoutes(v1 *gin.RouterGroup) {
    v1.GET("/templates", views.GetTemplatesView)
    v1.POST("/templates/deploy", views.DeployTemplateView)
    v1.GET("/minions", views.ListMinionsView)
    v1.POST("/deploy/mysql-cluster", views.AddMySQLClusterHandler)
}
```

## 6. 前端

### `src/api/salt.js`

```javascript
import { del, get, post, put } from '@/utils/request'

export const getTemplatesApi = () => get('/api/v1/salt/templates')
export const deployTaskApi = (data) => post('/api/v1/salt/templates/deploy', data)
export const getMinionsApi = () => get('/api/v1/salt/minions')
```

### `src/views/salt/route.js`

```javascript
const route = {
  name: 'view.salt',
  path: '/salt',
  component: () => import('./index.vue'),
  redirect: '/salt/deploy',
  meta: { title: 'Salt部署', icon: 'ThunderboltOutlined', keepAlive: true },
  children: [
    {
      name: 'view.salt.deploy',
      path: '/salt/deploy',
      component: () => import('./deploy/index.vue'),
      meta: { title: '部署管理', keepAlive: true },
    },
  ],
}
export default route
```

### `src/views/salt/index.vue`

```vue
<template><router-view /></template>
```

### `src/views/salt/deploy/index.vue`

```vue
<template>
  <div>
    <a-card title="Salt 部署管理">
      <a-form layout="vertical">
        <a-row :gutter="24">
          <a-col :span="8">
            <a-form-item label="部署模板" required>
              <a-select
                v-model:value="formData.template"
                placeholder="请选择部署模板"
                :options="templateOptions"
                @change="onTemplateChange"
              />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="目标主机" required>
              <a-select
                v-model:value="formData.targets"
                mode="multiple"
                placeholder="选择目标 minion"
                :options="minionOptions"
                :max-tag-count="3"
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <a-card title="配置参数" style="margin-top: 16px">
      <a-form layout="vertical" v-if="currentTemplate">
        <a-row :gutter="24">
          <a-col
            v-for="field in currentTemplate.fields"
            :key="field.key"
            :span="field.type === 'text' ? 24 : 12"
          >
            <a-form-item
              :label="field.label"
              :required="field.required"
              :help="field.description"
            >
              <a-input
                v-if="field.type === 'string'"
                v-model:value="field.value"
              />
              <a-input-number
                v-else-if="field.type === 'number'"
                v-model:value="field.value"
                style="width: 100%"
              />
              <a-switch
                v-else-if="field.type === 'boolean'"
                v-model:checked="field.value"
              />
              <a-textarea
                v-else-if="field.type === 'text'"
                v-model:value="field.value"
                :rows="4"
              />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <a-empty v-else description="请先选择部署模板" />
    </a-card>

    <div style="margin-top: 16px; text-align: right">
      <a-button
        type="primary"
        :loading="deploying"
        :disabled="!canDeploy"
        @click="handleDeploy"
      >
        <template #icon><ThunderboltOutlined /></template>
        {{ deploying ? '执行中...' : '执行部署' }}
      </a-button>
    </div>

    <a-card title="执行结果" style="margin-top: 16px" v-if="deployResult !== null">
      <a-result
        :status="deployResult.success ? 'success' : 'error'"
        :title="deployResult.success ? '部署成功' : '部署失败'"
        :sub-title="deployResult.message || ''"
      >
        <template #extra>
          <a-button @click="deployResult = null">清除结果</a-button>
        </template>
      </a-result>
      <a-collapse v-if="deployResult.detail">
        <a-collapse-panel header="查看详细输出">
          <pre style="max-height:400px;overflow:auto;background:#f5f5f5;padding:12px;border-radius:4px">{{ JSON.stringify(deployResult.detail, null, 2) }}</pre>
        </a-collapse-panel>
      </a-collapse>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ThunderboltOutlined } from '@ant-design/icons-vue'
import { getTemplatesApi, deployTaskApi, getMinionsApi } from '@/api/salt'

const formData = reactive({ template: undefined, targets: [] })
const deploying = ref(false)
const deployResult = ref(null)
const templates = ref([])
const minions = ref([])
const currentTemplate = ref(null)

const templateOptions = computed(() =>
  templates.value.map((t) => ({ value: t.name, label: t.title }))
)
const minionOptions = computed(() =>
  minions.value.map((m) => ({ value: m, label: m }))
)
const canDeploy = computed(() =>
  formData.template && formData.targets.length > 0
)

const onTemplateChange = () => {
  const tmpl = templates.value.find((t) => t.name === formData.template)
  if (!tmpl) { currentTemplate.value = null; return }
  const fields = JSON.parse(tmpl.fields_schema || '[]')
  const defaults = JSON.parse(tmpl.defaults || '{}')
  currentTemplate.value = {
    ...tmpl,
    fields: fields.map((f) => ({
      ...f,
      value: defaults[f.key] ?? f.default ?? '',
    })),
  }
  deployResult.value = null
}

const buildConfig = () => {
  const config = {}
  currentTemplate.value.fields.forEach((f) => { config[f.key] = f.value })
  return config
}

const handleDeploy = async () => {
  deploying.value = true
  deployResult.value = null
  try {
    const res = await deployTaskApi({
      template: formData.template,
      targets: formData.targets,
      config: buildConfig(),
    })
    deployResult.value = res.data || res
    if (res.code === '0000') message.success('部署完成')
  } catch (err) {
    deployResult.value = { success: false, message: err.message || '请求异常', detail: null }
    message.error('部署失败: ' + (err.message || '未知错误'))
  } finally {
    deploying.value = false
  }
}

onMounted(async () => {
  try {
    const [tplRes, minRes] = await Promise.all([
      getTemplatesApi(),
      getMinionsApi(),
    ])
    templates.value = tplRes.data || []
    minions.value = minRes.data?.minions || []
  } catch {
    message.warning('加载模板或 minion 列表失败')
  }
})
</script>
```

### `src/router/ark.js` — 注册路由

在 `asyncRoutes` 的 `children` 数组中导入并添加 saltRoute：

```javascript
import saltRoute from '@/views/salt/route'

// 在 asyncRoutes children 中加入
children: [
  // ... 现有路由,
  saltRoute,
],
```

---

## 验证方法

1. `cd backend && go build ./...`
2. `cd www && npm run build`
3. 启动后访问 `/salt/deploy` → 选择模板 → 填写配置 → 选目标 → 执行 → 看结果