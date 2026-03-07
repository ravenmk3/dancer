import type { Router, RouteLocationNormalized, NavigationGuardNext } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTabsStore } from '@/stores/tabs'

export function setupRouterGuard(router: Router) {
  // Before each guard
  router.beforeEach(async (to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    const authStore = useAuthStore()
    const tabsStore = useTabsStore()

    // Check if route requires authentication
    if (to.meta?.requiresAuth) {
      if (!authStore.isLoggedIn) {
        // Try to get user info if token exists
        if (authStore.token) {
          try {
            await authStore.getUserInfo()
          } catch {
            next('/login')
            return
          }
        } else {
          next('/login')
          return
        }
      }

      // Check permission
      if (to.meta?.permission === 'admin' && !authStore.isAdmin) {
        next('/403')
        return
      }
    }

    // Add tab
    tabsStore.addTab(to)

    next()
  })

  // After each hook
  router.afterEach((to: RouteLocationNormalized) => {
    // Set page title
    const title = to.meta?.title as string
    if (title) {
      document.title = `${title} - Dancer DNS`
    }
  })
}
