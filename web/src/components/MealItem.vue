<template>
  <div class="meal-item" @click="viewDetail">
    <!-- Icon -->
    <div class="meal-icon-wrapper">
      <span class="meal-icon">{{ getIcon(title) }}</span>
    </div>

    <!-- Content -->
    <div class="meal-content">
      <div class="meal-header">
        <span class="meal-type-badge">{{ title }}</span>
        <h3 class="recipe-name">{{ recipe?.name || '未安排' }}</h3>
      </div>
      
      <!-- Tags / Meta -->
      <div v-if="recipe" class="meal-meta">
        <div class="meta-tag">🔥 {{ Math.floor(recipe.energy) }} kcal</div>
        <div class="meta-tag">⏱️ {{ recipe.cooking_time }}分钟</div>
        <div class="meta-tag">💪 {{ Math.floor(recipe.protein) }}g 蛋白</div>
      </div>
      
      <!-- Fallback text if no recipe -->
      <div v-else class="no-recipe-text">
        点击添加食谱
      </div>
    </div>

    <!-- Action -->
    <div class="meal-action">
      <n-button text class="action-btn" @click.stop="viewDetail">
        <n-icon size="24" color="#d1d5db"><ChevronForward /></n-icon>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { NButton, NIcon } from 'naive-ui';
import { ChevronForward } from '@vicons/ionicons5';
import { useRouter } from 'vue-router';

const props = defineProps({
  title: String, // 早餐, 午餐...
  icon: String, // Optional, backward compatibility
  recipe: Object
});

const router = useRouter();

const viewDetail = () => {
  if (props.recipe?.id) {
    router.push(`/recipes/${props.recipe.id}`);
  }
};

const getIcon = (title) => {
  if (props.icon) return props.icon;
  const icons = {
    '早餐': '🌅',
    '午餐': '☀️',
    '晚餐': '🌙',
    '加餐': '🍎'
  };
  return icons[title] || '🍽️';
};
</script>

<style scoped>
.meal-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background-color: #fff;
  border-radius: 12px;
  transition: all 0.2s ease;
  cursor: pointer;
  border: 1px solid transparent;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.meal-item:hover {
  background-color: #f9fafb;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.meal-icon-wrapper {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f3f4f6;
  border-radius: 12px;
  margin-right: 16px;
  font-size: 24px;
  flex-shrink: 0;
}

.meal-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
}

.meal-header {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.meal-type-badge {
  font-size: 12px;
  color: #6b7280;
  font-weight: 500;
  letter-spacing: 0.05em;
}

.recipe-name {
  font-size: 16px;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
}

.meal-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
}

.meta-tag {
  font-size: 12px;
  color: #6b7280;
  background-color: #f3f4f6;
  padding: 2px 8px;
  border-radius: 4px;
  display: inline-flex;
  align-items: center;
}

.no-recipe-text {
  color: #9ca3af;
  font-style: italic;
  font-size: 14px;
}

.meal-action {
  margin-left: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.meal-item:hover .meal-action {
  opacity: 1;
}
</style>
