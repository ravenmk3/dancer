<template>
  <router-view />
</template>

<style lang="scss">
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
}

:root {
  --app-bg: #f5f7fa;
  --sidebar-bg: #ffffff;
  --sidebar-text: #303133;
  --navbar-bg: #ffffff;
}

[data-theme="dark"] {
  --app-bg: #141414;
  --sidebar-bg: #1d1d1d;
  --sidebar-text: #e0e0e0;
  --navbar-bg: #1d1d1d;
}
</style>
