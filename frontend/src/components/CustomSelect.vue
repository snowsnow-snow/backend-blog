<template>
  <div class="relative inline-block text-left" ref="containerRef">
    <!-- Trigger -->
    <button
      type="button"
      @click="isOpen = !isOpen"
      class="flex items-center justify-between gap-2 px-4 py-2 rounded-xl bg-muted/50 hover:bg-muted transition-all border border-border/60 focus:border-primary/30 outline-none w-full min-w-[120px]"
      :class="[containerClass]"
    >
      <span class="text-sm font-medium truncate">{{ selectedLabel || placeholder }}</span>
      <ChevronDown 
        class="w-4 h-4 text-muted-foreground transition-transform duration-300" 
        :class="{ 'rotate-180': isOpen }"
      />
    </button>

    <!-- Dropdown Menu -->
    <Transition
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="transform scale-95 opacity-0 -translate-y-2"
      enter-to-class="transform scale-100 opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in"
      leave-from-class="transform scale-100 opacity-100 translate-y-0"
      leave-to-class="transform scale-95 opacity-0 -translate-y-2"
    >
      <div
        v-if="isOpen"
        class="absolute z-50 mt-2 w-full min-w-[160px] origin-top-right rounded-2xl bg-background border shadow-xl shadow-foreground/5 backdrop-blur-xl p-1 outline-none right-0"
      >
        <div class="py-1 space-y-0.5">
          <button
            v-for="option in options"
            :key="option.value"
            @click="selectOption(option)"
            class="w-full text-left px-3 py-2 rounded-xl text-sm transition-all flex items-center justify-between group"
            :class="[
              modelValue === option.value 
                ? 'bg-primary text-primary-foreground font-bold' 
                : 'hover:bg-muted text-foreground'
            ]"
          >
            <span>{{ option.label }}</span>
            <Check v-if="modelValue === option.value" class="w-4 h-4" />
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ChevronDown, Check } from 'lucide-vue-next'

const props = defineProps({
  modelValue: [String, Number],
  options: {
    type: Array,
    required: true,
    // Expects array of { label: string, value: any }
  },
  placeholder: {
    type: String,
    default: 'Select an option'
  },
  containerClass: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const isOpen = ref(false)
const containerRef = ref(null)

const selectedLabel = computed(() => {
  const option = props.options.find(opt => opt.value === props.modelValue)
  return option ? option.label : ''
})

const selectOption = (option) => {
  emit('update:modelValue', option.value)
  emit('change', option.value)
  isOpen.value = false
}

const handleClickOutside = (event) => {
  if (containerRef.value && !containerRef.value.contains(event.target)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
