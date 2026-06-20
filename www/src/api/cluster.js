import { get, post, put, del } from "@/utils/request"

// 集群管理
export const GetClustersApi = (params) => get('/api/v1/cluster/clusters', params)
export const GetClusterDetailApi = (id) => get(`/api/v1/cluster/clusters/${id}`)
export const GetClusterInstancesApi = (id) => get(`/api/v1/cluster/clusters/instances/${id}`)
export const GetClusterDatabasesApi = (id) => get(`/api/v1/cluster/clusters/databases/${id}`)
export const GetClusterDatabaseDetailApi = (id) => get(`/api/v1/cluster/clusters/database-detail/${id}`)
export const GetClusterProxyListApi = (id) => get(`/api/v1/cluster/clusters/proxysql/${id}`)
export const SetInstanceStatusApi = (data) => post('/api/v1/cluster/clusters/instances/status', data)

// 数据库树形选项与表
export const GetDBOptionsApi = (username) => get(`/api/v1/cluster/clusters/dboptions/${username}`)
export const GetDBTablesApi = (sgId, dbId) => get(`/api/v1/cluster/clusters/dbtables/${sgId}/${dbId}`)
export const GetDBsApi = (id) => get(`/api/v1/cluster/databases/${id}`)

// SQL 查询与数据字典
export const ExecuteQueryApi = (data) => post('/api/v1/cluster/clusters/db/query/execute', data)
export const GetDataDictApi = (data) => post('/api/v1/cluster/clusters/db/query/datadict', data)
export const GetTableInfoApi = (params) => get('/api/v1/cluster/clusters/db/query/tableinfo', params)
export const GetQueryHistoryApi = (params) => get('/api/v1/cluster/clusters/db/query/history', params)

// SaltStack 自动化部署
export const AddMySQLClusterApi = (data) => post('/api/v1/cluster/clusters/addmysqlcluster', data)

// 备份管理
export const GetBackupConfigsApi = (params) => get('/api/v1/cluster/backup/configs', params)
export const CreateBackupConfigApi = (data) => post('/api/v1/cluster/backup/configs', data)
export const UpdateBackupConfigApi = (id, data) => put(`/api/v1/cluster/backup/configs/${id}`, data)
export const DeleteBackupConfigApi = (id) => del(`/api/v1/cluster/backup/configs/${id}`)
export const GetBackupTasksApi = (params) => get('/api/v1/cluster/backup/tasks', params)
export const CreateBackupTaskApi = (data) => post('/api/v1/cluster/backup/tasks', data)
export const UpdateBackupTaskStatusApi = (id, data) => put(`/api/v1/cluster/backup/tasks/${id}/status`, data)
export const GetBackupRecordsApi = (params) => get('/api/v1/cluster/backup/records', params)

// 备份模板管理
export const GetBackupTemplatesApi = (params) => get('/api/v1/cluster/backup/templates', params)
export const CreateBackupTemplateApi = (data) => post('/api/v1/cluster/backup/templates', data)
export const UpdateBackupTemplateApi = (id, data) => put(`/api/v1/cluster/backup/templates/${id}`, data)
export const DeleteBackupTemplateApi = (id) => del(`/api/v1/cluster/backup/templates/${id}`)
