<template>
  <router-link
      v-if="post && (post.slug || post.id)"
      :to="{ name: 'post-detail', params: { slug: post.slug || post.id } }"
      class="group block space-y-4"
  >

    <!-- Cover Image for Article or Gallery Preview -->
    <div
        v-if="post.postType === 'gallery' || post.coverImageId"
        class="aspect-video relative overflow-hidden rounded-lg bg-muted"
    >
      <img
          :src="coverUrl"
          :alt="post.title"
          class="absolute inset-0 h-full w-full object-cover transition-transform duration-500 group-hover:scale-105"
          v-if="coverUrl"
      />
      <div v-else class="flex items-center justify-center h-full text-muted-foreground">
        <ImageIcon v-if="post.postType === 'gallery'" class="w-8 h-8"/>
        <span v-else>No cover</span>
      </div>
      <!-- Type Badge -->
      <div
          class="absolute top-3 left-3 px-2 py-1 text-[10px] uppercase font-bold tracking-widest bg-background/90 text-foreground rounded border">
        {{ post.postType }}
      </div>
    </div>

    <!-- Text only/Article style -->
    <div class="space-y-2">
      <h2 class="text-xl font-medium group-hover:underline underline-offset-4 decoration-1">
        {{ post.title }}
      </h2>
      <p v-if="post.summary" class="text-muted-foreground line-clamp-2 leading-relaxed">
        {{ post.summary }}
      </p>
      <div class="flex items-center text-xs text-muted-foreground uppercase tracking-widest pt-2">
        <span>{{ formatDate(post.createdTime) }}</span>
      </div>
    </div>
  </router-link>
</template>

<script setup>
import {computed} from 'vue'
import {useI18n} from 'vue-i18n'
import dayjs from 'dayjs'

import {Image as ImageIcon} from 'lucide-vue-next'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const coverUrl = computed(() => {
  if (!props.post.coverImageId) return null
  return `${import.meta.env.VITE_API_BASE_URL}/media/${props.post.coverImageId}`
})


const {locale} = useI18n()

const formatDate = (date) => {
  // Access locale.value to make it reactive
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  return dayjs(date).locale(currentLocale).format('LL')
}

</script>
