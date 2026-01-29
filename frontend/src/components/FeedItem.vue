<template>
  <router-link :to="{ name: 'post-detail', params: { id: post.id } }" class="block group/post">
    <article class="relative">
      <header :class="post.postType === 'article' ? 'mb-0' : 'mb-8'">
        <div class="flex items-center gap-4 text-[10px] tracking-[0.3em] text-muted-foreground font-bold uppercase mb-6">
          <span>{{ formattedFullDate }}</span>
          <span v-if="post.category" class="w-1 h-1 bg-border rounded-full"></span>
          <span v-if="post.category">{{ post.category.name }}</span>
        </div>
        <div class="space-y-4">
          <h1 class="text-2xl md:text-4xl text-foreground font-normal leading-tight transition-opacity max-w-3xl">
            {{ post.title }}
          </h1>
          <p v-if="post.postType === 'article' && post.summary" class="text-muted-foreground text-lg md:text-xl leading-relaxed max-w-2xl font-light line-clamp-2">
            {{ post.summary }}
          </p>
        </div>
      </header>

      <!-- Content Area (Only for non-article types) -->
      <div v-if="post.postType !== 'article'" class="space-y-8">
        <!-- Media Grid (If Gallery) -->
        <div 
          v-if="post.mediaAssets?.length && isMedia" 
          class="grid gap-2"
          :class="limitedMedia.length === 1 ? 'grid-cols-1' : (limitedMedia.length === 2 ? 'grid-cols-2' : 'grid-cols-3')"
        >
          <div 
            v-for="(asset, index) in limitedMedia" 
            :key="asset.id"
            class="overflow-hidden aspect-square bg-muted relative group"
          >
            <img 
              :src="getMediaUrl(asset.id)" 
              class="w-full h-full bg-center bg-no-repeat bg-cover transition-transform duration-1000 group-hover:scale-105"
              loading="lazy"
            />
            
            <!-- +N Overlay for the last visible item if more exist -->
            <div v-if="index === 2 && post.mediaTotal > 3" class="absolute inset-0 image-overlay flex items-center justify-center transition-all duration-300 group-hover:bg-background/95">
               <span class="text-foreground text-[11px] font-bold tracking-[0.3em]">+{{ post.mediaTotal - 3 }} 更多</span>
            </div>
          </div>
        </div>

        <!-- Text Summary -->
        <p v-if="displayContent" class="text-muted-foreground text-xl leading-relaxed max-w-2xl font-light line-clamp-3">
          {{ displayContent }}
        </p>
      </div>
    </article>
  </router-link>
</template>

<script setup>
import {computed} from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import localizedFormat from 'dayjs/plugin/localizedFormat'
import 'dayjs/locale/zh-cn'
import MarkdownRenderer from './MarkdownRenderer.vue'

dayjs.extend(relativeTime)
dayjs.extend(localizedFormat)

const { locale } = useI18n()

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

// Logic to determine layout type
const displayContent = computed(() => {
  return props.post.summary || props.post.content || ''
})

const isMedia = computed(() => {
  if (props.post.postType === 'article') return false
  return props.post.postType === 'gallery' || props.post.postType === 'video' || (props.post.mediaAssets?.length > 0)
})

const formattedDate = computed(() => dayjs(props.post.createdTime).format('YYYY-MM-DD'))

const formattedFullDate = computed(() => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  return dayjs(props.post.createdTime).locale(currentLocale).format('LL')
})

// Media Helpers
const limitedMedia = computed(() => {
  if (!props.post.mediaAssets) return []
  return props.post.mediaAssets.slice(0, 3)
})


const getMediaUrl = (id) => {
  return `${import.meta.env.VITE_API_BASE_URL}/media/${id}`
}
</script>
