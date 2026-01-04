/**
 * 小沃API接口冒烟测试
 * 作者: 稳当 (SRE)
 * 用途: 快速验证核心功能是否正常工作
 */

const apiClient = require('./utils/apiClient');
const config = require('./utils/config');

describe('冒烟测试 - 核心功能验证', () => {
    let testRoom;
    let testSession;

    beforeAll(async () => {
        // 设置测试超时时间
        jest.setTimeout(30000);
        
        // 验证环境配置
        expect(config).toBeDefined();
        expect(config.getApiUrl()).toBeDefined();
        
        console.log('✅ 测试环境验证通过');
    });

    describe('API连接性测试', () => {
        test('TC_SMOKE_001: API服务可访问', async () => {
            const response = await apiClient.healthCheck();
            expect(response.status).toBe('success');
            console.log('✅ API服务连接正常');
        });

        test('TC_SMOKE_002: API响应时间正常', async () => {
            const startTime = Date.now();
            await apiClient.healthCheck();
            const responseTime = Date.now() - startTime;
            
            expect(responseTime).toBeLessThan(5000); // 5秒内响应
            console.log(`✅ API响应时间: ${responseTime}ms`);
        });
    });

    describe('房间管理核心功能', () => {
        test('TC_SMOKE_003: 能够创建公开房间', async () => {
            const roomData = {
                name: '冒烟测试房间',
                media_url: 'https://example.com/test-video.mp4',
                is_private: false
            };

            const response = await apiClient.createRoom(roomData);
            expect(response.status).toBe('success');
            expect(response.data).toHaveProperty('id');
            expect(response.data.name).toBe(roomData.name);
            
            testRoom = response.data;
            console.log('✅ 房间创建成功');
        });

        test('TC_SMOKE_004: 能够获取房间信息', async () => {
            const response = await apiClient.getRoom(testRoom.id);
            expect(response.status).toBe('success');
            expect(response.data.id).toBe(testRoom.id);
            console.log('✅ 房间信息获取正常');
        });
    });

    describe('会话管理核心功能', () => {
        test('TC_SMOKE_005: 能够创建用户会话', async () => {
            const userData = {
                nickname: '冒烟测试用户',
                avatar: 'https://example.com/avatar.jpg'
            };

            const response = await apiClient.createSession(userData);
            expect(response.status).toBe('success');
            expect(response.data).toHaveProperty('token');
            
            testSession = response.data;
            console.log('✅ 用户会话创建成功');
        });

        test('TC_SMOKE_006: 能够加入房间', async () => {
            const joinData = {
                room_id: testRoom.id,
                display_name: '冒烟测试用户'
            };

            const response = await apiClient.joinRoom(testRoom.id, joinData, testSession.token);
            expect(response.status).toBe('success');
            expect(response.data.room_id).toBe(testRoom.id);
            console.log('✅ 用户加入房间成功');
        });
    });

    describe('WebSocket连接核心功能', () => {
        let wsClient;

        beforeAll(async () => {
            const WebSocketClient = require('./utils/websocketClient');
            wsClient = new WebSocketClient();
        });

        test('TC_SMOKE_007: WebSocket能够连接', async () => {
            await expect(wsClient.connect(testRoom.id, testSession.token))
                .resolves.toBeDefined();
            
            expect(wsClient.isConnected).toBe(true);
            console.log('✅ WebSocket连接成功');
        });

        test('TC_SMOKE_008: 能够发送播放控制消息', async () => {
            expect(() => wsClient.sendPlay(0, 0)).not.toThrow();
            console.log('✅ 播放控制消息发送正常');
        });

        afterAll(async () => {
            if (wsClient && wsClient.isConnected) {
                wsClient.disconnect();
            }
        });
    });

    describe('数据一致性验证', () => {
        test('TC_SMOKE_009: 房间成员数量正确', async () => {
            const response = await apiClient.getRoomMembers(testRoom.id);
            expect(response.status).toBe('success');
            expect(response.data.members).toContainEqual(
                expect.objectContaining({
                    display_name: '冒烟测试用户'
                })
            );
            console.log('✅ 房间成员数据一致');
        });

        test('TC_SMOKE_010: 会话状态保持', async () => {
            const response = await apiClient.getSession(testSession.token);
            expect(response.status).toBe('success');
            expect(response.data.token).toBe(testSession.token);
            console.log('✅ 会话状态保持正常');
        });
    });

    describe('错误处理基础验证', () => {
        test('TC_SMOKE_011: 无效房间ID返回正确错误', async () => {
            await expect(apiClient.getRoom('invalid-id'))
                .rejects.toThrow();
            console.log('✅ 无效请求错误处理正常');
        });

        test('TC_SMOKE_012: 未授权访问被拒绝', async () => {
            await expect(apiClient.deleteRoom(testRoom.id, 'invalid-token'))
                .rejects.toThrow();
            console.log('✅ 访问控制正常');
        });
    });

    // 清理测试数据
    afterAll(async () => {
        try {
            // 离开房间
            if (testRoom && testSession) {
                await apiClient.leaveRoom(testRoom.id, testSession.token);
            }

            // 删除测试房间
            if (testRoom) {
                await apiClient.deleteRoom(testRoom.id, testSession.token);
            }

            console.log('✅ 清理测试数据完成');
        } catch (error) {
            console.warn('⚠️ 清理测试数据时出错:', error.message);
        }
    });
});

// 全局测试设置
beforeAll(() => {
    console.log('🧪 开始执行小沃API冒烟测试');
    console.log('📍 测试环境:', config.env);
    console.log('🔗 API地址:', config.getApiUrl());
});

afterAll(() => {
    console.log('🎉 小沃API冒烟测试完成');
});
