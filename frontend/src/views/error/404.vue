<template>
  <div class="error-page">
    <div class="error-content">
      <h1 class="error-code">404</h1>
      <p class="error-message">页面未找到</p>
      <el-button type="primary" @click="$router.push('/')">
        返回首页
      </el-button>
    </div>
  </div>
</template>

<style scoped lang="scss">
.error-page {
  min-height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;

  .error-content {
    text-align: center;

    .error-code {
      font-size: 120px;
      font-weight: bold;
      color: var(--el-color-primary);
      margin: 0 0 20px;
      line-height: 1;
    }

    .error-message {
      font-size: 20px;
      color: var(--el-text-color-secondary);
      margin-bottom: 30px;
    }
  }
}
</style>
