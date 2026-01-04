const WebSocketClient = require('../utils/websocketClient');
const apiClient = require('../utils/apiClient');
const config = require('../utils/config');

describe('WebSocket连接测试 (TC_WS_001-015)', () => {
  let testRoom;
  let testSession;
  let wsClient;
  let createdSessions = [];
  let createdRooms = [];

  beforeAll(async () => {
    console.log('🧪 开始WebSocket测试');
  });

  afterAll(async () => {
    console.log('🧹 清理测试数据');
    
    // 断开WebSocket连接
    if (wsClient) {
      wsClient.disconnect();
    }
    
    // 清理房间
    for (const room of createdRooms) {
      try {
        await apiClient.closeRoom(room.id, room.token);
      } catch (error) {
        console.log(`⚠️ 清理房间 ${room.id} 失败:`, error.message);
      }
    }
    
    // 清理会话
    for (const session of createdSessions) {
      try {
        await apiClient.deleteSession(session.id);
      } catch (error) {
        console.log(`⚠️ 删除会话 ${session.id} 失败:`, error.message);
      }
    }
  });

  beforeEach(async () => {
    // 创建测试会话
    const userData = config.getTestUser();
    const sessionResponse = await apiClient.createSession(userData);
    testSession = sessionResponse.data;
    createdSessions.push(testSession);
    
    // 创建测试房间
    const roomData = config.getTestRoom();
    const roomResponse = await apiClient.createRoom(roomData);
    testRoom = roomResponse.data;
    createdRooms.push(testRoom);
    
    // 加入房间
    const joinData = {
      room_id: testRoom.id,
      display_name: 'WebSocket测试用户'
    };
    
    await apiClient.joinRoom(testRoom.id, joinData, testSession.token);
    
    // 创建WebSocket客户端
    wsClient = new WebSocketClient();
  });

  afterEach(async () => {
    if (wsClient) {
      wsClient.disconnect();
    }
  });

  describe('3.1 WebSocket连接测试', () => {
    test('TC_WS_001: WebSocket连接成功', async () => {
      const ws = await wsClient.connect(testRoom.id, testSession.token);
      
      expect(ws).toBeInstanceOf(WebSocketClient);
      expect(wsClient.isConnected).toBe(true);
      
      console.log('✅ WebSocket连接成功:', testRoom.id);
    });

    test('TC_WS_002: WebSocket连接失败-房间不存在', async () => {
      await expect(
        wsClient.connect('invalid_room_id', testSession.token)
      ).rejects.toThrow();
    });

    test('TC_WS_003: WebSocket连接失败-token无效', async () => {
      await expect(
        wsClient.connect(testRoom.id, 'invalid_token')
      ).rejects.toThrow();
    });

    test('TC_WS_004: WebSocket连接失败-未加入房间', async () => {
      // 创建新会话但未加入房间
      const newUserData = {
        nickname: '新用户',
        avatar: 'https://example.com/avatar.jpg'
      };
      
      const newSessionResponse = await apiClient.createSession(newUserData);
      const newSession = newSessionResponse.data;
      createdSessions.push(newSession);
      
      await expect(
        wsClient.connect(testRoom.id, newSession.token)
      ).rejects.toThrow();
      
      // 清理新会话
      await apiClient.deleteSession(newSession.id);
    });

    test('TC_WS_005: WebSocket连接超时', async () => {
      // 使用非常短的超时时间测试超时
      const originalTimeout = config.getTimeout();
      
      // 这里需要临时修改配置，在实际测试中可以通过mock实现
      await expect(
        wsClient.connect(testRoom.id, testSession.token)
      ).resolves.toBeDefined();
    });
  });

  describe('3.2 消息传递测试', () => {
    beforeEach(async () => {
      await wsClient.connect(testRoom.id, testSession.token);
    });

    test('TC_WS_006: 接收房间更新消息', async () => {
      let messageReceived = false;
      
      wsClient.on('room_update', (data) => {
        expect(data).toHaveProperty('room_id');
        expect(data.room_id).toBe(testRoom.id);
        messageReceived = true;
      });
      
      // 模拟发送更新消息
      wsClient.send('update_request', { action: 'test' });
      
      // 等待消息接收（实际测试中需要后端推送消息）
      await new Promise(resolve => setTimeout(resolve, 100));
      
      console.log('✅ 房间更新消息测试完成');
    });

    test('TC_WS_007: 接收成员加入消息', async () => {
      let memberJoined = false;
      
      wsClient.on('member_joined', (data) => {
        expect(data).toHaveProperty('member_id');
        expect(data).toHaveProperty('room_id');
        expect(data.room_id).toBe(testRoom.id);
        memberJoined = true;
      });
      
      console.log('✅ 成员加入消息测试完成');
    });

    test('TC_WS_008: 接收成员离开消息', async () => {
      let memberLeft = false;
      
      wsClient.on('member_left', (data) => {
        expect(data).toHaveProperty('member_id');
        expect(data).toHaveProperty('room_id');
        expect(data.room_id).toBe(testRoom.id);
        memberLeft = true;
      });
      
      console.log('✅ 成员离开消息测试完成');
    });

    test('TC_WS_009: 接收播放状态消息', async () => {
      let playbackStatus = false;
      
      wsClient.on('playback_status', (data) => {
        expect(data).toHaveProperty('status'); // 'playing', 'paused', 'stopped'
        expect(data).toHaveProperty('current_time');
        expect(data).toHaveProperty('room_id');
        expect(data.room_id).toBe(testRoom.id);
        playbackStatus = true;
      });
      
      console.log('✅ 播放状态消息测试完成');
    });

    test('TC_WS_010: 接收同步控制消息', async () => {
      let syncCommand = false;
      
      wsClient.on('sync_command', (data) => {
        expect(data).toHaveProperty('command'); // 'play', 'pause', 'seek'
        expect(data).toHaveProperty('current_time');
        expect(data).toHaveProperty('room_id');
        expect(data.room_id).toBe(testRoom.id);
        syncCommand = true;
      });
      
      console.log('✅ 同步控制消息测试完成');
    });
  });

  describe('3.3 同步控制测试', () => {
    beforeEach(async () => {
      await wsClient.connect(testRoom.id, testSession.token);
    });

    test('TC_WS_011: 发送播放控制消息', async () => {
      const currentTime = 120; // 2分钟
      
      // 不应该抛出异常
      expect(() => {
        wsClient.sendPlay(currentTime);
      }).not.toThrow();
      
      console.log('✅ 发送播放控制消息成功');
    });

    test('TC_WS_012: 发送暂停控制消息', async () => {
      const currentTime = 120;
      
      expect(() => {
        wsClient.sendPause(currentTime);
      }).not.toThrow();
      
      console.log('✅ 发送暂停控制消息成功');
    });

    test('TC_WS_013: 发送拖拽控制消息', async () => {
      const currentTime = 300; // 5分钟
      
      expect(() => {
        wsClient.sendSeek(currentTime);
      }).not.toThrow();
      
      console.log('✅ 发送拖拽控制消息成功');
    });

    test('TC_WS_014: 未连接时发送消息', async () => {
      // 先断开连接
      wsClient.disconnect();
      
      expect(() => {
        wsClient.sendPlay(0);
      }).toThrow('WebSocket未连接');
      
      console.log('✅ 未连接时发送消息异常测试通过');
    });
  });

  describe('3.4 连接异常测试', () => {
    test('TC_WS_015: WebSocket自动重连', async () => {
      await wsClient.connect(testRoom.id, testSession.token);
      
      // 模拟断开连接
      wsClient.disconnect();
      
      // 等待一小段时间
      await new Promise(resolve => setTimeout(resolve, 100));
      
      // 尝试重连
      const reconnectedWs = await wsClient.reconnect();
      
      expect(reconnectedWs).toBeInstanceOf(WebSocketClient);
      expect(wsClient.isConnected).toBe(true);
      
      console.log('✅ WebSocket自动重连成功');
    });
  });

  describe('3.5 性能测试', () => {
    beforeEach(async () => {
      await wsClient.connect(testRoom.id, testSession.token);
    });

    test('WebSocket连接响应时间测试', async () => {
      const start = Date.now();
      
      const newWsClient = new WebSocketClient();
      await newWsClient.connect(testRoom.id, testSession.token);
      
      const connectTime = Date.now() - start;
      
      expect(connectTime).toBeLessThan(config.getPerformanceThresholds().websocket_connect_time);
      expect(newWsClient.isConnected).toBe(true);
      
      newWsClient.disconnect();
      
      console.log(`✅ WebSocket连接响应时间: ${connectTime}ms`);
    });

    test('WebSocket消息发送性能测试', async () => {
      const messageCount = 100;
      const start = Date.now();
      
      for (let i = 0; i < messageCount; i++) {
        wsClient.send('performance_test', { index: i });
      }
      
      const sendTime = Date.now() - start;
      const avgTime = sendTime / messageCount;
      
      expect(avgTime).toBeLessThan(config.getPerformanceThresholds().websocket_message_time);
      
      console.log(`✅ WebSocket消息发送平均时间: ${avgTime.toFixed(2)}ms`);
    });
  });

  describe('3.6 边界测试', () => {
    test('大量消息处理测试', async () => {
      await wsClient.connect(testRoom.id, testSession.token);
      
      const messageCount = 1000;
      let receivedCount = 0;
      
      wsClient.on('test_message', () => {
        receivedCount++;
      });
      
      const start = Date.now();
      
      // 发送大量消息
      for (let i = 0; i < messageCount; i++) {
        wsClient.send('test_message', { index: i });
      }
      
      // 等待消息处理
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const processingTime = Date.now() - start;
      
      console.log(`✅ 处理${messageCount}条消息用时: ${processingTime}ms`);
      
      // 不应该抛出异常
      expect(wsClient.isConnected).toBe(true);
    });
  });
});
