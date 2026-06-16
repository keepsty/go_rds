---
name: ant-design-vue-design-tokens
description: >
  Visual design token rules for Vue 3 + Ant Design Vue 4.x admin pages.
  Use when choosing colors, spacing, border-radius, or button styles.
  Enforces unified enterprise admin styling — no handmade color/spacing values.
  This is part of the ant-design-vue skill family — for layout/response,
  see ant-design-vue-responsive-layout.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# Design Tokens（企业后台设计规范）

你专注 **Vue 3 + Ant Design Vue 4.x** 的视觉设计规范。
目标是确保所有页面视觉统一，不出现散点样式。

---

## When to Apply

- 选择颜色、间距、圆角时
- 确定按钮样式时
- 统一页面视觉风格时
- 用户提到"样式"、"颜色"、"间距"、"圆角"、"设计规范"等关键词

## 不涵盖的内容

- 响应式布局 → 参考 [[ant-design-vue-responsive-layout]]
- 表单/弹窗 → 参考 [[ant-design-vue-form-modal]]

---

## 设计定位

这是一个**企业后台 UI 风格**，不是官网、营销页、作品集。

默认风格：

- 简洁、稳、清晰、有层级、易读、可维护、适度精致

要求：

- 页面信息分区明确
- 视觉层级清晰
- 表单和表格不拥挤
- 状态和操作一眼能懂
- 不堆无意义渐变、阴影、动画
- 不把营销站风格硬套到后台系统

---

## 间距规范

仅使用：

- `4 / 8 / 12 / 16 / 24 / 32`

禁止：

- 10、14、18、22 等散点间距

---

## 圆角规范

仅使用：

- `6 / 8 / 10 / 12`

禁止：

- 组件中手写临时圆角值

---

## 颜色规范

优先使用 Ant Design 语义色，不随意造色：

| 语义 | 颜色值 |
|---|---|
| Primary | `#1677ff` |
| Success | `#52c41a` |
| Warning | `#faad14` |
| Error | `#ff4d4f` |
| Border | `#d9d9d9` |
| 主文本 | `rgba(0,0,0,0.88)` |
| 次文本 | `rgba(0,0,0,0.65)` |

---

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
- Ant Design Vue 4.x
