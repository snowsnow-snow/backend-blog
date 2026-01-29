<template>
  <div class="flex items-center gap-1 bg-muted rounded-full p-1">
    <button
      v-for="lang in ['en', 'zh']"
      :key="lang"
      @click="changeLang(lang)"
      class="px-2 py-0.5 text-[10px] font-bold uppercase rounded-full transition-all"
      :class="locale === lang ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'"
    >
      {{ lang }}
    </button>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'

const { locale } = useI18n({ useScope: 'global' })

const changeLang = (lang) => {
  locale.value = lang
  localStorage.setItem('lang', lang)
  dayjs.locale(lang === 'zh' ? 'zh-cn' : 'en')
}

</script>
