<script setup lang="ts">
interface Props {
  description?: string
  image?: string
}

const props = withDefaults(defineProps<Props>(), {
  description: '暂无数据',
  image: undefined
})
</script>

<template>
  <el-empty :description="description" :image="image">
    <slot />
  </el-empty>
</template>

<style scoped lang="scss">
.app-empty {
  padding: 40px 0;
}
</style>
