<template>
  <div class="home-container">
    <!-- 顶部导航 -->
    <div class="top-bar">
      <div class="top-bar-content">
        <div class="brand">
          <span class="brand-icon">🥗</span>
          <span class="brand-name">NutriPlan</span>
        </div>
        <el-dropdown @command="handleCommand" trigger="click">
          <div class="user-menu">
            <el-avatar :size="36" :icon="UserFilled" style="background: #10b981;" />
            <span class="user-name">{{ authStore.user?.username }}</span>
            <el-icon class="dropdown-icon"><ArrowDown /></el-icon>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">个人档案</el-dropdown-item>
              <el-dropdown-item command="favorites">我的收藏</el-dropdown-item>
              <el-dropdown-item command="intake">饮食记录</el-dropdown-item>
              <el-dropdown-item command="weekly">周报告</el-dropdown-item>
              <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="main-wrapper">
      <!-- 欢迎区域 -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">
            {{ getGreeting() }}，<span class="highlight">{{ authStore.user?.username }}</span> 👋
          </h1>
          <p class="hero-subtitle">今天想吃点什么呢？</p>
        </div>
        <div class="quick-stats" v-if="authStore.hasProfile && authStore.profile">
          <div class="stat-item">
            <div class="stat-icon">🔥</div>
            <div class="stat-info">
              <div class="stat-value">{{ authStore.profile.tdee ? Math.floor(authStore.profile.tdee) : '--' }}</div>
              <div class="stat-label">每日目标热量</div>
            </div>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <div class="stat-icon">⚖️</div>
            <div class="stat-info">
              <div class="stat-value">{{ authStore.profile.current_weight || '--' }}</div>
              <div class="stat-label">当前体重 (kg)</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 功能卡片区 -->
      <div class="content-section">
        <div class="features-grid">
          <!-- 主要功能卡片 -->
          <div class="main-card big-card" @click="goToRecipes" :class="{ disabled: !authStore.hasProfile }">
            <div class="card-content">
              <div class="card-header">
                <span class="card-emoji">🍳</span>
                <h2 class="card-title-big">今日食谱</h2>
              </div>
              <p class="card-desc">{{ authStore.hasProfile ? '获取推荐的营养均衡食谱' : '完善档案后开启智能推荐' }}</p>
            </div>
          </div>

          <!-- 个人档案 -->
          <div class="main-card" @click="goToProfile">
            <div class="card-content">
              <div class="card-header-small">
                <span class="card-emoji-small">👤</span>
                <h3 class="card-title-small">个人档案</h3>
              </div>
              <p class="card-desc-small">{{ authStore.hasProfile ? '查看健康数据' : '完善基本信息' }}</p>
              <div class="card-status" v-if="!authStore.hasProfile">
                <el-tag type="warning" size="small" effect="plain">待完善</el-tag>
              </div>
            </div>
          </div>

          <!-- 饮食记录 -->
          <div class="main-card" @click="goToIntake">
            <div class="card-content">
              <div class="card-header-small">
                <span class="card-emoji-small">📝</span>
                <h3 class="card-title-small">饮食记录</h3>
              </div>
              <p class="card-desc-small">记录今天吃了什么</p>
            </div>
          </div>

          <!-- 周报告 -->
          <div class="main-card" @click="goToWeekly">
            <div class="card-content">
              <div class="card-header-small">
                <span class="card-emoji-small">📊</span>
                <h3 class="card-title-small">周报告</h3>
              </div>
              <p class="card-desc-small">查看本周营养趋势</p>
            </div>
          </div>

          <!-- 收藏夹 -->
          <div class="main-card" @click="goToFavorites">
            <div class="card-content">
              <div class="card-header-small">
                <span class="card-emoji-small">⭐</span>
                <h3 class="card-title-small">我的收藏</h3>
              </div>
              <p class="card-desc-small">喜欢的食谱都在这</p>
            </div>
          </div>

          <!-- 购物清单 -->
          <div class="main-card" @click="goToShopping">
            <div class="card-content">
              <div class="card-header-small">
                <span class="card-emoji-small">🛒</span>
                <h3 class="card-title-small">购物清单</h3>
              </div>
              <p class="card-desc-small">食材采购一键搞定</p>
            </div>
          </div>
        </div>

        <!-- 提示信息 -->
        <div class="tip-card" v-if="!authStore.hasProfile">
          <div class="tip-icon">💡</div>
          <div class="tip-content">
            <h4 class="tip-title">开始您的健康之旅</h4>
            <p class="tip-text">完善个人档案，获取专属营养方案</p>
            <el-button type="primary" size="small" @click="goToProfile" round>
              立即完善
            </el-button>
          </div>
        </div>
      </div>
    </div>
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

const goToShopping = () => {
  console.log('点击购物清单卡片')
  router.push('/shopping')
}

// 获取问候语
const getGreeting = () => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 9) return '早上好'
  if (hour < 12) return '上午好'
  if (hour < 14) return '中午好'
  if (hour < 18) return '下午好'
  if (hour < 22) return '晚上好'
  return '夜深了'
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
  background: #fafbfc;
  position: relative;
}

/* 顶部导航栏 */
.top-bar {
  background: #ffffff;
  border-bottom: 1px solid #e8eaed;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 2px rgba(0,0,0,0.05);
}

.top-bar-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 16px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.brand {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.brand-icon {
  font-size: 24px;
}

.brand-name {
  font-size: 20px;
  color: #1a1a1a;
  letter-spacing: -0.5px;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s;
}

.user-menu:hover {
  background: #f0f2f5;
}

.user-name {
  font-size: 14px;
  font-weight: 500;
  color: #2c3e50;
}

.dropdown-icon {
  font-size: 14px;
  color: #8b95a5;
}

/* 主内容包装 */
.main-wrapper {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

/* Hero区域 */
.hero-section {
  background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
  border-radius: 20px;
  padding: 40px;
  margin-bottom: 32px;
}

.hero-content {
  margin-bottom: 24px;
}

.hero-title {
  font-size: 32px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 8px 0;
  line-height: 1.2;
}

.highlight {
  color: #10b981;
}

.hero-subtitle {
  font-size: 16px;
  color: #6b7280;
  margin: 0;
}

/* 快速统计 */
.quick-stats {
  display: flex;
  gap: 24px;
  padding-top: 24px;
  border-top: 1px solid #e5e7eb;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 12px;
}

.stat-icon {
  font-size: 32px;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  line-height: 1;
}

.stat-label {
  font-size: 13px;
  color: #6b7280;
  margin-top: 4px;
}

.stat-divider {
  width: 1px;
  background: #e5e7eb;
}

/* 内容区域 */
.content-section {
  animation: fadeIn 0.5s ease-out;
}

/* 功能网格 */
.features-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 20px;
  margin-bottom: 24px;
}

/* 卡片样式 */
.main-card {
  background: #ffffff;
  border-radius: 16px;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid #e8eaed;
}

.main-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0,0,0,0.12);
  border-color: #10b981;
}

.main-card.disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.main-card.disabled:hover {
  transform: none;
  box-shadow: none;
  border-color: #e8eaed;
}

/* 大卡片（食谱推荐） */
.big-card {
  grid-column: span 2;
  background: linear-gradient(135deg, #10b98110 0%, #05966910 100%);
  border: 2px solid #10b98130;
}

.big-card:hover {
  border-color: #10b981;
}

.card-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.card-emoji {
  font-size: 40px;
}

.card-title-big {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0;
}

.card-desc {
  font-size: 15px;
  color: #6b7280;
  line-height: 1.5;
  margin-bottom: 16px;
}

.card-tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #10b981;
  color: white;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  width: fit-content;
}

.tag-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: white;
  animation: pulse 2s infinite;
}

/* 小卡片 */
.card-header-small {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
}

.card-emoji-small {
  font-size: 28px;
}

.card-title-small {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.card-desc-small {
  font-size: 14px;
  color: #6b7280;
  margin: 0;
}

.card-status {
  margin-top: 12px;
}

/* 提示卡片 */
.tip-card {
  background: linear-gradient(135deg, #fff7ed 0%, #fed7aa 100%);
  border: 2px solid #fed7aa;
  border-radius: 16px;
  padding: 24px;
  display: flex;
  gap: 20px;
  align-items: center;
}

.tip-icon {
  font-size: 48px;
  flex-shrink: 0;
}

.tip-content {
  flex: 1;
}

.tip-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 8px 0;
}

.tip-text {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 16px 0;
}

/* 动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* 响应式设计 */
@media (max-width: 1024px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .big-card {
    grid-column: span 2;
  }
}

@media (max-width: 768px) {
  .hero-section {
    padding: 24px;
  }
  
  .hero-title {
    font-size: 24px;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .big-card {
    grid-column: span 1;
  }
  
  .quick-stats {
    flex-direction: column;
    gap: 16px;
  }
  
  .stat-divider {
    display: none;
  }
  
  .tip-card {
    flex-direction: column;
    text-align: center;
  }
}

@media (max-width: 480px) {
  .top-bar-content {
    padding: 12px 16px;
  }
  
  .main-wrapper {
    padding: 20px 16px;
  }
  
  .hero-title {
    font-size: 20px;
  }
  
  .stat-value {
    font-size: 20px;
  }
}
</style>
