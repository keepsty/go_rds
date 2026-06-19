const route = {
  name: 'view.salt',
  path: '/salt',
  component: () => import('./index.vue'),
  redirect: '/salt/config',
  meta: { title: 'Salt部署', icon: 'ThunderboltOutlined', keepAlive: true },
  children: [
    { name: 'view.salt.config', path: '/salt/config', component: () => import('./config/index.vue'), meta: { title: '主机配置', keepAlive: true } },
    { name: 'view.salt.taskCreate', path: '/salt/task-create', component: () => import('./task-create/index.vue'), meta: { title: '创建任务', keepAlive: true } },
    { name: 'view.salt.taskRun', path: '/salt/task-run', component: () => import('./task-run/index.vue'), meta: { title: '运行任务', keepAlive: true } },
    { name: 'view.salt.history', path: '/salt/history', component: () => import('./history/index.vue'), meta: { title: '历史记录', keepAlive: true } },
  ],
}
export default route