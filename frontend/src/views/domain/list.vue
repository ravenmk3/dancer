<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppTable from '@/components/common/AppTable.vue'
import AppDialog from '@/components/common/AppDialog.vue'
import AppEmpty from '@/components/common/AppEmpty.vue'
import { getDomainListApi, createDomainApi, updateDomainApi, deleteDomainApi } from '@/api/domain'
import { getZoneListApi } from '@/api/zone'
import type { Domain, CreateDomainForm } from '@/types/domain'
import type { Zone } from '@/types/zone'
import type { Column } from '@/components/common/AppTable.vue'

// Data
const domains = ref<Domain[]>([])
const zones = ref<Zone[]>([])
const selectedZone = ref('')
const loading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const isEdit = ref(false)

// Form
const form = ref<CreateDomainForm>({
  zone: '',
  domain: '',
  ips: [''],
  ttl: 300
})

const formRules = {
  zone: [{ required: true, message: '请选择 Zone', trigger: 'change' }],
  domain: [
    { required: true, message: '请输入 Domain 名称', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9*@]([a-zA-Z0-9-]*[a-zA-Z0-9])?$/, message: '请输入有效的域名格式', trigger: 'blur' }
  ],
  ips: [
    { 
      validator: (rule: any, value: string[], callback: Function) => {
        if (!value || value.length === 0 || value.every(ip => !ip)) {
          callback(new Error('至少需要一个 IP 地址'))
        } else {
          const ipRegex = /^(\d{1,3}\.){3}\d{1,3}$/
          const invalidIps = value.filter(ip => ip && !ipRegex.test(ip))
          if (invalidIps.length > 0) {
            callback(new Error(`无效的 IP 地址: ${invalidIps.join(', ')}`))
          } else {
            callback()
          }
        }
      },
      trigger: 'blur'
    }
  ],
  ttl: [
    { required: true, message: '请输入 TTL', trigger: 'blur' },
    { type: 'number', min: 1, message: 'TTL 最小值为 1', trigger: 'blur' }
  ]
}

const formRef = ref()

// Table columns
const columns: Column[] = [
  { prop: 'domain', label: 'Domain', minWidth: '150' },
  { prop: 'name', label: '完整域名', minWidth: '200' },
  { 
    prop: 'ips', 
    label: 'IP 地址', 
    minWidth: '200',
    formatter: (row: any) => row.ips.join(', ')
  },
  { prop: 'ttl', label: 'TTL', width: '100' },
  { prop: 'record_count', label: '记录数', width: '80' },
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
  try {
    const res = await getZoneListApi()
    zones.value = res.zones || []
    if (zones.value.length > 0 && !selectedZone.value) {
      selectedZone.value = zones.value[0].zone
      fetchDomains()
    }
  } catch (error) {
    console.error(error)
  }
}

const fetchDomains = async () => {
  if (!selectedZone.value) return
  
  loading.value = true
  try {
    const res = await getDomainListApi(selectedZone.value)
    domains.value = res.domains || []
  } finally {
    loading.value = false
  }
}

const handleZoneChange = () => {
  fetchDomains()
}

const handleAdd = () => {
  if (!selectedZone.value) {
    ElMessage.warning('请先选择一个 Zone')
    return
  }
  isEdit.value = false
  form.value = {
    zone: selectedZone.value,
    domain: '',
    ips: [''],
    ttl: 300
  }
  dialogVisible.value = true
}

const handleEdit = (row: Domain) => {
  isEdit.value = true
  form.value = {
    zone: row.zone,
    domain: row.domain,
    ips: [...row.ips],
    ttl: row.ttl
  }
  dialogVisible.value = true
}

const handleDelete = async (row: Domain) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除 Domain "${row.name}" 吗？此操作不可恢复。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteDomainApi(row.zone, row.domain)
    ElMessage.success('删除成功')
    fetchDomains()
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
    
    // Filter out empty IPs
    const submitData = {
      ...form.value,
      ips: form.value.ips.filter(ip => ip.trim() !== '')
    }
    
    if (submitData.ips.length === 0) {
      ElMessage.error('至少需要一个 IP 地址')
      return
    }
    
    dialogLoading.value = true
    try {
      if (isEdit.value) {
        await updateDomainApi(submitData)
        ElMessage.success('更新成功')
      } else {
        await createDomainApi(submitData)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchDomains()
    } finally {
      dialogLoading.value = false
    }
  })
}

// IP management
const addIp = () => {
  form.value.ips.push('')
}

const removeIp = (index: number) => {
  if (form.value.ips.length > 1) {
    form.value.ips.splice(index, 1)
  }
}

onMounted(() => {
  fetchZones()
})
</script>

<template>
  <div class="domain-list">
    <div class="page-header">
      <div class="header-left">
        <h2>Domain 管理</h2>
        <el-select v-model="selectedZone" placeholder="选择 Zone" @change="handleZoneChange">
          <el-option
            v-for="zone in zones"
            :key="zone.zone"
            :label="zone.zone"
            :value="zone.zone"
          />
        </el-select>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新建 Domain
      </el-button>
    </div>

    <AppTable
      :data="domains"
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

    <AppEmpty v-if="!loading && domains.length === 0" description="暂无 Domain 数据" />

    <!-- Add/Edit Dialog -->
    <AppDialog
      v-model:visible="dialogVisible"
      :title="isEdit ? '编辑 Domain' : '新建 Domain'"
      :loading="dialogLoading"
      width="600px"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="所属 Zone" prop="zone">
          <el-select v-model="form.zone" placeholder="选择 Zone" style="width: 100%">
            <el-option
              v-for="zone in zones"
              :key="zone.zone"
              :label="zone.zone"
              :value="zone.zone"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="Domain" prop="domain">
          <el-input
            v-model="form.domain"
            placeholder="请输入 Domain 名称，如：www 或 @（根域名）"
          />
        </el-form-item>
        
        <el-form-item label="IP 地址" prop="ips">
          <div class="ip-list">
            <div
              v-for="(ip, index) in form.ips"
              :key="index"
              class="ip-item"
            >
              <el-input
                v-model="form.ips[index]"
                placeholder="如：192.168.1.1"
              />
              <el-button
                v-if="form.ips.length > 1"
                type="danger"
                link
                @click="removeIp(index)"
              >
                <el-icon><Delete /></el-icon>
              </el-button>
            </div>
            <el-button type="primary" link @click="addIp">
              <el-icon><Plus /></el-icon>
              添加 IP
            </el-button>
          </div>
        </el-form-item>
        
        <el-form-item label="TTL" prop="ttl">
          <el-input-number
            v-model="form.ttl"
            :min="1"
            :max="86400"
            style="width: 200px"
          />
          <span class="unit">秒</span>
        </el-form-item>
      </el-form>
    </AppDialog>
  </div>
</template>

<style scoped lang="scss">
.domain-list {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;

    .header-left {
      display: flex;
      align-items: center;
      gap: 20px;

      h2 {
        margin: 0;
        font-size: 20px;
        font-weight: 600;
      }

      .el-select {
        width: 200px;
      }
    }
  }

  .ip-list {
    .ip-item {
      display: flex;
      gap: 10px;
      margin-bottom: 10px;

      .el-input {
        flex: 1;
      }
    }
  }

  .unit {
    margin-left: 10px;
    color: var(--el-text-color-secondary);
  }
}
</style>
