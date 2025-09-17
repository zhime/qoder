import axios from 'axios'
import { useUserStore } from '@/store/user'
import { ElMessage } from 'element-plus'

const baseURL = import.meta.env.VITE_API_BASE_URL || '/api'

// 创建axios实例
const request = axios.create({
  baseURL,
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const userStore = useUserStore()
    if (userStore.accessToken) {
      config.headers.Authorization = `Bearer ${userStore.accessToken}`
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
    const res = response.data
    
    // 如果响应格式不是标准格式，直接返回
    if (typeof res !== 'object' || res === null) {
      return response
    }
    
    // 处理业务错误
    if (res.code && res.code !== 200) {
      ElMessage.error(res.message || '请求失败')
      return Promise.reject(new Error(res.message || '请求失败'))
    }
    
    return res
  },
  async (error) => {
    const userStore = useUserStore()
    
    if (error.response?.status === 401) {
      // 尝试刷新token
      try {
        await userStore.refreshAccessToken()
        // 重新发送原请求
        return request(error.config)
      } catch (refreshError) {
        // 刷新失败，跳转到登录页
        userStore.logout()
        window.location.href = '/login'
        ElMessage.error('登录已过期，请重新登录')
        return Promise.reject(refreshError)
      }
    }
    
    const message = error.response?.data?.message || error.message || '网络错误'
    ElMessage.error(message)
    return Promise.reject(error)
  }
)

export default request