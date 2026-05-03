<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加公告')
const isEdit = ref(false)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const noticeList = ref([])

const form = reactive({
  id: null,
  title: '',
  content: '',
  status: 1,
  sort: 0
})

const rules = {
  title: [
    { required: true, message: '请输入公告标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入公告内容', trigger: 'blur' }
  ]
}

// 获取公告列表
const fetchNoticeList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/admin/notice/list', { params })
    if (res.code === 200 || res.code === 0) {
      noticeList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取公告列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchNoticeList()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加公告'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑公告'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该公告吗？', '提示', {
      type: 'warning'
    })
    
    const res = await request.delete(`/admin/notice/delete/${row.id}`)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('删除成功')
      fetchNoticeList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!form.title || !form.content) {
    ElMessage.warning('请填写必填项')
    return
  }
  
  // 转换参数名称以匹配后端
  const submitData = {
    title: form.title,
    content: form.content,
    sort: form.sort,
    status: form.status,
    author: form.author
  }
  
  try {
    let res
    if (isEdit.value) {
      res = await request.put(`/admin/notice/update/${form.id}`, submitData)
    } else {
      res = await request.post('/admin/notice/create', submitData)
    }
    
    if (res.code === 200 || res.code === 0) {
      ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
      dialogVisible.value = false
      fetchNoticeList()
    }
  } catch (error) {
    console.error('提交失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    title: '',
    content: '',
    status: 1,
    sort: 0
  })
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchNoticeList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchNoticeList()
}

// 状态标签
const getStatusTag = (status) => {
  return status === 1 
    ? { type: 'success', text: '发布' }
    : { type: 'danger', text: '下架' }
}

onMounted(() => {
  fetchNoticeList()
})
</script>

<template>
  <div class="notice-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="公告标题" 
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
            <el-option label="发布" :value="1" />
            <el-option label="下架" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作按钮 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-button type="primary" icon="Plus" @click="handleAdd">添加公告</el-button>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table 
        :data="noticeList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column prop="sort" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTag(scope.row.status).type">
              {{ getStatusTag(scope.row.status).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleEdit(scope.row)">
              编辑
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

    <!-- 对话框 -->
    <el-dialog 
      v-model="dialogVisible" 
      :title="dialogTitle" 
      width="700px"
      @close="resetForm"
    >
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入公告标题" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="排序">
              <el-input-number v-model="form.sort" :min="0" style="width: 100%" />
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                提示：数字越小越靠前
              </div>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-radio-group v-model="form.status">
                <el-radio :value="1">发布</el-radio>
                <el-radio :value="0">下架</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="内容" prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="8"
            placeholder="请输入公告内容"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.notice-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
