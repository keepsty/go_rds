---
name: add-frontend-page
description: Add a frontend page (route + view + API module + store) matching backend module, following Ant Design Vue patterns
runAs: subagent
allowed-tools: read_file, write_file, edit_file, multi_edit, search_content, glob, list_directory
---
You are a Vue 3 frontend page generator for GoRDS (`www/`).

## Project Conventions

- **Framework**: Vue 3 + Composition API (`<script setup>`)
- **UI library**: Ant Design Vue 4 (`a-table`, `a-card`, `a-descriptions`, `a-tag`, etc.)
- **HTTP**: `@/utils/request` exports `get`, `post`, `put`, `del`
- **State**: Pinia store at `src/store/`
- **Router**: async routes in `src/router/ark.js` (Layout-based)
- **Permission**: route guard in `src/permission.js`
- **Style**: scoped CSS with CSS variables (`var(--ant-colorPrimary, #1677ff)`)
- **API modules**: one file per module in `src/api/{name}.js`

## Input

- `name`: route/page name — kebab-case (e.g. `cluster`)
- `Name`: PascalCase version (e.g. `Cluster`)
- `title`: display title (e.g. "集群管理")
- `icon`: Ant Design icon name (e.g. `ClusterOutlined`)
- `api_path`: backend API prefix (e.g. `/api/v1/cluster`)
- `feature`: description

## Steps

### Step 1: Create API module `src/api/{name}.js`

```js
import { get, post, put, del } from "@/utils/request"

// 列表
export const Get{Name}ListApi = (params) => get('{api_path}/{feature}s', params)
// 详情
export const Get{Name}DetailApi = (id) => get(`{api_path}/{feature}s/${id}`)
// 创建
export const Create{Name}Api = (data) => post('{api_path}/{feature}s', data)
// 更新
export const Update{Name}Api = (data) => put(`{api_path}/{feature}s/${data.id}`, data)
// 删除
export const Delete{Name}Api = (id) => del(`{api_path}/{feature}s/${id}`)
```

### Step 2: Create route file `src/views/{name}/route.js`

```js
const route = {
  name: '{name}',
  path: '/{name}',
  redirect: '/{name}/list',
  component: () => import('./index.vue'),
  meta: { title: '{title}', icon: '{icon}', keepAlive: true },
  children: [
    {
      name: '{name}.list',
      path: '/{name}/list',
      component: () => import('./list/{Name}List.vue'),
      meta: { title: '{title}列表', icon: 'UnorderedListOutlined', keepAlive: true },
    },
    {
      name: '{name}.detail',
      path: '/{name}/:id',
      component: () => import('./detail/{Name}Detail.vue'),
      meta: { title: '{title}详情', keepAlive: true, hidden: true },
    },
  ],
}

export default route
```

### Step 3: Create `src/views/{name}/index.vue`

```vue
<template>
  <router-view />
</template>
<script setup>
defineOptions({ name: '{Name}IndexView' })
</script>
```

### Step 4: Create list view `src/views/{name}/list/{Name}List.vue`

Template structure:
```vue
<template>
  <div class="gi-page-shell">
    <a-card class="xxx-list-card" title="{title}列表">
      <div class="gi-page-toolbar toolbar">
        <div class="toolbar-left">
          <a-button type="primary" @click="handleCreate">新建</a-button>
        </div>
        <a-input-search v-model:value="query.keyword" placeholder="搜索" allowClear
          class="toolbar-search" @search="loadData" />
      </div>
      <a-table
        size="middle"
        :columns="columns"
        :row-key="(r) => r.id"
        :data-source="dataList"
        :pagination="pagination"
        :loading="loading"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <router-link :to="{ name: '{name}.detail', params: { id: record.id } }">
              {{ record.name }}
            </router-link>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" @click="handleEdit(record)">编辑</a-button>
            <a-button type="link" danger @click="handleDelete(record)">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>
```

Script pattern:
```vue
<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { Get{Name}ListApi, Delete{Name}Api } from '@/api/{name}'

const router = useRouter()
const dataList = ref([])
const loading = ref(false)
const total = ref(0)

const query = reactive({ page: 1, size: 10, keyword: '' })

const columns = [
  { title: '名称', key: 'name', width: 200 },
  { title: '状态', key: 'status', width: 100 },
  { title: '创建时间', key: 'create_time', width: 180 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' },
]

const pagination = computed(() => ({
  current: query.page, pageSize: query.size, total: total.value,
  showSizeChanger: true, pageSizeOptions: ['10', '20', '50'],
  showTotal: (t) => `共 ${t} 条`,
}))

async function loadData() {
  loading.value = true
  try {
    const res = await Get{Name}ListApi({ page: query.page, size: query.size })
    dataList.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch { dataList.value = []
  } finally { loading.value = false }
}

function handleTableChange(pg) {
  query.page = pg.current; query.size = pg.pageSize; loadData()
}

onMounted(() => loadData())
</script>
```

### Step 5: Create detail view `src/views/{name}/detail/{Name}Detail.vue`

Use `a-descriptions` for showing data, `a-tabs` for sub-sections. Pattern:

```vue
<template>
  <a-card title="{title}详情" :loading="loading">
    <template #extra>
      <a-button type="link" @click="router.push({name:'{name}.list'})">← 返回</a-button>
    </template>
    <a-descriptions bordered :column="2" size="small">
      <a-descriptions-item label="名称">{{ detail.name }}</a-descriptions-item>
    </a-descriptions>
  </a-card>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Get{Name}DetailApi } from '@/api/{name}'
const route = useRoute(), router = useRouter()
const id = ref(parseInt(route.params.id, 10))
const loading = ref(false), detail = ref({})
async function load() {
  loading.value = true
  try { const res = await Get{Name}DetailApi(id.value); detail.value = res.data }
  finally { loading.value = false }
}
onMounted(() => load())
</script>
```

### Step 6: Register route in `src/router/ark.js`

```js
import {Name.toUpperCase()} from '@/views/{name}/route'
```
and include `{Name.toUpperCase()},` inside `children: [ ... ]` in `asyncRoutes`.

**Note**: Write the literal variable name, e.g. `import CLUSTER from '@/views/cluster/route'`.

### Step 7: Verify

```bash
npm --prefix www run build
```
