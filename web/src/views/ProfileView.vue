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
          <!-- 基础信息卡片 -->
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

          <!-- 健康指标卡片 -->
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
                <div class="metric-desc">基础代谢率 (kcal/天)</div>
              </div>
              <div class="metric-item">
                <div class="metric-label">TDEE</div>
                <div class="metric-value primary">{{ authStore.profile?.tdee || '--' }}</div>
                <div class="metric-desc">每日总消耗 (kcal/天)</div>
              </div>
            </div>
          </el-card>

          <!-- 健康目标卡片 -->
          <el-card shadow="hover" class="info-card">
            <template #header>
              <div class="card-header">
                <el-icon class="card-icon"><Aim /></el-icon>
                <span>健康目标</span>
              </div>
            </template>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="健康目标">
                <el-tag :type="getGoalTagType(authStore.profile?.health_goal)">
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

          <!-- 个性化设置卡片 -->
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
  Setting
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

onMounted(async () => {
  loading.value = true
  try {
    if (!authStore.profile) {
      await authStore.loadProfile()
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
  max-width: 1200px;
  margin: 0 auto;
}

.loading-container {
  padding: 40px;
}

.empty-state {
  padding: 60px 20px;
}

.profile-content {
  animation: fadeIn 0.5s ease-out;
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

.info-card {
  margin-bottom: 20px;
  border-radius: 12px;
  overflow: hidden;
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  padding: 20px 0;
}

.metric-item {
  text-align: center;
  padding: 20px;
  background: linear-gradient(135deg, #f5f7fa 0%, #ffffff 100%);
  border-radius: 12px;
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
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 8px;
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
</style>

