<template>
  <div class="space-y-4">
    <h3 class="text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground/50 px-3">热门标签</h3>
    <div class="space-y-1">
      <div 
        v-for="tag in tags" 
        :key="tag" 
        class="flex items-center justify-between p-3 rounded-xl cursor-pointer transition-all hover:bg-card hover:shadow-sm group"
      >
        <span class="text-sm text-foreground/70 group-hover:text-primary transition-colors">#{{ tag }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const tags = ['Minimalism', 'React', 'Photography', 'LifeStyle']
</script>
