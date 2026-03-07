<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getZoneListApi } from '@/api/zone'
import { getUserListApi } from '@/api/user'
import type { Zone } from '@/types/zone'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const zoneCount = ref(0)
const domainCount = ref(0)
const userCount = ref(0)
const loading = ref(false)

const stats = ref([
  { title: 'Zone 数量', value: 0, icon: 'Folder', color: '#409EFF' },
  { title: 'Domain 数量', value: 0, icon: 'Document', color: '#67C23A' },
  { title: '用户数量', value: 0, icon: 'UserFilled', color: '#E6A23C' }
])

const fetchStats = async () => {
  loading.value = true
  try {
    // 并行获取 Zone 和用户列表
    const [zoneRes, userRes] = await Promise.all([
      getZoneListApi().catch(() => ({ zones: [] })),
      getUserListApi().catch(() => ({ users: [] }))
    ])
    
    const zones = zoneRes.zones || []
    const users = userRes.users || []
    
    // 使用 zone 的 record_count 字段统计 domain 数量，避免并发请求
    const domainCountValue = zones.reduce((sum: number, zone: Zone) => sum + (zone.record_count || 0), 0)
    
    zoneCount.value = zones.length
    domainCount.value = domainCountValue
    userCount.value = users.length
    
    stats.value[0].value = zoneCount.value
    stats.value[1].value = domainCount.value
    stats.value[2].value = userCount.value
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<template>
  <div class="dashboard-page">
    <el-row :gutter="20">
      <el-col v-for="stat in stats" :key="stat.title" :span="8">
        <el-card shadow="hover" class="stat-card" v-loading="loading">
          <div class="stat-content">
            <div class="stat-icon" :style="{ backgroundColor: stat.color + '20', color: stat.color }">
              <el-icon :size="32">
                <component :is="stat.icon" />
              </el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-title">{{ stat.title }}</div>
              <div class="stat-value">{{ stat.value }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="welcome-card" shadow="never">
      <template #header>
        <div class="card-header">
          <span>欢迎使用 Dancer</span>
        </div>
      </template>
      <div class="welcome-content">
        <p>Dancer 是一个基于 etcd 的 DNS 管理工具，专为 CoreDNS 设计。</p>
        <p>您可以通过本系统管理 DNS Zone 和 Domain 记录。</p>
        <el-divider />
        <h4>快速开始</h4>
        <ul>
          <li>管理员用户可以在"DNS 管理"菜单中创建 Zone</li>
          <li>在"Domain 管理"中添加域名记录</li>
          <li>支持批量添加 IP 地址和自定义 TTL</li>
        </ul>
      </div>
    </el-card>
  </div>
</template>

<style scoped lang="scss">
.dashboard-page {
  .stat-card {
    margin-bottom: 20px;

    .stat-content {
      display: flex;
      align-items: center;

      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 15px;
      }

      .stat-info {
        .stat-title {
          font-size: 14px;
          color: var(--el-text-color-secondary);
          margin-bottom: 8px;
        }

        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: var(--el-text-color-primary);
        }
      }
    }
  }

  .welcome-card {
    .card-header {
      font-weight: 600;
      font-size: 16px;
    }

    .welcome-content {
      line-height: 1.8;
      color: var(--el-text-color-regular);

      h4 {
        margin: 20px 0 10px;
        color: var(--el-text-color-primary);
      }

      ul {
        padding-left: 20px;

        li {
          margin-bottom: 8px;
        }
      }
    }
  }
}
</style>
