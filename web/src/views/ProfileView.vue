<template>
  <div class="profile-view-container">
    <TopNavigation />
    
    <div class="page-header">
      <BackButton />
      <h1 class="page-title">个人档案</h1>
      <n-button type="primary" @click="goToEdit">
        <template #icon>
          <n-icon><CreateOutline /></n-icon>
        </template>
        更新档案
      </n-button>
    </div>

    <div class="main-content">
      <div v-if="loading" class="loading-container">
        <n-skeleton text :repeat="10" />
      </div>

      <div v-else-if="!authStore.hasProfile" class="empty-state">
        <n-empty description="您还没有完善健康档案">
          <template #extra>
            <n-button type="primary" @click="goToEdit">立即完善档案</n-button>
          </template>
        </n-empty>
      </div>

      <div v-else class="profile-content">
        <!-- 顶部：营养需求（最重要，突出显示） -->
        <div class="top-section" v-if="nutritionData">
          <n-card class="nutrition-card highlight-card">
            <template #header>
              <div class="card-header">
                <span class="card-icon">🔥</span>
                <span>每日营养需求</span>
              </div>
            </template>
            <div class="nutrition-content">
              <div class="target-calorie">
                <div class="target-label">目标热量</div>
                <div class="target-value">{{ nutritionData.target_calorie }} <span class="unit">kcal</span></div>
              </div>
              <div class="macros-grid">
                <div class="macro-item protein">
                  <div class="macro-icon">🥩</div>
                  <div class="macro-info">
                    <div class="macro-name">蛋白质</div>
                    <div class="macro-value">{{ nutritionData.protein_gram }}g</div>
                    <div class="macro-ratio">{{ nutritionData.protein_ratio.toFixed(0) }}%</div>
                  </div>
                </div>
                <div class="macro-item carb">
                  <div class="macro-icon">🍚</div>
                  <div class="macro-info">
                    <div class="macro-name">碳水化合物</div>
                    <div class="macro-value">{{ nutritionData.carb_gram }}g</div>
                    <div class="macro-ratio">{{ nutritionData.carb_ratio.toFixed(0) }}%</div>
                  </div>
                </div>
                <div class="macro-item fat">
                  <div class="macro-icon">🥑</div>
                  <div class="macro-info">
                    <div class="macro-name">脂肪</div>
                    <div class="macro-value">{{ nutritionData.fat_gram }}g</div>
                    <div class="macro-ratio">{{ nutritionData.fat_ratio.toFixed(0) }}%</div>
                  </div>
                </div>
              </div>
            </div>
          </n-card>
        </div>

        <!-- 中间：两列布局 - 基础信息和健康指标 -->
        <div class="middle-section">
          <n-card class="info-card">
            <template #header>
              <div class="card-header">
                <span class="card-icon">👤</span>
                <span>基础信息</span>
              </div>
            </template>
            <n-descriptions :column="2" bordered label-placement="left">
              <n-descriptions-item label="用户名">
                <n-tag type="primary" round>{{ authStore.user?.username }}</n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="邮箱">
                {{ authStore.profile?.email || '未设置' }}
              </n-descriptions-item>
              <n-descriptions-item label="性别">
                <n-tag :type="authStore.profile?.gender === '男' ? 'info' : 'error'" round>
                  {{ authStore.profile?.gender }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="年龄">
                {{ authStore.profile?.age }} 岁
              </n-descriptions-item>
              <n-descriptions-item label="身高">
                {{ authStore.profile?.height }} cm
              </n-descriptions-item>
              <n-descriptions-item label="体重">
                {{ authStore.profile?.weight }} kg
              </n-descriptions-item>
            </n-descriptions>
          </n-card>

          <n-card class="info-card">
            <template #header>
              <div class="card-header">
                <span class="card-icon">📊</span>
                <span>健康指标</span>
              </div>
            </template>
            <div class="health-metrics">
              <div class="metric-item">
                <div class="metric-label">BMI</div>
                <div class="metric-value" :class="getBMIClass(authStore.profile?.bmi)">
                  {{ authStore.profile?.bmi || '--' }}
                </div>
                <div class="metric-desc">{{ getBMIDesc(authStore.profile?.bmi) }}</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">BMR</div>
                <div class="metric-value primary">{{ authStore.profile?.bmr || '--' }}</div>
                <div class="metric-desc">基础代谢率</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">TDEE</div>
                <div class="metric-value primary">{{ authStore.profile?.tdee || '--' }}</div>
                <div class="metric-desc">每日总消耗</div>
              </div>
            </div>
          </n-card>
        </div>

        <!-- 底部：两列布局 - 健康目标和个性化设置 -->
        <div class="bottom-section">
          <n-card class="info-card">
            <template #header>
              <div class="card-header">
                <span class="card-icon">🎯</span>
                <span>健康目标</span>
              </div>
            </template>
            <n-descriptions :column="1" bordered label-placement="left">
              <n-descriptions-item label="健康目标">
                <n-tag :type="getGoalTagType(authStore.profile?.health_goal)" round size="large">
                  {{ authStore.profile?.health_goal || '未设置' }}
                </n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="目标体重">
                {{ authStore.profile?.target_weight ? `${authStore.profile.target_weight} kg` : '未设置' }}
              </n-descriptions-item>
              <n-descriptions-item label="活动水平">
                <n-tag type="info" round>{{ authStore.profile?.activity_level || '未设置' }}</n-tag>
              </n-descriptions-item>
              <n-descriptions-item label="每日用餐次数">
                {{ authStore.profile?.meal_times_per_day || 3 }} 次
              </n-descriptions-item>
            </n-descriptions>
          </n-card>

          <n-card class="info-card">
            <template #header>
              <div class="card-header">
                <span class="card-icon">⚙️</span>
                <span>个性化设置</span>
              </div>
            </template>
            <n-descriptions :column="1" bordered label-placement="left">
              <n-descriptions-item label="过敏源">
                <div class="text-content">
                  {{ authStore.profile?.allergies || '无' }}
                </div>
              </n-descriptions-item>
              <n-descriptions-item label="饮食偏好">
                <div class="text-content">
                  {{ authStore.profile?.dietary_prefs || '无' }}
                </div>
              </n-descriptions-item>
              <n-descriptions-item label="健康问题">
                <div class="text-content">
                  {{ authStore.profile?.health_conditions || '无' }}
                </div>
              </n-descriptions-item>
            </n-descriptions>
          </n-card>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons">
          <n-button type="primary" size="large" @click="goToEdit">
            <template #icon>
              <n-icon><CreateOutline /></n-icon>
            </template>
            更新档案
          </n-button>
          <n-button size="large" @click="goBack">
            返回首页
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NButton, NIcon, NSkeleton, NEmpty, NDescriptions, NDescriptionsItem, NTag } from 'naive-ui'
import { CreateOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/store/auth'
import { getNutritionRequirements } from '@/api/user'
import TopNavigation from '@/components/layout/TopNavigation.vue'
import BackButton from '@/components/layout/BackButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const nutritionData = ref(null)

onMounted(async () => {
  loading.value = true
  try {
    if (!authStore.profile) {
      await authStore.loadProfile()
    }
    // 加载营养需求数据
    if (authStore.hasProfile) {
      try {
        const response = await getNutritionRequirements()
        nutritionData.value = response
      } catch (error) {
        // 静默处理错误，不显示错误提示
        // 如果档案不完整或接口不存在，只是不显示营养需求卡片
        const errorMsg = error?.response?.data?.error || error.message
        const status = error?.response?.status
        console.log('营养需求数据不可用:', {
          status,
          error: errorMsg,
          url: '/api/v1/user/nutrition'
        })
        nutritionData.value = null
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

const getGoalTagType = (goal) => {
  const goalMap = {
    '减脂': 'error',
    '增肌': 'success',
    '控糖': 'warning',
    '维持健康': 'info'
  }
  return goalMap[goal] || 'default'
}
</script>

<style scoped>
.profile-view-container {
  min-height: 100vh;
  background: var(--bg-secondary);
}

.page-header {
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--spacing-lg) var(--spacing-xl);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin: 0;
  flex: 1;
  text-align: center;
}

.main-content {
  padding: 0 var(--spacing-xl) var(--spacing-3xl);
  max-width: 1400px;
  margin: 0 auto;
}

/* 响应式布局 */
@media (max-width: 1024px) {
  .middle-section,
  .bottom-section {
    grid-template-columns: 1fr;
  }
  
  .health-metrics {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .macros-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .health-metrics {
    grid-template-columns: 1fr;
  }
  
  .macros-grid {
    grid-template-columns: 1fr;
  }
  
  .main-content {
    padding: 20px;
  }
}

.loading-container {
  padding: var(--spacing-3xl);
}

.empty-state {
  padding: var(--spacing-4xl) var(--spacing-lg);
}

.profile-content {
  animation: fadeIn 0.5s ease-out;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-lg);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 顶部营养需求区域 */
.top-section {
  width: 100%;
}

.nutrition-card {
  border-radius: var(--radius-xl);
  overflow: hidden;
}

.highlight-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9ff 100%);
  box-shadow: var(--shadow-lg);
}

/* 中间区域：两列布局 */
.middle-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-lg);
}

/* 底部区域：两列布局 */
.bottom-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-lg);
}

.info-card {
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: var(--transition-normal);
}

.card-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
}

.card-icon {
  font-size: 24px;
}

.health-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-md);
  padding: var(--spacing-lg) 0;
}

.metric-item {
  text-align: center;
  padding: var(--spacing-md);
  background: linear-gradient(135deg, var(--bg-secondary) 0%, #ffffff 100%);
  border-radius: var(--radius-lg);
  border: 1px solid var(--border-color);
  transition: var(--transition-normal);
}

.metric-item:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.metric-label {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-sm);
  font-weight: var(--font-weight-medium);
}

.metric-value {
  font-size: 28px;
  font-weight: var(--font-weight-bold);
  margin-bottom: var(--spacing-xs);
  color: var(--text-primary);
}

.metric-value.primary {
  color: var(--color-primary);
}

.metric-value.normal {
  color: var(--color-success);
}

.metric-value.underweight {
  color: var(--color-info);
}

.metric-value.overweight {
  color: var(--color-warning);
}

.metric-value.obese {
  color: var(--color-error);
}

.metric-desc {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
}

.text-content {
  color: var(--text-secondary);
  line-height: 1.6;
  white-space: pre-wrap;
}

.action-buttons {
  display: flex;
  gap: var(--spacing-md);
  justify-content: center;
  margin-top: var(--spacing-2xl);
  padding-top: var(--spacing-2xl);
  border-top: 1px solid var(--border-color);
}

.nutrition-content {
  padding: var(--spacing-lg) 0;
}

.target-calorie {
  text-align: center;
  padding: var(--spacing-xl);
  background: linear-gradient(135deg, var(--color-primary) 0%, #34d399 100%);
  border-radius: var(--radius-xl);
  margin-bottom: var(--spacing-xl);
  color: white;
  box-shadow: var(--shadow-lg);
}

.target-label {
  font-size: var(--font-size-sm);
  opacity: 0.9;
  margin-bottom: var(--spacing-xs);
}

.target-value {
  font-size: 42px;
  font-weight: var(--font-weight-bold);
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.target-value .unit {
  font-size: var(--font-size-lg);
  margin-left: var(--spacing-xs);
}

.macros-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--spacing-md);
}

.macro-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: var(--spacing-xs);
  padding: var(--spacing-lg) var(--spacing-md);
  background: #ffffff;
  border-radius: var(--radius-xl);
  border: 2px solid transparent;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-sm);
}

.macro-item:hover {
  transform: translateY(-4px);
  box-shadow: var(--shadow-md);
}

.macro-item.protein {
  border-color: var(--color-protein);
}

.macro-item.carb {
  border-color: var(--color-carb);
}

.macro-item.fat {
  border-color: var(--color-fat);
}

.macro-icon {
  font-size: 40px;
  margin-bottom: var(--spacing-xs);
}

.macro-info {
  width: 100%;
}

.macro-name {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-xs);
  font-weight: var(--font-weight-medium);
}

.macro-value {
  font-size: var(--font-size-2xl);
  font-weight: var(--font-weight-bold);
  color: var(--text-primary);
  margin-bottom: var(--spacing-xs);
}

.macro-ratio {
  font-size: var(--font-size-xs);
  color: var(--text-tertiary);
  background: var(--bg-secondary);
  padding: 2px var(--spacing-xs);
  border-radius: var(--radius-full);
  display: inline-block;
}

/* 响应式布局 */
@media (max-width: 1024px) {
  .middle-section,
  .bottom-section {
    grid-template-columns: 1fr;
  }
  
  .health-metrics {
    grid-template-columns: repeat(3, 1fr);
  }
  
  .macros-grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 768px) {
  .page-header {
    padding: var(--spacing-md);
  }

  .main-content {
    padding: 0 var(--spacing-md) var(--spacing-2xl);
  }

  .health-metrics {
    grid-template-columns: 1fr;
  }
  
  .macros-grid {
    grid-template-columns: 1fr;
  }
}
</style>

