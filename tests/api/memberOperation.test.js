const apiClient = require('../utils/apiClient');
const config = require('../utils/config');

describe('成员操作测试 (TC_MO_001-016)', () => {
  let testRoom;
  let testSession;
  let createdSessions = [];
  let joinedRooms = [];

  beforeAll(async () => {
    console.log('🧪 开始成员操作测试');
  });

  afterAll(async () => {
    // 清理测试数据
    console.log('🧹 清理测试数据');
    
    // 离开房间
    for (const roomId of joinedRooms) {
      try {
        await apiClient.leaveRoom(roomId);
      } catch (error) {
        console.log(`⚠️ 离开房间 ${roomId} 失败:`, error.message);
      }
    }
    
    // 删除会话
    for (const session of createdSessions) {
      try {
        await apiClient.deleteSession(session.id);
      } catch (error) {
        console.log(`⚠️ 删除会话 ${session.id} 失败:`, error.message);
      }
    }
  });

  describe('2.1 创建会话测试', () => {
    test('TC_MO_001: 创建会话成功', async () => {
      const userData = config.getTestUser();
      
      const response = await apiClient.createSession(userData);
      
      expect(response.status).toBe('success');
      expect(response.data).toHaveProperty('id');
      expect(response.data).toHaveProperty('token');
      expect(response.data.nickname).toBe(userData.nickname);
      expect(response.data.avatar).toBe(userData.avatar);
      
      testSession = response.data;
      createdSessions.push(testSession);
      
      console.log('✅ 创建会话成功:', testSession.id);
    });

    test('TC_MO_002: 创建会话参数缺失', async () => {
      const invalidUserData = {};
      
      await expect(apiClient.createSession(invalidUserData)).rejects.toThrow();
    });

    test('TC_MO_003: 创建会话nickname长度超限', async () => {
      const invalidUserData = {
        ...config.getTestUser(),
        nickname: 'x'.repeat(51) // 超过50字符限制
      };
      
      await expect(apiClient.createSession(invalidUserData)).rejects.toThrow();
    });
  });

  describe('2.2 加入房间测试', () => {
    beforeEach(async () => {
      // 创建测试房间
      const roomData = config.getTestRoom();
      const roomResponse = await apiClient.createRoom(roomData);
      testRoom = roomResponse.data;
    });

    afterEach(async () => {
      // 清理测试房间
      if (testRoom) {
        try {
          await apiClient.closeRoom(testRoom.id, testRoom.token);
        } catch (error) {
          console.log(`⚠️ 清理房间 ${testRoom.id} 失败:`, error.message);
        }
      }
    });

    test('TC_MO_004: 加入公开房间成功', async () => {
      const joinData = {
        room_id: testRoom.id,
        display_name: '测试用户'
      };
      
      const response = await apiClient.joinRoom(testRoom.id, joinData, testSession.token);
      
      expect(response.status).toBe('success');
      expect(response.data).toHaveProperty('member_id');
      expect(response.data.room_id).toBe(testRoom.id);
      
      joinedRooms.push(testRoom.id);
      console.log('✅ 加入公开房间成功:', testRoom.id);
    });

    test('TC_MO_005: 加入私密房间成功', async () => {
      // 创建私密房间
      const privateRoomData = {
        ...config.getTestRoom(),
        is_private: true,
        password: 'test123'
      };
      
      const privateRoomResponse = await apiClient.createRoom(privateRoomData);
      const privateRoom = privateRoomResponse.data;
      
      const joinData = {
        room_id: privateRoom.id,
        display_name: '测试用户',
        password: 'test123'
      };
      
      const response = await apiClient.joinRoom(privateRoom.id, joinData, testSession.token);
      
      expect(response.status).toBe('success');
      expect(response.data.room_id).toBe(privateRoom.id);
      
      joinedRooms.push(privateRoom.id);
      
      // 清理私密房间
      try {
        await apiClient.closeRoom(privateRoom.id, privateRoom.token);
      } catch (error) {
        console.log(`⚠️ 清理私密房间 ${privateRoom.id} 失败:`, error.message);
      }
    });

    test('TC_MO_006: 加入房间room_id缺失', async () => {
      const joinData = {
        display_name: '测试用户'
      };
      
      await expect(apiClient.joinRoom('', joinData, testSession.token)).rejects.toThrow();
    });

    test('TC_MO_007: 加入不存在房间', async () => {
      const joinData = {
        room_id: 'invalid_room_id',
        display_name: '测试用户'
      };
      
      await expect(apiClient.joinRoom('invalid_room_id', joinData, testSession.token)).rejects.toThrow();
    });

    test('TC_MO_008: 加入私密房间密码错误', async () => {
      const joinData = {
        room_id: testRoom.id,
        display_name: '测试用户',
        password: 'wrong_password'
      };
      
      await expect(apiClient.joinRoom(testRoom.id, joinData, testSession.token)).rejects.toThrow();
    });

    test('TC_MO_009: 加入已满房间', async () => {
      // 创建小容量房间
      const smallRoomData = {
        ...config.getTestRoom(),
        max_users: 1
      };
      
      const smallRoomResponse = await apiClient.createRoom(smallRoomData);
      const smallRoom = smallRoomResponse.data;
      
      // 房间创建者自动加入
      joinedRooms.push(smallRoom.id);
      
      // 尝试加入
      const joinData = {
        room_id: smallRoom.id,
        display_name: '第二个用户'
      };
      
      await expect(apiClient.joinRoom(smallRoom.id, joinData, testSession.token)).rejects.toThrow();
      
      // 清理房间
      try {
        await apiClient.closeRoom(smallRoom.id, smallRoom.token);
      } catch (error) {
        console.log(`⚠️ 清理小房间 ${smallRoom.id} 失败:`, error.message);
      }
    });
  });

  describe('2.3 获取房间成员测试', () => {
    test('TC_MO_010: 获取房间成员成功', async () => {
      // 先加入房间
      const joinData = {
        room_id: testRoom.id,
        display_name: '测试用户'
      };
      
      await apiClient.joinRoom(testRoom.id, joinData, testSession.token);
      joinedRooms.push(testRoom.id);
      
      const response = await apiClient.getRoomMembers(testRoom.id);
      
      expect(response.status).toBe('success');
      expect(Array.isArray(response.data.members)).toBe(true);
      expect(response.data.members.length).toBeGreaterThan(0);
      
      console.log('✅ 获取房间成员成功，成员数:', response.data.members.length);
    });

    test('TC_MO_011: 获取不存在房间成员', async () => {
      await expect(apiClient.getRoomMembers('invalid_room_id')).rejects.toThrow();
    });
  });

  describe('2.4 离开房间测试', () => {
    test('TC_MO_012: 离开房间成功', async () => {
      // 先加入房间
      const joinData = {
        room_id: testRoom.id,
        display_name: '测试用户'
      };
      
      await apiClient.joinRoom(testRoom.id, joinData, testSession.token);
      
      const response = await apiClient.leaveRoom(testRoom.id, testSession.token);
      
      expect(response.status).toBe('success');
      
      // 从joinedRooms中移除
      joinedRooms = joinedRooms.filter(id => id !== testRoom.id);
      console.log('✅ 离开房间成功:', testRoom.id);
    });

    test('TC_MO_013: 未加入房间尝试离开', async () => {
      await expect(apiClient.leaveRoom('different_room_id', testSession.token)).rejects.toThrow();
    });

    test('TC_MO_014: 离开房间room_id缺失', async () => {
      await expect(apiClient.leaveRoom('', testSession.token)).rejects.toThrow();
    });
  });

  describe('2.5 会话管理测试', () => {
    test('TC_MO_015: 获取会话详情成功', async () => {
      const response = await apiClient.getSession(testSession.id);
      
      expect(response.status).toBe('success');
      expect(response.data.id).toBe(testSession.id);
      expect(response.data.nickname).toBe(testSession.nickname);
      
      console.log('✅ 获取会话详情成功:', testSession.id);
    });

    test('TC_MO_016: 会话心跳测试', async () => {
      const response = await apiClient.heartbeat(testSession.id);
      
      expect(response.status).toBe('success');
      expect(response.data).toHaveProperty('last_heartbeat');
      
      console.log('✅ 会话心跳成功:', testSession.id);
    });
  });

  describe('2.6 性能测试', () => {
    test('会话创建响应时间测试', async () => {
      const userData = config.getTestUser();
      
      const result = await apiClient.measureResponseTime('POST', '/api/v1/sessions', userData);
      
      expect(result.success).toBe(true);
      expect(result.responseTime).toBeLessThan(config.getPerformanceThresholds().api_response_time);
      console.log(`✅ 会话创建响应时间: ${result.responseTime}ms`);
    });

    test('加入房间响应时间测试', async () => {
      const joinData = {
        room_id: testRoom.id,
        display_name: '性能测试用户'
      };
      
      const result = await apiClient.measureResponseTime(
        'POST', 
        `/api/v1/rooms/${testRoom.id}/join`, 
        joinData
      );
      
      expect(result.success).toBe(true);
      expect(result.responseTime).toBeLessThan(config.getPerformanceThresholds().api_response_time);
      console.log(`✅ 加入房间响应时间: ${result.responseTime}ms`);
    });
  });
});
