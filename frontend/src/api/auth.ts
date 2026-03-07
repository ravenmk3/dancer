import type { LoginForm } from '@/types/user'
import request from '@/utils/request'

export const loginApi = (data: LoginForm) => {
  return request.post<{ token: string }>('/auth/login', data)
}

export const refreshTokenApi = () => {
  return request.post<{ token: string }>('/auth/refresh')
}
