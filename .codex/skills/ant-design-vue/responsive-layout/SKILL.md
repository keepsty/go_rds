---
name: ant-design-vue-responsive-layout
description: >
  Responsive and adaptive layout rules for Vue 3 + Ant Design Vue 4.x admin pages.
  Use when building page layouts, handling screen width compatibility,
  designing flex/grid arrangements, or adding media queries.
  This is part of the ant-design-vue skill family — for design tokens/colors,
  see the sibling skill ant-design-vue-design-tokens.
license: MIT
metadata:
  author: custom
  version: "1.0.0"
---

# Responsive & Adaptive Layout（响应式自适应布局）

你专注于 **Vue 3 + Ant Design Vue 4.x** 后台页面的响应式与自适应布局。
目标是确保页面在不同办公屏幕宽度下都可正常工作。

---

## When to Apply

- 新建页面布局结构时
- 需要兼容不同屏幕宽度时
- 优化现有页面的自适应表现时
- 用户提到"响应式"、"自适应"、"布局"、"flex"、"grid"、"屏幕适配"等关键词

## 不涵盖的内容

- 颜色/间距/圆角 Token → 参考 [[ant-design-vue-design-tokens]]
- 表格/表单页面结构 → 参考 [[ant-design-vue-crud-page]]

---

## 兼容目标

以桌面端后台系统为主，兼容常见办公屏幕宽度：

- 默认兼容：1440、1366、1280、1024
- 页面在窄屏下允许折行、换列、收起，**不允许主要操作消失**
- 尽量避免整体页面横向滚动；若表格列确实过多，只允许表格区域内部横向滚动

---

## 页面布局

- 使用 24 栅格：`a-row` + `a-col`
- `gutter` 优先 `16` 或 `24`
- 页面主区域优先使用 `flex`、`grid`、百分比宽度、`minmax`
- 避免大量固定像素宽度写死布局
- 卡片区、统计区、筛选区支持自动换行
- 表单 label 宽度保持统一，窄屏下允许改为纵向布局

---

## 查询区

- 查询区优先使用 `flex` + `flex-wrap`
- 不把所有查询项写死成一行
- 宽屏下多列排列，窄屏下自动折为 2 列或 1 列
- 查询按钮区固定在查询区末尾或单独成行

---

## 表格区

- 表格外层容器允许 `overflow-x: auto`
- 高优先级列固定展示，低优先级列在窄屏下可适度缩短宽度
- 操作列保持可见，不要被挤到不可点
- 长文本列使用省略、tooltip 或折叠展示

---

## 弹窗 / 抽屉

- Modal 宽度使用区间控制，不写死超大宽度
- 大表单优先 Drawer，避免窄屏弹窗塞满
- 弹窗内容区高度超出时内部滚动，不让整页抖动

---

## 按钮区

- 按钮区在窄屏下自动换行
- 批量操作、筛选、导入导出按钮过多时允许折叠到更多菜单

---

## 样式实现建议

- 优先使用 `flex` / `grid`
- 使用 `clamp()`、百分比、`min-width`、`max-width` 控制尺寸
- 必要时补充媒体查询，不要堆很多断点
- 断点建议围绕：`1200px`、`992px`、`768px`

---

## 技术约束

- Vue 3 + `script setup`
- JavaScript（不使用 TypeScript）
- Ant Design Vue 4.x
