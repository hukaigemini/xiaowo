import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useUserStore } from './user'

export interface Video {
  url: string
  name: string
  duration?: number
  currentTime?: number
}

export interface Room {
  id: string
  name: string
  creator: string
  maxMembers: number
  currentMembers: number
  createdAt: Date
  isActive: boolean
}

export interface RoomMember {
  userId: string
  username: string
  avatar?: string
  role: 'owner' | 'member'
  joinTime: Date
  lastSeen: Date
}

export const useRoomStore = defineStore('room', () => {
  // 状态
  const currentRoom = ref<Room | null>(null)
  const members = ref<RoomMember[]>([])
  const isLoading = ref(false)
  const currentVideo = ref<Video | null>(null)
  const isPlaying = ref(false)
  const playProgress = ref(0)

  // 计算属性
  const roomId = computed(() => currentRoom.value?.id || '')
  const roomName = computed(() => currentRoom.value?.name || '')
  const memberCount = computed(() => members.value.length)
  const isOwner = computed(() => {
    const userStore = useUserStore()
    return currentRoom.value?.creator === userStore.currentUser?.id
  })
  const memberList = computed(() => members.value)

  // 创建房间
  const createRoom = async (): Promise<Room> => {
    isLoading.value = true
    try {
      // 模拟API调用
      const room: Room = {
        id: 'room_' + Date.now(),
        name: '我的房间',
        creator: useUserStore().currentUser?.id || '',
        maxMembers: 7,
        currentMembers: 1,
        createdAt: new Date(),
        isActive: true,
      }
      
      currentRoom.value = room
      // 模拟添加创建者为第一个成员
      const userStore = useUserStore()
      if (userStore.currentUser) {
        members.value = [{
          userId: userStore.currentUser.id,
          username: userStore.currentUser.username,
          role: 'owner',
          joinTime: new Date(),
          lastSeen: new Date(),
        }]
      }
      
      return room
    } finally {
      isLoading.value = false
    }
  }

  // 加入房间
  const joinRoom = async (roomId: string) => {
    isLoading.value = true
    try {
      // 模拟API调用
      const room: Room = {
        id: roomId,
        name: '房间 ' + roomId,
        creator: 'user_owner',
        maxMembers: 7,
        currentMembers: 2,
        createdAt: new Date(),
        isActive: true,
      }
      
      currentRoom.value = room
      // 添加当前用户为成员
      const userStore = useUserStore()
      if (userStore.currentUser) {
        members.value.push({
          userId: userStore.currentUser.id,
          username: userStore.currentUser.username,
          role: 'member',
          joinTime: new Date(),
          lastSeen: new Date(),
        })
      }
      
    } finally {
      isLoading.value = false
    }
  }

  // 离开房间
  const leaveRoom = () => {
    currentRoom.value = null
    members.value = []
  }

  // 添加成员
  const addMember = (member: RoomMember) => {
    const existingIndex = members.value.findIndex(m => m.userId === member.userId)
    if (existingIndex >= 0) {
      members.value[existingIndex] = member
    } else {
      members.value.push(member)
    }
  }

  // 移除成员
  const removeMember = (userId: string) => {
    members.value = members.value.filter(m => m.userId !== userId)
  }

  // 更新成员状态
  const updateMember = (userId: string, updates: Partial<RoomMember>) => {
    const member = members.value.find(m => m.userId === userId)
    if (member) {
      Object.assign(member, updates)
    }
  }

  // 获取成员
  const getMember = (userId: string) => {
    return members.value.find(m => m.userId === userId)
  }

  // 视频相关方法
  const setCurrentVideo = (video: Video) => {
    currentVideo.value = video
    playProgress.value = 0
    isPlaying.value = false
  }

  const updatePlayProgress = (time: number) => {
    playProgress.value = time
    if (currentVideo.value) {
      currentVideo.value.currentTime = time
    }
  }

  const setPlaying = (playing: boolean) => {
    isPlaying.value = playing
  }

  // 初始化
  const initialize = () => {
    // 清理状态
    currentRoom.value = null
    members.value = []
    currentVideo.value = null
    isPlaying.value = false
    playProgress.value = 0
  }

  return {
    // 状态
    currentRoom,
    members,
    isLoading,
    currentVideo,
    isPlaying,
    playProgress,
    // 计算属性
    roomId,
    roomName,
    memberCount,
    isOwner,
    memberList,
    // 方法
    createRoom,
    joinRoom,
    leaveRoom,
    addMember,
    removeMember,
    updateMember,
    getMember,
    setCurrentVideo,
    updatePlayProgress,
    setPlaying,
    initialize,
  }
})