# NutriPlan 设计系统 - 阶段一完成

## ✅ 已完成内容

### 1. Naive UI 集成
- ✅ 安装 Naive UI 及相关依赖
- ✅ 配置全局主题覆盖
- ✅ 设置中文语言环境

### 2. 设计系统文件结构

```
src/
  styles/
    ├── variables.css    # CSS 变量（颜色、字体、间距等）
    ├── tokens.js        # JavaScript 设计 token
    ├── theme.js         # Naive UI 主题配置
    └── mixins.css       # 可复用样式混入
  
  components/
    ├── layout/          # 布局组件
    │   ├── PageLayout.vue
    │   ├── TopNavigation.vue
    │   └── BackButton.vue
    └── ui/              # 基础 UI 组件
        └── StatCard.vue
```

### 3. 设计 Token

#### 颜色系统
```css
/* 品牌色 */
--color-primary: #10b981 (绿色主色调)

/* 功能色 */
--color-success: #10b981
--color-warning: #f59e0b
--color-error: #ef4444
--color-info: #3b82f6

/* 营养素色彩 */
--color-protein: #ef4444  (蛋白质 - 红)
--color-carb: #f59e0b     (碳水 - 橙)
--color-fat: #8b5cf6      (脂肪 - 紫)
```

#### 间距系统
```css
--spacing-xs: 4px
--spacing-sm: 8px
--spacing-md: 12px
--spacing-lg: 16px
--spacing-xl: 24px
--spacing-2xl: 32px
```

#### 圆角系统
```css
--radius-sm: 6px
--radius-md: 8px
--radius-lg: 12px
--radius-xl: 16px
--radius-2xl: 20px
```

## 📦 创建的组件

### 1. 布局组件

#### PageLayout.vue
基础页面容器，提供统一的页面结构。

```vue
<template>
  <PageLayout>
    <TopNavigation />
    <!-- 页面内容 -->
  </PageLayout>
</template>
```

#### TopNavigation.vue
顶部导航栏组件，包含 Logo 和用户下拉菜单。

特性：
- 毛玻璃效果
- 粘性定位
- 响应式设计
- 统一的下拉菜单

使用示例：
```vue
<TopNavigation />
```

#### BackButton.vue
返回按钮组件

```vue
<BackButton text="返回主页" />
<BackButton :to="/home" />
```

### 2. UI 组件

#### StatCard.vue
统计卡片组件，用于展示关键数据。

```vue
<StatCard
  icon="🔥"
  iconBg="linear-gradient(135deg, #10b981 0%, #059669 100%)"
  :value="2028"
  label="每日目标热量"
  subtitle="kcal"
  clickable
  @click="handleClick"
/>
```

## 🎨 如何使用设计系统

### 在组件中使用 CSS 变量

```vue
<style scoped>
.my-card {
  background: var(--color-bg-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-xl);
  color: var(--color-text-primary);
  box-shadow: var(--shadow-md);
}
</style>
```

### 使用预定义的混入类

```vue
<div class="flex-between card-base card-hover">
  <!-- 内容 -->
</div>
```

### 在 JS 中使用设计 Token

```javascript
import { colors, spacing, radius } from '@/styles/tokens'

const cardStyle = {
  background: colors.primary.default,
  padding: spacing.xl,
  borderRadius: radius.lg,
}
```

### 使用 Naive UI 组件

```vue
<script setup>
import { NButton, NCard, NInput } from 'naive-ui'
</script>

<template>
  <n-button type="primary">主要按钮</n-button>
  <n-card title="卡片标题">卡片内容</n-card>
  <n-input placeholder="请输入" />
</template>
```

## 🔄 迁移指南

### Element Plus → Naive UI

| Element Plus | Naive UI | 备注 |
|-------------|----------|------|
| `el-button` | `n-button` | 属性基本一致 |
| `el-card` | `n-card` | 使用 `title` 属性设置标题 |
| `el-input` | `n-input` | 使用 `placeholder` 属性 |
| `el-select` | `n-select` | 使用 `options` 属性 |
| `el-dropdown` | `n-dropdown` | 使用 `options` 数组 |
| `ElMessage` | `useMessage()` | 需要在 setup 中使用 |

### 示例：按钮迁移

**之前（Element Plus）：**
```vue
<el-button type="primary" size="large">
  确认
</el-button>
```

**之后（Naive UI）：**
```vue
<n-button type="primary" size="large">
  确认
</n-button>
```

### 示例：消息提示迁移

**之前（Element Plus）：**
```javascript
import { ElMessage } from 'element-plus'

ElMessage.success('操作成功')
```

**之后（Naive UI）：**
```javascript
import { useMessage } from 'naive-ui'

const message = useMessage()
message.success('操作成功')
```

## 📱 响应式断点

```javascript
xs: 480px   // 手机
sm: 640px   // 大屏手机
md: 768px   // 平板
lg: 1024px  // 笔记本
xl: 1280px  // 桌面
2xl: 1536px // 大屏
```

使用示例：
```vue
<style scoped>
.container {
  padding: var(--spacing-md);
}

@media (min-width: 768px) {
  .container {
    padding: var(--spacing-xl);
  }
}
</style>
```

## 🎯 下一步计划（阶段二）

1. **迁移主要页面**
   - Home.vue 使用新的布局和组件
   - RecipeRecommend.vue 优化
   - IntakeRecord.vue 重构

2. **创建更多业务组件**
   - NutritionProgressBar.vue
   - MealCard.vue
   - RecipePlanCard.vue

3. **优化交互体验**
   - 添加加载动画
   - 优化过渡效果
   - 完善反馈机制

## 📝 注意事项

1. **渐进式迁移**：不要一次性替换所有 Element Plus 组件，逐步迁移
2. **保持一致性**：使用设计系统中定义的变量和组件
3. **测试兼容性**：确保在不同浏览器和设备上测试
4. **性能优化**：按需引入 Naive UI 组件

## 🚀 启动项目

```bash
# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build
```

---

**设计系统版本**: v1.0.0  
**更新日期**: 2026-01-24  
**负责人**: NutriPlan Team
