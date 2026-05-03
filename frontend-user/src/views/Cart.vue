<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const cartItems = ref([])
const loading = ref(false)
const allSelected = ref(false)
const selectAllLoading = ref(false)

// 获取购物车列表
const fetchCartList = async () => {
  loading.value = true
  try {
    const res = await request.get('/cart/list')
    cartItems.value = (res.data || []).map(item => ({
      ...item,
      selected: true
    }))
    updateAllSelected()
  } catch (error) {
    console.error('获取购物车失败:', error)
  } finally {
    loading.value = false
  }
}

// 计算选中的商品
const selectedItems = computed(() => {
  return cartItems.value.filter(item => item.selected)
})

// 计算总价格
const totalPrice = computed(() => {
  return selectedItems.value.reduce((total, item) => {
    return total + (item.vehicle?.price || 0) * item.quantity
  }, 0)
})

// 计算总数量
const totalQuantity = computed(() => {
  return selectedItems.value.reduce((total, item) => {
    return total + item.quantity
  }, 0)
})

// 更新全选状态
const updateAllSelected = () => {
  if (cartItems.value.length === 0) {
    allSelected.value = false
  } else {
    allSelected.value = cartItems.value.every(item => item.selected)
  }
}

// 全选/取消全选
const handleSelectAll = () => {
  selectAllLoading.value = true
  const newSelected = !allSelected.value
  cartItems.value.forEach(item => {
    item.selected = newSelected
  })
  allSelected.value = newSelected
  setTimeout(() => {
    selectAllLoading.value = false
  }, 100)
}

// 更新数量
const handleQuantityChange = async (item, newQuantity) => {
  try {
    await request.post('/cart/update', {
      id: item.id,
      quantity: newQuantity
    })
    item.quantity = newQuantity
  } catch (error) {
    console.error('更新数量失败:', error)
    ElMessage.error('更新失败')
  }
}

// 删除购物车项
const handleRemoveItem = async (item) => {
  try {
    await ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post('/cart/delete', { id: item.id })
    const index = cartItems.value.indexOf(item)
    if (index > -1) {
      cartItems.value.splice(index, 1)
    }
    updateAllSelected()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const handleRemoveSelected = async () => {
  if (selectedItems.value.length === 0) {
    ElMessage.warning('请选择要删除的商品')
    return
  }
  
  try {
    await ElMessageBox.confirm('确定要删除选中的商品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const ids = selectedItems.value.map(item => item.id)
    await request.post('/cart/delete', { ids })
    cartItems.value = cartItems.value.filter(item => !item.selected)
    updateAllSelected()
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 结算
const handleCheckout = () => {
  if (selectedItems.value.length === 0) {
    ElMessage.warning('请选择要结算的商品')
    return
  }
  router.push('/order/create')
}

onMounted(() => {
  fetchCartList()
})
</script>

<template>
  <div class="cart-container">
    <h2 class="page-title">购物车</h2>
    
    <!-- 购物车内容 -->
    <div class="cart-content">
      <el-table
        v-loading="loading"
        :data="cartItems"
        style="width: 100%"
      >
        <el-table-column width="50">
          <template #header>
            <el-checkbox 
              v-model="allSelected" 
              :loading="selectAllLoading"
              @change="handleSelectAll"
            />
          </template>
          <template #default="scope">
            <el-checkbox v-model="scope.row.selected" @change="updateAllSelected" />
          </template>
        </el-table-column>
        
        <el-table-column label="商品信息" min-width="300">
          <template #default="scope">
            <div class="product-info">
              <el-image 
                :src="scope.row.vehicle?.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
                fit="cover"
                style="width: 100px; height: 100px;"
              />
              <div class="product-detail">
                <h4 class="product-name">{{ scope.row.vehicle?.name }}</h4>
                <p class="product-desc">{{ scope.row.vehicle?.description }}</p>
              </div>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column label="单价" width="120">
          <template #default="scope">
            <span class="price">¥{{ scope.row.vehicle?.price }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="数量" width="150">
          <template #default="scope">
            <el-input-number
              v-model="scope.row.quantity"
              :min="1"
              :max="scope.row.vehicle?.stock || 100"
              size="small"
              @change="handleQuantityChange(scope.row, $event)"
            />
          </template>
        </el-table-column>
        
        <el-table-column label="小计" width="120">
          <template #default="scope">
            <span class="subtotal">¥{{ (scope.row.vehicle?.price * scope.row.quantity).toFixed(2) }}</span>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="80">
          <template #default="scope">
            <el-button 
              type="text" 
              danger
              @click="handleRemoveItem(scope.row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 空状态 -->
      <el-empty 
        v-if="cartItems.length === 0 && !loading" 
        description="购物车是空的"
      >
        <el-button type="primary" @click="router.push('/vehicle')">
          去逛逛
        </el-button>
      </el-empty>
    </div>
    
    <!-- 结算栏 -->
    <div class="cart-footer" v-if="cartItems.length > 0">
      <div class="footer-left">
        <el-checkbox 
          v-model="allSelected" 
          :loading="selectAllLoading"
          @change="handleSelectAll"
        >
          全选
        </el-checkbox>
        <el-button 
          type="text" 
          danger
          @click="handleRemoveSelected"
          :disabled="selectedItems.length === 0"
        >
          删除选中
        </el-button>
      </div>
      
      <div class="footer-right">
        <div class="summary">
          <span>已选 <span class="highlight">{{ totalQuantity }}</span> 件商品</span>
          <span style="margin-left: 30px;">
            合计：<span class="total-price">¥{{ totalPrice.toFixed(2) }}</span>
          </span>
        </div>
        <el-button 
          type="danger" 
          size="large"
          @click="handleCheckout"
          :disabled="selectedItems.length === 0"
        >
          结算
        </el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.cart-container {
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

.cart-content {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
}

.product-info {
  display: flex;
  align-items: center;
}

.product-detail {
  margin-left: 15px;
  flex: 1;
}

.product-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 5px;
  line-height: 1.4;
}

.product-desc {
  font-size: 12px;
  color: #909399;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.price {
  color: #f56c6c;
  font-weight: bold;
}

.subtotal {
  color: #f56c6c;
  font-weight: bold;
  font-size: 16px;
}

.cart-footer {
  background: #fff;
  border-radius: 8px;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  bottom: 0;
  z-index: 10;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.05);
}

.footer-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 30px;
}

.summary {
  font-size: 14px;
  color: #606266;
}

.highlight {
  color: #f56c6c;
  font-weight: bold;
}

.total-price {
  color: #f56c6c;
  font-size: 20px;
  font-weight: bold;
}
</style>
