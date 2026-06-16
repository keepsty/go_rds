---
name: ant-design-vue-db-ops-pages
description: >
  Specialized patterns for database, ops, approval-flow, and infrastructure admin pages
  with Vue 3 + Ant Design Vue 4.x. Use when building SQL audit pages, order/workflow pages,
  DB instance management, monitor/alert pages, or permission management UIs.
  This is part of the ant-design-vue skill family — for generic CRUD/table patterns,
  see ant-design-vue-crud-page.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# Database & Ops Admin Pages（数据库运维页面专项）

你专注于**数据库、工单、审批、配置、运维类**后台页面的开发。
这些页面有更高的信息密度和安全敏感度要求。

---

## When to Apply

- 开发数据库工单页
- 开发审批流相关页面
- 开发 SQL 审核 / 结果展示页
- 开发资源申请 / 扩容页
- 开发慢查询分析页
- 开发监控告警页
- 开发变更记录页
- 开发环境配置 / 权限管理页
- 用户提到"数据库"、"工单"、"审批"、"SQL"、"运维"等关键词

## 不涵盖的内容

- 通用 CRUD 列表页 → 参考 [[ant-design-vue-crud-page]]
- 表单弹窗 → 参考 [[ant-design-vue-form-modal]]
- 响应式布局 → 参考 [[ant-design-vue-responsive-layout]]

---

## 核心原则

1. 信息优先级明确
2. 结果展示比装饰更重要
3. 查询效率比视觉炫技更重要
4. 表格字段要真实可读
5. 操作风险要清晰提示
6. 结果页、详情页、审核页要重点突出关键信息

---

## 工单类页面

必须突出：

- 工单状态
- 执行环境
- 风险等级
- 审批链路

---

## SQL 结果展示页

- 重视可读性
- 避免字段堆叠混乱
- SQL 文本区域需要有合适的展示方式（等高宽、语法高亮等）

---

## 审批流页面

必须突出：

- 节点关系
- 当前节点
- 审批结论
- 各节点审批人/审批时间

---

## 风险与安全

- 环境必须明确区分：测试 / 预发 / 生产
- 高风险操作：按钮颜色、确认文案都要更直接
- 危险操作必须二次确认
- 高危按钮优先使用 `danger` + `modal.confirm()` 双重确认

---

## 典型场景清单

- 数据库工单页
- 审批流配置页
- SQL 审核页
- 资源申请页
- 扩容页
- 慢查询分析页
- 监控告警页
- 变更记录页
- 数据查询页
- 环境配置页
- 权限管理页

---

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
- Ant Design Vue 4.x
