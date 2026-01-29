<template>
  <div class="flex items-center gap-1 bg-muted/50 p-1 rounded-xl">
    <button
      v-for="mode in modes"
      :key="mode.id"
      @click="themeStore.setTheme(mode.id)"
      class="p-2 rounded-lg transition-all flex items-center justify-center group relative"
      :class="themeStore.themeMode === mode.id ? 'bg-background text-primary shadow-sm' : 'text-muted-foreground hover:bg-muted'"
      :title="$t('common.' + mode.id)"
    >
      <component :is="mode.icon" class="w-4 h-4" />
      <span v-if="showLabel" class="ml-2 text-xs font-medium hidden sm:inline">{{ $t('common.' + mode.id) }}</span>
    </button>
  </div>
</template>

<script setup>
import { useThemeStore } from '../store/theme'
import { Sun, Moon, Monitor } from 'lucide-vue-next'

const props = defineProps({
  showLabel: {
    type: Boolean,
    default: false
  }
})

const themeStore = useThemeStore()

const modes = [
  { id: 'light', icon: Sun },
  { id: 'dark', icon: Moon },
  { id: 'system', icon: Monitor }
]
</script>
