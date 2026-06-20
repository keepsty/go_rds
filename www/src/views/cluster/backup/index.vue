<template>
  <a-card title="备份管理">
    <a-tabs v-model:activeKey="activeTab">
      <!-- 备份配置 -->
      <a-tab-pane key="configs" tab="备份配置" forceRender>
        <div style="margin-bottom:12px;display:flex;gap:12px;flex-wrap:wrap">
          <a-select v-model:value="cfgFilter.db_type" placeholder="数据库类型" style="width:140px" allowClear @change="fetchConfigs">
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="MySQL">MySQL</a-select-option>
            <a-select-option value="Redis">Redis</a-select-option>
            <a-select-option value="TiDB">TiDB</a-select-option>
          </a-select>
          <a-input-search v-model:value="cfgFilter.search" placeholder="搜索配置名称" style="width:200px" @search="fetchConfigs" />
          <a-button type="primary" size="small" @click="openConfigModal()"><PlusOutlined />新增配置</a-button>
        </div>
        <a-table :dataSource="filteredConfigs" :columns="configColumns" :loading="cfgLoading" :pagination="false" rowKey="id" size="small">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'db_type'"><a-tag :color="{MySQL:'blue',Redis:'green',TiDB:'purple'}[record.db_type]">{{ record.db_type }}</a-tag></template>
            <template v-if="column.key === 'backup_type'"><a-tag>{{ record.backup_type }}</a-tag></template>
            <template v-if="column.key === 'status'"><a-tag :color="record.status==='enabled'?'green':'default'">{{ record.status==='enabled'?'启用':'禁用' }}</a-tag></template>
            <template v-if="column.key === 'action'">
              <a-button type="link" size="small" @click="openConfigModal(record)">编辑</a-button>
              <a-button type="link" size="small" danger @click="handleDeleteConfig(record)">删除</a-button>
              <a-button type="link" size="small" @click="openTaskModal(record)">创建任务</a-button>
            </template>
          </template>
        </a-table>
      </a-tab-pane>


      <!-- 备份任务（集群维度聚合） -->
      <a-tab-pane key="tasks" tab="备份任务">
        <!-- 集群概览 -->
        <template v-if="!drillInCluster">
          <div style="margin-bottom:12px;display:flex;gap:12px;flex-wrap:wrap">
            <a-select v-model:value="taskFilter.db_type" placeholder="数据库类型" style="width:140px" allowClear @change="filterTasks">
              <a-select-option value="">全部</a-select-option>
              <a-select-option value="MySQL">MySQL</a-select-option>
              <a-select-option value="Redis">Redis</a-select-option>
              <a-select-option value="TiDB">TiDB</a-select-option>
            </a-select>
            <a-input-search v-model:value="taskFilter.search" placeholder="搜索集群名称" style="width:200px" @search="filterTasks" />
          </div>
          <a-table :dataSource="clusterSummary" :columns="clusterColumns" :loading="taskLoading" :pagination="false" rowKey="key" size="middle">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'db_type'"><a-tag :color="{MySQL:'blue',Redis:'green',TiDB:'purple'}[record.db_type]">{{ record.db_type }}</a-tag></template>
              <template v-if="column.key === 'last_time'">{{ record.last_time || '-' }}</template>
              <template v-if="column.key === 'action'"><a-button type="link" @click="drillIn(record)">查看详情</a-button></template>
            </template>
          </a-table>
        </template>

        <!-- 集群详情（钻取） -->
        <template v-else>
          <a-button type="link" @click="drillIn(null)" style="margin-bottom:12px"><LeftOutlined /> 返回集群列表</a-button>
          <a-card size="small" style="margin-bottom:12px">
            <a-descriptions :title="drillInCluster.label" :column="4" size="small">
              <a-descriptions-item label="数据库类型">{{ drillInCluster.db_type }}</a-descriptions-item>
              <a-descriptions-item label="实例">{{ drillInCluster.instance }}</a-descriptions-item>
              <a-descriptions-item label="总任务数">{{ drillInCluster.total }}</a-descriptions-item>
              <a-descriptions-item label="成功/失败">{{ drillInCluster.success }}/{{ drillInCluster.failed }}</a-descriptions-item>
            </a-descriptions>
          </a-card>
          <a-table :dataSource="drillTasks" :columns="drillColumns" :pagination="false" rowKey="id" size="small">
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'"><a-tag :color="statusColor(record.status)">{{ statusLabel(record.status) }}</a-tag></template>
              <template v-if="column.key === 'file_name'">{{ record._records?.[0]?.file_name || '-' }}</template>
              <template v-if="column.key === 'file_size'">{{ record._records?.[0]?.file_size ? (record._records[0].file_size/1024/1024).toFixed(2)+'MB' : '-' }}</template>
              <template v-if="column.key === 'action'"><a-button type="link" size="small" @click="openDetailDrawer(record)">详情</a-button></template>
            </template>
          </a-table>
        </template>
      </a-tab-pane>
    </a-tabs>

    <!-- 配置编辑弹窗 -->
    <a-modal v-model:open="cfgModalOpen" :title="isEditConfig?'编辑配置':'新增配置'" width="600px" :confirm-loading="cfgSubmitting" @ok="handleConfigSubmit" destroyOnClose>
      <a-form layout="vertical">
        <a-form-item label="配置名称" required><a-input v-model:value="cfgForm.name" /></a-form-item>
        <a-form-item label="数据库类型" required><a-select v-model:value="cfgForm.db_type" :options="[{value:'MySQL',label:'MySQL'},{value:'Redis',label:'Redis'},{value:'TiDB',label:'TiDB'}]" /></a-form-item>
        <a-form-item label="实例标识" required><a-input v-model:value="cfgForm.instance_id" placeholder="主机:端口" /></a-form-item>
        <a-form-item label="备份类型" required><a-select v-model:value="cfgForm.backup_type" :options="[{value:'full',label:'全量'},{value:'incremental',label:'增量'}]" /></a-form-item>
        <a-row :gutter="16"><a-col :span="12"><a-form-item label="保留天数"><a-input-number v-model:value="cfgForm.retention_days" :min="1" style="width:100%" /></a-form-item></a-col><a-col :span="12"><a-form-item label="定时策略"><a-input v-model:value="cfgForm.schedule_cron" placeholder="0 2 * * *" /></a-form-item></a-col></a-row>
        <a-form-item label="存储路径"><a-input v-model:value="cfgForm.storage_path" /></a-form-item>
        <a-form-item label="状态" v-if="isEditConfig"><a-switch v-model:checked="cfgForm.statusEnabled" checked-children="启用" un-checked-children="禁用" /></a-form-item>
        <a-form-item label="备注"><a-textarea v-model:value="cfgForm.remark" /></a-form-item>
      </a-form>
    </a-modal>

    <!-- 创建任务弹窗 -->
    <a-modal v-model:open="taskModalOpen" title="创建备份任务" width="640px" :confirm-loading="taskSubmitting" @ok="handleTaskSubmit" destroyOnClose>
      <a-form layout="vertical" v-if="taskCfg">
        <a-form-item label="任务名称" required><a-input v-model:value="taskForm.name" :placeholder="taskCfg.name + '-备份'" /></a-form-item>
        <a-form-item label="备份模板"><a-select v-model:value="taskForm.template_id" placeholder="选择备份模板" :options="filteredTemplates.map(t=>({value:t.id,label:t.name}))" allowClear @change="onTaskTemplateChange" /></a-form-item>
        <a-divider v-if="taskFormTemplate">模板参数</a-divider>
        <a-row :gutter="16" v-if="taskFormTemplate">
          <a-col v-for="f in taskFormTemplate.fields" :key="f.key" :span="f.type === 'text' ? 24 : 12">
            <a-form-item :label="f.label" :required="f.required" :help="f.description">
              <a-input v-if="f.type==='string'" v-model:value="f.value" />
              <a-input-number v-else-if="f.type==='number'" v-model:value="f.value" style="width:100%" />
              <a-switch v-else-if="f.type==='boolean'" v-model:checked="f.value" />
              <a-textarea v-else-if="f.type==='text'" v-model:value="f.value" :rows="3" />
              <a-select v-else-if="f.type==='select'" v-model:value="f.value" :options="(f.options||[]).map(o=>({value:o,label:o}))" style="width:100%" />
            </a-form-item>
          </a-col>
          <a-empty v-if="!taskFormTemplate.fields.length" description="该模板无参数" />
        </a-row>
      </a-form>
    </a-modal>

    <!-- 任务日志弹窗 -->
    <a-modal v-model:open="logModalOpen" :title="logTitle" width="800px" :footer="null" destroyOnClose>
      <pre style="max-height:500px;overflow:auto;background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:4px;white-space:pre-wrap;word-break:break-all">{{ logContent }}</pre>
    </a-modal>

    <!-- 执行输出弹窗 -->
    <a-modal v-model:open="outputOpen" title="执行输出" width="800px" :footer="null">
      <pre style="max-height:500px;overflow:auto;background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:4px">{{ outputText }}</pre>
    </a-modal>

    <!-- 任务详情抽屉 -->
    <a-drawer v-model:open="detailDrawerOpen" :title="detailTitle" width="640px" destroyOnClose>
      <template v-if="detailTask">
        <a-descriptions :column="2" size="small" bordered>
          <a-descriptions-item label="任务名称" :span="2">{{ detailTask.name }}</a-descriptions-item>
          <a-descriptions-item label="数据库类型">{{ detailTask.db_type }}</a-descriptions-item>
          <a-descriptions-item label="实例">{{ detailTask.instance_id }}</a-descriptions-item>
          <a-descriptions-item label="备份方式">{{ detailTask.backup_type }}</a-descriptions-item>
          <a-descriptions-item label="状态"><a-tag :color="statusColor(detailTask.status)">{{ statusLabel(detailTask.status) }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="创建人">{{ detailTask.created_by }}</a-descriptions-item>
          <a-descriptions-item label="开始时间">{{ detailTask.started_at || '-' }}</a-descriptions-item>
          <a-descriptions-item label="完成时间">{{ detailTask.finished_at || '-' }}</a-descriptions-item>
          <a-descriptions-item label="备份文件名" :span="2">{{ detailRec?.file_name || '-' }}</a-descriptions-item>
          <a-descriptions-item label="文件大小">{{ detailRec?.file_size ? (detailRec.file_size/1024/1024).toFixed(2)+'MB' : '-' }}</a-descriptions-item>
          <a-descriptions-item label="存储位置">{{ detailRec?.file_path || '-' }}</a-descriptions-item>
          <a-descriptions-item label="加密密钥" :span="2">{{ detailTask.encrypt_key || '未加密' }}</a-descriptions-item>
          <a-descriptions-item label="恢复命令" :span="2"><code style="word-break:break-all">{{ buildRestoreCmd(detailTask, detailRec) }}</code></a-descriptions-item>
        </a-descriptions>
        <a-divider />
        <div style="font-weight:600;margin-bottom:8px">执行输出</div>
        <pre style="max-height:300px;overflow:auto;background:#1e1e1e;color:#d4d4d4;padding:12px;border-radius:4px">{{ detailRec?.output || '(无输出)' }}</pre>
      </template>
    </a-drawer>
  </a-card>
</template>

<script setup>
import { ref, reactive, computed, onMounted, h } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined, LeftOutlined } from '@ant-design/icons-vue'
import {
  GetBackupConfigsApi, CreateBackupConfigApi, UpdateBackupConfigApi, DeleteBackupConfigApi,
  GetBackupTasksApi, CreateBackupTaskApi, GetBackupRecordsApi, UpdateBackupTaskStatusApi,
  GetBackupTemplatesApi,
} from '@/api/cluster'

const activeTab = ref('configs')

// ── 备份配置 ──
const configColumns = [
  { title: '名称', dataIndex: 'name', width: 160 }, { title: '类型', key: 'db_type', width: 80 },
  { title: '实例', dataIndex: 'instance_id', width: 160 }, { title: '备份方式', key: 'backup_type', width: 80 },
  { title: '保留天数', dataIndex: 'retention_days', width: 70 }, { title: '存储路径', dataIndex: 'storage_path', ellipsis: true },
  { title: '状态', key: 'status', width: 60 }, { title: '操作', key: 'action', width: 200, fixed: 'right' },
]
const configs = ref([]); const cfgLoading = ref(false)
const cfgFilter = reactive({ db_type: '', search: '' })
const cfgModalOpen = ref(false); const isEditConfig = ref(false); const cfgSubmitting = ref(false); const editConfigId = ref(null)
const cfgForm = reactive({ name:'', db_type:'MySQL', instance_id:'', backup_type:'full', retention_days:7, schedule_cron:'', storage_path:'', statusEnabled:true, remark:'' })

const filteredConfigs = computed(() => {
  let items = configs.value
  if (cfgFilter.db_type) items = items.filter(i => i.db_type === cfgFilter.db_type)
  if (cfgFilter.search) items = items.filter(i => (i.name || '').includes(cfgFilter.search))
  return items
})

const fetchConfigs = async () => {
  cfgLoading.value = true; try { const r = await GetBackupConfigsApi(); configs.value = r.data || [] } catch { message.error('获取配置失败') } finally { cfgLoading.value = false }
}
const openConfigModal = (record) => {
  isEditConfig.value = !!record; editConfigId.value = record?.id || null
  Object.assign(cfgForm, record ? {
    name:record.name, db_type:record.db_type, instance_id:record.instance_id, backup_type:record.backup_type,
    retention_days:record.retention_days, schedule_cron:record.schedule_cron||'',
    storage_path:record.storage_path||'', statusEnabled:record.status!=='disabled', remark:record.remark||''
  } : { name:'', db_type:'MySQL', instance_id:'', backup_type:'full', retention_days:7, schedule_cron:'', storage_path:'', statusEnabled:true, remark:'' })
  cfgModalOpen.value = true
}
const handleConfigSubmit = async () => {
  if (!cfgForm.name) { message.warning('请输入名称'); return }
  cfgSubmitting.value = true
  try {
    const data = { name:cfgForm.name, db_type:cfgForm.db_type, instance_id:cfgForm.instance_id, backup_type:cfgForm.backup_type, retention_days:cfgForm.retention_days, schedule_cron:cfgForm.schedule_cron, storage_path:cfgForm.storage_path, remark:cfgForm.remark }
    if (isEditConfig.value) { data.status = cfgForm.statusEnabled ? 'enabled' : 'disabled'; await UpdateBackupConfigApi(editConfigId.value, data); message.success('更新成功') }
    else { await CreateBackupConfigApi(data); message.success('创建成功') }
    cfgModalOpen.value = false; await fetchConfigs()
  } catch (e) { message.error(e.message || '操作失败') } finally { cfgSubmitting.value = false }
}
const handleDeleteConfig = (record) => Modal.confirm({ title:'确认删除', content:`删除「${record.name}」？`, okText:'删除', okType:'danger', onOk:async()=>{ await DeleteBackupConfigApi(record.id); message.success('已删除'); await fetchConfigs() } })

// ── 备份模板 ──
const templates = ref([]); const tmplLoading = ref(false)

const fetchTemplates = async () => {
  tmplLoading.value = true; try { const r = await GetBackupTemplatesApi(); templates.value = r.data || [] } catch {} finally { tmplLoading.value = false }
}
const onTaskTemplateChange = () => {
  const t = templates.value.find(x => x.id === taskForm.template_id)
  if (!t) { taskFormTemplate.value = null; return }
  const fields = JSON.parse(t.config_schema || '[]'); const defaults = JSON.parse(t.default_config || '{}')
  taskFormTemplate.value = { ...t, fields: fields.map(f => ({ ...f, value: defaults[f.key] ?? f.default ?? '' })) }
}
const handleTaskSubmit = async () => {
  taskSubmitting.value = true
  try {
    const config_params = {}
    if (taskFormTemplate.value) taskFormTemplate.value.fields.forEach(f => { if (f.value !== '' && f.value !== undefined) config_params[f.key] = f.value })
    await CreateBackupTaskApi({ config_id:taskCfg.value.id, name:taskForm.name, backup_type:taskCfg.value.backup_type, config_params })
    message.success('任务已创建'); taskModalOpen.value = false; await fetchTasks()
  } catch (e) { message.error(e.message || '创建失败') } finally { taskSubmitting.value = false }
}

// ── 任务（集群维度 ──
const clusterColumns = [
  { title: '类型', key: 'db_type', width: 80 }, { title: '集群/实例', key: 'instance', ellipsis: true },
  { title: '配置数', key: 'config_count', width: 70 }, { title: '总任务', key: 'total', width: 70 },
  { title: '成功', key: 'success', width: 60 }, { title: '失败', key: 'failed', width: 60 },
  { title: '最后备份', key: 'last_time', width: 160 }, { title: '操作', key: 'action', width: 100 },
]
const drillColumns = [
  { title: '任务名称', dataIndex: 'name', width: 180 }, { title: '备份方式', dataIndex: 'backup_type', width: 80 },
  { title: '状态', key: 'status', width: 80 }, { title: '备份文件', key: 'file_name', ellipsis: true },
  { title: '文件大小', key: 'file_size', width: 100 }, { title: '完成时间', key: 'finished_at', width: 150 },
  { title: '操作', key: 'action', width: 60, fixed: 'right' },
]
const tasks = ref([]); const taskLoading = ref(false)
const taskFilter = reactive({ db_type: '', search: '' })
const drillInCluster = ref(null)  // null = cluster summary, obj = drilling into this cluster

const clusterSummary = computed(() => {
  const groups = {}
  tasks.value.forEach(t => {
    const key = `${t.db_type}|${t.instance_id}`
    if (!groups[key]) groups[key] = { key, db_type: t.db_type, instance: t.instance_id, label: `${t.db_type} - ${t.instance_id}`, total:0, success:0, failed:0, running:0, pending:0, config_count:0, last_time: null, last_ts: 0 }
    const g = groups[key]; g.total++
    if (t.status === 'success') g.success++
    else if (t.status === 'failed') g.failed++
    else if (t.status === 'running') g.running++
    else if (t.status === 'pending') g.pending++
    const tTime = t.finished_at || t.started_at || ''
    if (tTime) { const ts = new Date(tTime).getTime(); if (ts > g.last_ts) { g.last_ts = ts; g.last_time = tTime } }
  })
  let items = Object.values(groups)
  if (taskFilter.db_type) items = items.filter(i => i.db_type === taskFilter.db_type)
  if (taskFilter.search) items = items.filter(i => (i.instance||'').includes(taskFilter.search) || (i.db_type||'').includes(taskFilter.search))
  return items
})

const fetchTasks = async () => {
  taskLoading.value = true; try {
    const [t, r] = await Promise.all([GetBackupTasksApi(), GetBackupRecordsApi()])
    tasks.value = (t.data || []).map(t => ({...t, _records: (r.data||[]).filter(rec => rec.task_id === t.id)}))
  } catch (e) { message.error('获取任务失败: ' + (e.message || '')) } finally { taskLoading.value = false }
}
const filterTasks = () => {}
const drillIn = (cluster) => { drillInCluster.value = cluster }
const drillTasks = computed(() => {
  if (!drillInCluster.value) return []
  return tasks.value.filter(t => t.db_type === drillInCluster.value.db_type && t.instance_id === drillInCluster.value.instance)
})

const statusColor = (s) => ({ pending:'orange', running:'processing', success:'green', failed:'red' }[s] || 'default')
const statusLabel = (s) => ({ pending:'待执行', running:'执行中', success:'成功', failed:'失败' }[s] || s)

// ── 任务详情抽屉 ──
const detailDrawerOpen = ref(false); const detailTask = ref(null); const detailRec = ref(null)
const detailTitle = computed(() => detailTask.value ? `任务详情: ${detailTask.value.name}` : '')
const openDetailDrawer = (task) => {
  detailTask.value = task; detailRec.value = (task._records || [])[0] || null
  detailDrawerOpen.value = true
}
const buildRestoreCmd = (task, rec) => {
  if (!rec?.file_path) return '备份完成后可生成恢复命令'
  const fn = rec.file_name || 'backup'
  if (task.db_type === 'MySQL') return `xtrabackup --prepare --target-dir=./${fn}\nxtrabackup --copy-back --target-dir=./${fn} --datadir=/var/lib/mysql`
  if (task.db_type === 'Redis') return `cp ${rec.file_path} /var/lib/redis/dump.rdb\nredis-cli CONFIG SET dir /var/lib/redis\nredis-cli CONFIG SET dbfilename dump.rdb\nredis-cli DEBUG RELOAD`
  if (task.db_type === 'TiDB') return `dumpling --data-path ${rec.file_path} --host ${task.instance_id.split(':')[0]} --port ${task.instance_id.split(':')[1]||4000}`
  return '无恢复命令'
}

// ── 任务日志 ──
const logModalOpen = ref(false); const logTitle = ref(''); const logContent = ref('')
const showTaskLog = (record) => {
  logTitle.value = `任务日志: ${record.name}`
  if (record._records && record._records.length > 0) {
    logContent.value = record._records.map(r => {
      const t = r.started_at ? `[${r.started_at}]` : '[--]'
      const s = r.status === 'success' ? '✓' : r.status === 'failed' ? '✗' : '→'
      return `${t} ${s} 文件:${r.file_name||'-'} 大小:${r.file_size?(r.file_size/1024/1024).toFixed(2)+'MB':'-'}\n${r.output||'(无输出)'}`
    }).join('\n---\n')
  } else { logContent.value = '(暂无执行记录)' }
  logModalOpen.value = true
}

// ── 执行输出 ──
const outputOpen = ref(false); const outputText = ref('')
const showOutput = (record) => { outputText.value = record.output || '(无输出)'; outputOpen.value = true }

onMounted(() => { fetchConfigs(); fetchTemplates(); fetchTasks() })
</script>