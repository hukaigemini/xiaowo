const axios = require('axios');
const config = require('./config');

class ApiClient {
  constructor() {
    this.baseURL = config.getBaseUrl();
    this.timeout = config.getTimeout();
    this.client = axios.create({
      baseURL: this.baseURL,
      timeout: this.timeout,
      headers: {
        'Content-Type': 'application/json',
        'User-Agent': 'Xiaowo-API-Test/1.0.0'
      }
    });

    // 请求拦截器
    this.client.interceptors.request.use(
      (request) => {
        console.log(`🚀 API请求: ${request.method?.toUpperCase()} ${request.url}`);
        return request;
      },
      (error) => {
        console.error('❌ 请求错误:', error.message);
        return Promise.reject(error);
      }
    );

    // 响应拦截器
    this.client.interceptors.response.use(
      (response) => {
        console.log(`✅ API响应: ${response.status} ${response.config.url}`);
        return response;
      },
      (error) => {
        console.error(`❌ API错误: ${error.response?.status} ${error.config?.url}`);
        if (error.response) {
          console.error('错误详情:', error.response.data);
        }
        return Promise.reject(error);
      }
    );
  }

  // 创建会话
  async createSession(userData) {
    const response = await this.client.post('/api/v1/sessions', userData);
    return response.data;
  }

  // 获取会话
  async getSession(sessionId) {
    const response = await this.client.get(`/api/v1/sessions/${sessionId}`);
    return response.data;
  }

  // 更新会话
  async updateSession(sessionId, updates) {
    const response = await this.client.put(`/api/v1/sessions/${sessionId}`, updates);
    return response.data;
  }

  // 心跳
  async heartbeat(sessionId) {
    const response = await this.client.post(`/api/v1/sessions/${sessionId}/heartbeat`);
    return response.data;
  }

  // 验证会话
  async validateSession(sessionId) {
    const response = await this.client.get(`/api/v1/sessions/${sessionId}/validate`);
    return response.data;
  }

  // 删除会话
  async deleteSession(sessionId) {
    const response = await this.client.delete(`/api/v1/sessions/${sessionId}`);
    return response.data;
  }

  // 创建房间
  async createRoom(roomData, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post('/api/v1/rooms', roomData, { headers });
    return response.data;
  }

  // 获取房间列表
  async getRooms(params = {}) {
    const response = await this.client.get('/api/v1/rooms', { params });
    return response.data;
  }

  // 获取房间详情
  async getRoom(roomId) {
    const response = await this.client.get(`/api/v1/rooms/${roomId}`);
    return response.data;
  }

  // 更新房间
  async updateRoom(roomId, updates, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.put(`/api/v1/rooms/${roomId}`, updates, { headers });
    return response.data;
  }

  // 关闭房间
  async closeRoom(roomId, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.delete(`/api/v1/rooms/${roomId}`, { headers });
    return response.data;
  }

  // 获取房间成员
  async getRoomMembers(roomId) {
    const response = await this.client.get(`/api/v1/rooms/${roomId}/members`);
    return response.data;
  }

  // 加入房间
  async joinRoom(roomId, joinData, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post(`/api/v1/rooms/${roomId}/join`, joinData, { headers });
    return response.data;
  }

  // 离开房间
  async leaveRoom(roomId, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post(`/api/v1/rooms/${roomId}/leave`, {}, { headers });
    return response.data;
  }

  // 播放控制
  async playVideo(roomId, playData = {}, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post(`/api/v1/rooms/${roomId}/play`, playData, { headers });
    return response.data;
  }

  async pauseVideo(roomId, pauseData = {}, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post(`/api/v1/rooms/${roomId}/pause`, pauseData, { headers });
    return response.data;
  }

  async seekVideo(roomId, seekData, token = null) {
    const headers = token ? { Authorization: `Bearer ${token}` } : {};
    const response = await this.client.post(`/api/v1/rooms/${roomId}/seek`, seekData, { headers });
    return response.data;
  }

  // 获取播放状态
  async getPlaybackStatus(roomId) {
    const response = await this.client.get(`/api/v1/rooms/${roomId}/status`);
    return response.data;
  }

  // 获取响应时间
  async measureResponseTime(method, url, data = null) {
    const start = Date.now();
    try {
      const response = await this.client.request({
        method,
        url,
        data
      });
      const responseTime = Date.now() - start;
      return { responseTime, status: response.status, success: true };
    } catch (error) {
      const responseTime = Date.now() - start;
      return { 
        responseTime, 
        status: error.response?.status, 
        success: false, 
        error: error.message 
      };
    }
  }
}

module.exports = new ApiClient();
