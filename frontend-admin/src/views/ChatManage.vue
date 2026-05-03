<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '../utils/request'

const loading = ref(false)
const chatLoading = ref(false)
const chatVisible = ref(false)
const currentChat = ref(null)
const messageInput = ref('')
const messageContainer = ref(null)

const searchForm = reactive({
  keyword: '',
  status: ''
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const chatList = ref([])
const messages = ref([])

// 获取会话列表
const fetchChatList = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.currentPage,
      page_size: pagination.pageSize,
      ...searchForm
    }
    
    const res = await request.get('/admin/chat/users', { params })
    if (res.code === 200 || res.code === 0) {
      chatList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error('获取会话列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.currentPage = 1
  fetchChatList()
}

// 重置搜索
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  handleSearch()
}

// 打开聊天详情
const openChat = async (row) => {
  currentChat.value = row
  chatVisible.value = true
  messages.value = []
  
  try {
    const res = await request.get(`/admin/chat/messages/${row.id}`)
    if (res.code === 200 || res.code === 0) {
      messages.value = res.data || []
    }
  } catch (error) {
    console.error('获取消息失败:', error)
  }
  
  // 滚动到底部
  await nextTick()
  scrollToBottom()
}

// 滚动到底部
const scrollToBottom = () => {
  if (messageContainer.value) {
    messageContainer.value.scrollTop = messageContainer.value.scrollHeight
  }
}

// 发送消息
const sendMessage = async () => {
  if (!messageInput.value.trim()) {
    ElMessage.warning('请输入消息内容')
    return
  }
  
  if (!currentChat.value) {
    return
  }
  
  chatLoading.value = true
  try {
    const res = await request.post('/admin/chat/reply', {
      user_id: currentChat.value.id,
      content: messageInput.value
    })
    
    if (res.code === 200 || res.code === 0) {
      messages.value.push({
        id: Date.now(),
        content: messageInput.value,
        is_admin: true,
        created_at: new Date().toLocaleString()
      })
      messageInput.value = ''
      
      await nextTick()
      scrollToBottom()
    }
  } catch (error) {
    console.error('发送消息失败:', error)
  } finally {
    chatLoading.value = false
  }
}

// 分页变化
const handleSizeChange = (size) => {
  pagination.pageSize = size
  fetchChatList()
}

const handleCurrentChange = (page) => {
  pagination.currentPage = page
  fetchChatList()
}

// 状态标签
const getStatusTag = (status) => {
  return status === 1 
    ? { type: 'success', text: '已回复' }
    : { type: 'warning', text: '待回复' }
}

onMounted(() => {
  fetchChatList()
})
</script>

<template>
  <div class="chat-manage-container">
    <!-- 搜索区域 -->
    <el-card shadow="never" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input 
            v-model="searchForm.keyword" 
            placeholder="用户名/消息内容" 
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
            <el-option label="待回复" :value="0" />
            <el-option label="已回复" :value="1" />
          </el-select>
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
        :data="chatList" 
        v-loading="loading" 
        style="width: 100%"
        stripe
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="100" />
        <el-table-column prop="latestMessage" label="最新消息" min-width="200" show-overflow-tooltip />
        <el-table-column prop="unreadCount" label="未读" width="80">
          <template #default="scope">
            <el-badge :value="scope.row.unreadCount || 0" :hidden="!(scope.row.unreadCount > 0)">
              <span>{{ scope.row.unreadCount || 0 }}</span>
            </el-badge>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTag(scope.row.status).type">
              {{ getStatusTag(scope.row.status).text }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastMessageTime" label="最后消息时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="scope">
            <el-button type="primary" size="small" @click="openChat(scope.row)">
              回复
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

    <!-- 聊天对话框 -->
    <el-dialog 
      v-model="chatVisible" 
      title="客服聊天" 
      width="600px"
      :close-on-click-modal="false"
    >
      <div class="chat-container">
        <!-- 消息列表 -->
        <div ref="messageContainer" class="message-list">
          <el-empty 
            v-if="messages.length === 0 && !chatLoading" 
            description="暂无消息" 
          />
          
          <div 
            v-for="msg in messages" 
            :key="msg.id" 
            class="message-item"
            :class="{ 'admin-message': msg.isAdmin }"
          >
            <div class="message-avatar">
              <el-avatar :size="40">
                <el-icon v-if="msg.isAdmin"><UserFilled /></el-icon>
                <el-icon v-else><User /></el-icon>
              </el-avatar>
            </div>
            <div class="message-content">
              <div class="message-info">
                <span class="message-name">
                  {{ msg.isAdmin ? '管理员' : (msg.username || '用户') }}
                </span>
                <span class="message-time">{{ msg.createdAt }}</span>
              </div>
              <div class="message-bubble">
                {{ msg.content }}
              </div>
            </div>
          </div>
          
          <div v-if="chatLoading" class="loading-message">
            <i class="el-icon-loading"></i> 加载中...
          </div>
        </div>
        
        <!-- 消息输入 -->
        <div class="message-input-area">
          <el-input
            v-model="messageInput"
            type="textarea"
            :rows="2"
            placeholder="输入回复内容..."
            @keyup.enter.ctrl="sendMessage"
          />
          <div class="input-actions">
            <span style="color: #909399; font-size: 12px;">Ctrl+Enter 发送</span>
            <el-button type="primary" @click="sendMessage" :loading="chatLoading">
              发送
            </el-button>
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<style scoped>
.chat-manage-container {
  padding: 0;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.chat-container {
  display: flex;
  flex-direction: column;
  height: 500px;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 10px;
  background: #f5f7fa;
  border-radius: 4px;
  margin-bottom: 15px;
}

.message-item {
  display: flex;
  margin-bottom: 15px;
}

.admin-message {
  flex-direction: row-reverse;
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  margin: 0 10px;
}

.admin-message .message-content {
  text-align: right;
}

.message-info {
  margin-bottom: 5px;
}

.message-name {
  font-size: 12px;
  color: #909399;
  margin-right: 10px;
}

.message-time {
  font-size: 12px;
  color: #c0c4cc;
}

.message-bubble {
  display: inline-block;
  padding: 10px 15px;
  background: #fff;
  border-radius: 8px;
  word-break: break-all;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.admin-message .message-bubble {
  background: #409EFF;
  color: #fff;
}

.loading-message {
  text-align: center;
  color: #909399;
  padding: 20px;
}

.message-input-area {
  flex-shrink: 0;
}

.input-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}
</style>
