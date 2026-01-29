<template>
  <div class="relative group/gallery h-full w-full flex flex-col">
    <!-- Slider Container -->
    <div
        ref="slider"
        class="flex-1 flex flex-row overflow-x-auto overflow-y-hidden snap-x snap-mandatory scroll-smooth hide-scrollbar h-full"
        @scroll="handleScroll"
    >
      <div
          v-for="(asset, index) in mediaAssets"
          :key="asset.id"
          class="w-full h-full flex-none shrink-0 snap-center flex items-center justify-center p-4 sm:p-8"
      >
        <div class="relative w-full h-full flex flex-col items-center justify-center">
          <MediaRenderer
              :id="asset.id"
              :media-type="asset.mediaType"
              class=""
          />
          <div v-if="asset.metadata?.title"
               class="mt-6 text-[10px] font-bold uppercase tracking-[0.2em] text-muted-foreground/40 px-3 py-1">
            {{ asset.metadata.title }}
          </div>
        </div>
      </div>
    </div>

    <!-- Navigation Controls (Desktop Only) -->
    <button
        v-if="mediaAssets.length > 1"
        @click="scrollPrev"
        class="absolute left-4 md:-left-12 top-1/2 -translate-y-1/2 p-2 text-foreground/20 hover:text-primary hover:scale-125 opacity-0 group-hover/gallery:opacity-100 transition-all duration-300 z-10"
    >
      <ChevronLeft class="w-8 h-8"/>
    </button>
    <button
        v-if="mediaAssets.length > 1"
        @click="scrollNext"
        class="absolute right-4 md:-right-12 top-1/2 -translate-y-1/2 p-2 text-foreground/20 hover:text-primary hover:scale-125 opacity-0 group-hover/gallery:opacity-100 transition-all duration-300 z-10"
    >
      <ChevronRight class="w-8 h-8"/>
    </button>

    <!-- Indicators -->
    <div v-if="mediaAssets.length > 1"
         class="absolute -bottom-8 left-1/2 -translate-x-1/2 flex items-center gap-2 px-3 py-1.5">
      <button
          v-for="(_, index) in mediaAssets"
          :key="index"
          @click="scrollTo(index)"
          class="w-1 h-1 rounded-full transition-all duration-500"
          :class="activeIndex === index ? 'bg-primary w-3' : 'bg-primary/10 hover:bg-primary/30'"
      ></button>
    </div>
  </div>
</template>

<script setup>
import {ref, computed} from 'vue'
import {ChevronLeft, ChevronRight} from 'lucide-vue-next'
import MediaRenderer from './MediaRenderer.vue'

const props = defineProps({
  mediaAssets: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['active-change'])

const slider = ref(null)
const activeIndex = ref(0)

const handleScroll = (e) => {
  const scrollLeft = e.target.scrollLeft
  const itemWidth = e.target.clientWidth
  const newIndex = Math.round(scrollLeft / itemWidth)
  if (activeIndex.value !== newIndex) {
    activeIndex.value = newIndex
    emit('active-change', {index: newIndex, asset: props.mediaAssets[newIndex]})
  }
}

const scrollTo = (index) => {
  if (!slider.value) return
  const itemWidth = slider.value.clientWidth
  slider.value.scrollTo({
    left: index * itemWidth,
    behavior: 'smooth'
  })
}

const scrollPrev = () => {
  const prevIndex = activeIndex.value > 0 ? activeIndex.value - 1 : props.mediaAssets.length - 1
  scrollTo(prevIndex)
}

const scrollNext = () => {
  const nextIndex = activeIndex.value < props.mediaAssets.length - 1 ? activeIndex.value + 1 : 0
  scrollTo(nextIndex)
}
</script>

<style scoped>
.hide-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.hide-scrollbar::-webkit-scrollbar {
  display: none;
}
</style>
