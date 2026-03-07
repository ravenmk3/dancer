import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export type Theme = 'light' | 'dark' | 'auto'

export const useAppStore = defineStore('app', () => {
  // State
  const sidebarCollapsed = ref(false)
  const theme = ref<Theme>((localStorage.getItem('theme') as Theme) || 'auto')
  const currentTheme = ref<'light' | 'dark'>('light')

  // Getters
  const isDark = computed(() => currentTheme.value === 'dark')

  // Actions
  const toggleSidebar = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  const setTheme = (newTheme: Theme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme()
  }

  const applyTheme = () => {
    let resolvedTheme: 'light' | 'dark'
    
    if (theme.value === 'auto') {
      resolvedTheme = window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
    } else {
      resolvedTheme = theme.value
    }
    
    currentTheme.value = resolvedTheme
    document.documentElement.setAttribute('data-theme', resolvedTheme)
  }

  const initTheme = () => {
    applyTheme()
    // Listen for system theme changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
      if (theme.value === 'auto') {
        applyTheme()
      }
    })
  }

  return {
    sidebarCollapsed,
    theme,
    currentTheme,
    isDark,
    toggleSidebar,
    setTheme,
    applyTheme,
    initTheme
  }
})
