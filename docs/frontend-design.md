# Dancer DNS 管理前端设计文档

## 1. 概述

Dancer DNS 管理系统前端 - 基于 Vue 3 + Vite + Element Plus + Pinia 的后台管理界面。

**技术栈**：
- **框架**: Vue 3 (Composition API + `<script setup>`)
- **构建工具**: Vite 5
- **UI 组件库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **样式**: SCSS + CSS Variables (支持暗黑模式)

**架构特点**：
- 左侧导航菜单 + 右侧多标签页内容区
- 组件高内聚低耦合，细粒度封装
- 主题系统：亮色/暗黑双模式
- API 层统一封装，支持自动 Token 刷新
- 基于角色的权限控制（RBAC）

---

## 2. 目录结构

```
frontend/                          # 前端源代码根目录
├── public/                        # 静态资源（不经过构建）
│   └── favicon.ico
├── src/
│   ├── api/                       # API 接口封装
│   │   ├── auth.ts               # 认证相关接口
│   │   ├── user.ts               # 用户管理接口
│   │   ├── zone.ts               # Zone 管理接口
│   │   ├── domain.ts             # Domain 管理接口
│   │   └── health.ts             # 健康检查接口
│   │
│   ├── components/               # 组件库
│   │   ├── common/               # 通用组件
│   │   │   ├── AppTable.vue      # 表格组件（基于 ElTable 封装）
│   │   │   ├── AppForm.vue       # 表单组件
│   │   │   ├── AppDialog.vue     # 对话框组件
│   │   │   ├── AppSearch.vue     # 搜索组件
│   │   │   ├── AppPagination.vue # 分页组件
│   │   │   ├── AppCard.vue       # 卡片组件
│   │   │   ├── AppTag.vue        # 标签组件
│   │   │   ├── AppEmpty.vue      # 空状态组件
│   │   │   ├── AppLoading.vue    # 加载状态组件
│   │   │   └── AppIcon.vue       # 图标组件
│   │   │
│   │   ├── business/             # 业务组件
│   │   │   ├── DomainForm.vue    # Domain 表单
│   │   │   ├── ZoneForm.vue      # Zone 表单
│   │   │   ├── UserForm.vue      # 用户表单
│   │   │   ├── IPInput.vue       # IP 地址输入组件
│   │   │   └── TTLSelect.vue     # TTL 选择组件
│   │   │
│   │   └── layout/               # 布局组件
│   │       ├── AppLayout.vue     # 主布局
│   │       ├── Sidebar.vue       # 侧边栏
│   │       ├── Navbar.vue        # 顶部导航
│   │       ├── TabsBar.vue       # 标签页栏
│   │       ├── Breadcrumb.vue    # 面包屑
│   │       └── AppMain.vue       # 主内容区
│   │
│   ├── views/                    # 页面视图
│   │   ├── login/                # 登录页
│   │   │   └── index.vue
│   │   │
│   │   ├── dashboard/            # 仪表盘
│   │   │   └── index.vue
│   │   │
│   │   ├── zone/                 # Zone 管理
│   │   │   ├── list.vue
│   │   │   └── form.vue
│   │   │
│   │   ├── domain/               # Domain 管理
│   │   │   ├── list.vue
│   │   │   └── form.vue
│   │   │
│   │   ├── user/                 # 用户管理（Admin）
│   │   │   ├── list.vue
│   │   │   └── form.vue
│   │   │
│   │   ├── profile/              # 个人中心
│   │   │   └── index.vue
│   │   │
│   │   └── error/                # 错误页面
│   │       ├── 404.vue
│   │       └── 403.vue
│   │
│   ├── composables/              # 组合式函数
│   │   ├── useAuth.ts           # 认证相关
│   │   ├── useTable.ts          # 表格逻辑封装
│   │   ├── useForm.ts           # 表单逻辑封装
│   │   ├── useDialog.ts         # 对话框逻辑封装
│   │   ├── useTheme.ts          # 主题切换
│   │   ├── usePermission.ts     # 权限检查
│   │   └── useTabs.ts           # 标签页管理
│   │
│   ├── stores/                   # Pinia 状态管理
│   │   ├── auth.ts              # 认证状态
│   │   ├── user.ts              # 用户信息
│   │   ├── app.ts               # 应用状态（主题、侧边栏等）
│   │   ├── tabs.ts              # 标签页状态
│   │   └── permission.ts        # 权限状态
│   │
│   ├── router/                   # 路由配置
│   │   ├── index.ts             # 路由实例
│   │   ├── routes.ts            # 路由定义
│   │   ├── guards.ts            # 路由守卫
│   │   └── modules/             # 按模块拆分的路由
│   │       ├── auth.ts
│   │       ├── zone.ts
│   │       ├── domain.ts
│   │       └── user.ts
│   │
│   ├── utils/                    # 工具函数
│   │   ├── request.ts           # Axios 封装
│   │   ├── storage.ts           # 本地存储封装
│   │   ├── validate.ts          # 表单验证规则
│   │   ├── format.ts            # 格式化函数
│   │   ├── theme.ts             # 主题工具
│   │   └── constants.ts         # 常量定义
│   │
│   ├── styles/                   # 样式文件
│   │   ├── variables.scss       # SCSS 变量
│   │   ├── mixins.scss          # SCSS Mixins
│   │   ├── element-plus.scss    # Element Plus 覆盖
│   │   ├── transition.scss      # 过渡动画
│   │   ├── index.scss           # 入口样式
│   │   └── theme/               # 主题样式
│   │       ├── light.scss
│   │       └── dark.scss
│   │
│   ├── types/                    # TypeScript 类型定义
│   │   ├── api.ts               # API 响应类型
│   │   ├── user.ts              # 用户相关类型
│   │   ├── zone.ts              # Zone 相关类型
│   │   ├── domain.ts            # Domain 相关类型
│   │   ├── route.ts             # 路由类型
│   │   └── common.ts            # 通用类型
│   │
│   ├── App.vue                   # 根组件
│   └── main.ts                   # 入口文件
│
├── index.html                    # HTML 模板
├── vite.config.ts               # Vite 配置
├── tsconfig.json                # TypeScript 配置
├── package.json                 # 依赖配置
├── .env                         # 环境变量（开发）
├── .env.production              # 环境变量（生产）
└── eslint.config.js             # ESLint 配置

static/                          # 静态文件根目录（开发时使用）
├── index.html                   # 入口 HTML
├── assets/                      # 构建产物
│   ├── css/
│   ├── js/
│   └── fonts/
└── favicon.ico
```

---

## 3. 技术架构

### 3.1 项目架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         Vue 3 Application                        │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │   Views      │  │  Components  │  │   Composables        │  │
│  │   (Pages)    │  │   (UI/Biz)   │  │   (Logic)            │  │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────┘  │
│         │                 │                     │              │
│  ┌──────▼─────────────────▼─────────────────────▼───────────┐  │
│  │                    Pinia Stores                           │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │  │
│  │  │  Auth    │  │   App    │  │   Tabs   │  │Permission│  │  │
│  │  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │  │
│  └─────────────────────────┬─────────────────────────────────┘  │
│                            │                                    │
│  ┌─────────────────────────▼─────────────────────────────────┐  │
│  │                    API Layer (Axios)                       │  │
│  │         统一请求封装 / 拦截器 / Token 刷新 / 错误处理        │  │
│  └─────────────────────────┬─────────────────────────────────┘  │
│                            │                                    │
└────────────────────────────┼────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Backend API    │
                    └─────────────────┘
```

### 3.2 请求流程

```
┌─────────┐    ┌──────────────┐    ┌─────────────────┐    ┌─────────┐
│ Component│───>│  API Module  │───>│ Request (Axios) │───>│ Server  │
└─────────┘    └──────────────┘    └─────────────────┘    └─────────┘
                                                    │
                                                    ▼
┌─────────┐    ┌──────────────┐    ┌─────────────────┐
│ Component│<───│  Pinia Store │<───│  Response Data  │
└─────────┘    └──────────────┘    └─────────────────┘
```

---

## 4. 核心功能设计

### 4.1 认证与授权

**登录流程**：
1. 用户输入用户名/密码
2. 调用 `/api/auth/login` 获取 Token
3. Token 存储在 Pinia + localStorage
4. 获取当前用户信息并存储
5. 跳转到首页

**Token 刷新**：
- Token 即将过期时自动刷新
- 使用 `/api/auth/refresh` 接口
- 失败则跳转到登录页

**权限控制**：
- 路由级别：meta.requiresAuth / meta.requiresAdmin
- 组件级别：v-permission 指令
- 按钮级别：hasPermission() 方法

### 4.2 标签页系统

**功能**：
- 多标签页导航
- 标签可关闭（右键菜单：关闭、关闭其他、关闭全部）
- 刷新当前标签页
- 标签页状态持久化（可选）
- 离开页面前确认提示

**数据结构**：
```typescript
interface Tab {
  id: string           // 唯一标识（path + query 的 hash）
  title: string        // 显示标题
  path: string         // 路由路径
  query: object        // 查询参数
  closable: boolean    // 是否可关闭
  icon?: string        // 图标
  cached: boolean      // 是否缓存组件
}
```

### 4.3 主题系统

**实现方式**：
- CSS Variables 定义主题变量
- Element Plus 主题定制
- 跟随系统主题 / 手动切换
- localStorage 持久化主题偏好

**主题变量**：
```scss
// 亮色主题
:root {
  --app-bg: #f5f7fa;
  --sidebar-bg: #ffffff;
  --sidebar-text: #303133;
  --navbar-bg: #ffffff;
  --content-bg: #ffffff;
}

// 暗黑主题
[data-theme="dark"] {
  --app-bg: #141414;
  --sidebar-bg: #1d1d1d;
  --sidebar-text: #e0e0e0;
  --navbar-bg: #1d1d1d;
  --content-bg: #262626;
}
```

---

## 5. 组件设计

### 5.1 通用组件（Common）

#### AppTable - 表格组件

**功能**：
- 基于 ElTable 封装
- 自动处理加载状态
- 支持分页
- 支持多选
- 支持操作列
- 支持列配置持久化

**Props**：
```typescript
interface AppTableProps {
  data: any[]                    // 表格数据
  columns: Column[]              // 列配置
  loading?: boolean              // 加载状态
  pagination?: PaginationConfig  // 分页配置
  selection?: boolean            // 是否多选
  rowKey?: string                // 行唯一标识
  operations?: Operation[]       // 操作列
}
```

**使用示例**：
```vue
<AppTable
  :data="tableData"
  :columns="columns"
  :loading="loading"
  :pagination="pagination"
  @page-change="handlePageChange"
  @selection-change="handleSelectionChange"
>
  <template #actions="{ row }">
    <el-button @click="edit(row)">编辑</el-button>
    <el-button @click="delete(row)">删除</el-button>
  </template>
</AppTable>
```

#### AppForm - 表单组件

**功能**：
- 基于 ElForm 封装
- 支持表单配置化
- 自动处理验证
- 支持表单操作按钮
- 支持分组显示

**Props**：
```typescript
interface AppFormProps {
  model: object           // 表单数据
  fields: FormField[]     // 字段配置
  rules?: object          // 验证规则
  loading?: boolean       // 提交加载状态
  labelWidth?: string     // 标签宽度
  inline?: boolean        // 行内表单
}
```

#### AppDialog - 对话框组件

**功能**：
- 基于 ElDialog 封装
- 支持确定/取消按钮
- 自动处理加载状态
- 支持拖拽
- 支持全屏

**Props**：
```typescript
interface AppDialogProps {
  visible: boolean        // 显示状态
  title: string          // 标题
  width?: string         // 宽度
  fullscreen?: boolean   // 全屏
  loading?: boolean      // 确定按钮加载状态
  showFooter?: boolean   // 显示底部
}
```

#### AppSearch - 搜索组件

**功能**：
- 搜索表单封装
- 支持展开/收起
- 支持查询和重置
- 支持插槽扩展

### 5.2 业务组件（Business）

#### DomainForm - Domain 表单

**功能**：
- Zone 选择（下拉框）
- Domain 名称输入
- IP 地址列表（支持动态增删）
- TTL 选择

**使用场景**：
- Domain 创建/编辑

#### IPInput - IP 地址输入

**功能**：
- IP 格式验证
- 支持多个 IP
- 支持复制粘贴批量添加
- 支持 IP 去重

#### ZoneForm - Zone 表单

**功能**：
- Zone 名称输入
- 格式验证（FQDN）

#### UserForm - 用户表单

**功能**：
- 用户名输入
- 密码输入（新建时必填）
- 用户类型选择（admin/normal）

### 5.3 布局组件（Layout）

#### AppLayout - 主布局

**结构**：
```
┌──────────────────────────────────────────┐
│                 Navbar                   │
├──────────────┬───────────────────────────┤
│              │                           │
│   Sidebar    │      AppMain (Tabs)       │
│              │                           │
│              │   ┌──────────────────┐    │
│              │   │   TabsBar        │    │
│              │   ├──────────────────┤    │
│              │   │                  │    │
│              │   │   Router View    │    │
│              │   │   (Cached)       │    │
│              │   │                  │    │
│              │   └──────────────────┘    │
└──────────────┴───────────────────────────┘
```

#### Sidebar - 侧边栏

**功能**：
- 菜单折叠/展开
- 支持多级菜单
- 根据权限动态显示菜单
- 当前菜单高亮

**菜单配置**：
```typescript
interface MenuItem {
  path: string
  title: string
  icon?: string
  children?: MenuItem[]
  hidden?: boolean
  permission?: string
}
```

#### TabsBar - 标签页栏

**功能**：
- 显示已打开的标签页
- 支持关闭单个/关闭其他/关闭全部
- 支持刷新当前页
- 支持拖拽排序（可选）
- 超出时横向滚动

#### Navbar - 顶部导航

**功能**：
- 折叠/展开侧边栏按钮
- 面包屑导航
- 全屏切换
- 主题切换
- 用户信息下拉菜单
- 退出登录

---

## 6. 路由设计

### 6.1 路由结构

```typescript
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    hidden: true
  },
  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'dashboard' }
      }
    ]
  },
  {
    path: '/dns',
    component: Layout,
    meta: { title: 'DNS 管理', icon: 'dns' },
    children: [
      {
        path: 'zones',
        component: () => import('@/views/zone/list.vue'),
        meta: { title: 'Zone 管理', permission: 'admin' }
      },
      {
        path: 'domains',
        component: () => import('@/views/domain/list.vue'),
        meta: { title: 'Domain 管理' }
      }
    ]
  },
  {
    path: '/system',
    component: Layout,
    meta: { title: '系统管理', icon: 'system', permission: 'admin' },
    children: [
      {
        path: 'users',
        component: () => import('@/views/user/list.vue'),
        meta: { title: '用户管理' }
      }
    ]
  },
  {
    path: '/profile',
    component: Layout,
    hidden: true,
    children: [
      {
        path: '',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心' }
      }
    ]
  }
]
```

### 6.2 路由元信息

```typescript
interface RouteMeta {
  title: string           // 页面标题
  icon?: string          // 菜单图标
  hidden?: boolean       // 是否隐藏
  permission?: string    // 所需权限
  requiresAuth?: boolean // 是否需要登录
  keepAlive?: boolean    // 是否缓存
  activeMenu?: string    // 高亮菜单路径
}
```

### 6.3 路由守卫

**全局前置守卫**：
1. 检查是否需要登录
2. 检查 Token 是否有效
3. 检查用户权限
4. 生成动态菜单

---

## 7. 状态管理设计

### 7.1 Store 划分

```
stores/
├── auth.ts       # 认证状态（token, userInfo）
├── user.ts       # 用户信息
├── app.ts        # 应用状态（theme, sidebar, language）
├── tabs.ts       # 标签页状态
└── permission.ts # 权限状态（roles, routes）
```

### 7.2 Auth Store

```typescript
interface AuthState {
  token: string | null
  userInfo: UserInfo | null
  isLoggedIn: boolean
}

const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    token: localStorage.getItem('token'),
    userInfo: null,
    isLoggedIn: false
  }),
  
  actions: {
    async login(credentials: LoginForm) {
      const { data } = await loginApi(credentials)
      this.token = data.token
      localStorage.setItem('token', data.token)
      await this.getUserInfo()
    },
    
    async getUserInfo() {
      const { data } = await getUserInfoApi()
      this.userInfo = data
      this.isLoggedIn = true
    },
    
    logout() {
      this.token = null
      this.userInfo = null
      this.isLoggedIn = false
      localStorage.removeItem('token')
    },
    
    async refreshToken() {
      const { data } = await refreshTokenApi()
      this.token = data.token
      localStorage.setItem('token', data.token)
    }
  }
})
```

### 7.3 Tabs Store

```typescript
interface TabsState {
  tabs: Tab[]
  activeTab: string
}

const useTabsStore = defineStore('tabs', {
  state: (): TabsState => ({
    tabs: [],
    activeTab: ''
  }),
  
  actions: {
    addTab(tab: Tab) {
      if (!this.tabs.find(t => t.id === tab.id)) {
        this.tabs.push(tab)
      }
      this.activeTab = tab.id
    },
    
    removeTab(id: string) {
      const index = this.tabs.findIndex(t => t.id === id)
      this.tabs.splice(index, 1)
      // 激活相邻标签页
    },
    
    closeOthers(id: string) {
      const current = this.tabs.find(t => t.id === id)
      this.tabs = current ? [current] : []
    },
    
    closeAll() {
      this.tabs = []
      router.push('/dashboard')
    }
  }
})
```

---

## 8. API 层设计

### 8.1 请求封装

**axios 配置**：
```typescript
// utils/request.ts
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
  timeout: 10000
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const token = useAuthStore().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const { code, message, data } = response.data
    if (code === 'success') {
      return data
    }
    ElMessage.error(message)
    return Promise.reject(new Error(message))
  },
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore().logout()
      router.push('/login')
    }
    return Promise.reject(error)
  }
)
```

### 8.2 API 模块

```typescript
// api/auth.ts
export const loginApi = (data: LoginForm) => 
  request.post('/auth/login', data)

export const refreshTokenApi = () => 
  request.post('/auth/refresh')

// api/user.ts
export const getUserListApi = () => 
  request.post('/user/list')

export const createUserApi = (data: CreateUserForm) => 
  request.post('/user/create', data)

// api/zone.ts
export const getZoneListApi = () => 
  request.post('/dns/zones/list')

export const createZoneApi = (data: CreateZoneForm) => 
  request.post('/dns/zones/create', data)

// api/domain.ts
export const getDomainListApi = (data: { zone: string }) => 
  request.post('/dns/domains/list', data)

export const createDomainApi = (data: CreateDomainForm) => 
  request.post('/dns/domains/create', data)
```

---

## 9. 类型定义

### 9.1 API 类型

```typescript
// types/api.ts
export interface ApiResponse<T = any> {
  code: string
  message: string
  data?: T
}

export interface ListResponse<T> {
  list: T[]
  total: number
}
```

### 9.2 业务类型

```typescript
// types/user.ts
export interface User {
  id: string
  username: string
  user_type: 'admin' | 'normal'
  created_at: number
  updated_at: number
}

export interface LoginForm {
  username: string
  password: string
}

// types/zone.ts
export interface Zone {
  zone: string
  record_count: number
  created_at: number
  updated_at: number
}

export interface CreateZoneForm {
  zone: string
}

// types/domain.ts
export interface Domain {
  zone: string
  domain: string
  name: string
  ips: string[]
  ttl: number
  record_count: number
  created_at: number
  updated_at: number
}

export interface CreateDomainForm {
  zone: string
  domain: string
  ips: string[]
  ttl: number
}
```

---

## 10. 开发规范

### 10.1 命名规范

- **组件名**: PascalCase（AppTable, DomainForm）
- **文件名**: 组件 PascalCase，其他 camelCase
- **变量/函数**: camelCase
- **常量**: UPPER_SNAKE_CASE
- **类型/接口**: PascalCase
- **CSS 类名**: kebab-case

### 10.2 代码规范

- 使用 Composition API + `<script setup>`
- Props 定义使用 TypeScript 接口
- 组件使用 defineProps / defineEmits
- API 调用放在 composables 或 store 中
- 避免在模板中写复杂表达式
- 使用 computed 缓存计算属性

### 10.3 文件组织规范

```typescript
// 组件文件结构
<script setup lang="ts">
// 1. imports
import { ref, computed } from 'vue'
import type { PropType } from 'vue'

// 2. types/interfaces
interface Props {
  // ...
}

// 3. props/emits
const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

// 4. reactive state
const loading = ref(false)
const data = ref([])

// 5. computed
const formattedData = computed(() => {
  return data.value.map(...)
})

// 6. methods
const fetchData = async () => {
  // ...
}

// 7. lifecycle
onMounted(() => {
  fetchData()
})
</script>

<template>
  <!-- template content -->
</template>

<style scoped lang="scss">
/* styles */
</style>
```

---

## 11. 构建配置

### 11.1 Vite 配置

```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  },
  
  build: {
    outDir: '../static',
    assetsDir: 'assets',
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'vue-vendor': ['vue', 'vue-router', 'pinia']
        }
      }
    }
  }
})
```

### 11.2 环境变量

```
# .env
VITE_API_BASE_URL=/api
VITE_APP_TITLE=Dancer DNS

# .env.production
VITE_API_BASE_URL=/api
VITE_APP_TITLE=Dancer DNS
```

---

## 12. 开发计划

### Phase 1: 基础搭建（2天）

1. 初始化 Vite + Vue 3 项目
2. 安装依赖（Element Plus, Pinia, Vue Router, Axios）
3. 配置 TypeScript
4. 配置 ESLint + Prettier
5. 搭建目录结构
6. 配置路径别名
7. 配置代理

### Phase 2: 核心功能（3天）

1. **路由系统**
   - 配置路由
   - 实现路由守卫
   - 动态菜单生成

2. **状态管理**
   - 实现 Auth Store
   - 实现 App Store（主题、侧边栏）
   - 实现 Tabs Store

3. **布局组件**
   - AppLayout
   - Sidebar
   - Navbar
   - TabsBar

4. **API 层**
   - Axios 封装
   - 请求/响应拦截器
   - Token 刷新机制

5. **登录页面**
   - 登录表单
   - 表单验证
   - 登录逻辑

### Phase 3: 通用组件（2天）

1. AppTable - 表格组件
2. AppForm - 表单组件
3. AppDialog - 对话框组件
4. AppSearch - 搜索组件
5. AppPagination - 分页组件
6. AppLoading - 加载组件
7. AppEmpty - 空状态组件

### Phase 4: 业务页面（3天）

1. **仪表盘**
   - 统计卡片
   - Zone/Domain 数量展示

2. **Zone 管理**
   - Zone 列表页
   - Zone 创建/编辑弹窗
   - 删除确认

3. **Domain 管理**
   - Domain 列表页
   - Domain 创建/编辑弹窗
   - IP 输入组件

4. **用户管理（Admin）**
   - 用户列表页
   - 用户创建/编辑弹窗
   - 权限控制

5. **个人中心**
   - 用户信息展示
   - 修改密码

### Phase 5: 主题与优化（2天）

1. 主题系统实现
2. 暗黑模式切换
3. 样式统一和优化
4. 页面过渡动画
5. 响应式适配

### Phase 6: 集成与部署（1天）

1. Go embed 集成
2. 生产构建配置
3. 静态文件服务配置

---

## 13. 页面原型

### 13.1 登录页

```
┌──────────────────────────────────────┐
│                                      │
│           ┌──────────────────────┐   │
│           │   Dancer DNS         │   │
│           │                      │   │
│           │  ┌────────────────┐  │   │
│           │  │  用户名         │  │   │
│           │  └────────────────┘  │   │
│           │  ┌────────────────┐  │   │
│           │  │  密码           │  │   │
│           │  └────────────────┘  │   │
│           │                      │   │
│           │  ┌────────────────┐  │   │
│           │  │     登录        │  │   │
│           │  └────────────────┘  │   │
│           │                      │   │
│           └──────────────────────┘   │
│                                      │
└──────────────────────────────────────┘
```

### 13.2 主布局

```
┌────────────────────────────────────────────────────────────┐
│ ≡ │ Dancer DNS           👤 admin  ▼ │ 🌙 │ ⛶ │
├───────┬────────────────────────────────────────────────────┤
│       │  首页 / DNS管理 / Domain管理                       │
│  📊   ├────────────────────────────────────────────────────┤
│  Dashboard │  Domain管理 │ Zone管理 │ ×                    │
│       ├────────────────────────────────────────────────────┤
│  🌐   │                                                    │
│  DNS  │  ┌──────────────────────────────────────────────┐  │
│  Zone │  │ 搜索 🔍                    [+ 新建 Domain]   │  │
│  Domain│  └──────────────────────────────────────────────┘  │
│       │                                                    │
│  👥   │  ┌──────────────────────────────────────────────┐  │
│  用户 │  │ 域名        Zone        IP 地址      TTL   操作│  │
│  管理 │  │ www      example.com  192.168.1.1   300   ⚙️ │  │
│       │  │ @        example.com  192.168.1.2   300   ⚙️ │  │
│       │  │ api      example.com  10.0.0.1      600   ⚙️ │  │
│       │  └──────────────────────────────────────────────┘  │
│       │                                                    │
│       │  ┌──────────────────────────────────────────────┐  │
│       │  │ 共 10 条                  < 1 2 3 ... 10 >   │  │
│       │  └──────────────────────────────────────────────┘  │
└───────┴────────────────────────────────────────────────────┘
```

---

## 14. 注意事项

### 14.1 开发注意事项

1. **API 兼容性**
   - 所有 API 使用 POST 方法（除 /api/health 外）
   - Token 从 Authorization Header 获取
   - 错误处理使用后端返回的 code 和 message

2. **权限控制**
   - Admin 才能访问 Zone 管理和用户管理
   - 普通用户只能管理 Domain
   - 路由级别和按钮级别都要做权限控制

3. **表单验证**
   - Zone: 有效的 FQDN 格式
   - Domain: 子域名格式（支持 @ 表示根）
   - IP: IPv4 格式验证
   - TTL: 最小值为 1

4. **交互细节**
   - 删除操作需要确认对话框
   - 表单提交时按钮显示加载状态
   - 操作成功/失败显示 Toast 提示
   - 列表为空时显示空状态

### 14.2 性能优化

1. 路由懒加载
2. 组件按需引入
3. 使用 KeepAlive 缓存标签页
4. 图片懒加载（如有）
5. API 请求防抖/节流

---

## 15. 参考资源

- [Vue 3 文档](https://vuejs.org/)
- [Element Plus 文档](https://element-plus.org/)
- [Pinia 文档](https://pinia.vuejs.org/)
- [Vue Router 文档](https://router.vuejs.org/)
- [Vite 文档](https://vitejs.dev/)

---

**文档版本**: 1.0  
**创建日期**: 2026-03-07  
**作者**: OpenCode
