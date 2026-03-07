<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    title: item.meta.title as string,
    path: item.path
  }))
})
</script>

<template>
  <el-breadcrumb separator="/" class="breadcrumb">
    <el-breadcrumb-item
      v-for="(item, index) in breadcrumbs"
      :key="item.path"
    >
      <span
        :class="{ 'is-link': index < breadcrumbs.length - 1 }"
        @click="index < breadcrumbs.length - 1 && $router.push(item.path)"
      >
        {{ item.title }}
      </span>
    </el-breadcrumb-item>
  </el-breadcrumb>
</template>

<style scoped lang="scss">
.breadcrumb {
  margin-left: 15px;

  :deep(.el-breadcrumb__item) {
    .is-link {
      cursor: pointer;
      color: var(--el-text-color-regular);

      &:hover {
        color: var(--el-color-primary);
      }
    }

    &:last-child {
      .el-breadcrumb__inner {
        color: var(--el-text-color-secondary);
      }
    }
  }
}
</style>
