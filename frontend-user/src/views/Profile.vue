<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '../store/user'
import request from '../utils/request'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const activeTab = ref('info')

const userForm = reactive({
  username: '',
  phone: '',
  email: '',
  nickname: '',
  avatar: ''
})

const passwordForm = reactive({
  old_password: '',
  new_password: '',
  confirm_password: ''
})

const menuList = [
  {
    title: '我的订单',
    icon: 'Document',
    items: [
      { name: '全部订单', path: '/order' },
      { name: '待付款', path: '/order?status=pending_payment' },
      { name: '待发货', path: '/order?status=pending_shipment' },
      { name: '已发货', path: '/order?status=shipped' },
      { name: '已完成', path: '/order?status=completed' }
    ]
  },
  {
    title: '个人中心',
    icon: 'User',
    items: [
      { name: '个人信息', tab: 'info' },
      { name: '修改密码', tab: 'password' },
      { name: '收货地址', path: '/address' },
      { name: '我的收藏', path: '/favorite' }
    ]
  },
  {
    title: '客户服务',
    icon: 'Service',
    items: [
      { name: '在线客服', path: '/chat' }
    ]
  }
]

// 获取用户信息
const fetchUserInfo = async () => {
  loading.value = true
  try {
    const res = await request.get('/user/info')
    Object.assign(userForm, res.data)
  } catch (error) {
    console.error('获取用户信息失败:', error)
  } finally {
    loading.value = false
  }
}

// 更新用户信息
const handleUpdateInfo = async () => {
  try {
    await request.post('/user/update', userForm)
    ElMessage.success('更新成功')
  } catch (error) {
    console.error('更新信息失败:', error)
    ElMessage.error('更新失败')
  }
}

// 修改密码
const handleUpdatePassword = async () => {
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    ElMessage.warning('两次输入的密码不一致')
    return
  }
  
  try {
    await request.post('/user/change-password', {
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password
    })
    ElMessage.success('密码修改成功')
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error('修改失败')
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    userStore.logout()
    ElMessage.success('退出成功')
    router.push('/login')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('退出失败:', error)
    }
  }
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<template>
  <div class="profile-container">
    <el-row :gutter="20">
      <!-- 左侧菜单 -->
      <el-col :span="6">
        <div class="user-card">
          <el-avatar :size="80" icon="User" />
          <h3 class="username">{{ userForm.nickname || userForm.username }}</h3>
          <p class="user-phone">{{ userForm.phone || '未绑定手机号' }}</p>
        </div>
        
        <div class="menu-card">
          <div v-for="menu in menuList" :key="menu.title" class="menu-section">
            <div class="menu-title">
              <el-icon><component :is="menu.icon" /></el-icon>
              {{ menu.title }}
            </div>
            <div class="menu-items">
              <router-link
                v-for="item in menu.items"
                :key="item.name"
                :to="item.path || ''"
                class="menu-item"
                :class="{ active: item.tab === activeTab }"
                @click="item.tab && (activeTab = item.tab)"
              >
                {{ item.name }}
              </router-link>
            </div>
          </div>
          
          <div class="menu-section">
            <div class="menu-items">
              <div class="menu-item logout" @click="handleLogout">
                <el-icon><SwitchButton /></el-icon>
                退出登录
              </div>
            </div>
          </div>
        </div>
      </el-col>
      
      <!-- 右侧内容 -->
      <el-col :span="18">
        <div class="content-card">
          <el-tabs v-model="activeTab" style="display: none;">
            <el-tab-pane label="个人信息" name="info" />
            <el-tab-pane label="修改密码" name="password" />
          </el-tabs>
          
          <!-- 个人信息 -->
          <div v-if="activeTab === 'info'" v-loading="loading">
            <h3 class="content-title">个人信息</h3>
            <el-form :model="userForm" label-width="100px" style="max-width: 500px;">
              <el-form-item label="用户名">
                <el-input v-model="userForm.username" disabled />
              </el-form-item>
              <el-form-item label="昵称">
                <el-input v-model="userForm.nickname" placeholder="请输入昵称" />
              </el-form-item>
              <el-form-item label="手机号">
                <el-input v-model="userForm.phone" placeholder="请输入手机号" />
              </el-form-item>
              <el-form-item label="邮箱">
                <el-input v-model="userForm.email" placeholder="请输入邮箱" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleUpdateInfo">保存修改</el-button>
              </el-form-item>
            </el-form>
          </div>
          
          <!-- 修改密码 -->
          <div v-if="activeTab === 'password'">
            <h3 class="content-title">修改密码</h3>
            <el-form :model="passwordForm" label-width="120px" style="max-width: 500px;">
              <el-form-item label="原密码">
                <el-input 
                  v-model="passwordForm.old_password" 
                  type="password" 
                  show-password
                  placeholder="请输入原密码"
                />
              </el-form-item>
              <el-form-item label="新密码">
                <el-input 
                  v-model="passwordForm.new_password" 
                  type="password" 
                  show-password
                  placeholder="请输入新密码"
                />
              </el-form-item>
              <el-form-item label="确认新密码">
                <el-input 
                  v-model="passwordForm.confirm_password" 
                  type="password" 
                  show-password
                  placeholder="请再次输入新密码"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleUpdatePassword">确认修改</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.profile-container {
  padding: 20px 0;
}

.user-card {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 8px;
  padding: 30px 20px;
  text-align: center;
  color: #fff;
  margin-bottom: 15px;
}

.username {
  margin: 15px 0 8px;
  font-size: 18px;
}

.user-phone {
  font-size: 14px;
  opacity: 0.8;
}

.menu-card {
  background: #fff;
  border-radius: 8px;
  padding: 10px 0;
}

.menu-section {
  padding: 10px 0;
}

.menu-title {
  padding: 10px 20px;
  font-size: 16px;
  font-weight: 500;
  color: #303133;
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-items {
  padding: 5px 0;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  color: #606266;
  transition: all 0.3s;
  text-decoration: none;
  gap: 8px;
}

.menu-item:hover,
.menu-item.active {
  background: #ecf5ff;
  color: #409EFF;
}

.menu-item.logout {
  color: #f56c6c;
}

.menu-item.logout:hover {
  background: #fef0f0;
  color: #f56c6c;
}

.content-card {
  background: #fff;
  border-radius: 8px;
  padding: 30px;
  min-height: 500px;
}

.content-title {
  font-size: 18px;
  color: #303133;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}
</style>
