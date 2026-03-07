import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginForm, ChangePasswordForm } from '@/types/user'
import { loginApi, refreshTokenApi } from '@/api/auth'
import { getUserInfoApi, changePasswordApi } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
  // State - 使用 ref 确保响应式
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<User | null>(null)
  const isLoggedIn = computed(() => {
    const hasToken = !!token.value
    const hasUserInfo = !!userInfo.value
    console.log('[AuthStore] isLoggedIn computed:', { hasToken, hasUserInfo })
    return hasToken && hasUserInfo
  })
  const isAdmin = computed(() => userInfo.value?.user_type === 'admin')

  // Actions
  const login = async (credentials: LoginForm) => {
    console.log('[AuthStore] Login started...')
    
    try {
      // 1. 调用登录 API
      console.log('[AuthStore] Calling login API...')
      const res = await loginApi(credentials)
      console.log('[AuthStore] Login API response:', res)
      
      // 2. 检查响应数据
      if (!res || !res.token) {
        throw new Error('登录响应缺少 token')
      }
      
      // 3. 保存 token
      token.value = res.token
      localStorage.setItem('token', res.token)
      console.log('[AuthStore] Token saved, token.value:', token.value?.substring(0, 20) + '...')
      
      // 4. 获取用户信息 - 关键步骤
      console.log('[AuthStore] Fetching user info...')
      const user = await getUserInfo()
      console.log('[AuthStore] User info fetched successfully:', user)
      console.log('[AuthStore] isLoggedIn after login:', isLoggedIn.value)
      
      return true
    } catch (error: any) {
      console.error('[AuthStore] Login failed:', error)
      // 清理状态
      token.value = null
      userInfo.value = null
      localStorage.removeItem('token')
      throw error
    }
  }

  const getUserInfo = async () => {
    console.log('[AuthStore] Getting user info, current token:', !!token.value)
    try {
      const res = await getUserInfoApi()
      console.log('[AuthStore] User info API response:', res)
      
      if (!res) {
        throw new Error('获取用户信息失败：响应为空')
      }
      
      // 直接赋值，确保响应式更新
      userInfo.value = res
      console.log('[AuthStore] userInfo.value set to:', userInfo.value)
      
      return res
    } catch (error: any) {
      console.error('[AuthStore] Get user info failed:', error)
      throw error
    }
  }

  const logout = () => {
    console.log('[AuthStore] Logging out...')
    token.value = null
    userInfo.value = null
    localStorage.removeItem('token')
  }

  const refreshToken = async () => {
    try {
      const res = await refreshTokenApi()
      token.value = res.token
      localStorage.setItem('token', res.token)
      return true
    } catch {
      logout()
      return false
    }
  }

  const changePassword = async (data: ChangePasswordForm) => {
    await changePasswordApi(data)
    return true
  }

  // 初始化时尝试恢复登录状态
  const initAuth = async () => {
    const savedToken = localStorage.getItem('token')
    console.log('[AuthStore] initAuth, savedToken exists:', !!savedToken)
    
    if (savedToken) {
      console.log('[AuthStore] Found saved token, restoring session...')
      token.value = savedToken
      try {
        await getUserInfo()
        console.log('[AuthStore] Session restored, isLoggedIn:', isLoggedIn.value)
        return true
      } catch (error) {
        console.error('[AuthStore] Failed to restore session:', error)
        logout()
        return false
      }
    }
    return false
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    login,
    getUserInfo,
    logout,
    refreshToken,
    changePassword,
    initAuth
  }
})
