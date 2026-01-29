<template>
  <div class="p-6 rounded-2xl bg-card shadow-sm border border-muted/50">
    <div class="flex flex-col items-center text-center">
      <div
          class="w-20 h-20 rounded-full mb-4 border-2 border-transparent hover:border-primary transition-all overflow-hidden bg-muted">
        <img src="https://i.pravatar.cc/150?u=0" alt="Me" class="w-full h-full object-cover"/>
      </div>
      <h2 class="font-bold text-lg text-foreground">Alex Design</h2>
      <p class="text-sm mt-1 text-muted-foreground">极简主义探索者</p>

      <div class="flex gap-6 mt-6 text-sm font-medium">
        <div class="flex flex-col">
          <span class="text-foreground">142</span>
          <span class="text-muted-foreground text-[10px] uppercase tracking-wider">关注</span>
        </div>
        <div class="flex flex-col">
          <span class="text-foreground">3.5k</span>
          <span class="text-muted-foreground text-[10px] uppercase tracking-wider">粉丝</span>
        </div>
      </div>
    </div>
  </div>
</template>
