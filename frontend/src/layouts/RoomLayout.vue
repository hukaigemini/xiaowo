<template>
  <div class="h-screen bg-black flex flex-col overflow-hidden">
    <!-- 顶部工具栏 -->
    <header class="bg-gray-900 border-b border-gray-700 px-4 py-2 flex-shrink-0">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <router-link to="/" class="text-white hover:text-blue-400 transition-colors">
            ← 返回首页
          </router-link>
          <h1 class="text-white font-medium">{{ roomTitle }}</h1>
        </div>
        
        <div class="flex items-center space-x-2">
          <!-- 房间链接分享 -->
          <button
            @click="copyRoomLink"
            class="px-3 py-1 bg-blue-600 text-white text-sm rounded hover:bg-blue-700 transition-colors"
          >
            复制房间链接
          </button>
          
          <!-- 成员数显示 -->
          <div class="text-white text-sm">
            {{ userStore.userList.length }} 人在线
          </div>
        </div>
      </div>
    </header>

    <!-- 主要内容区域 -->
    <div class="flex-1 flex overflow-hidden">
      <!-- 视频播放器区域 -->
      <div class="flex-1 flex flex-col">
        <div class="flex-1 relative">
          <router-view v-slot="{ Component }">
            <transition name="fade" mode="out-in">
              <component :is="Component" />
            </transition>
          </router-view>
        </div>
      </div>
      
      <!-- 右侧成员列表 -->
      <div class="w-80 bg-gray-800 border-l border-gray-700 flex flex-col">
        <div class="p-4 border-b border-gray-700">
          <h2 class="text-white font-medium mb-2">房间成员</h2>
          <div class="text-gray-400 text-sm">
            {{ userStore.userList.length }} 人在线
          </div>
        </div>
        
        <div class="flex-1 overflow-y-auto">
          <div class="p-4 space-y-3">
            <div
              v-for="user in userStore.userList"
              :key="user.id"
              class="flex items-center space-x-3 p-2 rounded hover:bg-gray-700 transition-colors"
            >
              <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center">
                <span class="text-white text-sm font-medium">
                  {{ user.username.charAt(0).toUpperCase() }}
                </span>
              </div>
              <div class="flex-1">
                <div class="text-white text-sm">{{ user.username }}</div>
                <div class="text-gray-400 text-xs">
                  {{ user.isOnline ? '在线' : '离线' }}
                </div>
              </div>
              <div
                class="w-2 h-2 rounded-full"
                :class="user.isOnline ? 'bg-green-500' : 'bg-gray-500'"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/modules/user'
import { useRoomStore } from '@/store/modules/room'

const route = useRoute()
const userStore = useUserStore()
const roomStore = useRoomStore()

const roomTitle = computed(() => {
  const roomId = route.params.roomId
  return roomId ? `房间 ${roomId}` : '观影房间'
})

const copyRoomLink = async () => {
  const roomLink = `${window.location.origin}/room/${route.params.roomId}`
  try {
    await navigator.clipboard.writeText(roomLink)
    // TODO: 添加成功提示
    console.log('房间链接已复制到剪贴板')
  } catch (error) {
    console.error('复制失败:', error)
  }
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>