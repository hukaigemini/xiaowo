# 小窝监控系统文档

## 概述

小窝同步观影平台监控系统基于 Prometheus + Grafana + Loki + AlertManager 技术栈，提供全方位的系统可观测性能力。

## 架构图

```
┌─────────────┐    ┌──────────────┐    ┌─────────────┐
│   应用服务     │────│  Prometheus  │────│  AlertManager│
│ (Backend/    │    │   指标收集    │    │    告警管理   │
│  Frontend)   │    │              │    │             │
└─────────────┘    └──────────────┘    └─────────────┘
                           │                     │
                           │              ┌──────┴──────┐
                           ▼              ▼             ▼
                    ┌──────────────┐  ┌──────────┐ ┌──────────┐
                    │   Grafana    │  │   Email  │ │  Webhook │
                    │   可视化平台   │  │   通知    │ │  通知     │
                    └──────────────┘  └──────────┘ └──────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │     Loki     │
                    │   日志聚合     │
                    └──────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │   Promtail   │
                    │   日志收集     │
                    └──────────────┘
```

## 服务列表

### 核心监控服务

| 服务名称 | 端口 | 描述 | 访问地址 |
|---------|------|------|----------|
| Prometheus | 9090 | 指标收集和存储 | http://localhost:9090 |
| Grafana | 3001 | 数据可视化 | http://localhost:3001 |
| Loki | 3100 | 日志聚合 | http://localhost:3100 |
| AlertManager | 9093 | 告警管理 | http://localhost:9093 |
| Promtail | 9080 | 日志收集 | - |

### 指标收集器

| 服务名称 | 端口 | 描述 | 访问地址 |
|---------|------|------|----------|
| Node Exporter | 9100 | 系统指标 | http://localhost:9100 |
| Redis Exporter | 9121 | Redis指标 | http://localhost:9121 |
| cAdvisor | 8081 | 容器指标 | http://localhost:8081 |
| Blackbox Exporter | 9115 | 黑盒监控 | http://localhost:9115 |

## 快速开始

### 1. 启动监控系统

```bash
# 启动所有服务（包含监控系统）
docker-compose --profile monitoring up -d

# 仅启动基础应用服务（不包含监控系统）
docker-compose up -d

# 启动监控系统（如果基础服务已运行）
docker-compose --profile monitoring up -d prometheus grafana loki alertmanager
```

### 2. 访问监控界面

- **Grafana 仪表板**: http://localhost:3001
  - 默认用户名: `admin`
  - 默认密码: `admin123`

- **Prometheus 查询**: http://localhost:9090

- **AlertManager**: http://localhost:9093

- **Loki 日志**: 通过 Grafana Explore 页面访问

### 3. 验证监控状态

```bash
# 检查所有监控服务状态
docker-compose --profile monitoring ps

# 检查特定服务健康状态
curl http://localhost:9090/-/healthy
curl http://localhost:3001/api/health
curl http://localhost:9093/-/healthy
curl http://localhost:3100/ready
```

## 监控指标

### 应用级指标

#### HTTP 请求指标
- `http_requests_total`: HTTP 请求总数
- `http_request_duration_seconds`: HTTP 请求响应时间直方图
- `http_request_size_bytes`: HTTP 请求大小
- `http_response_size_bytes`: HTTP 响应大小

#### WebSocket 指标
- `xiaowo_websocket_connections_active`: 活跃 WebSocket 连接数
- `xiaowo_websocket_messages_total`: WebSocket 消息总数
- `xiaowo_websocket_message_duration_seconds`: 消息处理时间

#### 业务指标
- `xiaowo_active_rooms`: 活跃房间数量
- `xiaowo_online_users`: 在线用户数量
- `xiaowo_room_sync_total`: 房间同步总数
- `xiaowo_room_sync_success_total`: 房间同步成功数

### 系统级指标

#### CPU 和内存
- `node_cpu_seconds_total`: CPU 使用时间
- `node_memory_MemTotal_bytes`: 总内存
- `node_memory_MemAvailable_bytes`: 可用内存
- `node_memory_MemUsed_bytes`: 已使用内存

#### 磁盘
- `node_filesystem_avail_bytes`: 可用磁盘空间
- `node_filesystem_size_bytes`: 总磁盘空间

#### 网络
- `node_network_receive_bytes_total`: 网络接收字节数
- `node_network_transmit_bytes_total`: 网络发送字节数

### Redis 指标
- `redis_connected_clients`: Redis 客户端连接数
- `redis_used_memory_bytes`: Redis 内存使用量
- `redis_keyspace_hits_total`: Redis 键空间命中数

## 告警规则

### 基础设施告警

1. **服务不可用告警**
   - 触发条件: `up{job="xiaowo-backend"} == 0`
   - 严重程度: Critical
   - 描述: 后端服务不可用

2. **HTTP 错误率告警**
   - 触发条件: HTTP 5xx 错误率 > 5%
   - 严重程度: Warning
   - 描述: 后端 HTTP 错误率过高

### 性能告警

3. **响应时间告警**
   - 触发条件: P95 响应时间 > 2s
   - 严重程度: Warning
   - 描述: 后端响应时间过长

4. **资源使用率告警**
   - 触发条件: CPU 使用率 > 80% 或 内存使用率 > 80%
   - 严重程度: Warning
   - 描述: 系统资源使用率过高

### 业务告警

5. **活跃房间数异常**
   - 触发条件: 活跃房间数 = 0 且持续 5 分钟
   - 严重程度: Warning
   - 描述: 系统可能存在问题

6. **WebSocket 连接异常**
   - 触发条件: WebSocket 连接数 < 预期最小值
   - 严重程度: Warning
   - 描述: WebSocket 服务可能存在问题

## 仪表板使用

### 系统概览仪表板
- **位置**: Grafana 首页 → "小窝系统概览"
- **用途**: 整体系统运行状态概览
- **刷新频率**: 30秒

#### 关键面板
- 服务状态: 显示所有服务是否在线
- HTTP 请求率: 显示请求量和趋势
- 系统负载: CPU 和内存使用率
- 活跃房间数和在线用户数: 核心业务指标
- 错误率趋势: HTTP 5xx 错误率变化

### 应用性能仪表板
- **位置**: Grafana 首页 → "小窝应用性能监控"
- **用途**: 应用性能详细分析
- **刷新频率**: 15秒

#### 关键面板
- HTTP 响应时间 P95: 95%请求的响应时间
- HTTP 请求率: 按状态码和处理器分类
- WebSocket 连接状态: 实时连接数
- 消息处理延迟: WebSocket 消息处理时间
- 数据库查询时间: 数据库性能指标
- 房间同步成功率: 业务操作成功率

### 日志查询
1. 打开 Grafana → Explore
2. 选择数据源: Loki
3. 使用 LogQL 查询语言:

```logql
# 查询错误日志
{job="xiaowo-backend"} |= "ERROR"

# 查询特定时间范围的日志
{job="xiaowo-backend"} |= "ERROR" | json | level="error" | timestamp > "2023-12-01T00:00:00Z"

# 按标签过滤
{job="xiaowo-backend",level="error"}

# 统计日志数量
count_over_time({job="xiaowo-backend"}[5m])
```

## 环境变量

监控服务相关环境变量（在 `.env` 文件中配置）:

```bash
# 端口配置
PROMETHEUS_PORT=9090
GRAFANA_PORT=3001
LOKI_PORT=3100
ALERTMANAGER_PORT=9093
NODE_EXPORTER_PORT=9100
REDIS_EXPORTER_PORT=9121
CADVISOR_PORT=8081
BLACKBOX_EXPORTER_PORT=9115

# 安全配置
GRAFANA_ADMIN_PASSWORD=admin123
```

## 数据持久化

监控数据使用 Docker 卷进行持久化:

- `prometheus_data`: Prometheus 指标数据
- `grafana_data`: Grafana 仪表板和配置
- `loki_data`: Loki 日志数据
- `alertmanager_data`: AlertManager 配置和数据

## 性能优化

### Prometheus 配置优化
- 指标收集间隔: 15秒（默认）
- 数据保留时间: 15天（默认）
- 存储压缩: 启用

### Grafana 配置优化
- 数据库: SQLite（开发环境）
- 缓存: 启用查询缓存
- 插件: 安装必要插件

### Loki 配置优化
- 存储后端: 文件系统
- 索引: Boltdb-shipper
- 压缩: Snappy 压缩
- 查询缓存: 100MB

## 扩展功能

### 添加新的监控目标
1. 在 `monitoring/prometheus.yml` 中添加 scrape_config
2. 确保目标服务暴露 `/metrics` 端点
3. 重启 Prometheus 服务

### 自定义告警规则
1. 编辑 `monitoring/alert-rules.yml`
2. 添加新的告警规则
3. 在 AlertManager 中配置通知渠道

### 创建自定义仪表板
1. 在 Grafana 中创建新仪表板
2. 配置面板和数据源
3. 导出 JSON 配置
4. 保存到 `monitoring/grafana/dashboards/` 目录

## 最佳实践

1. **监控覆盖度**: 确保关键业务指标都有相应的监控
2. **告警准确性**: 避免误报，确保告警的可操作性
3. **数据保留**: 根据需要调整数据保留时间，平衡性能和成本
4. **定期检查**: 定期检查监控系统的健康状态
5. **文档更新**: 及时更新监控配置和文档

## 故障排查

详细的故障排查指南请参考 [TROUBLESHOOTING.md](./TROUBLESHOOTING.md)
