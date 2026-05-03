<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()

const categories = ref([])
const vehicles = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(12)
const total = ref(0)

// 筛选条件
const filter = ref({
  category_id: null,
  keyword: '',
  min_price: null,
  max_price: null,
  sort: 'default'
})

const fetchCategories = async () => {
  try {
    const res = await request.get('/category/list', { params: { status: 1 } })
    categories.value = [{ id: null, name: '全部' }, ...(res.data || [])]
  } catch (error) {
    console.error('获取分类失败:', error)
  }
}

const fetchVehicles = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
      status: 1
    }
    
    if (filter.value.category_id) params.category_id = filter.value.category_id
    if (filter.value.keyword) params.keyword = filter.value.keyword
    if (filter.value.min_price !== null) params.min_price = filter.value.min_price
    if (filter.value.max_price !== null) params.max_price = filter.value.max_price
    if (filter.value.sort) params.sort = filter.value.sort

    const res = await request.get('/vehicle/list', { params })
    vehicles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取电动车列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleVehicleClick = (vehicle) => {
  router.push(`/vehicle/${vehicle.id}`)
}

const handleCategoryChange = (categoryId) => {
  filter.value.category_id = categoryId
  page.value = 1
  fetchVehicles()
}

const handleSortChange = (sort) => {
  filter.value.sort = sort
  page.value = 1
  fetchVehicles()
}

const handlePriceSearch = () => {
  page.value = 1
  fetchVehicles()
}

const handlePageChange = (val) => {
  page.value = val
  fetchVehicles()
}

// 处理URL参数
const handleRouteParams = () => {
  if (route.query.category_id) {
    filter.value.category_id = Number(route.query.category_id)
  }
  if (route.query.keyword) {
    filter.value.keyword = route.query.keyword
  }
}

onMounted(async () => {
  handleRouteParams()
  await fetchCategories()
  await fetchVehicles()
})
</script>

<template>
  <div class="vehicle-list-container">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item>
        <router-link to="/">首页</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>电动车列表</el-breadcrumb-item>
    </el-breadcrumb>

    <!-- 分类筛选 -->
    <el-card class="filter-card">
      <div class="filter-section">
        <span class="filter-label">分类：</span>
        <el-radio-group v-model="filter.category_id" @change="handleCategoryChange">
          <el-radio-button v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </el-radio-button>
        </el-radio-group>
      </div>
      
      <div class="filter-section">
        <span class="filter-label">排序：</span>
        <el-radio-group v-model="filter.sort" @change="handleSortChange">
          <el-radio-button value="default">默认</el-radio-button>
          <el-radio-button value="price_asc">价格升序</el-radio-button>
          <el-radio-button value="price_desc">价格降序</el-radio-button>
          <el-radio-button value="sales_desc">销量优先</el-radio-button>
        </el-radio-group>
      </div>

      <div class="filter-section">
        <span class="filter-label">价格：</span>
        <el-input v-model.number="filter.min_price" placeholder="最低价" style="width: 100px;" />
        <span style="margin: 0 10px;">-</span>
        <el-input v-model.number="filter.max_price" placeholder="最高价" style="width: 100px;" />
        <el-button type="primary" size="small" @click="handlePriceSearch">确定</el-button>
      </div>
    </el-card>

    <!-- 商品列表 -->
    <div class="vehicle-list">
      <el-skeleton :loading="loading" :rows="4">
        <el-row :gutter="20">
          <el-col :xs="12" :sm="8" :md="6" v-for="vehicle in vehicles" :key="vehicle.id">
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

        <!-- 空状态 -->
        <el-empty v-if="vehicles.length === 0 && !loading" description="暂无商品" />

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
      </el-skeleton>
    </div>
  </div>
</template>

<style scoped>
.vehicle-list-container {
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

.filter-card {
  margin-bottom: 20px;
}

.filter-section {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.filter-section:last-child {
  margin-bottom: 0;
}

.filter-label {
  width: 60px;
  color: #909399;
  margin-right: 10px;
  flex-shrink: 0;
}

.vehicle-list {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}

.vehicle-card {
  background: #fff;
  border: 1px solid #ebeef5;
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
