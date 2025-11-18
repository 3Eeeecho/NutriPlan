<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h1 class="logo">NutriPlan</h1>
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
        <div class="welcome-card">
          <el-card shadow="hover" class="profile-card">
            <template #header>
              <div class="card-header">
                <el-icon class="card-icon"><Document /></el-icon>
                <span>健康档案</span>
              </div>
            </template>
            
            <div v-if="!authStore.hasProfile" class="empty-profile">
              <el-icon class="empty-icon"><Warning /></el-icon>
              <h3>您还没有完善健康档案</h3>
              <p>完善档案后，我们将为您提供个性化的营养计划</p>
              <el-button type="primary" size="large" @click="goToProfile">
                <el-icon><Edit /></el-icon>
                立即完善档案
              </el-button>
            </div>

            <div v-else class="profile-summary">
              <el-descriptions :column="2" border>
                <el-descriptions-item label="年龄">{{ authStore.profile?.age }} 岁</el-descriptions-item>
                <el-descriptions-item label="性别">{{ authStore.profile?.gender }}</el-descriptions-item>
                <el-descriptions-item label="身高">{{ authStore.profile?.height }} cm</el-descriptions-item>
                <el-descriptions-item label="体重">{{ authStore.profile?.weight }} kg</el-descriptions-item>
                <el-descriptions-item label="BMI">{{ authStore.profile?.bmi || '--' }}</el-descriptions-item>
                <el-descriptions-item label="健康目标">{{ authStore.profile?.health_goal || '--' }}</el-descriptions-item>
              </el-descriptions>
              <el-button type="primary" class="update-btn" @click="goToProfile">
                <el-icon><Edit /></el-icon>
                更新档案
              </el-button>
            </div>
          </el-card>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  UserFilled,
  User,
  ArrowDown,
  SwitchButton,
  Document,
  Warning,
  Edit
} from '@element-plus/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

onMounted(async () => {
  if (authStore.isAuthenticated) {
    // 如果没有 user 信息，尝试加载 profile（会同时更新 user）
    if (!authStore.user) {
      try {
        await authStore.loadProfile()
      } catch (error) {
        console.error('加载档案失败:', error)
      }
    } else if (!authStore.profile) {
      try {
        await authStore.loadProfile()
      } catch (error) {
        console.error('加载档案失败:', error)
      }
    }
  }
})

const goToProfile = () => {
  router.push('/profile')
}

const handleCommand = async (command) => {
  if (command === 'profile') {
    router.push('/profile/view')
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

.logo {
  font-size: 24px;
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
  padding: 30px;
  max-width: 1200px;
  margin: 0 auto;
}

.welcome-card {
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

.profile-card {
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

.empty-profile {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 64px;
  color: #c0c4cc;
  margin-bottom: 20px;
}

.empty-profile h3 {
  color: #333;
  margin-bottom: 10px;
}

.empty-profile p {
  color: #666;
  margin-bottom: 30px;
}

.profile-summary {
  padding: 20px 0;
}

.update-btn {
  margin-top: 20px;
  width: 100%;
}

:deep(.el-descriptions__label) {
  font-weight: 600;
}
</style>

