<template>
  <div class="instance-list">
    <div class="toolbar-actions">
      <a-button type="primary" size="small" @click="refreshList">刷新</a-button>
    </div>
    <a-table
      size="middle"
      :columns="columns"
      :row-key="(record) => record.id"
      :data-source="instanceList"
      :loading="loading"
      :pagination="false"
      :scroll="{ x: 900 }"
      :row-class-name="rowClassName"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'hostname'">
          <a-tooltip title="点击复制">
            <span class="click-copy" @click="copyText(record.hostname)">{{ record.hostname }}</span>
          </a-tooltip>
        </template>
        <template v-if="column.key === 'ip'">
          <a-tooltip title="点击复制">
            <span class="click-copy" @click="copyText(record.ip)">{{ record.ip }}</span>
          </a-tooltip>
        </template>
        <template v-if="column.key === 'role'">
          <a-tag :color="record.role === 1 ? 'blue' : 'green'">
            {{ record.role === 1 ? '主库' : '从库' }}
          </a-tag>
        </template>
        <template v-if="column.key === 'status'">
          <a-badge :status="record.status === 1 ? 'success' : 'default'" />
          {{ record.status === 1 ? 'Online' : 'Offline' }}
        </template>
        <template v-if="column.key === 'action'">
          <a-button-group size="small">
            <a-button v-if="record.role === 0 && record.status !== 1" type="primary" ghost @click="handleOnline(record)">
              上线
            </a-button>
            <a-button v-if="record.role === 0 && record.status === 1" type="warning" ghost @click="handleOffline(record)">
              下线
            </a-button>
          </a-button-group>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
defineOptions({ name: 'InstanceList' })

import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { GetClusterInstancesApi, SetInstanceStatusApi } from '@/api/cluster'

const props = defineProps({
  sgId: { type: Number, required: true },
})

const instanceList = ref([])
const loading = ref(false)

const columns = [
  { title: '主机名', key: 'hostname', width: 160 },
  { title: '配置', key: 'host_info', width: 160 },
  { title: 'IP地址', key: 'ip', width: 140 },
  { title: '端口', key: 'port', width: 80 },
  { title: '版本', key: 'mysql_version', width: 100 },
  { title: '角色', key: 'role', width: 80 },
  { title: '状态', key: 'status', width: 100 },
  { title: '实例用途', key: 'purpose', width: 120 },
  { title: '创建时间', key: 'create_time', width: 170 },
  { title: '操作', key: 'action', width: 160, fixed: 'right' },
]

function rowClassName(record) {
  if (record.role === 1) return 'master-row'
  return 'slave-row'
}

async function getInstanceList() {
  loading.value = true
  try {
    const res = await GetClusterInstancesApi(props.sgId)
    instanceList.value = res.data || []
  } catch {
    instanceList.value = []
  } finally {
    loading.value = false
  }
}

async function handleOnline(record) {
  try {
    await SetInstanceStatusApi({ id: record.id, status: 1 })
    message.success('上线成功')
    getInstanceList()
  } catch {
    message.error('上线失败')
  }
}

async function handleOffline(record) {
  try {
    await SetInstanceStatusApi({ id: record.id, status: 0 })
    message.success('下线成功')
    getInstanceList()
  } catch {
    message.error('下线失败')
  }
}

function copyText(text) {
  navigator.clipboard.writeText(text).then(() => {
    message.success('复制成功')
  })
}

function refreshList() {
  getInstanceList()
}

onMounted(() => {
  getInstanceList()
})
</script>

<style scoped>
.instance-list {
  margin-top: 8px;
}

.toolbar-actions {
  margin-bottom: 12px;
}

.click-copy {
  cursor: pointer;
  color: var(--ant-color-primary, #1677ff);
}

:deep(.master-row) {
  background-color: #e6f4ff;
}

:deep(.slave-row) {
  background-color: #f6ffed;
}
</style>
