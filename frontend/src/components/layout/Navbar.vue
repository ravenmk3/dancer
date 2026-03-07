<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useTabsStore } from '@/stores/tabs'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()
const tabsStore = useTabsStore()

const isCollapse = computed(() => appStore.sidebarCollapsed)
const userInfo = computed(() => authStore.userInfo)
const isDark = computed(() => appStore.isDark)

const toggleSidebar = () => {
  appStore.toggleSidebar()
}

const toggleTheme = () => {
  const newTheme = isDark.value ? 'light' : 'dark'
  appStore.setTheme(newTheme)
}

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'logout':
      authStore.logout()
      tabsStore.closeAll()
      router.push('/login')
      break
  }
}

const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
  } else {
    document.exitFullscreen()
  }
}
</script>

<template>
  <div class="navbar">
    <div class="left">
      <div class="collapse-btn" @click="toggleSidebar">
        <el-icon :size="20">
          <Fold v-if="!isCollapse" />
          <Expand v-else />
        </el-icon>
      </div>
      
      <breadcrumb />
    </div>
    
    <div class="right">
      <!-- Theme Toggle -->
      <div class="nav-item" @click="toggleTheme">
        <el-icon :size="18">
          <Sunny v-if="isDark" />
          <Moon v-else />
        </el-icon>
      </div>
      
      <!-- Fullscreen -->
      <div class="nav-item" @click="toggleFullscreen">
        <el-icon :size="18">
          <FullScreen />
        </el-icon>
      </div>
      
      <!-- User Dropdown -->
      <el-dropdown @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32" :icon="UserFilled" />
          <span class="username">{{ userInfo?.username }}</span>
          <el-icon><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<style scoped lang="scss">
.navbar {
  height: 60px;
  background-color: var(--navbar-bg);
  border-bottom: 1px solid var(--el-border-color-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;

  .left {
    display: flex;
    align-items: center;

    .collapse-btn {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      border-radius: 4px;
      transition: all 0.3s;
      color: var(--el-text-color-regular);

      &:hover {
        background-color: var(--el-fill-color-light);
        color: var(--el-color-primary);
      }
    }
  }

  .right {
    display: flex;
    align-items: center;
    gap: 10px;

    .nav-item {
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      border-radius: 4px;
      transition: all 0.3s;
      color: var(--el-text-color-regular);

      &:hover {
        background-color: var(--el-fill-color-light);
        color: var(--el-color-primary);
      }
    }

    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 5px 10px;
      cursor: pointer;
      border-radius: 4px;
      transition: all 0.3s;

      &:hover {
        background-color: var(--el-fill-color-light);
      }

      .username {
        font-size: 14px;
        color: var(--el-text-color-primary);
      }

      .el-icon {
        color: var(--el-text-color-secondary);
      }
    }
  }
}
</style>
