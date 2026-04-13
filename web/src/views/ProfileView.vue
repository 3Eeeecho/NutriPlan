<template>
  <div class="profile-view-page">
    <div class="profile-shell">
      <header class="page-header glass-card">
        <BackButton />
        <div class="page-title-wrap">
          <h1 class="page-title">个人主页</h1>
          <p class="page-subtitle">清晰了解身体状态，轻松管理每日营养计划</p>
        </div>
        <n-button type="primary" class="update-btn" @click="goToEdit">
          <template #icon>
            <n-icon><CreateOutline /></n-icon>
          </template>
          更新档案
        </n-button>
      </header>

      <main class="page-content">
        <section v-if="loading" class="loading-section">
          <n-skeleton text :repeat="10" />
        </section>

        <section v-else-if="!authStore.hasProfile" class="empty-section">
          <n-empty description="您还没有完善健康档案">
            <template #extra>
              <n-button type="primary" @click="goToEdit">立即完善档案</n-button>
            </template>
          </n-empty>
        </section>

        <section v-else class="profile-content">
          <div class="bento-grid">
            <article class="bento-card profile-card">
              <div class="profile-card__head">
                <div class="avatar-chip">{{ getUserInitial() }}</div>
                <div>
                  <p class="hero-greet">你好，{{ authStore.user?.username || 'Nutri 用户' }}</p>
                  <h2 class="hero-goal">今天也在靠近理想状态</h2>
                </div>
              </div>

              <div class="profile-meta-row">
                <span><svg class="meta-icon" viewBox="0 0 24 24" fill="none"><path d="M12 12a4 4 0 1 0 0-8 4 4 0 0 0 0 8Zm-7 8a7 7 0 0 1 14 0" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>{{ authStore.profile?.age || '--' }} 岁</span>
                <span class="meta-divider"></span>
                <span><svg class="meta-icon" viewBox="0 0 24 24" fill="none"><path d="M12 3v18" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><path d="M8 6h4M9 10h3M8 14h4M9 18h3" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/></svg>{{ authStore.profile?.height || '--' }} cm</span>
                <span class="meta-divider"></span>
                <span><svg class="meta-icon" viewBox="0 0 24 24" fill="none"><path d="M5 9a7 7 0 0 1 14 0v7H5V9Z" stroke="currentColor" stroke-width="1.8"/><path d="M12 11l2-2" stroke="currentColor" stroke-width="1.8" stroke-linecap="round"/><circle cx="12" cy="11" r="1" fill="currentColor"/></svg>{{ authStore.profile?.weight || '--' }} kg</span>
              </div>

              <div class="profile-tags">
                <span class="status-tag tag-goal">{{ authStore.profile?.health_goal || authStore.profile?.healthGoal || '维持健康' }}</span>
                <span class="status-tag tag-activity">{{ authStore.profile?.activity_level || authStore.profile?.activityLevel || '中度活动' }}</span>
              </div>

              <div class="profile-note">目标体重 {{ authStore.profile?.target_weight ?? authStore.profile?.targetWeight ?? '--' }} kg · 每日 {{ authStore.profile?.meal_times_per_day || authStore.profile?.mealTimesPerDay || 3 }} 餐</div>

              <div class="profile-mini-stats">
                <div class="mini-stat">
                  <span>BMI</span>
                  <strong>{{ authStore.profile?.bmi || '--' }}</strong>
                </div>
                <div class="mini-stat">
                  <span>BMR</span>
                  <strong>{{ authStore.profile?.bmr || '--' }}</strong>
                </div>
                <div class="mini-stat">
                  <span>TDEE</span>
                  <strong>{{ authStore.profile?.tdee || '--' }}</strong>
                </div>
              </div>
            </article>

            <article class="bento-card hero-nutrition-card soft-green">
              <div class="section-title-row">
                <h3>每日营养需求</h3>
                <n-tag round size="small" type="success">{{ authStore.profile?.health_goal || authStore.profile?.healthGoal || '个性化推荐' }}</n-tag>
              </div>

              <div class="hero-donut-layout">
                <div class="donut-cluster" aria-label="宏观元素比例环图">
                  <div class="macro-ring ring-fat" :style="{ '--pct': clampPercent(nutritionData?.fat_ratio ?? 0) }"></div>
                  <div class="macro-ring ring-carb" :style="{ '--pct': clampPercent(nutritionData?.carb_ratio ?? 0) }"></div>
                  <div class="macro-ring ring-protein" :style="{ '--pct': clampPercent(nutritionData?.protein_ratio ?? 0) }"></div>
                  <div class="donut-center">
                    <strong>{{ nutritionData?.target_calorie ?? '--' }}</strong>
                    <span>kcal</span>
                    <small>今日建议</small>
                  </div>
                </div>

                <div class="macro-legend">
                  <div v-for="item in macroItems" :key="item.key" class="legend-item">
                    <div class="legend-item__name"><span class="legend-dot" :class="`dot-${item.key}`"></span>{{ item.label }}</div>
                    <div class="legend-item__value">{{ item.gramLabel }}</div>
                    <div class="legend-item__ratio">{{ item.ratioLabel }}</div>
                  </div>
                </div>
              </div>
            </article>

            <article class="bento-card vital-card">
              <div class="section-title-row compact">
                <h3>核心健康指标</h3>
                <span class="subtle">BMI / BMR / TDEE</span>
              </div>

              <div class="bmi-head">
                <div>
                  <p class="bmi-label">BMI {{ getBMIDesc(authStore.profile?.bmi) }}</p>
                  <p class="bmi-value" :class="getBMIClass(authStore.profile?.bmi)">{{ authStore.profile?.bmi || '--' }}</p>
                </div>
              </div>

              <div class="bmi-gauge">
                <div class="bmi-track"></div>
                <div class="bmi-pointer" :style="{ left: `${bmiPosition}%` }"></div>
                <div
                  v-for="mark in bmiThresholdMarks"
                  :key="mark.value"
                  class="bmi-threshold-mark"
                  :style="{ left: `${mark.left}%` }"
                >
                  <span>{{ mark.value }}</span>
                </div>
              </div>
              <div class="bmi-scale">
                <span>偏瘦</span><span>正常</span><span>偏胖</span><span>肥胖</span>
              </div>

              <div class="vital-mini-grid">
                <div class="vital-mini-card">
                  <div class="vital-mini-card__label"><svg class="meta-icon" viewBox="0 0 24 24" fill="none"><path d="M12 3c2 3 5 5 5 9a5 5 0 1 1-10 0c0-4 3-6 5-9Z" stroke="currentColor" stroke-width="1.8"/></svg>BMR</div>
                  <div class="vital-mini-card__value">{{ authStore.profile?.bmr || '--' }}</div>
                  <div class="vital-mini-card__unit">基础代谢 kcal</div>
                </div>
                <div class="vital-mini-card">
                  <div class="vital-mini-card__label"><svg class="meta-icon" viewBox="0 0 24 24" fill="none"><path d="M4 15c3-5 5-6 8-6 3 0 5 1 8 6" stroke="currentColor" stroke-width="1.8"/><path d="M12 9V4" stroke="currentColor" stroke-width="1.8"/></svg>TDEE</div>
                  <div class="vital-mini-card__value">{{ authStore.profile?.tdee || '--' }}</div>
                  <div class="vital-mini-card__unit">日维持热量 kcal</div>
                </div>
              </div>
            </article>

            <article class="bento-card mode-card soft-green">
              <div class="section-title-row compact">
                <h3>饮食调节模式</h3>
                <span v-if="currentDietMode.until" class="subtle">截止 {{ formatDate(currentDietMode.until) }}</span>
              </div>

              <div class="mode-toolbar">
                <div class="diet-mode-days">
                  <span>天数</span>
                  <n-input-number v-model:value="modeDays" :min="1" :max="7" size="small" :disabled="modeSaving" />
                </div>
                <n-tag size="small" round :type="dietModeTagType(currentDietMode.mode)">{{ dietModeLabel(currentDietMode.mode) }}</n-tag>
              </div>

              <div class="segment-control" :class="{ disabled: modeSaving }">
                <div class="segment-slider" :style="segmentSliderStyle"></div>
                <button
                  v-for="item in dietModes"
                  :key="item.mode"
                  class="segment-item"
                  :class="{ active: currentDietMode.mode === item.mode }"
                  :disabled="modeSaving"
                  @click="changeDietMode(item.mode, modeDays, item.reason)"
                >
                  {{ item.label }}
                </button>
              </div>

              <transition name="mode-fade" mode="out-in">
                <div class="mode-explain-item active" :key="currentModeMeta.mode">
                  <strong>{{ currentModeMeta.label }}</strong>
                  <p>{{ currentModeMeta.desc }}</p>
                </div>
              </transition>
            </article>

            <article class="bento-card prefs-card">
              <div class="section-title-row compact">
                <h3>个性化偏好</h3>
                <span class="subtle">饮食与健康备注</span>
              </div>

              <div class="pref-list">
                <div class="pref-item">
                  <span class="pref-item__label">过敏源</span>
                  <span class="pref-item__value">{{ authStore.profile?.allergies || '无' }}</span>
                </div>
                <div class="pref-item">
                  <span class="pref-item__label">饮食偏好</span>
                  <span class="pref-item__value">{{ authStore.profile?.dietary_prefs || authStore.profile?.dietaryPrefs || '无' }}</span>
                </div>
                <div class="pref-item">
                  <span class="pref-item__label">健康问题</span>
                  <span class="pref-item__value">{{ authStore.profile?.health_conditions || authStore.profile?.healthConditions || '无' }}</span>
                </div>
              </div>
            </article>
          </div>

          <div class="action-buttons">
            <n-button type="primary" size="large" @click="goToEdit">
              <template #icon>
                <n-icon><CreateOutline /></n-icon>
              </template>
              更新档案
            </n-button>
            <n-button size="large" @click="goBack">返回首页</n-button>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NButton, NIcon, NSkeleton, NEmpty, NTag, NInputNumber } from 'naive-ui'
import { CreateOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/store/auth'
import { getNutritionRequirements, getDietMode, setDietMode } from '@/api/user'
import BackButton from '@/components/layout/BackButton.vue'
import { useMessage } from 'naive-ui'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const nutritionData = ref(null)
const message = useMessage()
const modeSaving = ref(false)
const modeDays = ref(2)
const currentDietMode = ref({ mode: 'normal', source: 'auto', reason: '', until: null })
const dietModes = [
  { mode: 'normal', label: '标准', reason: '手动恢复标准', desc: '按目标配比执行日常计划，适合常规状态。' },
  { mode: 'light_adjust', label: '轻调', reason: '手动轻调', desc: '小幅修正前一日偏差，降低波动更易坚持。' },
  { mode: 'bland', label: '清淡', reason: '生病/肠胃不适', desc: '减少刺激与油脂负担，优先温和与易消化。' },
  { mode: 'heavy_adjust', label: '重调', reason: '昨日重度偏离', desc: '在短期内加强约束，快速拉回营养节奏。' }
]

onMounted(async () => {
  loading.value = true
  try {
    if (!authStore.profile) {
      await authStore.loadProfile()
    }
    if (authStore.hasProfile) {
      try {
        const response = await getNutritionRequirements()
        nutritionData.value = response
      } catch (error) {
        const errorMsg = error?.response?.data?.error || error.message
        const status = error?.response?.status
        console.log('营养需求数据不可用:', {
          status,
          error: errorMsg,
          url: '/api/v1/user/nutrition'
        })
        nutritionData.value = null
      }

      try {
        const modeResp = await getDietMode()
        currentDietMode.value = {
          mode: modeResp?.mode || 'normal',
          source: modeResp?.source || 'auto',
          reason: modeResp?.reason || '',
          until: modeResp?.until || null
        }
        modeDays.value = 2
      } catch (error) {
        currentDietMode.value = { mode: 'normal', source: 'auto', reason: '', until: null }
        modeDays.value = 2
      }
    }
  } catch (error) {
    console.error('加载档案失败:', error)
  } finally {
    loading.value = false
  }
})

const goBack = () => {
  router.back()
}

const goToEdit = () => {
  router.push('/profile')
}

const getUserInitial = () => {
  const username = authStore.user?.username || 'N'
  return username.charAt(0).toUpperCase()
}

const getBMIClass = (bmi) => {
  if (!bmi) return ''
  if (bmi < 18.5) return 'underweight'
  if (bmi < 24) return 'normal'
  if (bmi < 28) return 'overweight'
  return 'obese'
}

const getBMIDesc = (bmi) => {
  if (!bmi) return '--'
  if (bmi < 18.5) return '偏瘦'
  if (bmi < 24) return '正常'
  if (bmi < 28) return '偏胖'
  return '肥胖'
}

const clampPercent = (value) => {
  const number = Number(value)
  if (Number.isNaN(number)) return 0
  return Math.max(0, Math.min(100, Math.round(number)))
}

const macroItems = computed(() => {
  const data = nutritionData.value || {}
  const normalizeGram = (value) => {
    const number = Number(value)
    if (Number.isNaN(number)) return '--'
    return `${Math.round(number)} g`
  }

  const normalizeRatio = (value) => {
    const number = Number(value)
    if (Number.isNaN(number)) return '--'
    return `${clampPercent(number)}%`
  }

  return [
    {
      key: 'protein',
      label: '蛋白质',
      gramLabel: normalizeGram(data.protein_gram),
      ratioLabel: normalizeRatio(data.protein_ratio)
    },
    {
      key: 'carb',
      label: '碳水',
      gramLabel: normalizeGram(data.carb_gram),
      ratioLabel: normalizeRatio(data.carb_ratio)
    },
    {
      key: 'fat',
      label: '脂肪',
      gramLabel: normalizeGram(data.fat_gram),
      ratioLabel: normalizeRatio(data.fat_ratio)
    }
  ]
})

const modeIndex = computed(() => {
  const index = dietModes.findIndex((item) => item.mode === currentDietMode.value.mode)
  return index < 0 ? 0 : index
})

const segmentSliderStyle = computed(() => {
  const width = 100 / dietModes.length
  return {
    width: `calc(${width}% - 8px)`,
    left: `calc(${modeIndex.value * width}% + 4px)`
  }
})

const bmiPosition = computed(() => {
  const value = Number(authStore.profile?.bmi)
  if (Number.isNaN(value)) return 0
  const mapped = ((value - BMI_RANGE_MIN) / (BMI_RANGE_MAX - BMI_RANGE_MIN)) * 100
  return Math.max(0, Math.min(100, mapped))
})

const currentModeMeta = computed(() => {
  return dietModes.find((item) => item.mode === currentDietMode.value.mode) || dietModes[0]
})

const BMI_RANGE_MIN = 15
const BMI_RANGE_MAX = 35

const bmiLeft = (value) => {
  const mapped = ((value - BMI_RANGE_MIN) / (BMI_RANGE_MAX - BMI_RANGE_MIN)) * 100
  return Math.max(0, Math.min(100, mapped))
}

const bmiThresholdMarks = computed(() => [
  { value: '18.5', left: bmiLeft(18.5) },
  { value: '24.0', left: bmiLeft(24) },
  { value: '28.0', left: bmiLeft(28) }
])

const dietModeLabel = (mode) => {
  const map = {
    normal: '标准模式',
    light_adjust: '轻调模式',
    bland: '清淡模式',
    heavy_adjust: '重调模式'
  }
  return map[mode] || '标准模式'
}

const dietModeTagType = (mode) => {
  const map = {
    normal: 'info',
    light_adjust: 'warning',
    bland: 'success',
    heavy_adjust: 'error'
  }
  return map[mode] || 'default'
}

const formatDate = (value) => {
  if (!value) return '--'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '--'
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

const changeDietMode = async (mode, days, reason) => {
  const safeDays = Math.max(1, Math.min(7, Number(days) || 1))
  modeSaving.value = true
  try {
    const resp = await setDietMode({ mode, days: safeDays, reason })
    currentDietMode.value = {
      mode: resp?.mode || mode,
      source: resp?.source || 'manual',
      reason: resp?.reason || reason || '',
      until: resp?.until || null
    }
    message.success(`已切换为${dietModeLabel(currentDietMode.value.mode)}`)
  } catch (error) {
    message.error(error?.response?.data?.error || '饮食模式更新失败')
  } finally {
    modeSaving.value = false
  }
}
</script>

<style scoped>
.profile-view-page {
  min-height: 100vh;
  --metric-font: 'DIN Alternate', 'Roboto Mono', 'SFMono-Regular', 'Consolas', 'Liberation Mono', monospace;
  font-family: 'Inter', 'PingFang SC', 'Microsoft YaHei', 'Noto Sans SC', sans-serif;
  text-rendering: optimizeLegibility;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  background-color: #f3f6f2;
  background-image:
    radial-gradient(circle at 10% 5%, rgba(255, 255, 255, 0.95) 0%, transparent 45%),
    radial-gradient(circle at 90% 95%, rgba(231, 242, 230, 0.65) 0%, transparent 45%);
}

.profile-shell {
  max-width: 1280px;
  margin: 0 auto;
  padding: 24px 20px 44px;
}

.glass-card,
.bento-card {
  background: #fff;
  border: none;
  border-radius: 16px;
  box-shadow: 0 8px 20px rgba(23, 44, 35, 0.06);
}

.bento-card {
  padding: 18px;
  height: 100%;
}

.soft-green {
  background: linear-gradient(180deg, #f8fcf7 0%, #f3f9f1 100%);
}

.page-header {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 16px;
  margin-bottom: 22px;
  padding: 14px 18px;
}

.page-title-wrap {
  display: grid;
  justify-items: center;
  gap: 4px;
}

.page-title {
  margin: 0;
  font-size: clamp(24px, 3vw, 34px);
  font-weight: 800;
  color: #112c22;
}

.page-subtitle {
  margin: 0;
  font-size: 14px;
  color: #5b7a6c;
}

.update-btn {
  border-radius: 999px;
  background: linear-gradient(135deg, #8dc75d 0%, #72aa46 100%);
  border: none;
  box-shadow: 0 6px 16px rgba(122, 180, 77, 0.24);
}

.page-content {
  display: grid;
  gap: 20px;
}

.loading-section,
.empty-section {
  padding: 30px 10px;
}

.profile-content {
  display: grid;
  gap: 20px;
  animation: fadeInUp 0.35s ease;
}

.bento-grid {
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  gap: 16px;
  align-items: stretch;
}

.profile-card {
  grid-column: span 4;
  display: grid;
  align-content: center;
  gap: 10px;
  min-height: 318px;
}

.hero-nutrition-card {
  grid-column: span 8;
  min-height: 318px;
}

.vital-card {
  grid-column: span 5;
  min-height: 290px;
}

.mode-card {
  grid-column: span 7;
  display: grid;
  align-content: center;
  gap: 10px;
  min-height: 290px;
}

.prefs-card {
  grid-column: span 12;
}

.profile-card__head {
  display: flex;
  align-items: center;
  gap: 14px;
}

.avatar-chip {
  width: 58px;
  height: 58px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  color: #fff;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #8dc75d 0%, #72aa46 100%);
}

.hero-greet {
  margin: 0;
  font-size: 13px;
  color: #688376;
}

.hero-goal {
  margin: 3px 0 0;
  font-size: 20px;
  color: #0f261e;
}

.profile-meta-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  color: #2f4c40;
  font-size: 14px;
}

.profile-meta-row span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  line-height: 1.2;
}

.meta-divider {
  width: 1px;
  height: 12px;
  background: #d9e6dc;
}

.meta-icon {
  width: 14px;
  height: 14px;
  color: #6c8a7b;
}

.profile-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  align-items: center;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  min-height: 32px;
  padding: 0 14px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 600;
  line-height: 1;
  white-space: nowrap;
  text-align: center;
}

.tag-goal {
  background: #e6f4de;
  color: #2f6a3a;
}

.tag-activity {
  background: #e8f0ff;
  color: #2e5dba;
}

.profile-note {
  font-size: 13px;
  color: #658274;
}

.profile-mini-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 8px;
}

.mini-stat {
  border-radius: 10px;
  background: #f4f8f3;
  box-shadow: inset 0 0 0 1px #e4ede3;
  padding: 7px 8px;
  display: grid;
  gap: 3px;
}

.mini-stat span {
  font-size: 11px;
  color: #6e867a;
  line-height: 1;
}

.mini-stat strong {
  font-size: 14px;
  color: #203b30;
  line-height: 1.1;
  font-family: var(--metric-font);
  font-variant-numeric: tabular-nums;
}

.section-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 12px;
}

.section-title-row h3 {
  margin: 0;
  font-size: 18px;
  color: #10261e;
}

.section-title-row.compact {
  margin-bottom: 14px;
}

.subtle {
  font-size: 12px;
  color: #779082;
  line-height: 1.35;
}

.hero-donut-layout {
  display: grid;
  grid-template-columns: minmax(240px, 1fr) 0.85fr;
  gap: 16px;
  align-items: center;
}

.donut-cluster {
  position: relative;
  width: min(320px, 100%);
  aspect-ratio: 1;
  margin: 0 auto;
  display: grid;
  place-items: center;
}

.macro-ring {
  position: absolute;
  border-radius: 50%;
  --pct: 0;
}

.ring-fat {
  width: 100%;
  height: 100%;
  background: conic-gradient(#67c1ff calc(var(--pct) * 1%), #e6f4fb 0);
}

.ring-carb {
  width: 78%;
  height: 78%;
  background: conic-gradient(#ffd05f calc(var(--pct) * 1%), #fcf6e6 0);
}

.ring-protein {
  width: 56%;
  height: 56%;
  background: conic-gradient(#7fc56e calc(var(--pct) * 1%), #eaf6e8 0);
}

.ring-fat::after,
.ring-carb::after,
.ring-protein::after {
  content: '';
  position: absolute;
  inset: 13%;
  border-radius: 50%;
  background: #f6faf5;
}

.donut-center {
  position: relative;
  z-index: 2;
  display: grid;
  justify-items: center;
  gap: 1px;
  color: #10261e;
}

.donut-center strong {
  font-size: clamp(34px, 4vw, 44px);
  line-height: 1;
  font-family: var(--metric-font);
  font-variant-numeric: tabular-nums;
}

.donut-center span {
  font-size: 14px;
  color: #6d887b;
}

.donut-center small {
  font-size: 12px;
  color: #8ba193;
}

.macro-legend {
  display: grid;
  gap: 10px;
}

.legend-item {
  border-radius: 12px;
  padding: 10px 12px;
  background: #ffffff;
  box-shadow: inset 0 0 0 1px #edf3ee;
  display: grid;
  grid-template-columns: 1fr auto auto;
  align-items: center;
  gap: 8px;
}

.legend-item__name {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #1f3b30;
  font-size: 13px;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.dot-protein { background: #7fc56e; }
.dot-carb { background: #ffd05f; }
.dot-fat { background: #67c1ff; }

.legend-item__value,
.legend-item__ratio {
  color: #224236;
  font-weight: 700;
  font-size: 13px;
  font-family: var(--metric-font);
  font-variant-numeric: tabular-nums;
  white-space: nowrap;
}

.bmi-head {
  margin-bottom: 8px;
}

.bmi-label {
  margin: 0;
  color: #607c6f;
  font-size: 13px;
}

.bmi-value {
  margin: 2px 0 0;
  font-size: 34px;
  line-height: 1;
  font-weight: 800;
  color: #10261e;
  font-family: var(--metric-font);
  font-variant-numeric: tabular-nums;
}

.bmi-value.normal { color: #22a06b; }
.bmi-value.underweight { color: #4f9de7; }
.bmi-value.overweight { color: #df9c1e; }
.bmi-value.obese { color: #e45c57; }

.bmi-gauge {
  position: relative;
  height: 40px;
  margin-bottom: 6px;
}

.bmi-track {
  position: absolute;
  left: 0;
  right: 0;
  top: 5px;
  height: 8px;
  border-radius: 999px;
  background: linear-gradient(90deg, #67c1ff 0%, #8edc80 35%, #ffd05f 68%, #f27a72 100%);
}

.bmi-pointer {
  position: absolute;
  top: 0;
  width: 14px;
  height: 14px;
  border-radius: 50%;
  border: 2px solid #fff;
  background: #0f261e;
  transform: translateX(-50%);
  box-shadow: 0 2px 10px rgba(14, 31, 24, 0.25);
}

.bmi-threshold-mark {
  position: absolute;
  top: 2px;
  transform: translateX(-50%);
  display: grid;
  justify-items: center;
  gap: 4px;
}

.bmi-threshold-mark::before {
  content: '';
  width: 1px;
  height: 15px;
  background: rgba(15, 38, 30, 0.28);
}

.bmi-threshold-mark span {
  font-size: 11px;
  color: #6b8578;
  line-height: 1;
  font-family: var(--metric-font);
}

.bmi-scale {
  display: flex;
  justify-content: space-between;
  color: #6f8a7d;
  font-size: 12px;
  margin-bottom: 14px;
}

.vital-mini-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.vital-mini-card {
  border-radius: 12px;
  background: #f9fcf8;
  box-shadow: inset 0 0 0 1px #e9f1ea;
  padding: 12px;
}

.vital-mini-card__label {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #607b6e;
}

.vital-mini-card__value {
  margin-top: 6px;
  font-size: 30px;
  font-weight: 700;
  color: #122a21;
  line-height: 1;
  font-family: var(--metric-font);
  font-variant-numeric: tabular-nums;
}

.vital-mini-card__unit {
  margin-top: 2px;
  font-size: 12px;
  color: #7f9489;
}

.mode-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 10px;
}

.diet-mode-days {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  color: #2b4a3d;
  font-size: 13px;
}

.diet-mode-days :deep(.n-input-number) {
  width: 92px;
}

.segment-control {
  position: relative;
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 4px;
  padding: 4px;
  border-radius: 999px;
  background: #edf4eb;
}

.segment-slider {
  position: absolute;
  top: 4px;
  bottom: 4px;
  border-radius: 999px;
  background: #ffffff;
  box-shadow: 0 6px 14px rgba(18, 42, 33, 0.1);
  transition: left 0.28s ease;
}

.segment-item {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  position: relative;
  z-index: 2;
  border: none;
  background: transparent;
  border-radius: 999px;
  height: 34px;
  line-height: 1;
  color: #39584c;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
}

.segment-item.active {
  color: #153226;
}

.segment-control.disabled,
.segment-item:disabled {
  cursor: not-allowed;
  opacity: 0.7;
}

.mode-explain-item {
  border-radius: 10px;
  padding: 8px 10px;
  background: #f4f8f2;
  box-shadow: inset 0 0 0 1px #e4ede4;
}

.mode-explain-item strong {
  display: block;
  font-size: 12px;
  color: #28473a;
  margin-bottom: 3px;
  line-height: 1.2;
}

.mode-explain-item p {
  margin: 0;
  font-size: 12px;
  line-height: 1.45;
  color: #678277;
}

.mode-explain-item.active {
  background: #eaf5e5;
  box-shadow: inset 0 0 0 1px #cde3c7;
}

.mode-fade-enter-active,
.mode-fade-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.mode-fade-enter-from,
.mode-fade-leave-to {
  opacity: 0;
  transform: translateY(4px);
}

.pref-list {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
}

.pref-item {
  border-radius: 12px;
  padding: 12px;
  background: #f9fcf8;
  box-shadow: inset 0 0 0 1px #e8f0ea;
  display: grid;
  gap: 6px;
}

.pref-item__label {
  color: #6b8478;
  font-size: 12px;
}

.pref-item__value {
  color: #162e24;
  font-size: 14px;
  line-height: 1.5;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 14px;
}

.action-buttons :deep(.n-button) {
  border-radius: 999px;
  min-width: 120px;
}

@media (max-width: 1100px) {
  .profile-card,
  .hero-nutrition-card,
  .vital-card,
  .mode-card,
  .prefs-card {
    grid-column: span 12;
  }

  .hero-donut-layout {
    grid-template-columns: 1fr;
  }

  .pref-list {
    grid-template-columns: 1fr;
  }

  .profile-mini-stats {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .profile-shell {
    padding: 14px 12px 28px;
  }

  .page-header {
    grid-template-columns: 1fr;
    justify-items: start;
    gap: 10px;
  }

  .page-title-wrap {
    justify-items: start;
  }

  .segment-control {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    border-radius: 14px;
  }

  .segment-slider {
    display: none;
  }

  .vital-mini-grid {
    grid-template-columns: 1fr;
  }

  .profile-mini-stats {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    display: grid;
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>

