<template>
  <div class="favorite-page">
    <div class="top-bar">
      <button class="back-btn" @click="goBack">← 返回</button>
      <h1>我的收藏</h1>
      <div class="placeholder"></div>
    </div>

    <div v-if="loading" class="loading">加载中…</div>
    <div v-else-if="favorites.length === 0" class="empty">还没有收藏食谱，去推荐页看看吧～</div>
    <div v-else class="grid">
      <RecipeCard
        v-for="item in favorites"
        :key="item.id"
        :title="item.meal_type || '食谱'"
        :recipe="item"
        :icon="pickIcon(item.meal_type)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import RecipeCard from '@/components/RecipeCard.vue'
import { getFavoriteList } from '@/api/recipeApi'
import { useRouter } from 'vue-router'

const favorites = ref([])
const loading = ref(false)
const router = useRouter()

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

const goBack = () => router.back()

const pickIcon = (meal) => {
  const map = {
    breakfast: '🥞',
    lunch: '🍱',
    dinner: '🍽️',
    snack: '🍪'
  }
  return map[meal] || '🍽️'
}

onMounted(loadFavorites)
</script>

<style scoped>
.favorite-page { max-width: 1080px; margin: 1.5rem auto; padding: 1rem; }
.top-bar { display:flex; align-items:center; justify-content:space-between; gap:1rem; margin-bottom:1rem; }
.back-btn { background: none; border: none; color: #555; cursor: pointer; }
.loading, .empty { text-align:center; color:#666; margin-top:2rem; }
.grid { display:grid; gap:1rem; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); }
</style>
