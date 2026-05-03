<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加管理员')
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

const adminList = ref([])

const form = reactive({
  id: null,
  username: '',
  password: '',
  realName: '',
  phone: '',
  status: 1
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ],
  realName: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ]
}

// 获取管理员列表
const fetchAdminList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/admin/list', { params })
    if (res.code === 200 || res.code === 0) {
      adminList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取管理员列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchAdminList()
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
  dialogTitle.value = '添加管理员'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑管理员'
  Object.assign(form, { ...row, password: '' })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该管理员吗？', '提示', {
      type: 'warning'
    })
    
    const res = await request.delete(`/admin/delete/${row.id}`)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('删除成功')
      fetchAdminList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!form.username || !form.realName) {
    ElMessage.warning('请填写必填项')
    return
  }
  
  if (!isEdit.value && !form.password) {
    ElMessage.warning('请输入密码')
    return
  }
  
  try {
    let res
    if (isEdit.value) {
      res = await request.put(`/admin/update/${form.id}`, form)
    } else {
      res = await request.post('/admin/create', form)
    }
    
    if (res.code === 200 || res.code === 0) {
      ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
      dialogVisible.value = false
      fetchAdminList()
    }
  } catch (error) {
    console.error('提交失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    username: '',
    password: '',
    realName: '',
    phone: '',
    status: 1
  })
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchAdminList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchAdminList()
}

// 状态标签
const getStatusTag = (status) => {
  return status === 1 
    ? { type: 'success', text: '正常' }
    : { type: 'danger', text: '禁用' }
}

onMounted(() => {
  fetchAdminList()
})
</script>

<template>
  <div class="admin-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="用户名/姓名" 
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
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
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
      <el-button type="primary" icon="Plus" @click="handleAdd">添加管理员</el-button>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table 
        :data="adminList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="realName" label="姓名" width="100" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="status" label="状态" width="80">
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
      width="500px"
      @close="resetForm"
    >
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef"
        label-width="100px"
      >
        <el-form-item label="用户名" prop="username">
          <el-input 
            v-model="form.username" 
            placeholder="请输入用户名"
            :disabled="isEdit"
          />
        </el-form-item>
        <el-form-item label="密码" :prop="isEdit ? undefined : 'password'">
          <el-input 
            v-model="form.password" 
            type="password"
            :placeholder="isEdit ? '不修改请留空' : '请输入密码'"
            show-password
          />
        </el-form-item>
        <el-form-item label="姓名" prop="realName">
          <el-input v-model="form.realName" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="状态" v-if="isEdit">
          <el-radio-group v-model="form.status">
            <el-radio :value="1">正常</el-radio>
            <el-radio :value="0">禁用</el-radio>
          </el-radio-group>
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
.admin-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
