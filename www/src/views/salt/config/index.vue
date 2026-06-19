<template>
  <div>
    <a-card title="主机配置">
      <template #extra>
        <a-button type="primary" @click="openModal()"><template #icon><PlusOutlined /></template>新增配置</a-button>
      </template>
      <a-table :dataSource="list" :columns="columns" :loading="loading" :pagination="false" rowKey="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'hosts'">
            <a-tag v-for="h in (typeof record.hosts === 'string' ? JSON.parse(record.hosts) : record.hosts || [])" :key="h.host">{{ h.host }}:{{ h.port || '-' }}</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" @click="openModal(record)">编辑</a-button>
            <a-button type="link" danger @click="handleDelete(record)">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-card>
    <a-modal v-model:open="modalOpen" :title="isEdit ? '编辑配置' : '新增配置'" width="640px" :confirm-loading="submitting" @ok="handleSubmit" destroyOnClose>
      <a-form layout="vertical">
        <a-form-item label="配置名称" required><a-input v-model:value="form.name" placeholder="如：生产MySQL集群" /></a-form-item>
        <a-form-item label="描述"><a-textarea v-model:value="form.description" :rows="2" /></a-form-item>
        <a-form-item label="主机列表（JSON）" required :help="hostHelpText">
          <a-textarea v-model:value="form.hosts" :rows="5" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getHostConfigsApi, createHostConfigApi, updateHostConfigApi, deleteHostConfigApi } from '@/api/salt'

const columns = [
  { title: '名称', dataIndex: 'name', key: 'name' },
  { title: '主机', key: 'hosts' },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '操作', key: 'action', width: 120, fixed: 'right' },
]
const hostHelpText = '格式: [{"host":"minion-01","port":3306,"role":"master"},...]'
const list = ref([]); const loading = ref(false); const modalOpen = ref(false)
const isEdit = ref(false); const submitting = ref(false); const editId = ref(null)
const form = ref({ name: '', hosts: '', description: '' })

const fetchData = async () => {
  loading.value = true; try { const r = await getHostConfigsApi(); list.value = r.data || [] } catch { message.error('获取失败') } finally { loading.value = false }
}
const openModal = (record) => {
  isEdit.value = !!record; editId.value = record?.id || null
  if (record) {
    form.value = { name: record.name, description: record.description || '', hosts: typeof record.hosts === 'string' ? record.hosts : JSON.stringify(record.hosts || []) }
  } else { form.value = { name: '', hosts: '[]', description: '' } }
  modalOpen.value = true
}
const handleSubmit = async () => {
  if (!form.value.name) { message.warning('请输入名称'); return }; submitting.value = true
  try {
    if (isEdit.value) { await updateHostConfigApi(editId.value, form.value) } else { await createHostConfigApi(form.value) }
    message.success(isEdit.value ? '更新成功' : '创建成功'); modalOpen.value = false; await fetchData()
  } catch (e) { message.error(e.message || '操作失败') } finally { submitting.value = false }
}
const handleDelete = (record) => Modal.confirm({ title: '确认删除', content: `删除「${record.name}」？`, okText: '删除', okType: 'danger', onOk: async () => { await deleteHostConfigApi(record.id); message.success('已删除'); await fetchData() } })
onMounted(fetchData)
</script>