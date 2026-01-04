/**
 * Jest 测试设置文件
 * 作者: 稳当 (SRE)
 * 功能: 全局测试配置和清理
 */

const path = require('path');

// 全局测试超时设置
jest.setTimeout(30000);

// 全局控制台日志捕获
const originalConsoleLog = console.log;
const originalConsoleWarn = console.warn;
const originalConsoleError = console.error;

beforeAll(() => {
  // 测试开始前的全局设置
  console.log('🧪 开始执行小沃API测试套件');
});

afterAll(() => {
  // 测试结束后的全局清理
  console.log('🎉 小沃API测试套件执行完成');
});

// 捕获并格式化控制台输出
beforeEach(() => {
  // 每个测试前的设置
});

afterEach(() => {
  // 每个测试后的清理
});

// 全局错误处理
process.on('unhandledRejection', (reason, promise) => {
  console.error('未处理的Promise拒绝:', promise, 'reason:', reason);
});

process.on('uncaughtException', (error) => {
  console.error('未捕获的异常:', error);
});

// 测试环境验证
global.testConfig = {
  isCI: process.env.CI === 'true' || process.env.NODE_ENV === 'ci',
  nodeEnv: process.env.NODE_ENV || 'development',
  apiUrl: process.env.API_URL || 'http://localhost:8080',
  timeout: parseInt(process.env.TEST_TIMEOUT) || 30000
};

// 全局测试工具函数
global.testUtils = {
  // 生成随机测试数据
  generateTestData: (prefix = 'test') => ({
    id: `${prefix}_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
    timestamp: Date.now(),
    random: Math.random()
  }),

  // 等待函数
  sleep: (ms) => new Promise(resolve => setTimeout(resolve, ms)),

  // 重试函数
  retry: async (fn, maxAttempts = 3, delay = 1000) => {
    for (let i = 0; i < maxAttempts; i++) {
      try {
        return await fn();
      } catch (error) {
        if (i === maxAttempts - 1) throw error;
        console.log(`重试 ${i + 1}/${maxAttempts}: ${error.message}`);
        await global.testUtils.sleep(delay);
      }
    }
  },

  // 验证响应格式
  validateApiResponse: (response) => {
    expect(response).toHaveProperty('status');
    expect(response.status).toMatch(/^(success|error)$/);
    
    if (response.status === 'success') {
      expect(response).toHaveProperty('data');
    } else {
      expect(response).toHaveProperty('message');
    }
  }
};

// 自定义匹配器
expect.extend({
  toBeValidApiResponse(received) {
    const isValid = received && 
                   typeof received.status === 'string' && 
                   ['success', 'error'].includes(received.status);
    
    return {
      message: () => `期望 ${received} 是有效的API响应格式`,
      pass: isValid
    };
  },

  toBeValidRoomData(received) {
    const isValid = received && 
                   typeof received.id === 'string' && 
                   typeof received.name === 'string' &&
                   typeof received.is_private === 'boolean';
    
    return {
      message: () => `期望 ${received} 是有效的房间数据`,
      pass: isValid
    };
  },

  toBeValidSessionData(received) {
    const isValid = received && 
                   typeof received.token === 'string' && 
                   received.token.length > 0;
    
    return {
      message: () => `期望 ${received} 是有效的会话数据`,
      pass: isValid
    };
  }
});

console.log('✅ Jest 测试环境设置完成');
