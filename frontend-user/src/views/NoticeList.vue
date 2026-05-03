<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '../utils/request'

const router = useRouter()

const notices = ref([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

// 获取公告列表
const fetchNotices = async () => {
  loading.value = true
  try {
    const res = await request.get('/notice/list', {
      params: {
        status: 1,
        page: page.value,
        page_size: pageSize.value
      }
    })
    notices.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('获取公告列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 查看公告详情
const handleViewDetail = (notice) => {
  router.push(`/notice/${notice.id}`)
}

// 分页
const handlePageChange = (val) => {
  page.value = val
  fetchNotices()
}

onMounted(() => {
  fetchNotices()
})
</script>

<template>
  <div class="notice-list-container">
    <h2 class="page-title">公告列表</h2>
    
    <div class="notice-list" v-loading="loading">
      <div 
        v-for="notice in notices" 
        :key="notice.id" 
        class="notice-card"
        @click="handleViewDetail(notice)"
      >
        <div class="notice-header">
          <h3 class="notice-title">{{ notice.title }}</h3>
          <span class="notice-time">{{ new Date(notice.created_at).toLocaleDateString() }}</span>
        </div>
        <div class="notice-summary">
          {{ notice.summary || notice.content?.substring(0, 200) + '...' }}
        </div>
        <div class="notice-footer">
          <span class="read-more">查看详情 →</span>
        </div>
      </div>
      
      <!-- 空状态 -->
      <el-empty 
        v-if="notices.length === 0 && !loading" 
        description="暂无公告"
      />
      
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
.notice-list-container {
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

.notice-list {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.notice-card {
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: all 0.3s;
}

.notice-card:last-child {
  border-bottom: none;
}

.notice-card:hover {
  background: #f5f7fa;
}

.notice-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.notice-title {
  font-size: 18px;
  color: #303133;
  margin: 0;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notice-time {
  font-size: 14px;
  color: #909399;
  margin-left: 20px;
  flex-shrink: 0;
}

.notice-summary {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  margin-bottom: 15px;
}

.notice-footer {
  text-align: right;
}

.read-more {
  font-size: 14px;
  color: #409EFF;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
