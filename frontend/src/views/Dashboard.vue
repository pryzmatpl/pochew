<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Navigation Header -->
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold text-gray-900">Read-It-Later</h1>
          </div>
          <div class="flex items-center space-x-4">
            <button
              @click="showAddModal = true"
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Add Article
            </button>
            <div class="relative">
              <button
                @click="showUserMenu = !showUserMenu"
                class="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
              >
                <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center text-white font-medium">
                  {{ authStore.userInitials }}
                </div>
              </button>
              <div
                v-if="showUserMenu"
                class="origin-top-right absolute right-0 mt-2 w-48 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 z-10"
              >
                <div class="py-1">
                  <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Profile</a>
                  <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-100">Settings</a>
                  <button
                    @click="logout"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Sign out
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="flex flex-col lg:flex-row gap-8">
        <!-- Sidebar -->
        <div class="lg:w-64 flex-shrink-0">
          <div class="bg-white rounded-lg shadow p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Filters</h3>
            <div class="space-y-3">
              <button
                @click="setFilter('all')"
                :class="[
                  'w-full text-left px-3 py-2 rounded-md text-sm',
                  currentFilter === 'all' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                ]"
              >
                All Articles ({{ articlesStore.stats.totalArticles }})
              </button>
              <button
                @click="setFilter('unread')"
                :class="[
                  'w-full text-left px-3 py-2 rounded-md text-sm',
                  currentFilter === 'unread' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                ]"
              >
                Unread ({{ articlesStore.stats.unreadArticles }})
              </button>
              <button
                @click="setFilter('favorites')"
                :class="[
                  'w-full text-left px-3 py-2 rounded-md text-sm',
                  currentFilter === 'favorites' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                ]"
              >
                Favorites ({{ articlesStore.stats.favoriteArticles }})
              </button>
              <button
                @click="setFilter('archived')"
                :class="[
                  'w-full text-left px-3 py-2 rounded-md text-sm',
                  currentFilter === 'archived' ? 'bg-blue-100 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                ]"
              >
                Archived ({{ articlesStore.stats.archivedArticles }})
              </button>
            </div>
          </div>
        </div>

        <!-- Main Content -->
        <div class="flex-1">
          <!-- Search Bar -->
          <div class="mb-6">
            <div class="relative">
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search articles..."
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
              <svg class="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>

          <!-- Loading State -->
          <div v-if="articlesStore.isLoading" class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>

          <!-- Error State -->
          <div v-else-if="articlesStore.error" class="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
            <p class="text-red-600">{{ articlesStore.error }}</p>
            <button 
              @click="loadArticles" 
              class="mt-2 text-sm text-red-600 hover:text-red-800 underline"
            >
              Try again
            </button>
          </div>

          <!-- Articles List -->
          <div v-else class="space-y-4">
            <ArticleCard
              v-for="article in filteredArticles"
              :key="article.id"
              :article="article"
              @updated="handleArticleUpdated"
              @deleted="handleArticleDeleted"
            />

            <!-- Empty State -->
            <div v-if="filteredArticles.length === 0" class="text-center py-12">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No articles found</h3>
              <p class="mt-1 text-sm text-gray-500">
                {{ searchQuery ? 'Try adjusting your search terms.' : 'Get started by adding your first article.' }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Article Modal -->
    <AddArticleModal
      :is-open="showAddModal"
      @close="showAddModal = false"
      @added="handleArticleAdded"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useArticlesStore } from '@/stores/articles'
import ArticleCard from '@/components/ArticleCard.vue'
import AddArticleModal from '@/components/AddArticleModal.vue'

const router = useRouter()
const authStore = useAuthStore()
const articlesStore = useArticlesStore()

// Reactive data
const showAddModal = ref(false)
const showUserMenu = ref(false)
const currentFilter = ref('all')
const searchQuery = ref('')

// Computed properties
const filteredArticles = computed(() => {
  let articles = articlesStore.getArticlesByFilter(currentFilter.value)
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    articles = articles.filter(article => 
      article.title?.toLowerCase().includes(query) ||
      article.description?.toLowerCase().includes(query) ||
      article.url?.toLowerCase().includes(query) ||
      article.tags?.some(tag => tag.toLowerCase().includes(query))
    )
  }
  
  return articles
})

// Methods
const setFilter = (filter) => {
  currentFilter.value = filter
}

const logout = async () => {
  await authStore.logout()
  router.push('/login')
}

const loadArticles = async () => {
  await articlesStore.fetchArticles()
  await articlesStore.fetchStats()
}

const handleArticleAdded = (article) => {
  // Article is already added to store by the modal
  // Just refresh stats
  articlesStore.updateStats()
}

const handleArticleUpdated = (articleId) => {
  // Article is already updated in store
  // Just refresh stats
  articlesStore.updateStats()
}

const handleArticleDeleted = (articleId) => {
  // Article is already removed from store
  // Just refresh stats
  articlesStore.updateStats()
}

// Close user menu when clicking outside
const closeUserMenu = (event) => {
  if (!event.target.closest('.relative')) {
    showUserMenu.value = false
  }
}

// Lifecycle hooks
onMounted(async () => {
  await loadArticles()
  document.addEventListener('click', closeUserMenu)
})

// Watch for search query changes with debounce
let searchTimeout
watch(searchQuery, (newQuery) => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(async () => {
    if (newQuery.trim()) {
      const result = await articlesStore.searchArticles(newQuery.trim())
      if (!result.success) {
        console.error('Search failed:', result.error)
      }
    }
  }, 300)
})

// Cleanup
import { onUnmounted } from 'vue'
onUnmounted(() => {
  document.removeEventListener('click', closeUserMenu)
  clearTimeout(searchTimeout)
})
</script> 