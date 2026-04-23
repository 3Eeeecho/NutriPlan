<template>
  <div class="page-container">
    <div class="content-wrapper">
      <div class="page-header">
        <div class="header-left">
          <n-button quaternary circle size="large" @click="router.push('/home')">
            <template #icon>
              <n-icon><ArrowBack /></n-icon>
            </template>
          </n-button>
          <div>
            <h1 class="page-title">饮食记录</h1>
            <p class="page-subtitle">{{ currentDate }}</p>
          </div>
        </div>
        <div class="header-right">
          <n-button quaternary circle @click="refreshView">
            <template #icon><n-icon><Refresh /></n-icon></template>
          </n-button>
        </div>
      </div>

      <n-grid :x-gap="24" :y-gap="24" cols="1 l:3" responsive="screen">
        <n-gi span="1">
          <div class="left-column">
            <n-card :bordered="false" class="dashboard-card">
              <div class="calorie-ring-wrapper">
                <n-progress
                  type="circle"
                  :percentage="energyPercentage"
                  :color="energyColor"
                  :stroke-width="12"
                  class="calorie-ring"
                >
                  <div class="ring-content">
                    <span class="ring-val">{{ Math.floor(nutritionStatus.total_energy || 0) }}</span>
                    <span class="ring-label">/ {{ Math.floor(nutritionStatus.target_energy || 2000) }} kcal</span>
                  </div>
                </n-progress>
                <div class="dashboard-title">今日热量摄入</div>
              </div>

              <div class="macro-bars">
                <div class="macro-item">
                  <div class="macro-header">
                    <span class="macro-label text-red">蛋白质</span>
                    <span class="macro-val">{{ nutritionStatus.total_protein?.toFixed(1) || 0 }} / {{ nutritionStatus.target_protein?.toFixed(1) || 0 }}g</span>
                  </div>
                  <n-progress
                    type="line"
                    :percentage="getPercentage(nutritionStatus.total_protein, nutritionStatus.target_protein)"
                    color="#ef4444"
                    :height="8"
                    :show-indicator="false"
                    class="rounded-progress"
                  />
                </div>
                <div class="macro-item">
                  <div class="macro-header">
                    <span class="macro-label text-amber">碳水</span>
                    <span class="macro-val">{{ nutritionStatus.total_carbohydrate?.toFixed(1) || 0 }} / {{ nutritionStatus.target_carbohydrate?.toFixed(1) || 0 }}g</span>
                  </div>
                  <n-progress
                    type="line"
                    :percentage="getPercentage(nutritionStatus.total_carbohydrate, nutritionStatus.target_carbohydrate)"
                    color="#f59e0b"
                    :height="8"
                    :show-indicator="false"
                    class="rounded-progress"
                  />
                </div>
                <div class="macro-item">
                  <div class="macro-header">
                    <span class="macro-label text-purple">脂肪</span>
                    <span class="macro-val">{{ nutritionStatus.total_fat?.toFixed(1) || 0 }} / {{ nutritionStatus.target_fat?.toFixed(1) || 0 }}g</span>
                  </div>
                  <n-progress
                    type="line"
                    :percentage="getPercentage(nutritionStatus.total_fat, nutritionStatus.target_fat)"
                    color="#8b5cf6"
                    :height="8"
                    :show-indicator="false"
                    class="rounded-progress"
                  />
                </div>
              </div>
            </n-card>

            <n-card :bordered="false" title="快捷操作" class="actions-card">
              <div class="quick-actions">
                <n-button block type="primary" size="large" class="camera-btn" @click="triggerAIUpload">
                  <template #icon><n-icon><Camera /></n-icon></template>
                  拍照识别
                </n-button>
                <n-button block ghost type="primary" size="large" @click="showAddModal = true">
                  <template #icon><n-icon><Add /></n-icon></template>
                  手动记一笔
                </n-button>
                <input
                  ref="fileInputRef"
                  type="file"
                  accept="image/*"
                  style="display: none"
                  @change="handleImageSelect"
                />
              </div>
            </n-card>
          </div>
        </n-gi>

        <n-gi span="2">
          <n-card :bordered="false" title="今日饮食记录" class="timeline-card">
            <template #header-extra>
              <n-tag :bordered="false" type="default">
                共 {{ nutritionStatus.records?.length || 0 }} 条记录
              </n-tag>
            </template>

            <div v-if="loading" class="loading-placeholder">
              <n-spin size="medium" />
            </div>

            <div v-else-if="!nutritionStatus.records || nutritionStatus.records.length === 0" class="empty-state">
              <n-empty description="今天还没有记录，去添加一餐吧">
                <template #extra>
                  <n-button type="primary" @click="triggerAIUpload">
                    <template #icon><n-icon><Camera /></n-icon></template>
                    拍照识别
                  </n-button>
                </template>
              </n-empty>
            </div>

            <n-timeline v-else class="meal-timeline">
              <n-timeline-item
                v-for="record in sortedRecords"
                :key="record.ID"
                :type="getMealTypeColor(record.mealType)"
                :title="getMealTypeLabel(record.mealType)"
                :content="record.foodName"
                :time="formatTime(record.CreatedAt)"
              >
                <template #default>
                  <div class="record-card">
                    <div class="record-main">
                      <div class="food-name">{{ record.foodName }}</div>
                      <div class="food-meta">
                        <n-tag size="small" :bordered="false" class="meta-tag">{{ record.intakeAmount }}g</n-tag>
                        <span class="meta-divider">|</span>
                        <span class="meta-val">{{ Math.floor(record.calculatedEnergy) }} kcal</span>
                      </div>
                    </div>
                    <div class="record-macros">
                      <div class="mini-macro">
                        <span class="mm-label">蛋</span>
                        <span class="mm-val">{{ record.calculatedProtein }}</span>
                      </div>
                      <div class="mini-macro">
                        <span class="mm-label">碳</span>
                        <span class="mm-val">{{ record.calculatedCarb }}</span>
                      </div>
                      <div class="mini-macro">
                        <span class="mm-label">脂</span>
                        <span class="mm-val">{{ record.calculatedFat }}</span>
                      </div>
                    </div>
                    <div class="record-actions">
                      <n-popconfirm @positive-click="handleDelete(record.ID)" positive-text="确定" negative-text="取消">
                        <template #trigger>
                          <n-button size="small" quaternary circle type="error">
                            <template #icon><n-icon><Trash /></n-icon></template>
                          </n-button>
                        </template>
                        确定删除这条记录吗？
                      </n-popconfirm>
                    </div>
                  </div>
                </template>
              </n-timeline-item>
            </n-timeline>
          </n-card>
        </n-gi>
      </n-grid>
    </div>

    <n-modal v-model:show="showAddModal" preset="card" title="添加饮食记录" style="width: 600px; max-width: 90vw;">
      <n-form :model="recordForm" label-placement="left" label-width="80" require-mark-placement="right-hanging">
        <n-grid :cols="2" :x-gap="12">
          <n-gi :span="2">
            <n-form-item label="餐点类型" path="meal_type">
              <n-select v-model:value="recordForm.meal_type" :options="mealOptions" placeholder="选择餐点" />
            </n-form-item>
          </n-gi>
          <n-gi :span="2">
            <n-form-item label="食物描述/名称" path="food_name">
              <n-input-group>
                <n-input v-model:value="recordForm.food_name" placeholder="例如：一碗牛肉面，少辣" />
                <n-button type="primary" ghost @click="handleAIAnalyze" :loading="analyzingText" :disabled="!recordForm.food_name">
                  <template #icon><n-icon><Sparkles /></n-icon></template>
                  AI 估算
                </n-button>
              </n-input-group>
            </n-form-item>
          </n-gi>
          <n-gi :span="2">
            <n-form-item label="分量 (g)" path="intake_amount">
              <n-input-number v-model:value="recordForm.intake_amount" :step="10" placeholder="估算重量" style="width: 100%" />
            </n-form-item>
          </n-gi>

          <n-gi :span="2">
            <n-divider dashed>营养成分</n-divider>
          </n-gi>

          <n-gi>
            <n-form-item label="热量" path="calculated_energy">
              <n-input-number v-model:value="recordForm.calculated_energy" placeholder="kcal" :show-button="false">
                <template #suffix>kcal</template>
              </n-input-number>
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="蛋白质" path="calculated_protein">
              <n-input-number v-model:value="recordForm.calculated_protein" placeholder="g" :show-button="false">
                <template #suffix>g</template>
              </n-input-number>
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="碳水" path="calculated_carb">
              <n-input-number v-model:value="recordForm.calculated_carb" placeholder="g" :show-button="false">
                <template #suffix>g</template>
              </n-input-number>
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="脂肪" path="calculated_fat">
              <n-input-number v-model:value="recordForm.calculated_fat" placeholder="g" :show-button="false">
                <template #suffix>g</template>
              </n-input-number>
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>
      <template #footer>
        <div class="modal-actions">
          <n-button @click="showAddModal = false">取消</n-button>
          <n-button type="primary" :loading="adding" @click="handleAddRecord">确认添加</n-button>
        </div>
      </template>
    </n-modal>

    <div v-if="recognizing" class="ai-loading-mask">
      <div class="ai-loading-content">
        <n-spin size="large" />
        <p>AI 正在分析食物营养...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import {
  NButton, NIcon, NGrid, NGi, NCard, NProgress, NTimeline, NTimelineItem,
  NTag, NEmpty, NSpin, NModal, NForm, NFormItem, NInput, NInputNumber,
  NSelect, NDivider, useMessage, NPopconfirm, NInputGroup
} from 'naive-ui'
import { ArrowBack, Refresh, Add, Camera, Trash, Sparkles } from '@vicons/ionicons5'
import { getTodayStatus, addIntakeRecord, deleteIntakeRecord } from '@/api/intakeApi'
import { recognizeFood, analyzeFoodText } from '@/api/foodRecognitionApi'

const router = useRouter()
const route = useRoute()
const message = useMessage()

const loading = ref(false)
const adding = ref(false)
const recognizing = ref(false)
const analyzingText = ref(false)
const showAddModal = ref(false)
const nutritionStatus = ref({})
const currentDate = ref('')
const fileInputRef = ref(null)
const aiData = ref(null)

const recordForm = ref({
  meal_type: 'breakfast',
  food_name: '',
  intake_amount: 100,
  calculated_energy: 0,
  calculated_protein: 0,
  calculated_carb: 0,
  calculated_fat: 0
})

const mealOptions = [
  { label: '早餐', value: 'breakfast' },
  { label: '午餐', value: 'lunch' },
  { label: '晚餐', value: 'dinner' },
  { label: '加餐', value: 'snack' }
]

const sortedRecords = computed(() => {
  if (!nutritionStatus.value.records) return []
  return [...nutritionStatus.value.records].reverse()
})

const energyPercentage = computed(() => {
  const current = nutritionStatus.value.total_energy || 0
  const target = nutritionStatus.value.target_energy || 2000
  return Math.min(100, (current / target) * 100)
})

const energyColor = computed(() => {
  if (energyPercentage.value > 100) return '#ef4444'
  return '#0d9488'
})

watch(() => recordForm.value.intake_amount, (newVal) => {
  if (aiData.value && newVal > 0) {
    const factor = newVal / 100
    recordForm.value.calculated_energy = Math.round(aiData.value.calories * factor)
    recordForm.value.calculated_protein = Number((aiData.value.protein * factor).toFixed(1))
    recordForm.value.calculated_carb = Number((aiData.value.carb * factor).toFixed(1))
    recordForm.value.calculated_fat = Number((aiData.value.fat * factor).toFixed(1))
  }
})

onMounted(() => {
  currentDate.value = new Date().toLocaleDateString('zh-CN', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
  loadTodayData()
  openAddFoodFromRoute()
})

const openAddFoodFromRoute = async () => {
  if (route.query?.action !== 'add-food') return

  const requestedMealType = String(route.query?.mealType || '').trim().toLowerCase()
  const mealType = ['breakfast', 'lunch', 'dinner', 'snack'].includes(requestedMealType)
    ? requestedMealType
    : 'breakfast'

  recordForm.value = {
    meal_type: mealType,
    food_name: '',
    intake_amount: 100,
    calculated_energy: 0,
    calculated_protein: 0,
    calculated_carb: 0,
    calculated_fat: 0
  }
  showAddModal.value = true

  await router.replace({ path: '/intake' })
}

const loadTodayData = async () => {
  loading.value = true
  try {
    const res = await getTodayStatus()
    nutritionStatus.value = res || {}
  } catch {
    message.error('加载今日数据失败')
  } finally {
    loading.value = false
  }
}

const refreshView = async () => {
  await loadTodayData()
  message.success('记录已刷新')
}

const triggerAIUpload = () => {
  fileInputRef.value?.click()
}

const handleImageSelect = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return

  recognizing.value = true
  try {
    const res = await recognizeFood(file)
    recordForm.value.food_name = res.dish_name
    recordForm.value.intake_amount = res.estimated_weight
    recordForm.value.calculated_energy = Math.floor(res.total_calories)
    recordForm.value.calculated_protein = Math.floor(res.total_protein)
    recordForm.value.calculated_carb = Math.floor(res.total_carbs)
    recordForm.value.calculated_fat = Math.floor(res.total_fat)

    aiData.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }

    showAddModal.value = true
    message.success(`识别成功：${res.dish_name}`)
  } catch {
    message.error('识别失败，请重试')
  } finally {
    recognizing.value = false
    event.target.value = ''
  }
}

const handleAIAnalyze = async () => {
  if (!recordForm.value.food_name) {
    message.warning('请输入食物描述')
    return
  }

  analyzingText.value = true
  try {
    const res = await analyzeFoodText(recordForm.value.food_name)
    recordForm.value.food_name = res.dish_name
    recordForm.value.intake_amount = res.estimated_weight
    recordForm.value.calculated_energy = Math.floor(res.total_calories)
    recordForm.value.calculated_protein = Math.floor(res.total_protein)
    recordForm.value.calculated_carb = Math.floor(res.total_carbs)
    recordForm.value.calculated_fat = Math.floor(res.total_fat)

    aiData.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }

    message.success('AI 分析完成')
  } catch {
    message.error('分析失败，请重试')
  } finally {
    analyzingText.value = false
  }
}

const handleAddRecord = async () => {
  if (!recordForm.value.food_name) {
    message.warning('请输入食物名称')
    return
  }

  adding.value = true
  try {
    await addIntakeRecord(recordForm.value)
    message.success('记录添加成功')
    showAddModal.value = false
    recordForm.value = {
      meal_type: 'breakfast',
      food_name: '',
      intake_amount: 100,
      calculated_energy: 0,
      calculated_protein: 0,
      calculated_carb: 0,
      calculated_fat: 0
    }
    aiData.value = null
    await loadTodayData()
  } catch {
    message.error('添加失败')
  } finally {
    adding.value = false
  }
}

const handleDelete = async (id) => {
  try {
    await deleteIntakeRecord(id)
    message.success('已删除')
    await loadTodayData()
  } catch {
    message.error('删除失败')
  }
}

const getPercentage = (val, target) => {
  if (!target) return 0
  return Math.min(100, (val / target) * 100)
}

const getMealTypeLabel = (type) => {
  const map = { breakfast: '早餐', lunch: '午餐', dinner: '晚餐', snack: '加餐' }
  return map[type] || type
}

const getMealTypeColor = (type) => {
  const map = { breakfast: 'info', lunch: 'warning', dinner: 'success', snack: 'default' }
  return map[type] || 'default'
}

const formatTime = (isoString) => {
  if (!isoString) return ''
  const d = new Date(isoString)
  return d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}
</script>

<style scoped>
.page-container {
  min-height: 100vh;
  background-color: #f5f7fa;
  position: relative;
}

.content-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 32px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1f2937;
  margin: 0;
  line-height: 1.2;
}

.page-subtitle {
  font-size: 14px;
  color: #6b7280;
  margin: 4px 0 0;
}

.dashboard-card,
.actions-card,
.timeline-card {
  border-radius: 16px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);
}

.actions-card {
  margin-top: 24px;
}

.calorie-ring-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24px;
}

.ring-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1;
}

.ring-val {
  font-size: 24px;
  font-weight: 800;
  color: #1f2937;
  white-space: nowrap;
}

.ring-label {
  font-size: 12px;
  color: #9ca3af;
  margin-top: 2px;
  white-space: nowrap;
}

.dashboard-title {
  margin-top: 12px;
  font-size: 14px;
  color: #4b5563;
  font-weight: 500;
}

.macro-bars {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.macro-header {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  margin-bottom: 4px;
}

.macro-label {
  font-weight: 600;
}

.text-red {
  color: #ef4444;
}

.text-amber {
  color: #f59e0b;
}

.text-purple {
  color: #8b5cf6;
}

.macro-val {
  color: #6b7280;
}

.rounded-progress :deep(.n-progress-graph-line-rail),
.rounded-progress :deep(.n-progress-graph-line-fill) {
  border-radius: 9999px;
}

.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.camera-btn {
  box-shadow: 0 10px 15px -3px rgba(13, 148, 136, 0.3), 0 4px 6px -2px rgba(13, 148, 136, 0.1);
  transition: transform 0.2s;
}

.camera-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 25px -5px rgba(13, 148, 136, 0.4), 0 10px 10px -5px rgba(13, 148, 136, 0.1);
}

.timeline-card {
  height: 100%;
}

.meal-timeline {
  padding: 8px 0;
}

.record-card {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto auto;
  gap: 16px;
  align-items: center;
  padding: 16px;
  background: #ffffff;
  border: 1px solid #eef2f7;
  border-radius: 14px;
}

.record-main {
  min-width: 0;
}

.food-name {
  font-size: 16px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 6px;
}

.food-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: #6b7280;
}

.meta-tag {
  background: #eef6f4;
  color: #0f766e;
}

.record-macros {
  display: flex;
  gap: 10px;
}

.mini-macro {
  min-width: 54px;
  padding: 8px 10px;
  background: #f8fafc;
  border-radius: 10px;
  text-align: center;
}

.mm-label {
  display: block;
  font-size: 12px;
  color: #6b7280;
}

.mm-val {
  display: block;
  margin-top: 2px;
  font-weight: 700;
  color: #111827;
}

.record-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.loading-placeholder,
.empty-state {
  padding: 48px 0;
  display: flex;
  justify-content: center;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.ai-loading-mask {
  position: fixed;
  inset: 0;
  background: rgba(245, 247, 250, 0.82);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(6px);
}

.ai-loading-content {
  background: #ffffff;
  padding: 32px 40px;
  border-radius: 20px;
  box-shadow: 0 16px 40px rgba(15, 23, 42, 0.1);
  text-align: center;
}

.ai-loading-content p {
  margin: 14px 0 0;
  color: #475569;
}

@media (max-width: 900px) {
  .content-wrapper {
    padding: 16px;
  }

  .page-header {
    align-items: flex-start;
    gap: 12px;
  }

  .record-card {
    grid-template-columns: 1fr;
  }

  .record-actions {
    justify-content: flex-start;
  }
}
</style>
