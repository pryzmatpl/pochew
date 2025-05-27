<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 overflow-y-auto">
    <!-- Backdrop -->
    <div 
      class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
      @click="closeModal"
    ></div>
    
    <!-- Modal -->
    <div class="flex min-h-full items-center justify-center p-4">
      <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full">
        <!-- Header -->
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">
            Add New Article
          </h3>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Form -->
        <form @submit.prevent="handleSubmit" class="p-6">
          <div class="space-y-4">
            <!-- URL Input -->
            <div>
              <label for="url" class="block text-sm font-medium text-gray-700 mb-2">
                Article URL
              </label>
              <input
                id="url"
                v-model="form.url"
                type="url"
                required
                placeholder="https://example.com/article"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :disabled="isLoading"
              />
              <p class="mt-1 text-xs text-gray-500">
                Enter the URL of the article you want to save
              </p>
            </div>
            
            <!-- Title Input (Optional) -->
            <div>
              <label for="title" class="block text-sm font-medium text-gray-700 mb-2">
                Custom Title (Optional)
              </label>
              <input
                id="title"
                v-model="form.title"
                type="text"
                placeholder="Leave empty to auto-extract from URL"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :disabled="isLoading"
              />
            </div>
            
            <!-- Tags Input -->
            <div>
              <label for="tags" class="block text-sm font-medium text-gray-700 mb-2">
                Tags (Optional)
              </label>
              <input
                id="tags"
                v-model="tagsInput"
                type="text"
                placeholder="technology, programming, web development"
                class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                :disabled="isLoading"
              />
              <p class="mt-1 text-xs text-gray-500">
                Separate tags with commas
              </p>
            </div>
            
            <!-- Options -->
            <div class="space-y-3">
              <div class="flex items-center">
                <input
                  id="markAsRead"
                  v-model="form.markAsRead"
                  type="checkbox"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  :disabled="isLoading"
                />
                <label for="markAsRead" class="ml-2 block text-sm text-gray-700">
                  Mark as read immediately
                </label>
              </div>
              
              <div class="flex items-center">
                <input
                  id="addToFavorites"
                  v-model="form.addToFavorites"
                  type="checkbox"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  :disabled="isLoading"
                />
                <label for="addToFavorites" class="ml-2 block text-sm text-gray-700">
                  Add to favorites
                </label>
              </div>
            </div>
          </div>
          
          <!-- Error Message -->
          <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-md">
            <p class="text-sm text-red-600">{{ error }}</p>
          </div>
          
          <!-- Success Message -->
          <div v-if="success" class="mt-4 p-3 bg-green-50 border border-green-200 rounded-md">
            <p class="text-sm text-green-600">{{ success }}</p>
          </div>
          
          <!-- Actions -->
          <div class="flex justify-end space-x-3 mt-6">
            <button
              type="button"
              @click="closeModal"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 transition-colors"
              :disabled="isLoading"
            >
              Cancel
            </button>
            <button
              type="submit"
              class="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
              :disabled="isLoading || !form.url"
            >
              <span v-if="isLoading" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Adding...
              </span>
              <span v-else>Add Article</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, watch } from 'vue'
import { useArticlesStore } from '@/stores/articles'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'added'])

const articlesStore = useArticlesStore()

const isLoading = ref(false)
const error = ref('')
const success = ref('')
const tagsInput = ref('')

const form = reactive({
  url: '',
  title: '',
  tags: [],
  markAsRead: false,
  addToFavorites: false
})

// Watch for tags input changes
watch(tagsInput, (newValue) => {
  form.tags = newValue
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag.length > 0)
})

// Reset form when modal opens/closes
watch(() => props.isOpen, (newValue) => {
  if (newValue) {
    resetForm()
  }
})

const resetForm = () => {
  form.url = ''
  form.title = ''
  form.tags = []
  form.markAsRead = false
  form.addToFavorites = false
  tagsInput.value = ''
  error.value = ''
  success.value = ''
}

const closeModal = () => {
  if (!isLoading.value) {
    emit('close')
  }
}

const handleSubmit = async () => {
  if (!form.url) return
  
  isLoading.value = true
  error.value = ''
  success.value = ''
  
  try {
    // First, try to capture the URL content
    const captureResult = await articlesStore.captureUrl(form.url)
    
    if (!captureResult.success) {
      error.value = captureResult.error || 'Failed to capture article content'
      return
    }
    
    // Prepare article data
    const articleData = {
      url: form.url,
      title: form.title || captureResult.data.title,
      description: captureResult.data.description,
      content: captureResult.data.content,
      tags: form.tags,
      isRead: form.markAsRead,
      isFavorite: form.addToFavorites,
      readingTime: captureResult.data.readingTime,
      wordCount: captureResult.data.wordCount,
      author: captureResult.data.author,
      publishedAt: captureResult.data.publishedAt
    }
    
    // Create the article
    const result = await articlesStore.createArticle(articleData)
    
    if (result.success) {
      success.value = 'Article added successfully!'
      emit('added', result.data)
      
      // Close modal after a short delay
      setTimeout(() => {
        closeModal()
      }, 1500)
    } else {
      error.value = result.error || 'Failed to add article'
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
    console.error('Error adding article:', err)
  } finally {
    isLoading.value = false
  }
}

// Handle escape key
const handleKeydown = (event) => {
  if (event.key === 'Escape' && props.isOpen && !isLoading.value) {
    closeModal()
  }
}

// Add event listener when component mounts
if (typeof window !== 'undefined') {
  window.addEventListener('keydown', handleKeydown)
}

// Clean up event listener when component unmounts
import { onUnmounted } from 'vue'
onUnmounted(() => {
  if (typeof window !== 'undefined') {
    window.removeEventListener('keydown', handleKeydown)
  }
})
</script> 