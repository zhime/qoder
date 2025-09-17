import request from './request'
import type { User } from '@/store/user'

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: User
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  access_token: string
  refresh_token: string
}

export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

// 登录
export const login = (username: string, password: string) => {
  return request.post<ApiResponse<LoginResponse>>('/auth/login', {
    username,
    password
  })
}

// 刷新令牌
export const refreshToken = (refreshToken: string) => {
  return request.post<ApiResponse<RefreshTokenResponse>>('/auth/refresh', {
    refresh_token: refreshToken
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request.get<ApiResponse<User>>('/auth/profile')
}

// 退出登录
export const logout = () => {
  return request.post('/auth/logout')
}