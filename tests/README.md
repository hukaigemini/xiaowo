# 小沃API接口自动化测试套件

**负责人**: 稳当 (SRE)  
**版本**: v1.0.0  
**更新时间**: 2025-12-30

## 📋 概述

小沃API接口自动化测试套件提供全面的API和WebSocket功能测试，确保系统稳定性和可靠性。

## 🏗️ 架构设计

### 测试分层
```
┌─────────────────────────────────────┐
│           测试报告层                 │
│  (HTML报告、覆盖率报告、CI集成)       │
├─────────────────────────────────────┤
│           测试执行层                 │
│  (测试脚本、执行器、CI/CD流水线)      │
├─────────────────────────────────────┤
│           测试逻辑层                 │
│  (业务测试、API测试、WebSocket测试)    │
├─────────────────────────────────────┤
│           工具支撑层                 │
│  (API客户端、WebSocket客户端、配置)   │
└─────────────────────────────────────┘
```

### 测试覆盖范围
- ✅ **房间管理**: 创建、更新、删除、查询
- ✅ **成员操作**: 加入、退出、会话管理
- ✅ **WebSocket**: 连接、消息同步、播放控制
- ✅ **错误处理**: 异常情况、边界条件
- ✅ **性能测试**: 响应时间、并发处理

## 🚀 快速开始

### 环境要求
- Node.js >= 16.0.0
- npm >= 8.0.0
- 后端服务运行在 http://localhost:8080

### 安装依赖
```bash
cd tests
npm install
```

### 运行测试

#### 方式一：使用测试执行脚本 (推荐)
```bash
# 快速冒烟测试
./run_tests.sh --mode quick

# 完整测试
./run_tests.sh --mode full --report --coverage

# CI模式
./run_tests.sh --mode ci

# 仅API测试
./run_tests.sh --mode api

# 仅WebSocket测试
./run_tests.sh --mode websocket
```

#### 方式二：直接使用npm脚本
```bash
# 运行所有测试
npm test

# 运行特定类型测试
npm run test:api
npm run test:websocket
npm run test:smoke

# 生成覆盖率报告
npm run test:coverage

# 生成HTML报告
npm run test:report
```

#### 方式三：使用Jest直接执行
```bash
# 运行所有测试
jest

# 运行特定文件
jest api/roomManagement.test.js

# 监视模式
jest --watch

# 并行执行
jest --maxWorkers=4
```

## 📁 项目结构

```
tests/
├── README.md                 # 测试套件文档
├── run_tests.sh             # 测试执行脚本
├── package.json             # 项目配置和依赖
├── setup.js                 # Jest全局配置
├── test_environment.json    # 测试环境配置
├── TEST_CASES.md           # 测试用例设计
├── utils/                  # 测试工具模块
│   ├── config.js           # 配置管理
│   ├── apiClient.js        # API客户端
│   └── websocketClient.js  # WebSocket客户端
├── api/                    # API测试用例
│   ├── roomManagement.test.js    # 房间管理测试
│   └── memberOperation.test.js   # 成员操作测试
├── websocket/              # WebSocket测试用例
│   └── websocket.test.js         # WebSocket功能测试
├── smoke.test.js           # 冒烟测试
├── reports/                # 测试报告目录
└── coverage/               # 覆盖率报告目录
```

## 🧪 测试用例设计

### 测试分类
- **P0 (阻断级)**: 核心功能失败，影响系统可用性
- **P1 (高优先级)**: 主要功能异常，影响用户体验
- **P2 (中优先级)**: 次要功能问题，不影响核心流程
- **P3 (低优先级)**: 优化项，不影响基本使用

### 测试场景覆盖
1. **正常流程测试**: 验证功能按预期工作
2. **异常流程测试**: 验证错误处理机制
3. **边界条件测试**: 验证极限情况处理
4. **性能测试**: 验证响应时间和并发能力
5. **安全测试**: 验证访问控制和权限

### 测试数据管理
- 使用独立测试环境
- 测试数据自动生成和清理
- 支持多环境配置 (dev/staging/prod)
- 测试数据隔离和复用

## 🔧 配置说明

### test_environment.json
```json
{
  "base_url": "http://localhost:8080",
  "api_version": "v1",
  "environments": {
    "development": {
      "base_url": "http://localhost:8080",
      "timeout": 5000,
      "retries": 3
    },
    "staging": {
      "base_url": "http://staging.xiaowo.com",
      "timeout": 10000,
      "retries": 5
    }
  }
}
```

### 环境变量
```bash
# 测试环境
export NODE_ENV=development

# API基础地址
export API_URL=http://localhost:8080

# 测试超时时间
export TEST_TIMEOUT=30000

# CI模式
export CI=true
```

## �� 测试报告

### 报告类型
1. **控制台报告**: 测试执行过程中的实时输出
2. **HTML报告**: 详细的测试结果和统计信息
3. **覆盖率报告**: 代码覆盖率和分支覆盖情况
4. **CI报告**: 持续集成环境中的标准化报告

### 报告位置
- HTML报告: `reports/test-report.html`
- 覆盖率报告: `coverage/lcov-report/index.html`
- 测试日志: `test_execution.log`

## �� CI/CD集成

### GitHub Actions工作流
项目配置了完整的CI/CD流水线，包括：

1. **代码质量检查**: Go格式检查、静态分析
2. **单元测试**: 后端Go测试、前端测试
3. **API自动化测试**: 完整的接口测试套件
4. **安全扫描**: Gosec安全漏洞扫描
5. **构建和部署**: Docker镜像构建和推送
6. **部署后验证**: 冒烟测试和健康检查

### 触发条件
- `push` 到 `main` 或 `develop` 分支
- 创建 `Pull Request`

### 环境配置
需要在GitHub Secrets中配置：
- `DOCKER_REGISTRY`: Docker镜像仓库地址
- `DOCKER_USERNAME`: Docker用户名
- `DOCKER_PASSWORD`: Docker密码
- `SLACK_WEBHOOK`: Slack通知Webhook

## 🚨 故障排查

### 常见问题

#### 1. 测试执行失败
```bash
# 检查依赖安装
npm install

# 检查测试环境
./run_tests.sh --mode quick --verbose

# 检查API服务状态
curl http://localhost:8080/health
```

#### 2. WebSocket连接失败
- 确认后端WebSocket服务正常运行
- 检查防火墙和网络配置
- 验证WebSocket URL格式

#### 3. 覆盖率报告为空
- 确认测试文件路径配置正确
- 检查Jest配置中的collectCoverageFrom
- 确保utils模块被正确导入

#### 4. CI环境测试失败
- 检查环境变量配置
- 确认服务启动顺序和等待时间
- 验证数据库连接配置

### 日志分析
```bash
# 查看详细执行日志
tail -f test_execution.log

# 查看Jest详细输出
./run_tests.sh --mode full --verbose

# 查看特定测试日志
jest --verbose --no-coverage api/roomManagement.test.js
```

## 🔍 最佳实践

### 测试编写规范
1. **描述性测试名称**: 使用清晰的测试用例描述
2. **单一职责**: 每个测试只验证一个功能点
3. **数据独立性**: 测试之间不应有依赖关系
4. **资源清理**: 测试后及时清理创建的资源
5. **断言充分**: 确保测试结果的准确性

### 性能优化
1. **并行执行**: 使用`--maxWorkers`提高执行效率
2. **测试隔离**: 避免测试间的状态共享
3. **数据库优化**: 使用测试专用数据库
4. **网络优化**: 合理设置超时和重试次数

### 维护建议
1. **定期更新测试用例**: 跟随功能迭代更新
2. **监控测试稳定性**: 关注flaky测试
3. **代码覆盖率监控**: 保持足够的覆盖率水平
4. **测试数据管理**: 定期清理过期测试数据

## 📞 技术支持

### 联系方式
- **负责人**: 稳当 (SRE)
- **技术栈**: Jest, Node.js, WebSocket, RESTful API
- **监控告警**: 集成到现有监控系统

### 相关文档
- [API接口文档](../docs/api/API_CONTRACT.md)
- [系统架构文档](../docs/architecture/)
- [部署指南](../docs/MONITORING.md)

---

**注意**: 本测试套件是小沃项目质量保障体系的重要组成部分，请确保在每次代码变更后运行相关测试。
