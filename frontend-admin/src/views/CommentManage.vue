<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const commentList = ref([])

// 获取评论列表
const fetchCommentList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/admin/comment/list', { params })
    if (res.code === 200 || res.code === 0) {
      commentList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取评论列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchCommentList()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  handleSearch()
}

// 审核通过
const handleApprove = async (row) => {
  try {
    const res = await request.put(`/admin/comment/status/${row.id}`, {
      status: 1
    })
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('审核通过')
      fetchCommentList()
    }
  } catch (error) {
    console.error('审核失败:', error)
  }
}

// 审核拒绝
const handleReject = async (row) => {
  try {
    const res = await request.put(`/admin/comment/status/${row.id}`, {
      status: 2
    })
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('已拒绝')
      fetchCommentList()
    }
  } catch (error) {
    console.error('操作失败:', error)
  }
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该评论吗？', '提示', {
      type: 'warning'
    })
    
    const res = await request.delete(`/admin/comment/delete/${row.id}`)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('删除成功')
      fetchCommentList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchCommentList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchCommentList()
}

// 状态标签
const getStatusTag = (status) => {
  const statusMap = {
    0: { type: 'info', text: '待审核' },
    1: { type: 'success', text: '已通过' },
    2: { type: 'danger', text: '已拒绝' }
  }
  return statusMap[status] || { type: 'info', text: '未知' }
}

// 星星显示
const renderStars = (rating) => {
  const fullStars = Math.floor(rating)
  const hasHalf = rating % 1 !== 0
  let stars = ''
  for (let i = 0; i < fullStars; i++) {
    stars += '★'
  }
  if (hasHalf) {
    stars += '☆'
  }
  return stars
}

onMounted(() => {
  fetchCommentList()
})
</script>

<template>
  <div class="comment-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="评论内容/电动车名称" 
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="searchForm.status" 
            placeholder="全部状态" 
            clearable
            style="width: 120px"
          >
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table 
        :data="commentList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column prop="vehicleName" label="电动车" min-width="150" />
        <el-table-column prop="rating" label="评分" width="120">
          <template #default="scope">
            <span style="color: #f59e0b;">{{ renderStars(scope.row.rating) }}</span>
            <span style="margin-left: 5px;">{{ scope.row.rating }}分</span>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTag(scope.row.status).type">
              {{ getStatusTag(scope.row.status).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button 
              v-if="scope.row.status === 0" 
              type="success" 
              size="small" 
              @click="handleApprove(scope.row)"
            >
              通过
            </el-button>
            <el-button 
              v-if="scope.row.status === 0" 
              type="warning" 
              size="small" 
              @click="handleReject(scope.row)"
            >
              拒绝
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<style scoped>
.comment-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
