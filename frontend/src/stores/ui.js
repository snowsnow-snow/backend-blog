import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
    state: () => ({
        activeTab: 'all',
        theme: localStorage.getItem('theme-mode') || 'system', // 'light', 'dark', 'system'
        toasts: [],
        mediaQuery: null
    }),
    getters: {
        darkMode: (state) => {
            if (state.theme === 'system') {
                return window.matchMedia('(prefers-color-scheme: dark)').matches
            }
            return state.theme === 'dark'
        }
    },
    actions: {
        addToast(message, type = 'info', duration = 3000) {
            const id = Date.now()
            this.toasts.push({ id, message, type })
            if (duration > 0) {
                setTimeout(() => {
                    this.removeToast(id)
                }, duration)
            }
        },
        removeToast(id) {
            this.toasts = this.toasts.filter(toast => toast.id !== id)
        },
        setActiveTab(tab) {
            this.activeTab = tab
        },
        setTheme(theme) {
            this.theme = theme
            localStorage.setItem('theme-mode', theme)
            this.applyTheme()
        },
        toggleTheme() {
            // Cycle: system -> light -> dark -> system
            const themes = ['system', 'light', 'dark']
            const currentIndex = themes.indexOf(this.theme)
            const nextIndex = (currentIndex + 1) % themes.length
            this.setTheme(themes[nextIndex])
        },
        applyTheme() {
            const isDark = this.darkMode
            if (isDark) {
                document.documentElement.classList.add('dark')
            } else {
                document.documentElement.classList.remove('dark')
            }
        },
        initTheme() {
            this.theme = localStorage.getItem('theme-mode') || 'system'
            this.applyTheme()

            // Listen for system theme changes
            if (!this.mediaQuery) {
                this.mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
                this.mediaQuery.addEventListener('change', () => {
                    if (this.theme === 'system') {
                        this.applyTheme()
                    }
                })
            }
        }
    }
})
