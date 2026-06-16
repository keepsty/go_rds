<template>
  <div class="proxysql-list">
    <a-table
      size="middle"
      :columns="columns"
      :row-key="(record) => record.id"
      :data-source="proxyList"
      :loading="loading"
      :pagination="false"
      :scroll="{ x: 900 }"
    >
      <template #expandableRow="{ record }">
        <a-table
          size="small"
          :columns="serverColumns"
          :data-source="record.proxy_mysql_server || []"
          :pagination="false"
          :row-key="(srv) => srv.proxysql_id + '-' + srv.mysql_hostname"
        />
      </template>
    </a-table>
  </div>
</template>

<script setup>
defineOptions({ name: 'ProxySQLList' })

import { ref, onMounted } from 'vue'
import { GetClusterProxyListApi } from '@/api/cluster'

const props = defineProps({
  sgId: { type: Number, required: true },
})

const proxyList = ref([])
const loading = ref(false)

const columns = [
  { title: '主机名', key: 'hostname', width: 160 },
  { title: '管理端口', key: 'admin_port', width: 100 },
  { title: '应用端口', key: 'app_port', width: 100 },
  { title: '权重', key: 'proxy_weight', width: 80 },
  { title: '版本', key: 'proxy_version', width: 100 },
  { title: '路由规则', key: 'rule_info', width: 200 },
]

const serverColumns = [
  { title: 'MySQL主机', key: 'mysql_hostname', width: 160 },
  { title: '端口', key: 'port', width: 80 },
  { title: '状态', key: 'status', width: 80 },
  { title: '权重', key: 'weight', width: 80 },
  { title: 'Hostgroup', key: 'hostgroup_id', width: 100 },
]

async function getProxyList() {
  loading.value = true
  try {
    const res = await GetClusterProxyListApi(props.sgId)
    proxyList.value = res.data || []
  } catch {
    proxyList.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getProxyList()
})
</script>

<style scoped>
.proxysql-list {
  margin-top: 8px;
}
</style>
