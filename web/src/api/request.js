import axios from 'axios'
import { createDiscreteApi } from 'naive-ui'
import { useAuthStore } from '@/store/auth'
import router from '@/router'

const { message } = createDiscreteApi(['message'])

// 创建 axios 实例
const request = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    // 如果有 token，添加到请求头
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    const { response } = error
    
    if (response) {
      switch (response.status) {
        case 401:
          message.error('登录已过期，请重新登录')
          const authStore = useAuthStore()
          authStore.logout()
          router.push('/login')
          break
        case 403:
          message.error(response.data?.error || '权限不足')
          router.push('/403')
          break
        case 400:
          message.error(response.data?.error || '请求参数错误')
          break
        case 404:
          // 404错误可能是路由不存在，但某些情况下是正常的（如营养需求接口在档案不完整时）
          // 不显示通用错误，让调用方自己处理
          break
        case 500:
          message.error('服务器错误，请稍后重试')
          break
        default:
          message.error(response.data?.error || '请求失败')
      }
    } else {
      message.error('网络错误，请检查网络连接')
    }
    
    return Promise.reject(error)
  }
)

export default request
