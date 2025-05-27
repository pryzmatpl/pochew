<template>
  <div class="bg-white rounded-lg shadow-sm border border-gray-200 hover:shadow-md transition-shadow duration-200">
    <!-- Article Header -->
    <div class="p-4 border-b border-gray-100">
      <div class="flex items-start justify-between">
        <div class="flex-1 min-w-0">
          <h3 class="text-lg font-semibold text-gray-900 truncate">
            <a 
              :href="article.url" 
              target="_blank" 
              rel="noopener noreferrer"
              class="hover:text-blue-600 transition-colors"
            >
              {{ article.title }}
            </a>
          </h3>
          <p class="text-sm text-gray-500 mt-1">
            {{ formatDomain(article.url) }} • {{ formatDate(article.createdAt) }}
          </p>
        </div>
        
        <!-- Action Buttons -->
        <div class="flex items-center space-x-2 ml-4">
          <button
            @click="toggleRead"
            :class="[
              'p-2 rounded-full transition-colors',
              article.isRead 
                ? 'text-green-600 bg-green-50 hover:bg-green-100' 
                : 'text-gray-400 hover:text-green-600 hover:bg-green-50'
            ]"
            :title="article.isRead ? 'Mark as unread' : 'Mark as read'"
          >
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>
          
          <button
            @click="toggleFavorite"
            :class="[
              'p-2 rounded-full transition-colors',
              article.isFavorite 
                ? 'text-yellow-500 bg-yellow-50 hover:bg-yellow-100' 
                : 'text-gray-400 hover:text-yellow-500 hover:bg-yellow-50'
            ]"
            :title="article.isFavorite ? 'Remove from favorites' : 'Add to favorites'"
          >
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          </button>
          
          <button
            @click="toggleArchive"
            :class="[
              'p-2 rounded-full transition-colors',
              article.isArchived 
                ? 'text-blue-600 bg-blue-50 hover:bg-blue-100' 
                : 'text-gray-400 hover:text-blue-600 hover:bg-blue-50'
            ]"
            :title="article.isArchived ? 'Unarchive' : 'Archive'"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4" />
            </svg>
          </button>
          
          <button
            @click="deleteArticle"
            class="p-2 rounded-full text-gray-400 hover:text-red-600 hover:bg-red-50 transition-colors"
            title="Delete article"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </div>
    </div>
    
    <!-- Article Content -->
    <div class="p-4">
      <p v-if="article.description" class="text-gray-700 text-sm leading-relaxed mb-3">
        {{ truncateText(article.description, 200) }}
      </p>
      
      <!-- Tags -->
      <div v-if="article.tags && article.tags.length > 0" class="flex flex-wrap gap-2 mb-3">
        <span
          v-for="tag in article.tags"
          :key="tag"
          class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
        >
          {{ tag }}
        </span>
      </div>
      
      <!-- Reading Progress -->
      <div v-if="article.readingProgress && article.readingProgress > 0" class="mb-3">
        <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
          <span>Reading Progress</span>
          <span>{{ Math.round(article.readingProgress) }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-1.5">
          <div 
            class="bg-blue-600 h-1.5 rounded-full transition-all duration-300"
            :style="{ width: `${article.readingProgress}%` }"
          ></div>
        </div>
      </div>
      
      <!-- Article Footer -->
      <div class="flex items-center justify-between text-xs text-gray-500">
        <div class="flex items-center space-x-4">
          <span v-if="article.readingTime">
            {{ article.readingTime }} min read
          </span>
          <span v-if="article.wordCount">
            {{ formatNumber(article.wordCount) }} words
          </span>
        </div>
        
        <div class="flex items-center space-x-2">
          <span v-if="article.isRead && article.readAt" class="text-green-600">
            Read {{ formatRelativeTime(article.readAt) }}
          </span>
          <span v-else-if="!article.isRead" class="text-orange-600">
            Unread
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useArticlesStore } from '@/stores/articles'

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['updated', 'deleted'])

const articlesStore = useArticlesStore()

const toggleRead = async () => {
  const result = await articlesStore.markAsRead(props.article.id)
  if (result.success) {
    emit('updated', props.article.id)
  }
}

const toggleFavorite = async () => {
  const result = await articlesStore.toggleFavorite(props.article.id)
  if (result.success) {
    emit('updated', props.article.id)
  }
}

const toggleArchive = async () => {
  const result = await articlesStore.toggleArchive(props.article.id)
  if (result.success) {
    emit('updated', props.article.id)
  }
}

const deleteArticle = async () => {
  if (confirm('Are you sure you want to delete this article?')) {
    const result = await articlesStore.deleteArticle(props.article.id)
    if (result.success) {
      emit('deleted', props.article.id)
    }
  }
}

const formatDomain = (url) => {
  try {
    return new URL(url).hostname.replace('www.', '')
  } catch {
    return 'Unknown'
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = (now - date) / (1000 * 60 * 60)
  
  if (diffInHours < 24) {
    return `${Math.floor(diffInHours)}h ago`
  } else if (diffInHours < 24 * 7) {
    return `${Math.floor(diffInHours / 24)}d ago`
  } else {
    return date.toLocaleDateString()
  }
}

const formatRelativeTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInMinutes = (now - date) / (1000 * 60)
  
  if (diffInMinutes < 60) {
    return `${Math.floor(diffInMinutes)}m ago`
  } else if (diffInMinutes < 60 * 24) {
    return `${Math.floor(diffInMinutes / 60)}h ago`
  } else {
    return `${Math.floor(diffInMinutes / (60 * 24))}d ago`
  }
}

const truncateText = (text, maxLength) => {
  if (!text || text.length <= maxLength) return text
  return text.substring(0, maxLength).trim() + '...'
}

const formatNumber = (num) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}
</script> 