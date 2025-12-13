<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1 class="logo">🍽️ NutriPlan</h1>
          <div class="header-right">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" :icon="UserFilled" />
                <span class="username">{{ authStore.user?.username || '用户' }}</span>
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">
                    <el-icon><User /></el-icon>
                    个人档案
                  </el-dropdown-item>
                  <el-dropdown-item command="favorites">
                    <el-icon><Star /></el-icon>
                    我的收藏
                  </el-dropdown-item>
                  <el-dropdown-item command="intake">
                    <el-icon><DataLine /></el-icon>
                    饮食记录
                  </el-dropdown-item>
                  <el-dropdown-item command="weekly">
                    <el-icon><TrendCharts /></el-icon>
                    周报告
                  </el-dropdown-item>
                  <el-dropdown-item divided command="logout">
                    <el-icon><SwitchButton /></el-icon>
                    退出登录
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </el-header>

      <el-main class="main-content">
        <!-- Welcome Section -->
        <div class="welcome-section" v-if="authStore.user">
          <h2 class="welcome-title">欢迎回来，{{ authStore.user.username }}！</h2>
          <p class="welcome-subtitle">开启您的健康营养之旅</p>
        </div>

        <!-- Feature Cards Grid -->
        <div class="features-grid">
          <!-- My Profile Card -->
          <div class="feature-card profile-card" @click="goToProfile">
            <div class="card-icon-wrapper profile">
              <el-icon class="card-icon"><User /></el-icon>
            </div>
            <h3 class="card-title">我的档案</h3>
            <p class="card-description">
              {{ authStore.hasProfile ? '管理健康数据与营养目标' : '完善档案以获取个性化推荐' }}
            </p>
            <div class="nutrition-preview" v-if="authStore.profile">
              <div class="preview-item">
                <span class="preview-label">每日热量目标</span>
                <span class="preview-value">{{ nutritionData?.target_calorie || (authStore.profile.tdee ? Math.floor(authStore.profile.tdee) : '--') }} kcal</span>
              </div>
            </div>
            <div class="card-status" v-if="!authStore.hasProfile">
              <el-tag type="warning" effect="plain">待完善</el-tag>
            </div>
            <div class="card-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>

          <!-- Recipe Recommendation Card -->
          <div 
            class="feature-card recipe-card" 
            @click="goToRecipes"
            :class="{ disabled: !authStore.hasProfile }"
          >
            <div class="card-icon-wrapper recipe">
              <el-icon class="card-icon"><KnifeFork /></el-icon>
            </div>
            <h3 class="card-title">智能食谱推荐</h3>
            <p class="card-description">
              {{ authStore.hasProfile ? '获取个性化的每日食谱方案' : '请先完善健康档案' }}
            </p>
            <div class="card-badge" v-if="authStore.hasProfile">
              <span class="badge-text">AI 推荐</span>
            </div>
            <div class="card-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>

          <!-- Intake Record Card -->
          <div class="feature-card intake-card" @click="goToIntake">
            <div class="card-icon-wrapper intake">
              <el-icon class="card-icon"><DataLine /></el-icon>
            </div>
            <h3 class="card-title">饮食记录</h3>
            <p class="card-description">记录每日饮食摄入与营养达标率</p>
            <div class="card-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>

          <!-- Weekly Report Card -->
          <div class="feature-card weekly-card" @click="goToWeekly">
            <div class="card-icon-wrapper weekly">
              <el-icon class="card-icon"><TrendCharts /></el-icon>
            </div>
            <h3 class="card-title">周报告</h3>
            <p class="card-description">查看每周饮食趋势分析</p>
            <div class="card-arrow">
              <el-icon><ArrowRight /></el-icon>
            </div>
          </div>

          <!-- Coming Soon Cards -->

          <div class="feature-card coming-soon-card disabled">
            <div class="card-icon-wrapper coming">
              <el-icon class="card-icon"><Medal /></el-icon>
            </div>
            <h3 class="card-title">健康目标</h3>
            <p class="card-description">设定和追踪您的目标</p>
            <el-tag type="info" size="small">即将推出</el-tag>
          </div>

          <div class="feature-card coming-soon-card disabled">
            <div class="card-icon-wrapper coming">
              <el-icon class="card-icon"><PieChart /></el-icon>
            </div>
            <h3 class="card-title">数据分析</h3>
            <p class="card-description">可视化您的营养数据</p>
            <el-tag type="info" size="small">即将推出</el-tag>
          </div>
        </div>

        <!-- Quick Tips -->
        <div class="tips-section" v-if="!authStore.hasProfile">
          <el-alert
            title="温馨提示"
            type="info"
            :closable="false"
            show-icon
          >
            <p>请先完善您的健康档案，以便我们为您提供个性化的营养推荐服务。</p>
            <el-button type="primary" size="small" @click="goToProfile" style="margin-top: 10px;">
              立即完善
            </el-button>
          </el-alert>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UserFilled,
  User,
  ArrowDown,
  ArrowRight,
  SwitchButton,
  Document,
  KnifeFork,
  TrendCharts,
  DataLine,
  Medal,
  PieChart,
  Star
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'
import { getNutritionRequirements } from '@/api/user'

const router = useRouter()
const authStore = useAuthStore()
const nutritionData = ref(null)

onMounted(async () => {
  if (authStore.isAuthenticated) {
    if (!authStore.user || !authStore.profile) {
      try {
        await authStore.loadProfile()
      } catch (error) {
        console.error('加载档案失败:', error)
      }
    }
    
    // 加载营养需求数据
    if (authStore.hasProfile) {
      try {
        const response = await getNutritionRequirements()
        nutritionData.value = response
      } catch (error) {
        console.log('营养需求数据不可用:', error)
      }
    }
  }
})

const goToProfile = () => {
  router.push('/profile/view')
}

const goToRecipes = () => {
  if (!authStore.hasProfile) {
    ElMessage.warning('请先完善健康档案')
    return
  }
  router.push('/recipes')
}

const goToIntake = () => {
  router.push('/intake')
}

const goToWeekly = () => {
  router.push('/weekly')
}

const goToNutrition = () => {
  if (!authStore.hasProfile) {
    ElMessage.warning('请先完善健康档案')
    return
  }
  // 可以创建一个专门的营养需求展示页面
  router.push('/profile/view')
}

const goToFavorites = () => {
  router.push('/favorites')
}

const handleCommand = async (command) => {
  if (command === 'profile') {
    router.push('/profile/view')
  } else if (command === 'favorites') {
    goToFavorites()
  } else if (command === 'intake') {
    goToIntake()
  } else if (command === 'weekly') {
    goToWeekly()
  } else if (command === 'logout') {
    try {
      await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
      authStore.logout()
      ElMessage.success('已退出登录')
      router.push('/login')
    } catch {
      // 用户取消
    }
  }
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.header {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 20px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
  padding: 0 30px;
}

.logo {
  font-size: 26px;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 8px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-weight: 500;
  color: #333;
}

.main-content {
  padding: 40px 30px;
  max-width: 1400px;
  margin: 0 auto;
}

.welcome-section {
  text-align: center;
  color: white;
  margin-bottom: 50px;
  animation: fadeIn 0.6s ease-out;
}

.welcome-title {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 10px 0;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.2);
}

.welcome-subtitle {
  font-size: 1.2rem;
  opacity: 0.9;
  margin: 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  grid-auto-rows: minmax(180px, auto);
  gap: 24px;
  animation: fadeInUp 0.8s ease-out;
}

.feature-card {
  background: white;
  border-radius: 24px;
  padding: 30px;
  cursor: pointer;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

/* Bento Grid Layout */
.recipe-card {
  grid-column: span 2;
  grid-row: span 2;
  background: linear-gradient(135deg, #fff 0%, #fff5f5 100%);
  border: 1px solid rgba(245, 87, 108, 0.1);
}

.profile-card {
  grid-column: span 1;
  grid-row: span 1;
}

.nutrition-card {
  grid-column: span 1;
  grid-row: span 1;
}

.coming-soon-card {
  grid-column: span 1;
  grid-row: span 1;
}

/* Special styling for the large recipe card */
.recipe-card .card-icon-wrapper {
  width: 80px;
  height: 80px;
  margin-bottom: 30px;
}

.recipe-card .card-title {
  font-size: 2rem;
  margin-bottom: 15px;
}

.recipe-card .card-description {
  font-size: 1.1rem;
  max-width: 80%;
}

.feature-card:hover:not(.disabled) {
  transform: translateY(-8px) scale(1.02);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.12);
  z-index: 1;
}

.feature-card.disabled {
  cursor: not-allowed;
  opacity: 0.7;
  background: #f8f9fa;
}

.card-icon-wrapper {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  transition: transform 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.feature-card:hover:not(.disabled) .card-icon-wrapper {
  transform: scale(1.1) rotate(5deg);
}

.card-icon-wrapper.profile {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  box-shadow: 0 8px 16px rgba(118, 75, 162, 0.2);
}

.card-icon-wrapper.recipe {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  box-shadow: 0 8px 16px rgba(245, 87, 108, 0.3);
}

.card-icon-wrapper.nutrition {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  box-shadow: 0 8px 16px rgba(0, 242, 254, 0.2);
}

.card-icon-wrapper.intake {
  background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);
  box-shadow: 0 8px 16px rgba(168, 237, 234, 0.3);
}

.card-icon-wrapper.weekly {
  background: linear-gradient(135deg, #ffecd2 0%, #fcb69f 100%);
  box-shadow: 0 8px 16px rgba(252, 182, 159, 0.3);
}

.card-icon-wrapper.coming {
  background: linear-gradient(135deg, #e0c3fc 0%, #8ec5fc 100%);
}

.card-icon {
  font-size: 28px;
  color: white;
}

.recipe-card .card-icon {
  font-size: 40px;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 8px 0;
}

.card-description {
  color: #718096;
  font-size: 0.9rem;
  line-height: 1.6;
  margin: 0;
  flex-grow: 1;
}

.card-status {
  margin-top: 15px;
}

.card-badge {
  position: absolute;
  top: 20px;
  right: 20px;
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
  color: white;
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  box-shadow: 0 4px 10px rgba(245, 87, 108, 0.3);
}

.nutrition-preview {
  background: #f7fafc;
  padding: 12px;
  border-radius: 12px;
  margin-top: 15px;
}

.preview-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.preview-label {
  font-size: 0.85rem;
  color: #718096;
  font-weight: 500;
}

.preview-value {
  font-size: 1rem;
  color: #2d3748;
  font-weight: 700;
}

.card-arrow {
  position: absolute;
  right: 24px;
  bottom: 24px;
  font-size: 24px;
  color: #cbd5e0;
  transition: all 0.3s;
}

.feature-card:hover:not(.disabled) .card-arrow {
  color: #667eea;
  transform: translateX(5px);
}

.tips-section {
  margin-top: 40px;
  animation: fadeIn 1s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Responsive Breakpoints */
@media (max-width: 1200px) {
  .features-grid {
    grid-template-columns: repeat(3, 1fr);
  }
  .recipe-card {
    grid-column: span 2;
    grid-row: span 2;
  }
}

@media (max-width: 900px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  .recipe-card {
    grid-column: span 2;
    grid-row: auto; /* Reset row span on smaller screens */
  }
}

@media (max-width: 600px) {
  .features-grid {
    grid-template-columns: 1fr;
  }
  .recipe-card {
    grid-column: auto;
  }
}
</style>
