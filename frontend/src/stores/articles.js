import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useArticlesStore = defineStore('articles', () => {
  // State
  const articles = ref([])
  const currentArticle = ref(null)
  const isLoading = ref(false)
  const error = ref(null)
  const stats = ref({
    totalArticles: 0,
    readArticles: 0,
    unreadArticles: 0,
    favoriteArticles: 0,
    archivedArticles: 0,
    storageUsed: 0
  })

  // Getters
  const unreadArticles = computed(() => 
    articles.value.filter(article => !article.isRead)
  )
  
  const favoriteArticles = computed(() => 
    articles.value.filter(article => article.isFavorite)
  )
  
  const archivedArticles = computed(() => 
    articles.value.filter(article => article.isArchived)
  )

  const getArticlesByFilter = computed(() => (filter) => {
    switch (filter) {
      case 'unread':
        return unreadArticles.value
      case 'favorites':
        return favoriteArticles.value
      case 'archived':
        return archivedArticles.value
      default:
        return articles.value.filter(article => !article.isArchived)
    }
  })

  // Actions
  const fetchArticles = async (params = {}) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.get('/api/v1/articles', { params })
      articles.value = response.data.articles || []
      updateStats()
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch articles'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const fetchArticle = async (id) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.get(`/api/v1/articles/${id}`)
      currentArticle.value = response.data
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch article'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const createArticle = async (articleData) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.post('/api/v1/articles', articleData)
      const newArticle = response.data
      articles.value.unshift(newArticle)
      updateStats()
      return { success: true, data: newArticle }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to create article'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const updateArticle = async (id, updateData) => {
    error.value = null
    
    try {
      const response = await axios.put(`/api/v1/articles/${id}`, updateData)
      const updatedArticle = response.data
      
      const index = articles.value.findIndex(article => article.id === id)
      if (index !== -1) {
        articles.value[index] = { ...articles.value[index], ...updatedArticle }
      }
      
      updateStats()
      return { success: true, data: updatedArticle }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update article'
      return { success: false, error: error.value }
    }
  }

  const deleteArticle = async (id) => {
    error.value = null
    
    try {
      await axios.delete(`/api/v1/articles/${id}`)
      articles.value = articles.value.filter(article => article.id !== id)
      updateStats()
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to delete article'
      return { success: false, error: error.value }
    }
  }

  const markAsRead = async (id) => {
    try {
      await axios.post(`/api/v1/articles/${id}/read`)
      const article = articles.value.find(a => a.id === id)
      if (article) {
        article.isRead = true
        article.readAt = new Date().toISOString()
      }
      updateStats()
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to mark as read'
      return { success: false, error: error.value }
    }
  }

  const toggleFavorite = async (id) => {
    try {
      await axios.post(`/api/v1/articles/${id}/favorite`)
      const article = articles.value.find(a => a.id === id)
      if (article) {
        article.isFavorite = !article.isFavorite
      }
      updateStats()
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to toggle favorite'
      return { success: false, error: error.value }
    }
  }

  const toggleArchive = async (id) => {
    try {
      await axios.post(`/api/v1/articles/${id}/archive`)
      const article = articles.value.find(a => a.id === id)
      if (article) {
        article.isArchived = !article.isArchived
      }
      updateStats()
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to toggle archive'
      return { success: false, error: error.value }
    }
  }

  const searchArticles = async (query, filters = {}) => {
    isLoading.value = true
    error.value = null
    
    try {
      const params = { q: query, ...filters }
      const response = await axios.get('/api/v1/search/articles', { params })
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.response?.data?.error || 'Search failed'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const captureUrl = async (url) => {
    isLoading.value = true
    error.value = null
    
    try {
      const response = await axios.post('/api/v1/capture/url', { url })
      return { success: true, data: response.data }
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to capture URL'
      return { success: false, error: error.value }
    } finally {
      isLoading.value = false
    }
  }

  const fetchStats = async () => {
    try {
      const response = await axios.get('/api/v1/users/stats')
      stats.value = response.data
      return { success: true }
    } catch (err) {
      console.error('Failed to fetch stats:', err)
      return { success: false }
    }
  }

  const updateStats = () => {
    stats.value = {
      totalArticles: articles.value.length,
      readArticles: articles.value.filter(a => a.isRead).length,
      unreadArticles: articles.value.filter(a => !a.isRead).length,
      favoriteArticles: articles.value.filter(a => a.isFavorite).length,
      archivedArticles: articles.value.filter(a => a.isArchived).length,
      storageUsed: articles.value.reduce((total, a) => total + (a.storageSize || 0), 0)
    }
  }

  const clearError = () => {
    error.value = null
  }

  const reset = () => {
    articles.value = []
    currentArticle.value = null
    error.value = null
    stats.value = {
      totalArticles: 0,
      readArticles: 0,
      unreadArticles: 0,
      favoriteArticles: 0,
      archivedArticles: 0,
      storageUsed: 0
    }
  }

  return {
    // State
    articles,
    currentArticle,
    isLoading,
    error,
    stats,
    
    // Getters
    unreadArticles,
    favoriteArticles,
    archivedArticles,
    getArticlesByFilter,
    
    // Actions
    fetchArticles,
    fetchArticle,
    createArticle,
    updateArticle,
    deleteArticle,
    markAsRead,
    toggleFavorite,
    toggleArchive,
    searchArticles,
    captureUrl,
    fetchStats,
    updateStats,
    clearError,
    reset
  }
}) 