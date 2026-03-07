import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginForm, ChangePasswordForm } from '@/types/user'
import { loginApi, refreshTokenApi } from '@/api/auth'
import { getUserInfoApi, changePasswordApi } from '@/api/user'

export const useAuthStore = defineStore('auth', () => {
  // State
  const token = ref<string | null>(localStorage.getItem('token'))
  const userInfo = ref<User | null>(null)
  const isLoggedIn = computed(() => !!token.value && !!userInfo.value)
  const isAdmin = computed(() => userInfo.value?.user_type === 'admin')

  // Actions
  const login = async (credentials: LoginForm) => {
    const res = await loginApi(credentials)
    token.value = res.token
    localStorage.setItem('token', res.token)
    await getUserInfo()
    return true
  }

  const getUserInfo = async () => {
    const res = await getUserInfoApi()
    userInfo.value = res
    return res
  }

  const logout = () => {
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

  return {
    token,
    userInfo,
    isLoggedIn,
    isAdmin,
    login,
    getUserInfo,
    logout,
    refreshToken,
    changePassword
  }
})
