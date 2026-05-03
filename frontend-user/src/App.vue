<script setup>
import { useRoute } from 'vue-router'
import { computed } from 'vue'

const route = useRoute()

// 不需要布局的页面（登录、注册等）
const noLayoutPages = ['/login', '/register']

const showLayout = computed(() => {
  return !noLayoutPages.includes(route.path)
})
</script>

<template>
  <!-- 有布局的页面 -->
  <template v-if="showLayout">
    <div class="app-container">
      <header class="app-header">
        <el-container>
          <el-header class="header-content">
            <div class="logo">
              <router-link to="/">
                <el-icon size="32" color="#409EFF"><ElectricVehicle /></el-icon>
                <span class="logo-text">电动车商城</span>
              </router-link>
            </div>
            
            <el-menu 
              :default-active="route.path" 
              mode="horizontal" 
              background-color="transparent"
              text-color="#303133"
              active-text-color="#409EFF"
              class="nav-menu"
            >
              <el-menu-item index="/">
                <router-link to="/">首页</router-link>
              </el-menu-item>
              <el-menu-item index="/vehicle">
                <router-link to="/vehicle">电动车</router-link>
              </el-menu-item>
              <el-menu-item index="/notice">
                <router-link to="/notice">公告</router-link>
              </el-menu-item>
            </el-menu>

            <div class="header-right">
              <el-input 
                placeholder="搜索电动车..." 
                style="width: 200px; margin-right: 20px;"
                clearable
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>

              <router-link to="/cart" class="cart-icon">
                <el-icon size="24"><ShoppingCart /></el-icon>
                <el-badge :value="0" class="item" style="margin-right: 20px;">
                </el-badge>
              </router-link>

              <template v-if="true">
                <router-link to="/chat" class="chat-icon">
                  <el-icon size="24"><ChatDotRound /></el-icon>
                </router-link>
                <el-dropdown style="margin-left: 20px;">
                  <span class="user-info">
                    <el-avatar :size="32" icon="User" />
                    <span class="username">用户</span>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item>
                        <router-link to="/profile">个人中心</router-link>
                      </el-dropdown-item>
                      <el-dropdown-item>
                        <router-link to="/order">我的订单</router-link>
                      </el-dropdown-item>
                      <el-dropdown-item>
                        <router-link to="/favorite">我的收藏</router-link>
                      </el-dropdown-item>
                      <el-dropdown-item divided>
                        <span style="color: #f56c6c;">退出登录</span>
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
              <template v-else>
                <router-link to="/login" class="login-btn">
                  <el-button type="primary">登录</el-button>
                </router-link>
                <router-link to="/register" class="register-btn">
                  <el-button>注册</el-button>
                </router-link>
              </template>
            </div>
          </el-header>
        </el-container>
      </header>

      <main class="app-main">
        <router-view />
      </main>

      <footer class="app-footer">
        <div class="footer-content">
          <div class="footer-section">
            <h3>关于我们</h3>
            <p>电动车商城是专业的电动车销售平台</p>
            <p>提供优质的电动车产品和服务</p>
          </div>
          <div class="footer-section">
            <h3>客户服务</h3>
            <p>客服热线：400-123-4567</p>
            <p>工作时间：9:00-18:00</p>
          </div>
          <div class="footer-section">
            <h3>关注我们</h3>
            <p>微信公众号：电动车商城</p>
            <p>官方微博：@电动车商城</p>
          </div>
        </div>
        <div class="footer-bottom">
          <p>© 2024 电动车商城 版权所有</p>
        </div>
      </footer>
    </div>
  </template>

  <!-- 无布局的页面（登录、注册等） -->
  <template v-else>
    <router-view />
  </template>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif;
  background-color: #f5f7fa;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.app-header {
  background-color: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
}

.logo {
  display: flex;
  align-items: center;
}

.logo a {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: #303133;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  margin-left: 10px;
  color: #409EFF;
}

.nav-menu {
  flex: 1;
  margin-left: 40px;
  border-bottom: none;
}

.nav-menu .el-menu-item {
  border-bottom: none;
}

.nav-menu a {
  text-decoration: none;
  color: inherit;
}

.header-right {
  display: flex;
  align-items: center;
}

.cart-icon, .chat-icon {
  color: #606266;
  margin-left: 20px;
  position: relative;
}

.cart-icon:hover, .chat-icon:hover {
  color: #409EFF;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.username {
  margin-left: 8px;
  color: #606266;
}

.login-btn, .register-btn {
  margin-left: 10px;
}

.app-main {
  flex: 1;
  max-width: 1200px;
  margin: 20px auto;
  width: 100%;
  padding: 0 20px;
}

.app-footer {
  background-color: #303133;
  color: #fff;
  margin-top: auto;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 40px 20px;
  display: flex;
  justify-content: space-between;
}

.footer-section {
  flex: 1;
}

.footer-section h3 {
  margin-bottom: 15px;
  font-size: 16px;
}

.footer-section p {
  margin-bottom: 8px;
  color: #909399;
  font-size: 14px;
}

.footer-bottom {
  border-top: 1px solid #4a4a4a;
  padding: 20px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}
</style>
