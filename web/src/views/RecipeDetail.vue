<template>
  <div class="page-container">
    <!-- Atmospheric Background Blobs -->
    <div class="blob-container">
      <div class="blob blob-1"></div>
      <div class="blob blob-2"></div>
    </div>

    <TopNavigation />

    <div class="content-wrapper">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <n-spin size="large" />
      </div>

      <!-- Error State -->
      <div v-else-if="!recipe || !recipe.id" class="error-state">
        <n-result status="404" title="未找到食谱" description="该食谱可能已被删除或不存在">
          <template #footer>
            <n-button @click="router.back()">返回上一页</n-button>
          </template>
        </n-result>
      </div>

      <!-- Main Content -->
      <div v-else class="animate-fade-in">
        <!-- Back Button -->
        <div class="back-btn-container">
          <n-button text class="back-btn" @click="router.back()">
            <template #icon>
              <n-icon><ArrowBack /></n-icon>
            </template>
            返回
          </n-button>
        </div>

        <!-- Header Section -->
        <div class="header-section">
          <div class="header-content">
            <h1 class="recipe-title">{{ recipe.name }}</h1>
            
            <!-- Meta Tags -->
            <div class="meta-tags">
              <div class="meta-tag">
                <n-icon class="icon-emerald"><TimeOutline /></n-icon>
                <span>{{ recipe.cooking_time || 15 }} 分钟</span>
              </div>
              <div class="meta-tag">
                <n-icon class="icon-orange"><FlameOutline /></n-icon>
                <span>{{ recipe.difficulty || '简单' }}</span>
              </div>
              <div class="meta-tag" v-if="recipe.meal_type">
                <n-icon class="icon-blue"><RestaurantOutline /></n-icon>
                <span>{{ recipe.meal_type }}</span>
              </div>
            </div>
          </div>

          <!-- Actions -->
          <div class="action-buttons">
            <n-button circle secondary type="default" size="large" @click="toggleFavorite">
              <template #icon>
                <n-icon :color="isFavorite ? '#ef4444' : undefined">
                  <Heart v-if="isFavorite" />
                  <HeartOutline v-else />
                </n-icon>
              </template>
            </n-button>
            <n-button circle secondary type="default" size="large" @click="shareRecipe">
              <template #icon>
                <n-icon><ShareSocialOutline /></n-icon>
              </template>
            </n-button>
          </div>
        </div>

        <!-- Hero Visual Section -->
        <div class="hero-visual" :style="{ background: getHeroBackground(recipe.name) }">
          <span class="hero-emoji">{{ getHeroEmoji(recipe.name) }}</span>
        </div>

        <!-- Nutrition Grid -->
        <div class="nutrition-grid">
          <!-- Calories -->
          <div class="nutrition-card">
            <div class="nutri-value-row">
              <span class="nutri-value text-emerald">{{ Math.floor(recipe.energy) }}</span>
              <span class="nutri-unit">kcal</span>
            </div>
            <div class="nutri-label">热量</div>
          </div>
          <!-- Protein -->
          <div class="nutrition-card">
            <div class="nutri-value-row">
              <span class="nutri-value text-red">{{ Math.floor(recipe.protein) }}</span>
              <span class="nutri-unit">g</span>
            </div>
            <div class="nutri-label">蛋白质</div>
          </div>
          <!-- Carbs -->
          <div class="nutrition-card">
            <div class="nutri-value-row">
              <span class="nutri-value text-amber">{{ Math.floor(recipe.carbohydrate) }}</span>
              <span class="nutri-unit">g</span>
            </div>
            <div class="nutri-label">碳水</div>
          </div>
          <!-- Fat -->
          <div class="nutrition-card">
            <div class="nutri-value-row">
              <span class="nutri-value text-purple">{{ Math.floor(recipe.fat) }}</span>
              <span class="nutri-unit">g</span>
            </div>
            <div class="nutri-label">脂肪</div>
          </div>
        </div>

        <div class="details-grid">
          <!-- Ingredients (Left Col) -->
          <div class="ingredients-col">
            <n-card :bordered="false" class="detail-card" title="所需食材">
              <template #header-extra>
                <n-icon size="20" class="icon-gray"><BasketOutline /></n-icon>
              </template>
              
              <div v-if="recipe.ingredients && recipe.ingredients.length" class="ingredients-list">
                <div 
                  v-for="(ing, idx) in recipe.ingredients" 
                  :key="idx"
                  class="ingredient-item"
                >
                  <div class="ing-name-wrapper">
                    <div class="ing-dot"></div>
                    <span class="ing-name">{{ parseIngredientName(ing) }}</span>
                  </div>
                  <span class="ing-amount">{{ parseIngredientAmount(ing) }}</span>
                </div>
              </div>
              <div v-else class="empty-text">暂无食材信息</div>
            </n-card>
          </div>

          <!-- Steps (Right 2 Cols) -->
          <div class="steps-col">
            <n-card :bordered="false" class="detail-card" title="烹饪步骤">
              <template #header-extra>
                <n-icon size="20" class="icon-gray"><ListOutline /></n-icon>
              </template>

              <div v-if="recipe.cooking_steps && recipe.cooking_steps.length" class="steps-list">
                <div 
                  v-for="(step, idx) in recipe.cooking_steps" 
                  :key="idx"
                  class="step-item group"
                >
                  <!-- Number Badge -->
                  <div class="step-badge-wrapper">
                    <div class="step-badge">
                      {{ idx + 1 }}
                    </div>
                  </div>
                  <!-- Step Text -->
                  <div class="step-text-wrapper">
                    <p class="step-text">
                      {{ step }}
                    </p>
                  </div>
                </div>
              </div>
              <div v-else class="empty-text">暂无烹饪步骤</div>
            </n-card>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { NButton, NIcon, NCard, NSpin, NResult, useMessage } from 'naive-ui';
import { 
  ArrowBack, 
  HeartOutline, 
  Heart, 
  ShareSocialOutline, 
  TimeOutline, 
  FlameOutline, 
  RestaurantOutline,
  BasketOutline,
  ListOutline
} from '@vicons/ionicons5';
import TopNavigation from '@/components/layout/TopNavigation.vue';
import { getRecipeDetail, addFavorite, removeFavorite } from '@/api/recipeApi';
import { useAuthStore } from '@/store/auth';

const route = useRoute();
const router = useRouter();
const message = useMessage();
const authStore = useAuthStore();

const id = route.params.id;
const recipe = ref(null);
const loading = ref(true);
const isFavorite = ref(false);

const load = async () => {
  loading.value = true;
  try {
    const res = await getRecipeDetail(id);
    recipe.value = res || {};
    isFavorite.value = !!recipe.value.is_favorite;
  } catch (e) {
    console.error(e);
    message.error('加载食谱详情失败');
  } finally {
    loading.value = false;
  }
};

const toggleFavorite = async () => {
  if (!authStore.isAuthenticated) {
    message.warning('请先登录');
    return;
  }
  try {
    if (isFavorite.value) {
      await removeFavorite(id);
      isFavorite.value = false;
      message.success('已取消收藏');
    } else {
      await addFavorite(id);
      isFavorite.value = true;
      message.success('已加入收藏');
    }
  } catch (e) {
    message.error('操作失败');
  }
};

const shareRecipe = async () => {
  const url = window.location.href;
  try {
    await navigator.clipboard.writeText(url);
    message.success('链接已复制到剪贴板');
  } catch (e) {
    message.info('请手动复制浏览器链接分享');
  }
};

const parseIngredientName = (ing) => {
  if (typeof ing === 'object' && ing !== null) return ing.name || '未知食材';
  // If string usually "Tomato" or "Tomato - 100g"
  if (typeof ing === 'string') {
    const parts = ing.split(/[-:]/);
    return parts[0].trim();
  }
  return String(ing);
};

const parseIngredientAmount = (ing) => {
  if (typeof ing === 'object' && ing !== null) return ing.amount || '';
  if (typeof ing === 'string') {
    const parts = ing.split(/[-:]/);
    return parts.length > 1 ? parts[1].trim() : '';
  }
  return '';
};

// Hero helpers
const getHeroEmoji = (name) => {
  if (!name) return '🥘';
  const map = {
    '鸡': '🍗', '牛': '🥩', '猪': '🍖', '鱼': '🐟', '虾': '🦐', 
    '蛋': '🥚', '菜': '🥬', '饭': '🍚', '面': '🍜', '汤': '🍲',
    '沙拉': '🥗', '面包': '🍞', '奶': '🥛', '果': '🍎'
  };
  for (const key in map) {
    if (name.includes(key)) return map[key];
  }
  return '🥘';
};

const getHeroBackground = (name) => {
  const colors = [
    'linear-gradient(135deg, #fef3c7 0%, #fffbeb 100%)', // amber-100
    'linear-gradient(135deg, #fee2e2 0%, #fef2f2 100%)', // red-100
    'linear-gradient(135deg, #d1fae5 0%, #ecfdf5 100%)', // emerald-100
    'linear-gradient(135deg, #dbeafe 0%, #eff6ff 100%)', // blue-100
    'linear-gradient(135deg, #f3e8ff 0%, #faf5ff 100%)'  // purple-100
  ];
  if (!name) return colors[0];
  const index = name.length % colors.length;
  return colors[index];
};

onMounted(() => {
  load();
});
</script>

<style scoped>
/* Hero Visual */
.hero-visual {
  height: 192px; /* h-48 */
  border-radius: 16px; /* rounded-2xl */
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24px;
  box-shadow: inset 0 2px 4px rgba(0,0,0,0.02);
}

.hero-emoji {
  font-size: 6rem; /* text-8xl approx */
  line-height: 1;
  filter: drop-shadow(0 4px 6px rgba(0,0,0,0.1));
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.page-container {
  min-height: 100vh;
  /* background-color: #F5F7FA; */ /* Removed solid color */
  background: linear-gradient(180deg, rgba(240, 253, 244, 0.6) 0%, #F5F7FA 40%, #F5F7FA 100%);
  position: relative;
  overflow-x: hidden; /* Prevent horizontal scroll from blobs */
}

.blob-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  pointer-events: none;
  overflow: hidden;
}

.blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.3;
}

.blob-1 {
  top: -10%;
  left: -5%;
  width: 500px;
  height: 500px;
  background-color: #bbf7d0; /* bg-green-200 */
}

.blob-2 {
  top: 5%;
  right: 10%;
  width: 400px;
  height: 400px;
  background-color: #fef9c3; /* bg-yellow-100 */
}

.content-wrapper {
  max-width: 1024px;
  margin: 0 auto;
  padding: 32px 16px;
  position: relative;
  z-index: 10;
}

.loading-state, .error-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 256px;
}

.animate-fade-in {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

/* Back Button */
.back-btn-container {
  margin-bottom: 16px;
}

.back-btn {
  color: #6b7280;
  transition: color 0.3s;
}

.back-btn:hover {
  color: #059669;
}

/* Header */
.header-section {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  align-items: flex-start;
  gap: 16px;
  margin-bottom: 32px;
}

@media (min-width: 768px) {
  .header-section {
    flex-direction: row;
    align-items: flex-end;
  }
}

.header-content {
  flex: 1;
}

.recipe-title {
  font-size: 2.25rem;
  font-weight: 700;
  color: #1f2937;
  margin-bottom: 12px;
  line-height: 1.2;
}

@media (min-width: 768px) {
  .recipe-title {
    font-size: 2.5rem;
  }
}

.meta-tags {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 16px;
  color: #6b7280;
  font-size: 0.875rem;
}

.meta-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  background-color: white;
  padding: 4px 12px;
  border-radius: 9999px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.icon-emerald { color: #10b981; }
.icon-orange { color: #f97316; }
.icon-blue { color: #3b82f6; }
.icon-gray { color: #9ca3af; }

.action-buttons {
  display: flex;
  gap: 12px;
}

/* Nutrition Grid */
.nutrition-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-bottom: 32px;
}

@media (min-width: 768px) {
  .nutrition-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

.nutrition-card {
  background-color: white;
  padding: 16px;
  border-radius: 16px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  border: 1px solid #f3f4f6;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  transition: box-shadow 0.3s;
}

.nutrition-card:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.nutri-label {
  font-size: 0.75rem;
  color: #9ca3af;
  font-weight: 500;
  margin-top: 4px;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.nutri-value-row {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.nutri-value {
  font-size: 1.75rem;
  font-weight: 700;
  line-height: 1;
  font-family: var(--font-family-number);
}

.text-emerald { color: #10b981; }
.text-red { color: #ef4444; }
.text-amber { color: #f59e0b; }
.text-purple { color: #8b5cf6; }

.nutri-unit {
  font-size: 0.875rem;
  color: #9ca3af;
  font-weight: normal;
}

/* Details Grid */
.details-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 32px;
}

@media (min-width: 1024px) {
  .details-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

.ingredients-col {
  grid-column: span 1;
}

.steps-col {
  grid-column: span 1;
}

@media (min-width: 1024px) {
  .steps-col {
    grid-column: span 2;
  }
}

.detail-card {
  border-radius: 16px;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  height: 100%;
}

.ingredients-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ingredient-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.ingredient-item:hover {
  background-color: #f9fafb;
}

.ing-name-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.ing-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background-color: #34d399;
}

.ing-name {
  color: #374151;
  font-weight: 500;
}

.ing-amount {
  color: #9ca3af;
  font-size: 0.875rem;
}

.empty-text {
  color: #9ca3af;
  font-style: italic;
  text-align: center;
  padding: 32px 0;
}

.steps-list {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.step-item {
  display: flex;
  gap: 16px;
}

.step-badge-wrapper {
  flex-shrink: 0;
  margin-top: 4px;
}

.step-badge {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #d1fae5;
  color: #059669;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.875rem;
  transition: all 0.3s;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.group:hover .step-badge {
  background-color: #10b981;
  color: white;
}

.step-text-wrapper {
  flex: 1;
  padding-top: 4px;
}

.step-text {
  color: #374151;
  line-height: 1.625;
  font-size: 1rem;
  transition: color 0.3s;
}

.group:hover .step-text {
  color: #111827;
}
</style>
