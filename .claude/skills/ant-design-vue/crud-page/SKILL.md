
# CRUD Table Page（企业后台列表页开发）

你是一个专注于 **Vue 3 + JavaScript + Ant Design Vue 4.x** 列表页开发专家。目标是输出结构清晰、交互统一、可直接落地的后台列表页面代码。

---

## When to Apply

- 新建后台列表页 / CRUD 页面
- 开发或改造 table 页面
- 优化现有列表页结构

## 不涵盖的内容

- 表单与弹窗细节 → 参考 **go-rds-backend-structure** 的 form-modal
- 父子组件数据传递 → 参考 **component-binding**
- 响应式适配 → 参考 **responsive-layout**

---

## 页面结构规范

```
page/ → 查询区 → 操作区 → 表格区 → 分页区 → 新增/编辑弹窗区
```

- 侧边栏宽度建议 `200px ~ 240px`
- 主内容区优先用 `a-card` 承载
- 页面内边距建议 `16px ~ 24px`

---

## Best Practices

### 表格 (Table)
- 使用标准 `columns` 配置，`rowKey` 必填
- 操作列固定右侧（放最后）
- 长文本做省略或 tooltip
- 状态字段尽量可视化（tag / switch）
- 危险操作必须 `modal.confirm()`

### 查询区
- 查询条件优先展示高频项，低频折叠
- 查询表单和编辑表单分开管理
- 查询区优先使用 `flex` + `flex-wrap`

### 数据流
- 查询、新增、编辑、删除方法分开
- 加载态 `tableLoading`、提交中 `submitLoading` 明确
- 成功后给反馈，刷新列表或关闭弹窗

---

## 状态与方法命名

- 状态：`tableData` / `queryForm` / `pagination` / `tableLoading` / `submitLoading`
- 方法：`getList` / `handleSearch` / `handleReset` / `handleCreate` / `handleEdit` / `handleDelete`

---

## 技术约束

- Vue 3 + `script setup`，JavaScript（非 TypeScript），Ant Design Vue 4.x
