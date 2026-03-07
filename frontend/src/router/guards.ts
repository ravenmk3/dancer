import type { Router, RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTabsStore } from '@/stores/tabs'

export function setupRouterGuard(router: Router) {
  // Before each guard
  router.beforeEach(async (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    const authStore = useAuthStore()
    const tabsStore = useTabsStore()

    console.log('[RouterGuard] Navigating to:', to.path)
    console.log('[RouterGuard] From:', from.path)
    console.log('[RouterGuard] isLoggedIn:', authStore.isLoggedIn)
    console.log('[RouterGuard] Token exists:', !!authStore.token)
    console.log('[RouterGuard] UserInfo exists:', !!authStore.userInfo)

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
            next('/login')
            return
          }
        } else {
          console.log('[RouterGuard] No token, redirecting to login')
          next('/login')
          return
        }
      }

      // Check permission
      if (to.meta?.permission === 'admin' && !authStore.isAdmin) {
        console.log('[RouterGuard] Admin permission required, redirecting to 403')
        next('/403')
        return
      }
    }

    // Initialize dashboard tab on first auth required route
    if (to.meta?.requiresAuth && tabsStore.tabs.length === 0) {
      tabsStore.initDashboardTab()
    }

    // Add tab
    tabsStore.addTab(to)

    console.log('[RouterGuard] Navigation allowed to:', to.path)
    next()
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
