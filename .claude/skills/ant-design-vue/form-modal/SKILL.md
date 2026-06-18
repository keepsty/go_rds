
# Form & Modal（企业后台表单与弹窗）

你是一个专注于 **Vue 3 + JavaScript + Ant Design Vue 4.x** 表单与弹窗开发专家。

---

## When to Apply

- 新建表单页面 / 组件
- 新增/编辑弹窗（Modal）
- 详情抽屉（Drawer）
- 表单校验逻辑

---

## Best Practices

### 表单
- 字段命名与接口字段保持一致
- 校验规则集中维护，提交前统一 `validateFields`
- 编辑态必须回填，不随意 `resetFields`
- 查询表单和编辑表单分开管理

### 弹窗 / 抽屉
- 简单新增/编辑用 **Modal**，信息量大用 **Drawer**
- 新增/编辑优先共用同一个弹窗
- 使用 `destroyOnClose`，标题明确区分 新增/编辑/查看
- Modal 宽度区间控制，大表单优先 Drawer

### 按钮规范
- 主操作 `type="primary"`，危险操作 `danger`
- 同一区域最多两个主按钮，统一 `size="middle"`
- 禁止手写按钮颜色、圆角、高度

### 反馈规范
- 轻提示 `message.success()`，显著错误 `notification.error()`
- 不可逆动作 `modal.confirm()`
- 加载反馈 `a-spin`，提交中 `submitLoading`

### 状态命名
- `modalOpen` / `drawerOpen` / `formState` / `currentRow` / `submitLoading`
- 方法：`openModal` / `closeModal` / `handleSubmit`

---

## 技术约束

- Vue 3 + `script setup`，JavaScript，Ant Design Vue 4.x
