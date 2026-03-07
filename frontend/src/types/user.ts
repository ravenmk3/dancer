export type UserType = 'admin' | 'normal'

export interface User {
  id: string
  username: string
  user_type: UserType
  created_at: number
  updated_at: number
}

export interface LoginForm {
  username: string
  password: string
}

export interface CreateUserForm {
  username: string
  password: string
  user_type: UserType
}

export interface UpdateUserForm {
  id: string
  username?: string
  password?: string
  user_type?: UserType
}

export interface ChangePasswordForm {
  old_password: string
  new_password: string
}
