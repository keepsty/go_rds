<template>
  <div>
    <a-card title="Salt 部署模版管理">
      <template #extra>
        <a-button type="primary" @click="openModal()">
          <template #icon><PlusOutlined /></template>
          新增模版
        </a-button>
      </template>
      <a-table
        :dataSource="tableData"
        :columns="columns"
        :loading="loading"
        :pagination="false"
        rowKey="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'fields_schema'">
            <a-tag>{{ JSON.parse(record.fields_schema || '[]').length }} 个字段</a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" @click="openModal(record)">编辑</a-button>
            <a-button type="link" danger @click="handleDelete(record)">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="modalOpen"
      :title="isEdit ? '编辑模版' : '新增模版'"
      width="720px"
      :confirm-loading="submitting"
      @ok="handleSubmit"
      destroyOnClose
    >
      <a-form layout="vertical">
        <a-form-item label="标识" required v-if="!isEdit">
          <a-input v-model:value="formState.name" placeholder="唯一标识，如 cmd-run" />
        </a-form-item>
        <a-form-item label="标题" required>
          <a-input v-model:value="formState.title" placeholder="远程执行命令" />
        </a-form-item>
        <a-form-item label="描述">
          <a-textarea v-model:value="formState.description" :rows="2" />
        </a-form-item>
        <a-form-item label="字段定义（JSON）" required>
          <a-textarea
            v-model:value="formState.fields_schema"
            :rows="6"
            placeholder='[{"key":"command","label":"命令","type":"text","required":true}]'
          />
        </a-form-item>
        <a-form-item label="默认值（JSON）">
          <a-textarea
            v-model:value="formState.defaults"
            :rows="3"
            placeholder='{"saltenv":"base"}'
          />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getTemplatesApi, createTemplateApi, updateTemplateApi, deleteTemplateApi } from '@/api/salt'

const columns = [
  { title: '标识', dataIndex: 'name', key: 'name', width: 160 },
  { title: '标题', dataIndex: 'title', key: 'title' },
  { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
  { title: '字段', key: 'fields_schema', width: 80 },
  { title: '操作', key: 'action', width: 120, fixed: 'right' },
]

const tableData = ref([])
const loading = ref(false)
const modalOpen = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref(null)
const formState = reactive({
  name: '', title: '', description: '', fields_schema: '', defaults: '',
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getTemplatesApi()
    tableData.value = res.data || []
  } catch { message.error('获取模版列表失败') }
  finally { loading.value = false }
}

const openModal = (record) => {
  isEdit.value = !!record
  editId.value = record?.id || null
  if (record) {
    formState.name = record.name
    formState.title = record.title
    formState.description = record.description || ''
    formState.fields_schema = typeof record.fields_schema === 'string' ? record.fields_schema : JSON.stringify(record.fields_schema || [])
    formState.defaults = typeof record.defaults === 'string' ? record.defaults : JSON.stringify(record.defaults || {})
  } else {
    formState.name = ''; formState.title = ''; formState.description = ''; formState.fields_schema = ''; formState.defaults = ''
  }
  modalOpen.value = true
}

const handleSubmit = async () => {
  if (!formState.name && !isEdit.value) { message.warning('请输入标识'); return }
  if (!formState.title) { message.warning('请输入标题'); return }
  if (!formState.fields_schema) { message.warning('请输入字段定义'); return }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateTemplateApi(editId.value, {
        title: formState.title, description: formState.description,
        fields_schema: formState.fields_schema, defaults: formState.defaults,
      })
      message.success('更新成功')
    } else {
      await createTemplateApi({
        name: formState.name, title: formState.title, description: formState.description,
        fields_schema: formState.fields_schema, defaults: formState.defaults,
      })
      message.success('创建成功')
    }
    modalOpen.value = false
    await fetchData()
  } catch (err) { message.error(err.message || '操作失败') }
  finally { submitting.value = false }
}

const handleDelete = (record) => {
  Modal.confirm({
    title: '确认删除',
    content: `确定删除模版「${record.title}」吗？`,
    okText: '删除',
    okType: 'danger',
    onOk: async () => {
      try {
        await deleteTemplateApi(record.id)
        message.success('删除成功')
        await fetchData()
      } catch (err) { message.error(err.message || '删除失败') }
    },
  })
}

onMounted(fetchData)
</script>