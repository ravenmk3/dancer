<script setup lang="ts">
import { reactive, ref } from 'vue'

export interface SearchField {
  key: string
  label: string
  type?: 'input' | 'select' | 'date'
  placeholder?: string
  options?: { label: string; value: any }[]
  width?: string
}

interface Props {
  fields: SearchField[]
  showReset?: boolean
  showExpand?: boolean
  defaultExpanded?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  showReset: true,
  showExpand: false,
  defaultExpanded: false
})

const emit = defineEmits<{
  (e: 'search', params: Record<string, any>): void
  (e: 'reset'): void
}>()

const form = reactive<Record<string, any>>({})
const expanded = ref(props.defaultExpanded)
const formRef = ref()

const handleSearch = () => {
  const params: Record<string, any> = {}
  props.fields.forEach(field => {
    if (form[field.key] !== undefined && form[field.key] !== '') {
      params[field.key] = form[field.key]
    }
  })
  emit('search', params)
}

const handleReset = () => {
  props.fields.forEach(field => {
    form[field.key] = undefined
  })
  emit('reset')
}
</script>

<template>
  <el-form
    ref="formRef"
    :model="form"
    class="app-search"
    inline
  >
    <el-form-item
      v-for="field in fields"
      :key="field.key"
      :label="field.label"
    >
      <el-input
        v-if="field.type === 'input' || !field.type"
        v-model="form[field.key]"
        :placeholder="field.placeholder || `请输入${field.label}`"
        clearable
      />
      <el-select
        v-else-if="field.type === 'select'"
        v-model="form[field.key]"
        :placeholder="field.placeholder || `请选择${field.label}`"
        clearable
      >
        <el-option
          v-for="opt in field.options"
          :key="opt.value"
          :label="opt.label"
          :value="opt.value"
        />
      </el-select>
    </el-form-item>
    
    <el-form-item>
      <el-button type="primary" @click="handleSearch">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
      <el-button v-if="showReset" @click="handleReset">
        <el-icon><Refresh /></el-icon>
        重置
      </el-button>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="scss">
.app-search {
  margin-bottom: 20px;
  padding: 20px;
  background-color: var(--el-fill-color-light);
  border-radius: 4px;

  :deep(.el-form-item) {
    margin-bottom: 0;
  }
}
</style>
