<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const addresses = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editingAddress = ref(null)

const form = ref({
  id: null,
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  is_default: false
})

const rules = {
  name: [{ required: true, message: '请输入收货人姓名', trigger: 'blur' }],
  phone: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  province: [{ required: true, message: '请输入省份', trigger: 'blur' }],
  city: [{ required: true, message: '请输入城市', trigger: 'blur' }],
  district: [{ required: true, message: '请输入区县', trigger: 'blur' }],
  detail: [{ required: true, message: '请输入详细地址', trigger: 'blur' }]
}

// 获取地址列表
const fetchAddresses = async () => {
  loading.value = true
  try {
    const res = await request.get('/address/list')
    addresses.value = res.data || []
  } catch (error) {
    console.error('获取地址列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 打开对话框
const openDialog = (address = null) => {
  editingAddress.value = address
  if (address) {
    form.value = { ...address }
  } else {
    form.value = {
      id: null,
      name: '',
      phone: '',
      province: '',
      city: '',
      district: '',
      detail: '',
      is_default: false
    }
  }
  dialogVisible.value = true
}

// 保存地址
const handleSave = async () => {
  try {
    if (editingAddress.value) {
      await request.post('/address/update', form.value)
      ElMessage.success('修改成功')
    } else {
      await request.post('/address/add', form.value)
      ElMessage.success('添加成功')
    }
    dialogVisible.value = false
    fetchAddresses()
  } catch (error) {
    console.error('保存地址失败:', error)
    ElMessage.error('保存失败')
  }
}

// 删除地址
const handleDelete = async (address) => {
  try {
    await ElMessageBox.confirm('确定要删除该地址吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await request.post('/address/delete', { id: address.id })
    ElMessage.success('删除成功')
    fetchAddresses()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 设置默认地址
const handleSetDefault = async (address) => {
  try {
    await request.post('/address/set-default', { id: address.id })
    ElMessage.success('设置成功')
    fetchAddresses()
  } catch (error) {
    console.error('设置默认地址失败:', error)
    ElMessage.error('设置失败')
  }
}

onMounted(() => {
  fetchAddresses()
})
</script>

<template>
  <div class="address-list-container">
    <div class="header-section">
      <h2 class="page-title">地址管理</h2>
      <el-button type="primary" @click="openDialog()">
        <el-icon><Plus /></el-icon>
        添加地址
      </el-button>
    </div>
    
    <!-- 地址列表 -->
    <div class="address-list" v-loading="loading">
      <div v-for="address in addresses" :key="address.id" class="address-card">
        <div class="address-info">
          <div class="address-header">
            <span class="name">{{ address.name }}</span>
            <span class="phone">{{ address.phone }}</span>
            <el-tag v-if="address.is_default" type="success" size="small">默认</el-tag>
          </div>
          <div class="address-detail">
            {{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}
          </div>
        </div>
        <div class="address-actions">
          <el-button 
            v-if="!address.is_default" 
            type="text" 
            size="small"
            @click="handleSetDefault(address)"
          >
            设为默认
          </el-button>
          <el-button type="primary" link size="small" @click="openDialog(address)">
            编辑
          </el-button>
          <el-button type="danger" link size="small" @click="handleDelete(address)">
            删除
          </el-button>
        </div>
      </div>
      
      <!-- 空状态 -->
      <el-empty 
        v-if="addresses.length === 0 && !loading" 
        description="暂无收货地址"
      >
        <el-button type="primary" @click="openDialog()">
          添加地址
        </el-button>
      </el-empty>
    </div>
    
    <!-- 添加/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="editingAddress ? '编辑地址' : '添加地址'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="收货人" prop="name">
          <el-input v-model="form.name" placeholder="请输入收货人姓名" />
        </el-form-item>
        
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入手机号" />
        </el-form-item>
        
        <el-form-item label="省份" prop="province">
          <el-input v-model="form.province" placeholder="请输入省份" />
        </el-form-item>
        
        <el-form-item label="城市" prop="city">
          <el-input v-model="form.city" placeholder="请输入城市" />
        </el-form-item>
        
        <el-form-item label="区县" prop="district">
          <el-input v-model="form.district" placeholder="请输入区县" />
        </el-form-item>
        
        <el-form-item label="详细地址" prop="detail">
          <el-input
            v-model="form.detail"
            type="textarea"
            :rows="2"
            placeholder="请输入详细地址"
          />
        </el-form-item>
        
        <el-form-item label="默认地址">
          <el-switch v-model="form.is_default" active-text="是" inactive-text="否" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSave">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.address-list-container {
  padding: 20px 0;
}

.header-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 24px;
  color: #303133;
  border-bottom: 2px solid #409EFF;
  padding-bottom: 10px;
  display: inline-block;
}

.address-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 15px;
  border: 1px solid #ebeef5;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.address-info {
  flex: 1;
}

.address-header {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}

.name {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  margin-right: 20px;
}

.phone {
  color: #606266;
  margin-right: 15px;
}

.address-detail {
  color: #606266;
  line-height: 1.6;
}

.address-actions {
  display: flex;
  gap: 10px;
}
</style>
