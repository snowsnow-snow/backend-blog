import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useThemeStore = defineStore('theme', () => {
    const themeMode = ref(localStorage.getItem('theme-mode') || 'system')

    const setTheme = (mode) => {
        themeMode.value = mode
        localStorage.setItem('theme-mode', mode)
    }

    const applyTheme = () => {
        const isDark = themeMode.value === 'dark' ||
            (themeMode.value === 'system' && window.matchMedia('(prefers-color-scheme: dark)').matches)

        if (isDark) {
            document.documentElement.classList.add('dark')
        } else {
            document.documentElement.classList.remove('dark')
        }
    }

    watch(themeMode, applyTheme, { immediate: true })

    // Also listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
        if (themeMode.value === 'system') {
            applyTheme()
        }
    })

    return { themeMode, setTheme }
})
