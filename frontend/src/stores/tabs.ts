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
  fixed?: boolean
}

// Fixed tabs that cannot be closed
const FIXED_TABS = ['/dashboard']

export const useTabsStore = defineStore('tabs', () => {
  // State
  const tabs = ref<Tab[]>([])
  const activeTab = ref<string>('')

  // Getters
  const activeTabInfo = computed(() => {
    return tabs.value.find(tab => tab.id === activeTab.value)
  })

  // Initialize with dashboard tab
  const initDashboardTab = () => {
    const dashboardTab: Tab = {
      id: '/dashboard' + JSON.stringify({}),
      title: '仪表盘',
      path: '/dashboard',
      query: {},
      closable: false,
      icon: 'DataLine',
      cached: true,
      fixed: true
    }
    
    const exists = tabs.value.find(tab => tab.id === dashboardTab.id)
    if (!exists) {
      tabs.value.unshift(dashboardTab)
    }
  }

  // Actions
  const addTab = (route: RouteLocationNormalized) => {
    const { path, query, meta } = route
    const id = path + JSON.stringify(query)
    
    // Don't add login page
    if (path === '/login') {
      return
    }

    // Ensure dashboard exists
    initDashboardTab()

    // Check if tab already exists
    const existingTab = tabs.value.find(tab => tab.id === id)
    if (existingTab) {
      activeTab.value = id
      return
    }

    const isFixed = FIXED_TABS.includes(path)
    const tab: Tab = {
      id,
      title: (meta?.title as string) || '未命名',
      path,
      query: { ...query },
      closable: !isFixed,
      icon: meta?.icon as string,
      cached: meta?.keepAlive !== false,
      fixed: isFixed
    }

    tabs.value.push(tab)
    activeTab.value = id
  }

  const removeTab = (id: string) => {
    const tab = tabs.value.find(tab => tab.id === id)
    // Cannot remove fixed tabs
    if (!tab || tab.fixed || !tab.closable) return

    const index = tabs.value.findIndex(tab => tab.id === id)
    if (index === -1) return

    const isActive = activeTab.value === id
    tabs.value.splice(index, 1)

    // Activate adjacent tab if removing active tab
    if (isActive && tabs.value.length > 0) {
      const newIndex = Math.min(index, tabs.value.length - 1)
      const nextTab = tabs.value[newIndex]
      if (nextTab) {
        activeTab.value = nextTab.id
      }
    }
  }

  const closeOthers = (id: string) => {
    const current = tabs.value.find(tab => tab.id === id)
    if (!current) return
    
    // Keep fixed tabs and current tab
    tabs.value = tabs.value.filter(tab => tab.fixed || tab.id === id)
    activeTab.value = id
  }

  const closeAll = () => {
    // Only keep fixed tabs
    tabs.value = tabs.value.filter(tab => tab.fixed)
    if (tabs.value.length > 0) {
      const firstTab = tabs.value[0]
      if (firstTab) {
        activeTab.value = firstTab.id
      }
    }
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
    refreshTab,
    initDashboardTab
  }
})
