<template>
  <div class="min-h-[60vh] flex items-center justify-center">
    <div class="w-full max-w-sm space-y-8 p-8 border rounded-2xl bg-card">
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-semibold">{{ $t('login.title') }}</h1>
        <p class="text-sm text-muted-foreground">{{ $t('login.subtitle') }}</p>
      </div>

      <form @submit.prevent="handleLogin" class="space-y-6">
        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground px-1">{{ $t('login.username') }}</label>
          <input
            v-model="username"
            type="text"
            required
            class="w-full px-4 py-3 rounded-xl border bg-background focus:ring-2 focus:ring-primary outline-none transition-all"
            placeholder="admin"
          />
        </div>
        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground px-1">{{ $t('login.password') }}</label>
          <input
            v-model="password"
            type="password"
            required
            class="w-full px-4 py-3 rounded-xl border bg-background focus:ring-2 focus:ring-primary outline-none transition-all"
            placeholder="••••••••"
          />
        </div>

        <p v-if="error" class="text-xs text-destructive text-center">{{ error }}</p>

        <button
          type="submit"
          :disabled="loading"
          class="w-full py-3 px-4 bg-primary text-primary-foreground rounded-xl font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {{ loading ? $t('login.signingIn') : $t('login.continue') }}
        </button>
      </form>

    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

import api from '../api/client'

const router = useRouter()
const username = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const { t } = useI18n({ useScope: 'global' })


const handleLogin = async () => {
  loading.value = true
  error.value = t('login.error') // Default generic error
  try {
    const response = await api.post('/login', {
      username: username.value,
      password: password.value
    })
    // Assuming backend returns { token: '...' } in response.data
    const token = response.token || response.data?.token
    if (token) {
      localStorage.setItem('token', token)
      router.push({ name: 'admin-dashboard' })
    } else {
      error.value = 'Invalid login response.'
    }
  } catch (err) {
    error.value = t('login.error')
    console.error(err)
  } finally {
    loading.value = false
  }
}

</script>
