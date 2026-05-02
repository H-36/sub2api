<template>
  <AppLayout>
    <div class="image-playground-page">
      <iframe
        ref="iframeRef"
        :src="iframeSrc"
        class="image-playground-frame"
        title="Image Playground"
        referrerpolicy="same-origin"
        allow="clipboard-read; clipboard-write"
        @load="syncFrameTheme"
      />
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'

type PageTheme = 'light' | 'dark'

const initialTheme = detectTheme()
const iframeRef = ref<HTMLIFrameElement | null>(null)
const pageTheme = ref<PageTheme>(initialTheme)
const pageBackground = ref(detectBackground(initialTheme))

let themeObserver: MutationObserver | null = null

function detectTheme(): PageTheme {
  if (typeof document === 'undefined') return 'light'
  return document.documentElement.classList.contains('dark') ? 'dark' : 'light'
}

function detectBackground(theme: PageTheme): string {
  if (typeof window === 'undefined') {
    return theme === 'dark' ? 'rgb(2, 6, 23)' : 'rgb(249, 250, 251)'
  }

  const background = window.getComputedStyle(document.body).backgroundColor
  if (!background || background === 'rgba(0, 0, 0, 0)' || background === 'transparent') {
    return theme === 'dark' ? 'rgb(2, 6, 23)' : 'rgb(249, 250, 251)'
  }
  return background
}

function updateThemeState() {
  const nextTheme = detectTheme()
  pageTheme.value = nextTheme
  pageBackground.value = detectBackground(nextTheme)
}

function buildIframeSrc() {
  if (typeof window === 'undefined') {
    return '/image-playground-app/'
  }

  const params = new URLSearchParams({
    apiUrl: `${window.location.origin}/v1`,
    provider: 'openai',
    apiMode: 'images',
    sub2apiImagePlayground: '1',
    theme: pageTheme.value,
    sub2apiBg: pageBackground.value
  })

  return `/image-playground-app/?${params.toString()}`
}

const iframeSrc = ref(buildIframeSrc())

function syncFrameTheme() {
  if (typeof window === 'undefined') return
  const targetWindow = iframeRef.value?.contentWindow
  if (!targetWindow) return

  targetWindow.postMessage(
    {
      type: 'sub2api:image-playground-theme',
      theme: pageTheme.value,
      background: pageBackground.value
    },
    window.location.origin,
  )
}

onMounted(() => {
  updateThemeState()
  iframeSrc.value = buildIframeSrc()
  syncFrameTheme()

  themeObserver = new MutationObserver(() => {
    updateThemeState()
    syncFrameTheme()
  })
  themeObserver.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['class'],
  })
})

onUnmounted(() => {
  if (themeObserver) {
    themeObserver.disconnect()
    themeObserver = null
  }
})
</script>

<style scoped>
.image-playground-page {
  height: calc(100vh - 64px - 2rem);
  min-height: 640px;
}

.image-playground-frame {
  display: block;
  width: 100%;
  height: 100%;
  border: 0;
  background: transparent;
}

@media (min-width: 768px) {
  .image-playground-page {
    height: calc(100vh - 64px - 3rem);
  }
}

@media (min-width: 1024px) {
  .image-playground-page {
    height: calc(100vh - 64px - 4rem);
  }
}
</style>
