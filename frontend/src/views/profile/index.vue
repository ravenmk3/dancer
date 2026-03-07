<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import type { User, ChangePasswordForm } from '@/types/user'

const router = useRouter()
const authStore = useAuthStore()

const userInfo = ref<User | null>(null)
const loading = ref(false)
const passwordForm = ref<ChangePasswordForm>({
  old_password: '',
  new_password: ''
})
const passwordFormRef = ref()

const passwordRules = {
  old_password: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  new_password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为 6 个字符', trigger: 'blur' }
  ]
}

onMounted(() => {
  userInfo.value = authStore.userInfo
})

const handleChangePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    
    loading.value = true
    try {
      await authStore.changePassword(passwordForm.value)
      ElMessage.success('密码修改成功')
      passwordForm.value = { old_password: '', new_password: '' }
    } finally {
      loading.value = false
    }
  })
}

const formatTime = (timestamp: number) => {
  return new Date(timestamp * 1000).toLocaleString()
}
</script>

<template>
  <div class="profile-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>个人信息</span>
            </div>
          </template>
          
          <div class="user-info" v-if="userInfo">
            <div class="info-item">
              <label>用户名：</label>
              <span>{{ userInfo.username }}</span>
            </div>
            <div class="info-item">
              <label>用户类型：</label>
              <el-tag :type="userInfo.user_type === 'admin' ? 'danger' : 'info'">
                {{ userInfo.user_type === 'admin' ? '管理员' : '普通用户' }}
              </el-tag>
            </div>
            <div class="info-item">
              <label>创建时间：</label>
              <span>{{ formatTime(userInfo.created_at) }}</span>
            </div>
            <div class="info-item">
              <label>更新时间：</label>
              <span>{{ formatTime(userInfo.updated_at) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>修改密码</span>
            </div>
          </template>
          
          <el-form
            ref="passwordFormRef"
            :model="passwordForm"
            :rules="passwordRules"
            label-width="100px"
          >
            <el-form-item label="旧密码" prop="old_password">
              <el-input
                v-model="passwordForm.old_password"
                type="password"
                show-password
                placeholder="请输入旧密码"
              />
            </el-form-item>
            
            <el-form-item label="新密码" prop="new_password">
              <el-input
                v-model="passwordForm.new_password"
                type="password"
                show-password
                placeholder="请输入新密码"
              />
            </el-form-item>
            
            <el-form-item>
              <el-button
                type="primary"
                :loading="loading"
                @click="handleChangePassword"
              >
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped lang="scss">
.profile-page {
  .card-header {
    font-weight: 600;
  }

  .user-info {
    .info-item {
      display: flex;
      align-items: center;
      margin-bottom: 15px;

      label {
        width: 80px;
        color: var(--el-text-color-secondary);
      }

      span {
        color: var(--el-text-color-primary);
      }
    }
  }
}
</style>
