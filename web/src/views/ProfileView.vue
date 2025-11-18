<template>
  <div class="profile-view-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button text @click="goBack">
            <el-icon><ArrowLeft /></el-icon>
            返回
          </el-button>
          <h1 class="page-title">个人档案</h1>
          <el-button type="primary" @click="goToEdit">
            <el-icon><Edit /></el-icon>
            更新档案
          </el-button>
        </div>
      </el-header>

      <el-main class="main-content">
        <div v-if="loading" class="loading-container">
          <el-skeleton :rows="10" animated />
        </div>

        <div v-else-if="!authStore.hasProfile" class="empty-state">
          <el-empty description="您还没有完善健康档案">
            <el-button type="primary" @click="goToEdit">立即完善档案</el-button>
          </el-empty>
        </div>

        <div v-else class="profile-content">
          <!-- 顶部：营养需求（最重要，突出显示） -->
          <div class="top-section" v-if="nutritionData">
            <el-card shadow="hover" class="nutrition-card highlight-card">
              <template #header>
                <div class="card-header highlight-header">
                  <el-icon class="card-icon"><Food /></el-icon>
                  <span>每日营养需求</span>
                </div>
              </template>
              <div class="nutrition-content">
                <div class="target-calorie">
                  <div class="target-label">目标热量</div>
                  <div class="target-value">{{ nutritionData.target_calorie }} kcal</div>
                </div>
                <div class="macros-grid">
                  <div class="macro-item protein">
                    <div class="macro-icon">🥩</div>
                    <div class="macro-info">
                      <div class="macro-name">蛋白质</div>
                      <div class="macro-value">{{ nutritionData.protein_gram }}g</div>
                      <div class="macro-ratio">{{ nutritionData.protein_ratio.toFixed(1) }}%</div>
                    </div>
                  </div>
                  <div class="macro-item carb">
                    <div class="macro-icon">🍚</div>
                    <div class="macro-info">
                      <div class="macro-name">碳水化合物</div>
                      <div class="macro-value">{{ nutritionData.carb_gram }}g</div>
                      <div class="macro-ratio">{{ nutritionData.carb_ratio.toFixed(1) }}%</div>
                    </div>
                  </div>
                  <div class="macro-item fat">
                    <div class="macro-icon">🥑</div>
                    <div class="macro-info">
                      <div class="macro-name">脂肪</div>
                      <div class="macro-value">{{ nutritionData.fat_gram }}g</div>
                      <div class="macro-ratio">{{ nutritionData.fat_ratio.toFixed(1) }}%</div>
                    </div>
                  </div>
                </div>
              </div>
            </el-card>
          </div>

          <!-- 中间：两列布局 - 基础信息和健康指标 -->
          <div class="middle-section">
            <el-card shadow="hover" class="info-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><User /></el-icon>
                  <span>基础信息</span>
                </div>
              </template>
              <el-descriptions :column="2" border>
                <el-descriptions-item label="用户名">
                  <el-tag>{{ authStore.user?.username }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="邮箱">
                  {{ authStore.profile?.email || '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="性别">
                  <el-tag :type="authStore.profile?.gender === '男' ? 'primary' : 'danger'">
                    {{ authStore.profile?.gender }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="年龄">
                  {{ authStore.profile?.age }} 岁
                </el-descriptions-item>
                <el-descriptions-item label="身高">
                  {{ authStore.profile?.height }} cm
                </el-descriptions-item>
                <el-descriptions-item label="体重">
                  {{ authStore.profile?.weight }} kg
                </el-descriptions-item>
              </el-descriptions>
            </el-card>

            <el-card shadow="hover" class="info-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><DataAnalysis /></el-icon>
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
            </el-card>
          </div>

          <!-- 底部：两列布局 - 健康目标和个性化设置 -->
          <div class="bottom-section">
            <el-card shadow="hover" class="info-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Aim /></el-icon>
                  <span>健康目标</span>
                </div>
              </template>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="健康目标">
                  <el-tag :type="getGoalTagType(authStore.profile?.health_goal)" size="large">
                    {{ authStore.profile?.health_goal || '未设置' }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="目标体重">
                  {{ authStore.profile?.target_weight ? `${authStore.profile.target_weight} kg` : '未设置' }}
                </el-descriptions-item>
                <el-descriptions-item label="活动水平">
                  <el-tag type="info">{{ authStore.profile?.activity_level || '未设置' }}</el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="每日用餐次数">
                  {{ authStore.profile?.meal_times_per_day || 3 }} 次
                </el-descriptions-item>
              </el-descriptions>
            </el-card>

            <el-card shadow="hover" class="info-card">
              <template #header>
                <div class="card-header">
                  <el-icon class="card-icon"><Setting /></el-icon>
                  <span>个性化设置</span>
                </div>
              </template>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="过敏源">
                  <div class="text-content">
                    {{ authStore.profile?.allergies || '无' }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="饮食偏好">
                  <div class="text-content">
                    {{ authStore.profile?.dietary_prefs || '无' }}
                  </div>
                </el-descriptions-item>
                <el-descriptions-item label="健康问题">
                  <div class="text-content">
                    {{ authStore.profile?.health_conditions || '无' }}
                  </div>
                </el-descriptions-item>
              </el-descriptions>
            </el-card>
          </div>

          <!-- 操作按钮 -->
          <div class="action-buttons">
            <el-button type="primary" size="large" @click="goToEdit">
              <el-icon><Edit /></el-icon>
              更新档案
            </el-button>
            <el-button size="large" @click="goBack">
              返回首页
            </el-button>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowLeft,
  Edit,
  User,
  DataAnalysis,
  Aim,
  Setting,
  Food
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { getNutritionRequirements } from '@/api/user'

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
    '减脂': 'danger',
    '增肌': 'success',
    '控糖': 'warning',
    '维持健康': 'info'
  }
  return goalMap[goal] || ''
}
</script>

<style scoped>
.profile-view-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
}

.header {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 30px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.main-content {
  padding: 30px;
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
  padding: 40px;
}

.empty-state {
  padding: 60px 20px;
}

.profile-content {
  animation: fadeIn 0.5s ease-out;
  display: flex;
  flex-direction: column;
  gap: 20px;
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
  border-radius: 16px;
  overflow: hidden;
  border: 2px solid #667eea;
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.2);
}

.highlight-card {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9ff 100%);
}

.highlight-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 12px 20px;
  margin: -20px -20px 20px -20px;
}

.highlight-header .card-icon {
  color: white;
}

.highlight-header span {
  color: white;
}

/* 中间区域：两列布局 */
.middle-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

/* 底部区域：两列布局 */
.bottom-section {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-card {
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
}

.card-icon {
  font-size: 20px;
  color: #667eea;
}

.health-metrics {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  padding: 20px 0;
}

.metric-item {
  text-align: center;
  padding: 16px;
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  border-radius: 10px;
  border: 1px solid #e4e7ed;
  transition: all 0.3s ease;
}

.metric-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.metric-label {
  font-size: 14px;
  color: #666;
  margin-bottom: 10px;
  font-weight: 500;
}

.metric-value {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 6px;
  color: #333;
}

.metric-value.primary {
  color: #667eea;
}

.metric-value.normal {
  color: #67c23a;
}

.metric-value.underweight {
  color: #409eff;
}

.metric-value.overweight {
  color: #e6a23c;
}

.metric-value.obese {
  color: #f56c6c;
}

.metric-desc {
  font-size: 12px;
  color: #999;
}

.text-content {
  color: #666;
  line-height: 1.6;
  white-space: pre-wrap;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: center;
  margin-top: 30px;
  padding-top: 30px;
  border-top: 1px solid #e4e7ed;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
  width: 120px;
}

:deep(.el-descriptions__content) {
  color: #333;
}

.nutrition-content {
  padding: 20px 0;
}

.target-calorie {
  text-align: center;
  padding: 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 12px;
  margin-bottom: 24px;
  color: white;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.3);
}

.target-label {
  font-size: 14px;
  opacity: 0.9;
  margin-bottom: 8px;
}

.target-value {
  font-size: 42px;
  font-weight: 700;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.macros-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.macro-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 8px;
  padding: 20px 16px;
  background: #ffffff;
  border-radius: 12px;
  border: 2px solid transparent;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.macro-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.macro-item.protein {
  border-color: #409eff;
}

.macro-item.carb {
  border-color: #67c23a;
}

.macro-item.fat {
  border-color: #e6a23c;
}

.macro-icon {
  font-size: 40px;
  margin-bottom: 4px;
}

.macro-info {
  width: 100%;
}

.macro-name {
  font-size: 13px;
  color: #666;
  margin-bottom: 6px;
  font-weight: 500;
}

.macro-value {
  font-size: 24px;
  font-weight: 700;
  color: #333;
  margin-bottom: 4px;
}

.macro-ratio {
  font-size: 12px;
  color: #999;
  background: #f0f0f0;
  padding: 2px 8px;
  border-radius: 10px;
  display: inline-block;
}
</style>

