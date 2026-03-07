<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTabsStore } from '@/stores/tabs'

const route = useRoute()
const tabsStore = useTabsStore()

const cachedViews = computed(() => {
  return tabsStore.tabs.filter(tab => tab.cached).map(tab => tab.path)
})
</script>

<template>
  <div class="app-main">
    <router-view v-slot="{ Component }">
      <keep-alive :include="cachedViews">
        <component :is="Component" :key="route.fullPath" />
      </keep-alive>
    </router-view>
  </div>
</template>

<style scoped lang="scss">
.app-main {
  flex: 1;
  padding: 20px;
  overflow: auto;
  background-color: var(--app-bg);
}
</style>
