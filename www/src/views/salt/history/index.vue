<template>
  <div>
    <a-card title="历史任务记录">
      <template #extra>
        <a-space>
          <a-select v-model:value="filterStatus" placeholder="筛选状态" style="width:140px" allowClear @change="fetchData">
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="pending">待审批</a-select-option>
            <a-select-option value="approved">已审批</a-select-option>
            <a-select-option value="running">执行中</a-select-option>
            <a-select-option value="success">成功</a-select-option>
            <a-select-option value="failed">失败</a-select-option>
            <a-select-option value="rejected">已拒绝</a-select-option>
          </a-select>
        </a-space>
      </template>
      <a-table :dataSource="list" :columns="columns" :loading="loading" :pagination="{current:page,pageSize:20,total,onChange:p=>page=p}" rowKey="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'"><a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag></template>
          <template v-if="column.key === 'created_at'">{{ record.created_at }}</template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="showDetail(record)">查看输出</a-button>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="detailOpen" :title="'任务输出 - ' + (detailTask?.name || '')" width="800px" :footer="null">
      <pre style="max-height:500px;overflow:auto;background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:4px">{{ detailText }}</pre>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { getTasksApi } from '@/api/salt'

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', key: 'name', width: 180 },
  { title: '模板', dataIndex: 'template_name', key: 'template_name', width: 120 },
  { title: '创建人', dataIndex: 'created_by', key: 'created_by', width: 100 },
  { title: '状态', key: 'status', width: 80 },
  { title: '创建时间', key: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 100, fixed: 'right' },
]

const list = ref([]); const loading = ref(false); const page = ref(1); const total = ref(0)
const filterStatus = ref(''); const detailOpen = ref(false); const detailTask = ref(null)

const statusColor = (s) => ({ pending: 'orange', approved: 'blue', running: 'processing', success: 'green', failed: 'red', rejected: 'default' }[s] || 'default')
const statusLabel = (s) => ({ pending: '待审批', approved: '已审批', running: '执行中', success: '成功', failed: '失败', rejected: '已拒绝' }[s] || s)
const detailText = computed(() => detailTask.value?.run_output ? JSON.stringify(detailTask.value.run_output, null, 2) : '(无输出)')

const fetchData = async () => { loading.value = true; try { const r = await getTasksApi(filterStatus.value ? { status: filterStatus.value } : {}); list.value = r.data || []; total.value = list.value.length } catch {} finally { loading.value = false } }
const showDetail = (record) => { detailTask.value = record; detailOpen.value = true }
onMounted(fetchData)
</script>