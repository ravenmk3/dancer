<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTabsStore } from '@/stores/tabs'
import type { Tab } from '@/stores/tabs'

const route = useRoute()
const router = useRouter()
const tabsStore = useTabsStore()

const contextMenuVisible = ref(false)
const contextMenuPosition = ref({ x: 0, y: 0 })
const contextMenuTab = ref<Tab | null>(null)

const tabs = computed(() => tabsStore.tabs)
const activeTab = computed(() => tabsStore.activeTab)

const handleTabClick = (tab: Tab) => {
  router.push({ path: tab.path, query: tab.query })
}

const handleTabRemove = (id: string) => {
  tabsStore.removeTab(id)
}

const handleContextMenu = (e: MouseEvent, tab: Tab) => {
  e.preventDefault()
  contextMenuVisible.value = true
  contextMenuPosition.value = { x: e.clientX, y: e.clientY }
  contextMenuTab.value = tab
}

const closeContextMenu = () => {
  contextMenuVisible.value = false
  contextMenuTab.value = null
}

const closeCurrent = () => {
  if (contextMenuTab.value) {
    tabsStore.removeTab(contextMenuTab.value.id)
  }
  closeContextMenu()
}

const closeOthers = () => {
  if (contextMenuTab.value) {
    tabsStore.closeOthers(contextMenuTab.value.id)
  }
  closeContextMenu()
}

const closeAll = () => {
  tabsStore.closeAll()
  closeContextMenu()
}

const refreshCurrent = () => {
  if (contextMenuTab.value) {
    tabsStore.refreshTab(contextMenuTab.value.id)
  }
  closeContextMenu()
}

// Watch for route changes
watch(() => route.path, () => {
  tabsStore.addTab(route)
})

// Close context menu on click outside
window.addEventListener('click', closeContextMenu)
</script>

<template>
  <div class="tabs-bar">
    <el-scrollbar>
      <div class="tabs-wrapper">
        <div
          v-for="tab in tabs"
          :key="tab.id"
          :class="['tab-item', { active: activeTab === tab.id }]"
          @click="handleTabClick(tab)"
          @contextmenu.prevent="handleContextMenu($event, tab)"
        >
          <el-icon v-if="tab.icon">
            <component :is="tab.icon" />
          </el-icon>
          <span class="tab-title">{{ tab.title }}</span>
          <el-icon
            v-if="tab.closable"
            class="tab-close"
            @click.stop="handleTabRemove(tab.id)"
          >
            <Close />
          </el-icon>
        </div>
      </div>
    </el-scrollbar>

    <!-- Context Menu -->
    <div
      v-if="contextMenuVisible"
      class="context-menu"
      :style="{ left: contextMenuPosition.x + 'px', top: contextMenuPosition.y + 'px' }"
    >
      <div class="menu-item" @click="refreshCurrent">
        <el-icon><RefreshRight /></el-icon>
        <span>刷新</span>
      </div>
      <div class="menu-item" @click="closeCurrent">
        <el-icon><Close /></el-icon>
        <span>关闭</span>
      </div>
      <div class="menu-item" @click="closeOthers">
        <el-icon><CircleClose /></el-icon>
        <span>关闭其他</span>
      </div>
      <div class="menu-item" @click="closeAll">
        <el-icon><FolderDelete /></el-icon>
        <span>关闭全部</span>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.tabs-bar {
  height: 40px;
  background-color: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color-light);
  position: relative;
}

.tabs-wrapper {
  display: flex;
  align-items: center;
  height: 40px;
  padding: 0 10px;
}

.tab-item {
  display: flex;
  align-items: center;
  padding: 0 15px;
  height: 32px;
  margin-right: 5px;
  background-color: var(--el-fill-color-light);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  font-size: 13px;
  color: var(--el-text-color-regular);

  &:hover {
    background-color: var(--el-fill-color);
  }

  &.active {
    background-color: var(--el-color-primary);
    color: #fff;

    .tab-close {
      color: #fff;
    }
  }

  .el-icon {
    margin-right: 5px;
    font-size: 14px;
  }

  .tab-title {
    max-width: 120px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .tab-close {
    margin-left: 8px;
    margin-right: 0;
    font-size: 12px;
    width: 14px;
    height: 14px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s;

    &:hover {
      background-color: rgba(0, 0, 0, 0.1);
    }
  }
}

.context-menu {
  position: fixed;
  z-index: 3000;
  background-color: var(--el-bg-color);
  border: 1px solid var(--el-border-color-light);
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  padding: 5px 0;
  min-width: 120px;

  .menu-item {
    display: flex;
    align-items: center;
    padding: 8px 15px;
    cursor: pointer;
    font-size: 13px;
    color: var(--el-text-color-regular);
    transition: all 0.3s;

    &:hover {
      background-color: var(--el-fill-color-light);
      color: var(--el-color-primary);
    }

    .el-icon {
      margin-right: 8px;
      font-size: 14px;
    }
  }
}
</style>
