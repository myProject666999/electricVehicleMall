<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import request from '../utils/request'

const route = useRoute()
const router = useRouter()

const notice = ref(null)
const loading = ref(false)

// 获取公告详情
const fetchNoticeDetail = async () => {
  loading.value = true
  try {
    const res = await request.get(`/notice/detail/${route.params.id}`)
    notice.value = res.data
  } catch (error) {
    console.error('获取公告详情失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchNoticeDetail()
})
</script>

<template>
  <div class="notice-detail-container">
    <!-- 面包屑导航 -->
    <el-breadcrumb class="breadcrumb" separator="/">
      <el-breadcrumb-item>
        <router-link to="/">首页</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>
        <router-link to="/notice">公告列表</router-link>
      </el-breadcrumb-item>
      <el-breadcrumb-item>{{ notice?.title || '公告详情' }}</el-breadcrumb-item>
    </el-breadcrumb>
    
    <div class="notice-content" v-loading="loading">
      <div v-if="notice">
        <h1 class="notice-title">{{ notice.title }}</h1>
        <div class="notice-meta">
          <span>发布时间：{{ new Date(notice.created_at).toLocaleString() }}</span>
          <span v-if="notice.author">作者：{{ notice.author }}</span>
        </div>
        <div class="notice-body" v-html="notice.content">
        </div>
      </div>
      
      <el-empty v-if="!notice && !loading" description="公告不存在" />
    </div>
    
    <div class="action-bar">
      <el-button @click="router.back()">
        <el-icon><ArrowLeft /></el-icon>
        返回列表
      </el-button>
    </div>
  </div>
</template>

<style scoped>
.notice-detail-container {
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

.notice-content {
  background: #fff;
  border-radius: 8px;
  padding: 30px;
  min-height: 400px;
}

.notice-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 15px;
  text-align: center;
  line-height: 1.5;
}

.notice-meta {
  text-align: center;
  color: #909399;
  font-size: 14px;
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid #ebeef5;
}

.notice-meta span {
  margin: 0 15px;
}

.notice-body {
  font-size: 16px;
  color: #303133;
  line-height: 1.8;
}

.notice-body :deep(p) {
  margin-bottom: 15px;
  text-indent: 2em;
}

.notice-body :deep(img) {
  max-width: 100%;
  height: auto;
  display: block;
  margin: 20px auto;
}

.action-bar {
  margin-top: 20px;
  text-align: center;
}
</style>
