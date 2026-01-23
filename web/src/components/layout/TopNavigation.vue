<template>
  <div class="top-navigation">
    <div class="nav-content">
      <!-- Logo 区域 -->
      <div class="brand" @click="goHome">
        <span class="brand-icon">🥗</span>
        <span class="brand-name">NutriPlan</span>
      </div>

      <!-- 用户区域 -->
      <div class="user-section">
        <n-dropdown
          :options="dropdownOptions"
          @select="handleSelect"
          placement="bottom-end"
        >
          <div class="user-menu">
            <n-avatar
              round
              :size="36"
              :style="{ background: '#10b981' }"
            >
              <n-icon size="20">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                </svg>
              </n-icon>
            </n-avatar>
            <span class="user-name">{{ username }}</span>
            <n-icon size="16" class="dropdown-icon">
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
                <path fill="currentColor" d="M7 10l5 5 5-5z"/>
              </svg>
            </n-icon>
          </div>
        </n-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, h } from 'vue'
import { useRouter } from 'vue-router'
import { NDropdown, NAvatar, NIcon } from 'naive-ui'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()

const username = computed(() => authStore.user?.username || '用户')

const dropdownOptions = [
  {
    label: '个人档案',
    key: 'profile',
    icon: () => h('span', '👤'),
  },
  {
    label: '我的收藏',
    key: 'favorites',
    icon: () => h('span', '⭐'),
  },
  {
    label: '饮食记录',
    key: 'intake',
    icon: () => h('span', '📝'),
  },
  {
    label: '周报告',
    key: 'weekly',
    icon: () => h('span', '📊'),
  },
  {
    type: 'divider',
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h('span', '🚪'),
  },
]

const handleSelect = (key) => {
  switch (key) {
    case 'profile':
      router.push('/profile/view')
      break
    case 'favorites':
      router.push('/favorites')
      break
    case 'intake':
      router.push('/intake')
      break
    case 'weekly':
      router.push('/weekly')
      break
    case 'logout':
      authStore.logout()
      router.push('/login')
      break
  }
}

const goHome = () => {
  router.push('/home')
}
</script>

<style scoped>
.top-navigation {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--color-border-primary);
  position: sticky;
  top: 0;
  z-index: var(--z-index-sticky);
  box-shadow: var(--shadow-xs);
}

.nav-content {
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
  cursor: pointer;
  transition: opacity var(--transition-fast);
}

.brand:hover {
  opacity: 0.8;
}

.brand-icon {
  font-size: 24px;
}

.brand-name {
  font-size: 20px;
  font-weight: var(--font-weight-bold);
  color: var(--color-text-primary);
  letter-spacing: -0.5px;
}

.user-menu {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 12px;
  border-radius: var(--radius-full);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.user-menu:hover {
  background: var(--color-bg-hover);
}

.user-name {
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-medium);
  color: var(--color-text-primary);
}

.dropdown-icon {
  color: var(--color-text-tertiary);
}

@media (max-width: 768px) {
  .nav-content {
    padding: 12px 16px;
  }

  .brand-name {
    font-size: 18px;
  }

  .user-name {
    display: none;
  }
}
</style>
