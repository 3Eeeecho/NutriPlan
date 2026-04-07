<template>
  <div class="nutrition-dashboard">
    <header class="top-header glass-card">
      <div class="brand-wrap">
        <h1 class="brand">基于用户健康数据与营养模型的饮食推荐系统</h1>
      </div>
    </header>

    <main class="dashboard-grid">
      <section class="center-panel">
        <p v-if="加载中" class="loading-tip">正在同步饮食数据...</p>

        <div class="record-ops glass-card">
          <div class="ops-left">
            <button class="icon-circle soft" aria-label="上一天记录" @click="上一天"><n-icon><ChevronBackOutline /></n-icon></button>
            <button class="icon-circle soft" aria-label="下一天记录" :disabled="是否今天" @click="下一天"><n-icon><ChevronForwardOutline /></n-icon></button>
            <p class="ops-date">{{ 日期标题 }}</p>
          </div>
          <label class="calendar-picker">
            <n-icon><CalendarOutline /></n-icon>
            <input type="date" :value="日期输入值" :max="今日日期输入值" @change="选择日期" />
          </label>
        </div>

        <div class="timeline">
          <article
            v-for="meal in 时间线餐次"
            :key="meal.type"
            class="timeline-row"
            draggable="true"
            @dragstart="开始拖拽餐次(meal.type)"
            @dragover.prevent
            @drop="放置餐次(meal.type)"
            @dragend="结束拖拽餐次"
          >
            <div class="time-col">
              <span class="dot" />
              <div class="time-editor" @click.stop>
                <button v-if="编辑时间餐次 !== meal.type" class="time-display-btn" @click="开始编辑时间(meal.type)">{{ meal.time }}</button>
                <input
                  v-else
                  type="time"
                  class="time-input"
                  :value="meal.time"
                  @input="(event) => 更新餐次时间(meal.type, event.target.value)"
                  @blur="结束编辑时间"
                  @keyup.enter="结束编辑时间"
                />
              </div>
            </div>

            <div class="timeline-content glass-card meal-card" :class="`meal-${meal.type}`">
              <div class="meal-head">
                <div class="meal-title-wrap">
                  <input
                    v-if="编辑名称餐次 === meal.type"
                    v-model="编辑名称内容"
                    class="meal-name-input"
                    maxlength="16"
                    @blur="保存名称编辑(meal.type)"
                    @keyup.enter="保存名称编辑(meal.type)"
                  />
                  <h3 v-else>{{ meal.name }}</h3>
                  <button
                    v-if="编辑名称餐次 !== meal.type"
                    class="rename-link"
                    aria-label="编辑餐次名称"
                    title="改名"
                    @click.stop="开始编辑名称(meal.type, meal.name)"
                  >
                    改名
                  </button>
                </div>
                <p class="kcal">{{ meal.kcal || '--' }}</p>
              </div>

              <div class="meal-hero">
                <div class="hero-badge">{{ meal.icon }}</div>
                <div class="hero-illustration">{{ meal.hero }}</div>
              </div>

              <div v-if="meal.type === 'snack'" class="snack-group">
                <article v-for="snack in snackCards" :key="snack.label" class="snack-mini-card">
                  <span class="mini-emoji">{{ snack.emoji }}</span>
                  <span class="mini-label">{{ snack.label }}</span>
                  <button class="mini-plus" aria-label="添加">+</button>
                </article>
              </div>
              <div v-else class="food-hints">
                <span v-for="hint in meal.hints" :key="hint" class="food-chip">{{ hint }}</span>
              </div>

              <div class="bar-track">
                <div class="bar-fill" :style="{ width: meal.progress + '%' }" />
              </div>

              <div class="meal-actions">
                <button class="record-btn-main" @click="去记录餐次">
                  <n-icon><AddCircleOutline /></n-icon>
                  <span>{{ meal.type === 'snack' ? '记录餐次' : '记录今日餐次' }}</span>
                </button>
                <button class="detail-link" @click.stop="切换详情餐次(meal.type)">查看详情</button>
              </div>

              <div v-if="meal.variant === 'logged'" class="footer-tip">状态：已记录</div>
            </div>
          </article>
        </div>
      </section>

      <section class="right-panel">
        <div class="overview-card glass-card">
          <span class="floating-red top-right"><n-icon><SparklesOutline /></n-icon></span>
          <div class="ring-wrap">
            <div class="calorie-ring" :style="ringStyle">
              <div class="ring-center">
                <h3>{{ 剩余热量 }}</h3>
                <p>剩余千卡</p>
              </div>
            </div>
            <p class="goal-text">每日目标：{{ 每日目标热量 }} 千卡</p>
          </div>

          <div class="macro-stack">
            <article v-for="macro in 宏量卡片" :key="macro.label" class="macro-card">
              <div class="macro-left">
                <span class="macro-icon" :style="{ background: macro.soft }">{{ macro.icon }}</span>
                <div>
                  <p class="macro-label">{{ macro.label }}</p>
                  <p class="macro-target">目标：{{ macro.target }}</p>
                </div>
              </div>
              <p class="macro-value">{{ macro.value }}</p>
            </article>
          </div>
        </div>

        <div class="details-card glass-card">
          <div class="details-head">
            <h3>{{ 当前详情标题 }}</h3>
            <p>{{ 当前详情时间 }}</p>
          </div>

          <p v-if="!当前详情食物.length" class="sample-note">今日暂无记录，以下为示例食谱</p>

          <div class="micro-grid">
            <div v-for="line in 详情营养条" :key="line.label" class="micro-item">
              <div class="micro-top">
                <span>{{ line.label }}</span>
                <span>{{ line.value }}</span>
              </div>
              <div class="micro-track">
                <div class="micro-fill" :style="{ width: line.progress + '%', background: line.color }" />
              </div>
            </div>
          </div>

          <div class="food-table">
            <div class="table-head">
              <span>食物项</span>
              <span>份量</span>
              <span>热量</span>
            </div>

            <div v-for="food in 展示详情食物" :key="food.id" class="food-row">
              <div class="food-name">
                <span class="thumb" :style="{ background: food.bg }">{{ food.emoji }}</span>
                <span>{{ food.name }}</span>
                <span v-if="food.sample" class="sample-tag">示例</span>
              </div>
              <span>{{ food.portion }}</span>
              <span class="strong-text">{{ food.kcal }}</span>
            </div>
          </div>

          <div class="add-more">
            <input type="text" placeholder="为这餐添加更多食物..." />
            <button class="plus-btn" aria-label="添加食物">+</button>
          </div>
        </div>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { NIcon } from 'naive-ui'
import {
  AddCircleOutline,
  CalendarOutline,
  ChevronBackOutline,
  ChevronForwardOutline,
  SparklesOutline
} from '@vicons/ionicons5'
import { getTodayStatus } from '@/api/intakeApi'
import { getNutritionRequirements } from '@/api/user'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

const 加载中 = ref(false)
const 已选日期 = ref(new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate()))
const 编辑时间餐次 = ref('')
const 编辑名称餐次 = ref('')
const 编辑名称内容 = ref('')
const 拖拽源餐次 = ref('')
const 每日状态 = ref({
  records: [],
  totalEnergy: 0,
  totalProtein: 0,
  totalCarbohydrate: 0,
  totalFat: 0,
  targetEnergy: 2500,
  targetProtein: 150,
  targetCarbohydrate: 220,
  targetFat: 60
})
const 详情餐次类型 = ref('breakfast')

const 餐次时间本地键 = 'nutriplan_home_meal_times'
const 餐次顺序本地键 = 'nutriplan_home_meal_order'
const 餐次名称本地键 = 'nutriplan_home_meal_names'
const 默认餐次时间 = {
  breakfast: '08:15',
  lunch: '12:30',
  snack: '16:00',
  dinner: '19:30'
}
const 餐次时间映射 = ref({ ...默认餐次时间 })

const 餐次配置 = [
  { type: 'breakfast', name: '早餐', time: '08:15', icon: '🍳', hero: '🍳🍞', hints: ['🍓', '🥣', '🥛'], ratio: 0.16, emoji: '🍳', bg: 'linear-gradient(135deg, #ffdb96, #ffb56f)' },
  { type: 'lunch', name: '午餐', time: '12:30', icon: '🥗', hero: '🥗🍽️', hints: ['🥦', '🍗', '🍚'], ratio: 0.32, emoji: '🥗', bg: 'linear-gradient(135deg, #bdeecb, #88dbc8)' },
  { type: 'snack', name: '下午加餐', time: '16:00', icon: '🫐', hero: '🥜🍎', hints: ['🥜', '🍓', '🍵'], ratio: 0.2, emoji: '🍎', bg: 'linear-gradient(135deg, #ffe7c8, #ffd6ad)' },
  { type: 'dinner', name: '晚餐', time: '19:30', icon: '🥩', hero: '🥩🍷', hints: ['🥔', '🥬', '🍷'], ratio: 0.32, emoji: '🍽️', bg: 'linear-gradient(135deg, #d9ecff, #bcdcff)' }
]
const snackCards = [
  { label: '坚果/种子', emoji: '🥜' },
  { label: '水果', emoji: '🍓' },
  { label: '茶点', emoji: '🍵' }
]

const 默认餐次顺序 = 餐次配置.map((item) => item.type)
const 默认餐次名称 = Object.fromEntries(餐次配置.map((item) => [item.type, item.name]))
const 餐次顺序 = ref([...默认餐次顺序])
const 餐次名称映射 = ref({ ...默认餐次名称 })

const 安全数值 = (value) => {
  const number = Number(value)
  return Number.isFinite(number) ? number : 0
}

const 转日期参数 = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const 格式化时间 = (value) => {
  if (!value) return '记录于 --:--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '记录于 --:--'
  const hour = String(date.getHours()).padStart(2, '0')
  const minute = String(date.getMinutes()).padStart(2, '0')
  return `记录于 ${hour}:${minute}`
}

const 归一化餐次 = (value) => {
  if (value === 1 || value === '1' || value === '早餐' || value === 'breakfast') return 'breakfast'
  if (value === 2 || value === '2' || value === '午餐' || value === 'lunch') return 'lunch'
  if (value === 3 || value === '3' || value === '晚餐' || value === 'dinner') return 'dinner'
  if (value === 4 || value === '4' || value === '加餐' || value === '下午加餐' || value === 'snack') return 'snack'
  return ''
}

const 获取餐次记录 = (type) => {
  return (每日状态.value.records || []).filter((item) => 归一化餐次(item.mealType || item.meal_type) === type)
}

const 加载餐次时间配置 = () => {
  try {
    const raw = localStorage.getItem(餐次时间本地键)
    if (!raw) return
    const parsed = JSON.parse(raw)
    餐次时间映射.value = { ...默认餐次时间, ...parsed }
  } catch (error) {
    console.warn('读取餐次时间配置失败', error)
  }
}

const 保存餐次时间配置 = () => {
  localStorage.setItem(餐次时间本地键, JSON.stringify(餐次时间映射.value))
}

const 加载餐次个性化配置 = () => {
  try {
    const orderRaw = localStorage.getItem(餐次顺序本地键)
    if (orderRaw) {
      const parsedOrder = JSON.parse(orderRaw)
      const validOrder = parsedOrder.filter((type) => 默认餐次顺序.includes(type))
      const missing = 默认餐次顺序.filter((type) => !validOrder.includes(type))
      餐次顺序.value = [...validOrder, ...missing]
    }

    const nameRaw = localStorage.getItem(餐次名称本地键)
    if (nameRaw) {
      const parsedName = JSON.parse(nameRaw)
      餐次名称映射.value = { ...默认餐次名称, ...parsedName }
    }
  } catch (error) {
    console.warn('读取餐次个性化配置失败', error)
  }
}

const 保存餐次顺序配置 = () => {
  localStorage.setItem(餐次顺序本地键, JSON.stringify(餐次顺序.value))
}

const 保存餐次名称配置 = () => {
  localStorage.setItem(餐次名称本地键, JSON.stringify(餐次名称映射.value))
}

const 开始编辑时间 = (mealType) => {
  编辑时间餐次.value = mealType
}

const 更新餐次时间 = (mealType, value) => {
  if (!value) return
  餐次时间映射.value = {
    ...餐次时间映射.value,
    [mealType]: value
  }
}

const 结束编辑时间 = () => {
  if (!编辑时间餐次.value) return
  保存餐次时间配置()
  编辑时间餐次.value = ''
}

const 开始编辑名称 = (mealType, currentName) => {
  编辑名称餐次.value = mealType
  编辑名称内容.value = currentName || 餐次名称映射.value[mealType] || ''
}

const 保存名称编辑 = (mealType) => {
  const nextName = 编辑名称内容.value.trim()
  if (!nextName) {
    编辑名称餐次.value = ''
    编辑名称内容.value = ''
    return
  }
  餐次名称映射.value = {
    ...餐次名称映射.value,
    [mealType]: nextName
  }
  保存餐次名称配置()
  编辑名称餐次.value = ''
  编辑名称内容.value = ''
}

const 开始拖拽餐次 = (mealType) => {
  拖拽源餐次.value = mealType
}

const 放置餐次 = (targetType) => {
  const sourceType = 拖拽源餐次.value
  if (!sourceType || sourceType === targetType) return

  const nextOrder = [...餐次顺序.value]
  const sourceIndex = nextOrder.indexOf(sourceType)
  const targetIndex = nextOrder.indexOf(targetType)
  if (sourceIndex < 0 || targetIndex < 0) return

  nextOrder.splice(sourceIndex, 1)
  nextOrder.splice(targetIndex, 0, sourceType)
  餐次顺序.value = nextOrder
  保存餐次顺序配置()
}

const 结束拖拽餐次 = () => {
  拖拽源餐次.value = ''
}

const 从日期字符串创建日期 = (value) => {
  if (!value) return null
  const [year, month, day] = value.split('-').map((item) => Number(item))
  if (!year || !month || !day) return null
  const date = new Date(year, month - 1, day)
  return Number.isNaN(date.getTime()) ? null : date
}

const 今日日期输入值 = computed(() => 转日期参数(new Date()))
const 日期输入值 = computed(() => 转日期参数(已选日期.value))

const 选择日期 = async (event) => {
  const nextDate = 从日期字符串创建日期(event?.target?.value)
  if (!nextDate) return
  已选日期.value = nextDate
  await 同步看板数据()
}

const 日期标题 = computed(() => {
  const now = new Date()
  const isToday = now.toDateString() === 已选日期.value.toDateString()
  const month = 已选日期.value.getMonth() + 1
  const day = 已选日期.value.getDate()
  return isToday ? `今天, ${month}月${day}日` : `${month}月${day}日`
})

const 是否今天 = computed(() => new Date().toDateString() === 已选日期.value.toDateString())

const 每日目标热量 = computed(() => Math.round(安全数值(每日状态.value.targetEnergy)))
const 已摄入热量 = computed(() => Math.round(安全数值(每日状态.value.totalEnergy)))
const 剩余热量 = computed(() => Math.max(0, 每日目标热量.value - 已摄入热量.value).toLocaleString('zh-CN'))

const 时间线餐次 = computed(() => {
  const 配置映射 = Object.fromEntries(餐次配置.map((item) => [item.type, item]))
  return 餐次顺序.value.map((type) => {
    const 基础配置 = 配置映射[type]
    if (!基础配置) return null

    const config = {
      ...基础配置,
      name: 餐次名称映射.value[type] || 基础配置.name
    }
    const records = 获取餐次记录(config.type)
    const consumed = Math.round(records.reduce((sum, item) => sum + 安全数值(item.calculatedEnergy || item.calculated_energy), 0))
    const target = Math.max(1, Math.round(每日目标热量.value * config.ratio))
    const progress = Math.min(100, Math.round((consumed / target) * 100))
    return {
      ...config,
      time: 餐次时间映射.value[config.type] || config.time,
      records,
      consumed,
      target,
      kcal: consumed > 0 ? `${consumed} / ${target} 千卡` : '',
      progress,
      variant: consumed > 0 ? (config.type === 'lunch' ? 'logged' : 'filled') : 'empty'
    }
  }).filter(Boolean)
})

const 宏量卡片 = computed(() => {
  const totalCarbs = Math.round(安全数值(每日状态.value.totalCarbohydrate))
  const totalProtein = Math.round(安全数值(每日状态.value.totalProtein))
  const totalFat = Math.round(安全数值(每日状态.value.totalFat))

  const targetCarbs = Math.round(安全数值(每日状态.value.targetCarbohydrate))
  const targetProtein = Math.round(安全数值(每日状态.value.targetProtein))
  const targetFat = Math.round(安全数值(每日状态.value.targetFat))

  return [
    { label: '碳水', value: `${totalCarbs}g`, target: `${targetCarbs}g`, icon: '🌾', soft: 'rgba(255, 186, 108, 0.22)' },
    { label: '蛋白质', value: `${totalProtein}g`, target: `${targetProtein}g`, icon: '💧', soft: 'rgba(116, 163, 255, 0.2)' },
    { label: '脂肪', value: `${totalFat}g`, target: `${targetFat}g`, icon: '🟡', soft: 'rgba(255, 215, 83, 0.22)' }
  ]
})

const 当前详情餐次 = computed(() => {
  const current = 时间线餐次.value.find((item) => item.type === 详情餐次类型.value)
  if (current) return current
  return 时间线餐次.value.find((item) => item.records.length > 0) || 时间线餐次.value[1] || 时间线餐次.value[0]
})

const 当前详情标题 = computed(() => `${当前详情餐次.value?.name || '餐次'}详情`)

const 当前详情时间 = computed(() => {
  const record = 当前详情餐次.value?.records?.[0]
  return 格式化时间(record?.createdAt || record?.CreatedAt)
})

const 详情营养条 = computed(() => {
  const records = 当前详情餐次.value?.records || []
  const ratio = 当前详情餐次.value?.ratio || 0.25
  const protein = records.reduce((sum, item) => sum + 安全数值(item.calculatedProtein || item.calculated_protein), 0)
  const carbs = records.reduce((sum, item) => sum + 安全数值(item.calculatedCarb || item.calculated_carb), 0)
  const fat = records.reduce((sum, item) => sum + 安全数值(item.calculatedFat || item.calculated_fat), 0)

  const targetProtein = Math.max(1, 安全数值(每日状态.value.targetProtein) * ratio)
  const targetCarbs = Math.max(1, 安全数值(每日状态.value.targetCarbohydrate) * ratio)
  const targetFat = Math.max(1, 安全数值(每日状态.value.targetFat) * ratio)

  return [
    { label: '蛋白质', value: `${Math.round(protein)}g`, progress: Math.min(100, Math.round((protein / targetProtein) * 100)), color: '#4e8fff' },
    { label: '碳水', value: `${Math.round(carbs)}g`, progress: Math.min(100, Math.round((carbs / targetCarbs) * 100)), color: '#ff8d2f' },
    { label: '脂肪', value: `${Math.round(fat)}g`, progress: Math.min(100, Math.round((fat / targetFat) * 100)), color: '#f5c84b' }
  ]
})

const 当前详情食物 = computed(() => {
  const meal = 当前详情餐次.value
  const fallback = { emoji: meal?.emoji || '🍽️', bg: meal?.bg || 'linear-gradient(135deg, #7aa694, #4f7a69)' }
  return (meal?.records || []).map((item, index) => ({
    id: item.id || item.ID || `${meal.type}-${index}`,
    name: item.foodName || item.food_name || '未知食物',
    portion: `${Math.round(安全数值(item.intakeAmount || item.intake_amount))}g`,
    kcal: `${Math.round(安全数值(item.calculatedEnergy || item.calculated_energy))} 千卡`,
    emoji: fallback.emoji,
    bg: fallback.bg,
    sample: false
  }))
})

const 示例食谱映射 = {
  breakfast: [
    { name: '牛油果全麦吐司', portion: '130g', kcal: '285 千卡', emoji: '🥑' },
    { name: '蓝莓酸奶碗', portion: '180g', kcal: '210 千卡', emoji: '🫐' }
  ],
  lunch: [
    { name: '鸡胸肉藜麦沙拉', portion: '260g', kcal: '420 千卡', emoji: '🥗' },
    { name: '南瓜浓汤', portion: '220g', kcal: '165 千卡', emoji: '🎃' }
  ],
  snack: [
    { name: '混合坚果', portion: '35g', kcal: '195 千卡', emoji: '🥜' },
    { name: '苹果切片', portion: '120g', kcal: '62 千卡', emoji: '🍎' }
  ],
  dinner: [
    { name: '香煎三文鱼配芦笋', portion: '240g', kcal: '468 千卡', emoji: '🐟' },
    { name: '黑椒菌菇意面', portion: '210g', kcal: '338 千卡', emoji: '🍝' }
  ]
}

const 示例详情食物 = computed(() => {
  const meal = 当前详情餐次.value
  const list = 示例食谱映射[meal?.type] || []
  return list.map((item, index) => ({
    id: `sample-${meal?.type || 'meal'}-${index}`,
    name: item.name,
    portion: item.portion,
    kcal: item.kcal,
    emoji: item.emoji,
    bg: meal?.bg || 'linear-gradient(135deg, #7aa694, #4f7a69)',
    sample: true
  }))
})

const 展示详情食物 = computed(() => {
  return 当前详情食物.value.length ? 当前详情食物.value : 示例详情食物.value
})

const ringStyle = computed(() => {
  const goal = Math.max(1, 每日目标热量.value)
  const progress = Math.min(100, Math.round((已摄入热量.value / goal) * 100))
  return {
    background: `conic-gradient(#2f7c63 0 ${progress}%, #dbe8e1 ${progress}% 100%)`
  }
})

const 同步看板数据 = async () => {
  加载中.value = true
  try {
    if (authStore.isAuthenticated && !authStore.profile) {
      await authStore.loadProfile()
    }

    const [status, nutrition] = await Promise.all([
      getTodayStatus(转日期参数(已选日期.value)),
      getNutritionRequirements().catch(() => null)
    ])

    每日状态.value = {
      records: status?.records || [],
      totalEnergy: 安全数值(status?.total_energy ?? status?.totalEnergy),
      totalProtein: 安全数值(status?.total_protein ?? status?.totalProtein),
      totalCarbohydrate: 安全数值(status?.total_carbohydrate ?? status?.totalCarbohydrate),
      totalFat: 安全数值(status?.total_fat ?? status?.totalFat),
      targetEnergy: 安全数值(status?.target_energy ?? status?.targetEnergy ?? nutrition?.target_calorie ?? nutrition?.targetCalorie) || 2500,
      targetProtein: 安全数值(status?.target_protein ?? status?.targetProtein ?? nutrition?.target_protein ?? nutrition?.targetProtein) || 150,
      targetCarbohydrate: 安全数值(status?.target_carbohydrate ?? status?.targetCarbohydrate ?? nutrition?.target_carbohydrate ?? nutrition?.targetCarbohydrate) || 220,
      targetFat: 安全数值(status?.target_fat ?? status?.targetFat ?? nutrition?.target_fat ?? nutrition?.targetFat) || 60
    }

    const hasCurrent = 时间线餐次.value.some((item) => item.type === 详情餐次类型.value)
    if (!hasCurrent) {
      详情餐次类型.value = 'breakfast'
    }
  } catch (error) {
    console.error('加载首页看板数据失败', error)
    每日状态.value = {
      ...每日状态.value,
      records: []
    }
  } finally {
    加载中.value = false
  }
}

const 上一天 = async () => {
  已选日期.value = new Date(已选日期.value.getFullYear(), 已选日期.value.getMonth(), 已选日期.value.getDate() - 1)
  await 同步看板数据()
}

const 下一天 = async () => {
  if (是否今天.value) return
  已选日期.value = new Date(已选日期.value.getFullYear(), 已选日期.value.getMonth(), 已选日期.value.getDate() + 1)
  await 同步看板数据()
}

const 切换详情餐次 = (mealType) => {
  详情餐次类型.value = mealType
}

const 去记录餐次 = () => {
  router.push('/intake')
}

onMounted(async () => {
  加载餐次时间配置()
  加载餐次个性化配置()
  await 同步看板数据()
})
</script>

<style scoped>
:global(body) {
  font-family: 'Inter', 'Segoe UI', Roboto, Arial, sans-serif;
}

.nutrition-dashboard {
  min-height: 100vh;
  padding: 16px;
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(circle at 88% 10%, rgba(242, 252, 247, 0.95), transparent 34%),
    radial-gradient(circle at 10% 74%, rgba(226, 242, 234, 0.9), transparent 35%),
    #edf6f2;
  color: #16382f;
}

.nutrition-dashboard::before,
.nutrition-dashboard::after {
  content: '';
  position: absolute;
  width: 360px;
  height: 360px;
  border-radius: 46% 54% 63% 37% / 43% 38% 62% 57%;
  pointer-events: none;
  z-index: 0;
  filter: blur(2px);
}

.nutrition-dashboard::before {
  right: -120px;
  top: 180px;
  background: radial-gradient(circle at 35% 35%, rgba(123, 214, 184, 0.22), rgba(93, 161, 138, 0));
}

.nutrition-dashboard::after {
  left: -160px;
  bottom: 40px;
  background: radial-gradient(circle at 50% 50%, rgba(158, 213, 191, 0.2), rgba(118, 184, 159, 0));
}

.glass-card {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.9), rgba(248, 253, 250, 0.96));
  border: 1px solid rgba(22, 56, 47, 0.08);
  border-radius: 0;
  box-shadow: 0 16px 34px rgba(18, 64, 50, 0.1);
}

.top-header,
.center-panel,
.timeline-content,
.overview-card,
.details-card,
.food-table {
  border-radius: 0;
  clip-path: polygon(0 0, calc(100% - 14px) 0, 100% 14px, 100% 100%, 14px 100%, 0 calc(100% - 14px));
}

.top-header {
  min-height: 96px;
  padding: 16px 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 14px;
  position: relative;
  z-index: 2;
}

.brand-wrap {
  width: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  gap: 6px;
}

.brand {
  margin: 0;
  font-size: 2rem;
  letter-spacing: 0.03em;
  color: #184c3d;
  font-weight: 800;
}

.brand-date {
  margin: 0;
  font-size: 1rem;
  color: #5f8578;
  font-weight: 600;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
}

.icon-circle {
  width: 40px;
  height: 40px;
  border: 1px solid rgba(43, 116, 92, 0.16);
  border-radius: 50%;
  background: linear-gradient(145deg, #ffffff, #ebf7f1);
  color: #245f4d;
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: 0.24s ease;
  box-shadow: 0 8px 18px rgba(44, 113, 90, 0.16);
}

.icon-circle:hover {
  transform: translateY(-1px);
}

.icon-circle:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  transform: none;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(540px, 1.1fr) minmax(360px, 0.9fr);
  align-items: start;
  gap: 24px;
  position: relative;
  z-index: 1;
}

.center-panel {
  background: linear-gradient(165deg, rgba(248, 253, 250, 0.64), rgba(242, 250, 246, 0.3));
  border: 1px solid rgba(29, 90, 71, 0.07);
  border-radius: 0;
  padding: 18px 16px 22px;
  backdrop-filter: blur(2px);
}

.loading-tip {
  margin: 0 0 12px;
  color: #5d8778;
  font-size: 0.86rem;
  font-weight: 600;
}

.record-ops {
  margin-bottom: 14px;
  padding: 10px 12px;
  border: 1px solid rgba(33, 97, 77, 0.12);
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.ops-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ops-date {
  margin: 0 0 0 6px;
  font-size: 0.9rem;
  color: #3a6e5d;
  font-weight: 700;
}

.calendar-picker {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #2a6552;
  font-size: 0.86rem;
  font-weight: 700;
}

.calendar-picker input {
  border: 1px solid rgba(45, 112, 90, 0.2);
  background: #f5fbf8;
  color: #2f6554;
  height: 34px;
  padding: 0 8px;
}

.timeline {
  padding-left: 22px;
  border-left: 2px solid rgba(76, 151, 124, 0.42);
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  left: -2px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: linear-gradient(180deg, rgba(255, 190, 94, 0.88), rgba(96, 201, 170, 0.9), rgba(84, 174, 143, 0.9));
  filter: drop-shadow(0 0 8px rgba(86, 175, 145, 0.5));
}

.timeline-row {
  display: grid;
  grid-template-columns: 98px minmax(280px, 1fr);
  gap: 12px;
  position: relative;
  cursor: grab;
  transition: transform 0.28s ease;
}

.timeline-row:hover {
  transform: translateX(4px) scale(1.01);
}

.timeline-row:active {
  cursor: grabbing;
}

.time-col {
  padding-top: 10px;
  position: relative;
}

.dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #2f7c63;
  position: absolute;
  left: -25px;
  top: 13px;
  border: 2px solid #e8f3ee;
  box-shadow: 0 0 0 4px rgba(110, 192, 166, 0.22), 0 0 10px rgba(86, 166, 139, 0.42);
}

.time-editor {
  min-height: 24px;
  display: flex;
  align-items: center;
}

.time-display-btn {
  border: none;
  background: transparent;
  color: #5f8a7b;
  font-size: 0.84rem;
  font-weight: 700;
  padding: 0;
  cursor: pointer;
}

.time-input {
  width: 92px;
  border: 1px solid rgba(47, 124, 99, 0.25);
  border-radius: 8px;
  background: #f4faf6;
  color: #3f6f60;
  padding: 3px 6px;
  font-size: 0.78rem;
}

.timeline-content {
  padding: 16px;
  border-radius: 0;
  position: relative;
  overflow: hidden;
}

.meal-card::before {
  content: '';
  position: absolute;
  width: 120px;
  height: 120px;
  right: -26px;
  top: -34px;
  border-radius: 50%;
  background: radial-gradient(circle at 40% 40%, rgba(255, 255, 255, 0.5), rgba(255, 255, 255, 0));
}

.meal-card::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  border: 1px solid rgba(255, 255, 255, 0.34);
  pointer-events: none;
}

.meal-breakfast {
  background: linear-gradient(145deg, rgba(255, 243, 201, 0.95), rgba(255, 214, 153, 0.96));
}

.meal-lunch {
  background: linear-gradient(145deg, rgba(220, 247, 229, 0.95), rgba(183, 239, 231, 0.96));
}

.meal-snack {
  background: linear-gradient(145deg, rgba(255, 232, 198, 0.95), rgba(255, 214, 171, 0.96));
}

.meal-dinner {
  background: linear-gradient(145deg, rgba(218, 237, 255, 0.95), rgba(192, 221, 255, 0.96));
}

.meal-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 10px;
}

.meal-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meal-head h3 {
  margin: 0;
  font-size: 1.24rem;
  color: #204c3f;
}

.meal-name-input {
  width: 150px;
  border: 1px solid rgba(47, 124, 99, 0.28);
  border-radius: 8px;
  background: #f4faf6;
  color: #2c6251;
  padding: 5px 8px;
  font-size: 0.96rem;
  font-weight: 700;
}

.rename-link {
  border: none;
  background: transparent;
  color: #3d7d67;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  padding: 0;
}

.kcal {
  margin: 0;
  color: #4d7c6d;
  font-weight: 700;
  font-size: 0.86rem;
}

.meal-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.hero-badge {
  width: 42px;
  height: 42px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.65);
  display: grid;
  place-items: center;
  font-size: 1.1rem;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.7);
}

.hero-illustration {
  font-size: 2.1rem;
  filter: drop-shadow(0 4px 8px rgba(49, 74, 66, 0.2));
}

.food-hints {
  display: flex;
  gap: 8px;
  margin-bottom: 10px;
}

.food-chip {
  width: 32px;
  height: 32px;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.72);
  display: grid;
  place-items: center;
  font-size: 1rem;
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.72);
}

.snack-group {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
  margin-bottom: 10px;
}

.snack-mini-card {
  border-radius: 0;
  background: rgba(255, 255, 255, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.72);
  min-height: 94px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.mini-emoji {
  font-size: 1.2rem;
}

.mini-label {
  font-size: 0.75rem;
  font-weight: 700;
  color: #4b6f64;
}

.mini-plus {
  border: none;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  background: linear-gradient(145deg, #2f8e6f, #4ab98d);
  color: #fff;
  font-weight: 800;
  cursor: pointer;
}

.bar-track {
  height: 10px;
  border-radius: 999px;
  background: rgba(236, 243, 239, 0.82);
  overflow: hidden;
  margin-bottom: 12px;
}

.bar-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2f7c63, #52ad8b);
}

.meal-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.record-btn-main {
  border: none;
  border-radius: 0;
  min-height: 42px;
  padding: 0 14px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: white;
  font-weight: 700;
  background: linear-gradient(135deg, #2f8e6f, #30b487);
  box-shadow: 0 8px 16px rgba(44, 132, 102, 0.28);
  cursor: pointer;
}

.detail-link {
  border: none;
  background: transparent;
  color: #2b6d59;
  font-weight: 700;
  font-size: 0.86rem;
  cursor: pointer;
}

.footer-tip {
  margin-top: 8px;
  font-size: 0.76rem;
  color: #346b58;
  font-weight: 700;
  text-align: right;
}

.right-panel {
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-height: 0;
  margin-top: 8px;
}

.overview-card {
  padding: 16px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
  position: relative;
}

.floating-red {
  position: absolute;
  width: 26px;
  height: 26px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  color: white;
  background: linear-gradient(145deg, #f14652, #e12f3c);
  box-shadow: 0 8px 16px rgba(241, 70, 82, 0.35);
}

.floating-red.top-right {
  right: 14px;
  top: -10px;
}

.floating-red.lower-right {
  right: 16px;
  bottom: -12px;
}

.ring-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.calorie-ring {
  width: 190px;
  height: 190px;
  border-radius: 50%;
  display: grid;
  place-items: center;
}

.ring-center {
  width: 142px;
  height: 142px;
  border-radius: 50%;
  background: #f6fbf8;
  display: grid;
  place-items: center;
  text-align: center;
}

.ring-center h3 {
  margin: 0;
  font-size: 2rem;
  color: #1f5a46;
}

.ring-center p {
  margin: 2px 0 0;
  font-size: 0.8rem;
  color: #668f80;
  font-weight: 700;
}

.goal-text {
  margin: 8px 0 0;
  font-size: 0.95rem;
  color: #5a8375;
  font-weight: 600;
}

.macro-stack {
  display: flex;
  flex-direction: column;
  gap: 10px;
  justify-content: center;
}

.macro-card {
  border-radius: 0;
  background: #f8fcfa;
  border: 1px solid rgba(35, 89, 71, 0.08);
  padding: 10px 12px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.macro-left {
  display: flex;
  gap: 10px;
  align-items: center;
}

.macro-icon {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  display: grid;
  place-items: center;
}

.macro-label {
  margin: 0;
  color: #245946;
  font-weight: 700;
}

.macro-target {
  margin: 2px 0 0;
  color: #6d9384;
  font-size: 0.77rem;
}

.macro-value {
  margin: 0;
  color: #1f5947;
  font-weight: 800;
  font-size: 1.16rem;
}

.details-card {
  padding: 16px;
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 520px;
}

.details-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.details-head h3 {
  margin: 0;
  color: #235747;
  font-size: 1.62rem;
}

.details-head p {
  margin: 0;
  color: #6b9283;
  font-size: 0.86rem;
}

.sample-note {
  margin: 0 0 10px;
  color: #5f8779;
  font-size: 0.82rem;
  font-weight: 700;
}

.micro-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 16px;
}

.micro-top {
  display: flex;
  justify-content: space-between;
  color: #507c6d;
  font-size: 0.82rem;
  font-weight: 700;
  margin-bottom: 6px;
}

.micro-track {
  height: 5px;
  border-radius: 999px;
  overflow: hidden;
  background: #dde8e2;
}

.micro-fill {
  height: 100%;
}

.food-table {
  border-radius: 0;
  background: #f8fcfa;
  border: 1px solid rgba(30, 86, 68, 0.08);
  padding: 10px 12px;
  flex: 1;
}

.table-head,
.food-row {
  display: grid;
  grid-template-columns: 1.6fr 0.7fr 0.7fr;
  align-items: center;
}

.table-head {
  color: #709384;
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 6px;
  font-weight: 700;
}

.food-row {
  height: 54px;
  border-top: 1px solid rgba(28, 84, 67, 0.08);
  font-size: 0.95rem;
  color: #295c4c;
}

.food-empty {
  color: #7d9f92;
}

.food-name {
  display: flex;
  align-items: center;
  gap: 9px;
  font-weight: 600;
}

.sample-tag {
  font-size: 0.68rem;
  padding: 2px 6px;
  border: 1px solid rgba(74, 131, 112, 0.25);
  color: #4f7f6f;
  background: rgba(232, 246, 239, 0.9);
}

.thumb {
  width: 32px;
  height: 32px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  color: #f4fffb;
}

.strong-text {
  font-weight: 800;
}

.add-more {
  margin-top: 12px;
  border-radius: 0;
  height: 48px;
  border: 1px solid rgba(40, 97, 78, 0.18);
  background: #eef6f1;
  display: flex;
  align-items: center;
  padding: 6px;
}

.add-more input {
  border: none;
  outline: none;
  background: transparent;
  padding: 0 12px;
  color: #2c6251;
  flex: 1;
}

.plus-btn {
  width: 34px;
  height: 34px;
  border-radius: 0;
  border: none;
  background: #2f7c63;
  color: white;
  font-size: 1.3rem;
  cursor: pointer;
}

@media (max-width: 1360px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .right-panel {
    margin-top: 0;
  }
}

@media (max-width: 980px) {
  .top-header {
    min-height: 112px;
  }

  .brand {
    font-size: 1.65rem;
  }

  .top-actions {
    position: static;
    transform: none;
    margin-top: 4px;
  }

  .brand-wrap {
    gap: 8px;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .timeline-row {
    grid-template-columns: 72px 1fr;
    transform: none !important;
  }

  .overview-card {
    grid-template-columns: 1fr;
  }

  .snack-group {
    grid-template-columns: 1fr;
  }

  .record-ops {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
