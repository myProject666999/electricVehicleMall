<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加电动车')
const isEdit = ref(false)

const searchForm = reactive({
  keyword: '',
  categoryId: '',
  status: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const vehicleList = ref([])
const categoryList = ref([])

const form = reactive({
  id: null,
  name: '',
  categoryId: '',
  price: '',
  stock: '',
  description: '',
  image: '',
  status: 1,
  specification: ''
})

const rules = {
  name: [
    { required: true, message: '请输入电动车名称', trigger: 'blur' }
  ],
  categoryId: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存', trigger: 'blur' }
  ]
}

// 获取电动车列表
const fetchVehicleList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/vehicle/list', { params })
    if (res.code === 200 || res.code === 0) {
      vehicleList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取电动车列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategoryList = async () => {
  try {
    const res = await request.get('/category/list')
    if (res.code === 200 || res.code === 0) {
      categoryList.value = res.data || []
    }
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchVehicleList()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.categoryId = ''
  searchForm.status = ''
  handleSearch()
}

// 添加
const handleAdd = () => {
  isEdit.value = false
  dialogTitle.value = '添加电动车'
  resetForm()
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  dialogTitle.value = '编辑电动车'
  Object.assign(form, { ...row })
  dialogVisible.value = true
}

// 删除
const handleDelete = async (row) => {
  try {
    await ElMessageBox.confirm('确定要删除该电动车吗？', '提示', {
      type: 'warning'
    })
    
    const res = await request.delete(`/admin/vehicle/delete/${row.id}`)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('删除成功')
      fetchVehicleList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 批量删除
const handleBatchDelete = () => {
  ElMessage.warning('批量删除功能开发中')
}

// 提交表单
const handleSubmit = async () => {
  // 转换参数名称以匹配后端
  const submitData = {
    name: form.name,
    category_id: form.categoryId,
    price: form.price,
    stock: form.stock,
    description: form.description,
    image: form.image,
    status: form.status,
    specifications: form.specification,
    sort: 0,
    is_recommend: 0
  }
  
  try {
    let res
    if (isEdit.value) {
      res = await request.put(`/admin/vehicle/update/${form.id}`, submitData)
    } else {
      res = await request.post('/admin/vehicle/create', submitData)
    }
    
    if (res.code === 200 || res.code === 0) {
      ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
      dialogVisible.value = false
      fetchVehicleList()
    }
  } catch (error) {
    console.error('提交失败:', error)
  }
}

// 重置表单
const resetForm = () => {
  Object.assign(form, {
    id: null,
    name: '',
    categoryId: '',
    price: '',
    stock: '',
    description: '',
    image: '',
    status: 1,
    specification: ''
  })
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchVehicleList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchVehicleList()
}

// 状态标签
const getStatusTag = (status) => {
  return status === 1 
    ? { type: 'success', text: '上架' }
    : { type: 'danger', text: '下架' }
}

onMounted(() => {
  fetchCategoryList()
  fetchVehicleList()
})
</script>

<template>
  <div class="vehicle-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="电动车名称" 
            clearable
            style="width: 200px"
          />
        </el-form-item>
        <el-form-item label="分类">
          <el-select 
            v-model="searchForm.categoryId" 
            placeholder="全部分类" 
            clearable
            style="width: 150px"
          >
            <el-option 
              v-for="cat in categoryList" 
              :key="cat.id" 
              :label="cat.name" 
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select 
            v-model="searchForm.status" 
            placeholder="全部状态" 
            clearable
            style="width: 120px"
          >
            <el-option label="上架" :value="1" />
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
      <el-button type="primary" icon="Plus" @click="handleAdd">添加电动车</el-button>
      <el-button type="danger" icon="Delete" @click="handleBatchDelete">批量删除</el-button>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never">
      <el-table 
        :data="vehicleList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="image" label="图片" width="80">
          <template #default="scope">
            <el-image 
              :src="scope.row.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle&image_size=square'" 
              :preview-src-list="[scope.row.image || 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=electric%20vehicle&image_size=square']"
              fit="cover"
              style="width: 50px; height: 50px"
            />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="150" />
        <el-table-column prop="price" label="价格" width="120">
          <template #default="scope">
            ¥{{ scope.row.price?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="stock" label="库存" width="80" />
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
      width="700px"
      @close="resetForm"
    >
      <el-form 
        :model="form" 
        :rules="rules" 
        ref="formRef"
        label-width="100px"
      >
        <el-form-item label="名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入电动车名称" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="分类" prop="categoryId">
              <el-select 
                v-model="form.categoryId" 
                placeholder="请选择分类"
                style="width: 100%"
              >
                <el-option 
                  v-for="cat in categoryList" 
                  :key="cat.id" 
                  :label="cat.name" 
                  :value="cat.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="form.status">
                <el-radio :value="1">上架</el-radio>
                <el-radio :value="0">下架</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="价格" prop="price">
              <el-input-number 
                v-model="form.price" 
                :min="0" 
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="库存" prop="stock">
              <el-input-number 
                v-model="form.stock" 
                :min="0"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="图片" prop="image">
          <el-input v-model="form.image" placeholder="请输入图片URL" />
        </el-form-item>
        <el-form-item label="规格" prop="specification">
          <el-input 
            v-model="form.specification" 
            type="textarea" 
            :rows="2"
            placeholder="请输入规格信息"
          />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input 
            v-model="form.description" 
            type="textarea" 
            :rows="4"
            placeholder="请输入电动车描述"
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
.vehicle-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
