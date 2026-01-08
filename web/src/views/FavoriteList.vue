<template>
  <div class="favorite-page">
    <!-- 顶部导航栏 -->
    <div class="top-nav">
      <el-button @click="goBack" class="back-button" circle>
        <el-icon><ArrowLeft /></el-icon>
      </el-button>
      <div class="user-info-nav">
        <el-dropdown @command="handleCommand">
          <span class="user-info-display">
            <el-avatar :size="32" :icon="UserFilled" />
            <span class="username">{{ authStore.user?.username || '用户' }}</span>
            <el-icon><ArrowDown /></el-icon>
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="home">
                <el-icon><HomeFilled /></el-icon>
                返回首页
              </el-dropdown-item>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>
                个人档案
              </el-dropdown-item>
              <el-dropdown-item command="intake">
                <el-icon><DataLine /></el-icon>
                饮食记录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <div class="top-bar">
      <h1>我的收藏</h1>
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
import { ArrowLeft, UserFilled, ArrowDown, HomeFilled, User, DataLine } from '@element-plus/icons-vue'
import RecipeCard from '@/components/RecipeCard.vue'
import { getFavoriteList } from '@/api/recipeApi'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'

const favorites = ref([])
const loading = ref(false)
const router = useRouter()
const authStore = useAuthStore()

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

const handleCommand = (command) => {
  switch (command) {
    case 'home':
      router.push('/home')
      break
    case 'profile':
      router.push('/profile/view')
      break
    case 'intake':
      router.push('/intake')
      break
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

onMounted(loadFavorites)
</script>

<style scoped>
.favorite-page { max-width: 1080px; margin: 1.5rem auto; padding: 1rem; }

/* 顶部导航栏样式 */
.top-nav {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 10px 0;
}

.back-button {
  background-color: #f5f5f5;
  color: #333;
  border: none;
  width: 40px;
  height: 40px;
  transition: all 0.3s;
}

.back-button:hover {
  background-color: #e0e0e0;
  transform: translateX(-3px);
}

.user-info-nav {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info-display {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 20px;
  transition: all 0.3s;
}

.user-info-display:hover {
  background-color: #f5f5f5;
}

.username {
  font-size: 14px;
  font-weight: 500;
  color: #333;
}

.top-bar { display:flex; align-items:center; justify-content:center; gap:1rem; margin-bottom:1rem; }
.back-btn { background: none; border: none; color: #555; cursor: pointer; }
.loading, .empty { text-align:center; color:#666; margin-top:2rem; }
.grid { display:grid; gap:1rem; grid-template-columns: repeat(auto-fill, minmax(260px, 1fr)); }
</style>
