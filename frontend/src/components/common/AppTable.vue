<script setup lang="ts">
import type { TableColumnCtx } from 'element-plus'

export interface Column {
  prop?: string
  label: string
  width?: string | number
  minWidth?: string | number
  fixed?: boolean | 'left' | 'right'
  sortable?: boolean | 'custom'
  formatter?: (row: any, column: TableColumnCtx<any>, cellValue: any) => string
  slot?: string
  type?: 'selection' | 'index' | 'expand'
}

interface Props {
  data: any[]
  columns: Column[]
  loading?: boolean
  pagination?: {
    page: number
    pageSize: number
    total: number
    pageSizes?: number[]
  }
  selection?: boolean
  rowKey?: string
  showIndex?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  selection: false,
  rowKey: 'id',
  showIndex: false
})

const emit = defineEmits<{
  (e: 'page-change', page: number): void
  (e: 'size-change', size: number): void
  (e: 'selection-change', selection: any[]): void
}>()

const handleCurrentChange = (page: number) => {
  emit('page-change', page)
}

const handleSizeChange = (size: number) => {
  emit('size-change', size)
}

const handleSelectionChange = (selection: any[]) => {
  emit('selection-change', selection)
}
</script>

<template>
  <div class="app-table">
    <el-table
      v-loading="loading"
      :data="data"
      :row-key="rowKey"
      @selection-change="handleSelectionChange"
    >
      <!-- Selection column -->
      <el-table-column v-if="selection" type="selection" width="55" />
      
      <!-- Index column -->
      <el-table-column v-if="showIndex" type="index" width="50" label="#" />
      
      <!-- Data columns -->
      <el-table-column
        v-for="column in columns"
        :key="column.prop || column.type || column.slot"
        :prop="column.prop"
        :label="column.label"
        :width="column.width"
        :min-width="column.minWidth"
        :fixed="column.fixed"
        :sortable="column.sortable"
        :formatter="column.formatter"
        :type="column.type"
      >
        <template v-if="column.slot" #default="scope">
          <slot :name="column.slot" v-bind="scope" />
        </template>
      </el-table-column>
    </el-table>
    
    <!-- Pagination -->
    <div v-if="pagination" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="pagination.pageSizes || [10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.app-table {
  .pagination-wrapper {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
