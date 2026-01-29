<template>
  <div class="h-screen flex flex-col bg-background overflow-hidden">
    <!-- Header -->
    <header class="flex-none flex justify-between items-center bg-background py-4 border-b px-4 sm:px-8 z-10">
      <div class="flex items-center gap-4">
        <router-link to="/admin" class="p-2 hover:bg-muted rounded-full transition-colors">
          <ArrowLeft class="w-5 h-5" />
        </router-link>
        <h1 class="text-xl font-semibold">{{ isEdit ? $t('editor.editPost') : $t('editor.newPost') }}</h1>
      </div>
      <div class="flex items-center gap-3">
        <CustomSelect
          v-model="post.status"
          :options="[
            { label: $t('common.draft'), value: 'draft' },
            { label: $t('common.published'), value: 'published' }
          ]"
        />
        <button
          @click="savePost"
          :disabled="saving"
          class="bg-primary text-primary-foreground px-6 py-2 rounded-lg font-medium hover:opacity-90 transition-opacity disabled:opacity-50"
        >
          {{ saving ? $t('common.saving') : $t('common.save') }}
        </button>
      </div>

    </header>

    <div class="flex-1 flex flex-col overflow-y-auto px-4 sm:px-8 py-8 w-full max-w-4xl mx-auto space-y-8 custom-scrollbar">
      <!-- Metadata Section -->
      <div class="space-y-6 flex-none">
        <div v-if="post.postType !== 'article'" class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.title') }}</label>
          <input
            v-model="post.title"
            type="text"
            :placeholder="$t('editor.titlePlaceholder')"
            class="w-full text-3xl font-bold bg-transparent border-none outline-none placeholder:text-muted-foreground/30"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="space-y-2">
            <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('editor.postType') }}</label>
            <div class="flex gap-2 pt-1">
              <button
                v-for="type in ['article', 'gallery']"
                :key="type"
                @click="post.postType = type"
                class="flex-1 px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all border"
                :class="post.postType === type ? 'bg-foreground text-background border-foreground' : 'bg-background text-muted-foreground border-border/40 hover:border-border'"
              >
                {{ $t('editor.' + type) }}
              </button>
            </div>
          </div>

          <div class="space-y-2">
            <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.category') }}</label>
            <div class="flex flex-wrap gap-2 pt-1">
              <button
                @click="post.categoryId = ''"
                class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all border"
                :class="post.categoryId === '' ? 'bg-foreground text-background border-foreground' : 'bg-background text-muted-foreground border-border/40 hover:border-border'"
              >
                {{ $t('common.noSelection') || 'None' }}
              </button>
              <button
                v-for="cat in categories"
                :key="cat.id"
                @click="post.categoryId = cat.id"
                class="px-4 py-1.5 rounded-full text-[10px] font-bold uppercase tracking-widest transition-all border"
                :class="post.categoryId === cat.id ? 'bg-foreground text-background border-foreground' : 'bg-background text-muted-foreground border-border/40 hover:border-border'"
              >
                {{ cat.name }}
              </button>
            </div>
          </div>
        </div>

        <div class="space-y-2">
          <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.summary') }}</label>
          <textarea 
            v-model="post.summary" 
            :placeholder="$t('editor.summaryPlaceholder')" 
            class="w-full p-3 rounded-xl border border-border bg-card/30 text-sm min-h-[60px] h-[60px] resize-none focus:ring-2 focus:ring-primary/20 outline-none transition-all placeholder:text-muted-foreground/50"
          ></textarea>
        </div>
      </div>

      <!-- Content Area -->
      <div v-if="post.postType === 'article'" class="flex-1 flex flex-col min-h-0 space-y-4">
        <label class="flex-none text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.content') }}</label>
        
        <div class="flex-none flex items-center gap-4 h-24">
          <div 
            @click="triggerFileInput"
            class="flex-1 h-full border-2 border-dashed rounded-2xl flex flex-col items-center justify-center space-y-1 border-muted hover:border-primary/50 transition-colors cursor-pointer group bg-muted/5"
          >
            <Plus v-if="!uploading" class="w-4 h-4 text-muted-foreground transition-transform group-hover:scale-110" />
            <div v-else class="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
            <p class="text-[10px] font-medium uppercase tracking-wider text-muted-foreground">{{ uploading ? $t('editor.uploading') : $t('editor.uploadImage') }}</p>
          </div>

          <div v-if="uploadedImages.length" class="flex-none w-2/3 h-full bg-muted/20 rounded-2xl p-2 overflow-x-auto overflow-y-hidden custom-scrollbar flex items-center gap-3">
             <div v-for="(img, idx) in uploadedImages" :key="img.id" class="flex-none group relative w-16 h-16 rounded-lg overflow-hidden bg-muted">
                <img :src="getMediaUrl(img.id)" class="w-full h-full object-cover" />
                <button 
                  @click="uploadedImages.splice(idx, 1)" 
                  class="absolute top-1 right-1 p-1 bg-background/80 rounded-full hover:bg-destructive hover:text-white transition-all opacity-0 group-hover:opacity-100"
                >
                  <X class="w-3 h-3" />
                </button>
                <button 
                  @click="copyMarkdown(img.id)"
                  class="absolute inset-0 bg-primary/80 text-white text-[8px] font-bold uppercase flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                >
                  {{ $t('editor.copyMarkdown') }}
                </button>
             </div>
          </div>
        </div>

        <textarea
          ref="editorTextarea"
          v-model="post.content"
          :placeholder="$t('editor.contentPlaceholder')"
          class="flex-1 bg-card/30 border border-border rounded-2xl p-6 outline-none resize-none leading-relaxed text-lg focus:ring-2 focus:ring-primary/5 focus:border-primary/30 transition-all placeholder:text-muted-foreground/50 overflow-y-auto"
        ></textarea>
      </div>


      <div v-else-if="post.postType === 'gallery'" class="flex-1 flex flex-col min-h-0 space-y-4">
        <label class="flex-none text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('editor.gallery') }}</label>
        
        <div 
          @click="triggerFileInput"
          class="flex-none h-24 border-2 border-dashed rounded-2xl flex flex-col items-center justify-center space-y-1 border-muted hover:border-primary/50 transition-colors cursor-pointer group bg-muted/5"
        >
          <Plus v-if="!uploading" class="w-4 h-4 text-muted-foreground transition-transform group-hover:scale-110" />
          <div v-else class="w-4 h-4 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
          <p class="text-[10px] font-medium uppercase tracking-wider text-muted-foreground">{{ uploading ? $t('editor.uploading') : $t('editor.uploadPhotos') }}</p>
        </div>
        
        <div v-if="post.mediaAssets?.length" class="flex-1 overflow-y-auto custom-scrollbar p-1">
          <div class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-5 gap-4">
            <div v-for="(asset, idx) in post.mediaAssets" :key="idx" class="aspect-square relative rounded-xl overflow-hidden bg-muted group shadow-sm border border-border/50">
               <img :src="getMediaUrl(asset.id)" class="w-full h-full object-cover transition-transform group-hover:scale-105" />
               <button @click="post.mediaAssets.splice(idx, 1)" class="absolute top-2 right-2 p-1.5 bg-background/80 rounded-full hover:bg-destructive hover:text-white transition-all opacity-0 group-hover:opacity-100 shadow-lg">
                 <X class="w-4 h-4" />
               </button>
            </div>
          </div>
        </div>
      </div>

      <input 
        ref="fileInput"
        type="file" 
        multiple 
        class="hidden" 
        accept="image/*"
        @change="handleFileUpload" 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowLeft, Plus, X } from 'lucide-vue-next'
import api from '../../api/client'

import { useUiStore } from '../../stores/ui'
import { useI18n } from 'vue-i18n'
import CustomSelect from '../../components/CustomSelect.vue'

const router = useRouter()
const route = useRoute()
const uiStore = useUiStore()
const { t } = useI18n()
const isEdit = computed(() => !!route.params.id)
const fileInput = ref(null)
const editorTextarea = ref(null)

const categories = ref([])
const post = ref({
  title: '',
  summary: '',
  content: '',
  postType: 'article',
  categoryId: '',
  status: 'draft',
  mediaAssets: []
})

watch(() => post.value.content, (newContent) => {
  if (post.value.postType !== 'article' || !newContent) return
  
  // Extract first line
  const firstLine = newContent.split('\n')[0].trim()
  if (firstLine) {
    // Remove leading # and space
    post.value.title = firstLine.replace(/^#+\s*/, '')
  } else {
    post.value.title = ''
  }
})

const saving = ref(false)
const uploading = ref(false)
const uploadedImages = ref([])

const copyMarkdown = (id) => {
  const markdown = `![](${getMediaUrl(id)})`
  navigator.clipboard.writeText(markdown)
  uiStore.addToast(t('editor.markdownCopied'), 'success')
}

const getMediaUrl = (id) => {
  return `${import.meta.env.VITE_API_BASE_URL}/media/${id}`
}

const fetchCategories = async () => {
  try {
    const response = await api.get('/api/categories')
    categories.value = response || []
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

const fetchPost = async () => {
  if (!isEdit.value) return
  try {
    const response = await api.get(`/api/posts/${route.params.id}`)
    post.value = response
    if (post.value.postType === 'article') {
      uploadedImages.value = post.value.mediaAssets || []
    }
  } catch (err) {
    console.error(err)
  }
}

const triggerFileInput = () => {
  if (uploading.value) return
  fileInput.value.click()
}

const handleFileUpload = async (event) => {
  const files = event.target.files
  if (!files.length) return
  
  uploading.value = true
  try {
    // Select multi-files, but upload individually as requested
    for (let i = 0; i < files.length; i++) {
      const formData = new FormData()
      formData.append('file', files[i])
      
      // Separate request for each file
      const response = await api.post('/api/upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      
      // axios interceptor returns res.data
      if (response && response.id) {
        if (post.value.postType === 'gallery') {
          post.value.mediaAssets.push(response)
        } else {
          uploadedImages.value.push(response)
        }
      }
    }
  } catch (err) {
    console.error('Upload failed:', err)
    uiStore.addToast(t('editor.uploadFail'), 'error')
  } finally {
    uploading.value = false
    // Clear the input so the same file can be selected again
    event.target.value = ''
  }
}

const savePost = async () => {
  if (!post.value.title) {
    uiStore.addToast(t('editor.titleRequired') || 'Title is required', 'warning')
    return
  }
  
  if (post.value.postType === 'article') {
    post.value.mediaAssets = uploadedImages.value
  }
  
  saving.value = true
  try {
    if (isEdit.value) {
      await api.put(`/api/posts/${route.params.id}`, post.value)
    } else {
      await api.post('/api/posts', post.value)
    }
    uiStore.addToast(t('editor.success'), 'success')
    router.push('/admin')
  } catch (err) {
    console.error(err)
    uiStore.addToast(t('editor.fail'), 'error')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  fetchPost()
  fetchCategories()
})
</script>
