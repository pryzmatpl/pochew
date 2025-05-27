import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token'))
  const isLoading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const userInitials = computed(() => {
    if (!user.value) return 'U'
    const firstName = user.value.firstName || user.value.username || ''
    const lastName = user.value.lastName || ''
    return (firstName.charAt(0) + lastName.charAt(0)).toUpperCase() || 'U'
  })

  // Actions
  const login = async (credentials) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.post('/api/v1/auth/login', credentials)
      const { token: authToken, user: userData } = response.data
      
      token.value = authToken
      user.value = userData
      localStorage.setItem('token', authToken)
      
      // Set default authorization header
      axios.defaults.headers.common['Authorization'] = `Bearer ${authToken}`
      
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.post('/api/v1/auth/register', userData)
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.response?.data?.error || 'Registration failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    try {
      if (token.value) {
        await axios.post('/api/v1/auth/logout')
      }
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      delete axios.defaults.headers.common['Authorization']
    }
  }

  const fetchUser = async () => {
    if (!token.value) return
    
    try {
      const response = await axios.get('/api/v1/users/me')
      user.value = response.data
    } catch (err) {
      console.error('Failed to fetch user:', err)
      // If token is invalid, logout
      if (err.response?.status === 401) {
        logout()
      }
    }
  }

  const refreshToken = async () => {
    try {
      const response = await axios.post('/api/v1/auth/refresh')
      const { token: newToken } = response.data
      
      token.value = newToken
      localStorage.setItem('token', newToken)
      axios.defaults.headers.common['Authorization'] = `Bearer ${newToken}`
      
      return true
    } catch (err) {
      console.error('Token refresh failed:', err)
      logout()
      return false
    }
  }

  const updateProfile = async (profileData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.put('/api/v1/users/me', profileData)
      user.value = { ...user.value, ...response.data }
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Profile update failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const changePassword = async (passwordData) => {
    isLoading.value = true
    error.value = null
    
    try {
      await axios.post('/api/v1/users/change-password', passwordData)
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Password change failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  // Initialize auth state
  const init = () => {
    if (token.value) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
      fetchUser()
    }
  }

  return {
    // State
    user,
    token,
    isLoading,
    error,
    
    // Getters
    isAuthenticated,
    userInitials,
    
    // Actions
    login,
    register,
    logout,
    fetchUser,
    refreshToken,
    updateProfile,
    changePassword,
    init
  }
}) 