# 监控系统故障排查指南

## 概述

本指南提供了小窝监控系统常见问题的诊断和解决方法，帮助快速定位和解决监控相关问题。

## 快速诊断命令

### 服务健康检查

```bash
# 检查所有监控服务状态
docker-compose --profile monitoring ps

# 检查服务日志
docker-compose --profile monitoring logs prometheus
docker-compose --profile monitoring logs grafana
docker-compose --profile monitoring logs loki

# 检查端口监听状态
netstat -tlnp | grep -E "(9090|3001|3100|9093)"

# 检查磁盘空间
df -h
```

### 网络连接检查

```bash
# 检查服务间网络连通性
docker-compose --profile monitoring exec prometheus wget -qO- http://loki:3100/ready
docker-compose --profile monitoring exec grafana wget -qO- http://prometheus:9090/-/healthy

# 检查DNS解析
docker-compose --profile monitoring exec prometheus nslookup loki
```

## 常见问题及解决方案

### 1. Prometheus 无法收集指标

#### 问题症状
- Prometheus 目标页面显示所有服务为 down
- 仪表板无数据显示
- 告警规则触发 "TargetDown" 告警

#### 诊断步骤

1. **检查 Prometheus 配置**
   ```bash
   # 验证配置文件语法
   docker-compose --profile monitoring exec prometheus promtool check config /etc/prometheus/prometheus.yml
   docker-compose --profile monitoring exec prometheus promtool check rules /etc/prometheus/alert-rules.yml
   ```

2. **检查网络连接**
   ```bash
   # 从 Prometheus 容器内部测试网络连接
   docker-compose --profile monitoring exec prometheus wget -qO- http://backend:8080/metrics
   docker-compose --profile monitoring exec prometheus wget -qO- http://redis-exporter:9121/metrics
   ```

3. **检查目标服务状态**
   ```bash
   # 检查后端服务
   curl http://localhost:8080/health
   curl http://localhost:8080/metrics
   
   # 检查 Redis Exporter
   curl http://localhost:9121/metrics
   ```

#### 解决方案

1. **修复网络配置**
   ```bash
   # 重启监控服务
   docker-compose --profile monitoring restart prometheus
   
   # 检查 Docker 网络
   docker network ls
   docker network inspect xiaowo-app-network
   ```

2. **修复配置错误**
   ```yaml
   # 在 prometheus.yml 中确保正确的服务名和端口
   scrape_configs:
     - job_name: 'xiaowo-backend'
       static_configs:
         - targets: ['backend:8080']  # 使用容器名，不是 localhost
   ```

### 2. Grafana 无法连接数据源

#### 问题症状
- Grafana 仪表板显示 "No data"
- 数据源连接失败
- 查询报错

#### 诊断步骤

1. **检查 Grafana 日志**
   ```bash
   docker-compose --profile monitoring logs grafana
   ```

2. **验证数据源配置**
   ```bash
   # 从 Grafana 容器内部测试 Prometheus 连接
   docker-compose --profile monitoring exec grafana wget -qO- http://prometheus:9090/api/v1/query?query=up
   ```

3. **检查防火墙和安全组**
   ```bash
   # 检查端口访问
   curl -u admin:admin123 http://localhost:3001/api/health
   ```

#### 解决方案

1. **重置 Grafana 数据源**
   ```bash
   # 删除并重新创建数据源
   curl -X DELETE -u admin:admin123 http://localhost:3001/api/datasources/name/Prometheus
   curl -X POST -u admin:admin123 -H "Content-Type: application/json" \
     -d '{"name":"Prometheus","type":"prometheus","url":"http://prometheus:9090","isDefault":true}' \
     http://localhost:3001/api/datasources
   ```

2. **修复网络问题**
   ```bash
   # 确保 Grafana 和 Prometheus 在同一网络
   docker-compose --profile monitoring restart grafana
   ```

### 3. Loki 无法收集日志

#### 问题症状
- Grafana Explore 页面无日志数据
- Promtail 日志显示错误
- 日志查询返回空结果

#### 诊断步骤

1. **检查 Promtail 配置**
   ```bash
   # 验证 Promtail 配置文件
   docker-compose --profile monitoring exec promtail cat /etc/promtail/config.yml
   ```

2. **检查日志文件路径**
   ```bash
   # 检查容器日志路径
   ls -la /var/lib/docker/containers/
   ls -la /var/log/xiaowo/
   ```

3. **测试日志收集**
   ```bash
   # 从 Promtail 容器测试 Loki 连接
   docker-compose --profile monitoring exec promtail wget -qO- http://loki:3100/ready
   ```

#### 解决方案

1. **修复 Promtail 配置**
   ```yaml
   # 在 promtail-config.yml 中确保正确的路径
   scrape_configs:
     - job_name: xiaowo-backend
       static_configs:
         - targets:
             - localhost
           labels:
             job: xiaowo-backend
             __path__: /var/log/xiaowo/backend/*.log  # 确保路径正确
   ```

2. **创建日志目录**
   ```bash
   # 创建应用日志目录
   sudo mkdir -p /var/log/xiaowo/{backend,frontend}
   sudo chown -R 10001:10001 /var/log/xiaowo/
   ```

3. **重启日志收集服务**
   ```bash
   docker-compose --profile monitoring restart promtail loki
   ```

### 4. AlertManager 告警不发送

#### 问题症状
- 告警触发但未收到通知
- AlertManager 显示告警状态正常
- 邮件发送失败

#### 诊断步骤

1. **检查 AlertManager 配置**
   ```bash
   # 验证配置文件
   docker-compose --profile monitoring exec alertmanager amtool config check
   ```

2. **测试 SMTP 连接**
   ```bash
   # 测试邮件发送
   docker-compose --profile monitoring exec alertmanager amtool config routes test
   ```

3. **检查告警规则**
   ```bash
   # 在 Prometheus 中测试告警规则
   curl http://localhost:9090/api/v1/query?query='up{job="xiaowo-backend"}==0'
   ```

#### 解决方案

1. **配置邮件服务**
   ```yaml
   # 在 alertmanager.yml 中配置正确的 SMTP
   global:
     smtp_smarthost: 'smtp.gmail.com:587'
     smtp_from: 'your-email@gmail.com'
     smtp_auth_username: 'your-email@gmail.com'
     smtp_auth_password: 'your-app-password'
   ```

2. **修复路由配置**
   ```yaml
   # 确保告警路由配置正确
   route:
     receiver: 'default-receiver'
     routes:
       - match:
           severity: critical
         receiver: 'critical-alerts'
   ```

### 5. 磁盘空间不足

#### 问题症状
- Prometheus 无法写入数据
- Loki 日志收集停止
- 容器重启频繁

#### 诊断步骤

1. **检查磁盘使用情况**
   ```bash
   # 检查磁盘使用率
   df -h
   
   # 检查 Docker 卷大小
   docker system df
   
   # 检查具体目录大小
   du -sh /var/lib/docker/volumes/*
   ```

2. **检查监控数据大小**
   ```bash
   # 检查 Prometheus 数据目录
   docker-compose --profile monitoring exec prometheus du -sh /prometheus
   
   # 检查 Loki 数据目录
   docker-compose --profile monitoring exec loki du -sh /loki
   ```

#### 解决方案

1. **清理旧数据**
   ```bash
   # 清理未使用的 Docker 资源
   docker system prune -f
   
   # 清理 Prometheus 旧数据（保留最近 7 天）
   docker-compose --profile monitoring exec prometheus find /prometheus -name "*.db" -mtime +7 -delete
   ```

2. **调整数据保留策略**
   ```yaml
   # 在 prometheus.yml 中添加数据保留配置
   global:
     retention: 7d  # 保留 7 天数据
   ```

3. **扩容磁盘**
   ```bash
   # 扩展磁盘空间或移动数据到更大磁盘
   ```

### 6. 服务启动失败

#### 问题症状
- 容器启动后立即退出
- 健康检查失败
- 端口冲突错误

#### 诊断步骤

1. **查看容器日志**
   ```bash
   # 查看启动失败的详细日志
   docker-compose --profile monitoring logs --tail=50 prometheus
   docker-compose --profile monitoring logs --tail=50 grafana
   ```

2. **检查端口占用**
   ```bash
   # 检查端口是否被占用
   netstat -tlnp | grep -E "(9090|3001|3100|9093)"
   
   # 检查 Docker 端口映射
   docker-compose --profile monitoring ps
   ```

3. **检查配置文件权限**
   ```bash
   # 检查配置文件权限
   ls -la monitoring/
   docker-compose --profile monitoring exec prometheus ls -la /etc/prometheus/
   ```

#### 解决方案

1. **修复端口冲突**
   ```bash
   # 停止占用端口的进程
   sudo lsof -ti:9090 | xargs kill -9
   
   # 或者修改端口映射
   # 在 docker-compose.yml 中修改端口配置
   ```

2. **修复权限问题**
   ```bash
   # 修复配置文件权限
   chmod 644 monitoring/*.yml
   chmod -R 755 monitoring/grafana/
   
   # 确保 Grafana 用户有权限访问数据目录
   docker-compose --profile monitoring exec grafana chown -R 472:472 /var/lib/grafana
   ```

3. **清理损坏的数据**
   ```bash
   # 删除损坏的卷数据
   docker-compose --profile monitoring down
   docker volume rm xiaowo_prometheus_data xiaowo_grafana_data
   docker-compose --profile monitoring up -d
   ```

## 性能优化问题

### 1. 查询响应慢

#### 症状
- Grafana 仪表板加载缓慢
- Prometheus 查询超时
- Loki 日志查询慢

#### 解决方案

1. **优化 Prometheus 查询**
   ```promql
   # 使用适当的查询时间范围
   rate(http_requests_total[5m])  # 使用聚合函数
   
   # 避免长时间范围的全量扫描
   # 改用较小的查询窗口
   ```

2. **调整 Prometheus 配置**
   ```yaml
   # 在 prometheus.yml 中优化配置
   global:
     scrape_interval: 30s  # 减少收集频率
     evaluation_interval: 30s
   ```

3. **优化 Loki 查询**
   ```logql
   # 添加时间范围限制
   {job="xiaowo-backend"} |= "ERROR" | timestamp > "2023-12-01T00:00:00Z" | timestamp < "2023-12-01T23:59:59Z"
   
   # 使用标签过滤减少扫描范围
   {job="xiaowo-backend",level="error"}
   ```

### 2. 内存使用过高

#### 症状
- 容器内存使用率持续增长
- 系统响应变慢
- OOM Kill 发生

#### 解决方案

1. **限制容器内存使用**
   ```yaml
   # 在 docker-compose.yml 中添加内存限制
   services:
     prometheus:
       deploy:
         resources:
           limits:
             memory: 2G
           reservations:
             memory: 1G
   ```

2. **优化 Grafana 配置**
   ```ini
   # 在 grafana.ini 中设置
   [database]
   max_idle_conn = 2
   max_open_conn = 0
   conn_max_lifetime = 14400
   
   [query_history]
   enabled = false
   
   [dashboard_previews]
   enabled = false
   ```

## 监控数据异常

### 1. 指标值异常

#### 症状
- 指标显示负值或极大值
- 计数器重置
- 数据丢失

#### 诊断步骤

1. **检查应用指标暴露**
   ```bash
   # 检查后端服务指标
   curl http://localhost:8080/metrics | grep -E "(http_requests|websocket)"
   
   # 检查 Redis 指标
   curl http://localhost:9121/metrics | grep connected_clients
   ```

2. **验证指标计算**
   ```bash
   # 使用 PromQL 测试查询
   curl 'http://localhost:9090/api/v1/query?query=rate(http_requests_total[5m])'
   ```

#### 解决方案

1. **修复应用代码**
   ```go
   // 确保指标正确递增
   httpRequestsTotal.WithLabelValues(method, route, statusCode).Inc()
   
   // 避免指标重置
   if counter < previous {
       // 处理计数器重置情况
   }
   ```

2. **调整告警阈值**
   ```yaml
   # 在 alert-rules.yml 中调整阈值
   - alert: HighHTTPErrorRate
     expr: |
       (
         rate(http_requests_total{job="xiaowo-backend",status=~"5.."}[5m]) 
         / 
         rate(http_requests_total{job="xiaowo-backend"}[5m])
       ) * 100 > 10  # 提高阈值
   ```

## 维护和更新

### 定期维护任务

1. **每日检查**
   ```bash
   # 检查服务健康状态
   ./scripts/health-check.sh
   
   # 检查磁盘空间
   df -h | awk '$5 > 80 {print "WARNING: " $1 " is " $5 " full"}'
   
   # 检查告警状态
   curl -s http://localhost:9093/api/v1/alerts | jq '.data[].state'
   ```

2. **每周维护**
   ```bash
   # 清理旧日志
   docker system prune -f
   
   # 更新镜像
   docker-compose --profile monitoring pull
   docker-compose --profile monitoring up -d
   
   # 备份配置
   tar -czf monitoring-backup-$(date +%Y%m%d).tar.gz monitoring/
   ```

3. **每月检查**
   ```bash
   # 检查监控覆盖率
   # 验证所有关键指标都有相应告警
   
   # 性能评估
   # 检查查询响应时间和系统资源使用率
   
   # 文档更新
   # 更新监控文档和故障排查指南
   ```

### 版本升级

1. **升级前准备**
   ```bash
   # 备份现有配置
   cp -r monitoring/ monitoring-backup-$(date +%Y%m%d)/
   
   # 检查新版本兼容性
   docker-compose --profile monitoring config
   ```

2. **执行升级**
   ```bash
   # 拉取新版本镜像
   docker-compose --profile monitoring pull
   
   # 滚动更新服务
   docker-compose --profile monitoring up -d
   ```

3. **验证升级**
   ```bash
   # 验证服务正常
   ./scripts/health-check.sh
   
   # 检查数据完整性
   # 验证仪表板显示正常
   ```

## 紧急处理流程

### 1. 服务完全不可用

1. **立即响应**
   ```bash
   # 检查所有服务状态
   docker-compose --profile monitoring ps
   
   # 查看关键日志
   docker-compose --profile monitoring logs --tail=100
   ```

2. **快速恢复**
   ```bash
   # 重启所有监控服务
   docker-compose --profile monitoring restart
   
   # 如果问题持续，回滚到上一版本
   docker-compose --profile monitoring down
   git checkout HEAD~1
   docker-compose --profile monitoring up -d
   ```

### 2. 数据丢失

1. **数据恢复**
   ```bash
   # 检查数据卷
   docker volume ls | grep xiaowo
   
   # 恢复备份数据
   docker run --rm -v xiaowo_prometheus_data:/data -v $(pwd):/backup alpine \
     sh -c "cp /backup/prometheus-backup/* /data/"
   ```

2. **预防措施**
   - 定期备份监控配置和数据
   - 设置监控数据保留策略
   - 实施告警规则监控监控服务本身

## 联系支持

如果按照本指南无法解决问题，请：

1. 收集以下信息：
   - 详细的错误日志
   - 系统配置信息
   - 复现步骤
   - 环境信息

2. 提交问题报告，包含：
   - 问题描述
   - 故障排查步骤
   - 已尝试的解决方案
   - 相关日志和截图

## 相关资源

- [Prometheus 官方文档](https://prometheus.io/docs/)
- [Grafana 官方文档](https://grafana.com/docs/)
- [Loki 官方文档](https://grafana.com/docs/loki/)
- [AlertManager 官方文档](https://prometheus.io/docs/alerting/latest/alertmanager/)

---

**注意**: 本指南持续更新，建议定期查看最新版本。
