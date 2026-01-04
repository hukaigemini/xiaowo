const fs = require('fs');
const path = require('path');

// 加载测试环境配置
const testConfig = JSON.parse(
  fs.readFileSync(path.join(__dirname, '../test_environment.json'), 'utf8')
);

class TestConfig {
  constructor() {
    this.env = process.env.NODE_ENV || 'development';
    this.config = testConfig.environments[this.env] || testConfig.environments.development;
    this.testData = testConfig.test_data;
  }

  getBaseUrl() {
    return this.config.base_url;
  }

  getApiUrl(endpoint = '') {
    return `${this.getBaseUrl()}/api/${this.config.api_version}${endpoint}`;
  }

  getWebSocketUrl(roomId, token) {
    const wsProtocol = this.getBaseUrl().startsWith('https') ? 'wss' : 'ws';
    const baseWsUrl = this.getBaseUrl().replace(/^https?/, wsProtocol);
    return `${baseWsUrl}/ws/room/${roomId}?token=${token}`;
  }

  getTimeout() {
    return this.config.timeout;
  }

  getRetries() {
    return this.config.retries;
  }

  getTestUser(index = 0) {
    return this.testData.users[index] || this.testData.users[0];
  }

  getTestRoom(index = 0) {
    return this.testData.rooms[index] || this.testData.rooms[0];
  }

  getPerformanceThresholds() {
    return this.testData.performance_thresholds;
  }

  getWebSocketSettings() {
    return this.testData.websocket_settings;
  }
}

module.exports = new TestConfig();
