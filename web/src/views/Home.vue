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

        <div class="recipe-action-bar glass-card" v-if="是否今天">
          <div class="action-info">
            <h3 class="plan-status">{{ 今日方案 ? '今日食谱已生成' : '开始规划健康饮食' }}</h3>
            <p class="plan-desc">{{ 今日方案 ? '按计划饮食，保持活力每一天' : '为您量身定制的营养均衡食谱方案' }}</p>
          </div>
          <button class="commercial-btn" @click="去生成食谱">
            <n-icon><SparklesOutline /></n-icon>
            {{ 今日方案 ? '重新生成食谱' : '智能生成今日食谱' }}
          </button>
        </div>

        <div class="timeline">
          <article
            v-for="meal in 时间线餐次"
            :key="meal.type"
            class="timeline-row"
            :class="{ 'is-drag-over': 拖拽目标餐次 === meal.type, 'is-active-detail': 详情餐次类型 === meal.type }"
            :draggable="可拖拽餐次 === meal.type"
            @dragstart="开始拖拽餐次(meal.type)"
            @dragover.prevent="进入拖拽目标(meal.type)"
            @dragleave="离开拖拽目标(meal.type)"
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

            <div class="timeline-content glass-card meal-card" :class="`meal-${meal.type}`" @click="切换详情餐次(meal.type)">
              <div class="meal-head">
                <div class="meal-title-wrap">
                  <input
                    v-if="编辑名称餐次 === meal.type"
                    v-model="编辑名称内容"
                    class="meal-name-input"
                    maxlength="16"
                    @blur="保存名称编辑(meal.type)"
                    @keyup.enter="保存名称编辑(meal.type)"
                    @click.stop
                  />
                  <h3
                    v-else
                    class="editable-meal-name"
                    title="点击修改名称"
                    @click.stop="开始编辑名称(meal.type, meal.name)"
                  >
                    {{ meal.name }}
                  </h3>
                </div>
                <div class="head-right">
                  <p class="kcal">{{ meal.kcal || '--' }}</p>
                  <div class="drag-handle" title="按住拖拽排序" @mousedown.stop="允许拖拽(meal.type)" @mouseup.stop="禁止拖拽" @mouseleave="禁止拖拽" @click.stop>
                    <n-icon><MenuOutline /></n-icon>
                  </div>
                </div>
              </div>

              <div class="meal-hero">
                <div class="hero-badge">{{ meal.icon }}</div>
              </div>

              <div class="food-hints">
                <span v-for="hint in meal.hints" :key="hint" class="food-chip">{{ hint }}</span>
              </div>

              <div class="bar-track">
                <div class="bar-fill" :style="{ width: meal.progress + '%' }" />
              </div>

              <div class="meal-actions">
                <button class="record-btn-main" @click.stop="打开记录悬浮栏(meal.type)">
                  <n-icon><AddCircleOutline /></n-icon>
                  <span>{{ meal.type === 'snack' ? '记录餐次' : '记录今日餐次' }}</span>
                </button>
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

        <div class="details-card glass-card" :class="{ 'pulse-anim': 闪烁详情动画 }">
          <div class="details-head">
            <h3>{{ 当前详情标题 }}</h3>
            <p>{{ 当前详情时间 }}</p>
          </div>

          <div v-if="!当前详情食物.length" class="sample-note-wrapper">
            <p class="sample-note">
              {{ 今日方案 ? '今日暂无记录，以下为今日推荐食谱' : '今日暂无记录，以下为推荐食谱' }}
            </p>
            <button v-if="是否今天 && !今日方案" class="generate-btn" @click="去生成食谱">
              <n-icon><SparklesOutline /></n-icon> 智能生成今日食谱
            </button>
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
              <span class="check-col"></span>
              <span>食物项</span>
              <span>份量</span>
              <span>热量</span>
              <span>蛋白质</span>
              <span>碳水</span>
              <span>脂肪</span>
              <span class="action-col"></span>
            </div>

            <div v-for="food in 展示详情食物" :key="food.id" class="food-row">
              <label class="check-wrap">
                <input
                  type="checkbox"
                  class="food-check-input"
                  :checked="!food.sample || 本地勾选食物.some(f => f.id === food.id)"
                  :disabled="!food.sample"
                  @change="切换食物状态(food)"
                />
              </label>
              <div class="food-name">
                <span class="thumb" :style="{ background: food.image ? 'transparent' : food.bg }">
                  <img v-if="food.image" :src="food.image" alt="食物图" />
                  <span v-else>{{ food.emoji }}</span>
                </span>
                <span>{{ food.name }}</span>
              </div>
              <span>{{ food.portion }}</span>
              <span class="strong-text">{{ food.kcal }}</span>
              <span class="sub-text">{{ food.protein }}</span>
              <span class="sub-text">{{ food.carbs }}</span>
              <span class="sub-text">{{ food.fat }}</span>
              <div class="action-wrap">
                <n-popconfirm
                  v-if="!food.sample && food.realId"
                  @positive-click="删除真实饮食记录(food.realId)"
                  positive-text="删除"
                  negative-text="取消"
                >
                  <template #trigger>
                    <button class="row-delete-btn" aria-label="删除食物">
                      <n-icon><TrashOutline /></n-icon>
                    </button>
                  </template>
                  确定要删除这道食物记录吗？
                </n-popconfirm>
              </div>
            </div>
          </div>

          <div class="add-more">
            <input type="text" placeholder="为这餐添加更多食物..." @click="打开记录悬浮栏(当前详情餐次.type)" readonly />
            <button class="plus-btn" aria-label="添加食物" @click="打开记录悬浮栏(当前详情餐次.type)">+</button>
          </div>
        </div>
      </section>
    </main>

    <!-- 右侧悬浮栏：快速添加饮食记录 -->
    <n-drawer v-model:show="显示添加记录面板" :width="500" placement="right">
      <n-drawer-content title="快速添加饮食记录" closable>
        <template #header>
          <div class="drawer-header-brand">
            <h3 style="margin: 0; color: #235747; font-size: 1.2rem;">添加健康一餐</h3>
          </div>
        </template>
        
        <n-form :model="记录表单" label-placement="top" size="large">
          <n-grid :cols="2" :x-gap="16">
            <n-gi :span="2">
              <n-form-item label="餐点类型">
                <n-select v-model:value="记录表单.meal_type" :options="餐次选项" placeholder="选择餐点" />
              </n-form-item>
            </n-gi>
            <n-gi :span="2">
              <n-form-item label="食物名称 / 描述">
                <n-input-group>
                  <n-input v-model:value="记录表单.food_name" placeholder="例如：一碗牛肉面，少辣" />
                  <n-button type="primary" color="#8ec662" @click="执行AI估算" :loading="分析中" :disabled="!记录表单.food_name">
                    <template #icon><n-icon><SparklesOutline /></n-icon></template>
                    AI估量
                  </n-button>
                  <n-button type="default" @click="triggerAIUpload">
                    <template #icon><n-icon><CameraOutline /></n-icon></template>
                    拍照
                  </n-button>
                  <input type="file" ref="fileInputRef" accept="image/*" style="display:none" @change="处理图片识别" />
                </n-input-group>
              </n-form-item>
            </n-gi>
            <n-gi :span="2">
              <n-form-item label="分量 (g)">
                <n-input-number v-model:value="记录表单.intake_amount" :step="10" placeholder="估算重量" style="width: 100%" />
              </n-form-item>
            </n-gi>
            
            <n-gi :span="2">
              <n-divider dashed style="color:#5f8578; font-size: 0.85rem;">营养成分 (系统可选填/AI自动填)</n-divider>
            </n-gi>

            <n-gi>
              <n-form-item label="热量">
                <n-input-number v-model:value="记录表单.calculated_energy" placeholder="kcal" :show-button="false">
                  <template #suffix>kcal</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="蛋白质">
                <n-input-number v-model:value="记录表单.calculated_protein" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="碳水">
                <n-input-number v-model:value="记录表单.calculated_carb" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
            <n-gi>
              <n-form-item label="脂肪">
                <n-input-number v-model:value="记录表单.calculated_fat" placeholder="g" :show-button="false">
                  <template #suffix>g</template>
                </n-input-number>
              </n-form-item>
            </n-gi>
          </n-grid>
        </n-form>

        <template #footer>
          <div style="display: flex; gap: 12px; justify-content: flex-end;">
            <n-button @click="显示添加记录面板 = false">暂不添加</n-button>
            <n-button color="#8ec662" @click="保存饮食记录" :loading="保存中" style="color: white; font-weight: 600;">
              确认记录
            </n-button>
          </div>
        </template>
      </n-drawer-content>
    </n-drawer>
    
    <!-- AI 识别 Loading 遮罩 -->
    <div v-if="识别中" class="ai-loading-mask">
      <div class="ai-loading-inner">
        <n-spin size="large" stroke="#8ec662" />
        <p style="color: #245f4d; font-weight: 600; margin-top: 16px;">AI 正在分析食物营养...</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { 
  ref, computed, h, onMounted, onUnmounted, watch, nextTick
} from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { NIcon, useMessage, NDrawer, NDrawerContent, NForm, NFormItem, NInput, NInputGroup, NButton, NInputNumber, NSelect, NDivider, NGrid, NGi, NSpin, NPopconfirm } from 'naive-ui'
import {
  AddCircleOutline,
  CalendarOutline,
  ChevronBackOutline,
  ChevronForwardOutline,
  SparklesOutline,
  MenuOutline,
  CameraOutline,
  FlashOutline,
  TrashOutline
} from '@vicons/ionicons5'
import { getTodayStatus, addIntakeRecord, deleteIntakeRecord } from '@/api/intakeApi'
import { getNutritionRequirements } from '@/api/user'
import { getSelectedRecipePlan, getRecipeRecommendations, selectRecipePlan } from '@/api/recipeApi'
import { recognizeFood, analyzeFoodText } from '@/api/foodRecognitionApi'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

// ---- 侧边栏/添加记录相关的状态与逻辑 ----
const 显示添加记录面板 = ref(false)
const 分析中 = ref(false)
const 识别中 = ref(false)
const 保存中 = ref(false)
const fileInputRef = ref(null)
const ai数据缓存 = ref(null)

const 餐次选项 = [
  { label: '🌅 早餐', value: 'breakfast' },
  { label: '☀️ 午餐', value: 'lunch' },
  { label: '🌙 晚餐', value: 'dinner' },
  { label: '🍎 加餐', value: 'snack' }
]

const 记录表单 = ref({
  meal_type: 'breakfast',
  food_name: '',
  intake_amount: 100,
  calculated_energy: 0,
  calculated_protein: 0,
  calculated_carb: 0,
  calculated_fat: 0
})

const 打开记录悬浮栏 = (mealType) => {
  let mappedType = mealType
  if (!['breakfast', 'lunch', 'dinner', 'snack'].includes(mealType)) {
    mappedType = 'breakfast'
  }
  
  记录表单.value = {
    meal_type: mappedType,
    food_name: '',
    intake_amount: 100,
    calculated_energy: 0,
    calculated_protein: 0,
    calculated_carb: 0,
    calculated_fat: 0
  }
  ai数据缓存.value = null
  显示添加记录面板.value = true
}

const 执行AI估算 = async () => {
  if (!记录表单.value.food_name) return
  分析中.value = true
  try {
    const res = await analyzeFoodText(记录表单.value.food_name)
    记录表单.value.food_name = res.dish_name
    记录表单.value.intake_amount = res.estimated_weight
    记录表单.value.calculated_energy = Math.floor(res.total_calories)
    记录表单.value.calculated_protein = Math.floor(res.total_protein)
    记录表单.value.calculated_carb = Math.floor(res.total_carbs)
    记录表单.value.calculated_fat = Math.floor(res.total_fat)

    ai数据缓存.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }
    message.success('AI文字估算完成')
  } catch (error) {
    message.error('AI估算失败，请手工填写。')
  } finally {
    分析中.value = false
  }
}

const triggerAIUpload = () => {
  fileInputRef.value?.click()
}

const 处理图片识别 = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  
  识别中.value = true
  try {
    const res = await recognizeFood(file)
    记录表单.value.food_name = res.dish_name
    记录表单.value.intake_amount = res.estimated_weight
    记录表单.value.calculated_energy = Math.floor(res.total_calories)
    记录表单.value.calculated_protein = Math.floor(res.total_protein)
    记录表单.value.calculated_carb = Math.floor(res.total_carbs)
    记录表单.value.calculated_fat = Math.floor(res.total_fat)

    ai数据缓存.value = {
      calories: res.calories_per_100g,
      protein: res.protein_per_100g,
      carb: res.carbs_per_100g,
      fat: res.fat_per_100g
    }
    
    if (!显示添加记录面板.value) {
      显示添加记录面板.value = true
    }
    message.success('AI视觉分析完成')
  } catch (e) {
    message.error('图片分析失败，请重试')
  } finally {
    识别中.value = false
    if (fileInputRef.value) {
      fileInputRef.value.value = ''
    }
  }
}

const 保存饮食记录 = async () => {
  if (!记录表单.value.food_name) {
    message.warning('请输入食物名称或描述')
    return
  }
  
  保存中.value = true
  try {
    const isTodaySyncing = 是否今天.value // 记录一下是不是今天
    await addIntakeRecord(记录表单.value)
    message.success('成功记录饮食')
    显示添加记录面板.value = false
    
    // 如果记录的是正在查看的那天，重新拉取看板数据刷新
    if (isTodaySyncing) {
      await 同步看板数据()
    }
  } catch (error) {
    message.error('保存饮食记录失败')
  } finally {
    保存中.value = false
  }
}

const 删除真实饮食记录 = async (recordId) => {
  if (!recordId) return
  
  try {
    加载中.value = true
    const isTodaySyncing = 是否今天.value
    await deleteIntakeRecord(recordId)
    message.success('已删除记录')
    
    if (isTodaySyncing) {
      await 同步看板数据()
    }
  } catch (error) {
    console.error('删除饮食记录失败', error)
    message.error('删除失败，请稍后重试')
  } finally {
    加载中.value = false
  }
}

watch(() => 记录表单.value.intake_amount, (newVal) => {
  if (ai数据缓存.value && newVal > 0) {
    const factor = newVal / 100
    记录表单.value.calculated_energy = Math.round(ai数据缓存.value.calories * factor)
    记录表单.value.calculated_protein = Number((ai数据缓存.value.protein * factor).toFixed(1))
    记录表单.value.calculated_carb = Number((ai数据缓存.value.carb * factor).toFixed(1))
    记录表单.value.calculated_fat = Number((ai数据缓存.value.fat * factor).toFixed(1))
  }
})
// ---------------------------------------------

const 加载中 = ref(false)
const 今日方案 = ref(null)
const 已选日期 = ref(new Date(new Date().getFullYear(), new Date().getMonth(), new Date().getDate()))
const 编辑时间餐次 = ref('')
const 编辑名称餐次 = ref('')
const 编辑名称内容 = ref('')
const 拖拽源餐次 = ref('')
const 拖拽目标餐次 = ref('')
const 可拖拽餐次 = ref('')
const 闪烁详情动画 = ref(false)
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
  { type: 'breakfast', name: '早餐', time: '08:15', icon: '🍳', hero: '🍳🍞', hints: ['🍓', '🥣', '🥛'], ratio: 0.16, emoji: '🍳', bg: 'linear-gradient(135deg, #fdf5d3, #f5e094)' },
  { type: 'lunch', name: '午餐', time: '12:30', icon: '🥗', hero: '🥗🍽️', hints: ['🥦', '🍗', '🍚'], ratio: 0.32, emoji: '🥗', bg: 'linear-gradient(135deg, #e5f5e5, #b8e0b8)' },
  { type: 'snack', name: '下午加餐', time: '16:00', icon: '🫐', hero: '🥜🍎', hints: ['🥜', '🍓', '🍵'], ratio: 0.2, emoji: '🍎', bg: 'linear-gradient(135deg, #fdf5d3, #f5e094)' },
  { type: 'dinner', name: '晚餐', time: '19:30', icon: '🥩', hero: '🥩🍷', hints: ['🥔', '🥬', '🍷'], ratio: 0.32, emoji: '🍽️', bg: 'linear-gradient(135deg, #dff0de, #a7d7a7)' }
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

const 获取计划餐次食谱列表 = (plan, mealType) => {
  if (!plan || !mealType) return []

  const itemsKey = `${mealType}_items`
  if (Array.isArray(plan[itemsKey]) && plan[itemsKey].length > 0) {
    return plan[itemsKey].filter((item) => item && item.name)
  }

  const singleRecipe = plan[mealType]
  if (singleRecipe && singleRecipe.name) {
    return [singleRecipe]
  }

  return []
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

const 允许拖拽 = (type) => {
  可拖拽餐次.value = type
}

const 禁止拖拽 = () => {
  可拖拽餐次.value = ''
}

const 进入拖拽目标 = (mealType) => {
  if (拖拽源餐次.value && 拖拽源餐次.value !== mealType) {
    拖拽目标餐次.value = mealType
  }
}

const 离开拖拽目标 = (mealType) => {
  if (拖拽目标餐次.value === mealType) {
    拖拽目标餐次.value = ''
  }
}

const 放置餐次 = (targetType) => {
  const sourceType = 拖拽源餐次.value
  拖拽目标餐次.value = ''
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
  拖拽目标餐次.value = ''
  可拖拽餐次.value = ''
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

const 本地勾选食物 = ref([])
const 切换食物状态 = (food) => {
  if (!food.sample) return
  const parsedKcal = parseFloat(food.kcal) || 0
  const parsedProtein = parseFloat(food.protein) || 0
  const parsedCarbs = parseFloat(food.carbs) || 0
  const parsedFat = parseFloat(food.fat) || 0
  const idx = 本地勾选食物.value.findIndex(f => f.id === food.id)
  if (idx !== -1) {
    本地勾选食物.value.splice(idx, 1)
  } else {
    本地勾选食物.value.push({ id: food.id, kcal: parsedKcal, protein: parsedProtein, carbs: parsedCarbs, fat: parsedFat })
  }
}

const 本地勾选摄入总量 = computed(() => {
  return 本地勾选食物.value.reduce((sum, item) => sum + item.kcal, 0)
})
const 本地勾选碳水总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.carbs, 0))
const 本地勾选蛋白质总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.protein, 0))
const 本地勾选脂肪总量 = computed(() => 本地勾选食物.value.reduce((sum, item) => sum + item.fat, 0))

const 每日目标热量 = computed(() => Math.round(安全数值(每日状态.value.targetEnergy)))
const 已摄入热量 = computed(() => Math.round(安全数值(每日状态.value.totalEnergy)) + 本地勾选摄入总量.value)

const 真实剩余热量 = computed(() => Math.max(0, 每日目标热量.value - 已摄入热量.value))
const 真实进度条 = computed(() => {
  const goal = Math.max(1, 每日目标热量.value)
  return Math.min(100, Math.round((已摄入热量.value / goal) * 100))
})

const 动画剩余热量 = ref(0)
const 动画进度条 = ref(0)
let 动画标识 = null

watch([真实剩余热量, 真实进度条], ([新热量, 新进度]) => {
  const 起始热量 = 动画剩余热量.value
  const 起始进度 = 动画进度条.value
  const 差值热量 = 新热量 - 起始热量
  const 差值进度 = 新进度 - 起始进度
  
  if (差值热量 === 0 && 差值进度 === 0 && 动画剩余热量.value !== 0) return

  const 开始时间 = performance.now()
  const 持续时间 = 600

  if (动画标识) cancelAnimationFrame(动画标识)

  const 执行动画 = (当前时间) => {
    const 经过时间 = 当前时间 - 开始时间
    const 进度 = Math.min(经过时间 / 持续时间, 1)
    const 缓动 = 1 - Math.pow(1 - 进度, 4) // easeOutQuart 缓动函数
    
    动画剩余热量.value = 起始热量 + 差值热量 * 缓动
    动画进度条.value = 起始进度 + 差值进度 * 缓动

    if (进度 < 1) {
      动画标识 = requestAnimationFrame(执行动画)
    } else {
      动画剩余热量.value = 新热量
      动画进度条.value = 新进度
    }
  }
  动画标识 = requestAnimationFrame(执行动画)
}, { immediate: true })

const 剩余热量 = computed(() => Math.round(动画剩余热量.value).toLocaleString('zh-CN'))

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

    let defaultHints = []
    if (是否今天.value && 今日方案.value) {
      const planItems = 获取计划餐次食谱列表(今日方案.value, config.type)
      if (planItems.length > 0) {
        defaultHints = planItems.map((item) => item.name).filter(Boolean)
      }
    }

    if (defaultHints.length === 0) {
      defaultHints = (示例食谱映射[config.type] || []).map(item => item.name)
    } else {
      defaultHints = defaultHints.slice(0, 4)
    }

    let hints = []
    if (records.length > 0) {
      const addedNames = records.map(r => r.foodName || r.food_name).filter(Boolean)
      hints = Array.from(new Set([...addedNames, ...defaultHints])).slice(0, 4)
    } else {
      hints = defaultHints
    }

    return {
      ...config,
      time: 餐次时间映射.value[config.type] || config.time,
      records,
      consumed,
      target,
      hints,
      kcal: consumed > 0 ? `${consumed} / ${target} 千卡` : '',
      progress,
      variant: consumed > 0 ? (config.type === 'lunch' ? 'logged' : 'filled') : 'empty'
    }
  }).filter(Boolean)
})

const 宏量卡片 = computed(() => {
  const totalCarbs = Math.round(安全数值(每日状态.value.totalCarbohydrate) + 本地勾选碳水总量.value)
  const totalProtein = Math.round(安全数值(每日状态.value.totalProtein) + 本地勾选蛋白质总量.value)
  const totalFat = Math.round(安全数值(每日状态.value.totalFat) + 本地勾选脂肪总量.value)

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
  let protein = records.reduce((sum, item) => sum + 安全数值(item.calculatedProtein || item.calculated_protein), 0)
  let carbs = records.reduce((sum, item) => sum + 安全数值(item.calculatedCarb || item.calculated_carb), 0)
  let fat = records.reduce((sum, item) => sum + 安全数值(item.calculatedFat || item.calculated_fat), 0)

  if (展示详情食物.value) {
    const checkedSample = 展示详情食物.value.filter(food => food.sample && 本地勾选食物.value.some(f => f.id === food.id))
    protein += checkedSample.reduce((sum, item) => sum + (parseFloat(item.protein) || 0), 0)
    carbs += checkedSample.reduce((sum, item) => sum + (parseFloat(item.carbs) || 0), 0)
    fat += checkedSample.reduce((sum, item) => sum + (parseFloat(item.fat) || 0), 0)
  }

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
  return (meal?.records || []).map((item, index) => {
    const foodName = item.foodName || item.food_name || '未知食物'
    return {
      id: item.id || item.ID || `${meal.type}-${index}`,
      realId: item.id || item.ID, // 提取真实 ID
      name: foodName,
      portion: `${Math.round(安全数值(item.intakeAmount || item.intake_amount))}g`,
      kcal: `${Math.round(安全数值(item.calculatedEnergy || item.calculated_energy))} 千卡`,
      protein: `${Math.round(安全数值(item.calculatedProtein || item.calculated_protein))}g`,
      carbs: `${Math.round(安全数值(item.calculatedCarb || item.calculated_carb))}g`,
      fat: `${Math.round(安全数值(item.calculatedFat || item.calculated_fat))}g`,
      emoji: fallback.emoji,
      bg: fallback.bg,
      image: item.imageUrl || item.image_url || 食物图片映射[foodName] || defaultFoodThumb,
      sample: false
    }
  })
})

const 食物图片映射 = {
  '牛油果全麦吐司': 'https://images.unsplash.com/photo-1588137378633-dea1336ce1e2?auto=format&fit=crop&w=150&q=80',
  '蓝莓酸奶碗': 'https://images.unsplash.com/photo-1493770348161-369560ae357d?auto=format&fit=crop&w=150&q=80',
  '鸡胸肉藜麦沙拉': 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&w=150&q=80',
  '南瓜浓汤': 'https://images.unsplash.com/photo-1476718406336-bb5a9690ee2a?auto=format&fit=crop&w=150&q=80',
  '混合坚果': 'https://images.unsplash.com/photo-1536591375315-e2343b67bf96?auto=format&fit=crop&w=150&q=80',
  '苹果切片': 'https://images.unsplash.com/photo-1568702846914-9661f04d9c73?auto=format&fit=crop&w=150&q=80',
  '香煎三文鱼配芦笋': 'https://images.unsplash.com/photo-1467003909585-2f8a72700288?auto=format&fit=crop&w=150&q=80',
  '黑椒菌菇意面': 'https://images.unsplash.com/photo-1473093295043-cdd812d0e601?auto=format&fit=crop&w=150&q=80'
}
const defaultFoodThumb = 'https://images.unsplash.com/photo-1490645935967-10de6ba17061?auto=format&fit=crop&w=150&q=80'

const 示例食谱映射 = {
  breakfast: [
    { name: '牛油果全麦吐司', portion: '130g', kcal: '285 千卡', protein: '8g', carbs: '32g', fat: '15g', emoji: '🥑' },
    { name: '蓝莓酸奶碗', portion: '180g', kcal: '210 千卡', protein: '12g', carbs: '24g', fat: '7g', emoji: '🫐' }
  ],
  lunch: [
    { name: '鸡胸肉藜麦沙拉', portion: '260g', kcal: '420 千卡', protein: '42g', carbs: '38g', fat: '12g', emoji: '🥗' },
    { name: '南瓜浓汤', portion: '220g', kcal: '165 千卡', protein: '4g', carbs: '28g', fat: '5g', emoji: '🎃' }
  ],
  snack: [
    { name: '混合坚果', portion: '35g', kcal: '195 千卡', protein: '7g', carbs: '9g', fat: '16g', emoji: '🥜' },
    { name: '苹果切片', portion: '120g', kcal: '62 千卡', protein: '0g', carbs: '14g', fat: '0g', emoji: '🍎' }
  ],
  dinner: [
    { name: '香煎三文鱼配芦笋', portion: '240g', kcal: '468 千卡', protein: '38g', carbs: '8g', fat: '28g', emoji: '🐟' },
    { name: '黑椒菌菇意面', portion: '210g', kcal: '338 千卡', protein: '14g', carbs: '52g', fat: '8g', emoji: '🍝' }
  ]
}

const 示例详情食物 = computed(() => {
  const meal = 当前详情餐次.value
  let list = []

  if (是否今天.value && 今日方案.value) {
    const planItems = 获取计划餐次食谱列表(今日方案.value, meal?.type)
    if (planItems.length > 0) {
      list = planItems.map((recipe) => ({
        name: recipe.name,
        portion: '1份',
        kcal: `${Math.round(recipe.energy || 0)} 千卡`,
        protein: `${Math.round(recipe.protein || 0)}g`,
        carbs: `${Math.round(recipe.carbohydrate || 0)}g`,
        fat: `${Math.round(recipe.fat || 0)}g`,
        emoji: meal?.emoji || '🍽️',
        image: recipe.imageUrl || recipe.image_url || defaultFoodThumb
      }))
    }
  }

  if (!list.length) {
    list = 示例食谱映射[meal?.type] || []
  }

  return list.map((item, index) => ({
    id: `sample-${meal?.type || 'meal'}-${index}`,
    name: item.name,
    portion: item.portion,
    kcal: item.kcal,
    protein: item.protein,
    carbs: item.carbs,
    fat: item.fat,
    emoji: item.emoji || meal?.emoji || '🍽️',
    bg: meal?.bg || 'linear-gradient(135deg, #7aa694, #4f7a69)',
    image: item.image || 食物图片映射[item.name] || defaultFoodThumb,
    sample: true
  }))
})

const 展示详情食物 = computed(() => {
  const actual = 当前详情食物.value || []
  const samples = 示例详情食物.value || []
  
  // 保留原有的示例食谱样貌，每次添加的食物追加/叠加在前面，保持 UI 不会空洞
  const actualNames = new Set(actual.map(item => item.name))
  const filteredSamples = samples.filter(item => !actualNames.has(item.name))

  return [...actual, ...filteredSamples]
})

const ringStyle = computed(() => {
  const progress = 动画进度条.value
  return {
    background: `conic-gradient(#8ec662 0 ${progress}%, #e7f2da ${progress}% 100%)`
  }
})

const 同步看板数据 = async () => {
  加载中.value = true
  try {
    if (authStore.isAuthenticated && !authStore.profile) {
      await authStore.loadProfile()
    }

    const promises = [
      getTodayStatus(转日期参数(已选日期.value)),
      getNutritionRequirements().catch(() => null)
    ]
    if (是否今天.value) {
      promises.push(getSelectedRecipePlan().catch(() => null))
    }

    const res = await Promise.all(promises)
    const status = res[0]
    const nutrition = res[1]
    const plan = res[2]

    if (是否今天.value && plan && plan.id) {
      今日方案.value = plan
    } else {
      今日方案.value = null
    }

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
    本地勾选食物.value = []
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
  if (详情餐次类型.value === mealType) return
  详情餐次类型.value = mealType
  闪烁详情动画.value = true
  setTimeout(() => { 闪烁详情动画.value = false }, 350)
}

const 去生成食谱 = async () => {
  if (!是否今天.value) return
  加载中.value = true
  try {
    const res = await getRecipeRecommendations(1)
    if (res && res.plans && res.plans.length > 0) {
      await selectRecipePlan(res.plans[0])
      message.success('已为您生成新的今日食谱')
      await 同步看板数据()
    } else {
      message.warning('暂无合适的食谱推荐')
    }
  } catch (error) {
    console.error('智能生成食谱失败', error)
    message.error('生成食谱失败，请稍后重试')
  } finally {
    加载中.value = false
  }
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
  padding: 24px;
  position: relative;
  overflow: hidden;
  /* 亲自然：淡色背景加微光感，模拟阳光洒落 */
  background-color: #f4f6f3;
  background-image: 
    radial-gradient(circle at 15% 10%, rgba(255, 255, 255, 0.9) 0%, transparent 40%),
    radial-gradient(circle at 85% 90%, rgba(240, 246, 237, 0.7) 0%, transparent 45%);
  color: #1a2f24; /* 高对比深绿植色 */
}

.glass-card {
  background: #ffffff;
  border: none; /* 移除生硬边框，改用自然光影 */
  border-radius: 20px; /* 有机柔和的圆角 */
  box-shadow: 0 10px 40px -10px rgba(32, 59, 45, 0.06); /* 柔和植物阴影 */
  transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1); /* 150-220ms的自然节奏 */
}

.glass-card:hover {
  /* Hover时如光线移动，亮度微升，轻微上浮 */
  transform: translateY(-2px);
  box-shadow: 0 16px 48px -12px rgba(32, 59, 45, 0.1);
}

.top-header,
.center-panel,
.timeline-content,
.overview-card,
.details-card,
.food-table {
  border-radius: 16px;
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
  font-size: 2.2rem;
  letter-spacing: 0.03em;
  color: #184c3d;
  font-weight: 800;
}

.brand-date {
  margin: 0;
  font-size: 1.15rem;
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
  border: none;
  border-radius: 50%;
  background: linear-gradient(145deg, #ffffff 0%, #f2f7f2 100%);
  color: #2b5643; /* 深叶绿 */
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
  box-shadow: 0 4px 14px rgba(32, 59, 45, 0.08); /* 柔和植物影 */
}

.icon-circle:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 22px rgba(32, 59, 45, 0.12); /* 偏光效果 */
}

.icon-circle:active {
  transform: translateY(1px) scale(0.96); /* 下沉动效 */
  box-shadow: 0 2px 8px rgba(32, 59, 45, 0.06);
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
  background: transparent;
  border: none;
  padding: 0;
  backdrop-filter: none;
}

.loading-tip {
  margin: 0 0 12px;
  color: #5d8778;
  font-size: 1rem;
  font-weight: 600;
}

.record-ops {
  margin-bottom: 24px;
  padding: 12px 18px;
  border: 1px solid #e9f0eb;
  border-radius: 12px;
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
  color: #2b5745;
  font-size: 1rem;
  font-weight: 700;
}

.ops-subtitle {
  margin: 0 0 0 6px;
  color: #7b9f8f;
  font-size: 0.95rem;
}

.calendar-picker {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #2a6552;
  font-size: 1rem;
  font-weight: 700;
}

.calendar-picker input {
  border: 1px solid rgba(45, 112, 90, 0.2);
  background: #f5fbf8;
  color: #2f6554;
  height: 34px;
  padding: 0 8px;
}

.recipe-action-bar {
  margin-bottom: 24px;
  padding: 20px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-left: 6px solid #8ec662;
}

.action-info .plan-status {
  margin: 0 0 4px;
  font-size: 1.1rem;
  color: #2c3e35;
  font-weight: 700;
}

.action-info .plan-desc {
  margin: 0;
  font-size: 0.85rem;
  color: #6c8a7b;
}

.commercial-btn {
  border: none;
  /* 亲自然：叶绿自然渐变 */
  background: linear-gradient(135deg, #8dc75d 0%, #72aa46 100%);
  color: #ffffff;
  padding: 10px 20px;
  border-radius: 20px; /* 有机圆角 */
  font-size: 1rem;
  font-weight: 600;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  box-shadow: 0 6px 16px rgba(122, 180, 77, 0.25);
  transition: all 0.18s cubic-bezier(0.2, 0.8, 0.2, 1); /* 在 150-220ms 范围内 */
}

.commercial-btn:hover {
  background: linear-gradient(135deg, #95cc68 0%, #79b04d 100%); /* 亮度微升 */
  transform: translateY(-2px); /* 悬浮感 */
  box-shadow: 0 10px 24px rgba(122, 180, 77, 0.4); 
}

.commercial-btn:active {
  transform: translateY(1px) scale(0.97); /* Active 轻微下沉 */
  box-shadow: 0 2px 8px rgba(122, 180, 77, 0.15); /* 阴影收拢 */
}

.timeline {
  padding-left: 22px;
  border-left: none;
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: relative;
}

.timeline::before {
  content: '';
  position: absolute;
  left: 3px;
  top: 0;
  bottom: 0;
  width: 2px;
  background: #ddead1; /* 亲自然：植物茎秆色 */
}

.timeline-row {
  display: grid;
  grid-template-columns: 98px minmax(280px, 1fr);
  gap: 12px;
  position: relative;
  transition: transform 0.22s cubic-bezier(0.2, 0.8, 0.2, 1);
}

.timeline-row:hover {
  transform: translateX(4px) scale(1.005);
}

.timeline-row.is-drag-over .timeline-content {
  border: 2px dashed #8dc75d;
  background: #f4faf0;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(122, 180, 77, 0.15);
}

.timeline-row.is-active-detail .timeline-content {
  box-shadow: 0 0 0 2px rgba(132, 189, 87, 0.4); /* 植物环绕感 */
}

.time-col {
  padding-top: 10px;
  position: relative;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: #ffffff;
  position: absolute;
  left: -23px;
  top: 15px;
  border: 2px solid #8ec662;
  box-shadow: 0 0 0 4px #f8faf9;
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
  cursor: pointer;
}

.meal-card::before,
.meal-card::after {
  display: none;
}

.meal-breakfast {
  background: #ffffff;
  border-left: 6px solid #f5e094;
}

.meal-lunch {
  background: #ffffff;
  border-left: 6px solid #8ec662;
}

.meal-snack {
  background: #ffffff;
  border-left: 6px solid #f5e094;
}

.meal-dinner {
  background: #ffffff;
  border-left: 6px solid #8ec662;
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
  font-size: 1.35rem;
  color: #204c3f;
}

.editable-meal-name {
  cursor: pointer;
  transition: color 0.2s ease;
}

.meal-name-input {
  width: 150px;
  border: 1px solid rgba(47, 124, 99, 0.28);
  border-radius: 8px;
  background: #f4faf6;
  color: #2c6251;
  padding: 5px 8px;
  font-size: 1.1rem;
  font-weight: 700;
}

.head-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.drag-handle {
  cursor: grab;
  font-size: 1.25rem;
  color: #a0bcae;
  display: flex;
  padding: 4px;
  border-radius: 4px;
  transition: background 0.2s, color 0.2s;
}

.drag-handle:hover {
  background: #f0f5f2;
  color: #5f8a7b;
}

.drag-handle:active {
  cursor: grabbing;
}

.kcal {
  margin: 0;
  color: #4d7c6d;
  font-weight: 700;
  font-size: 1rem;
}

.meal-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.hero-badge {
  width: 46px;
  height: 46px;
  border-radius: 12px;
  background: white;
  display: grid;
  place-items: center;
  font-size: 1.6rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
}

.food-hints {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 12px;
}

.food-chip {
  padding: 6px 12px;
  border-radius: 16px;
  background: #f4faf0;
  display: inline-flex;
  align-items: center;
  font-size: 1rem;
  font-weight: 600;
  color: #3d6a4e;
  border: 1px solid #ddecce;
  white-space: nowrap;
}

.bar-track {
  height: 8px;
  border-radius: 999px;
  background: #f0f5f2;
  overflow: hidden;
  margin-bottom: 12px;
}

.bar-fill {
  height: 100%;
  border-radius: inherit;
  background: #8ec662;
}

.meal-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.record-btn-main {
  border: 1px solid #ddecce;
  border-radius: 8px;
  min-height: 38px;
  padding: 0 16px;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #4a7c59;
  font-weight: 600;
  background: #f4faf0;
  box-shadow: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.record-btn-main:hover {
  background: #e7f2da;
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
  font-size: 2.2rem;
  color: #1f5a46;
}

.ring-center p {
  margin: 2px 0 0;
  font-size: 0.95rem;
  color: #668f80;
  font-weight: 700;
}

.goal-text {
  margin: 8px 0 0;
  font-size: 1.1rem;
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
  border-radius: 12px;
  background: #ffffff;
  border: 1px solid #e9f0eb;
  padding: 12px 14px;
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
  font-size: 1.1rem;
}

.macro-target {
     margin: 2px 0 0;
  color: #6d9384;
  font-size: 0.9rem;
}

.macro-value {
  margin: 0;
  color: #1f5947;
  font-weight: 800;
  font-size: 1.3rem;
}

.details-card {
  padding: 16px;
  position: relative;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 520px;
}

@keyframes detail-pulse {
  0% { opacity: 0.8; transform: scale(0.995); }
  50% { opacity: 1; transform: scale(1.008); box-shadow: 0 12px 32px rgba(142, 198, 98, 0.15); border-color: #8ec662; }
  100% { opacity: 1; transform: scale(1); }
}

.pulse-anim {
  animation: detail-pulse 0.35s cubic-bezier(0.2, 0.8, 0.2, 1);
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
  font-size: 1.8rem;
}

.details-head p {
  margin: 0;
  color: #6b9283;
  font-size: 1rem;
}

.sample-note-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  color: #6b9283;
  font-size: 0.95rem;
}

.sample-note {
  margin: 0;
  color: #5f8779;
  font-size: 0.82rem;
  font-weight: 700;
}

.generate-btn {
  border: none;
  background: #8ec662;
  color: white;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.82rem;
  display: flex;
  align-items: center;
  gap: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 600;
  box-shadow: 0 4px 8px rgba(142, 198, 98, 0.2);
}

.generate-btn:hover {
  transform: translateY(-1px);
  background: #7bb550;
  box-shadow: 0 6px 12px rgba(142, 198, 98, 0.3);
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
  font-size: 0.95rem;
  font-weight: 700;
  margin-bottom: 6px;
}

.micro-track {
  height: 6px;
  border-radius: 999px;
  overflow: hidden;
  background: #f0f5f2;
}

.micro-fill {
  height: 100%;
}

.food-table {
  border-radius: 12px;
  background: #ffffff;
  border: 1px solid #e9f0eb;
  padding: 10px 12px;
  flex: 1;
}

.table-head,
.food-row {
  display: grid;
  grid-template-columns: 28px 1.9fr 0.6fr 0.8fr 0.6fr 0.6fr 0.6fr 32px;
  align-items: center;
  gap: 8px;
}

.check-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.food-check-input {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: #8ec662;
}

.table-head {
  color: #709384;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 6px;
  font-weight: 700;
}

.food-row {
  height: 72px;
  border-top: 1px solid #f0f5f2;
  font-size: 1.05rem;
  color: #2c3e35;
}

.action-col {
  text-align: center;
}

.action-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
}

.row-delete-btn {
  background: transparent;
  border: none;
  color: #a0bcae;
  font-size: 1.3rem;
  cursor: pointer;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.row-delete-btn:hover {
  background: #fdf5f5;
  color: #e12f3c;
}

.food-empty {
  color: #7d9f92;
}

.food-name {
  display: flex;
  align-items: center;
  gap: 12px;
  font-weight: 600;
  font-size: 1.05rem;
}

.thumb {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: grid;
  place-items: center;
  color: #f4fffb;
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
}

.thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.strong-text {
  font-weight: 800;
}

.sub-text {
  font-size: 0.9rem;
  color: #7b9f8f;
}

.unit-text {
  font-size: 0.85rem;
  color: #9cbcae;
  margin-left: 2px;
}

/* AI 识别 Loading 遮罩 */
.ai-loading-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(244, 246, 243, 0.85); /* 亲自然半透 */
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(6px);
}

.ai-loading-inner {
  text-align: center;
  background: white;
  padding: 32px 48px;
  border-radius: 20px;
  box-shadow: 0 16px 48px rgba(32, 59, 45, 0.1);
}
</style>
