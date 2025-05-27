import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// Create axios instance
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
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

// Response interceptor to handle auth errors
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      // Try to refresh token
      const authStore = useAuthStore()
      const refreshResult = await authStore.refreshToken()
      
      if (refreshResult.success) {
        // Retry original request with new token
        const token = localStorage.getItem('token')
        originalRequest.headers.Authorization = `Bearer ${token}`
        return api(originalRequest)
      } else {
        // Refresh failed, logout user
        authStore.logout()
        window.location.href = '/login'
      }
    }
    
    return Promise.reject(error)
  }
)

// Set default axios instance for stores
axios.defaults = api.defaults
axios.interceptors = api.interceptors

export default api 