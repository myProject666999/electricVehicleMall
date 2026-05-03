<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const stats = ref({
  totalUsers: 0,
  totalOrders: 0,
  totalVehicles: 0,
  totalRevenue: 0,
  pendingOrders: 0,
  todayOrders: 0
})

const orderChartData = ref([])
const recentOrders = ref([])

// 统计卡片配置
const statCards = [
  { 
    title: '总用户数', 
    key: 'totalUsers', 
    icon: 'User', 
    color: '#409EFF',
    suffix: ''
  },
  { 
    title: '总订单数', 
    key: 'totalOrders', 
    icon: 'Document', 
    color: '#67C23A',
    suffix: ''
  },
  { 
    title: '电动车数量', 
    key: 'totalVehicles', 
    icon: 'Truck', 
    color: '#E6A23C',
    suffix: ''
  },
  { 
    title: '总收入', 
    key: 'totalRevenue', 
    icon: 'Money', 
    color: '#F56C6C',
    suffix: ' 元'
  }
]

// 获取统计数据
const fetchStats = async () => {
  loading.value = true
  try {
    const res = await request.get('/admin/statistics')
    if (res.code === 200 || res.code === 0) {
      stats.value = {
        totalUsers: res.data.total_orders || 0,
        totalOrders: res.data.total_orders || 0,
        totalVehicles: 0,
        totalRevenue: res.data.today_sales || 0,
        pendingOrders: res.data.pending_ship || 0,
        todayOrders: 0
      }
    }
  } catch (error) {
    console.error('获取统计数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取最近订单
const fetchRecentOrders = async () => {
  try {
    const res = await request.get('/admin/order/list', { params: { page: 1, pageSize: 5 } })
    if (res.code === 200 || res.code === 0) {
      recentOrders.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取最近订单失败:', error)
  }
}

// 订单状态标签
const getStatusTag = (status) => {
  const statusMap = {
    0: { type: 'info', text: '待支付' },
    1: { type: 'warning', text: '待发货' },
    2: { type: 'primary', text: '已发货' },
    3: { type: 'success', text: '已完成' },
    4: { type: 'danger', text: '已取消' }
  }
  return statusMap[status] || { type: 'info', text: '未知' }
}

onMounted(() => {
  fetchStats()
  fetchRecentOrders()
})
</script>

<template>
  <div class="dashboard-container">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6" v-for="card in statCards" :key="card.title">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-value">
                {{ stats[card.key] ? stats[card.key].toLocaleString() : 0 }}
                <span v-if="card.suffix">{{ card.suffix }}</span>
              </div>
              <div class="stat-title">{{ card.title }}</div>
            </div>
            <div class="stat-icon" :style="{ background: card.color }">
              <el-icon size="32" color="#fff">
                <component :is="card.icon" />
              </el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 快捷操作 -->
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :xs="24" :sm="24" :lg="24">
        <el-card shadow="hover">
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" icon="Truck" @click="$router.push('/vehicle')">
              电动车管理
            </el-button>
            <el-button type="success" icon="Document" @click="$router.push('/order')">
              订单管理
            </el-button>
            <el-button type="warning" icon="User" @click="$router.push('/user')">
              用户管理
            </el-button>
            <el-button type="info" icon="Bell" @click="$router.push('/notice')">
              公告管理
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 最近订单 -->
    <el-row :gutter="20">
      <el-col :xs="24" :sm="24" :lg="24">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>最近订单</span>
              <el-button type="text" @click="$router.push('/order')">查看全部</el-button>
            </div>
          </template>
          
          <el-table :data="recentOrders" v-loading="loading" style="width: 100%">
            <el-table-column prop="orderNo" label="订单号" width="180" />
            <el-table-column prop="username" label="用户" width="120" />
            <el-table-column prop="totalAmount" label="金额" width="120">
              <template #default="scope">
                ¥{{ scope.row.totalAmount?.toFixed(2) || '0.00' }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusTag(scope.row.status).type">
                  {{ getStatusTag(scope.row.status).text }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="创建时间" width="180" />
          </el-table>
          
          <el-empty v-if="recentOrders.length === 0 && !loading" description="暂无订单数据" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.dashboard-container {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 8px;
}

.stat-value span {
  font-size: 16px;
  font-weight: normal;
}

.stat-title {
  font-size: 14px;
  color: #909399;
}

.stat-icon {
  width: 70px;
  height: 70px;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  opacity: 0.9;
}

.quick-actions {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
