<template>
  <div class="fixed top-12 right-6 z-[9999] flex flex-col gap-4 pointer-events-none max-w-[calc(100vw-3rem)] sm:max-w-md">
    <transition-group name="toast">
      <div
        v-for="toast in uiStore.toasts"
        :key="toast.id"
        class="pointer-events-auto group flex items-center gap-4 px-6 py-4 rounded-[24px] border border-white/10 bg-background/60 backdrop-blur-3xl shadow-[0_20px_50px_-12px_rgba(0,0,0,0.15)] ring-1 ring-black/5 transition-all duration-500"
        :class="[
          toast.type === 'success' ? 'text-emerald-600 dark:text-emerald-400' : '',
          toast.type === 'error' ? 'text-rose-600 dark:text-rose-400' : '',
          toast.type === 'warning' ? 'text-amber-600 dark:text-amber-400' : '',
          toast.type === 'info' ? 'text-foreground' : ''
        ]"
      >
        <div 
          class="w-10 h-10 rounded-full flex items-center justify-center shrink-0 transition-transform group-hover:scale-110 duration-500"
          :class="[
            toast.type === 'success' ? 'bg-emerald-500/10' : '',
            toast.type === 'error' ? 'bg-rose-500/10' : '',
            toast.type === 'warning' ? 'bg-amber-500/10' : '',
            toast.type === 'info' ? 'bg-foreground/5' : ''
          ]"
        >
          <component 
            :is="getIcon(toast.type)" 
            class="w-5 h-5" 
            stroke-width="2.5"
          />
        </div>
        
        <div class="flex-1 min-w-0 py-1">
          <p class="text-xs font-bold uppercase tracking-[0.15em] opacity-40 mb-1">
            {{ toast.type }}
          </p>
          <p class="text-sm font-semibold leading-relaxed tracking-tight text-foreground/90">
            {{ toast.message }}
          </p>
        </div>

        <button 
          @click="uiStore.removeToast(toast.id)"
          class="p-2 rounded-full hover:bg-foreground/5 transition-all opacity-0 group-hover:opacity-100 -mr-2"
        >
          <X class="w-4 h-4 text-muted-foreground" />
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup>
import { 
  CheckCircle2, AlertCircle, Info, AlertTriangle, X 
} from 'lucide-vue-next'
import { useUiStore } from '../stores/ui'

const uiStore = useUiStore()

const getIcon = (type) => {
  switch (type) {
    case 'success': return CheckCircle2
    case 'error': return AlertCircle
    case 'warning': return AlertTriangle
    default: return Info
  }
}
</script>

<style scoped>
.toast-enter-active {
  transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}
.toast-leave-active {
  transition: all 0.4s cubic-bezier(0.7, 0, 0.84, 0);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(40px) scale(0.9);
}

.toast-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.toast-move {
  transition: transform 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}
</style>
