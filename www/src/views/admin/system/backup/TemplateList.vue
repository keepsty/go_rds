<template>
  <div>
    <a-card title="备份模板管理">
      <template #extra>
        <a-button type="primary" @click="openModal()"><PlusOutlined />新增模板</a-button>
      </template>
      <a-table :dataSource="list" :columns="columns" :loading="loading" :pagination="false" rowKey="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'db_type'"><a-tag :color="{MySQL:'blue',Redis:'green',TiDB:'purple'}[record.db_type]">{{ record.db_type }}</a-tag></template>
          <template v-if="column.key === 'fields'"><a-tag>{{ fieldCount(record) }} 个字段</a-tag></template>
          <template v-if="column.key === 'action'">
            <a-button type="link" @click="openModal(record)">编辑</a-button>
            <a-button type="link" danger @click="handleDelete(record)">删除</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="modalOpen" :title="isEdit?'编辑模板':'新增模板'" width="800px" :confirm-loading="submitting" @ok="handleSubmit" destroyOnClose>
      <a-form layout="vertical">
        <a-form-item label="模板名称" required><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item label="数据库类型" required>
          <a-select v-model:value="form.db_type" :options="[{value:'MySQL',label:'MySQL'},{value:'Redis',label:'Redis'},{value:'TiDB',label:'TiDB'}]" />
        </a-form-item>
        <a-form-item label="描述"><a-textarea v-model:value="form.description" :rows="2" /></a-form-item>
        <a-form-item label="字段定义（YAML格式）" required help="每行一个字段，从 - key:xxx 开始缩进">
          <a-textarea v-model:value="form.config_schema" :rows="12" style="font-family:monospace;font-size:13px;line-height:1.6" />
        </a-form-item>
        <a-collapse ghost v-if="form.config_schema.trim()">
          <a-collapse-panel header="校验结果（可预览解析后的JSON）">
            <pre style="max-height:200px;overflow:auto;background:#f5f5f5;padding:8px;border-radius:4px;font-size:13px">{{ previewJson }}</pre>
          </a-collapse-panel>
        </a-collapse>
        <a-form-item label="默认值（JSON）"><a-textarea v-model:value="form.default_config" :rows="4" style="font-family:monospace;font-size:13px" /></a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { GetBackupTemplatesApi, CreateBackupTemplateApi, UpdateBackupTemplateApi, DeleteBackupTemplateApi } from '@/api/cluster'

const columns = [
  { title: '名称', dataIndex: 'name', width: 180 }, { title: '类型', key: 'db_type', width: 80 },
  { title: '描述', dataIndex: 'description', ellipsis: true }, { title: '字段', key: 'fields', width: 80 },
  { title: '操作', key: 'action', width: 120, fixed: 'right' },
]
const list = ref([]); const loading = ref(false)
const modalOpen = ref(false); const isEdit = ref(false); const submitting = ref(false); const editId = ref(null)
const form = ref({ name:'', db_type:'MySQL', description:'', config_schema:'', default_config:'' })

// ── JSON ↔ YAML 转换 ──
function jsonToYaml(jsonStr) {
  try {
    const arr = JSON.parse(jsonStr)
    if (!Array.isArray(arr)) return jsonStr
    return arr.map(f => {
      const lines = [`- key: ${f.key}`, `  label: ${f.label}`, `  type: ${f.type}`]
      if (f.required) lines.push('  required: true')
      if (f.default !== undefined && f.default !== '') lines.push(`  default: ${f.default}`)
      if (f.description) lines.push(`  description: "${f.description}"`)
      if (f.options && f.options.length) lines.push(`  options: [${f.options.join(', ')}]`)
      return lines.join('\n')
    }).join('\n')
  } catch { return jsonStr }
}

function yamlToJson(yamlStr) {
  try {
    const lines = yamlStr.split('\n').filter(l => l.trim())
    const fields = []; let current = null
    for (const line of lines) {
      const trimmed = line.trim()
      if (trimmed.startsWith('- key:')) {
        if (current) fields.push(current)
        current = { key: trimmed.replace('- key:', '').trim() }
      } else if (current) {
        const colonIdx = trimmed.indexOf(':')
        if (colonIdx === -1) continue
        const k = trimmed.slice(0, colonIdx).trim()
        let v = trimmed.slice(colonIdx + 1).trim()
        if (v.startsWith('"') && v.endsWith('"')) v = v.slice(1, -1)
        if (k === 'required' || k === 'default') {
          if (v === 'true') v = true
          else if (v === 'false') v = false
          else if (!isNaN(Number(v))) v = Number(v)
        }
        if (k === 'options') {
          v = v.replace(/^\[|\]$/g, '').split(',').map(s => s.trim())
        }
        current[k] = v
      }
    }
    if (current) fields.push(current)
    return JSON.stringify(fields)
  } catch { return '[]' }
}

const fieldCount = (r) => { try { const arr = JSON.parse(typeof r.config_schema === 'string' ? r.config_schema : '[]'); return Array.isArray(arr) ? arr.length : 0 } catch { return 0 } }

const previewJson = computed(() => {
  try {
    const json = yamlToJson(form.value.config_schema)
    const arr = JSON.parse(json)
    return JSON.stringify(arr, null, 2)
  } catch { return form.value.config_schema ? '(YAML格式有误)' : '(空)' }
})

const fetchData = async () => {
  loading.value = true; try { const r = await GetBackupTemplatesApi(); list.value = r.data || [] } catch { message.error('获取失败') } finally { loading.value = false }
}

const openModal = (record) => {
  isEdit.value = !!record; editId.value = record?.id || null
  if (record) {
    const rawSchema = typeof record.config_schema === 'string' ? record.config_schema : JSON.stringify(record.config_schema || [])
    form.value = {
      name: record.name, db_type: record.db_type, description: record.description || '',
      config_schema: jsonToYaml(rawSchema),
      default_config: typeof record.default_config === 'string' ? record.default_config : JSON.stringify(record.default_config || {}, null, 2),
    }
  } else { form.value = { name:'', db_type:'MySQL', description:'', config_schema:'', default_config:'' } }
  modalOpen.value = true
}

const handleSubmit = async () => {
  if (!form.value.name || !form.value.config_schema.trim()) { message.warning('请填写必填项'); return }
  const jsonSchema = yamlToJson(form.value.config_schema)
  try {
    JSON.parse(jsonSchema) // validate
  } catch {
    message.error('YAML 格式错误，请检查缩进和格式')
    return
  }
  submitting.value = true
  try {
    const data = { name:form.value.name, db_type:form.value.db_type, description:form.value.description, config_schema: jsonSchema, default_config: form.value.default_config }
    if (isEdit.value) { await UpdateBackupTemplateApi(editId.value, data); message.success('更新成功') }
    else { await CreateBackupTemplateApi(data); message.success('创建成功') }
    modalOpen.value = false; await fetchData()
  } catch (e) { message.error(e.message || '操作失败') } finally { submitting.value = false }
}

const handleDelete = (record) => Modal.confirm({ title:'确认删除', content:`删除模板「${record.name}」？`, okText:'删除', okType:'danger', onOk:async()=>{ await DeleteBackupTemplateApi(record.id); message.success('已删除'); await fetchData() } })

onMounted(fetchData)
</script>