<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import request from '../utils/request'

const messageList = ref([])
const newMessage = ref('')
const loading = ref(false)
const messageContainer = ref(null)

// 获取消息列表
const fetchMessages = async () => {
  loading.value = true
  try {
    const res = await request.get('/chat/list')
    messageList.value = res.data || []
  } catch (error) {
    console.error('获取消息列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 滚动到底部
const scrollToBottom = () => {
  if (messageContainer.value) {
    setTimeout(() => {
      messageContainer.value.scrollTop = messageContainer.value.scrollHeight
    }, 100)
  }
}

// 发送消息
const handleSendMessage = async () => {
  if (!newMessage.value.trim()) {
    return
  }
  
  // 添加到本地消息列表
  const tempMessage = {
    id: Date.now(),
    content: newMessage.value,
    is_admin: false,
    created_at: new Date().toISOString()
  }
  
  messageList.value.push(tempMessage)
  const content = newMessage.value
  newMessage.value = ''
  
  scrollToBottom()
  
  try {
    const res = await request.post('/chat/send', { content })
    // 更新消息ID
    const index = messageList.value.findIndex(msg => msg.id === tempMessage.id)
    if (index > -1) {
      messageList.value[index] = res.data
    }
  } catch (error) {
    console.error('发送消息失败:', error)
    ElMessage.error('发送失败')
    // 移除临时消息
    const index = messageList.value.findIndex(msg => msg.id === tempMessage.id)
    if (index > -1) {
      messageList.value.splice(index, 1)
    }
  }
}

// 模拟客服回复（演示用）
const simulateAdminReply = () => {
  setTimeout(() => {
    messageList.value.push({
      id: Date.now(),
      content: '您好，请问有什么可以帮助您的？',
      is_admin: true,
      created_at: new Date().toISOString()
    })
    scrollToBottom()
  }, 1000)
}

onMounted(() => {
  fetchMessages()
  scrollToBottom()
})
</script>

<template>
  <div class="chat-container">
    <h2 class="page-title">在线客服</h2>
    
    <div class="chat-box">
      <!-- 消息列表 -->
      <div class="message-list" ref="messageContainer" v-loading="loading">
        <div v-if="messageList.length === 0 && !loading" class="empty-messages">
          <el-empty description="暂无消息，开始对话吧" />
        </div>
        
        <div 
          v-for="message in messageList" 
          :key="message.id" 
          class="message-item"
          :class="{ 'admin-message': message.is_admin }"
        >
          <div class="message-avatar">
            <el-avatar 
              :size="40" 
              :icon="message.is_admin ? 'Service' : 'User'"
              :style="{ background: message.is_admin ? '#409EFF' : '#67c23a' }"
            />
          </div>
          <div class="message-content">
            <div class="message-header">
              <span class="sender">{{ message.is_admin ? '客服' : '我' }}</span>
              <span class="time">{{ new Date(message.created_at).toLocaleString() }}</span>
            </div>
            <div class="message-bubble">
              {{ message.content }}
            </div>
          </div>
        </div>
      </div>
      
      <!-- 输入区域 -->
      <div class="input-area">
        <el-input
          v-model="newMessage"
          type="textarea"
          :rows="3"
          placeholder="请输入消息..."
          @keyup.enter.ctrl="handleSendMessage"
        />
        <div class="action-bar">
          <span class="tip">按 Ctrl + Enter 发送</span>
          <el-button type="primary" @click="handleSendMessage">
            <el-icon><Promotion /></el-icon>
            发送
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.chat-container {
  padding: 20px 0;
  height: calc(100vh - 140px);
  display: flex;
  flex-direction: column;
}

.page-title {
  font-size: 24px;
  color: #303133;
  margin-bottom: 20px;
  border-bottom: 2px solid #409EFF;
  padding-bottom: 10px;
  display: inline-block;
}

.chat-box {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.message-list {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
}

.empty-messages {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.message-item {
  display: flex;
  margin-bottom: 20px;
  animation: fadeIn 0.3s ease;
}

.message-item.admin-message {
  flex-direction: row-reverse;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-avatar {
  flex-shrink: 0;
}

.message-content {
  max-width: 70%;
  margin: 0 15px;
}

.message-item.admin-message .message-content {
  text-align: right;
}

.message-header {
  margin-bottom: 8px;
}

.sender {
  font-size: 14px;
  color: #909399;
  margin-right: 10px;
}

.time {
  font-size: 12px;
  color: #c0c4cc;
}

.message-bubble {
  display: inline-block;
  padding: 12px 16px;
  border-radius: 8px;
  background: #f5f7fa;
  color: #303133;
  line-height: 1.6;
  max-width: 100%;
  word-break: break-word;
}

.message-item.admin-message .message-bubble {
  background: #ecf5ff;
  color: #409EFF;
}

.input-area {
  border-top: 1px solid #ebeef5;
  padding: 20px;
  background: #fafafa;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 10px;
}

.tip {
  font-size: 12px;
  color: #909399;
}
</style>
