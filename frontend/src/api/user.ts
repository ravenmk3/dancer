import type { User, CreateUserForm, UpdateUserForm, ChangePasswordForm } from '@/types/user'
import request from '@/utils/request'

export const getUserInfoApi = () => {
  return request.post<User>('/me')
}

export const changePasswordApi = (data: ChangePasswordForm) => {
  return request.post('/me/change-password', data)
}

export const getUserListApi = () => {
  return request.post<{ users: User[] }>('/user/list')
}

export const createUserApi = (data: CreateUserForm) => {
  return request.post<User>('/user/create', data)
}

export const updateUserApi = (data: UpdateUserForm) => {
  return request.post('/user/update', data)
}

export const deleteUserApi = (id: string) => {
  return request.post('/user/delete', { id })
}
