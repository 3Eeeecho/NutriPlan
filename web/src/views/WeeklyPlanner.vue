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
          <div class="week-range-controls">
            <n-button quaternary circle @click="goPrevWeek">
              <template #icon><n-icon><ChevronBackOutline /></n-icon></template>
            </n-button>
            <n-date-picker
              v-model:value="selectedWeekRange"
              type="daterange"
              :clearable="false"
              @update:value="handleWeekRangeChange"
            />
            <n-button quaternary circle @click="goNextWeek">
              <template #icon><n-icon><ChevronForwardOutline /></n-icon></template>
            </n-button>
          </div>
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
                <span class="meal-title" @dblclick.stop="openMealNoteEditor(day.date, meal.type)">{{ getMealTypeLabel(meal.type) }}</span>
                <div class="meal-title-actions">
                  <span class="meal-total">{{ meal.energy }} 千卡</span>
                  <n-dropdown
                    trigger="click"
                    :options="getMealActionOptions(day.date, meal.type)"
                    @select="(key) => handleMealAction(key, day.date, meal.type)"
                  >
                    <n-icon class="meal-action-icon"><EllipsisVerticalOutline /></n-icon>
                  </n-dropdown>
                </div>
              </div>

              <div v-if="isMealNoteEditing(day.date, meal.type)" class="meal-note-editor">
                <n-input
                  v-model:value="mealNoteDraftMap[getMealNoteKey(day.date, meal.type)]"
                  type="textarea"
                  :autosize="{ minRows: 2, maxRows: 6 }"
                  placeholder="输入该餐备注..."
                  class="meal-note-input"
                />
                <div class="meal-note-actions">
                  <button class="note-icon-btn" @click="deleteMealNoteInline(day.date, meal.type)">
                    <n-icon><CloseOutline /></n-icon>
                  </button>
                  <button class="note-icon-btn save" @click="saveMealNoteInline(day.date, meal.type)">
                    <n-icon><SaveOutline /></n-icon>
                  </button>
                </div>
              </div>
              <div
                v-else-if="hasMealNote(day.date, meal.type)"
                class="meal-note-display"
                @dblclick.stop="openMealNoteEditor(day.date, meal.type)"
              >
                {{ mealNoteMap[getMealNoteKey(day.date, meal.type)] }}
              </div>

              <div
                v-for="item in meal.items"
                :key="item.ID || item.id"
                class="meal-item-card"
                :class="{ 'recommend-card': item.isRecommendationRecipe }"
              >
                <div v-if="item.imageUrl" class="meal-item-image-wrap">
                  <img :src="item.imageUrl" alt="meal-image" class="meal-item-image" />
                </div>
                <div class="meal-item-name">{{ item.foodName || item.food_name }}</div>
                <div class="meal-item-meta">
                  <span>{{ Math.round(item.intakeAmount || item.intake_amount || 0) }} {{ item.intakeUnit || 'g' }}</span>
                  <span>{{ Math.round(item.calculatedEnergy || item.calculated_energy || 0) }} kcal</span>
                  <span>蛋 {{ Number(item.calculatedProtein || item.calculated_protein || 0).toFixed(1) }}</span>
                  <span>碳 {{ Number(item.calculatedCarb || item.calculated_carb || 0).toFixed(1) }}</span>
                  <span>脂 {{ Number(item.calculatedFat || item.calculated_fat || 0).toFixed(1) }}</span>
                  <span v-if="item.isRecommendationRecipe" class="recommend-tag">推荐食谱</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <n-modal v-model:show="addFoodModalVisible" preset="card" title="添加自定义食物" style="width: 460px; max-width: 92vw;">
        <n-form :model="addFoodForm" label-placement="left" label-width="88">
          <n-form-item label="食物名称">
            <n-input v-model:value="addFoodForm.food_name" placeholder="例如：无糖酸奶" />
          </n-form-item>
          <n-form-item label="摄入量(g)">
            <n-input-number v-model:value="addFoodForm.intake_amount" :step="10" :min="1" style="width: 100%;" />
          </n-form-item>
          <n-form-item label="热量(kcal)">
            <n-input-number v-model:value="addFoodForm.calculated_energy" :step="5" :min="0" style="width: 100%;" />
          </n-form-item>
          <n-form-item label="蛋白(g)">
            <n-input-number v-model:value="addFoodForm.calculated_protein" :step="0.5" :min="0" style="width: 100%;" />
          </n-form-item>
          <n-form-item label="碳水(g)">
            <n-input-number v-model:value="addFoodForm.calculated_carb" :step="0.5" :min="0" style="width: 100%;" />
          </n-form-item>
          <n-form-item label="脂肪(g)">
            <n-input-number v-model:value="addFoodForm.calculated_fat" :step="0.5" :min="0" style="width: 100%;" />
          </n-form-item>
        </n-form>
        <template #footer>
          <div class="modal-footer">
            <n-button @click="addFoodModalVisible = false">取消</n-button>
            <n-button type="primary" :loading="submitting" @click="handleSubmitCustomFood">保存</n-button>
          </div>
        </template>
      </n-modal>

      <n-modal v-model:show="deleteFoodModalVisible" preset="card" title="删除食物" style="width: 460px; max-width: 92vw;">
        <n-form label-placement="left" label-width="90">
          <n-form-item label="选择食物">
            <n-select
              v-model:value="selectedDeleteRecordId"
              :options="deleteFoodOptions"
              placeholder="请选择要删除的食物"
            />
          </n-form-item>
        </n-form>
        <template #footer>
          <div class="modal-footer">
            <n-button @click="deleteFoodModalVisible = false">取消</n-button>
            <n-popconfirm
              @positive-click="handleDeleteFood"
              positive-text="确认删除"
              negative-text="取消"
            >
              <template #trigger>
                <n-button type="error" :loading="deleting" :disabled="!selectedDeleteRecordId">删除</n-button>
              </template>
              确定删除该食物记录吗？
            </n-popconfirm>
          </div>
        </template>
      </n-modal>

      <n-modal v-model:show="deleteRecommendModalVisible" preset="card" title="删除推荐食谱" style="width: 420px; max-width: 90vw;">
        <div style="color: #475569; line-height: 1.7;">确定删除该餐次的推荐食谱吗？删除后将不再显示在本周计划中。</div>
        <template #footer>
          <div class="modal-footer">
            <n-button @click="deleteRecommendModalVisible = false">取消</n-button>
            <n-button type="error" @click="handleDeleteRecommendation">确认删除</n-button>
          </div>
        </template>
      </n-modal>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, onActivated } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NIcon, NSpin, NDropdown, NModal, NForm, NFormItem, NInput, NInputNumber, NSelect, NPopconfirm, NDatePicker, useMessage } from 'naive-ui'
import { ChevronBackOutline, ChevronForwardOutline, RefreshOutline, EllipsisVerticalOutline, CloseOutline, SaveOutline } from '@vicons/ionicons5'
import { getWeeklyReport, getTodayStatus, addIntakeRecord, deleteIntakeRecord } from '@/api/intakeApi'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const message = useMessage()
const authStore = useAuthStore()

const loading = ref(false)
const weeklyReport = ref({ daily_data: [] })
const recommendationCache = ref({ date: '', data: null })
const addFoodModalVisible = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const deleteFoodModalVisible = ref(false)
const deleteRecommendModalVisible = ref(false)
const selectedMealContext = ref({ date: '', mealType: '' })
const selectedDeleteRecordId = ref(null)
const deleteFoodOptions = ref([])
const mealNoteMap = ref({})
const mealNoteDraftMap = ref({})
const mealNoteEditingMap = ref({})
const selectedWeekRange = ref(null)
const RECOMMENDATION_CACHE_KEY = 'nutriplan_daily_recommendation'
const WEEKLY_NOTE_KEY_PREFIX = 'nutriplan_weekly_meal_note'

const addFoodForm = ref({
  food_name: '',
  intake_amount: 100,
  calculated_energy: 0,
  calculated_protein: 0,
  calculated_carb: 0,
  calculated_fat: 0
})

const getMealActionOptions = (date, mealType) => {
  const hasDeletableFood = getBackendMealRecords(date, mealType).length > 0
  const noteKey = getMealNoteKey(date, mealType)
  const hasNote = !!mealNoteMap.value[noteKey]
  const hasRecommendation = buildRecommendationRecordsForDate(date).some((item) => normalizeMealType(item.mealType || item.meal_type) === mealType)

  return [
    { label: '添加食物', key: 'add-custom-food' },
    { label: '删除食物', key: 'delete-food', disabled: !hasDeletableFood },
    { label: hasNote ? '修改笔记' : '添加笔记', key: 'add-note' },
    { label: '删除推荐食谱', key: 'delete-recommend', disabled: !hasRecommendation }
  ]
}

const getMealNoteKey = (date, mealType) => `${normalizeDateKey(date)}_${mealType}`

const hasMealNote = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  return !!mealNoteMap.value[key]
}

const isMealNoteEditing = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  return !!mealNoteEditingMap.value[key]
}

const getTodayKey = () => {
  const today = new Date()
  return `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
}

const getWeekStart = (date) => {
  const current = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  let weekday = current.getDay()
  if (weekday === 0) weekday = 7
  current.setDate(current.getDate() - (weekday - 1))
  return current
}

const formatApiDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const initWeekRange = () => {
  const start = getWeekStart(new Date())
  const end = new Date(start)
  end.setDate(start.getDate() + 6)
  selectedWeekRange.value = [start.getTime(), end.getTime()]
}

const normalizeDateKey = (value) => {
  if (!value) return ''
  const match = String(value).match(/^(\d{4})-(\d{1,2})-(\d{1,2})$/)
  if (!match) return String(value)
  const year = match[1]
  const month = String(match[2]).padStart(2, '0')
  const day = String(match[3]).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const getNoteStorageKey = () => `${WEEKLY_NOTE_KEY_PREFIX}_${authStore.user?.id || 'guest'}`

const loadMealNotes = () => {
  try {
    const raw = localStorage.getItem(getNoteStorageKey())
    mealNoteMap.value = raw ? JSON.parse(raw) : {}
    mealNoteDraftMap.value = { ...mealNoteMap.value }
    mealNoteEditingMap.value = {}
  } catch (error) {
    mealNoteMap.value = {}
    mealNoteDraftMap.value = {}
    mealNoteEditingMap.value = {}
  }
}

const saveMealNotes = () => {
  localStorage.setItem(getNoteStorageKey(), JSON.stringify(mealNoteMap.value || {}))
}

const loadRecommendationCache = () => {
  try {
    const raw = localStorage.getItem(RECOMMENDATION_CACHE_KEY)
    if (!raw) {
      recommendationCache.value = { date: '', data: null }
      return
    }
    const parsed = JSON.parse(raw)
    recommendationCache.value = {
      date: normalizeDateKey(parsed?.date || ''),
      data: parsed?.data || null
    }
  } catch (error) {
    recommendationCache.value = { date: '', data: null }
  }
}

const buildWeeklyReportByDailyStatus = async (startDate, endDate) => {
  const dailyData = []
  const cursor = new Date(startDate.getFullYear(), startDate.getMonth(), startDate.getDate())
  const endKey = formatApiDate(endDate)

  while (formatApiDate(cursor) <= endKey) {
    const dateKey = formatApiDate(cursor)
    try {
      const status = await getTodayStatus(dateKey)
      dailyData.push({
        date: dateKey,
        total_energy: Number(status?.total_energy ?? status?.totalEnergy ?? 0),
        total_protein: Number(status?.total_protein ?? status?.totalProtein ?? 0),
        total_carbohydrate: Number(status?.total_carbohydrate ?? status?.totalCarbohydrate ?? 0),
        total_fat: Number(status?.total_fat ?? status?.totalFat ?? 0),
        records: status?.records || []
      })
    } catch (error) {
      dailyData.push({
        date: dateKey,
        total_energy: 0,
        total_protein: 0,
        total_carbohydrate: 0,
        total_fat: 0,
        records: []
      })
    }
    cursor.setDate(cursor.getDate() + 1)
  }

  return {
    start_date: formatApiDate(startDate),
    end_date: formatApiDate(endDate),
    daily_data: dailyData
  }
}

const buildRecommendationRecordsForDate = (dateKey) => {
  if (normalizeDateKey(recommendationCache.value.date) !== normalizeDateKey(dateKey) || !recommendationCache.value.data) return []

  const source = recommendationCache.value.data
  const mealOrder = ['breakfast', 'lunch', 'dinner', 'snack']

  return mealOrder
    .filter((type) => !!source[type])
    .map((type) => {
      const meal = source[type]
      return {
        id: `recommend-${dateKey}-${type}-${meal.id || meal.name || 'unknown'}`,
        mealType: type,
        foodName: meal.name || '推荐食谱',
        intakeAmount: 1,
        intakeUnit: '份',
        calculatedEnergy: Number(meal.energy || meal.calories || 0),
        calculatedProtein: Number(meal.protein || meal.total_protein || 0),
        calculatedCarb: Number(meal.carbohydrate || meal.carbs || meal.total_carbohydrate || 0),
        calculatedFat: Number(meal.fat || meal.total_fat || 0),
        imageUrl: meal.image_url || meal.imageUrl || '',
        isRecommendationRecipe: true
      }
    })
}

const weekColumns = computed(() => {
  const list = weeklyReport.value?.daily_data || []
  const todayKey = getTodayKey()

  return list.map((day) => {
    const date = day.date
    const dateObj = new Date(date)
    const backendRecords = Array.isArray(day.records) ? day.records : []
    const recommendationRecords = buildRecommendationRecordsForDate(date)
    const records = [...recommendationRecords, ...backendRecords]
    const recommendationEnergy = recommendationRecords.reduce((sum, item) => sum + Number(item.calculatedEnergy || 0), 0)

    return {
      date,
      isToday: date === todayKey,
      displayDate: `${dateObj.getMonth() + 1}月${dateObj.getDate()}日`,
      weekLabel: dateObj.toLocaleDateString('zh-CN', { weekday: 'long' }),
      totalEnergy: Math.round((day.total_energy || day.totalEnergy || 0) + recommendationEnergy),
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
      const items = [...grouped[type]].sort((first, second) => {
        const firstWeight = first.isRecommendationRecipe ? -1 : 0
        const secondWeight = second.isRecommendationRecipe ? -1 : 0
        return firstWeight - secondWeight
      })
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
    if (!Array.isArray(selectedWeekRange.value) || selectedWeekRange.value.length !== 2) {
      initWeekRange()
    }

    const [startTimestamp, endTimestamp] = selectedWeekRange.value
    const startDate = new Date(startTimestamp)
    const endDate = new Date(endTimestamp)

    if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
      initWeekRange()
    }

    const [safeStartTimestamp, safeEndTimestamp] = selectedWeekRange.value
    const safeStartDate = new Date(safeStartTimestamp)
    const safeEndDate = new Date(safeEndTimestamp)

    loadRecommendationCache()
    loadMealNotes()

    const expectedStart = formatApiDate(safeStartDate)
    const expectedEnd = formatApiDate(safeEndDate)

    let res = await getWeeklyReport({
      start_date: formatApiDate(safeStartDate),
      end_date: formatApiDate(safeEndDate)
    })

    const actualStart = normalizeDateKey(res?.start_date || res?.startDate || '')
    const actualEnd = normalizeDateKey(res?.end_date || res?.endDate || '')
    const isRangeMatched = actualStart === expectedStart && actualEnd === expectedEnd

    if (!isRangeMatched) {
      res = await buildWeeklyReportByDailyStatus(safeStartDate, safeEndDate)
    }

    weeklyReport.value = res || { daily_data: [] }
  } catch (error) {
    message.error('加载周食谱失败')
  } finally {
    loading.value = false
  }
}

const handleWeekRangeChange = async (value) => {
  if (!Array.isArray(value) || value.length !== 2) {
    initWeekRange()
    await loadWeeklyData()
    return
  }

  const [startTimestamp, endTimestamp] = value
  const startDate = new Date(startTimestamp)
  const endDate = new Date(endTimestamp)

  if (Number.isNaN(startDate.getTime()) || Number.isNaN(endDate.getTime())) {
    initWeekRange()
    await loadWeeklyData()
    return
  }

  const diffDays = Math.round((endDate.getTime() - startDate.getTime()) / (24 * 60 * 60 * 1000))

  if (diffDays !== 6) {
    const fixedEnd = new Date(startDate)
    fixedEnd.setDate(startDate.getDate() + 6)
    selectedWeekRange.value = [startDate.getTime(), fixedEnd.getTime()]
    message.warning('周视图固定为连续 7 天，已自动调整')
  } else {
    selectedWeekRange.value = value
  }

  await loadWeeklyData()
}

const shiftWeek = async (step) => {
  if (!Array.isArray(selectedWeekRange.value) || selectedWeekRange.value.length !== 2) {
    initWeekRange()
  }

  const [startTimestamp, endTimestamp] = selectedWeekRange.value
  const delta = step * 7 * 24 * 60 * 60 * 1000
  selectedWeekRange.value = [startTimestamp + delta, endTimestamp + delta]
  await loadWeeklyData()
}

const goPrevWeek = async () => {
  await shiftWeek(-1)
}

const goNextWeek = async () => {
  await shiftWeek(1)
}

const openCustomFoodModal = (date, mealType) => {
  if (date !== getTodayKey()) {
    message.warning('当前仅支持向今天添加自定义食物')
    return
  }

  selectedMealContext.value = { date, mealType }
  addFoodForm.value = {
    food_name: '',
    intake_amount: 100,
    calculated_energy: 0,
    calculated_protein: 0,
    calculated_carb: 0,
    calculated_fat: 0
  }
  addFoodModalVisible.value = true
}

const getBackendMealRecords = (date, mealType) => {
  const dailyData = weeklyReport.value?.daily_data || []
  const targetDay = dailyData.find((item) => normalizeDateKey(item.date) === normalizeDateKey(date))
  const records = Array.isArray(targetDay?.records) ? targetDay.records : []

  return records
    .filter((record) => normalizeMealType(record.mealType || record.meal_type) === mealType)
    .map((record) => ({
      id: Number(record.ID || record.id),
      name: record.foodName || record.food_name || `记录 #${record.ID || record.id}`,
      energy: Math.round(Number(record.calculatedEnergy || record.calculated_energy || 0))
    }))
    .filter((record) => Number.isFinite(record.id) && record.id > 0)
}

const openDeleteFoodModal = (date, mealType) => {
  const records = getBackendMealRecords(date, mealType)
  if (records.length === 0) {
    message.info('该餐次暂无可删除的食物')
    return
  }

  deleteFoodOptions.value = records.map((item) => ({
    label: `${item.name}（${item.energy} kcal）`,
    value: item.id
  }))
  selectedDeleteRecordId.value = deleteFoodOptions.value[0]?.value || null
  selectedMealContext.value = { date, mealType }
  deleteFoodModalVisible.value = true
}

const openMealNoteEditor = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  mealNoteDraftMap.value = {
    ...mealNoteDraftMap.value,
    [key]: mealNoteMap.value[key] || ''
  }
  mealNoteEditingMap.value = {
    ...mealNoteEditingMap.value,
    [key]: true
  }
}

const cancelMealNoteEdit = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  mealNoteDraftMap.value = {
    ...mealNoteDraftMap.value,
    [key]: mealNoteMap.value[key] || ''
  }

  const next = { ...mealNoteEditingMap.value }
  delete next[key]
  mealNoteEditingMap.value = next
}

const saveMealNoteInline = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  const content = String(mealNoteDraftMap.value[key] || '').trim()

  if (content) {
    mealNoteMap.value = { ...mealNoteMap.value, [key]: content }
  } else {
    const next = { ...mealNoteMap.value }
    delete next[key]
    mealNoteMap.value = next
  }

  saveMealNotes()
  cancelMealNoteEdit(date, mealType)
  message.success('备注已保存')
}

const deleteMealNoteInline = (date, mealType) => {
  const key = getMealNoteKey(date, mealType)
  const next = { ...mealNoteMap.value }
  delete next[key]
  mealNoteMap.value = next

  const nextDraft = { ...mealNoteDraftMap.value }
  delete nextDraft[key]
  mealNoteDraftMap.value = nextDraft

  const nextEditing = { ...mealNoteEditingMap.value }
  delete nextEditing[key]
  mealNoteEditingMap.value = nextEditing

  saveMealNotes()
  message.success('备注已删除')
}

const openDeleteRecommendModal = (date, mealType) => {
  selectedMealContext.value = { date, mealType }
  deleteRecommendModalVisible.value = true
}

const handleDeleteRecommendation = async () => {
  const { date, mealType } = selectedMealContext.value
  if (!date || !mealType) return

  if (normalizeDateKey(recommendationCache.value.date) !== normalizeDateKey(date) || !recommendationCache.value.data) {
    deleteRecommendModalVisible.value = false
    message.info('当前日期没有可删除的推荐食谱')
    return
  }

  const nextData = { ...(recommendationCache.value.data || {}) }
  delete nextData[mealType]

  recommendationCache.value = {
    date: recommendationCache.value.date,
    data: nextData
  }

  localStorage.setItem(RECOMMENDATION_CACHE_KEY, JSON.stringify({
    date: recommendationCache.value.date,
    data: nextData
  }))

  deleteRecommendModalVisible.value = false
  await loadWeeklyData()
  message.success('已删除该餐次推荐食谱')
}

const handleDeleteFood = async () => {
  if (!selectedDeleteRecordId.value) {
    message.warning('请选择要删除的食物')
    return
  }

  deleting.value = true
  try {
    await deleteIntakeRecord(selectedDeleteRecordId.value)
    message.success('已删除食物')
    deleteFoodModalVisible.value = false
    await loadWeeklyData()
  } catch (error) {
    message.error('删除失败，请稍后重试')
  } finally {
    deleting.value = false
  }
}

const handleMealAction = (actionKey, date, mealType) => {
  if (actionKey === 'add-custom-food') {
    openCustomFoodModal(date, mealType)
    return
  }

  if (actionKey === 'delete-food') {
    openDeleteFoodModal(date, mealType)
    return
  }

  if (actionKey === 'add-note') {
    openMealNoteEditor(date, mealType)
    return
  }

  if (actionKey === 'delete-recommend') {
    openDeleteRecommendModal(date, mealType)
  }
}

const handleSubmitCustomFood = async () => {
  if (!addFoodForm.value.food_name) {
    message.warning('请先填写食物名称')
    return
  }

  submitting.value = true
  try {
    await addIntakeRecord({
      meal_type: selectedMealContext.value.mealType,
      food_name: addFoodForm.value.food_name,
      intake_amount: Number(addFoodForm.value.intake_amount || 0),
      calculated_energy: Number(addFoodForm.value.calculated_energy || 0),
      calculated_protein: Number(addFoodForm.value.calculated_protein || 0),
      calculated_carb: Number(addFoodForm.value.calculated_carb || 0),
      calculated_fat: Number(addFoodForm.value.calculated_fat || 0)
    })

    message.success('已添加到今日饮食记录')
    addFoodModalVisible.value = false
    await loadWeeklyData()
  } catch (error) {
    message.error('添加失败，请稍后重试')
  } finally {
    submitting.value = false
  }
}

const handleWindowFocus = () => {
  loadWeeklyData()
}

onMounted(() => {
  initWeekRange()
  loadWeeklyData()
  window.addEventListener('focus', handleWindowFocus)
})

onActivated(() => {
  if (!Array.isArray(selectedWeekRange.value) || selectedWeekRange.value.length !== 2) {
    initWeekRange()
  }
  loadWeeklyData()
})

onUnmounted(() => {
  window.removeEventListener('focus', handleWindowFocus)
})
</script>

<style scoped>
.weekly-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 0;
}

.weekly-container {
  max-width: none;
  width: 100%;
  margin: 0;
  padding: 18px;
  background: #141923;
  border-radius: 0;
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

.right-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.week-range-controls {
  display: flex;
  align-items: center;
  gap: 8px;
}

.week-range-controls :deep(.n-date-picker) {
  width: 250px;
}

.week-range-controls :deep(.n-input__input-el) {
  color: #e2e8f0;
}

.week-range-controls :deep(.n-input) {
  background: #0f172a;
  border-color: #334155;
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

.meal-title-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meal-title {
  color: #f3f4f6;
  font-weight: 700;
  font-size: 34px;
  cursor: pointer;
}

.meal-total {
  color: #94a3b8;
  font-size: 12px;
}

.meal-note {
  margin-bottom: 8px;
  font-size: 12px;
  color: #93c5fd;
  border: 1px dashed rgba(147, 197, 253, 0.35);
  border-radius: 6px;
  padding: 4px 8px;
}

.meal-note-editor {
  margin-bottom: 8px;
  border: 1px solid #334155;
  border-radius: 8px;
  background: #151c28;
  padding: 8px;
}

.meal-note-display {
  margin-bottom: 8px;
  border: 1px solid rgba(96, 165, 250, 0.35);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.75);
  color: #c7d2fe;
  padding: 8px 10px;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  cursor: pointer;
}

.meal-note-input :deep(textarea) {
  background: transparent;
  color: #e2e8f0;
}

.meal-note-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: 8px;
}

.note-icon-btn {
  border: none;
  background: transparent;
  color: #cbd5e1;
  font-size: 18px;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
}

.note-icon-btn:hover {
  background: rgba(148, 163, 184, 0.15);
}

.note-icon-btn.save {
  color: #60a5fa;
}

.meal-action-icon {
  font-size: 18px;
  color: #94a3b8;
  cursor: pointer;
  border-radius: 6px;
  padding: 2px;
  transition: all 0.2s ease;
}

.meal-action-icon:hover {
  background: rgba(148, 163, 184, 0.14);
  color: #e2e8f0;
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

.recommend-card {
  border-color: #2c695a;
  background: #0f1f1c;
}

.meal-item-image-wrap {
  width: 100%;
  height: 92px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 8px;
  background: #0b1220;
}

.meal-item-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.recommend-tag {
  color: #60a5fa;
  border: 1px solid rgba(96, 165, 250, 0.4);
  padding: 0 6px;
  border-radius: 999px;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
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
