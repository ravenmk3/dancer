<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppTable from '@/components/common/AppTable.vue'
import AppDialog from '@/components/common/AppDialog.vue'
import AppSearch from '@/components/common/AppSearch.vue'
import AppEmpty from '@/components/common/AppEmpty.vue'
import { getZoneListApi, createZoneApi, deleteZoneApi } from '@/api/zone'
import type { Zone, CreateZoneForm } from '@/types/zone'
import type { Column } from '@/components/common/AppTable.vue'

// Data
const zones = ref<Zone[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const isEdit = ref(false)

// Form
const form = ref<CreateZoneForm>({
  zone: ''
})

const formRules = {
  zone: [
    { required: true, message: '请输入 Zone 名称', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9][a-zA-Z0-9-]*(\.[a-zA-Z0-9][a-zA-Z0-9-]*)+$/, message: '请输入有效的域名格式', trigger: 'blur' }
  ]
}

const formRef = ref()

// Table columns
const columns: Column[] = [
  { prop: 'zone', label: 'Zone 名称', minWidth: '200' },
  { prop: 'record_count', label: 'Domain 数量', width: '120' },
  { 
    prop: 'created_at', 
    label: '创建时间', 
    width: '180',
    formatter: (row: any) => new Date(row.created_at * 1000).toLocaleString()
  },
  { 
    prop: 'updated_at', 
    label: '更新时间', 
    width: '180',
    formatter: (row: any) => new Date(row.updated_at * 1000).toLocaleString()
  },
  { slot: 'actions', label: '操作', width: '150', fixed: 'right' }
]

// Methods
const fetchZones = async () => {
  loading.value = true
  try {
    const res = await getZoneListApi()
    zones.value = res.zones || []
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  form.value = { zone: '' }
  dialogVisible.value = true
}

const handleEdit = (row: Zone) => {
  isEdit.value = true
  form.value = { zone: row.zone }
  dialogVisible.value = true
}

const handleDelete = async (row: Zone) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Zone "${row.zone}" 吗？此操作将同时删除该 Zone 下的所有 Domain 记录，且不可恢复。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteZoneApi(row.zone)
    ElMessage.success('删除成功')
    fetchZones()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    
    dialogLoading.value = true
    try {
      if (isEdit.value) {
        // Update API not implemented in backend, use create for now
        ElMessage.info('编辑功能暂未实现')
      } else {
        await createZoneApi(form.value)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchZones()
    } finally {
      dialogLoading.value = false
    }
  })
}

onMounted(() => {
  fetchZones()
})
</script>

<template>
  <div class="zone-list">
    <div class="page-header">
      <h2>Zone 管理</h2>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新建 Zone
      </el-button>
    </div>

    <AppTable
      :data="zones"
      :columns="columns"
      :loading="loading"
    >
      <template #actions="{ row }">
        <el-button
          type="primary"
          link
          size="small"
          @click="handleEdit(row)"
        >
          编辑
        </el-button>
        <el-button
          type="danger"
          link
          size="small"
          @click="handleDelete(row)"
        >
          删除
        </el-button>
      </template>
    </AppTable>

    <AppEmpty v-if="!loading && zones.length === 0" description="暂无 Zone 数据" />

    <!-- Add/Edit Dialog -->
    <AppDialog
      v-model:visible="dialogVisible"
      :title="isEdit ? '编辑 Zone' : '新建 Zone'"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="Zone 名称" prop="zone">
          <el-input
            v-model="form.zone"
            placeholder="请输入 Zone 名称，如：example.com"
          />
        </el-form-item>
      </el-form>
    </AppDialog>
  </div>
</template>

<style scoped lang="scss">
.zone-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    h2 {
      margin: 0;
      font-size: 20px;
      font-weight: 600;
    }
  }
}
</style>
