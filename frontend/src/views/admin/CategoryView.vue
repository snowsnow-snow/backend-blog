<template>
  <div class="space-y-8">
    <div class="flex justify-between items-end">
      <div>
        <h1 class="text-3xl font-semibold tracking-tight">{{ $t('nav.categories') }}</h1>
        <p class="text-muted-foreground mt-1">管理您的博客分类，让内容井然有序。</p>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
      <!-- Add Category Form -->
      <div class="md:col-span-1">
        <div class="p-6 rounded-2xl border bg-card/50 backdrop-blur-sm space-y-4 shadow-sm">
          <h3 class="font-bold text-sm uppercase tracking-widest opacity-40">{{ $t('common.add') }}</h3>
          
          <div class="space-y-4">
            <div class="space-y-1">
              <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.name') }}</label>
              <input 
                v-model="newCategory.name"
                type="text" 
                class="w-full bg-background border rounded-xl px-4 py-2.5 text-sm focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                placeholder="例如：旅行"
              />
            </div>
            
            <div class="space-y-1">
              <label class="text-xs font-bold uppercase tracking-widest text-muted-foreground/60">{{ $t('common.description') }}</label>
              <textarea 
                v-model="newCategory.description"
                class="w-full bg-background border rounded-xl px-4 py-2 text-sm focus:ring-2 focus:ring-primary/20 outline-none transition-all min-h-[100px] resize-none"
                placeholder="描述这个分类..."
              ></textarea>
            </div>

            <button 
              @click="createCategory"
              :disabled="loading || !newCategory.name"
              class="w-full bg-foreground text-background py-3 rounded-xl font-bold uppercase text-xs tracking-widest hover:scale-[1.02] active:scale-[0.98] transition-all disabled:opacity-50 disabled:hover:scale-100 shadow-lg shadow-foreground/10"
            >
              {{ loading ? $t('common.saving') : $t('common.add') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Categories List -->
      <div class="md:col-span-2">
        <div class="border rounded-2xl overflow-hidden bg-card/30 backdrop-blur-sm shadow-sm">
          <table class="w-full text-left border-collapse">
            <thead class="bg-muted/50 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground/50 border-b">
              <tr>
                <th class="px-6 py-4">{{ $t('common.name') }}</th>
                <th class="px-6 py-4 text-right">{{ $t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-border/50">
              <tr v-for="cat in categories" :key="cat.id" class="hover:bg-muted/20 transition-colors group">
                <td class="px-6 py-4">
                  <div class="font-medium text-sm">{{ cat.name }}</div>
                  <div class="text-[10px] text-muted-foreground/60 mt-0.5">{{ cat.description || '无描述' }}</div>
                </td>
                <td class="px-6 py-4 text-right">
                  <button 
                    @click="deleteCategory(cat.id)"
                    class="text-destructive/40 hover:text-destructive p-2 rounded-full hover:bg-destructive/5 transition-all outline-none"
                  >
                    <Trash2 :size="16" />
                  </button>
                </td>
              </tr>
              <tr v-if="categories.length === 0">
                <td colspan="2" class="px-6 py-20 text-center text-muted-foreground/40 italic text-sm">暂无分类</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Trash2 } from 'lucide-vue-next'
import api from '../../api/client'
import { useUiStore } from '../../stores/ui'

const uiStore = useUiStore()
const categories = ref([])
const loading = ref(false)
const newCategory = ref({
  name: '',
  description: ''
})

const fetchCategories = async () => {
  try {
    const response = await api.get('/api/categories')
    categories.value = response || []
  } catch (err) {
    console.error('Failed to fetch categories:', err)
  }
}

const createCategory = async () => {
  loading.value = true
  try {
    await api.post('/api/categories', newCategory.value)
    newCategory.value = { name: '', description: '' }
    fetchCategories()
    uiStore.showToast('分类创建成功', 'success')
  } catch (err) {
    console.error('Failed to create category:', err)
  } finally {
    loading.value = false
  }
}

const deleteCategory = async (id) => {
  if (!confirm('确定要删除这个分类吗？')) return
  try {
    await api.delete(`/api/categories/${id}`)
    fetchCategories()
    uiStore.showToast('分类已删除', 'success')
  } catch (err) {
    console.error('Failed to delete category:', err)
  }
}

onMounted(fetchCategories)
</script>
