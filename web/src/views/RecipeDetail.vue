<template>
  <div class="recipe-detail">
    <button class="back-btn" @click="goBack">← 返回</button>

    <header class="header">
      <h1>{{ recipe.name }}</h1>
      <p v-if="recipe.description" class="desc">{{ recipe.description }}</p>
      <div class="actions">
        <button class="fav-btn" @click="toggleFavorite">
          <span v-if="isFavorite">💖 取消收藏</span>
          <span v-else>🤍 收藏</span>
        </button>
        <button class="share-btn" @click="shareRecipe">🔗 分享</button>
        <button class="share-btn" @click="goFavorites">⭐ 我的收藏</button>
      </div>
    </header>

    <section class="meta">
      <div class="nutrition">
        <div class="item">热量: <strong>{{ Math.floor(recipe.energy || 0) }} kcal</strong></div>
        <div class="item">蛋白质: <strong>{{ Math.floor(recipe.protein || 0) }} g</strong></div>
        <div class="item">碳水: <strong>{{ Math.floor(recipe.carbohydrate || 0) }} g</strong></div>
        <div class="item">脂肪: <strong>{{ Math.floor(recipe.fat || 0) }} g</strong></div>
      </div>
      <div class="cook-info">
        <div>烹饪时间: {{ recipe.cooking_time || '-' }} 分钟</div>
        <div>难度: {{ recipe.difficulty || '未知' }}</div>
      </div>
    </section>

    <section class="ingredients" v-if="recipe.ingredients && recipe.ingredients.length">
      <h3>食材</h3>
      <ul>
        <li v-for="(ing, idx) in recipe.ingredients" :key="idx">{{ formatIngredient(ing) }}</li>
      </ul>
    </section>

    <section class="steps" v-if="recipe.cooking_steps && recipe.cooking_steps.length">
      <h3>烹饪步骤</h3>
      <ol>
        <li v-for="(step, idx) in recipe.cooking_steps" :key="idx">{{ step }}</li>
      </ol>
    </section>

    <div v-if="loading" class="loading">加载中…</div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getRecipeDetail, addFavorite, removeFavorite } from '@/api/recipeApi'

const route = useRoute()
const router = useRouter()
const id = route.params.id

const recipe = ref({})
const loading = ref(false)
const isFavorite = ref(false)

const load = async () => {
  loading.value = true
  try {
    const res = await getRecipeDetail(id)
      // axios interceptor already returns response body
      recipe.value = res || {}
      isFavorite.value = !!recipe.value.is_favorite
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  load()
})

const goBack = () => router.back()
const goFavorites = () => router.push({ name: 'FavoriteList' })

const toggleFavorite = async () => {
  try {
    if (isFavorite.value) {
      await removeFavorite(id)
      isFavorite.value = false
    } else {
      await addFavorite(id)
      isFavorite.value = true
    }
  } catch (e) {
    console.error(e)
  }
}

const shareRecipe = async () => {
  const url = window.location.href
  try {
    await navigator.clipboard.writeText(url)
    alert('已复制链接到剪贴板')
  } catch (e) {
    // fallback
    const dummy = document.createElement('textarea')
    document.body.appendChild(dummy)
    dummy.value = url
    dummy.select()
    document.execCommand('copy')
    document.body.removeChild(dummy)
    alert('已复制链接到剪贴板')
  }
}

const formatIngredient = (ing) => {
  if (!ing) return ''
  if (typeof ing === 'string') return ing
  if (typeof ing === 'object') {
    const name = ing.name || ''
    const amount = ing.amount || ''
    return amount ? `${name} - ${amount}` : name
  }
  return String(ing)
}
</script>

<style scoped>
.recipe-detail { 
  max-width: 1000px; 
  margin: 0 auto; 
  padding: 2rem; 
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  min-height: 100vh;
}

.back-btn { 
  background: rgba(255,255,255,0.2); 
  border: none; 
  color: #fff; 
  cursor: pointer; 
  padding: 0.5rem 1rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
  font-size: 0.95rem;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(255,255,255,0.3);
  transform: translateX(-3px);
}

.header { 
  background: white;
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  margin-bottom: 1.5rem;
}

.header h1 { 
  margin: 0 0 0.5rem 0; 
  font-size: 2rem;
  color: #2c3e50;
  font-weight: 700;
}

.desc { 
  color: #666; 
  margin: 0.5rem 0 1.5rem 0;
  font-size: 1rem;
  line-height: 1.6;
}

.actions { 
  display: flex; 
  gap: 0.75rem; 
  flex-wrap: wrap;
  margin-top: 1rem;
}

.fav-btn, .share-btn { 
  min-width: 130px; 
  padding: 0.65rem 1.2rem; 
  border-radius: 10px; 
  border: none;
  cursor: pointer; 
  font-weight: 600;
  font-size: 0.95rem;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1); 
  transition: all 0.2s ease;
}

.fav-btn:hover, .share-btn:hover { 
  transform: translateY(-2px); 
  box-shadow: 0 6px 16px rgba(0,0,0,0.15);
}

.fav-btn { 
  background: linear-gradient(135deg, #ff6b9d 0%, #ff8e53 100%);
  color: white;
}

.share-btn { 
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
}

.meta { 
  background: white;
  padding: 1.5rem;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  margin-bottom: 1.5rem;
}

.nutrition { 
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.nutrition .item { 
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 1rem;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.08);
  font-size: 0.9rem;
  color: #2c3e50;
}

.nutrition .item strong {
  display: block;
  font-size: 1.3rem;
  margin-top: 0.3rem;
  color: #667eea;
}

.cook-info {
  display: flex;
  gap: 2rem;
  padding-top: 1rem;
  border-top: 2px solid #f0f0f0;
  color: #555;
  font-size: 0.95rem;
}

.ingredients, .steps {
  background: white;
  padding: 2rem;
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  margin-bottom: 1.5rem;
}

.ingredients h3, .steps h3 {
  margin: 0 0 1.2rem 0;
  font-size: 1.4rem;
  color: #2c3e50;
  padding-bottom: 0.8rem;
  border-bottom: 3px solid #667eea;
}

.ingredients ul { 
  padding-left: 1.5rem;
  margin: 0;
}

.ingredients li {
  padding: 0.5rem 0;
  font-size: 1rem;
  color: #555;
  line-height: 1.6;
}

.steps ol { 
  padding-left: 1.5rem;
  margin: 0;
  counter-reset: step-counter;
  list-style: none;
}

.steps li {
  padding: 1rem 0;
  font-size: 1rem;
  color: #555;
  line-height: 1.8;
  position: relative;
  padding-left: 3rem;
  counter-increment: step-counter;
}

.steps li::before {
  content: counter(step-counter);
  position: absolute;
  left: 0;
  top: 0.8rem;
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  font-size: 0.9rem;
}

.loading { 
  text-align: center; 
  color: white;
  padding: 3rem;
  font-size: 1.1rem;
}
</style>
