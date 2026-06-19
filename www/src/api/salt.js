import { del, get, post, put } from '@/utils/request'

// 普通用户可用
export const getTemplatesApi = () => get('/api/v1/salt/templates')
export const deployTaskApi = (data) => post('/api/v1/salt/templates/deploy', data)
export const getMinionsApi = () => get('/api/v1/salt/minions')

export const getHostConfigsApi = () => get('/api/v1/salt/host-configs')

export const getTasksApi = (params) => get('/api/v1/salt/tasks', params)
export const createTaskApi = (data) => post('/api/v1/salt/tasks', data)
export const runTaskApi = (id) => post(`/api/v1/salt/tasks/${id}/run`)

// 管理员专用
export const createTemplateApi = (data) => post('/api/v1/admin/salt/templates', data)
export const updateTemplateApi = (id, data) => put(`/api/v1/admin/salt/templates/${id}`, data)
export const deleteTemplateApi = (id) => del(`/api/v1/admin/salt/templates/${id}`)

export const createHostConfigApi = (data) => post('/api/v1/admin/salt/host-configs', data)
export const updateHostConfigApi = (id, data) => put(`/api/v1/admin/salt/host-configs/${id}`, data)
export const deleteHostConfigApi = (id) => del(`/api/v1/admin/salt/host-configs/${id}`)

export const approveTaskApi = (id, data) => put(`/api/v1/admin/salt/tasks/${id}/approve`, data)