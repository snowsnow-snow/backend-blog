<template>
  <div class="relative group/media flex items-center justify-center w-full h-full">
    <template v-if="mediaType === 'video'">
      <video
          :src="mediaUrl"
          controls
          class="max-w-full max-h-full object-contain"
          @error="handleError"
      ></video>
    </template>

    <template v-else>
      <img
          v-if="!error"
          :src="mediaUrl"
          :alt="id"
          class="block max-w-full max-h-full object-contain transition-opacity duration-300"
          loading="lazy"
          @error="handleError"
      />
    </template>

    <div v-if="error"
         class="absolute inset-0 flex flex-col items-center justify-center bg-muted text-destructive text-xs p-4">
      <svg xmlns="http://www.w3.org/2000/svg" class="h-8 w-8 mb-2 opacity-50" fill="none" viewBox="0 0 24 24"
           stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
      </svg>
      <span>Failed to load media</span>
    </div>
  </div>
</template>

<script setup>
import {computed, ref} from 'vue'

const props = defineProps({
  id: {
    type: [String, Number],
    required: true
  },
  mediaType: {
    type: String,
    default: 'image'
  }
})

const error = ref(false)

// 构建媒体资源 URL
const mediaUrl = computed(() => {
  // 确保环境变量 VITE_API_BASE_URL 已正确配置
  const baseUrl = import.meta.env.VITE_API_BASE_URL || ''
  return `${baseUrl}/media/${props.id}`
})

const handleError = () => {
  error.value = true
}
</script>

<style scoped>
/* 如果需要针对不同比例做特殊微调，可以在此添加 */
</style>