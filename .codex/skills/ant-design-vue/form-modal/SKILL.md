---
name: ant-design-vue-form-modal
description: >
  Build production-grade forms, modals, and drawers with Vue 3 + JavaScript + Ant Design Vue 4.x.
  Use when creating form components, modal forms, drawer forms, edit dialogs,
  search filters, or form validation logic.
  This is part of the ant-design-vue skill family — for table pages, component binding,
  or responsive layout, see the sibling skills.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# Form & Modal（企业后台表单与弹窗开发）

你是一个专注于 **Vue 3 + JavaScript + Ant Design Vue 4.x** 表单与弹窗开发专家。
目标是输出交互自然、校验完整、可复用的表单弹窗代码。

---

## When to Apply

- 新建表单页面 / 组件
- 开发新增/编辑弹窗（Modal）
- 开发详情抽屉（Drawer）
- 表单校验逻辑设计
- 用户提到"表单"、"弹窗"、"抽屉"、"Modal"、"Drawer"、"校验"等关键词

## 不涵盖的内容

- 表格列表页 → 参考 [[ant-design-vue-crud-page]]
- 父子组件数据传递 → 参考 [[ant-design-vue-component-binding]]
- 响应式适配 → 参考 [[ant-design-vue-responsive-layout]]

---

## Best Practices

### 表单 (Form)

- 优先使用 `a-form`
- 表单字段命名与接口字段保持一致
- 校验规则集中维护
- 提交前统一校验：走 `validateFields`
- 编辑态必须正确回填（回填前先整理字段映射）
- 必填项明确，自动补齐规则
- 不随意调用 `resetFields` 做粗暴清空
- 查询表单和编辑表单分开管理

### 弹窗 / 抽屉

- 简单新增/编辑优先使用 **Modal**
- 信息量大、需要看上下文时用 **Drawer**
- 能共用一个弹窗就不要拆多个（新增/编辑优先共用）
- 使用 `destroyOnClose`
- footer 只保留核心按钮
- 标题明确区分 新增 / 编辑 / 查看
- 弹窗状态和表单状态分离
- 提交动作统一在确定按钮中处理
- Modal 宽度使用区间控制，不写死超大宽度
- 大表单优先 Drawer，避免窄屏弹窗塞满
- 弹窗内容区高度超出时内部滚动

---

## 按钮规范

- 主操作：`a-button type="primary"`
- 次操作：默认 `a-button`
- 危险操作：`a-button danger`
- 同一视觉区域最多两个主按钮
- 文案要直白，不要写虚词
- 默认统一 `size="middle"`
- 禁止手动设置按钮颜色、圆角、高度

---

## 反馈规范

- 轻提示：`message.success()`
- 显著错误：`notification.error()`
- 不可逆动作：`modal.confirm()`
- 加载反馈：`a-spin`
- 提交中 `submitLoading` 状态明确
- 成功后给反馈，并刷新列表或关闭弹窗

---

## 状态命名建议

优先使用：

- `modalOpen`
- `drawerOpen`
- `formState`
- `currentRow`
- `submitLoading`
- `queryForm`

### 方法命名建议

优先使用：

- `openModal`
- `closeModal`
- `handleSubmit`
- `handleCreate`
- `handleEdit`

避免：

- `doIt` / `submitFn` / `clickOk`

---

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
- Ant Design Vue 4.x
