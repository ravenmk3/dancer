import axios from 'axios'
import type { AxiosInstance, AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types/api'
import router from '@/router'

const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    console.log('[Request]', config.method?.toUpperCase(), config.url)
    return config
  },
  (error) => {
    console.error('[Request Error]', error)
    return Promise.reject(error)
  }
)

// Response interceptor
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    console.log('[Response]', response.config.url, 'Status:', response.status)
    console.log('[Response Data]', JSON.stringify(response.data, null, 2))
    
    const { code, message, data } = response.data
    console.log('[Response Parsed]', { code, message, hasData: !!data })
    
    if (code === 'success') {
      console.log('[Response] Success, returning data:', data)
      return data
    }
    
    console.error('[Response] Error code:', code, 'Message:', message)
    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message || '请求失败'))
  },
  (error) => {
    console.error('[Response Error]', error)
    if (error.response) {
      const status = error.response.status
      const data = error.response.data as ApiResponse
      console.error('[Response Error Details]', { status, data })
      
      if (status === 401) {
        localStorage.removeItem('token')
        router.push('/login')
        ElMessage.error('登录已过期，请重新登录')
      } else if (status === 403) {
        ElMessage.error(data.message || '没有权限执行此操作')
      } else {
        ElMessage.error(data.message || '请求失败')
      }
    } else {
      console.error('[Response Error] No response:', error.message)
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

export default request
