<template>
  <div class="w-full">
    <div class="py-4">
      <!-- Introductory Section Removed -->

      <div v-if="loading" class="space-y-32">
        <div v-for="i in 3" :key="i" class="animate-pulse space-y-8">
          <div class="h-4 bg-muted/50 rounded w-24"></div>
          <div class="h-10 bg-muted/50 rounded w-3/4"></div>
          <div class="h-48 bg-muted/20 rounded-xl w-full"></div>
        </div>
      </div>

      <!-- ... (Error and Empty states remain same but can be tweaked if needed) ... -->
      <div v-else-if="error" class="text-center py-32 flex flex-col items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-rose-500/5 flex items-center justify-center ring-1 ring-rose-500/10">
          <AlertCircle class="w-8 h-8 text-rose-500/40" />
        </div>
        <div class="space-y-2">
          <h3 class="text-lg font-bold tracking-tight">无法获取内容</h3>
          <p class="text-sm text-muted-foreground/60 max-w-[240px] mx-auto leading-relaxed">
            数据在传输过程中遇到了问题，这可能只是暂时的。
          </p>
        </div>
        <button 
          @click="fetchPosts" 
          class="flex items-center gap-2 px-6 py-2.5 rounded-full bg-foreground text-background text-xs font-bold uppercase tracking-widest hover:scale-105 active:scale-95 transition-all shadow-lg"
        >
          <RotateCcw class="w-3.5 h-3.5" />
          <span>{{ $t('common.retry') }}</span>
        </button>
      </div>

      <div v-else-if="filteredPosts.length === 0" class="text-center py-32 flex flex-col items-center gap-6">
        <div class="w-16 h-16 rounded-full bg-muted/40 flex items-center justify-center ring-1 ring-border/5">
          <Ghost class="w-8 h-8 text-muted-foreground/20" />
        </div>
        <div class="space-y-1">
          <h3 class="text-lg serif-font opacity-40 italic">Nothingness</h3>
          <p class="text-xs font-bold uppercase tracking-[0.2em] text-muted-foreground/30">
            暂无相关记录
          </p>
        </div>
      </div>
      <div v-else class="space-y-0">
        <template v-for="(post, index) in filteredPosts" :key="post.id">
          <FeedItem :post="post" />
          <div v-if="index < filteredPosts.length - 1" class="post-separator"></div>
        </template>

        <!-- End of Feed -->
        <div v-if="filteredPosts.length > 0"
             class="py-12 text-center text-[10px] uppercase tracking-[0.3em] text-muted-foreground/40">
          — 已经到底了 —
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HomeView'
}
</script>

<script setup>
import { ref, onMounted, computed, onActivated, nextTick } from 'vue'
import { onBeforeRouteLeave } from 'vue-router'
import { AlertCircle, RotateCcw, Ghost } from 'lucide-vue-next'
import api from '../api/client'
import FeedItem from '../components/FeedItem.vue'
import { useUiStore } from '../stores/ui'

const uiStore = useUiStore()
const posts = ref([])
const loading = ref(true)
const error = ref(null)
const scrollTop = ref(0)

// Save scroll position before leaving
onBeforeRouteLeave((to, from, next) => {
  scrollTop.value = window.scrollY
  next()
})

// Restore scroll position when reactivated by keep-alive
onActivated(() => {
  if (scrollTop.value > 0) {
    nextTick(() => {
      window.scrollTo({
        top: scrollTop.value,
        behavior: 'instant'
      })
    })
  }
})

const filteredPosts = computed(() => {
  if (uiStore.activeTab === 'all') return posts.value
  return posts.value.filter(post => {
    return post.categoryId === uiStore.activeTab
  })
})

const fetchPosts = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await api.get('/posts')
    posts.value = response.list || []
  } catch (err) {
    console.error('Failed to fetch posts:', err)
    error.value = 'Failed to load posts. Please try again later.'
  } finally {
    loading.value = false
    // Ensure we are at the top on initial load if no scroll was intended
    if (scrollTop.value === 0) {
      window.scrollTo(0, 0)
    }
  }
}

onMounted(() => {
  // Only fetch if we don't have posts yet (first load)
  if (posts.value.length === 0) {
    fetchPosts()
  }
})
</script>
