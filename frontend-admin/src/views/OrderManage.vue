<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const detailVisible = ref(false)
const shipVisible = ref(false)
const currentOrder = ref({})

const searchForm = reactive({
  orderNo: '',
  status: '',
  startDate: '',
  endDate: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const orderList = ref([])

// 获取订单列表
const fetchOrderList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/admin/order/list', { params })
    if (res.code === 200 || res.code === 0) {
      orderList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取订单列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchOrderList()
}

// 重置搜索
const handleReset = () => {
  searchForm.orderNo = ''
  searchForm.status = ''
  searchForm.startDate = ''
  searchForm.endDate = ''
  handleSearch()
}

// 查看详情
const handleDetail = (row) => {
  currentOrder.value = row
  detailVisible.value = true
}

// 发货
const handleShip = (row) => {
  currentOrder.value = { ...row, logisticsCompany: '', trackingNo: '' }
  shipVisible.value = true
}

// 提交发货
const submitShip = async () => {
  if (!currentOrder.value.logisticsCompany || !currentOrder.value.trackingNo) {
    ElMessage.warning('请填写物流公司和物流单号')
    return
  }
  
  try {
    const res = await request.post(`/admin/order/ship/${currentOrder.value.id}`, {
      logistics_com: currentOrder.value.logisticsCompany,
      logistics_no: currentOrder.value.trackingNo
    })
    
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('发货成功')
      shipVisible.value = false
      fetchOrderList()
    }
  } catch (error) {
    console.error('发货失败:', error)
  }
}

// 取消订单
const handleCancel = async (row) => {
  try {
    await ElMessageBox.confirm('确定要取消该订单吗？', '提示', {
      type: 'warning'
    })
    
    const res = await request.put(`/admin/orders/${row.id}/cancel`)
    if (res.code === 200 || res.code === 0) {
      ElMessage.success('取消成功')
      fetchOrderList()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消失败:', error)
    }
  }
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchOrderList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchOrderList()
}

// 订单状态标签
const getStatusTag = (status) => {
  const statusMap = {
    0: { type: 'info', text: '待支付' },
    1: { type: 'warning', text: '待发货' },
    2: { type: 'primary', text: '已发货' },
    3: { type: 'success', text: '已完成' },
    4: { type: 'danger', text: '已取消' }
  }
  return statusMap[status] || { type: 'info', text: '未知' }
}

onMounted(() => {
  fetchOrderList()
})
</script>

<template>
  <div class="order-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="订单号">
          <el-input 
            v-model="searchForm.orderNo" 
            placeholder="请输入订单号" 
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
            <el-option label="待支付" :value="0" />
            <el-option label="待发货" :value="1" />
            <el-option label="已发货" :value="2" />
            <el-option label="已完成" :value="3" />
            <el-option label="已取消" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="searchForm.startDate"
            type="date"
            placeholder="开始日期"
            style="width: 150px"
          />
          <span style="margin: 0 10px;">至</span>
          <el-date-picker
            v-model="searchForm.endDate"
            type="date"
            placeholder="结束日期"
            style="width: 150px"
          />
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
        :data="orderList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column prop="orderNo" label="订单号" width="200" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column prop="totalAmount" label="订单金额" width="120">
          <template #default="scope">
            ¥{{ scope.row.totalAmount?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTag(scope.row.status).type">
              {{ getStatusTag(scope.row.status).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="logisticsCompany" label="物流公司" width="120">
          <template #default="scope">
            {{ scope.row.logisticsCompany || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="trackingNo" label="物流单号" width="150">
          <template #default="scope">
            {{ scope.row.trackingNo || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="handleDetail(scope.row)">
              详情
            </el-button>
            <el-button 
              v-if="scope.row.status === 1" 
              type="success" 
              size="small" 
              @click="handleShip(scope.row)"
            >
              发货
            </el-button>
            <el-button 
              v-if="scope.row.status === 0 || scope.row.status === 1" 
              type="danger" 
              size="small" 
              @click="handleCancel(scope.row)"
            >
              取消
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

    <!-- 订单详情对话框 -->
    <el-dialog 
      v-model="detailVisible" 
      title="订单详情" 
      width="800px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ currentOrder.orderNo }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusTag(currentOrder.status).type">
            {{ getStatusTag(currentOrder.status).text }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="用户">{{ currentOrder.username }}</el-descriptions-item>
        <el-descriptions-item label="收货人">{{ currentOrder.consigneeName }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ currentOrder.consigneePhone }}</el-descriptions-item>
        <el-descriptions-item label="收货地址">{{ currentOrder.consigneeAddress }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">¥{{ currentOrder.totalAmount?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="订单备注">{{ currentOrder.remark || '无' }}</el-descriptions-item>
        <el-descriptions-item label="物流公司">{{ currentOrder.logisticsCompany || '-' }}</el-descriptions-item>
        <el-descriptions-item label="物流单号">{{ currentOrder.trackingNo || '-' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentOrder.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="支付时间">{{ currentOrder.paidAt || '-' }}</el-descriptions-item>
      </el-descriptions>
      
      <el-divider>订单商品</el-divider>
      
      <el-table :data="currentOrder.items || []" style="width: 100%">
        <el-table-column prop="vehicleName" label="商品名称" />
        <el-table-column prop="price" label="单价" width="120">
          <template #default="scope">
            ¥{{ scope.row.price?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="80" />
        <el-table-column prop="subtotal" label="小计" width="120">
          <template #default="scope">
            ¥{{ scope.row.subtotal?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
      </el-table>
      
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 发货对话框 -->
    <el-dialog 
      v-model="shipVisible" 
      title="订单发货" 
      width="500px"
    >
      <el-form label-width="100px">
        <el-form-item label="订单号">
          <el-input :value="currentOrder.orderNo" disabled />
        </el-form-item>
        <el-form-item label="物流公司" required>
          <el-select 
            v-model="currentOrder.logisticsCompany" 
            placeholder="请选择物流公司"
            style="width: 100%"
          >
            <el-option label="顺丰速运" value="顺丰速运" />
            <el-option label="中通快递" value="中通快递" />
            <el-option label="圆通速递" value="圆通速递" />
            <el-option label="申通快递" value="申通快递" />
            <el-option label="韵达速递" value="韵达速递" />
            <el-option label="京东物流" value="京东物流" />
          </el-select>
        </el-form-item>
        <el-form-item label="物流单号" required>
          <el-input 
            v-model="currentOrder.trackingNo" 
            placeholder="请输入物流单号"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="shipVisible = false">取消</el-button>
        <el-button type="primary" @click="submitShip">确认发货</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.order-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
