<template>
  <div>
    <a-card title="创建部署任务">
      <a-form layout="vertical">
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="任务名称" required><a-input v-model:value="form.name" placeholder="如：生产MySQL初始化部署" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="部署模板" required>
              <a-select v-model:value="form.template_name" placeholder="选择模板" :options="templates.map(t=>({value:t.name,label:t.title}))" @change="onTemplateChange" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="主机配置" required>
              <a-select v-model:value="form.host_config_id" placeholder="选择主机配置" :options="hostConfigs.map(h=>({value:h.id,label:h.name}))" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-divider>模板参数</a-divider>
        <a-row :gutter="24" v-if="currentTemplate">
          <a-col v-for="f in currentTemplate.fields" :key="f.key" :span="f.type === 'text' ? 24 : 12">
            <a-form-item :label="f.label" :required="f.required" :help="f.description">
              <a-input v-if="f.type==='string'" v-model:value="f.value" />
              <a-input-number v-else-if="f.type==='number'" v-model:value="f.value" style="width:100%" />
              <a-switch v-else-if="f.type==='boolean'" v-model:checked="f.value" />
              <a-textarea v-else-if="f.type==='text'" v-model:value="f.value" :rows="3" />
            </a-form-item>
          </a-col>
          <a-empty v-if="!currentTemplate.fields.length" description="该模板无参数" />
        </a-row>
        <a-empty v-else description="请先选择部署模板" />
      </a-form>
      <div style="text-align:right;margin-top:16px">
        <a-button type="primary" :loading="submitting" :disabled="!canSubmit" @click="handleSubmit">创建任务</a-button>
      </div>
    </a-card>
  </div>
</template>
<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getTemplatesApi, getHostConfigsApi, createTaskApi } from '@/api/salt'

const templates = ref([]); const hostConfigs = ref([]); const submitting = ref(false)
const currentTemplate = ref(null)
const form = reactive({ name: '', template_name: undefined, host_config_id: undefined, config_params: {} })
const canSubmit = computed(() => form.name && form.template_name && form.host_config_id)

const onTemplateChange = () => {
  const t = templates.value.find(x => x.name === form.template_name); if (!t) { currentTemplate.value = null; return }
  const fields = JSON.parse(t.fields_schema || '[]'); const defaults = JSON.parse(t.defaults || '{}')
  currentTemplate.value = { ...t, fields: fields.map(f => ({ ...f, value: defaults[f.key] ?? f.default ?? '' })) }
}

const handleSubmit = async () => {
  submitting.value = true
  try {
    const config = {}; (currentTemplate.value?.fields || []).forEach(f => { config[f.key] = f.value })
    await createTaskApi({ name: form.name, template_name: form.template_name, host_config_id: form.host_config_id, config_params: config })
    message.success('任务已创建，等待管理员审批'); form.name = ''; form.template_name = undefined; form.host_config_id = undefined; currentTemplate.value = null
  } catch (e) { message.error(e.message || '创建失败') }
  finally { submitting.value = false }
}

onMounted(async () => {
  try { const [t, h] = await Promise.all([getTemplatesApi(), getHostConfigsApi()]); templates.value = t.data || []; hostConfigs.value = h.data || [] } catch { message.warning('加载数据失败') }
})
</script>