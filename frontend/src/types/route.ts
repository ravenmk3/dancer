import type { RouteRecordRaw } from 'vue-router'

export interface MenuItem {
  path: string
  title: string
  icon?: string
  children?: MenuItem[]
  hidden?: boolean
  permission?: string
}

export interface RouteMeta {
  title: string
  icon?: string
  hidden?: boolean
  permission?: string
  requiresAuth?: boolean
  keepAlive?: boolean
  activeMenu?: string
}

export interface AppRouteRecordRaw extends Omit<RouteRecordRaw, 'meta' | 'children'> {
  meta?: RouteMeta
  children?: AppRouteRecordRaw[]
}
