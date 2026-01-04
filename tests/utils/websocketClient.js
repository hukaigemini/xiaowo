const WebSocket = require('ws');
const config = require('./config');

class WebSocketClient {
  constructor() {
    this.ws = null;
    this.roomId = null;
    this.token = null;
    this.messageHandlers = new Map();
    this.connectionHandlers = [];
    this.disconnectionHandlers = [];
    this.isConnected = false;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.reconnectDelay = 1000;
  }

  // 连接到房间WebSocket
  async connect(roomId, token) {
    return new Promise((resolve, reject) => {
      try {
        this.roomId = roomId;
        this.token = token;
        
        const wsUrl = config.getWebSocketUrl(roomId, token);
        console.log(`🔗 连接到WebSocket: ${wsUrl}`);
        
        this.ws = new WebSocket(wsUrl);
        
        // 连接超时
        const timeout = setTimeout(() => {
          reject(new Error('WebSocket连接超时'));
        }, config.getTimeout());
        
        this.ws.on('open', () => {
          clearTimeout(timeout);
          console.log('✅ WebSocket连接成功');
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.setupEventHandlers();
          resolve(this);
        });
        
        this.ws.on('error', (error) => {
          clearTimeout(timeout);
          console.error('❌ WebSocket连接错误:', error.message);
          reject(error);
        });
        
        this.ws.on('close', (code, reason) => {
          console.log(`🔌 WebSocket连接关闭: ${code} ${reason}`);
          this.isConnected = false;
          this.handleDisconnection();
        });
        
      } catch (error) {
        reject(error);
      }
    });
  }
  
  setupEventHandlers() {
    this.ws.on('message', (data) => {
      try {
        const message = JSON.parse(data.toString());
        console.log(`📨 WebSocket消息:`, message);
        this.handleMessage(message);
      } catch (error) {
        console.error('❌ 解析WebSocket消息失败:', error.message);
      }
    });
  }
  
  handleMessage(message) {
    // 处理不同类型的消息
    const { type, data } = message;
    
    switch (type) {
      case 'room_update':
        this.notifyHandlers('room_update', data);
        break;
      case 'member_joined':
        this.notifyHandlers('member_joined', data);
        break;
      case 'member_left':
        this.notifyHandlers('member_left', data);
        break;
      case 'playback_status':
        this.notifyHandlers('playback_status', data);
        break;
      case 'sync_command':
        this.notifyHandlers('sync_command', data);
        break;
      case 'error':
        this.notifyHandlers('error', data);
        break;
      default:
        console.log('📨 未知消息类型:', type);
    }
  }
  
  // 发送消息
  send(type, data = {}) {
    if (!this.isConnected) {
      throw new Error('WebSocket未连接');
    }
    
    const message = JSON.stringify({ type, data });
    this.ws.send(message);
    console.log(`📤 发送WebSocket消息: ${type}`, data);
  }
  
  // 播放控制消息
  sendPlay(currentTime = 0, position = 0) {
    this.send('play_command', { currentTime, position });
  }
  
  sendPause(currentTime = 0, position = 0) {
    this.send('pause_command', { currentTime, position });
  }
  
  sendSeek(currentTime, position = 0) {
    this.send('seek_command', { currentTime, position });
  }
  
  // 消息处理器
  on(type, handler) {
    if (!this.messageHandlers.has(type)) {
      this.messageHandlers.set(type, []);
    }
    this.messageHandlers.get(type).push(handler);
  }
  
  notifyHandlers(type, data) {
    const handlers = this.messageHandlers.get(type) || [];
    handlers.forEach(handler => {
      try {
        handler(data);
      } catch (error) {
        console.error(`❌ 消息处理器错误 (${type}):`, error.message);
      }
    });
  }
  
  // 连接状态处理器
  onConnect(handler) {
    this.connectionHandlers.push(handler);
  }
  
  onDisconnect(handler) {
    this.disconnectionHandlers.push(handler);
  }
  
  handleDisconnection() {
    this.disconnectionHandlers.forEach(handler => {
      try {
        handler();
      } catch (error) {
        console.error('❌ 断开连接处理器错误:', error.message);
      }
    });
  }
  
  // 自动重连
  async reconnect() {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      throw new Error('达到最大重连次数');
    }
    
    this.reconnectAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1);
    
    console.log(`🔄 尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts}), 延迟 ${delay}ms`);
    
    await new Promise(resolve => setTimeout(resolve, delay));
    
    return this.connect(this.roomId, this.token);
  }
  
  // 断开连接
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
      this.isConnected = false;
    }
  }
  
  // 等待特定消息
  waitForMessage(type, timeout = 5000) {
    return new Promise((resolve, reject) => {
      const timer = setTimeout(() => {
        reject(new Error(`等待消息 ${type} 超时`));
      }, timeout);
      
      const handler = (data) => {
        clearTimeout(timer);
        this.messageHandlers.get(type)?.forEach(h => {
          this.messageHandlers.set(type, this.messageHandlers.get(type).filter(hh => hh !== handler));
        });
        resolve(data);
      };
      
      this.on(type, handler);
    });
  }
  
  // 等待连接建立
  waitForConnect(timeout = 5000) {
    return new Promise((resolve, reject) => {
      if (this.isConnected) {
        resolve();
        return;
      }
      
      const timer = setTimeout(() => {
        reject(new Error('等待连接超时'));
      }, timeout);
      
      this.onConnect(() => {
        clearTimeout(timer);
        resolve();
      });
    });
  }
  
  // 等待断开连接
  waitForDisconnect(timeout = 5000) {
    return new Promise((resolve, reject) => {
      if (!this.isConnected) {
        resolve();
        return;
      }
      
      const timer = setTimeout(() => {
        reject(new Error('等待断开连接超时'));
      }, timeout);
      
      this.onDisconnect(() => {
        clearTimeout(timer);
        resolve();
      });
    });
  }
}

module.exports = WebSocketClient;
