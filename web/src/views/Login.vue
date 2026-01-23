<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <div class="logo">🥗</div>
        <h1 class="title">NutriPlan</h1>
        <p class="subtitle">智能营养计划系统</p>
      </div>

      <n-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        size="large"
      >
        <n-form-item path="username">
          <n-input
            v-model:value="loginForm.username"
            placeholder="请输入用户名"
            clearable
          >
            <template #prefix>
              <n-icon :component="PersonOutline" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item path="password">
          <n-input
            v-model:value="loginForm.password"
            type="password"
            placeholder="请输入密码"
            show-password-on="click"
            @keyup.enter="handleLogin"
          >
            <template #prefix>
              <n-icon :component="LockClosedOutline" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item>
          <n-button
            type="primary"
            size="large"
            block
            :loading="loading"
            @click="handleLogin"
            class="login-button"
          >
            {{ loading ? '登录中...' : '登录' }}
          </n-button>
        </n-form-item>

        <div class="login-footer">
          <span>还没有账号？</span>
          <n-button text type="primary" @click="goToRegister">
            立即注册
          </n-button>
        </div>
      </n-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { NForm, NFormItem, NInput, NButton, NIcon, useMessage } from 'naive-ui'
import { PersonOutline, LockClosedOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()
const message = useMessage()

const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true
    
    await authStore.login(loginForm.username, loginForm.password)
    message.success('登录成功')
    
    // 获取重定向路径
    const redirect = router.currentRoute.value.query.redirect || '/home'
    router.push(redirect)
  } catch (error) {
    if (error.errors) {
      // 表单验证错误
      return
    }
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  router.push('/register')
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-lg);
  background: linear-gradient(135deg, var(--color-primary) 0%, #764ba2 100%);
}

.login-box {
  width: 100%;
  max-width: 420px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: var(--radius-2xl);
  padding: var(--spacing-3xl);
  box-shadow: var(--shadow-2xl);
  animation: slideUp 0.5s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-header {
  text-align: center;
  margin-bottom: var(--spacing-3xl);
}

.logo {
  font-size: 64px;
  margin-bottom: var(--spacing-md);
}

.title {
  font-size: var(--font-size-3xl);
  font-weight: var(--font-weight-bold);
  background: linear-gradient(135deg, var(--color-primary) 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 var(--spacing-sm) 0;
}

.subtitle {
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  margin: 0;
}

.login-button {
  height: 48px;
  font-size: var(--font-size-base);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-lg);
  margin-top: var(--spacing-md);
}

.login-footer {
  text-align: center;
  margin-top: var(--spacing-lg);
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--spacing-xs);
}

@media (max-width: 480px) {
  .login-box {
    padding: var(--spacing-xl);
  }
  
  .logo {
    font-size: 48px;
  }
  
  .title {
    font-size: var(--font-size-2xl);
  }
}
</style>

