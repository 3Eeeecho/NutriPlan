import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as loginApi, getUserProfile } from '@/api/user'

const normalizeRole = (role) => {
  if (role === 'admin' || role === 'user') return role
  return 'guest'
}

const decodeTokenRole = (token) => {
  if (!token) return ''
  try {
    const payload = token.split('.')[1]
    if (!payload) return ''
    const normalized = payload.replace(/-/g, '+').replace(/_/g, '/')
    const padded = normalized.padEnd(normalized.length + ((4 - (normalized.length % 4)) % 4), '=')
    return JSON.parse(atob(padded))?.role || ''
  } catch {
    return ''
  }
}

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))
  const profile = ref(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)
  const role = computed(() => {
    const storedRole = normalizeRole(user.value?.role)
    if (storedRole !== 'guest') return storedRole

    const tokenRole = normalizeRole(decodeTokenRole(token.value))
    if (tokenRole !== 'guest') return tokenRole

    return token.value ? 'user' : 'guest'
  })
  const isAdmin = computed(() => role.value === 'admin')
  const isUser = computed(() => role.value === 'user' || role.value === 'admin')
  const hasProfile = computed(() => {
    if (!profile.value) return false
    return !!(profile.value.age && profile.value.height && profile.value.weight)
  })

  // 登录
  const login = async (username, password) => {
    try {
      const response = await loginApi({ username, password })
      token.value = response.token
      user.value = {
        ...response.user,
        role: response.user?.role || 'user'
      }
      localStorage.setItem('token', response.token)
      localStorage.setItem('user', JSON.stringify(user.value))
      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = () => {
    token.value = ''
    user.value = null
    profile.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  // 加载用户档案
  const loadProfile = async () => {
    try {
      const response = await getUserProfile()
      profile.value = response.user
      // 如果 user 信息为空，从 profile 中更新（包含 username）
      if (!user.value && response.user) {
        user.value = {
          id: response.user.id,
          username: response.user.username,
          email: response.user.email,
          role: response.user.role || 'user'
        }
        localStorage.setItem('user', JSON.stringify(user.value))
      } else if (response.user && user.value) {
        user.value = {
          ...user.value,
          role: response.user.role || user.value.role || 'user'
        }
        localStorage.setItem('user', JSON.stringify(user.value))
      }
      return response.user
    } catch (error) {
      throw error
    }
  }

  // 更新用户档案（本地状态）
  const updateProfile = (newProfile) => {
    profile.value = { ...profile.value, ...newProfile }
  }

  return {
    token,
    user,
    profile,
    isAuthenticated,
    role,
    isAdmin,
    isUser,
    hasProfile,
    login,
    logout,
    loadProfile,
    updateProfile
  }
})

