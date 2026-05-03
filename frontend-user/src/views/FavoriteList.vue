<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const favorites = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 获取收藏列表
const fetchFavorites = async () => {
  loading.value = true
  try {
    const res = await request.get('/favorite/list', {
      params: {
        page: page.value,
        page_size: pageSize.value
      }
    })
    favorites.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取收藏列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 取消收藏
const handleRemoveFavorite = async (favorite) => {
  try {
    await ElMessageBox.confirm('确定要取消收藏吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post('/favorite/delete', { id: favorite.id })
    const index = favorites.value.indexOf(favorite)
    if (index > -1) {
      favorites.value.splice(index, 1)
    }
    total.value--
    ElMessage.success('已取消收藏')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消收藏失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 跳转商品详情
const handleVehicleClick = (vehicle) => {
  if (vehicle) {
    router.push(`/vehicle/${vehicle.id}`)
  }
}

// 分页
const handlePageChange = (val) => {
  page.value = val
  fetchFavorites()
}

onMounted(() => {
  fetchFavorites()
})
</script>

<template>
  <div class="favorite-list-container">
    <h2 class="page-title">我的收藏</h2>
    
    <!-- 收藏列表 -->
    <div class="favorite-list" v-loading="loading">
      <el-row :gutter="20">
        <el-col :xs="12" :sm="8" :md="6" v-for="favorite in favorites" :key="favorite.id">
          <div class="favorite-card">
            <div class="vehicle-image" @click="handleVehicleClick(favorite.vehicle)">
              <el-image 
                :src="favorite.vehicle?.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
                fit="cover"
                style="width: 100%; height: 200px;"
              />
            </div>
            <div class="vehicle-info">
              <h4 class="vehicle-name" @click="handleVehicleClick(favorite.vehicle)">
                {{ favorite.vehicle?.name }}
              </h4>
              <div class="vehicle-price">
                <span class="price">¥{{ favorite.vehicle?.price }}</span>
                <span v-if="favorite.vehicle?.original_price > favorite.vehicle?.price" class="original-price">
                  ¥{{ favorite.vehicle?.original_price }}
                </span>
              </div>
              <div class="vehicle-action">
                <el-button 
                  type="danger" 
                  size="small"
                  @click="handleRemoveFavorite(favorite)"
                >
                  <el-icon><Delete /></el-icon>
                  取消收藏
                </el-button>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
      
      <!-- 空状态 -->
      <el-empty 
        v-if="favorites.length === 0 && !loading" 
        description="暂无收藏商品"
      >
        <el-button type="primary" @click="router.push('/vehicle')">
          去逛逛
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
.favorite-list-container {
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

.favorite-list {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}

.favorite-card {
  background: #fff;
  border: 1px solid #ebeef5;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: 20px;
  transition: all 0.3s;
}

.favorite-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.vehicle-image {
  position: relative;
  height: 200px;
  background: #f5f7fa;
  cursor: pointer;
}

.vehicle-info {
  padding: 15px;
}

.vehicle-name {
  font-size: 14px;
  color: #303133;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  cursor: pointer;
}

.vehicle-price {
  margin-bottom: 10px;
}

.vehicle-price .price {
  font-size: 16px;
  color: #f56c6c;
  font-weight: bold;
}

.vehicle-price .original-price {
  font-size: 12px;
  color: #909399;
  text-decoration: line-through;
  margin-left: 10px;
}

.vehicle-action {
  display: flex;
  justify-content: flex-end;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
