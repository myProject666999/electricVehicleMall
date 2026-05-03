import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../store/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/vehicle',
    name: 'VehicleList',
    component: () => import('../views/VehicleList.vue'),
    meta: { title: '电动车列表' }
  },
  {
    path: '/vehicle/:id',
    name: 'VehicleDetail',
    component: () => import('../views/VehicleDetail.vue'),
    meta: { title: '电动车详情' }
  },
  {
    path: '/notice',
    name: 'NoticeList',
    component: () => import('../views/NoticeList.vue'),
    meta: { title: '公告列表' }
  },
  {
    path: '/notice/:id',
    name: 'NoticeDetail',
    component: () => import('../views/NoticeDetail.vue'),
    meta: { title: '公告详情' }
  },
  {
    path: '/cart',
    name: 'Cart',
    component: () => import('../views/Cart.vue'),
    meta: { title: '购物车', requiresAuth: true }
  },
  {
    path: '/order',
    name: 'OrderList',
    component: () => import('../views/OrderList.vue'),
    meta: { title: '订单列表', requiresAuth: true }
  },
  {
    path: '/order/:id',
    name: 'OrderDetail',
    component: () => import('../views/OrderDetail.vue'),
    meta: { title: '订单详情', requiresAuth: true }
  },
  {
    path: '/order/create',
    name: 'CreateOrder',
    component: () => import('../views/CreateOrder.vue'),
    meta: { title: '创建订单', requiresAuth: true }
  },
  {
    path: '/address',
    name: 'AddressList',
    component: () => import('../views/AddressList.vue'),
    meta: { title: '地址管理', requiresAuth: true }
  },
  {
    path: '/favorite',
    name: 'FavoriteList',
    component: () => import('../views/FavoriteList.vue'),
    meta: { title: '我的收藏', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  },
  {
    path: '/chat',
    name: 'Chat',
    component: () => import('../views/Chat.vue'),
    meta: { title: '在线客服', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || '电动车销售系统'

  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    const userStore = useUserStore()
    if (!userStore.isLoggedIn) {
      // 未登录，跳转到登录页
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }
  }

  next()
})

export default router
