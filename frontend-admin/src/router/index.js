import { createRouter, createWebHistory } from 'vue-router'
import { useAdminStore } from '../store/admin'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '管理员登录', requiresAuth: false }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('../views/Layout.vue'),
    meta: { title: '管理后台', requiresAuth: true },
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '仪表盘', icon: 'DataBoard' }
      },
      {
        path: '/admin',
        name: 'AdminManage',
        component: () => import('../views/AdminManage.vue'),
        meta: { title: '管理员管理', icon: 'UserFilled' }
      },
      {
        path: '/user',
        name: 'UserManage',
        component: () => import('../views/UserManage.vue'),
        meta: { title: '用户管理', icon: 'User' }
      },
      {
        path: '/category',
        name: 'CategoryManage',
        component: () => import('../views/CategoryManage.vue'),
        meta: { title: '分类管理', icon: 'Menu' }
      },
      {
        path: '/vehicle',
        name: 'VehicleManage',
        component: () => import('../views/VehicleManage.vue'),
        meta: { title: '电动车管理', icon: 'Truck' }
      },
      {
        path: '/order',
        name: 'OrderManage',
        component: () => import('../views/OrderManage.vue'),
        meta: { title: '订单管理', icon: 'Document' }
      },
      {
        path: '/notice',
        name: 'NoticeManage',
        component: () => import('../views/NoticeManage.vue'),
        meta: { title: '公告管理', icon: 'Bell' }
      },
      {
        path: '/banner',
        name: 'BannerManage',
        component: () => import('../views/BannerManage.vue'),
        meta: { title: '轮播图管理', icon: 'Picture' }
      },
      {
        path: '/comment',
        name: 'CommentManage',
        component: () => import('../views/CommentManage.vue'),
        meta: { title: '评论管理', icon: 'ChatDotSquare' }
      },
      {
        path: '/chat',
        name: 'ChatManage',
        component: () => import('../views/ChatManage.vue'),
        meta: { title: '客服回复', icon: 'Service' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = to.meta.title || '电动车销售系统后台'

  // 检查是否需要登录
  if (to.meta.requiresAuth !== false) {
    const adminStore = useAdminStore()
    if (!adminStore.isLoggedIn) {
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
