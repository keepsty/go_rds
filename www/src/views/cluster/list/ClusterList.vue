<template>
  <div class="gi-page-shell cluster-list-page">
    <a-card class="cluster-list-card" title="集群列表">
      <div class="gi-page-toolbar toolbar">
        <div class="toolbar-left">
          <a-button type="primary" @click="handleCreate">创建集群</a-button>
        </div>
        <a-input-search
          v-model:value="queryInfo.query"
          placeholder="请输入集群名称搜索"
          allowClear
          class="toolbar-search"
          @search="getClusterList"
        />
      </div>

      <div class="table-wrapper">
        <a-table
          size="middle"
          :columns="columns"
          :row-key="(record) => record.id"
          :data-source="clusterList"
          :pagination="pagination"
          :loading="loading"
          @change="handleTableChange"
          :scroll="{ x: 1000 }"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'name'">
              <router-link class="name-link" :to="{ name: 'cluster.detail', params: { sg_id: record.id } }">
                {{ record.name }}
              </router-link>
            </template>
            <template v-if="column.key === 'environment'">
              <a-tag :color="envColorMap[record.environment] || 'default'">
                {{ envMap[record.environment] || record.environment }}
              </a-tag>
            </template>
            <template v-if="column.key === 'ha_type'">
              {{ haTypeMap[record.ha_type] || record.ha_type }}
            </template>
            <template v-if="column.key === 'middleware'">
              {{ middlewareMap[record.middleware] || record.middleware }}
            </template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" @click="handleDetail(record.id)">
                详情
              </a-button>
              <a-button type="link" size="small" @click="handleEdit(record.id)">
                编辑
              </a-button>
            </template>
          </template>
        </a-table>
      </div>
    </a-card>
  </div>
</template>

<script setup>
defineOptions({ name: 'ClusterListView' })

import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GetClustersApi } from '@/api/cluster'

const router = useRouter()

const clusterList = ref([])
const loading = ref(false)
const total = ref(0)

const queryInfo = reactive({
  page: 1,
  size: 10,
  query: '',
})

const columns = [
  { title: '#', key: 'index', width: 60, customRender: ({ index }) => index + 1 },
  { title: '集群名', key: 'name', width: 180 },
  { title: '环境', key: 'environment', width: 100 },
  { title: '产品线', key: 'prod_name', width: 120 },
  { title: '中间件', key: 'middleware', width: 100 },
  { title: 'HA类型', key: 'ha_type', width: 100 },
  { title: '高峰期', key: 'peak_time', width: 120 },
  { title: '创建时间', key: 'create_time', width: 180 },
  { title: '负责DBA', key: 'dba_user', width: 120 },
  { title: '操作', key: 'action', width: 150, fixed: 'right' },
]

const envMap = { 0: 'prod', 1: 'rc', 2: 'k8s', 3: 'press' }
const envColorMap = { 0: 'red', 1: 'orange', 2: 'blue', 3: 'purple' }
const haTypeMap = { 0: 'MHA', 1: 'ORC', 2: 'MGR' }
const middlewareMap = { 0: 'ProxySQL', 1: 'Zebra', 2: 'MGW' }

const pagination = computed(() => ({
  current: queryInfo.page,
  pageSize: queryInfo.size,
  total: total.value,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100'],
  showTotal: (total) => `共 ${total} 条`,
}))

async function getClusterList() {
  loading.value = true
  try {
    const res = await GetClustersApi({
      page: queryInfo.page,
      size: queryInfo.size,
    })
    clusterList.value = res.data?.clusters || []
    total.value = res.data?.total || 0
  } catch {
    clusterList.value = []
  } finally {
    loading.value = false
  }
}

function handleTableChange(pagination) {
  queryInfo.page = pagination.current
  queryInfo.size = pagination.pageSize
  getClusterList()
}

function handleDetail(id) {
  router.push({ name: 'cluster.detail', params: { sg_id: id } })
}

function handleEdit(id) {
  console.log('edit cluster:', id)
}

function handleCreate() {
  console.log('create cluster')
}

onMounted(() => {
  getClusterList()
})
</script>

<style scoped>
.cluster-list-card {
  border-radius: 12px;
  border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  box-shadow: 0 2px 8px rgb(0 0 0 / 5%);
}

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 8px;
}

.toolbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.toolbar-search {
  width: 300px;
}

.table-wrapper {
  margin-top: 8px;
}

.name-link {
  font-weight: 500;
}
</style>
