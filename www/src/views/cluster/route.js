const route = {
  name: 'cluster',
  path: '/cluster',
  redirect: '/cluster/list',
  component: () => import('./index.vue'),
  meta: { title: '集群管理', icon: 'ClusterOutlined', keepAlive: true },
  children: [
    {
      name: 'cluster.list',
      path: '/cluster/list',
      component: () => import('./list/ClusterList.vue'),
      meta: { title: '集群列表', icon: 'UnorderedListOutlined', keepAlive: true },
    },
    {
      name: 'cluster.detail',
      path: '/cluster/:sg_id',
      component: () => import('./detail/ClusterDetail.vue'),
      meta: { title: '集群详情', keepAlive: true, hidden: true },
    },
    {
      name: 'cluster.backup',
      path: '/cluster/backup',
      component: () => import('./backup/index.vue'),
      meta: { title: '备份管理', icon: 'CloudUploadOutlined', keepAlive: true },
    },
  ],
}

export default route
