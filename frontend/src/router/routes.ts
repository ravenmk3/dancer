import type { AppRouteRecordRaw } from '@/types/route'
import AppLayout from '@/components/layout/AppLayout.vue'

export const constantRoutes: AppRouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/index.vue'),
    meta: { hidden: true, title: '登录' }
  },
  {
    path: '/',
    component: AppLayout,
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/index.vue'),
        meta: { title: '仪表盘', icon: 'DataLine', requiresAuth: true }
      }
    ]
  },
  {
    path: '/dns',
    component: AppLayout,
    meta: { title: 'DNS 管理', icon: 'Connection', requiresAuth: true },
    children: [
      {
        path: 'zones',
        name: 'Zones',
        component: () => import('@/views/zone/list.vue'),
        meta: { title: 'Zone 管理', icon: 'Folder', permission: 'admin', requiresAuth: true }
      },
      {
        path: 'domains',
        name: 'Domains',
        component: () => import('@/views/domain/list.vue'),
        meta: { title: 'Domain 管理', icon: 'Document', requiresAuth: true }
      }
    ]
  },
  {
    path: '/system',
    component: AppLayout,
    meta: { title: '系统管理', icon: 'Setting', permission: 'admin', requiresAuth: true },
    children: [
      {
        path: 'users',
        name: 'Users',
        component: () => import('@/views/user/list.vue'),
        meta: { title: '用户管理', icon: 'User', requiresAuth: true }
      }
    ]
  },
  {
    path: '/profile',
    component: AppLayout,
    meta: { hidden: true, requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Profile',
        component: () => import('@/views/profile/index.vue'),
        meta: { title: '个人中心', requiresAuth: true }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: { hidden: true, title: '页面未找到' }
  }
]
