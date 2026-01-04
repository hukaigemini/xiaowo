<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-bg text-text-main font-sans overflow-hidden selection:bg-primary selection:text-white relative">
    
    <!-- 装饰背景：柔和的漫反射光晕 -->
    <div class="fixed inset-0 pointer-events-none z-0 overflow-hidden">
      <div class="absolute top-[-10%] left-[-10%] w-[600px] h-[600px] bg-primary/10 rounded-full blur-[100px]"></div>
      <div class="absolute bottom-[-10%] right-[-10%] w-[500px] h-[500px] bg-secondary/10 rounded-full blur-[80px]"></div>
    </div>

    <!-- 首页：欢迎卡片 -->
    <div class="relative z-10 w-full max-w-md px-6 animate-fade-in">
      
      <!-- 主卡片 -->
      <div class="bg-surface rounded-4xl p-8 shadow-card border-2 border-white/50 relative overflow-hidden mt-8">
        
        <!-- 顶部装饰线 -->
        <div class="absolute top-0 left-0 w-full h-2 bg-gradient-to-r from-primary/20 via-primary/40 to-primary/20"></div>

        <!-- 标题区 -->
        <div class="text-center mb-10 mt-2">
          <div class="inline-flex items-center justify-center w-20 h-20 bg-white rounded-full shadow-lg mb-6 transform hover:scale-105 hover:shadow-xl transition-all duration-300 cursor-default">
            <svg class="w-10 h-10 text-primary" viewBox="0 0 24 24" fill="currentColor">
              <!-- 小窝logo：舒适的沙发 -->
              <!-- 沙发主体 -->
              <path d="M3 14c0-3.5 2.5-6 7-6h4c4.5 0 7 2.5 7 6v2H3v-2z"/>
              <!-- 沙发靠背 -->
              <rect x="2" y="6" width="20" height="10" rx="2" fill="currentColor"/>
              <!-- 沙发扶手 -->
              <rect x="1" y="10" width="4" height="6" rx="1" fill="currentColor"/>
              <rect x="19" y="10" width="4" height="6" rx="1" fill="currentColor"/>
              <!-- 沙发腿 -->
              <rect x="4" y="16" width="2" height="3" rx="1" fill="currentColor" opacity="0.7"/>
              <rect x="18" y="16" width="2" height="3" rx="1" fill="currentColor" opacity="0.7"/>
              <!-- 温馨装饰：坐垫褶皱 -->
              <path d="M8 10h8M8 12h6" stroke="currentColor" stroke-width="1" opacity="0.6"/>
            </svg>
          </div>
          <h1 class="text-3xl font-bold tracking-tight text-text-main mb-2">小窝</h1>
          <p class="text-text-muted font-medium px-4">"再远，也要窝在一起"</p>
        </div>

        <!-- MVP阶段移除快速回房提示 -->

        <!-- 核心操作区 -->
        <div class="space-y-6">
          
          <!-- 创建房间 -->
          <button 
            @click="createRoom"
            class="w-full py-4 bg-[#FF9F76] hover:bg-[#FF8C5A] text-white font-bold text-lg rounded-2xl shadow-[0_8px_20px_-6px_rgba(255,159,118,0.5)] hover:shadow-[0_12px_25px_-8px_rgba(255,159,118,0.6)] hover:-translate-y-1 transition-all duration-300 active:scale-95 flex items-center justify-center space-x-2 group cursor-pointer"
          >
            <span>✨ 搭个温馨小窝</span>
            <svg class="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
            </svg>
          </button>

          <!-- 分隔符 -->
          <div class="relative flex py-2 items-center">
            <div class="flex-grow border-t-2 border-[#E6D5B8]"></div>
            <span class="flex-shrink-0 mx-4 text-text-muted text-sm font-medium">或</span>
            <div class="flex-grow border-t-2 border-[#E6D5B8]"></div>
          </div>

          <!-- 加入房间 -->
          <div class="space-y-3">
            <div class="relative group mt-1 mb-1">
              <input 
                v-model="roomId"
                type="text" 
                placeholder="输入暗号 (房间号)..." 
                class="w-full px-6 py-3 bg-bg rounded-2xl border-2 border-transparent focus:border-primary focus:outline-none text-lg text-text-main placeholder-[#C5B4A0] transition-all shadow-inner-warm text-center font-medium tracking-widest"
                @keyup.enter="joinRoom"
              >
            </div>
            <button 
              @click="joinRoom"
              :disabled="!roomId.trim()"
              class="w-full py-3 bg-white border-2 border-[#E6D5B8] text-text-muted font-bold rounded-2xl hover:bg-[#FFFBF0] hover:text-text-main hover:border-primary/30 transition-all duration-300 active:scale-95 shadow-sm disabled:opacity-50 disabled:cursor-not-allowed cursor-pointer"
            >
              敲门进入
            </button>
          </div>
        </div>
      </div>

      <!-- 底部信息 -->
      <div class="mt-8 text-center">
        <p class="text-text-muted/60 text-xs font-medium tracking-wider mb-4">🔒 只有你们知道这里。就像在家里关上门一样安全。</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useRoomStore } from '@/store/modules/room'
import { useUserStore } from '@/store/modules/user'

const router = useRouter()
const roomStore = useRoomStore()
const userStore = useUserStore()

// 响应式数据
const roomId = ref('')

// 初始化用户状态
onMounted(async () => {
  await userStore.initialize()
})

// 进入房间
const enterRoom = async (targetRoomId?: string) => {
  try {
    // 保存房间信息到LocalStorage
    const finalRoomId = targetRoomId || roomId.value || generateRoomId()
    localStorage.setItem('xiaowo_room', JSON.stringify({
      roomId: finalRoomId,
      lastVisit: Date.now()
    }))
    
    // 路由跳转
    await router.push(`/room/${finalRoomId}`)
  } catch (error) {
    console.error('进入房间失败:', error)
  }
}

// 创建房间
const createRoom = async () => {
  try {
    await userStore.initialize()
    if (!userStore.currentUser) {
      // 如果没有用户，先创建用户
      await userStore.createTempUser('观众' + Math.floor(Math.random() * 1000))
    }
    
    const newRoom = await roomStore.createRoom()
    await enterRoom(newRoom.id)
  } catch (error) {
    console.error('创建房间失败:', error)
  }
}

// 加入房间
const joinRoom = async () => {
  if (!roomId.value.trim()) return
  
  try {
    await userStore.initialize()
    if (!userStore.currentUser) {
      await userStore.createTempUser('观众' + Math.floor(Math.random() * 1000))
    }
    
    await roomStore.joinRoom(roomId.value.trim())
    await enterRoom(roomId.value.trim())
  } catch (error) {
    console.error('加入房间失败:', error)
  }
}

// 生成随机房间号
const generateRoomId = (): string => {
  const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789'
  let result = ''
  for (let i = 0; i < 6; i++) {
    result += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  return result
}
</script>

<style scoped>
/* 自定义动画 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    transform: translateY(10px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

@keyframes pulseWarm {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-fade-in {
  animation: fadeIn 0.5s ease-in-out;
}

.animate-slide-up {
  animation: slideUp 0.3s ease-out;
}

.animate-pulse-warm {
  animation: pulseWarm 2s infinite;
}
</style>