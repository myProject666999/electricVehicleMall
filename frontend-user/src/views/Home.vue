<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const router = useRouter()

const banners = ref([])
const categories = ref([])
const recommendVehicles = ref([])
const hotVehicles = ref([])
const notices = ref([])
const loading = ref(true)

const fetchBanners = async () => {
  try {
    const res = await request.get('/banner/list', { params: { status: 1 } })
    banners.value = res.data || []
  } catch (error) {
    console.error('获取轮播图失败:', error)
  }
}

const fetchCategories = async () => {
  try {
    const res = await request.get('/category/list', { params: { status: 1 } })
    categories.value = res.data || []
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

const fetchVehicles = async () => {
  try {
    // 获取推荐电动车
    const recommendRes = await request.get('/vehicle/list', { 
      params: { is_recommend: 1, status: 1, page_size: 8 } 
    })
    recommendVehicles.value = recommendRes.data?.list || []
    
    // 获取热销电动车
    const hotRes = await request.get('/vehicle/list', { 
      params: { status: 1, page_size: 8 } 
    })
    hotVehicles.value = hotRes.data?.list || []
  } catch (error) {
    console.error('获取电动车失败:', error)
  }
}

const fetchNotices = async () => {
  try {
    const res = await request.get('/notice/list', { 
      params: { status: 1, page_size: 5 } 
    })
    notices.value = res.data?.list || []
  } catch (error) {
    console.error('获取公告失败:', error)
  }
}

const handleVehicleClick = (vehicle) => {
  router.push(`/vehicle/${vehicle.id}`)
}

const handleCategoryClick = (category) => {
  router.push(`/vehicle?category_id=${category.id}`)
}

const handleNoticeClick = (notice) => {
  router.push(`/notice/${notice.id}`)
}

onMounted(async () => {
  try {
    await Promise.all([
      fetchBanners(),
      fetchCategories(),
      fetchVehicles(),
      fetchNotices()
    ])
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="home-container">
    <!-- 轮播图 -->
    <div class="banner-section">
      <el-carousel height="400px" indicator-position="outside">
        <el-carousel-item v-for="banner in banners" :key="banner.id">
          <div class="banner-item" :style="{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }">
            <div class="banner-content">
              <h2>{{ banner.title }}</h2>
              <p v-if="banner.image" class="banner-desc">点击查看详情</p>
              <el-button type="primary" size="large" v-if="banner.vehicle_id" @click="handleVehicleClick({ id: banner.vehicle_id })">
                立即查看
              </el-button>
            </div>
          </div>
        </el-carousel-item>
        <!-- 默认轮播图 -->
        <el-carousel-item v-if="banners.length === 0">
          <div class="banner-item" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%)">
            <div class="banner-content">
              <h2>欢迎来到电动车商城</h2>
              <p>精选优质电动车，品质保障，售后无忧</p>
              <el-button type="primary" size="large" @click="router.push('/vehicle')">
                立即选购
              </el-button>
            </div>
          </div>
        </el-carousel-item>
      </el-carousel>
    </div>

    <!-- 分类导航 -->
    <div class="category-section">
      <div class="section-title">
        <h3>商品分类</h3>
      </div>
      <div class="category-list">
        <div 
          v-for="category in categories" 
          :key="category.id" 
          class="category-item"
          @click="handleCategoryClick(category)"
        >
          <el-icon size="40" color="#409EFF">
            <component :is="getCategoryIcon(category.name)" />
          </el-icon>
          <span>{{ category.name }}</span>
        </div>
      </div>
    </div>

    <!-- 推荐电动车 -->
    <div class="vehicle-section">
      <div class="section-title">
        <h3>推荐商品</h3>
        <router-link to="/vehicle" class="more-link">查看更多 →</router-link>
      </div>
      <el-row :gutter="20">
        <el-col :xs="12" :sm="8" :md="6" v-for="vehicle in recommendVehicles" :key="vehicle.id">
          <div class="vehicle-card" @click="handleVehicleClick(vehicle)">
            <div class="vehicle-image">
              <el-image 
                :src="vehicle.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
                fit="cover"
                style="width: 100%; height: 200px;"
              />
              <div v-if="vehicle.is_recommend" class="tag-recommend">推荐</div>
            </div>
            <div class="vehicle-info">
              <h4 class="vehicle-name">{{ vehicle.name }}</h4>
              <div class="vehicle-price">
                <span class="price">¥{{ vehicle.price }}</span>
                <span v-if="vehicle.original_price > vehicle.price" class="original-price">¥{{ vehicle.original_price }}</span>
              </div>
              <div class="vehicle-meta">
                <span>销量: {{ vehicle.sales }}</span>
                <span>库存: {{ vehicle.stock }}</span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 热销电动车 -->
    <div class="vehicle-section" v-if="hotVehicles.length > 0">
      <div class="section-title">
        <h3>热销商品</h3>
        <router-link to="/vehicle" class="more-link">查看更多 →</router-link>
      </div>
      <el-row :gutter="20">
        <el-col :xs="12" :sm="8" :md="6" v-for="vehicle in hotVehicles" :key="vehicle.id">
          <div class="vehicle-card" @click="handleVehicleClick(vehicle)">
            <div class="vehicle-image">
              <el-image 
                :src="vehicle.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle%20product%20photo&image_size=square'"
                fit="cover"
                style="width: 100%; height: 200px;"
              />
            </div>
            <div class="vehicle-info">
              <h4 class="vehicle-name">{{ vehicle.name }}</h4>
              <div class="vehicle-price">
                <span class="price">¥{{ vehicle.price }}</span>
              </div>
              <div class="vehicle-meta">
                <span>销量: {{ vehicle.sales }}</span>
              </div>
            </div>
          </div>
        </el-col>
      </el-row>
    </div>

    <!-- 公告列表 -->
    <div class="notice-section">
      <div class="section-title">
        <h3>最新公告</h3>
        <router-link to="/notice" class="more-link">查看更多 →</router-link>
      </div>
      <el-card>
        <ul class="notice-list">
          <li v-for="notice in notices" :key="notice.id" @click="handleNoticeClick(notice)">
            <span class="notice-title">{{ notice.title }}</span>
            <span class="notice-date">{{ new Date(notice.created_at).toLocaleDateString() }}</span>
          </li>
          <li v-if="notices.length === 0">
            <el-empty description="暂无公告" :image-size="60" />
          </li>
        </ul>
      </el-card>
    </div>
  </div>
</template>

<script>
const getCategoryIcon = (name) => {
  const icons = {
    '电动自行车': 'Bicycle',
    '电动摩托车': 'Motorcycle',
    '电动三轮车': 'Truck',
    '电动四轮车': 'Van'
  }
  return icons[name] || 'ElectricVehicle'
}

export default {
  methods: {
    getCategoryIcon
  }
}
</script>

<style scoped>
.home-container {
  padding-bottom: 40px;
}

/* 轮播图 */
.banner-section {
  margin-bottom: 30px;
}

.banner-item {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
  position: relative;
}

.banner-content {
  text-align: center;
  color: #fff;
  z-index: 10;
}

.banner-content h2 {
  font-size: 36px;
  margin-bottom: 15px;
}

.banner-content p {
  font-size: 18px;
  margin-bottom: 25px;
  opacity: 0.9;
}

/* 分类导航 */
.category-section {
  margin-bottom: 30px;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}

.section-title {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-title h3 {
  font-size: 20px;
  color: #303133;
  border-left: 4px solid #409EFF;
  padding-left: 12px;
}

.more-link {
  color: #909399;
  text-decoration: none;
  font-size: 14px;
}

.more-link:hover {
  color: #409EFF;
}

.category-list {
  display: flex;
  justify-content: space-around;
  flex-wrap: wrap;
}

.category-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 20px 30px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 8px;
}

.category-item:hover {
  background: #f5f7fa;
  transform: translateY(-5px);
}

.category-item span {
  margin-top: 10px;
  color: #606266;
  font-size: 14px;
}

/* 电动车卡片 */
.vehicle-section {
  margin-bottom: 30px;
}

.vehicle-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  margin-bottom: 20px;
}

.vehicle-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-5px);
}

.vehicle-image {
  position: relative;
  height: 200px;
  background: #f5f7fa;
}

.tag-recommend {
  position: absolute;
  top: 10px;
  left: 10px;
  background: #f56c6c;
  color: #fff;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
}

.vehicle-info {
  padding: 15px;
}

.vehicle-name {
  font-size: 16px;
  color: #303133;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.vehicle-price {
  margin-bottom: 10px;
}

.vehicle-price .price {
  font-size: 18px;
  color: #f56c6c;
  font-weight: bold;
}

.vehicle-price .original-price {
  font-size: 14px;
  color: #909399;
  text-decoration: line-through;
  margin-left: 10px;
}

.vehicle-meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
}

/* 公告列表 */
.notice-section {
  margin-bottom: 30px;
}

.notice-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.notice-list li {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: color 0.3s;
}

.notice-list li:last-child {
  border-bottom: none;
}

.notice-list li:hover .notice-title {
  color: #409EFF;
}

.notice-title {
  color: #303133;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.notice-date {
  color: #909399;
  font-size: 12px;
  margin-left: 20px;
}
</style>
