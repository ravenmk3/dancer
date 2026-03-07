import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { RouteLocationNormalized } from 'vue-router'

export interface Tab {
  id: string
  title: string
  path: string
  query: Record<string, any>
  closable: boolean
  icon?: string
  cached: boolean
}

export const useTabsStore = defineStore('tabs', () => {
  // State
  const tabs = ref<Tab[]>([])
  const activeTab = ref<string>('')

  // Getters
  const activeTabInfo = computed(() => {
    return tabs.value.find(tab => tab.id === activeTab.value)
  })

  // Actions
  const addTab = (route: RouteLocationNormalized) => {
    const { path, query, meta } = route
    const id = path + JSON.stringify(query)
    
    // Don't add login page or hidden routes
    if (path === '/login' || meta?.hidden) {
      return
    }

    // Check if tab already exists
    const existingTab = tabs.value.find(tab => tab.id === id)
    if (existingTab) {
      activeTab.value = id
      return
    }

    const tab: Tab = {
      id,
      title: (meta?.title as string) || '未命名',
      path,
      query: { ...query },
      closable: path !== '/dashboard',
      icon: meta?.icon as string,
      cached: meta?.keepAlive !== false
    }

    tabs.value.push(tab)
    activeTab.value = id
  }

  const removeTab = (id: string) => {
    const index = tabs.value.findIndex(tab => tab.id === id)
    if (index === -1) return

    const isActive = activeTab.value === id
    tabs.value.splice(index, 1)

    // Activate adjacent tab if removing active tab
    if (isActive && tabs.value.length > 0) {
      const newIndex = Math.min(index, tabs.value.length - 1)
      activeTab.value = tabs.value[newIndex].id
    }
  }

  const closeOthers = (id: string) => {
    const current = tabs.value.find(tab => tab.id === id)
    if (!current) return
    tabs.value = [current]
    activeTab.value = id
  }

  const closeAll = () => {
    tabs.value = []
    activeTab.value = ''
  }

  const setActiveTab = (id: string) => {
    activeTab.value = id
  }

  const refreshTab = (id: string) => {
    // This will be handled by the component to force re-render
    const tab = tabs.value.find(t => t.id === id)
    if (tab) {
      tab.cached = false
      setTimeout(() => {
        tab.cached = true
      }, 0)
    }
  }

  return {
    tabs,
    activeTab,
    activeTabInfo,
    addTab,
    removeTab,
    closeOthers,
    closeAll,
    setActiveTab,
    refreshTab
  }
})
