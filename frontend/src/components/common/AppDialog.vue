<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  visible: boolean
  title: string
  width?: string | number
  fullscreen?: boolean
  loading?: boolean
  showFooter?: boolean
  confirmText?: string
  cancelText?: string
  confirmType?: 'primary' | 'success' | 'warning' | 'danger'
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  width: '500px',
  fullscreen: false,
  loading: false,
  showFooter: true,
  confirmText: '确定',
  cancelText: '取消',
  confirmType: 'primary',
  disabled: false
})

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
  (e: 'close'): void
}>()

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const handleClose = () => {
  emit('update:visible', false)
  emit('close')
}

const handleConfirm = () => {
  emit('confirm')
}

const handleCancel = () => {
  emit('update:visible', false)
  emit('cancel')
}
</script>

<template>
  <el-dialog
    v-model="dialogVisible"
    :title="title"
    :width="width"
    :fullscreen="fullscreen"
    :close-on-click-modal="false"
    destroy-on-close
    @close="handleClose"
  >
    <slot />
    
    <template v-if="showFooter" #footer>
      <span class="dialog-footer">
        <el-button @click="handleCancel">{{ cancelText }}</el-button>
        <el-button
          :type="confirmType"
          :loading="loading"
          :disabled="disabled"
          @click="handleConfirm"
        >
          {{ confirmText }}
        </el-button>
      </span>
    </template>
  </el-dialog>
</template>

<style scoped lang="scss">
.app-dialog {
  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }
}
</style>
