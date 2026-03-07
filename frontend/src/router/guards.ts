import type { Router, RouteLocationNormalized } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTabsStore } from '@/stores/tabs'

export function setupRouterGuard(router: Router) {
  // Before each guard
  router.beforeEach(async (to: RouteLocationNormalized) => {
    const authStore = useAuthStore()
    const tabsStore = useTabsStore()

    console.log('[RouterGuard] Navigating to:', to.path)

    // Check if route requires authentication
    if (to.meta?.requiresAuth) {
      console.log('[RouterGuard] Route requires auth')
      
      if (!authStore.isLoggedIn) {
        console.log('[RouterGuard] Not logged in, checking token...')
        
        // Try to get user info if token exists
        if (authStore.token) {
          console.log('[RouterGuard] Token found, trying to get user info...')
          try {
            await authStore.getUserInfo()
            console.log('[RouterGuard] User info restored')
          } catch (error) {
            console.error('[RouterGuard] Failed to get user info:', error)
            return '/login'
          }
        } else {
          console.log('[RouterGuard] No token, redirecting to login')
          return '/login'
        }
      }

      // Check permission
      if (to.meta?.permission === 'admin' && !authStore.isAdmin) {
        console.log('[RouterGuard] Admin permission required, redirecting to 403')
        return '/403'
      }
    }

    // Initialize dashboard tab on first auth required route
    if (to.meta?.requiresAuth && tabsStore.tabs.length === 0) {
      tabsStore.initDashboardTab()
    }

    // Add tab
    tabsStore.addTab(to)

    console.log('[RouterGuard] Navigation allowed to:', to.path)
    return true
  })

  // After each hook
  router.afterEach((to: RouteLocationNormalized) => {
    // Set page title
    const title = to.meta?.title as string
    if (title) {
      document.title = `${title} - Dancer`
    }
  })
}
