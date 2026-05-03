import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAdminStore } from '../store/admin'

// 创建axios实例
const request = axios.create({
  baseURL: 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  (config) => {
    const adminStore = useAdminStore()
    if (adminStore.token) {
      config.headers.Authorization = `Bearer ${adminStore.token}`
    }
    return config
  },
  (error) => {
    console.error('请求错误:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response) => {
    const res = response.data
    
    // 检查响应状态
    if (res.code === 200 || res.code === 0) {
      return res
    }
    
    // 处理错误
    ElMessage.error(res.message || '请求失败')
    return Promise.reject(new Error(res.message || '请求失败'))
  },
  (error) => {
    console.error('响应错误:', error)
    
    // 处理HTTP错误状态码
    if (error.response) {
      switch (error.response.status) {
        case 401:
          // 未授权，跳转到登录页
          ElMessageBox.confirm('登录状态已过期，请重新登录', '提示', {
            confirmButtonText: '重新登录',
            cancelButtonText: '取消',
            type: 'warning'
          }).then(() => {
            const adminStore = useAdminStore()
            adminStore.logout()
            window.location.href = '/login'
          })
          break
        case 403:
          ElMessage.error('没有权限访问该资源')
          break
        case 404:
          ElMessage.error('请求的资源不存在')
          break
        case 500:
          ElMessage.error('服务器内部错误')
          break
        default:
          ElMessage.error(error.response.data?.message || '请求失败')
      }
    } else if (error.message.includes('timeout')) {
      ElMessage.error('请求超时，请检查网络连接')
    } else {
      ElMessage.error('网络错误，请稍后重试')
    }
    
    return Promise.reject(error)
  }
)

export default request
