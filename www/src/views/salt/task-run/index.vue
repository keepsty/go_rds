<template>
  <div>
    <a-card title="运行任务">
      <a-table :dataSource="pendingList" :columns="columns" :loading="loading" :pagination="false" rowKey="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'"><a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag></template>
          <template v-if="column.key === 'action'">
            <a-button v-if="record.status === 'pending'" type="primary" size="small" @click="handleApprove(record)">审批通过</a-button>
            <a-button v-else-if="record.status === 'approved'" type="primary" size="small" :loading="runningId === record.id" @click="handleRun(record)">执行部署</a-button>
            <a-button v-else-if="record.status === 'running'" type="default" size="small" disabled>执行中...</a-button>
            <span v-else>-</span>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-card title="执行输出" style="margin-top:16px" v-if="outputTask">
      <a-button v-if="outputTask.status==='running'" type="danger" size="small" style="margin-bottom:8px">正在执行...</a-button>
      <pre style="max-height:400px;overflow:auto;background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:4px">{{ outputText }}</pre>
    </a-card>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { getTasksApi, approveTaskApi, runTaskApi, getTemplatesApi, getHostConfigsApi } from '@/api/salt'

const columns = [
  { title: '任务名称', dataIndex: 'name', key: 'name', width: 180 },
  { title: '模板', dataIndex: 'template_name', key: 'template_name', width: 120 },
  { title: '主机配置ID', dataIndex: 'host_config_id', key: 'host_config_id', width: 100 },
  { title: '创建人', dataIndex: 'created_by', key: 'created_by', width: 100 },
  { title: '状态', key: 'status', width: 100 },
  { title: '操作', key: 'action', width: 120, fixed: 'right' },
]

const allTasks = ref([]); const loading = ref(false); const runningId = ref(null); const outputTask = ref(null)
const pendingList = computed(() => allTasks.value.filter(t => ['pending', 'approved', 'running'].includes(t.status)))
const outputText = computed(() => outputTask.value?.run_output ? JSON.stringify(outputTask.value.run_output, null, 2) : '(无输出)')

const statusColor = (s) => ({ pending: 'orange', approved: 'blue', running: 'processing', success: 'green', failed: 'red', rejected: 'default' }[s] || 'default')
const statusLabel = (s) => ({ pending: '待审批', approved: '已审批', running: '执行中', success: '成功', failed: '失败', rejected: '已拒绝' }[s] || s)

const fetchData = async () => { loading.value = true; try { const r = await getTasksApi(); allTasks.value = r.data || [] } catch { message.error('获取失败') } finally { loading.value = false } }
const handleApprove = (record) => Modal.confirm({ title: '审批', content: `审批通过任务「${record.name}」？`, okText: '通过', onOk: async () => { await approveTaskApi(record.id, { action: 'approve' }); message.success('已审批'); await fetchData() } })
const handleRun = async (record) => {
  runningId.value = record.id; outputTask.value = record
  try {
    await runTaskApi(record.id)
    message.success('执行完成'); await fetchData()
    const r = await getTasksApi(); const updated = (r.data || []).find(t => t.id === record.id); if (updated) outputTask.value = updated
  } catch (e) { message.error(e.message || '执行失败'); await fetchData() }
  finally { runningId.value = null }
}
onMounted(fetchData)
</script>