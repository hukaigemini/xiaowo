const apiClient = require('../utils/apiClient');
const config = require('../utils/config');

describe('房间管理测试 (TC_RM_001-013)', () => {
  let testRoom;
  let testRoomToken;
  let createdRooms = [];

  beforeAll(async () => {
    console.log('🧪 开始房间管理测试');
  });

  afterAll(async () => {
    // 清理测试数据
    console.log('🧹 清理测试数据');
    for (const room of createdRooms) {
      try {
        await apiClient.closeRoom(room.id, room.token);
      } catch (error) {
        console.log(`⚠️ 清理房间 ${room.id} 失败:`, error.message);
      }
    }
  });

  describe('1.1 创建房间测试', () => {
    test('TC_RM_001: 创建公开房间成功', async () => {
      const roomData = config.getTestRoom();
      
      const response = await apiClient.createRoom(roomData);
      
      expect(response.status).toBe('success');
      expect(response.data).toHaveProperty('id');
      expect(response.data).toHaveProperty('token');
      expect(response.data.name).toBe(roomData.name);
      expect(response.data.is_private).toBe(roomData.is_private);
      
      testRoom = response.data;
      testRoomToken = response.data.token;
      createdRooms.push(testRoom);
      
      console.log('✅ 创建公开房间成功:', testRoom.id);
    });

    test('TC_RM_002: 创建私密房间成功', async () => {
      const privateRoomData = {
        ...config.getTestRoom(),
        is_private: true,
        password: 'test123'
      };
      
      const response = await apiClient.createRoom(privateRoomData);
      
      expect(response.status).toBe('success');
      expect(response.data.is_private).toBe(true);
      expect(response.data).toHaveProperty('token');
      
      createdRooms.push(response.data);
      console.log('✅ 创建私密房间成功:', response.data.id);
    });

    test('TC_RM_003: 创建房间参数缺失', async () => {
      const invalidRoomData = {
        description: '缺少name和media_url'
      };
      
      await expect(apiClient.createRoom(invalidRoomData)).rejects.toThrow();
    });

    test('TC_RM_004: 创建房间name长度超限', async () => {
      const invalidRoomData = {
        ...config.getTestRoom(),
        name: 'x'.repeat(101) // 超过100字符限制
      };
      
      await expect(apiClient.createRoom(invalidRoomData)).rejects.toThrow();
    });

    test('TC_RM_005: 创建房间max_users越界', async () => {
      const invalidRoomData = {
        ...config.getTestRoom(),
        max_users: 0
      };
      
      await expect(apiClient.createRoom(invalidRoomData)).rejects.toThrow();
    });

    test('TC_RM_006: 创建私密房间未提供密码', async () => {
      const invalidRoomData = {
        ...config.getTestRoom(),
        is_private: true,
        password: ''
      };
      
      await expect(apiClient.createRoom(invalidRoomData)).rejects.toThrow();
    });
  });

  describe('1.2 获取房间信息测试', () => {
    test('TC_RM_007: 获取房间详情成功', async () => {
      const response = await apiClient.getRoom(testRoom.id);
      
      expect(response.status).toBe('success');
      expect(response.data.id).toBe(testRoom.id);
      expect(response.data).toHaveProperty('name');
      expect(response.data).toHaveProperty('media_url');
      expect(response.data).toHaveProperty('created_at');
      
      console.log('✅ 获取房间详情成功:', testRoom.id);
    });

    test('TC_RM_008: 获取不存在房间', async () => {
      await expect(apiClient.getRoom('invalid_room_id')).rejects.toThrow();
    });

    test('TC_RM_009: 获取房间列表成功', async () => {
      const response = await apiClient.getRooms({ page: 1, limit: 10 });
      
      expect(response.status).toBe('success');
      expect(response.data).toHaveProperty('rooms');
      expect(response.data).toHaveProperty('pagination');
      expect(Array.isArray(response.data.rooms)).toBe(true);
      
      console.log('✅ 获取房间列表成功');
    });

    test('TC_RM_010: 分页参数验证', async () => {
      await expect(apiClient.getRooms({ page: -1, limit: 1000 })).rejects.toThrow();
    });
  });

  describe('1.3 更新房间测试', () => {
    test('TC_RM_011: 房间创建者更新成功', async () => {
      const updates = {
        name: '更新后的房间名称',
        description: '更新后的房间描述'
      };
      
      const response = await apiClient.updateRoom(testRoom.id, updates, testRoomToken);
      
      expect(response.status).toBe('success');
      expect(response.data.name).toBe(updates.name);
      expect(response.data.description).toBe(updates.description);
      
      console.log('✅ 更新房间成功:', testRoom.id);
    });

    test('TC_RM_012: 非创建者尝试更新', async () => {
      const updates = { name: '尝试修改' };
      
      await expect(apiClient.updateRoom(testRoom.id, updates)).rejects.toThrow();
    });

    test('TC_RM_013: 更新不存在房间', async () => {
      const updates = { name: '更新不存在房间' };
      
      await expect(apiClient.updateRoom('invalid_room_id', updates)).rejects.toThrow();
    });
  });

  describe('1.4 删除房间测试', () => {
    test('TC_RM_014: 删除房间成功', async () => {
      const roomToDelete = createdRooms.find(r => r.id !== testRoom.id);
      
      if (roomToDelete) {
        const response = await apiClient.closeRoom(roomToDelete.id, roomToDelete.token);
        
        expect(response.status).toBe('success');
        console.log('✅ 删除房间成功:', roomToDelete.id);
      }
    });

    test('TC_RM_015: 非创建者尝试删除', async () => {
      await expect(apiClient.closeRoom(testRoom.id)).rejects.toThrow();
    });

    test('TC_RM_016: 删除不存在房间', async () => {
      await expect(apiClient.closeRoom('invalid_room_id')).rejects.toThrow();
    });
  });

  describe('1.5 性能测试', () => {
    test('房间创建响应时间测试', async () => {
      const roomData = config.getTestRoom();
      
      const result = await apiClient.measureResponseTime('POST', '/api/v1/rooms', roomData);
      
      expect(result.success).toBe(true);
      expect(result.responseTime).toBeLessThan(config.getPerformanceThresholds().api_response_time);
      console.log(`✅ 房间创建响应时间: ${result.responseTime}ms`);
    });
  });
});
