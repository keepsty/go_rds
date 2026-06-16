---
name: ant-design-vue-component-binding
description: >
  Parent-child component data binding patterns for Vue 3 + Ant Design Vue 4.x.
  Use when designing component interfaces, deciding props/emits/v-model,
  splitting a page into components, or refactoring data flow between components.
  This is part of the ant-design-vue skill family — for forms/modals, table pages,
  or responsive layout, see the sibling skills.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# Component Binding（父子组件数据绑定规则）

你专注于 **Vue 3 + JavaScript** 的父子组件数据绑定模式。
目标是确保组件间数据流清晰、单向、可维护。

---

## When to Apply

- 需要拆页面为多个组件时
- 设计组件 Props / Emits 接口时
- 处理弹窗表单的 open/close 控制时
- 处理查询组件与表格的数据传递时
- 用户提到"组件通信"、"v-model"、"props"、"emit"等关键词

## 不涵盖的内容

- 表单弹窗本身 → 参考 [[ant-design-vue-form-modal]]
- 表格列表页 → 参考 [[ant-design-vue-crud-page]]

---

## 总原则

- 父组件负责：数据源、状态控制、接口调用
- 子组件负责：展示、局部交互、事件抛出
- 坚持**单向数据流**，避免子组件直接改父组件状态

---

## Props

- 子组件通过 `defineProps` 接收数据
- **不允许**直接修改 `props`
- 传入对象或数组时，子组件也不要直接改原值

## Emit

- 子组件通过 `defineEmits` 通知父组件
- 事件命名要直白：
  - `submit`
  - `cancel`
  - `change`
  - `update:open`
  - `update:value`

## v-model

- 组件只有一个核心双向绑定值时，优先使用 `v-model`
- 有多个双向字段时，使用 `v-model:xxx`
- 弹窗开关统一推荐：
  - 父组件：`v-model:open`
  - 子组件：接收 `open`，派发 `update:open`

---

## 典型模式

### 弹窗表单

- 父组件：
  - 管 `open`
  - 管 `currentRow`
  - 管提交接口
  - 管刷新列表
- 子组件：
  - 接收 `open`、`formData`
  - 抛出 `submit`、`cancel`、`update:open`
  - 负责展示、编辑、校验
  - **子组件表单不要自己偷偷发请求并改父组件列表**

### 表格 + 查询

- 父组件：
  - 管 `queryForm`
  - 管 `tableData`
  - 管分页和请求
- 子组件：
  - 查询组件只抛查询条件，不直接控制表格刷新
  - 表格组件只抛分页、选择、行操作事件

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
