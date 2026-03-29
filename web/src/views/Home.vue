<template>
  <div class="nutrition-dashboard">
    <header class="top-header glass-card">
      <div class="brand-wrap">
        <h1 class="brand">基于用户健康数据与营养模型的饮食推荐系统</h1>
      </div>

      <div class="top-actions">
      </div>
    </header>

    <main class="dashboard-grid">
      <section class="center-panel">
        <div class="section-head">
          <div>
            <p class="kicker">营养记录</p>
            <h2 class="headline">{{ 日期标题 }}</h2>
          </div>
          <div class="date-controls">
            <button class="icon-circle soft" aria-label="上一天" @click="上一天"><n-icon><ChevronBackOutline /></n-icon></button>
            <button class="icon-circle soft" aria-label="下一天" :disabled="是否今天" @click="下一天"><n-icon><ChevronForwardOutline /></n-icon></button>
          </div>
        </div>

        <p v-if="加载中" class="loading-tip">正在同步饮食数据...</p>

        <div class="timeline">
          <article
            v-for="meal in 时间线餐次"
            :key="meal.type"
            class="timeline-row"
            :class="meal.variant"
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

            <div class="timeline-content glass-card">
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
                    class="rename-btn"
                    aria-label="编辑餐次名称"
                    title="编辑餐次名称"
                    @click.stop="开始编辑名称(meal.type, meal.name)"
                  >
                    <n-icon><CreateOutline /></n-icon>
                  </button>
                </div>
                <p v-if="meal.kcal" class="kcal">{{ meal.kcal }}</p>
              </div>

              <template v-if="meal.variant === 'filled'">
                <div class="meal-inline">
                  <div class="mini-icon">{{ meal.icon }}</div>
                  <div class="bar-track">
                    <div class="bar-fill" :style="{ width: meal.progress + '%' }" />
                  </div>
                  <button class="icon-circle xsmall" aria-label="编辑">
                    <n-icon><CreateOutline /></n-icon>
                  </button>
                </div>
                <button class="detail-btn detail-unified" @click.stop="切换详情餐次(meal.type)">查看详情</button>
              </template>

              <template v-else-if="meal.variant === 'logged'">
                <div class="logged-card">
                  <div class="status-row">
                    <span class="status-badge">状态：已记录</span>
                    <span class="status-kcal">{{ meal.kcal }}</span>
                  </div>
                  <div class="bar-track strong">
                    <div class="bar-fill" :style="{ width: meal.progress + '%' }" />
                  </div>
                  <button class="detail-btn detail-unified" @click.stop="切换详情餐次(meal.type)">查看详情</button>
                </div>
              </template>

              <template v-else>
                <div class="meal-actions">
                  <button class="ghost-record" @click="去记录餐次">+ 记录餐次</button>
                  <button class="detail-btn empty detail-unified" @click.stop="切换详情餐次(meal.type)">查看详情</button>
                </div>
              </template>
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
          <span class="floating-red lower-right"><n-icon><HeartOutline /></n-icon></span>
          <div class="details-head">
            <h3>{{ 当前详情标题 }}</h3>
            <p>{{ 当前详情时间 }}</p>
          </div>

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

            <div v-if="!当前详情食物.length" class="food-row food-empty">
              <div class="food-name">暂无记录</div>
              <span>-</span>
              <span class="strong-text">-</span>
            </div>

            <div v-for="food in 当前详情食物" :key="food.id" class="food-row">
              <div class="food-name">
                <span class="thumb" :style="{ background: food.bg }">{{ food.emoji }}</span>
                <span>{{ food.name }}</span>
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
  ChevronBackOutline,
  ChevronForwardOutline,
  CreateOutline,
  HeartOutline,
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
  { type: 'breakfast', name: '早餐', time: '08:15', icon: '🥣', ratio: 0.16, emoji: '🍳', bg: 'linear-gradient(135deg, #78b7ff, #4d84d7)' },
  { type: 'lunch', name: '午餐', time: '12:30', icon: '🥗', ratio: 0.32, emoji: '🥗', bg: 'linear-gradient(135deg, #8cc7a5, #5e9b77)' },
  { type: 'snack', name: '下午加餐', time: '16:00', icon: '🍎', ratio: 0.2, emoji: '🍎', bg: 'linear-gradient(135deg, #f7ba80, #e08d4f)' },
  { type: 'dinner', name: '晚餐', time: '19:30', icon: '🍽️', ratio: 0.32, emoji: '🍽️', bg: 'linear-gradient(135deg, #afcf78, #7ca451)' }
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
    餐次时间映射.value = {
      ...默认餐次时间,
      ...parsed
    }
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
      餐次名称映射.value = {
        ...默认餐次名称,
        ...parsedName
      }
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

const 日期标题 = computed(() => {
  const now = new Date()
  const isToday = now.toDateString() === 已选日期.value.toDateString()
  const month = 已选日期.value.getMonth() + 1
  const day = 已选日期.value.getDate()
  return isToday ? `今天，${month}月${day}日` : `${month}月${day}日`
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
    bg: fallback.bg
  }))
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
  background:
    radial-gradient(circle at 90% 8%, rgba(231, 244, 238, 0.9), transparent 30%),
    radial-gradient(circle at 12% 70%, rgba(206, 231, 220, 0.7), transparent 32%),
    #edf4f0;
  color: #16382f;
}

.glass-card {
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.88), rgba(248, 253, 250, 0.92));
  border: 1px solid rgba(22, 56, 47, 0.08);
  border-radius: 22px;
  box-shadow: 0 12px 30px rgba(18, 64, 50, 0.09);
}

.top-header {
  height: 72px;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 14px;
}

.brand-wrap {
  width: 100%;
  display: flex;
  justify-content: center;
}

.brand {
  margin: 0;
  font-size: 1.7rem;
  letter-spacing: -0.02em;
  color: #1b5f4b;
  text-align: center;
}

.global-nav {
  display: flex;
  gap: 30px;
}

.nav-link {
  text-decoration: none;
  color: #4f7f71;
  font-weight: 600;
  position: relative;
}

.nav-link.active,
.nav-link:hover {
  color: #1e6f57;
}

.nav-link.active::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: -10px;
  height: 2px;
  border-radius: 999px;
  background: #2f7c63;
}

.top-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.loading-tip {
  margin: 0 0 10px;
  color: #5d8778;
  font-size: 0.86rem;
  font-weight: 600;
}

.icon-circle {
  width: 38px;
  height: 38px;
  border: 1px solid rgba(30, 111, 87, 0.18);
  border-radius: 50%;
  background: #f6fcf9;
  color: #245f4d;
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: 0.2s ease;
}

.icon-circle:hover {
  transform: translateY(-1px);
  background: #eaf4ef;
}

.icon-circle:disabled {
  opacity: 0.45;
  cursor: not-allowed;
  transform: none;
}

.icon-circle.soft {
  width: 34px;
  height: 34px;
  background: #edf5f1;
}

.icon-circle.xsmall {
  width: 30px;
  height: 30px;
}

.profile-pill {
  height: 42px;
  border-radius: 999px;
  background: #f2faf6;
  border: 1px solid rgba(30, 111, 87, 0.16);
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 4px 6px 4px 14px;
}

.profile-meta {
  line-height: 1.05;
}

.profile-name {
  margin: 0;
  font-size: 0.84rem;
  color: #214f42;
  font-weight: 700;
}

.profile-tier {
  margin: 2px 0 0;
  font-size: 0.7rem;
  color: #709486;
}

.avatar {
  width: 30px;
  height: 30px;
  border-radius: 50%;
  background: linear-gradient(145deg, #1b6b53, #2f7c63);
  color: white;
  display: grid;
  place-items: center;
  font-weight: 700;
  text-transform: uppercase;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: minmax(460px, 0.92fr) minmax(420px, 1fr);
  gap: 18px;
}

.sidebar {
  min-height: calc(100vh - 118px);
  padding: 18px 16px;
  display: flex;
  flex-direction: column;
}

.org-card {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 14px;
}

.org-logo {
  width: 42px;
  height: 42px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  background: #2f7c63;
  color: white;
}

.org-title {
  margin: 0;
  font-weight: 700;
  color: #1f5445;
}

.org-subtitle {
  margin: 2px 0 0;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 0.64rem;
  color: #6d9385;
}

.menu-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-top: 8px;
}

.menu-item {
  border: none;
  border-radius: 14px;
  background: transparent;
  color: #467264;
  height: 45px;
  padding: 0 12px;
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 0.96rem;
  font-weight: 600;
  cursor: pointer;
  text-align: left;
}

.menu-item:hover {
  background: #e9f3ee;
}

.menu-item.active {
  color: #225944;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px rgba(37, 101, 80, 0.18), 0 6px 14px rgba(34, 88, 69, 0.08);
}

.record-btn {
  margin-top: auto;
  border: none;
  height: 48px;
  border-radius: 999px;
  background: linear-gradient(140deg, #2f7c63, #3a9a78);
  color: white;
  font-weight: 700;
  font-size: 0.95rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(43, 117, 91, 0.3);
}

.center-panel {
  background: rgba(250, 254, 252, 0.56);
  border: 1px solid rgba(29, 90, 71, 0.08);
  border-radius: 24px;
  padding: 20px;
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.kicker {
  margin: 0;
  font-size: 0.88rem;
  color: #6f9587;
  font-weight: 600;
}

.headline {
  margin: 4px 0 0;
  font-size: 1.45rem;
  color: #214f42;
}

.date-controls {
  display: flex;
  gap: 8px;
}

.timeline {
  padding-left: 18px;
  border-left: 2px solid rgba(52, 125, 99, 0.2);
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.timeline-row {
  display: grid;
  grid-template-columns: 92px 1fr;
  gap: 10px;
  position: relative;
  cursor: grab;
}

.timeline-row:active {
  cursor: grabbing;
}

.time-col {
  padding-top: 10px;
  position: relative;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #2f7c63;
  position: absolute;
  left: -24px;
  top: 13px;
  border: 2px solid #e8f3ee;
}

.time {
  font-size: 0.78rem;
  color: #709486;
  font-weight: 600;
}

.timeline-content {
  padding: 14px;
  max-width: 720px;
}

.time-editor {
  min-height: 24px;
  display: flex;
  align-items: center;
}

.time-display-btn {
  border: none;
  background: transparent;
  color: #709486;
  font-size: 0.78rem;
  font-weight: 600;
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

.meal-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 12px;
}

.meal-title-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.meal-head h3 {
  margin: 0;
  font-size: 1.26rem;
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

.rename-btn {
  border: none;
  border-radius: 999px;
  width: 28px;
  height: 28px;
  padding: 0;
  background: rgba(47, 124, 99, 0.12);
  color: #2a6a56;
  display: grid;
  place-items: center;
  cursor: pointer;
}

.rename-btn:hover {
  background: rgba(47, 124, 99, 0.2);
}

.kcal {
  margin: 0;
  color: #5f8678;
  font-weight: 700;
  font-size: 0.92rem;
}

.meal-inline {
  display: flex;
  align-items: center;
  gap: 10px;
}

.mini-icon {
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: #e4f2eb;
  display: grid;
  place-items: center;
}

.bar-track {
  flex: 1;
  height: 10px;
  border-radius: 999px;
  background: #deebe4;
  overflow: hidden;
}

.bar-track.strong {
  background: rgba(222, 235, 228, 0.35);
}

.bar-fill {
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2f7c63, #52ad8b);
}

.logged-card {
  border-radius: 24px;
  background: linear-gradient(145deg, #21684f, #2f7c63);
  color: white;
  padding: 16px;
}

.status-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.status-badge {
  font-size: 0.78rem;
  opacity: 0.92;
}

.status-kcal {
  font-size: 0.85rem;
  font-weight: 700;
}

.detail-btn {
  margin-top: 12px;
  margin-left: auto;
  display: block;
  border: none;
  border-radius: 999px;
  height: 34px;
  min-width: 118px;
  background: rgba(255, 255, 255, 0.2);
  color: white;
  font-weight: 700;
  cursor: pointer;
}

.detail-btn.detail-unified {
  min-width: 126px;
}

.ghost-record {
  width: 100%;
  height: 54px;
  border-radius: 999px;
  border: 1.5px dashed rgba(53, 113, 91, 0.34);
  background: rgba(255, 255, 255, 0.55);
  color: #547f70;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
}

.meal-actions {
  display: grid;
  gap: 10px;
}

.detail-btn.empty {
  margin-top: 0;
  min-width: 100%;
  height: 42px;
  background: rgba(66, 120, 99, 0.12);
  color: #2d6553;
}

.right-panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
  min-height: 0;
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
  letter-spacing: 0.08em;
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
  border-radius: 18px;
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
  border-radius: 16px;
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
  border-radius: 999px;
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

.add-more input::placeholder {
  color: #7b9e90;
}

.plus-btn {
  width: 34px;
  height: 34px;
  border-radius: 50%;
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
    grid-column: auto;
  }
}

@media (max-width: 980px) {
  .global-nav {
    display: none;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .sidebar {
    min-height: auto;
  }

  .overview-card {
    grid-template-columns: 1fr;
  }

  .timeline-row {
    grid-template-columns: 72px 1fr;
  }
}
</style>
