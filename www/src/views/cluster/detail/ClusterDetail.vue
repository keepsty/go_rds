<template>
  <div class="gi-page-shell cluster-detail-page">
    <a-card class="cluster-detail-card" title="集群详情" :loading="loading">
      <template #extra>
        <a-button type="link" @click="goBack">← 返回列表</a-button>
      </template>

      <!-- 描述信息 -->
      <a-descriptions bordered :column="3" size="small">
        <a-descriptions-item label="集群名" :span="1">{{ clusterDetail.sg_name }}</a-descriptions-item>
        <a-descriptions-item label="集群描述" :span="2">{{ clusterDetail.cluster_description }}</a-descriptions-item>
        <a-descriptions-item label="产品线">{{ clusterDetail.prod_name }}</a-descriptions-item>
        <a-descriptions-item label="业务负责人">{{ clusterDetail.rd_owner }}</a-descriptions-item>
        <a-descriptions-item label="负责DBA">{{ clusterDetail.dba_user }}</a-descriptions-item>
        <a-descriptions-item label="集群等级">{{ clusterDetail.service_level }}</a-descriptions-item>
        <a-descriptions-item label="环境">
          <a-tag :color="envColorMap[clusterDetail.environment] || 'default'">
            {{ envMap[clusterDetail.environment] || clusterDetail.environment }}
          </a-tag>
        </a-descriptions-item>
        <a-descriptions-item label="中间件">{{ middlewareMap[clusterDetail.middleware] || clusterDetail.middleware }}</a-descriptions-item>
        <a-descriptions-item label="HA类型">{{ haTypeMap[clusterDetail.ha_type] || clusterDetail.ha_type }}</a-descriptions-item>
        <a-descriptions-item label="创建时间">{{ clusterDetail.create_time }}</a-descriptions-item>
        <a-descriptions-item label="集群高峰期">{{ clusterDetail.peak_time }}</a-descriptions-item>
        <a-descriptions-item label="DNS名称">{{ clusterDetail.dns_name }}</a-descriptions-item>
        <a-descriptions-item label="VIP">{{ clusterDetail.vip }}</a-descriptions-item>
      </a-descriptions>

      <!-- Tabs 切换 -->
      <a-tabs v-model:activeKey="activeTab" class="detail-tabs" :style="{ marginTop: '16px' }">
        <a-tab-pane key="instances" tab="实例列表">
          <InstanceList :sg-id="sgId" />
        </a-tab-pane>
        <a-tab-pane key="databases" tab="库表管理">
          <DatabaseList :sg-id="sgId" />
        </a-tab-pane>
        <a-tab-pane key="proxysql" tab="中间件管理">
          <ProxySQLList :sg-id="sgId" />
        </a-tab-pane>
      </a-tabs>
    </a-card>
  </div>
</template>

<script setup>
defineOptions({ name: 'ClusterDetailView' })

import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { GetClusterDetailApi } from '@/api/cluster'
import InstanceList from './components/InstanceList.vue'
import DatabaseList from './components/DatabaseList.vue'
import ProxySQLList from './components/ProxySQLList.vue'

const route = useRoute()
const router = useRouter()

const sgId = ref(parseInt(route.params.sg_id, 10))
const loading = ref(false)
const activeTab = ref('instances')

const clusterDetail = ref({
  sg_name: '',
  cluster_description: '',
  prod_name: '',
  rd_owner: '',
  dba_user: '',
  service_level: '',
  environment: 0,
  middleware: 0,
  ha_type: 0,
  create_time: '',
  peak_time: '',
  dns_name: '',
  vip: '',
})

const envMap = { 0: 'prod', 1: 'rc', 2: 'k8s', 3: 'press' }
const envColorMap = { 0: 'red', 1: 'orange', 2: 'blue', 3: 'purple' }
const haTypeMap = { 0: 'MHA', 1: 'ORC', 2: 'MGR' }
const middlewareMap = { 0: 'ProxySQL', 1: 'Zebra', 2: 'MGW' }

async function getClusterDetail() {
  loading.value = true
  try {
    const res = await GetClusterDetailApi(sgId.value)
    if (res.data && res.data.length > 0) {
      clusterDetail.value = res.data[0]
    }
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push({ name: 'cluster.list' })
}

onMounted(() => {
  getClusterDetail()
})
</script>

<style scoped>
.cluster-detail-card {
  border-radius: 12px;
  border: 1px solid var(--ant-colorBorderSecondary, #f0f0f0);
  box-shadow: 0 2px 8px rgb(0 0 0 / 5%);
}

.detail-tabs :deep(.ant-tabs-nav) {
  margin-bottom: 10px;
}
</style>
