<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()

const vehicle = ref(null)
const comments = ref([])
const loading = ref(true)
const quantity = ref(1)
const activeTab = ref('info')
const commentPage = ref(1)
const commentTotal = ref(0)

// 获取电动车详情
const fetchVehicleDetail = async () => {
  try {
    const res = await request.get(`/vehicle/detail/${route.params.id}`)
    vehicle.value = res.data
  } catch (error) {
    console.error('获取电动车详情失败:', error)
    ElMessage.error('获取商品详情失败')
  }
}

// 获取评论列表
const fetchComments = async () => {
  try {
    const res = await request.get('/comment/list', {
      params: {
        vehicle_id: route.params.id,
        page: commentPage.value,
        page_size: 10
      }
    })
    comments.value = res.data?.list || []
    commentTotal.value = res.data?.total || 0
  } catch (error) {
    console.error('获取评论列表失败:', error)
  }
}

// 添加购物车
const handleAddToCart = async () => {
  try {
    await request.post('/cart/add', {
      vehicle_id: route.params.id,
      quantity: quantity.value
    })
    ElMessage.success('已添加到购物车')
  } catch (error) {
    if (error.response?.status === 401) {
      ElMessage.warning('请先登录')
      router.push('/login')
      return
    }
    console.error('添加购物车失败:', error)
    ElMessage.error('添加购物车失败')
  }
}

// 立即购买
const handleBuyNow = () => {
  handleAddToCart().then(() => {
    router.push('/cart')
  })
}

// 添加收藏
const handleAddFavorite = async () => {
  try {
    await request.post('/favorite/add', {
      vehicle_id: route.params.id
    })
    ElMessage.success('收藏成功')
  } catch (error) {
    if (error.response?.status === 401) {
      ElMessage.warning('请先登录')
      router.push('/login')
      return
    }
    console.error('添加收藏失败:', error)
    ElMessage.error('收藏失败')
  }
}

onMounted(async () => {
  await Promise.all([
    fetchVehicleDetail(),
    fetchComments()
  ])
  loading.value = false
})
</script>

<template>
  <div class="vehicle-detail-container" v-loading="loading">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item>
        <router-link to="/">首页</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>
        <router-link to="/vehicle">电动车列表</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>{{ vehicle?.name || '商品详情' }}</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 商品信息 -->
    <div class="vehicle-detail">
      <el-row :gutter="40">
        <!-- 商品图片 -->
        <el-col :md="10">
          <div class="vehicle-image-box">
            <el-image 
              :src="vehicle?.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
              fit="contain"
              style="width: 100%; height: 400px;"
              preview-src-list="[vehicle?.image]"
            />
          </div>
        </el-col>

        <!-- 商品信息 -->
        <el-col :md="14">
          <div class="vehicle-info-box">
            <h1 class="vehicle-name">{{ vehicle?.name }}</h1>
            
            <div class="vehicle-price-section">
              <span class="price-label">价格</span>
              <span class="price">¥{{ vehicle?.price }}</span>
              <span v-if="vehicle?.original_price > vehicle?.price" class="original-price">
                ¥{{ vehicle?.original_price }}
              </span>
              <span v-if="vehicle?.original_price > vehicle?.price" class="discount">
                {{ Math.round((1 - vehicle.price / vehicle.original_price) * 100) }}% OFF
              </span>
            </div>

            <div class="vehicle-meta">
              <div class="meta-item">
                <span class="meta-label">销量：</span>
                <span class="meta-value">{{ vehicle?.sales }}</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">库存：</span>
                <span class="meta-value">{{ vehicle?.stock }} 台</span>
              </div>
              <div class="meta-item">
                <span class="meta-label">评分：</span>
                <el-rate v-model="vehicle?.rating" disabled />
              </div>
            </div>

            <div class="vehicle-action">
              <div class="quantity-section">
                <span class="label">数量：</span>
                <el-input-number 
                  v-model="quantity" 
                  :min="1" 
                  :max="vehicle?.stock || 100"
                  size="large"
                />
              </div>

              <div class="action-buttons">
                <el-button 
                  type="danger" 
                  size="large"
                  @click="handleBuyNow"
                >
                  <el-icon><Money /></el-icon>
                  立即购买
                </el-button>
                <el-button 
                  type="primary" 
                  size="large"
                  @click="handleAddToCart"
                >
                  <el-icon><ShoppingCart /></el-icon>
                  加入购物车
                </el-button>
                <el-button 
                  size="large"
                  @click="handleAddFavorite"
                >
                  <el-icon><Heart /></el-icon>
                  收藏
                </el-button>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 标签页 -->
    <div class="detail-tabs">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="商品详情" name="info">
          <div class="tab-content">
            <div class="product-description" v-if="vehicle?.description">
              <h3>商品描述</h3>
              <div v-html="vehicle.description"></div>
            </div>
            <div class="product-specs">
              <h3>商品参数</h3>
              <el-table :data="vehicle?.specs || []" style="width: 100%">
                <el-table-column property="name" label="参数名称" />
                <el-table-column property="value" label="参数值" />
              </el-table>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane :label="`评价 (${commentTotal})`" name="comments">
          <div class="tab-content">
            <div class="comments-section">
              <div v-if="comments.length > 0">
                <div class="comment-item" v-for="comment in comments" :key="comment.id">
                  <div class="comment-header">
                    <el-avatar :size="40" icon="User" />
                    <div class="user-info">
                      <span class="username">{{ comment.user_name || '匿名用户' }}</span>
                      <el-rate v-model="comment.rating" disabled />
                    </div>
                    <span class="comment-time">{{ new Date(comment.created_at).toLocaleString() }}</span>
                  </div>
                  <div class="comment-content">{{ comment.content }}</div>
                  <div v-if="comment.images" class="comment-images">
                    <el-image 
                      v-for="(img, index) in comment.images.split(',')" 
                      :key="index"
                      :src="img"
                      :preview-src-list="comment.images.split(',')"
                      fit="cover"
                      style="width: 100px; height: 100px; margin-right: 10px;"
                    />
                  </div>
                  <div v-if="comment.admin_reply" class="admin-reply">
                    <strong>商家回复：</strong>
                    <span>{{ comment.admin_reply }}</span>
                  </div>
                </div>
              </div>
              <el-empty v-else description="暂无评价" />
            </div>
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<style scoped>
.vehicle-detail-container {
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

.vehicle-detail {
  background: #fff;
  padding: 30px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.vehicle-image-box {
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f7fa;
}

.vehicle-info-box {
  padding: 0 20px;
}

.vehicle-name {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  line-height: 1.4;
}

.vehicle-price-section {
  background: #fff5f5;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.price-label {
  color: #909399;
  margin-right: 10px;
}

.price {
  font-size: 28px;
  color: #f56c6c;
  font-weight: bold;
}

.original-price {
  font-size: 16px;
  color: #909399;
  text-decoration: line-through;
  margin-left: 15px;
}

.discount {
  background: #f56c6c;
  color: #fff;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  margin-left: 10px;
}

.vehicle-meta {
  margin-bottom: 20px;
}

.meta-item {
  margin-bottom: 10px;
}

.meta-label {
  color: #909399;
  width: 60px;
  display: inline-block;
}

.meta-value {
  color: #303133;
}

.vehicle-action {
  border-top: 1px solid #ebeef5;
  padding-top: 20px;
}

.quantity-section {
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}

.quantity-section .label {
  color: #909399;
  margin-right: 15px;
}

.action-buttons {
  display: flex;
  gap: 15px;
}

.action-buttons .el-button {
  padding: 12px 30px;
}

.detail-tabs {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
}

.tab-content {
  padding: 20px;
}

.product-description,
.product-specs {
  margin-bottom: 30px;
}

.product-description h3,
.product-specs h3 {
  font-size: 18px;
  color: #303133;
  margin-bottom: 15px;
  border-left: 4px solid #409EFF;
  padding-left: 10px;
}

.comments-section {
  min-height: 200px;
}

.comment-item {
  padding: 20px 0;
  border-bottom: 1px solid #ebeef5;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.user-info {
  margin-left: 15px;
  flex: 1;
}

.username {
  display: block;
  color: #303133;
  margin-bottom: 5px;
}

.comment-time {
  color: #909399;
  font-size: 12px;
}

.comment-content {
  color: #606266;
  line-height: 1.6;
  margin-bottom: 10px;
}

.comment-images {
  margin-bottom: 10px;
}

.admin-reply {
  background: #f5f7fa;
  padding: 10px;
  border-radius: 4px;
  color: #606266;
}

.admin-reply strong {
  color: #409EFF;
}
</style>
