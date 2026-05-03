<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()

const order = ref(null)
const loading = ref(false)

// 订单状态映射
const statusMap = {
  pending_payment: { text: '待付款', color: 'warning' },
  pending_shipment: { text: '待发货', color: 'info' },
  shipped: { text: '已发货', color: 'primary' },
  completed: { text: '已完成', color: 'success' },
  cancelled: { text: '已取消', color: 'danger' }
}

// 获取订单详情
const fetchOrderDetail = async () => {
  loading.value = true
  try {
    const res = await request.get(`/order/detail/${route.params.id}`)
    order.value = res.data
  } catch (error) {
    console.error('获取订单详情失败:', error)
  } finally {
    loading.value = false
  }
}

// 取消订单
const handleCancelOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post('/order/cancel', { id: order.value.id })
    ElMessage.success('订单已取消')
    fetchOrderDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消订单失败:', error)
      ElMessage.error('取消订单失败')
    }
  }
}

// 支付订单
const handlePayOrder = async () => {
  try {
    await ElMessageBox.confirm('确定要支付该订单吗？', '提示', {
      confirmButtonText: '立即支付',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await request.post('/order/pay', { id: order.value.id })
    ElMessage.success('支付成功')
    fetchOrderDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('支付失败:', error)
      ElMessage.error('支付失败')
    }
  }
}

// 确认收货
const handleConfirmReceive = async () => {
  try {
    await ElMessageBox.confirm('确定已收到货物吗？', '提示', {
      confirmButtonText: '确认收货',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    await request.post('/order/receive', { id: order.value.id })
    ElMessage.success('确认收货成功')
    fetchOrderDetail()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('确认收货失败:', error)
      ElMessage.error('确认收货失败')
    }
  }
}

onMounted(() => {
  fetchOrderDetail()
})
</script>

<template>
  <div class="order-detail-container">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item>
        <router-link to="/">首页</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>
        <router-link to="/order">我的订单</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>订单详情</el-breadcrumb-item>
    </el-breadcrumb>
    
    <div class="order-detail" v-loading="loading">
      <div v-if="order">
        <!-- 订单状态 -->
        <div class="order-status-section">
          <div class="status-info">
            <el-icon size="48" :color="getStatusColor(order.status)">
              <component :is="getStatusIcon(order.status)" />
            </el-icon>
            <div>
              <h2 class="status-text">{{ statusMap[order.status]?.text || order.status }}</h2>
              <p class="status-desc" v-if="order.status === 'pending_payment'">
                请在30分钟内完成支付，超时订单将自动取消
              </p>
            </div>
          </div>
          <div class="status-actions">
            <template v-if="order.status === 'pending_payment'">
              <el-button type="danger" size="large" @click="handlePayOrder">
                立即支付
              </el-button>
              <el-button size="large" @click="handleCancelOrder">
                取消订单
              </el-button>
            </template>
            <template v-if="order.status === 'shipped'">
              <el-button type="success" size="large" @click="handleConfirmReceive">
                确认收货
              </el-button>
            </template>
          </div>
        </div>
        
        <!-- 订单信息 -->
        <div class="order-info-section">
          <h3 class="section-title">订单信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">
              {{ order.order_no }}
            </el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ new Date(order.created_at).toLocaleString() }}
            </el-descriptions-item>
            <el-descriptions-item label="支付方式">
              {{ order.payment_method || '待支付' }}
            </el-descriptions-item>
            <el-descriptions-item label="支付时间" v-if="order.paid_at">
              {{ new Date(order.paid_at).toLocaleString() }}
            </el-descriptions-item>
            <el-descriptions-item label="发货时间" v-if="order.shipped_at">
              {{ new Date(order.shipped_at).toLocaleString() }}
            </el-descriptions-item>
            <el-descriptions-item label="完成时间" v-if="order.completed_at">
              {{ new Date(order.completed_at).toLocaleString() }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
        
        <!-- 收货地址 -->
        <div class="address-section">
          <h3 class="section-title">收货地址</h3>
          <div class="address-info">
            <div class="address-header">
              <span class="name">{{ order.address_name }}</span>
              <span class="phone">{{ order.address_phone }}</span>
            </div>
            <div class="address-detail">
              {{ order.address_province }}{{ order.address_city }}{{ order.address_district }}{{ order.address_detail }}
            </div>
          </div>
        </div>
        
        <!-- 商品列表 -->
        <div class="items-section">
          <h3 class="section-title">商品信息</h3>
          <div class="items-list">
            <div 
              v-for="item in order.items" 
              :key="item.id" 
              class="item-card"
            >
              <el-image 
                :src="item.vehicle?.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
                fit="cover"
                style="width: 100px; height: 100px;"
              />
              <div class="item-info">
                <h4 class="item-name">{{ item.vehicle?.name }}</h4>
                <div class="item-price">
                  <span>¥{{ item.price }}</span>
                  <span>x {{ item.quantity }}</span>
                  <span class="subtotal">小计：¥{{ (item.price * item.quantity).toFixed(2) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- 物流信息 -->
        <div class="shipping-section" v-if="order.logistics_company || order.logistics_no">
          <h3 class="section-title">物流信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="物流公司" v-if="order.logistics_company">
              {{ order.logistics_company }}
            </el-descriptions-item>
            <el-descriptions-item label="物流单号" v-if="order.logistics_no">
              {{ order.logistics_no }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
        
        <!-- 订单金额 -->
        <div class="amount-section">
          <h3 class="section-title">订单金额</h3>
          <div class="amount-summary">
            <div class="amount-row">
              <span class="label">商品金额：</span>
              <span class="value">¥{{ order.subtotal_amount?.toFixed(2) }}</span>
            </div>
            <div class="amount-row" v-if="order.shipping_fee > 0">
              <span class="label">运费：</span>
              <span class="value">¥{{ order.shipping_fee?.toFixed(2) }}</span>
            </div>
            <div class="amount-row" v-if="order.discount_amount > 0">
              <span class="label">优惠：</span>
              <span class="value discount">-¥{{ order.discount_amount?.toFixed(2) }}</span>
            </div>
            <div class="amount-row total">
              <span class="label">实付金额：</span>
              <span class="value">¥{{ order.total_amount?.toFixed(2) }}</span>
            </div>
          </div>
        </div>
        
        <!-- 订单备注 -->
        <div class="remark-section" v-if="order.remark">
          <h3 class="section-title">订单备注</h3>
          <div class="remark-content">{{ order.remark }}</div>
        </div>
      </div>
      
      <el-empty v-if="!order && !loading" description="订单不存在" />
    </div>
    
    <div class="action-bar">
      <el-button @click="router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>
  </div>
</template>

<script>
const getStatusColor = (status) => {
  const colors = {
    pending_payment: '#e6a23c',
    pending_shipment: '#409EFF',
    shipped: '#409EFF',
    completed: '#67c23a',
    cancelled: '#f56c6c'
  }
  return colors[status] || '#909399'
}

const getStatusIcon = (status) => {
  const icons = {
    pending_payment: 'Wallet',
    pending_shipment: 'Box',
    shipped: 'Truck',
    completed: 'CircleCheck',
    cancelled: 'CircleClose'
  }
  return icons[status] || 'Document'
}

export default {
  methods: {
    getStatusColor,
    getStatusIcon
  }
}
</script>

<style scoped>
.order-detail-container {
  padding: 20px 0;
}

.breadcrumb {
  margin-bottom: 20px;
}

.breadcrumb a {
  color: #606266;
  text-decoration: none;
}

.breadcrumb a:hover {
  color: #409EFF;
}

.order-detail {
  background: #fff;
  border-radius: 8px;
  padding: 30px;
}

.order-status-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  margin-bottom: 20px;
}

.status-info {
  display: flex;
  align-items: center;
  color: #fff;
}

.status-text {
  font-size: 24px;
  margin: 0 0 5px 20px;
}

.status-desc {
  font-size: 14px;
  margin: 0 0 0 20px;
  opacity: 0.8;
}

.section-title {
  font-size: 16px;
  color: #303133;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
}

.order-info-section,
.address-section,
.items-section,
.shipping-section,
.amount-section,
.remark-section {
  margin-bottom: 25px;
}

.address-info {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
}

.address-header {
  margin-bottom: 10px;
}

.name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-right: 20px;
}

.phone {
  color: #606266;
}

.address-detail {
  color: #606266;
}

.items-list {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
}

.item-card {
  display: flex;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #ebeef5;
}

.item-card:last-child {
  border-bottom: none;
}

.item-info {
  margin-left: 20px;
  flex: 1;
}

.item-name {
  font-size: 16px;
  color: #303133;
  margin-bottom: 10px;
}

.item-price {
  font-size: 14px;
  color: #606266;
}

.item-price .subtotal {
  float: right;
  color: #f56c6c;
  font-weight: bold;
}

.amount-summary {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
}

.amount-row {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  margin-bottom: 15px;
}

.amount-row:last-child {
  margin-bottom: 0;
}

.amount-row .label {
  color: #606266;
  font-size: 14px;
}

.amount-row .value {
  color: #303133;
  font-size: 14px;
  min-width: 100px;
  text-align: right;
}

.amount-row.total .label {
  font-size: 16px;
  font-weight: 500;
}

.amount-row.total .value {
  font-size: 20px;
  color: #f56c6c;
  font-weight: bold;
}

.amount-row .value.discount {
  color: #67c23a;
}

.remark-content {
  background: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  color: #606266;
}

.action-bar {
  margin-top: 20px;
  text-align: center;
}
</style>
