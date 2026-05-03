<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const addresses = ref([])
const selectedAddress = ref(null)
const cartItems = ref([])
const selectedCartIds = ref([])
const remark = ref('')
const loading = ref(false)
const submitting = ref(false)

// 获取地址列表
const fetchAddresses = async () => {
  try {
    const res = await request.get('/address/list')
    addresses.value = res.data || []
    // 默认选中默认地址
    selectedAddress.value = addresses.value.find(addr => addr.is_default) || addresses.value[0] || null
  } catch (error) {
    console.error('获取地址列表失败:', error)
  }
}

// 获取购物车列表
const fetchCartItems = async () => {
  loading.value = true
  try {
    const res = await request.get('/cart/list')
    cartItems.value = (res.data || []).filter(item => item.selected)
    selectedCartIds.value = cartItems.value.map(item => item.id)
  } catch (error) {
    console.error('获取购物车失败:', error)
  } finally {
    loading.value = false
  }
}

// 计算总价格
const totalAmount = computed(() => {
  return cartItems.value.reduce((total, item) => {
    return total + (item.vehicle?.price || 0) * item.quantity
  }, 0)
})

// 计算总数量
const totalQuantity = computed(() => {
  return cartItems.value.reduce((total, item) => {
    return total + item.quantity
  }, 0)
})

// 提交订单
const handleSubmit = async () => {
  if (!selectedAddress.value) {
    ElMessage.warning('请选择收货地址')
    return
  }
  
  if (cartItems.value.length === 0) {
    ElMessage.warning('购物车为空')
    return
  }
  
  submitting.value = true
  try {
    const res = await request.post('/order/create', {
      address_id: selectedAddress.value.id,
      cart_ids: selectedCartIds.value,
      remark: remark.value
    })
    
    ElMessage.success('订单创建成功')
    router.push(`/order/${res.data.id}`)
  } catch (error) {
    console.error('创建订单失败:', error)
    ElMessage.error('创建订单失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchAddresses()
  fetchCartItems()
})
</script>

<template>
  <div class="create-order-container">
    <h2 class="page-title">确认订单</h2>
    
    <div class="order-content" v-loading="loading">
      <!-- 收货地址 -->
      <div class="section">
        <h3 class="section-title">收货地址</h3>
        <div class="address-list">
          <div 
            v-for="address in addresses" 
            :key="address.id" 
            class="address-card"
            :class="{ selected: selectedAddress?.id === address.id }"
            @click="selectedAddress = address"
          >
            <div class="address-header">
              <span class="name">{{ address.name }}</span>
              <span class="phone">{{ address.phone }}</span>
              <el-tag v-if="address.is_default" type="success" size="small">默认</el-tag>
            </div>
            <div class="address-detail">
              {{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}
            </div>
          </div>
          
          <div class="no-address" v-if="addresses.length === 0">
            <el-empty description="暂无收货地址" :image-size="60">
              <el-button type="primary" @click="router.push('/address')">
                添加收货地址
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>
      
      <!-- 商品列表 -->
      <div class="section">
        <h3 class="section-title">商品清单</h3>
        <div class="items-list">
          <div 
            v-for="item in cartItems" 
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
                <span class="price">¥{{ item.vehicle?.price }}</span>
                <span class="quantity">x {{ item.quantity }}</span>
              </div>
            </div>
            <div class="item-subtotal">
              ¥{{ (item.vehicle?.price * item.quantity).toFixed(2) }}
            </div>
          </div>
          
          <div class="no-items" v-if="cartItems.length === 0">
            <el-empty description="购物车为空" :image-size="60">
              <el-button type="primary" @click="router.push('/vehicle')">
                去购物
              </el-button>
            </el-empty>
          </div>
        </div>
      </div>
      
      <!-- 订单备注 -->
      <div class="section">
        <h3 class="section-title">订单备注</h3>
        <el-input
          v-model="remark"
          type="textarea"
          :rows="3"
          placeholder="选填，请输入订单备注信息..."
          maxlength="200"
          show-word-limit
        />
      </div>
      
      <!-- 订单金额 -->
      <div class="section amount-section">
        <div class="amount-summary">
          <div class="amount-row">
            <span class="label">商品金额：</span>
            <span class="value">¥{{ totalAmount.toFixed(2) }}</span>
          </div>
          <div class="amount-row">
            <span class="label">商品数量：</span>
            <span class="value">{{ totalQuantity }} 件</span>
          </div>
          <div class="amount-row total">
            <span class="label">实付金额：</span>
            <span class="value">¥{{ totalAmount.toFixed(2) }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- 底部操作栏 -->
    <div class="action-bar">
      <div class="amount-info">
        <span>共 <span class="highlight">{{ totalQuantity }}</span> 件商品</span>
        <span style="margin-left: 30px;">
          合计：<span class="total-price">¥{{ totalAmount.toFixed(2) }}</span>
        </span>
      </div>
      <el-button
        type="danger"
        size="large"
        :loading="submitting"
        :disabled="cartItems.length === 0 || !selectedAddress"
        @click="handleSubmit"
      >
        提交订单
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.create-order-container {
  padding: 20px 0;
  padding-bottom: 80px;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  border-bottom: 2px solid #409EFF;
  padding-bottom: 10px;
  display: inline-block;
}

.order-content {
  background: #fff;
  border-radius: 8px;
  padding: 30px;
  margin-bottom: 20px;
}

.section {
  margin-bottom: 30px;
}

.section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 16px;
  color: #303133;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 1px solid #ebeef5;
}

.address-list {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
}

.address-card {
  flex: 1;
  min-width: 280px;
  max-width: calc(50% - 10px);
  padding: 20px;
  border: 2px solid #ebeef5;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.address-card:hover {
  border-color: #409EFF;
}

.address-card.selected {
  border-color: #409EFF;
  background: #ecf5ff;
}

.address-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-right: 15px;
}

.phone {
  color: #606266;
  margin-right: 15px;
}

.address-detail {
  color: #606266;
  font-size: 14px;
}

.no-address {
  width: 100%;
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

.item-price .price {
  color: #f56c6c;
  font-weight: bold;
  margin-right: 20px;
}

.item-subtotal {
  font-size: 16px;
  color: #f56c6c;
  font-weight: bold;
  min-width: 120px;
  text-align: right;
}

.no-items {
  width: 100%;
}

.amount-section {
  background: #f5f7fa;
  border-radius: 8px;
  padding: 20px;
}

.amount-summary {
  max-width: 400px;
  margin-left: auto;
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
  min-width: 120px;
  text-align: right;
}

.amount-row.total .label {
  font-size: 18px;
  font-weight: 500;
}

.amount-row.total .value {
  font-size: 24px;
  color: #f56c6c;
  font-weight: bold;
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  z-index: 100;
}

.amount-info {
  font-size: 16px;
  color: #606266;
}

.highlight {
  color: #f56c6c;
  font-weight: bold;
}

.total-price {
  color: #f56c6c;
  font-size: 24px;
  font-weight: bold;
}
</style>
