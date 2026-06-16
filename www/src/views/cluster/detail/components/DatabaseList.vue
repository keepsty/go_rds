<template>
  <div class="database-list">
    <div class="toolbar-actions">
      <a-button type="primary" size="small" @click="refreshList">刷新</a-button>
    </div>
    <a-table
      size="middle"
      :columns="columns"
      :row-key="(record) => record.id"
      :data-source="databaseList"
      :loading="loading"
      :pagination="false"
      :scroll="{ x: 800 }"
    >
      <template #expandableRow="{ record }">
        <a-table
          size="small"
          :columns="tableColumns"
          :data-source="record.tables || []"
          :pagination="false"
          :row-key="(tb) => tb.id"
        >
          <template #bodyCell="{ column: col, record: tb }">
            <template v-if="col.key === 'table_name'">
              <a @click="showTableInfo(record, tb)">{{ tb.table_name }}</a>
            </template>
          </template>
        </a-table>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'">
          <span class="db-name">{{ record.name }}</span>
        </template>
        <template v-if="column.key === 'action'">
          <a-button type="link" size="small" @click="handleDBOperation(record)">
            库迁移
          </a-button>
          <a-button type="link" size="small" @click="showDBDict(record)">
            数据字典
          </a-button>
        </template>
      </template>
    </a-table>
  </div>
</template>

<script setup>
defineOptions({ name: 'DatabaseList' })

import { ref, onMounted } from 'vue'
import { GetClusterDatabasesApi } from '@/api/cluster'

const props = defineProps({
  sgId: { type: Number, required: true },
})

const databaseList = ref([])
const loading = ref(false)

const columns = [
  { title: '数据库名', key: 'name', width: 200 },
  { title: '数据空间', key: 'database_size', width: 120 },
  { title: '字符集', key: 'database_charset', width: 120 },
  { title: '业务RD', key: 'rd_user', width: 160 },
  { title: '创建时间', key: 'create_time', width: 170 },
  { title: '操作', key: 'action', width: 180, fixed: 'right' },
]

const tableColumns = [
  { title: '表名', key: 'table_name', width: 180 },
  { title: '表大小(G)', key: 'table_size', width: 100 },
  { title: '空洞大小(G)', key: 'free_size', width: 100 },
  { title: '字符集', key: 'table_collation', width: 100 },
  { title: '行数', key: 'table_rows', width: 80 },
  { title: '自增值', key: 'auto_increase', width: 80 },
]

async function getDatabaseList() {
  loading.value = true
  try {
    const res = await GetClusterDatabasesApi(props.sgId)
    databaseList.value = res.data || []
  } catch {
    databaseList.value = []
  } finally {
    loading.value = false
  }
}

function showTableInfo(db, table) {
  console.log('show table info:', db.name, table.table_name)
}

function showDBDict(db) {
  console.log('show data dict:', db.name)
}

function handleDBOperation(record) {
  console.log('migrate db:', record.name)
}

function refreshList() {
  getDatabaseList()
}

onMounted(() => {
  getDatabaseList()
})
</script>

<style scoped>
.database-list {
  margin-top: 8px;
}

.toolbar-actions {
  margin-bottom: 12px;
}

.db-name {
  font-weight: 500;
}
</style>
