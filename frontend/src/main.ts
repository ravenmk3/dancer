import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

import App from './App.vue'
import router from './router'
import { useAuthStore } from '@/stores/auth'

// Import global styles
import './styles/index.scss'

const app = createApp(App)

// Register Element Plus
app.use(ElementPlus)

// Register all Element Plus icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Use Pinia and Router
app.use(createPinia())
app.use(router)

// Initialize auth state before mounting
const initApp = async () => {
  console.log('[App] Initializing...')
  
  const authStore = useAuthStore()
  
  // Try to restore session if token exists
  await authStore.initAuth()
  
  console.log('[App] Auth initialized, isLoggedIn:', authStore.isLoggedIn)
  
  // Mount app
  app.mount('#app')
  console.log('[App] Mounted successfully')
}

initApp().catch(error => {
  console.error('[App] Failed to initialize:', error)
  app.mount('#app')
})
