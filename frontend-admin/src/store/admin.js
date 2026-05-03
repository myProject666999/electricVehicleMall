import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAdminStore = defineStore('admin', () => {
  // 状态
  const token = ref(localStorage.getItem('admin_token') || '')
  const adminInfo = ref(JSON.parse(localStorage.getItem('admin_info') || 'null'))

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)

  // 登录
  const login = (newToken, info) => {
    token.value = newToken
    adminInfo.value = info
    localStorage.setItem('admin_token', newToken)
    if (info) {
      localStorage.setItem('admin_info', JSON.stringify(info))
    }
  }

  // 退出登录
  const logout = () => {
    token.value = ''
    adminInfo.value = null
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_info')
  }

  // 更新管理员信息
  const updateAdminInfo = (info) => {
    adminInfo.value = info
    if (info) {
      localStorage.setItem('admin_info', JSON.stringify(info))
    }
  }

  return {
    token,
    adminInfo,
    isLoggedIn,
    login,
    logout,
    updateAdminInfo
  }
})
