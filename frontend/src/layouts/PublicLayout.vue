<template>
  <div class="min-h-screen transition-colors duration-300 bg-background text-foreground">
    <!-- --- Navigation Bar --- -->
    <header v-if="!isPostDetail" class="sticky top-0 z-50 w-full bg-background/90 backdrop-blur-md">
      <div class="mx-auto max-w-[1024px] px-8 py-10 flex items-center justify-between">
        <div class="flex items-center gap-3">
        </div>
        <nav class="flex items-center gap-6 sm:gap-12">
          <button
              v-for="tab in dynamicTabs"
              :key="tab.id"
              @click="uiStore.setActiveTab(tab.id)"
              class="text-[11px] font-bold tracking-[0.2em] uppercase transition-all hover-link"
              :class="uiStore.activeTab === tab.id ? 'text-foreground' : 'text-muted-foreground'"
          >
            {{ tab.label }}
          </button>
          
          <!-- Theme Toggle (Segmented) -->
          <div class="flex items-center gap-1 bg-muted/30 p-1 rounded-full border border-border/10">
            <button
              v-for="mode in ['light', 'dark', 'system']"
              :key="mode"
              @click="uiStore.setTheme(mode)"
              class="p-1.5 rounded-full transition-all flex items-center justify-center group relative"
              :class="uiStore.theme === mode ? 'bg-background text-primary shadow-sm' : 'text-muted-foreground/60 hover:text-foreground'"
              :title="$t('common.' + mode)"
            >
              <component :is="mode === 'light' ? Sun : (mode === 'dark' ? Moon : Monitor)" :size="13" />
            </button>
          </div>
        </nav>
      </div>
    </header>

    <!-- --- Main Layout --- -->
    <main class="mx-auto w-full max-w-[1024px] flex-1 px-8">
      <!-- Central Feed -->
      <div class="w-full">
        <router-view v-slot="{ Component, route }">
          <transition name="fade" mode="out-in">
            <keep-alive :include="['HomeView']">
              <component :is="Component" :key="route.meta.keepAlive ? 'home' : route.fullPath"/>
            </keep-alive>
          </transition>
        </router-view>

      </div>
    </main>
    <Toast/>
  </div>
</template>

<script setup>
import {ref, computed, onMounted} from 'vue'
import {useRoute} from 'vue-router'
import { Share2, MoreHorizontal, Search, Settings, Feather, Menu, X, Sun, Moon, Monitor } from 'lucide-vue-next'
import { useUiStore } from '../stores/ui'
import api from '../api/client'
import Toast from '../components/Toast.vue'
import ProfileCard from '../components/ProfileCard.vue'
import HotTags from '../components/HotTags.vue'

const uiStore = useUiStore()
const route = useRoute()
const isMobileMenuOpen = ref(false)
const categories = ref([])

const isPostDetail = computed(() => route.name === 'post-detail')

const dynamicTabs = computed(() => {
  const base = [{ id: 'all', label: '全部' }]
  const dynamic = categories.value.map(c => ({
    id: c.id,
    label: c.name
  }))
  return [...base, ...dynamic]
})

const fetchCategories = async () => {
  try {
    const response = await api.get('/api/categories')
    categories.value = response || []
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

onMounted(() => {
  uiStore.initTheme()
  fetchCategories()
})

</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.98);
}

.fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(1.02);
}

.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: all 0.2s ease-out;
}

.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
  transform: translateY(5px);
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
