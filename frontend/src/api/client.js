import axios from 'axios'

const api = axios.create({
    baseURL: import.meta.env.VITE_API_BASE_URL,
    timeout: 10000,

    headers: {
        'Content-Type': 'application/json'
    }
})

// Request interceptor for JWT
api.interceptors.request.use(
    (config) => {
        const token = localStorage.getItem('token')
        if (token) {
            config.headers.Authorization = `Bearer ${token}`
        }
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// Response interceptor for error handling
api.interceptors.response.use(
    (response) => {
        const res = response.data
        if (typeof res.code !== 'undefined' && res.code !== 0) {
            const errorMsg = res.msg || 'API Error'
            // We can't use useUiStore() outside of setup/components directly without a slight trick
            // or by importing the store instance after pinia is initialized.
            // But since this is a simple app, we can import it and it should work if called within an action.
            import('../stores/ui').then(m => {
                const uiStore = m.useUiStore()
                uiStore.addToast(errorMsg, 'error')
            })
            return Promise.reject(new Error(errorMsg))
        }
        return res.data !== undefined ? res.data : res
    },

    (error) => {
        if (error.response?.status === 401) {
            localStorage.removeItem('token')
            window.location.href = '/login'
        } else {
            const errorMsg = error.response?.data?.msg || error.message || 'Network Error'
            import('../stores/ui').then(m => {
                const uiStore = m.useUiStore()
                uiStore.addToast(errorMsg, 'error')
            })
        }
        return Promise.reject(error)
    }
)

export default api
