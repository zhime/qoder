import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, refreshToken, getUserInfo } from '@/api/auth'

export interface User {
  id: number
  username: string
  email: string
  role: string
  status: number
}

export const useUserStore = defineStore('user', () => {
  const accessToken = ref<string>('')
  const refreshTokenValue = ref<string>('')
  const user = ref<User | null>(null)

  const isLoggedIn = computed(() => !!accessToken.value)

  // 从localStorage恢复状态
  const loadFromStorage = () => {
    const token = localStorage.getItem('access_token')
    const refresh = localStorage.getItem('refresh_token')
    const userInfo = localStorage.getItem('user_info')
    
    if (token) accessToken.value = token
    if (refresh) refreshTokenValue.value = refresh
    if (userInfo) user.value = JSON.parse(userInfo)
  }

  // 保存到localStorage
  const saveToStorage = () => {
    if (accessToken.value) {
      localStorage.setItem('access_token', accessToken.value)
    }
    if (refreshTokenValue.value) {
      localStorage.setItem('refresh_token', refreshTokenValue.value)
    }
    if (user.value) {
      localStorage.setItem('user_info', JSON.stringify(user.value))
    }
  }

  // 清除存储
  const clearStorage = () => {
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user_info')
  }

  // 登录
  const loginUser = async (username: string, password: string) => {
    try {
      const response = await login(username, password)
      accessToken.value = response.data.access_token
      refreshTokenValue.value = response.data.refresh_token
      user.value = response.data.user
      saveToStorage()
      return response
    } catch (error) {
      throw error
    }
  }

  // 刷新令牌
  const refreshAccessToken = async () => {
    try {
      if (!refreshTokenValue.value) {
        throw new Error('No refresh token available')
      }
      const response = await refreshToken(refreshTokenValue.value)
      accessToken.value = response.data.access_token
      refreshTokenValue.value = response.data.refresh_token
      saveToStorage()
      return response
    } catch (error) {
      logout()
      throw error
    }
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    try {
      const response = await getUserInfo()
      user.value = response.data
      saveToStorage()
      return response
    } catch (error) {
      throw error
    }
  }

  // 登出
  const logout = () => {
    accessToken.value = ''
    refreshTokenValue.value = ''
    user.value = null
    clearStorage()
  }

  // 初始化时从storage加载
  loadFromStorage()

  return {
    accessToken,
    refreshTokenValue,
    user,
    isLoggedIn,
    loginUser,
    refreshAccessToken,
    fetchUserInfo,
    logout
  }
})