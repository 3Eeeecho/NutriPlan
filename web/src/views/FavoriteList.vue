<template>
  <PageLayout>
    <div class="favorite-container">
      <section class="hero-section">
        <div class="stats-row">
          <div class="stat-card">
            <div class="stat-label">总食谱数量</div>
            <div class="stat-value">{{ favoriteStats.total }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">已分类食谱</div>
            <div class="stat-value">{{ favoriteStats.categorizedCount }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">高蛋白食谱</div>
            <div class="stat-value">{{ favoriteStats.highProteinCount }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-label">覆盖餐别</div>
            <div class="stat-value">{{ favoriteStats.coveredMealTypes }}<span class="stat-unit">/4</span></div>
          </div>
        </div>
      </section>

      <section class="toolbar-section">
        <n-input
          v-model:value="keyword"
          clearable
          placeholder="搜索食谱名称或描述"
          class="search-input"
        />
        <n-select
          v-model:value="mealFilter"
          :options="mealFilterOptions"
          class="filter-select"
        />
        <n-select
          v-model:value="sortBy"
          :options="sortOptions"
          class="sort-select"
        />
        <n-button quaternary @click="resetFilters">重置</n-button>
      </section>

      <div class="result-bar">
        <span class="result-count">共 {{ sortedFavorites.length }} 条结果</span>
        <n-tag v-if="mealFilter !== 'all'" size="small" round :bordered="false" type="success">
          {{ formatMealType(mealFilter) }}
        </n-tag>
      </div>

      <n-spin :show="loading">
        <div v-if="!loading && favorites.length === 0" class="empty-state">
          <n-empty description="还没有收藏食谱" />
        </div>

        <div v-else-if="!loading && sortedFavorites.length === 0" class="empty-state">
          <n-empty description="没有符合筛选条件的食谱">
            <template #extra>
              <n-button tertiary type="primary" @click="resetFilters">清空筛选</n-button>
            </template>
          </n-empty>
        </div>

        <div v-else class="recipe-grid">
          <RecipeCard
            v-for="item in sortedFavorites"
            :key="item.id"
            :title="formatMealType(item.resolvedMealType)"
            :recipe="item"
            :icon="pickIcon(item.resolvedMealType)"
          />
        </div>
      </n-spin>
    </div>
  </PageLayout>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { NSpin, NEmpty, NButton, NInput, NSelect, NTag } from 'naive-ui'
import { getFavoriteList } from '@/api/recipeApi'
import PageLayout from '@/components/layout/PageLayout.vue'
import RecipeCard from '@/components/RecipeCard.vue'

const favorites = ref([])
const loading = ref(false)
const keyword = ref('')
const mealFilter = ref('all')
const sortBy = ref('energy_desc')

const sortOptions = [
  { label: '热量从高到低', value: 'energy_desc' },
  { label: '热量从低到高', value: 'energy_asc' },
  { label: '蛋白质从高到低', value: 'protein_desc' },
  { label: '名称 A-Z', value: 'name_asc' }
]

const normalizeMealType = (value = '') => {
  const text = String(value).trim().toLowerCase()
  const map = {
    breakfast: 'breakfast',
    lunch: 'lunch',
    dinner: 'dinner',
    snack: 'snack',
    早餐: 'breakfast',
    午餐: 'lunch',
    晚餐: 'dinner',
    加餐: 'snack'
  }
  return map[text] || ''
}

const inferMealTypeByName = (name = '') => {
  const text = String(name).toLowerCase()
  if (!text) return ''
  if (text.includes('早餐') || text.includes('早饭') || text.includes('粥') || text.includes('三明治') || text.includes('面包') || text.includes('牛奶')) {
    return 'breakfast'
  }
  if (text.includes('午餐') || text.includes('午饭') || text.includes('便当') || text.includes('盖饭')) {
    return 'lunch'
  }
  if (text.includes('晚餐') || text.includes('晚饭') || text.includes('夜宵') || text.includes('汤')) {
    return 'dinner'
  }
  if (text.includes('加餐') || text.includes('零食') || text.includes('酸奶') || text.includes('水果') || text.includes('坚果')) {
    return 'snack'
  }
  return ''
}

const normalizedFavorites = computed(() => favorites.value.map((item) => ({
  ...item,
  resolvedMealType: normalizeMealType(item.meal_type || item.mealType) || inferMealTypeByName(item.name),
  energy: Number(item.energy || 0),
  protein: Number(item.protein || 0),
  carbohydrate: Number(item.carbohydrate || 0),
  fat: Number(item.fat || 0)
})))

const mealFilterOptions = computed(() => {
  const countMap = normalizedFavorites.value.reduce((acc, item) => {
    const type = item.resolvedMealType
    if (type) {
      acc[type] = (acc[type] || 0) + 1
    }
    return acc
  }, {})

  return [
    { label: `全部（${normalizedFavorites.value.length}）`, value: 'all' },
    { label: `早餐（${countMap.breakfast || 0}）`, value: 'breakfast' },
    { label: `午餐（${countMap.lunch || 0}）`, value: 'lunch' },
    { label: `晚餐（${countMap.dinner || 0}）`, value: 'dinner' },
    { label: `加餐（${countMap.snack || 0}）`, value: 'snack' }
  ]
})

const filteredFavorites = computed(() => {
  const text = keyword.value.trim().toLowerCase()

  return normalizedFavorites.value.filter((item) => {
    const matchType = mealFilter.value === 'all' || item.resolvedMealType === mealFilter.value
    const sourceText = `${item.name || ''} ${item.description || ''}`.toLowerCase()
    const matchKeyword = !text || sourceText.includes(text)
    return matchType && matchKeyword
  })
})

const sortedFavorites = computed(() => {
  const list = [...filteredFavorites.value]

  if (sortBy.value === 'energy_desc') {
    return list.sort((first, second) => second.energy - first.energy)
  }
  if (sortBy.value === 'energy_asc') {
    return list.sort((first, second) => first.energy - second.energy)
  }
  if (sortBy.value === 'protein_desc') {
    return list.sort((first, second) => second.protein - first.protein)
  }
  if (sortBy.value === 'name_asc') {
    return list.sort((first, second) => (first.name || '').localeCompare(second.name || '', 'zh-CN'))
  }
  return list
})

const favoriteStats = computed(() => {
  const total = normalizedFavorites.value.length
  const categorizedCount = normalizedFavorites.value.filter((item) => !!item.resolvedMealType).length
  const highProteinCount = normalizedFavorites.value.filter((item) => item.protein >= 20).length
  const coveredMealTypeSet = new Set(
    normalizedFavorites.value
      .map((item) => item.resolvedMealType)
      .filter((type) => !!type)
  )

  return {
    total,
    categorizedCount,
    highProteinCount,
    coveredMealTypes: coveredMealTypeSet.size
  }
})

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

const resetFilters = () => {
  keyword.value = ''
  mealFilter.value = 'all'
  sortBy.value = 'energy_desc'
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
  max-width: 1240px;
  margin: 0 auto;
  padding: 20px 24px 28px;
}

.hero-section {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.08), rgba(59, 130, 246, 0.04));
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 14px;
  margin-bottom: 18px;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.stat-card {
  border-radius: 12px;
  border: 1px solid var(--color-border-primary);
  background: rgba(255, 255, 255, 0.78);
  padding: 10px 12px;
}

.stat-label {
  font-size: 12px;
  color: var(--color-text-secondary);
  margin-bottom: 6px;
}

.stat-value {
  font-size: 20px;
  font-weight: 800;
  color: var(--color-text-primary);
  line-height: 1;
}

.stat-unit {
  margin-left: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
}

.toolbar-section {
  display: grid;
  grid-template-columns: 1.3fr 180px 180px auto;
  gap: 10px;
  margin-bottom: 12px;
}

.search-input,
.filter-select,
.sort-select {
  width: 100%;
}

.result-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.result-count {
  font-size: 13px;
  color: var(--color-text-secondary);
}

.empty-state {
  padding: 60px 0;
  display: flex;
  justify-content: center;
}

.recipe-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(290px, 1fr));
  gap: 18px;
}

@media (max-width: 768px) {
  .favorite-container {
    padding: 14px;
  }

  .hero-section {
    padding: 14px;
  }

  .stats-row {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .toolbar-section {
    grid-template-columns: 1fr;
  }

  .result-bar {
    margin-bottom: 10px;
  }

  .recipe-grid {
    grid-template-columns: 1fr;
    gap: 16px;
  }
}
</style>
