import { defineStore } from 'pinia'
import api from '../api/client'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: null,
        token: localStorage.getItem('token') || null,
        loading: false,
        error: null
    }),
    getters: {
        isAuthenticated: (state) => !!state.token
    },
    actions: {
        async login(credentials) {
            this.loading = true
            this.error = null
            try {
                const response = await api.post('/login', credentials)
                this.token = response.token
                this.user = { username: response.username }
                localStorage.setItem('token', this.token)
                return true
            } catch (err) {
                this.error = err.response?.data?.message || 'Login failed'
                return false
            } finally {
                this.loading = false
            }
        },
        logout() {
            this.user = null
            this.token = null
            localStorage.removeItem('token')
        }
    }
})
