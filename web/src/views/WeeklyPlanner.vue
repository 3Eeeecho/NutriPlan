<template>
  <div class="weekly-page">
    <div class="weekly-container">
      <div class="weekly-header">
        <div class="left-actions">
          <n-button quaternary circle @click="router.push('/home')">
            <template #icon><n-icon><ChevronBackOutline /></n-icon></template>
          </n-button>
          <div class="view-toggle">
            <n-button size="small" @click="router.push('/home')">日</n-button>
            <n-button size="small" type="primary">周</n-button>
          </div>
          <h2 class="title">本周计划</h2>
        </div>

        <div class="right-actions">
          <n-button @click="loadWeeklyData" :loading="loading" type="primary" ghost>
            <template #icon><n-icon><RefreshOutline /></n-icon></template>
            重新加载
          </n-button>
        </div>
      </div>

      <div v-if="loading" class="loading-wrap">
        <n-spin size="large" />
      </div>

      <div v-else class="week-columns">
        <div
          v-for="day in weekColumns"
          :key="day.date"
          class="day-column"
          :class="{ 'is-today': day.isToday }"
        >
          <div class="day-header">
            <div class="day-date">{{ day.displayDate }}</div>
            <div class="day-week-label">{{ day.weekLabel }}</div>
            <div class="day-kcal">{{ day.totalEnergy }} 千卡</div>
          </div>

          <div v-if="day.records.length === 0" class="day-empty">
            暂无已规划餐次
          </div>

          <div v-else class="meal-section-list">
            <div
              v-for="meal in getMealSections(day.records)"
              :key="`${day.date}-${meal.type}`"
              class="meal-section"
            >
              <div class="meal-title-row">
                <span class="meal-title">{{ getMealTypeLabel(meal.type) }}</span>
                <span class="meal-total">{{ meal.energy }} 千卡</span>
              </div>

              <div
                v-for="item in meal.items"
                :key="item.ID || item.id"
                class="meal-item-card"
              >
                <div class="meal-item-name">{{ item.foodName || item.food_name }}</div>
                <div class="meal-item-meta">
                  <span>{{ Math.round(item.intakeAmount || item.intake_amount || 0) }} {{ item.intakeUnit || 'g' }}</span>
                  <span>{{ Math.round(item.calculatedEnergy || item.calculated_energy || 0) }} kcal</span>
                  <span>蛋 {{ Number(item.calculatedProtein || item.calculated_protein || 0).toFixed(1) }}</span>
                  <span>碳 {{ Number(item.calculatedCarb || item.calculated_carb || 0).toFixed(1) }}</span>
                  <span>脂 {{ Number(item.calculatedFat || item.calculated_fat || 0).toFixed(1) }}</span>
                  <span v-if="item.isCompletedRecipe" class="done-tag">已完成食谱</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NIcon, NSpin, useMessage } from 'naive-ui'
import { ChevronBackOutline, RefreshOutline } from '@vicons/ionicons5'
import { getWeeklyReport } from '@/api/intakeApi'
import { useAuthStore } from '@/store/auth'
import { buildCompletedRecordsForDate } from '@/utils/completedRecipes'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const weeklyReport = ref({ daily_data: [] })

const weekColumns = computed(() => {
  const list = weeklyReport.value?.daily_data || []
  const today = new Date()
  const todayKey = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`

  return list.map((day) => {
    const date = day.date
    const dateObj = new Date(date)
    const backendRecords = Array.isArray(day.records) ? day.records : []
    const completedRecords = buildCompletedRecordsForDate(authStore.user?.id, date)
    const records = [...backendRecords, ...completedRecords]
    const completedEnergy = completedRecords.reduce((sum, item) => sum + Number(item.calculatedEnergy || 0), 0)

    return {
      date,
      isToday: date === todayKey,
      displayDate: `${dateObj.getMonth() + 1}月${dateObj.getDate()}日`,
      weekLabel: dateObj.toLocaleDateString('zh-CN', { weekday: 'long' }),
      totalEnergy: Math.round((day.total_energy || day.totalEnergy || 0) + completedEnergy),
      records
    }
  })
})

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
  return map[text] || 'snack'
}

const getMealSections = (records = []) => {
  const order = ['breakfast', 'lunch', 'dinner', 'snack']
  const grouped = { breakfast: [], lunch: [], dinner: [], snack: [] }

  records.forEach((record) => {
    const type = normalizeMealType(record.mealType || record.meal_type)
    grouped[type].push(record)
  })

  return order
    .map((type) => {
      const items = grouped[type]
      const energy = Math.round(items.reduce((sum, item) => sum + Number(item.calculatedEnergy || item.calculated_energy || 0), 0))
      return { type, items, energy }
    })
    .filter((group) => group.items.length > 0)
}

const getMealTypeLabel = (type) => {
  const map = {
    breakfast: '早餐',
    lunch: '午餐',
    dinner: '晚餐',
    snack: '加餐'
  }
  return map[type] || '餐次'
}

const loadWeeklyData = async () => {
  loading.value = true
  try {
    const res = await getWeeklyReport()
    weeklyReport.value = res || { daily_data: [] }
  } catch (error) {
    message.error('加载周食谱失败')
  } finally {
    loading.value = false
  }
}

onMounted(loadWeeklyData)
</script>

<style scoped>
.weekly-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.weekly-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 18px;
  background: #141923;
  border-radius: 16px;
  border: 1px solid #242c3b;
}

.weekly-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.left-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.title {
  margin: 0;
  color: #f3f4f6;
  font-size: 32px;
  font-weight: 700;
}

.view-toggle {
  display: flex;
  gap: 6px;
}

.week-columns {
  display: grid;
  grid-template-columns: repeat(7, minmax(220px, 1fr));
  gap: 12px;
  overflow-x: auto;
  padding-bottom: 8px;
}

.day-column {
  min-height: 620px;
  background: #1b2230;
  border: 1px solid #2a3345;
  border-radius: 12px;
  padding: 12px;
}

.day-column.is-today {
  border-color: #f97316;
  box-shadow: inset 0 0 0 1px rgba(249, 115, 22, 0.35);
}

.day-header {
  margin-bottom: 12px;
}

.day-date {
  font-size: 14px;
  color: #e5e7eb;
}

.day-week-label {
  margin-top: 2px;
  font-size: 33px;
  font-weight: 700;
  color: #60a5fa;
}

.day-column.is-today .day-week-label {
  color: #fb923c;
}

.day-kcal {
  margin-top: 6px;
  color: #a7f3d0;
  font-size: 13px;
}

.day-empty {
  border: 1px dashed #334155;
  border-radius: 10px;
  padding: 22px 8px;
  text-align: center;
  color: #94a3b8;
  font-size: 13px;
}

.meal-section-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.meal-section {
  background: #232b3a;
  border: 1px solid #334155;
  border-radius: 10px;
  padding: 8px;
}

.meal-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.meal-title {
  color: #f3f4f6;
  font-weight: 700;
  font-size: 27px;
}

.meal-total {
  color: #94a3b8;
  font-size: 12px;
}

.meal-item-card {
  background: #111827;
  border: 1px solid #263043;
  border-radius: 8px;
  padding: 8px;
}

.meal-item-card + .meal-item-card {
  margin-top: 8px;
}

.meal-item-name {
  color: #60a5fa;
  font-weight: 600;
  font-size: 15px;
  margin-bottom: 6px;
}

.meal-item-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  color: #cbd5e1;
  font-size: 12px;
}

.done-tag {
  color: #34d399;
  border: 1px solid rgba(52, 211, 153, 0.35);
  padding: 0 6px;
  border-radius: 999px;
}

.loading-wrap {
  display: flex;
  justify-content: center;
  padding: 56px 0;
}

@media (max-width: 1024px) {
  .weekly-container {
    border-radius: 12px;
    padding: 12px;
  }

  .weekly-header {
    align-items: flex-start;
    flex-direction: column;
    gap: 10px;
  }

  .title {
    font-size: 24px;
  }
}
</style>
