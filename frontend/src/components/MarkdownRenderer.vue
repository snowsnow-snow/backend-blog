<template>
  <div class="prose dark:prose-invert max-w-none" v-html="renderedContent"></div>
</template>

<script setup>
import {computed} from 'vue'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github-dark.css' // Could also dynamically switch

const props = defineProps({
  content: {
    type: String,
    default: ''
  },
  removeFirstTitle: {
    type: Boolean,
    default: false
  }
})

const md = new MarkdownIt({
  html: true,
  linkify: true,
  typographer: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, {language: lang}).value;
      } catch (__) {
      }
    }
    return ''; // use external default escaping
  }
})

const renderedContent = computed(() => {
  let contentToRender = props.content
  if (props.removeFirstTitle && contentToRender) {
    const lines = contentToRender.split('\n')
    // Find first non-empty line
    let firstLineIdx = -1
    for (let i = 0; i < lines.length; i++) {
        if (lines[i].trim()) {
            firstLineIdx = i
            break
        }
    }
    
    // If first non-empty line starts with # (Markdown header), remove it
    if (firstLineIdx !== -1 && lines[firstLineIdx].trim().startsWith('#')) {
      lines.splice(firstLineIdx, 1)
      contentToRender = lines.join('\n')
    }
  }
  return md.render(contentToRender)
})
</script>

<style>
@reference "../style.css";

/* Base prose styles if tailwind-typography isn't available, 
   but we'll assume standard minimalist prose styles here */
.prose h1 {
  @apply text-2xl font-bold mb-4 text-foreground;
}

.prose h2 {
  @apply text-xl font-bold mb-3 mt-6 text-foreground;
}

.prose h3 {
  @apply text-lg font-bold mb-2 mt-4 text-foreground;
}

.prose p {
  @apply mb-3 leading-relaxed text-muted-foreground/70 font-medium;
}

.prose pre {
  @apply bg-muted/30 border border-border/10 p-4 rounded-xl overflow-x-auto my-4;
}

.prose code {
  @apply font-mono text-xs bg-muted/50 px-1 py-0.5 rounded;
}

.prose blockquote {
  @apply border-l-2 border-border/50 pl-4 italic my-4 text-muted-foreground/60;
}

.prose img {
  @apply rounded-xl mx-auto my-6;
}

.prose ul, .prose ol {
  @apply mb-3 pl-5 text-muted-foreground/70;
}

.prose li {
  @apply mb-1;
}

.prose a {
  @apply text-foreground/80 underline underline-offset-4 decoration-border hover:text-foreground transition-all;
}
</style>
