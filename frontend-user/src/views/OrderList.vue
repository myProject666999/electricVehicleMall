<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const orders = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const activeTab = ref('all')

// 订单状态映射
const statusMap = {
  pending_payment: { text: '待付款', color: 'warning' },
  pending_shipment: { text: '待发货', color: 'info' },
  shipped: { text: '已发货', color: 'primary' },
  completed: { text: '已完成', color: 'success' },
  cancelled: { text: '已取消', color: 'danger' }
}

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (activeTab.value !== 'all') {
      params.status = activeTab.value
    }
    
    const res = await request.get('/order/list', { params })
    orders.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 切换标签
const handleTabChange = (name) => {
  activeTab.value = name
  page.value = 1
  fetchOrderList()
}

// 分页
const handlePageChange = (val) => {
  page.value = val
  fetchOrderList()
}

// 查看订单详情
const handleViewDetail = (order) => {
  router.push(`/order/${order.id}`)
}

// 取消订单
const handleCancelOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post('/order/cancel', { id: order.id })
    ElMessage.success('订单已取消')
    fetchOrderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error)
      ElMessage.error('取消订单失败')
    }
  }
}

// 支付订单
const handlePayOrder = async (order) => {
  try {
    await ElMessageBox.confirm('确定要支付该订单吗？', '提示', {
      confirmButtonText: '立即支付',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await request.post('/order/pay', { id: order.id })
    ElMessage.success('支付成功')
    fetchOrderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('支付失败:', error)
      ElMessage.error('支付失败')
    }
  }
}

// 确认收货
const handleConfirmReceive = async (order) => {
  try {
    await ElMessageBox.confirm('确定已收到货物吗？', '提示', {
      confirmButtonText: '确认收货',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await request.post('/order/receive', { id: order.id })
    ElMessage.success('确认收货成功')
    fetchOrderList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认收货失败:', error)
      ElMessage.error('确认收货失败')
    }
  }
}

onMounted(() => {
  fetchOrderList()
})
</script>

<template>
  <div class="order-list-container">
    <h2 class="page-title">我的订单</h2>
    
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange" style="margin-bottom: 20px;">
      <el-tab-pane label="全部订单" name="all" />
      <el-tab-pane label="待付款" name="pending_payment" />
      <el-tab-pane label="待发货" name="pending_shipment" />
      <el-tab-pane label="已发货" name="shipped" />
      <el-tab-pane label="已完成" name="completed" />
    </el-tabs>
    
    <!-- 订单列表 -->
    <div class="order-list" v-loading="loading">
      <div v-for="order in orders" :key="order.id" class="order-card">
        <!-- 订单头部 -->
        <div class="order-header">
          <div class="order-info">
            <span class="order-no">订单号：{{ order.order_no }}</span>
            <span class="order-time">{{ new Date(order.created_at).toLocaleString() }}</span>
          </div>
          <el-tag :type="statusMap[order.status]?.color || 'info'">
            {{ statusMap[order.status]?.text || order.status }}
          </el-tag>
        </div>
        
        <!-- 订单商品 -->
        <div class="order-items">
          <div 
            v-for="item in order.items" 
            :key="item.id" 
            class="order-item"
            @click="handleViewDetail(order)"
          >
            <el-image 
              :src="item.vehicle?.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
              fit="cover"
              style="width: 80px; height: 80px;"
            />
            <div class="item-info">
              <h4 class="item-name">{{ item.vehicle?.name }}</h4>
              <div class="item-price">
                <span>¥{{ item.price }}</span>
                <span>x {{ item.quantity }}</span>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 订单底部 -->
        <div class="order-footer">
          <div class="order-summary">
            <span>共 {{ order.items?.length || 0 }} 件商品</span>
            <span style="margin-left: 20px;">
              合计：<span class="total-price">¥{{ order.total_amount }}</span>
              <span class="shipping-fee" v-if="order.shipping_fee > 0">(含运费¥{{ order.shipping_fee }})</span>
            </span>
          </div>
          
          <div class="order-actions">
            <el-button 
              type="primary" 
              size="small"
              @click="handleViewDetail(order)"
            >
              查看详情
            </el-button>
            
            <!-- 待付款 -->
            <template v-if="order.status === 'pending_payment'">
              <el-button 
                type="danger" 
                size="small"
                @click="handlePayOrder(order)"
              >
                立即支付
              </el-button>
              <el-button 
                type="text" 
                size="small"
                @click="handleCancelOrder(order)"
              >
                取消订单
              </el-button>
            </template>
            
            <!-- 已发货 -->
            <template v-if="order.status === 'shipped'">
              <el-button 
                type="success" 
                size="small"
                @click="handleConfirmReceive(order)"
              >
                确认收货
              </el-button>
            </template>
          </div>
        </div>
      </div>
      
      <!-- 空状态 -->
      <el-empty 
        v-if="orders.length === 0 && !loading" 
        description="暂无订单"
      >
        <el-button type="primary" @click="router.push('/vehicle')">
          去购物
        </el-button>
      </el-empty>
      
      <!-- 分页 -->
      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="page"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<style scoped>
.order-list-container {
  padding: 20px 0;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  border-bottom: 2px solid #409EFF;
  padding-bottom: 10px;
  display: inline-block;
}

.order-card {
  background: #fff;
  border-radius: 8px;
  margin-bottom: 15px;
  overflow: hidden;
  border: 1px solid #ebeef5;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  background: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.order-info {
  display: flex;
  gap: 20px;
}

.order-no {
  color: #303133;
  font-weight: 500;
}

.order-time {
  color: #909399;
  font-size: 14px;
}

.order-items {
  padding: 15px 20px;
}

.order-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
  cursor: pointer;
  padding: 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

.order-item:last-child {
  margin-bottom: 0;
}

.order-item:hover {
  background: #f5f7fa;
}

.item-info {
  margin-left: 15px;
  flex: 1;
}

.item-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 8px;
  line-height: 1.4;
}

.item-price {
  font-size: 14px;
  color: #909399;
}

.item-price span:first-child {
  color: #f56c6c;
  font-weight: bold;
  margin-right: 15px;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-top: 1px solid #ebeef5;
  background: #fafafa;
}

.order-summary {
  font-size: 14px;
  color: #606266;
}

.total-price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.shipping-fee {
  color: #909399;
  font-size: 12px;
  margin-left: 5px;
}

.order-actions {
  display: flex;
  gap: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
