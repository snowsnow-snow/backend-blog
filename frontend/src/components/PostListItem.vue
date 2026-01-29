<template>
  <router-link
    v-if="post && (post.slug || post.id)"
    :to="{ name: 'post-detail', params: { slug: post.slug || post.id } }"
    class="flex items-center gap-6 group py-4 border-b border-muted/20 last:border-0 hover:bg-muted/5 px-4 -mx-4 rounded-xl transition-all duration-300"
  >
    <!-- Thumbnail -->
    <div 
      v-if="coverUrl"
      class="w-12 h-12 flex-shrink-0 relative overflow-hidden rounded-lg bg-muted"
    >
      <img
        :src="coverUrl"
        :alt="post.title"
        class="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-110"
      />
    </div>

    <!-- Date & Title -->
    <div class="flex items-baseline gap-6 min-w-0 flex-grow">
      <span class="text-[10px] font-mono text-muted-foreground/60 tracking-wider uppercase whitespace-nowrap pt-0.5">
        {{ formatDate(post.createdTime) }}
      </span>
      <h2 class="text-base font-normal text-foreground/90 group-hover:text-primary group-hover:translate-x-1 transition-all duration-300 truncate">
        {{ post.title }}
      </h2>
    </div>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const { locale } = useI18n()

const coverUrl = computed(() => {
  if (!props.post.coverImageId) return null
  return `${import.meta.env.VITE_API_BASE_URL}/media/${props.post.coverImageId}`
})

const formatDate = (date) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  return dayjs(date).locale(currentLocale).format('YYYY.MM.DD')
}
</script>
