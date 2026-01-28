<template>
  <div class="recipe-card" @click="viewDetail">
    <div class="card-image-placeholder" :style="{ background: getGradient(icon) }">
      <span class="meal-icon">{{ icon }}</span>
    </div>
    
    <div class="card-content">
      <div class="header">
        <h4 class="recipe-name">{{ recipe.name }}</h4>
        <n-tag size="small" :bordered="false" type="success" round class="meal-tag">
          {{ title }}
        </n-tag>
      </div>
      
      <p v-if="recipe.description" class="recipe-desc text-ellipsis">{{ recipe.description }}</p>
      
      <div class="nutrition-grid">
        <div class="nutrient-item">
          <div class="value-row">
            <span class="value energy">{{ Math.floor(recipe.energy) }}</span>
            <span class="unit">kcal</span>
          </div>
          <span class="label">热量</span>
        </div>
        <div class="nutrient-item">
          <div class="value-row">
            <span class="value protein">{{ Math.floor(recipe.protein) }}</span>
            <span class="unit">g</span>
          </div>
          <span class="label">蛋白质</span>
        </div>
        <div class="nutrient-item">
          <div class="value-row">
            <span class="value carb">{{ Math.floor(recipe.carbohydrate) }}</span>
            <span class="unit">g</span>
          </div>
          <span class="label">碳水</span>
        </div>
        <div class="nutrient-item">
          <div class="value-row">
            <span class="value fat">{{ Math.floor(recipe.fat) }}</span>
            <span class="unit">g</span>
          </div>
          <span class="label">脂肪</span>
        </div>
      </div>

      <div class="footer-info">
        <div class="meta">
          <span class="time" v-if="recipe.cooking_time">
            <n-icon><TimeOutline /></n-icon> {{ recipe.cooking_time }}m
          </span>
          <span class="difficulty" v-if="recipe.difficulty">
            {{ getDifficultyEmoji(recipe.difficulty) }}
          </span>
        </div>
        
        <n-button 
          circle 
          secondary 
          :type="localFavorite ? 'error' : 'default'" 
          @click.stop="toggleFavorite"
          class="fav-btn"
        >
          <template #icon>
            <n-icon><Heart v-if="localFavorite" /><HeartOutline v-else /></n-icon>
          </template>
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { NTag, NButton, NIcon, useMessage } from 'naive-ui'
import { TimeOutline, Heart, HeartOutline } from '@vicons/ionicons5'
import { addFavorite, removeFavorite } from '@/api/recipeApi'

const props = defineProps({
  title: {
    type: String,
    required: true
  },
  recipe: {
    type: Object,
    required: true
  },
  icon: {
    type: String,
    default: '🍽️'
  }
})

const router = useRouter()
const message = useMessage()
const localFavorite = ref(props.recipe.is_favorite || false)

function getGradient(icon) {
  // 根据图标或类型生成柔和的背景渐变
  return 'linear-gradient(135deg, #e0f2fe 0%, #f0fdf4 100%)'
}

function getDifficultyEmoji(difficulty) {
  const map = {
    '简单': '🟢',
    '中等': '🟡',
    '困难': '🔴'
  };
  return map[difficulty] || '⚪';
}

function viewDetail() {
  router.push({ name: 'RecipeDetail', params: { id: props.recipe.id } })
}

async function toggleFavorite(e) {
  try {
    if (localFavorite.value) {
      await removeFavorite(props.recipe.id)
      localFavorite.value = false
      message.success('已取消收藏')
    } else {
      await addFavorite(props.recipe.id)
      localFavorite.value = true
      message.success('已加入收藏')
    }
  } catch (err) {
    console.error(err)
    message.error('操作失败')
  }
}
</script>

<style scoped>
.recipe-card {
  background: white;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid var(--color-border-primary, #e5e7eb);
  transition: all 0.3s ease;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.recipe-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 24px -8px rgba(0, 0, 0, 0.1);
  border-color: var(--color-primary, #10b981);
}

.card-image-placeholder {
  height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  position: relative;
}

.card-content {
  padding: 16px;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 8px;
  gap: 8px;
}

.recipe-name {
  font-size: 16px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.recipe-desc {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 16px;
  line-height: 1.5;
  height: 40px; /* 固定高度保持卡片整齐 */
}

.text-ellipsis {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 4px;
  margin-bottom: 16px;
  padding: 8px;
  background: #f9fafb;
  border-radius: 8px;
}

.nutrient-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.value-row {
  display: flex;
  align-items: baseline;
  gap: 2px;
}

.nutrient-item .value {
  font-size: 14px;
  font-weight: 700;
}

.nutrient-item .unit {
  font-size: 10px;
  color: #9ca3af;
}

.nutrient-item .label {
  font-size: 11px;
  color: #6b7280;
}

.value.energy { color: #10b981; }
.value.protein { color: #ef4444; }
.value.carb { color: #f59e0b; }
.value.fat { color: #8b5cf6; }

.footer-info {
  margin-top: auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #f3f4f6;
}

.meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: #6b7280;
  align-items: center;
}

.meta .time {
  display: flex;
  align-items: center;
  gap: 4px;
}
</style>
