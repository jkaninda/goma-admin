import { defineStore } from 'pinia'
import { ref, computed, watch, onUnmounted } from 'vue'

type ThemeMode = 'light' | 'dark' | 'system'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(
    (localStorage.getItem('goma_theme') as ThemeMode) || 'system'
  )

  const systemDark = ref(
    window.matchMedia('(prefers-color-scheme: dark)').matches
  )

  const isDark = computed(() => {
    if (mode.value === 'system') return systemDark.value
    return mode.value === 'dark'
  })

  function applyTheme() {
    document.documentElement.setAttribute(
      'data-theme',
      isDark.value ? 'dark' : 'light'
    )
  }

  function setMode(m: ThemeMode) {
    mode.value = m
    localStorage.setItem('goma_theme', m)
  }

  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')

  function onSystemChange(e: MediaQueryListEvent) {
    systemDark.value = e.matches
  }

  mediaQuery.addEventListener('change', onSystemChange)

  onUnmounted(() => {
    mediaQuery.removeEventListener('change', onSystemChange)
  })

  watch(isDark, () => applyTheme(), { immediate: true })

  return { mode, isDark, setMode }
})
