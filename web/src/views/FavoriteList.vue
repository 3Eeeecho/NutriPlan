<template>
  <PageLayout>
    <TopNavigation />
    
    <div class="favorite-container">
      <div class="page-header">
        <div class="header-left">
          <BackButton />
        </div>
        <h1 class="page-title">我的收藏</h1>
        <div class="header-right"></div>
      </div>

      <n-spin :show="loading">
        <div v-if="!loading && favorites.length === 0" class="empty-state">
          <n-empty description="还没有收藏食谱，去推荐页看看吧～">
            <template #extra>
              <n-button type="primary" @click="router.push('/recommend')">
                去逛逛
              </n-button>
            </template>
          </n-empty>
        </div>

        <div v-else class="recipe-grid">
          <RecipeCard
            v-for="item in favorites"
            :key="item.id"
            :title="formatMealType(item.meal_type)"
            :recipe="item"
            :icon="pickIcon(item.meal_type)"
          />
        </div>
      </n-spin>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NSpin, NEmpty, NButton } from 'naive-ui'
import { getFavoriteList } from '@/api/recipeApi'
import PageLayout from '@/components/layout/PageLayout.vue'
import TopNavigation from '@/components/layout/TopNavigation.vue'
import BackButton from '@/components/layout/BackButton.vue'
import RecipeCard from '@/components/RecipeCard.vue'

const router = useRouter()
const favorites = ref([])
const loading = ref(false)

const loadFavorites = async () => {
  loading.value = true
  try {
    const res = await getFavoriteList()
    favorites.value = res.recipes || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const pickIcon = (meal) => {
  const map = {
    breakfast: '🥞',
    lunch: '🍱',
    dinner: '🍽️',
    snack: '🍪'
  }
  return map[meal] || '🍽️'
}

const formatMealType = (type) => {
  const map = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐',
    snack: '加餐'
  }
  return map[type] || '食谱'
}

onMounted(loadFavorites)
</script>

<style scoped>
.favorite-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32px;
}

.header-left, .header-right {
  width: 80px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  text-align: center;
}

.empty-state {
  padding: 60px 0;
  display: flex;
  justify-content: center;
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

@media (max-width: 768px) {
  .favorite-container {
    padding: 16px;
  }
  
  .page-header {
    margin-bottom: 24px;
  }
  
  .recipe-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
