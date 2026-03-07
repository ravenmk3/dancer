<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { RouteLocationMatched } from 'vue-router'
import { constantRoutes } from '@/router/routes'
import type { MenuItem } from '@/types/route'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const isCollapse = computed(() => appStore.sidebarCollapsed)

// Generate menu from routes
const menus = computed<MenuItem[]>(() => {
  const generateMenus = (routes: typeof constantRoutes, parentPath = ''): MenuItem[] => {
    const result: MenuItem[] = []
    
    for (const item of routes) {
      if (item.meta?.hidden) continue
      
      // Check permission
      if (item.meta?.permission === 'admin' && !authStore.isAdmin) continue
      
      // Build full path by combining parent path with current path
      const fullPath = item.path.startsWith('/') ? item.path : `${parentPath}/${item.path}`
      
      const menu: MenuItem = {
        path: item.redirect as string || fullPath,
        title: item.meta?.title as string || '',
        icon: item.meta?.icon as string,
        permission: item.meta?.permission as string
      }
      
      if (item.children && item.children.length > 0) {
        const childMenus = generateMenus(item.children as typeof constantRoutes, fullPath)
        if (childMenus.length > 0) {
          menu.children = childMenus
        }
      }
      
      result.push(menu)
    }
    
    return result
  }
  
  return generateMenus(constantRoutes.filter(r => r.path !== '/login' && r.path !== '/profile'))
})

const activeMenu = computed(() => {
  const { path } = route
  return path
})

const handleMenuClick = (menu: MenuItem) => {
  router.push(menu.path)
}

// Check if menu is active
const isMenuActive = (menu: MenuItem): boolean => {
  return route.path === menu.path || route.path.startsWith(menu.path + '/')
}
</script>

<template>
  <div class="sidebar" :class="{ collapsed: isCollapse }">
    <div class="logo">
      <div class="logo-icon">D</div>
      <span v-if="!isCollapse" class="logo-text">Dancer DNS</span>
    </div>
    
    <el-scrollbar class="menu-scrollbar">
      <el-menu
        :default-active="activeMenu"
        :collapse="isCollapse"
        :collapse-transition="false"
        background-color="transparent"
        text-color="var(--sidebar-text)"
        active-text-color="var(--el-color-primary)"
      >
        <template v-for="menu in menus" :key="menu.path">
          <!-- Submenu -->
          <el-sub-menu v-if="menu.children && menu.children.length" :index="menu.path">
            <template #title>
              <el-icon v-if="menu.icon">
                <component :is="menu.icon" />
              </el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <el-menu-item
              v-for="child in menu.children"
              :key="child.path"
              :index="child.path"
              @click="handleMenuClick(child)"
            >
              <el-icon v-if="child.icon">
                <component :is="child.icon" />
              </el-icon>
              <span>{{ child.title }}</span>
            </el-menu-item>
          </el-sub-menu>
          
          <!-- Menu Item -->
          <el-menu-item v-else :index="menu.path" @click="handleMenuClick(menu)">
            <el-icon v-if="menu.icon">
              <component :is="menu.icon" />
            </el-icon>
            <span>{{ menu.title }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>
  </div>
</template>

<style scoped lang="scss">
.sidebar {
  width: 210px;
  height: 100%;
  background-color: var(--sidebar-bg);
  border-right: 1px solid var(--el-border-color-light);
  transition: width 0.3s;
  display: flex;
  flex-direction: column;

  &.collapsed {
    width: 64px;
  }
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 15px;
  border-bottom: 1px solid var(--el-border-color-light);

  .logo-icon {
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 20px;
    font-weight: bold;
    flex-shrink: 0;
  }

  .logo-text {
    margin-left: 10px;
    font-size: 18px;
    font-weight: 600;
    color: var(--sidebar-text);
    white-space: nowrap;
    overflow: hidden;
  }
}

.menu-scrollbar {
  flex: 1;
  overflow: hidden;

  :deep(.el-menu) {
    border-right: none;
    background-color: transparent;

    .el-menu-item,
    .el-sub-menu__title {
      height: 50px;
      line-height: 50px;

      &:hover {
        background-color: var(--el-fill-color-light);
      }

      &.is-active {
        background-color: var(--el-color-primary-light-9);
      }
    }

    .el-sub-menu {
      .el-menu {
        background-color: transparent;
      }

      .el-menu-item {
        height: 40px;
        line-height: 40px;
      }
    }
  }
}
</style>
