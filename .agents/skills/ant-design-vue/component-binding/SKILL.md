
# Component Binding（父子组件数据绑定）

你专注于 **Vue 3 + JavaScript** 的父子组件数据绑定模式。

---

## When to Apply

- 需要拆页面为多个组件时
- 设计 Props / Emits 接口时
- 处理弹窗 open/close 控制时
- 处理查询组件与表格数据传递时

---

## 总原则

- 父组件：数据源、状态控制、接口调用
- 子组件：展示、局部交互、事件抛出
- 坚持**单向数据流**，子组件不直接改父组件状态

---

## Props / Emit / v-model

| 模式 | 说明 |
|------|------|
| `defineProps` | 子组件接收数据，**不允许**直接修改 |
| `defineEmits` | 子组件通知父组件：`submit` / `cancel` / `change` / `update:open` |
| `v-model` | 单个双向值用 `v-model`，多个用 `v-model:xxx` |
| 弹窗开关 | 父管 `v-model:open`，子收 `open` + 派发 `update:open` |

---

## 典型模式

### 弹窗表单
- **父组件**：管 `open`、`currentRow`、提交接口、刷新列表
- **子组件**：收 `open`、`formData`，抛 `submit`/`cancel`/`update:open`，负责展示和校验
- 子组件表单**不要自己发请求改父组件列表**

### 表格 + 查询
- **父组件**：管 `queryForm`、`tableData`、分页和请求
- **子组件**：查询组件只抛条件，表格组件只抛分页/选择/行操作事件

---

## 技术约束

- Vue 3 + `script setup`，JavaScript
