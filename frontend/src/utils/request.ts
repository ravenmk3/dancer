import axios from 'axios'
import type { AxiosResponse, InternalAxiosRequestConfig } from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types/api'
import router from '@/router'

// Create axios instance
const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
axiosInstance.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    const { code, message, data } = response.data

    if (code === 'success') {
      return data
    }

    ElMessage.error(message || '请求失败')
    return Promise.reject(new Error(message || '请求失败'))
  },
  (error) => {
    if (error.response) {
      const status = error.response.status
      const data = error.response.data as ApiResponse

      if (status === 401) {
        // 登录接口的 401 错误由登录页面处理，不在这里显示消息
        const isLoginRequest = error.config?.url?.includes('/auth/login')
        if (!isLoginRequest) {
          localStorage.removeItem('token')
          router.push('/login')
          ElMessage.error('登录已过期，请重新登录')
        }
        return Promise.reject(error)
      } else if (status === 403) {
        ElMessage.error(data.message || '没有权限执行此操作')
      } else {
        ElMessage.error(data.message || '请求失败')
      }
    } else {
      ElMessage.error('网络错误，请检查网络连接')
    }
    return Promise.reject(error)
  }
)

// Custom request function that properly types the response
const request = <T = any>(config: InternalAxiosRequestConfig): Promise<T> => {
  return axiosInstance(config) as Promise<T>
}

// Add convenience methods
request.get = <T = any>(url: string, config?: Omit<InternalAxiosRequestConfig, 'url' | 'method'>): Promise<T> => {
  return axiosInstance.get(url, config) as Promise<T>
}

request.post = <T = any>(url: string, data?: any, config?: Omit<InternalAxiosRequestConfig, 'url' | 'method' | 'data'>): Promise<T> => {
  return axiosInstance.post(url, data, config) as Promise<T>
}

request.put = <T = any>(url: string, data?: any, config?: Omit<InternalAxiosRequestConfig, 'url' | 'method' | 'data'>): Promise<T> => {
  return axiosInstance.put(url, data, config) as Promise<T>
}

request.delete = <T = any>(url: string, config?: Omit<InternalAxiosRequestConfig, 'url' | 'method'>): Promise<T> => {
  return axiosInstance.delete(url, config) as Promise<T>
}

export default request
