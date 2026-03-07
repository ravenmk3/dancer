<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppTable from '@/components/common/AppTable.vue'
import AppDialog from '@/components/common/AppDialog.vue'
import AppEmpty from '@/components/common/AppEmpty.vue'
import { getUserListApi, createUserApi, updateUserApi, deleteUserApi } from '@/api/user'
import type { User, CreateUserForm, UpdateUserForm, UserType } from '@/types/user'
import type { Column } from '@/components/common/AppTable.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

// Data
const users = ref<User[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const isEdit = ref(false)
const currentUserId = ref('')

// Form
const form = ref<CreateUserForm>({
  username: '',
  password: '',
  user_type: 'normal'
})

const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '长度在 3 到 32 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为 6 个字符', trigger: 'blur' }
  ],
  user_type: [{ required: true, message: '请选择用户类型', trigger: 'change' }]
}

const formRef = ref()

// Table columns
const columns: Column[] = [
  { prop: 'username', label: '用户名', minWidth: '150' },
  { 
    prop: 'user_type', 
    label: '用户类型', 
    width: '120',
    formatter: (row: User) => row.user_type === 'admin' ? '管理员' : '普通用户'
  },
  { 
    prop: 'created_at', 
    label: '创建时间', 
    width: '180',
    formatter: (row: User) => new Date(row.created_at * 1000).toLocaleString()
  },
  { 
    prop: 'updated_at', 
    label: '更新时间', 
    width: '180',
    formatter: (row: User) => new Date(row.updated_at * 1000).toLocaleString()
  },
  { slot: 'actions', label: '操作', width: '150', fixed: 'right' }
]

// Methods
const fetchUsers = async () => {
  loading.value = true
  try {
    const res = await getUserListApi()
    users.value = res.users || []
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  isEdit.value = false
  currentUserId.value = ''
  form.value = {
    username: '',
    password: '',
    user_type: 'normal'
  }
  dialogVisible.value = true
}

const handleEdit = (row: User) => {
  isEdit.value = true
  currentUserId.value = row.id
  form.value = {
    username: row.username,
    password: '', // Password is optional for update
    user_type: row.user_type
  }
  dialogVisible.value = true
}

const handleDelete = async (row: User) => {
  // Cannot delete yourself
  if (row.id === authStore.userInfo?.id) {
    ElMessage.error('不能删除当前登录用户')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${row.username}" 吗？此操作不可恢复。`,
      '警告',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await deleteUserApi(row.id)
    ElMessage.success('删除成功')
    fetchUsers()
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
        const updateData: UpdateUserForm = {
          id: currentUserId.value,
          username: form.value.username,
          user_type: form.value.user_type
        }
        if (form.value.password) {
          updateData.password = form.value.password
        }
        await updateUserApi(updateData)
        ElMessage.success('更新成功')
      } else {
        await createUserApi(form.value)
        ElMessage.success('创建成功')
      }
      dialogVisible.value = false
      fetchUsers()
    } finally {
      dialogLoading.value = false
    }
  })
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="user-list">
    <div class="page-header">
      <h2>用户管理</h2>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新建用户
      </el-button>
    </div>

    <AppTable
      :data="users"
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

    <AppEmpty v-if="!loading && users.length === 0" description="暂无用户数据" />

    <!-- Add/Edit Dialog -->
    <AppDialog
      v-model:visible="dialogVisible"
      :title="isEdit ? '编辑用户' : '新建用户'"
      :loading="dialogLoading"
      @confirm="handleSubmit"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="formRules"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input
            v-model="form.username"
            placeholder="请输入用户名"
            :disabled="isEdit"
          />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            show-password
          />
          <div v-if="isEdit" class="form-tip">留空表示不修改密码</div>
        </el-form-item>
        
        <el-form-item label="用户类型" prop="user_type">
          <el-radio-group v-model="form.user_type">
            <el-radio value="normal">普通用户</el-radio>
            <el-radio value="admin">管理员</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </AppDialog>
  </div>
</template>

<style scoped lang="scss">
.user-list {
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

  .form-tip {
    font-size: 12px;
    color: var(--el-text-color-secondary);
    margin-top: 5px;
  }
}
</style>
