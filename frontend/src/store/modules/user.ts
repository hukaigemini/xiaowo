import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface User {
  id: string
  username: string
  avatar?: string
  isOnline: boolean
  joinTime: Date
}

export const useUserStore = defineStore('user', () => {
  // 状态
  const currentUser = ref<User | null>(null)
  const users = ref<User[]>([])

  // 计算属性
  const isLoggedIn = computed(() => currentUser.value !== null)
  const userList = computed(() => users.value)

  // 初始化用户状态
  const initialize = async () => {
    try {
      // 从本地存储恢复用户状态
      const savedUser = localStorage.getItem('xiaowo_user')
      if (savedUser) {
        currentUser.value = JSON.parse(savedUser)
      }
    } catch (error) {
      console.error('初始化用户状态失败:', error)
    }
  }

  // 设置当前用户
  const setCurrentUser = (user: User) => {
    currentUser.value = user
    // 保存到本地存储
    localStorage.setItem('xiaowo_user', JSON.stringify(user))
  }

  // 添加用户到房间
  const addUser = (user: User) => {
    const existingIndex = users.value.findIndex(u => u.id === user.id)
    if (existingIndex >= 0) {
      users.value[existingIndex] = user
    } else {
      users.value.push(user)
    }
  }

  // 从房间移除用户
  const removeUser = (userId: string) => {
    users.value = users.value.filter(u => u.id !== userId)
  }

  // 更新用户状态
  const updateUser = (userId: string, updates: Partial<User>) => {
    const user = users.value.find(u => u.id === userId)
    if (user) {
      Object.assign(user, updates)
    }
  }

  // 清除用户状态
  const clearUser = () => {
    currentUser.value = null
    users.value = []
    localStorage.removeItem('xiaowo_user')
  }

  // 生成随机用户ID
  const generateUserId = () => {
    return 'user_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
  }

  // 创建临时用户
  const createTempUser = (username: string): User => {
    return {
      id: generateUserId(),
      username,
      isOnline: true,
      joinTime: new Date(),
    }
  }

  return {
    // 状态
    currentUser,
    users,
    // 计算属性
    isLoggedIn,
    userList,
    // 方法
    initialize,
    setCurrentUser,
    addUser,
    removeUser,
    updateUser,
    clearUser,
    createTempUser,
  }
})