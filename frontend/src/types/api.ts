export interface ApiResponse<T = any> {
  code: string
  message: string
  data?: T
}

export interface ListResponse<T> {
  list: T[]
  total: number
}
