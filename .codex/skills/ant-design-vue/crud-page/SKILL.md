---
name: ant-design-vue-crud-page
description: >
  Build production-grade CRUD table pages with Vue 3 + JavaScript + Ant Design Vue 4.x.
  Use when creating or refactoring list pages, table pages, search + table + modal pages,
  paginated data views, or any backend admin list interface.
  This is part of the ant-design-vue skill family — for forms/modals, component binding,
  or responsive layout, see the sibling skills.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# CRUD Table Page（企业后台列表页开发）

你是一个专注于 **Vue 3 + JavaScript + Ant Design Vue 4.x** 列表页开发专家。
目标是输出结构清晰、交互统一、可直接落地的后台列表页面代码。

---

## When to Apply

- 新建后台列表页 / CRUD 页面
- 开发或改造 table 页面
- 优化现有列表页结构
- 用户提到"列表"、"表格"、"分页"、"CRUD"、"管理页面"等关键词

## 不涵盖的内容

- 表单与弹窗细节 → 参考 [[ant-design-vue-form-modal]]
- 父子组件数据传递 → 参考 [[ant-design-vue-component-binding]]
- 响应式适配 → 参考 [[ant-design-vue-responsive-layout]]

---

## 页面结构规范

### 推荐页面蓝图

```
page/
├── 查询区
├── 操作区
├── 表格区
├── 分页区
└── 新增/编辑弹窗区
```

### 推荐后台布局

```text
Layout
 |- Sider
 |- Header
 `- Content
```

- 使用 `Layout / Layout.Sider / Layout.Header / Layout.Content`
- 侧边栏宽度建议 `200px ~ 240px`
- 主内容区优先用 `a-card` 承载
- 页面内边距建议 `16px ~ 24px`

---

## Best Practices

### 表格 (Table)

- 使用标准 `columns` 配置
- `rowKey` 必填
- 操作列固定右侧（放最后）
- 长文本做省略或 tooltip
- 状态字段尽量可视化（tag / switch / 明确文案）
- 时间字段统一格式
- 危险操作必须二次确认（`modal.confirm()`）
- 分页参数统一维护
- 列标题明确，不做无意义复杂渲染

### 查询区

- 查询条件优先展示高频项
- 低频筛选可以折叠或后置
- 查询表单和编辑表单分开管理
- 查询区优先使用 `flex` + `flex-wrap`

### 请求与数据流

- 查询、详情、新增、编辑、删除方法分开
- 请求方法命名按业务语义来
- 加载态 `tableLoading` 明确
- 提交中 `submitLoading` 明确
- 成功后给反馈，并刷新列表或关闭弹窗

---

## 状态命名建议

优先使用：

- `tableData`
- `queryForm`
- `pagination`
- `tableLoading`
- `submitLoading`

### 方法命名建议

优先使用：

- `getList`
- `handleSearch`
- `handleReset`
- `handleCreate`
- `handleEdit`
- `handleDelete`
- `handleView`

避免：

- `doIt` / `submitFn` / `clickOk` / `getDataInfo` / `handleEverything`

---

## 逻辑分层要求

页面逻辑至少分成：

1. 页面状态
2. 数据请求
3. 表格操作
4. 表单操作
5. 生命周期初始化

不要把逻辑写成一坨。

---

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
- Ant Design Vue 4.x

---

## Output Format

当用户要求开发列表页时：

1. **先给结论** — 功能、结构、是否拆组件、关键交互
2. **再给文件结构** — 代码放在哪
3. **再给完整代码** — Vue 页面、接口占位、关键方法
4. **最后补充接入说明** — 待替换字段、权限对接、可抽公共组件

---

## 禁止事项

- 不要把内容堆在一个 Card 里
- 不要凭空捏造后端接口字段
- 不要生成花哨、视觉很满但信息混乱的页面
- 不要把所有逻辑揉进模板
