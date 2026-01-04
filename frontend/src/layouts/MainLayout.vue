<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- 主导航栏 -->
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo 和标题 -->
          <div class="flex items-center">
            <router-link to="/" class="flex items-center space-x-2">
              <div class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center">
                <span class="text-white font-bold text-sm">小</span>
              </div>
              <span class="text-xl font-semibold text-gray-900">小窝同步观影</span>
            </router-link>
          </div>
          
          <!-- 用户状态 -->
          <div class="flex items-center space-x-4">
            <div v-if="userStore.isLoggedIn" class="flex items-center space-x-2">
              <span class="text-sm text-gray-700">{{ userStore.currentUser?.username }}</span>
              <div class="w-6 h-6 bg-green-500 rounded-full"></div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- 主内容区域 -->
    <main class="flex-1">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- 页脚 -->
    <footer class="bg-white border-t border-gray-200 py-4">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-sm text-gray-500">
        © 2025 小窝同步观影平台 - 与朋友一起享受观影时光
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
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