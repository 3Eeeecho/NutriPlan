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
        <h3 class="recipe-name">{{ mealTitle }}</h3>
      </div>
      
      <!-- Tags / Meta -->
      <div v-if="hasRecipes" class="meal-meta">
        <div class="meta-tag">🍽️ {{ mealRecipes.length }} 道</div>
        <div class="meta-tag">🔥 {{ Math.floor(totalEnergy) }} kcal</div>
        <div class="meta-tag">💪 {{ Math.floor(totalProtein) }}g 蛋白</div>
      </div>

      <div v-if="mealRecipes.length > 1" class="meal-combo-list">
        <span v-for="item in mealRecipes" :key="item.id || item.name" class="combo-chip">{{ item.name }}</span>
      </div>
      
      <!-- Fallback text if no recipe -->
      <div v-else class="no-recipe-text">
        点击添加食谱
      </div>
    </div>

    <!-- Action -->
    <div class="meal-action">
      <n-button 
        v-if="hasRecipes"
        circle 
        secondary 
        :type="isSynced ? 'default' : 'success'" 
        class="sync-btn" 
        :class="{ 'synced': isSynced }"
        @click.stop="handleSync"
        :title="isSynced ? '已打卡' : '打卡并同步到饮食记录'"
        :disabled="isSynced"
      >
        <template #icon>
          <n-icon>
            <CheckmarkCircle v-if="isSynced" />
            <CheckmarkCircleOutline v-else />
          </n-icon>
        </template>
      </n-button>
      <n-button text class="action-btn" @click.stop="viewDetail">
        <n-icon size="24" color="#d1d5db"><ChevronForward /></n-icon>
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { NButton, NIcon } from 'naive-ui';
import { ChevronForward, CheckmarkCircleOutline, CheckmarkCircle } from '@vicons/ionicons5';
import { useRouter } from 'vue-router';

const props = defineProps({
  title: String, // 早餐, 午餐...
  icon: String, // Optional, backward compatibility
  recipe: Object,
  recipes: {
    type: Array,
    default: () => []
  },
  isSynced: Boolean // New prop
});

const emit = defineEmits(['sync']);

const router = useRouter();

const mealRecipes = computed(() => {
  if (Array.isArray(props.recipes) && props.recipes.length > 0) {
    return props.recipes.filter(item => item && item.name);
  }
  return props.recipe ? [props.recipe] : [];
});

const hasRecipes = computed(() => mealRecipes.value.length > 0);

const mealTitle = computed(() => {
  if (!hasRecipes.value) return '未安排';
  if (mealRecipes.value.length === 1) return mealRecipes.value[0]?.name || '未安排';
  return mealRecipes.value.map(item => item.name).join(' + ');
});

const totalEnergy = computed(() => mealRecipes.value.reduce((sum, item) => sum + (Number(item.energy) || 0), 0));
const totalProtein = computed(() => mealRecipes.value.reduce((sum, item) => sum + (Number(item.protein) || 0), 0));

const viewDetail = () => {
  const target = mealRecipes.value[0];
  if (target?.id) {
    router.push(`/recipes/${target.id}`);
  }
};

const handleSync = () => {
  if (props.isSynced) return;
  const target = mealRecipes.value[0];
  console.log('MealItem handleSync clicked', target);
  if (target) {
    emit('sync', { recipe: target, type: props.title });
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

.meal-combo-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
}

.combo-chip {
  font-size: 12px;
  color: #374151;
  background-color: #eef2ff;
  padding: 3px 8px;
  border-radius: 999px;
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
  display: flex;
  align-items: center;
  gap: 8px;
  opacity: 1; /* Always visible for better UX, or keep hover effect */
}

.sync-btn {
  transition: all 0.2s;
}

.sync-btn:hover {
  transform: scale(1.1);
}
</style>
