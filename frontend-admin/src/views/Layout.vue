<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '../store/admin'
import { ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()

const isCollapse = ref(false)
const currentPath = ref(route.path)

// 菜单配置
const menuItems = [
  {
    title: '仪表盘',
    icon: 'DataBoard',
    path: '/dashboard'
  },
  {
    title: '管理员管理',
    icon: 'UserFilled',
    path: '/admin'
  },
  {
    title: '用户管理',
    icon: 'User',
    path: '/user'
  },
  {
    title: '分类管理',
    icon: 'Menu',
    path: '/category'
  },
  {
    title: '电动车管理',
    icon: 'Truck',
    path: '/vehicle'
  },
  {
    title: '订单管理',
    icon: 'Document',
    path: '/order'
  },
  {
    title: '公告管理',
    icon: 'Bell',
    path: '/notice'
  },
  {
    title: '轮播图管理',
    icon: 'Picture',
    path: '/banner'
  },
  {
    title: '评论管理',
    icon: 'ChatDotSquare',
    path: '/comment'
  },
  {
    title: '客服回复',
    icon: 'Service',
    path: '/chat'
  }
]

// 切换菜单
const handleMenuClick = (path) => {
  currentPath.value = path
  router.push(path)
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    adminStore.logout()
    router.push('/login')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('退出失败:', error)
    }
  }
}

onMounted(() => {
  currentPath.value = route.path
})
</script>

<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside 
      :width="isCollapse ? '64px' : '220px'" 
      class="layout-aside"
    >
      <div class="logo-section">
        <el-icon size="32" color="#409EFF"><ElectricVehicle /></el-icon>
        <span v-show="!isCollapse" class="logo-text">管理后台</span>
      </div>
      
      <el-menu
        :default-active="currentPath"
        :collapse="isCollapse"
        :collapse-transition="false"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item
          v-for="item in menuItems"
          :key="item.path"
          :index="item.path"
          @click="handleMenuClick(item.path)"
        >
          <el-icon><component :is="item.icon" /></el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <!-- 主内容区 -->
    <el-container class="layout-main">
      <!-- 顶部导航栏 -->
      <el-header class="layout-header">
        <div class="header-left">
          <el-icon 
            class="collapse-icon" 
            @click="isCollapse = !isCollapse"
          >
            <component :is="isCollapse ? 'Expand' : 'Fold'" />
          </el-icon>
          <span class="page-title">{{ route.meta.title || '仪表盘' }}</span>
        </div>
        
        <div class="header-right">
          <el-dropdown>
            <span class="user-info">
              <el-avatar :size="32" icon="UserFilled" />
              <span class="username">{{ adminStore.adminInfo?.username || '管理员' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <!-- 内容区 -->
      <el-main class="layout-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-container {
  height: 100vh;
}

.layout-aside {
  background-color: #304156;
  transition: width 0.3s;
}

.logo-section {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 60px;
  padding: 10px;
  border-bottom: 1px solid #1f2d3d;
}

.logo-text {
  color: #fff;
  font-size: 18px;
  font-weight: bold;
  margin-left: 10px;
}

.layout-main {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 20px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-icon {
  font-size: 20px;
  cursor: pointer;
  color: #606266;
  margin-right: 20px;
}

.collapse-icon:hover {
  color: #409EFF;
}

.page-title {
  font-size: 16px;
  color: #303133;
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 5px 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  margin: 0 8px;
  color: #606266;
  font-size: 14px;
}

.layout-content {
  background: #f5f7fa;
  overflow-y: auto;
  padding: 20px;
}
</style>
