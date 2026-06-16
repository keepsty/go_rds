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
