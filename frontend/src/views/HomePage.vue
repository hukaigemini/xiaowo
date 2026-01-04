<template>
  <div class="home-page min-h-screen flex flex-col items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
    <div class="container max-w-4xl mx-auto px-6 text-center">
      <!-- Logo和标题 -->
      <div class="mb-12 animate-fade-in">
        <div class="w-20 h-20 mx-auto mb-6 bg-primary-600 rounded-full flex items-center justify-center">
          <svg class="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
          </svg>
        </div>
        <h1 class="text-5xl font-bold text-gray-900 mb-4">小窝同步观影</h1>
        <p class="text-xl text-gray-600 mb-8">与朋友一起享受同步观影的美好时光</p>
      </div>

      <!-- 操作按钮 -->
      <div class="flex flex-col sm:flex-row gap-4 justify-center items-center mb-12 animate-slide-up">
        <button 
          @click="createRoom"
          :disabled="isCreating"
          class="btn-primary text-lg px-8 py-4 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isCreating" class="flex items-center">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            创建中...
          </span>
          <span v-else class="flex items-center">
            <svg class="w-6 h-6 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
            </svg>
            创建房间
          </span>
        </button>
        
        <button 
          @click="showJoinDialog = true"
          class="btn-outline text-lg px-8 py-4"
        >
          <svg class="w-6 h-6 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
          </svg>
          加入房间
        </button>
      </div>

      <!-- 特性介绍 -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div class="card text-center">
          <div class="w-12 h-12 mx-auto mb-4 bg-primary-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">实时同步</h3>
          <p class="text-gray-600">所有成员观看进度完全同步，打造沉浸式观影体验</p>
        </div>
        
        <div class="card text-center">
          <div class="w-12 h-12 mx-auto mb-4 bg-primary-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8h2a2 2 0 012 2v6a2 2 0 01-2 2h-2v4l-4-4H9a1.994 1.994 0 01-1.414-.586m0 0L11 14h4a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2v4l.586-.586z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">实时聊天</h3>
          <p class="text-gray-600">边看边聊，与朋友分享观影心得和即时互动</p>
        </div>
        
        <div class="card text-center">
          <div class="w-12 h-12 mx-auto mb-4 bg-primary-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z"></path>
            </svg>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">多平台支持</h3>
          <p class="text-gray-600">支持电脑、平板、手机，随时随地享受同步观影</p>
        </div>
      </div>
    </div>

    <!-- 加入房间对话框 -->
    <div v-if="showJoinDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showJoinDialog = false">
      <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4 animate-slide-up">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">加入房间</h3>
        <form @submit.prevent="joinRoom">
          <div class="mb-4">
            <label for="roomId" class="block text-sm font-medium text-gray-700 mb-2">房间号</label>
            <input 
              id="roomId"
              v-model="joinForm.roomId"
              type="text" 
              placeholder="请输入房间号"
              class="input"
              required
            >
          </div>
          <div class="mb-4">
            <label for="username" class="block text-sm font-medium text-gray-700 mb-2">昵称</label>
            <input 
              id="username"
              v-model="joinForm.username"
              type="text" 
              placeholder="请输入你的昵称"
              class="input"
              required
            >
          </div>
          <div class="flex gap-3 justify-end">
            <button 
              type="button"
              @click="showJoinDialog = false"
              class="btn-secondary"
            >
              取消
            </button>
            <button 
              type="submit"
              :disabled="isJoining"
              class="btn-primary disabled:opacity-50"
            >
              {{ isJoining ? '加入中...' : '加入房间' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { useRoomStore } from '@/store/modules/room'

const router = useRouter()
const userStore = useUserStore()
const roomStore = useRoomStore()

// 响应式数据
const isCreating = ref(false)
const isJoining = ref(false)
const showJoinDialog = ref(false)

const joinForm = ref({
  roomId: '',
  username: ''
})

// 方法
const createRoom = async () => {
  isCreating.value = true
  try {
    const room = await roomStore.createRoom()
    router.push(`/room/${room.id}`)
  } catch (error) {
    console.error('创建房间失败:', error)
    // TODO: 显示错误提示
  } finally {
    isCreating.value = false
  }
}

const joinRoom = async () => {
  isJoining.value = true
  try {
    await roomStore.joinRoom(joinForm.value.roomId)
    router.push(`/room/${joinForm.value.roomId}`)
  } catch (error) {
    console.error('加入房间失败:', error)
    // TODO: 显示错误提示
  } finally {
    isJoining.value = false
    showJoinDialog.value = false
  }
}
</script>

<style scoped>
.home-page {
  min-height: 100vh;
}
</style>