<template>
  <div class="space-y-8">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">{{ $t('nav.dashboard') }}</h1>
      </div>

      <router-link
          to="/admin/posts/new"
          class="bg-primary text-primary-foreground px-4 py-2 rounded-lg font-medium hover:opacity-90 transition-opacity whitespace-nowrap"
      >
        {{ $t('nav.newPost') }}
      </router-link>

    </div>

    <!-- Stats or summary could go here -->

    <div class="border rounded-2xl overflow-hidden bg-card">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse min-w-[600px]">
          <thead class="bg-muted/50 text-xs font-bold uppercase tracking-widest text-muted-foreground border-b">
          <tr>
            <th class="px-6 py-4">{{ $t('common.status') }}</th>
            <th class="px-6 py-4 min-w-[300px]">{{ $t('common.title') }}</th>
            <th class="px-6 py-4">{{ $t('common.type') }}</th>
            <th class="px-6 py-4 min-w-[140px]">{{ $t('common.date') }}</th>

            <th class="px-6 py-4 text-right">{{ $t('common.actions') }}</th>
          </tr>
          </thead>

          <tbody class="divide-y">
          <tr v-for="post in posts" :key="post.id" class="hover:bg-muted/30 transition-colors">
            <td class="px-6 py-4">
                <span
                    class="px-2 py-0.5 rounded-full text-[10px] uppercase font-bold border"
                    :class="post.status === 'published' ? 'bg-green-100 text-green-700 dark:bg-green-900/30' : 'bg-orange-100 text-orange-700 dark:bg-orange-900/30'"
                >
                  {{ post.status }}
                </span>
            </td>
            <td class="px-6 py-4 font-medium">{{ post.title }}</td>
            <td class="px-6 py-4"><span class="text-xs text-muted-foreground">{{ post.postType }}</span></td>
            <td class="px-6 py-4 text-sm text-muted-foreground whitespace-nowrap">{{ formatDate(post.createdTime) }}</td>

            <td class="px-6 py-4 text-right space-x-3">
              <router-link :to="`/admin/posts/${post.id}/edit`" class="text-primary hover:underline text-sm font-medium">
                {{ $t('common.edit') }}
              </router-link>
              <button @click="deletePost(post.id)" class="text-destructive hover:underline text-sm font-medium">
                {{ $t('common.delete') }}
              </button>
            </td>

          </tr>
          <tr v-if="posts.length === 0">
            <td colspan="5" class="px-6 py-20 text-center text-muted-foreground italic">{{ $t('home.noPosts') }}</td>
          </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref, onMounted} from 'vue'
import {useI18n} from 'vue-i18n'
import dayjs from 'dayjs'

import api from '../../api/client'

const {t, locale} = useI18n()
const posts = ref([])


const fetchPosts = async () => {
  try {
    const response = await api.get('/api/posts')
    posts.value = response.list || []

  } catch (err) {
    console.error('Failed to fetch posts:', err)
  }
}

const deletePost = async (id) => {
  if (!confirm('Are you sure you want to delete this post?')) return
  try {
    await api.delete(`/api/posts/${id}`)
    fetchPosts()
  } catch (err) {
    alert('Failed to delete post.')
  }
}

const formatDate = (date) => {
  const currentLocale = locale.value === 'zh' ? 'zh-cn' : 'en'
  return dayjs(date).locale(currentLocale).format('ll')
}


onMounted(fetchPosts)
</script>
