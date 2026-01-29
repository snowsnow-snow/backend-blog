<template>
  <div class="h-screen flex bg-background overflow-hidden relative">
    <!-- Mobile Header -->
    <header class="sm:hidden fixed top-0 left-0 right-0 h-16 border-b bg-background/80 backdrop-blur-xl z-30 flex items-center justify-between px-4">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-primary"></div>
        <span class="font-bold tracking-tight">ADMIN</span>
      </div>
      <button @click="isSidebarOpen = true" class="p-2 hover:bg-muted rounded-xl transition-colors">
        <Menu class="w-6 h-6" />
      </button>
    </header>

    <!-- Sidebar Backdrop -->
    <Transition name="fade">
      <div 
        v-if="isSidebarOpen" 
        @click="isSidebarOpen = false"
        class="sm:hidden fixed inset-0 bg-background/60 backdrop-blur-sm z-40"
      ></div>
    </Transition>

    <!-- Sidebar -->
    <aside 
      class="fixed sm:static inset-y-0 left-0 w-64 border-r flex flex-col bg-background transition-transform duration-300 z-50 pointer-events-auto"
      :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full sm:translate-x-0'"
    >
      <div class="p-6 border-b flex items-center justify-between gap-3 overflow-hidden shrink-0">
        <div class="flex items-center gap-3">
          <div class="w-8 h-8 rounded-lg bg-primary shrink-0"></div>
          <span class="font-bold tracking-tight inline">ADMIN</span>
        </div>
        <button @click="isSidebarOpen = false" class="sm:hidden p-2 hover:bg-muted rounded-lg transition-colors">
          <ChevronLeft class="w-5 h-5 text-muted-foreground" />
        </button>
      </div>

      <nav class="flex-grow p-4 space-y-2 overflow-y-auto scrollbar-hide">
        <router-link
          v-for="item in navItems"
          :key="item.name"
          :to="item.to"
          @click="isSidebarOpen = false"
          class="flex items-center gap-3 p-3 rounded-xl transition-all duration-300 group hover:bg-muted"
          :active-class="item.name === 'dashboard' ? '' : 'bg-gradient-to-r from-primary to-primary/80 text-primary-foreground shadow-lg shadow-primary/25 font-bold scale-[1.02]'"
          exact-active-class="bg-gradient-to-r from-primary to-primary/80 text-primary-foreground shadow-lg shadow-primary/25 font-bold scale-[1.02]"
        >
          <component :is="item.icon" class="w-5 h-5 shrink-0" />
          <span class="font-medium">{{ item.label }}</span>
        </router-link>
      </nav>

      <div class="p-6 border-t space-y-4 shrink-0">
        <div class="flex flex-col gap-4 items-start">
          <LanguageSwitcher />
          <ThemeToggle />
        </div>
        
        <button @click="logout" class="flex items-center gap-3 py-3 w-full rounded-xl hover:bg-destructive/10 text-destructive transition-colors group">
          <LogOut class="w-5 h-5 shrink-0" />
          <span class="font-medium">{{ $t('nav.logout') }}</span>
        </button>
      </div>
    </aside>

    <main class="flex-grow overflow-y-auto mt-16 sm:mt-0">
      <div class="max-w-5xl mx-auto p-4 sm:p-8 lg:p-12">
        <router-view />
      </div>
    </main>
    <Toast />
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ref, computed } from 'vue'
import { LayoutDashboard, FileText, LogOut, Tags, Menu, ChevronLeft, Languages, Palette } from 'lucide-vue-next'
import ThemeToggle from '../components/ThemeToggle.vue'
import LanguageSwitcher from '../components/LanguageSwitcher.vue'
import Toast from '../components/Toast.vue'

const router = useRouter()
const { t } = useI18n({ useScope: 'global' })
const isSidebarOpen = ref(false)


const navItems = computed(() => [
  { name: 'dashboard', label: t('nav.dashboard'), to: '/admin', icon: LayoutDashboard },
  { name: 'new-post', label: t('nav.newPost'), to: '/admin/posts/new', icon: FileText },
  { name: 'categories', label: t('nav.categories'), to: '/admin/categories', icon: Tags }
])


const logout = () => {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>
