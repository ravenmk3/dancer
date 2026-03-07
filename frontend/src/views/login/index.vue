<script setup lang="ts">
import { ref, reactive, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import type { LoginForm } from '@/types/user'

const router = useRouter()
const authStore = useAuthStore()

const loading = ref(false)
const loginForm = reactive<LoginForm>({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 32, message: '长度在 3 到 32 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为 6 个字符', trigger: 'blur' }
  ]
}

const formRef = ref()

const handleLogin = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid: boolean) => {
    if (!valid) return
    
    loading.value = true
    console.log('[Login] ========== Login Process Started ==========')
    
    try {
      // 1. 执行登录
      console.log('[Login] Step 1: Calling authStore.login()')
      const loginResult = await authStore.login(loginForm)
      console.log('[Login] Step 1 Complete: loginResult =', loginResult)
      
      // 2. 等待 DOM 更新，确保 computed 属性已更新
      console.log('[Login] Step 2: Waiting for nextTick...')
      await nextTick()
      console.log('[Login] Step 2 Complete: nextTick done')
      
      // 3. 验证登录状态
      console.log('[Login] Step 3: Checking auth state...')
      console.log('[Login] - authStore.token exists:', !!authStore.token)
      console.log('[Login] - authStore.userInfo exists:', !!authStore.userInfo)
      console.log('[Login] - authStore.isLoggedIn:', authStore.isLoggedIn)
      
      if (!authStore.isLoggedIn) {
        throw new Error('登录状态验证失败：isLoggedIn 为 false')
      }
      
      // 4. 显示成功消息
      console.log('[Login] Step 4: Showing success message')
      ElMessage.success('登录成功')
      
      // 5. 导航到首页 - 使用 push 而不是 replace
      console.log('[Login] Step 5: Navigating to /dashboard')
      await router.push('/dashboard')
      console.log('[Login] ========== Login Process Complete ==========')
      
    } catch (error: any) {
      console.error('[Login] ========== Login Failed ==========')
      console.error('[Login] Error:', error)
      console.error('[Login] Error message:', error?.message)
      console.error('[Login] Error response:', error?.response?.data)
      
      // 显示具体错误信息
      const errorMessage = error?.response?.data?.message || error?.message || '登录失败'
      ElMessage.error(errorMessage)
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <div class="login-page">
    <div class="login-box">
      <div class="login-header">
        <div class="logo">D</div>
        <h1 class="title">Dancer DNS</h1>
        <p class="subtitle">DNS 管理系统</p>
      </div>
      
      <el-form
        ref="formRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            :prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <el-button
          type="primary"
          size="large"
          class="login-btn"
          :loading="loading"
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<style scoped lang="scss">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, var(--el-color-primary-light-9) 0%, var(--el-color-primary-light-7) 100%);

  .login-box {
    width: 400px;
    padding: 40px;
    background-color: var(--el-bg-color);
    border-radius: 12px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);

    .login-header {
      text-align: center;
      margin-bottom: 30px;

      .logo {
        width: 64px;
        height: 64px;
        background: linear-gradient(135deg, var(--el-color-primary) 0%, var(--el-color-primary-light-3) 100%);
        border-radius: 16px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin: 0 auto 20px;
        color: #fff;
        font-size: 32px;
        font-weight: bold;
      }

      .title {
        font-size: 24px;
        font-weight: 600;
        color: var(--el-text-color-primary);
        margin: 0 0 8px;
      }

      .subtitle {
        font-size: 14px;
        color: var(--el-text-color-secondary);
        margin: 0;
      }
    }

    .login-form {
      .el-input {
        :deep(.el-input__inner) {
          height: 44px;
        }
      }

      .login-btn {
        width: 100%;
        height: 44px;
        margin-top: 10px;
        font-size: 16px;
      }
    }
  }
}
</style>
