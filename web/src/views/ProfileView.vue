<template>
  <div class="profile-view-page">
    <div class="profile-shell">
      <header class="page-header">
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
          <div class="hero-card">
            <div class="hero-main">
              <div class="avatar-chip">{{ getUserInitial() }}</div>
              <div>
                <p class="hero-greet">你好，{{ authStore.user?.username || 'Nutri 用户' }}</p>
                <h2 class="hero-goal">{{ authStore.profile?.health_goal || '维持健康' }} · 今日继续加油</h2>
                <p class="hero-note">
                  目标体重：{{ authStore.profile?.target_weight ?? authStore.profile?.targetWeight ?? '--' }} kg ·
                  每日用餐：{{ authStore.profile?.meal_times_per_day || authStore.profile?.mealTimesPerDay || 3 }} 次
                </p>
              </div>
            </div>

            <div class="hero-metrics">
              <div class="metric-pill">
                <span class="metric-pill__label">BMI</span>
                <strong class="metric-pill__value" :class="getBMIClass(authStore.profile?.bmi)">
                  {{ authStore.profile?.bmi || '--' }}
                </strong>
                <span class="metric-pill__desc">{{ getBMIDesc(authStore.profile?.bmi) }}</span>
              </div>
              <div class="metric-pill">
                <span class="metric-pill__label">BMR</span>
                <strong class="metric-pill__value">{{ authStore.profile?.bmr || '--' }}</strong>
                <span class="metric-pill__desc">基础代谢</span>
              </div>
              <div class="metric-pill">
                <span class="metric-pill__label">TDEE</span>
                <strong class="metric-pill__value">{{ authStore.profile?.tdee || '--' }}</strong>
                <span class="metric-pill__desc">维持热量</span>
              </div>
            </div>
          </div>

          <n-card v-if="nutritionData" class="nutrition-card" :bordered="false">
            <div class="section-title-row">
              <h3>每日营养需求</h3>
                <n-tag round size="small" type="success">
                  {{ authStore.profile?.health_goal || authStore.profile?.healthGoal || '个性化推荐' }}
              </n-tag>
            </div>

            <div class="nutrition-summary">
              <div class="calorie-block">
                <div class="calorie-label">目标热量</div>
                <div class="calorie-value">{{ nutritionData.target_calorie }}<span>kcal</span></div>
              </div>

              <div class="macro-bars">
                <div class="macro-row protein">
                  <div class="macro-row__head">
                    <span>蛋白质</span>
                    <span>{{ nutritionData.protein_gram }}g · {{ nutritionData.protein_ratio.toFixed(0) }}%</span>
                  </div>
                  <div class="macro-row__track">
                    <div class="macro-row__fill" :style="{ width: `${nutritionData.protein_ratio}%` }"></div>
                  </div>
                </div>

                <div class="macro-row carb">
                  <div class="macro-row__head">
                    <span>碳水化合物</span>
                    <span>{{ nutritionData.carb_gram }}g · {{ nutritionData.carb_ratio.toFixed(0) }}%</span>
                  </div>
                  <div class="macro-row__track">
                    <div class="macro-row__fill" :style="{ width: `${nutritionData.carb_ratio}%` }"></div>
                  </div>
                </div>

                <div class="macro-row fat">
                  <div class="macro-row__head">
                    <span>脂肪</span>
                    <span>{{ nutritionData.fat_gram }}g · {{ nutritionData.fat_ratio.toFixed(0) }}%</span>
                  </div>
                  <div class="macro-row__track">
                    <div class="macro-row__fill" :style="{ width: `${nutritionData.fat_ratio}%` }"></div>
                  </div>
                </div>
              </div>
            </div>
          </n-card>

          <div class="content-grid">
            <n-card class="info-card" :bordered="false">
              <template #header>
                <div class="card-header"><span>基础信息</span></div>
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
                    {{ authStore.profile?.gender || '未设置' }}
                  </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="年龄">{{ authStore.profile?.age || '--' }} 岁</n-descriptions-item>
                <n-descriptions-item label="身高">{{ authStore.profile?.height || '--' }} cm</n-descriptions-item>
                <n-descriptions-item label="体重">{{ authStore.profile?.weight || '--' }} kg</n-descriptions-item>
              </n-descriptions>
            </n-card>

            <n-card class="info-card" :bordered="false">
              <template #header>
                <div class="card-header"><span>目标与活动</span></div>
              </template>
              <n-descriptions :column="1" bordered label-placement="left">
                <n-descriptions-item label="健康目标">
                  <n-tag :type="getGoalTagType(authStore.profile?.health_goal || authStore.profile?.healthGoal)" round>
                    {{ authStore.profile?.health_goal || authStore.profile?.healthGoal || '未设置' }}
                  </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="目标体重">
                  {{ (authStore.profile?.target_weight ?? authStore.profile?.targetWeight) ? `${authStore.profile?.target_weight ?? authStore.profile?.targetWeight} kg` : '未设置' }}
                </n-descriptions-item>
                <n-descriptions-item label="活动水平">
                  <n-tag type="info" round>{{ authStore.profile?.activity_level || authStore.profile?.activityLevel || '未设置' }}</n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="每日用餐次数">
                  {{ authStore.profile?.meal_times_per_day || authStore.profile?.mealTimesPerDay || 3 }} 次
                </n-descriptions-item>
              </n-descriptions>
            </n-card>

            <n-card class="info-card" :bordered="false">
              <template #header>
                <div class="card-header"><span>个性化偏好</span></div>
              </template>
              <n-descriptions :column="1" bordered label-placement="left">
                <n-descriptions-item label="过敏源">
                  <div class="text-content">{{ authStore.profile?.allergies || '无' }}</div>
                </n-descriptions-item>
                <n-descriptions-item label="饮食偏好">
                  <div class="text-content">{{ authStore.profile?.dietary_prefs || authStore.profile?.dietaryPrefs || '无' }}</div>
                </n-descriptions-item>
                <n-descriptions-item label="健康问题">
                  <div class="text-content">{{ authStore.profile?.health_conditions || authStore.profile?.healthConditions || '无' }}</div>
                </n-descriptions-item>
              </n-descriptions>
            </n-card>

            <n-card class="info-card" :bordered="false">
              <template #header>
                <div class="card-header"><span>身体指标解读</span></div>
              </template>
              <div class="insight-grid">
                <div class="insight-item">
                  <div class="insight-item__label">BMI 状态</div>
                  <div class="insight-item__value">{{ getBMIDesc(authStore.profile?.bmi) }}</div>
                  <div class="insight-item__desc">建议维持规律作息与稳定运动节奏。</div>
                </div>
                <div class="insight-item">
                  <div class="insight-item__label">代谢能力</div>
                  <div class="insight-item__value">{{ authStore.profile?.bmr || '--' }} kcal</div>
                  <div class="insight-item__desc">基础代谢是每日饮食规划的参考基线。</div>
                </div>
                <div class="insight-item">
                  <div class="insight-item__label">维持热量</div>
                  <div class="insight-item__value">{{ authStore.profile?.tdee || '--' }} kcal</div>
                  <div class="insight-item__desc">结合活动量，动态调整更容易长期坚持。</div>
                </div>
              </div>
            </n-card>
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
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { NCard, NButton, NIcon, NSkeleton, NEmpty, NDescriptions, NDescriptionsItem, NTag } from 'naive-ui'
import { CreateOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/store/auth'
import { getNutritionRequirements } from '@/api/user'
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

const getGoalTagType = (goal) => {
  const goalMap = {
    减脂: 'error',
    增肌: 'success',
    控糖: 'warning',
    维持健康: 'info'
  }
  return goalMap[goal] || 'default'
}
</script>

<style scoped>
.profile-view-page {
  min-height: 100vh;
  background: radial-gradient(circle at 0% 0%, #edf9ef 0%, #f4f8ff 36%, #eef3fa 100%);
}

.profile-shell {
  max-width: 1240px;
  margin: 0 auto;
  padding: 20px 20px 44px;
}

.page-header {
  display: grid;
  grid-template-columns: auto 1fr auto;
  align-items: center;
  gap: 14px;
  margin-bottom: 18px;
}

.page-title-wrap {
  display: grid;
  justify-items: center;
  gap: 2px;
}

.page-title {
  margin: 0;
  font-size: clamp(24px, 3.3vw, 32px);
  font-weight: 700;
  letter-spacing: -0.01em;
  color: #1f2937;
}

.page-subtitle {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
}

.update-btn {
  border-radius: 10px;
}

.page-content {
  display: grid;
  gap: 18px;
}

.loading-section,
.empty-section {
  padding: 30px 8px;
}

.profile-content {
  display: grid;
  gap: 18px;
  animation: fadeInUp 0.35s ease;
}

.hero-card {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 14px;
  padding: 18px;
  border-radius: 20px;
  background: linear-gradient(145deg, #f8fff9 0%, #f5faff 100%);
  border: 1px solid rgba(130, 184, 150, 0.26);
  box-shadow: 0 12px 30px rgba(21, 42, 36, 0.07);
}

.hero-main {
  display: flex;
  align-items: center;
  gap: 14px;
}

.avatar-chip {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: grid;
  place-items: center;
  color: #fff;
  font-size: 24px;
  font-weight: 700;
  background: linear-gradient(135deg, #39a36b 0%, #2a8969 100%);
  box-shadow: 0 10px 20px rgba(43, 137, 105, 0.28);
}

.hero-greet {
  margin: 0;
  color: #4b5563;
  font-size: 13px;
}

.hero-goal {
  margin: 4px 0;
  font-size: clamp(18px, 2.4vw, 24px);
  color: #0f172a;
  letter-spacing: -0.01em;
}

.hero-note {
  margin: 0;
  font-size: 13px;
  color: #6b7280;
}

.hero-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 10px;
}

.metric-pill {
  padding: 12px 10px;
  border-radius: 14px;
  border: 1px solid #dce6df;
  background: rgba(255, 255, 255, 0.72);
  text-align: center;
}

.metric-pill__label {
  display: block;
  font-size: 11px;
  color: #6b7280;
  margin-bottom: 4px;
}

.metric-pill__value {
  display: block;
  font-size: 20px;
  font-weight: 700;
  color: #111827;
  line-height: 1.1;
}

.metric-pill__value.normal {
  color: #10b981;
}

.metric-pill__value.underweight {
  color: #3b82f6;
}

.metric-pill__value.overweight {
  color: #f59e0b;
}

.metric-pill__value.obese {
  color: #ef4444;
}

.metric-pill__desc {
  display: block;
  margin-top: 4px;
  font-size: 11px;
  color: #6b7280;
}

.nutrition-card {
  border-radius: 18px;
  background: linear-gradient(135deg, #ffffff 0%, #f4fff8 100%);
  box-shadow: 0 10px 24px rgba(18, 34, 24, 0.08);
}

.section-title-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 14px;
}

.section-title-row h3 {
  margin: 0;
  font-size: 19px;
  color: #111827;
}

.nutrition-summary {
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
}

.calorie-block {
  border-radius: 14px;
  padding: 16px;
  background: linear-gradient(160deg, #2faa67 0%, #43b27a 100%);
  color: #fff;
  box-shadow: 0 12px 24px rgba(48, 170, 110, 0.28);
}

.calorie-label {
  font-size: 13px;
  opacity: 0.9;
}

.calorie-value {
  margin-top: 8px;
  font-size: 38px;
  font-weight: 700;
  line-height: 1;
}

.calorie-value span {
  margin-left: 6px;
  font-size: 14px;
  font-weight: 500;
}

.macro-bars {
  display: grid;
  gap: 12px;
}

.macro-row {
  border-radius: 12px;
  border: 1px solid #e6ecf2;
  background: #fff;
  padding: 12px;
}

.macro-row__head {
  display: flex;
  justify-content: space-between;
  gap: 8px;
  font-size: 13px;
  color: #4b5563;
  margin-bottom: 8px;
}

.macro-row__track {
  height: 8px;
  border-radius: 99px;
  overflow: hidden;
  background: #edf2f7;
}

.macro-row__fill {
  height: 100%;
  border-radius: inherit;
}

.macro-row.protein .macro-row__fill {
  background: linear-gradient(90deg, #f87171 0%, #ef4444 100%);
}

.macro-row.carb .macro-row__fill {
  background: linear-gradient(90deg, #60a5fa 0%, #3b82f6 100%);
}

.macro-row.fat .macro-row__fill {
  background: linear-gradient(90deg, #fbbf24 0%, #f59e0b 100%);
}

.content-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.info-card {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.82);
  box-shadow: 0 10px 22px rgba(15, 23, 42, 0.06);
  backdrop-filter: blur(6px);
}

.card-header {
  font-size: 18px;
  font-weight: 700;
  color: #1f2937;
}

.text-content {
  color: #4b5563;
  line-height: 1.65;
  white-space: pre-wrap;
}

.insight-grid {
  display: grid;
  gap: 10px;
}

.insight-item {
  border: 1px solid #e5ebf0;
  border-radius: 12px;
  background: #fbfdff;
  padding: 12px;
}

.insight-item__label {
  font-size: 12px;
  color: #64748b;
}

.insight-item__value {
  margin-top: 2px;
  font-size: 20px;
  font-weight: 700;
  color: #0f172a;
}

.insight-item__desc {
  margin-top: 4px;
  font-size: 12px;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 2px;
}

@media (max-width: 1024px) {
  .hero-card {
    grid-template-columns: 1fr;
  }

  .nutrition-summary {
    grid-template-columns: 1fr;
  }

  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .profile-shell {
    padding: 12px 12px 28px;
  }

  .page-header {
    grid-template-columns: 1fr;
    justify-items: start;
    gap: 10px;
  }

  .page-title-wrap {
    justify-items: start;
  }

  .hero-main {
    align-items: flex-start;
  }

  .hero-metrics {
    grid-template-columns: 1fr;
  }

  .calorie-value {
    font-size: 32px;
  }

  .action-buttons {
    justify-content: stretch;
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

