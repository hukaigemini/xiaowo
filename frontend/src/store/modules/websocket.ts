import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface WebSocketMessage {
  type: string
  data: any
  timestamp: number
}

export const useWebSocketStore = defineStore('websocket', () => {
  // 状态
  const socket = ref<WebSocket | null>(null)
  const isConnected = ref(false)
  const isConnecting = ref(false)
  const lastMessage = ref<WebSocketMessage | null>(null)
  const messages = ref<WebSocketMessage[]>([])

  // 连接 WebSocket
  const connect = async (url: string): Promise<void> => {
    if (isConnected.value || isConnecting.value) {
      return
    }

    isConnecting.value = true

    try {
      socket.value = new WebSocket(url)

      socket.value.onopen = () => {
        isConnected.value = true
        isConnecting.value = false
        console.log('WebSocket connected')
      }

      socket.value.onmessage = (event) => {
        try {
          const message: WebSocketMessage = {
            type: 'message',
            data: JSON.parse(event.data),
            timestamp: Date.now()
          }
          lastMessage.value = message
          messages.value.push(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      socket.value.onclose = () => {
        isConnected.value = false
        isConnecting.value = false
        console.log('WebSocket disconnected')
      }

      socket.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        isConnecting.value = false
      }
    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
      isConnecting.value = false
    }
  }

  // 断开连接
  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
    }
    isConnected.value = false
    isConnecting.value = false
  }

  // 发送消息
  const sendMessage = (message: any): boolean => {
    if (!isConnected.value || !socket.value) {
      console.warn('WebSocket is not connected')
      return false
    }

    try {
      socket.value.send(JSON.stringify(message))
      return true
    } catch (error) {
      console.error('Failed to send WebSocket message:', error)
      return false
    }
  }

  return {
    // 状态
    socket,
    isConnected,
    isConnecting,
    lastMessage,
    messages,
    // 方法
    connect,
    disconnect,
    sendMessage,
  }
})