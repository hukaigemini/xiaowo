# 小窝同步观影平台 - 快速启动指南

## 🚀 一键启动开发环境

### 环境要求
- Docker 20.10+
- Docker Compose 2.0+

### 快速启动步骤

```bash
# 1. 准备环境配置
cp .env.example .env

# 2. 启动所有服务
docker-compose up -d

# 3. 查看服务状态
docker-compose ps

# 4. 访问应用
# 前端: http://localhost:3000
# 后端API: http://localhost:8080
```

## 📊 服务架构

| 服务 | 端口 | 描述 | 状态检查 |
|------|------|------|----------|
| Frontend | 3000 | Vue3 前端应用 | http://localhost:3000 |
| Backend | 8080 | Go API 服务 | http://localhost:8080/health |
| Redis | 6379 | 会话存储 | docker-compose exec redis redis-cli ping |

## 🔧 常用命令

```bash
# 查看日志
docker-compose logs -f

# 重启服务
docker-compose restart

# 停止服务
docker-compose down

# 清理数据（慎用）
docker-compose down -v
```

## 📱 使用说明

1. **创建房间**: 点击首页"创建房间"按钮
2. **加入房间**: 使用邀请链接或房间ID
3. **播放视频**: 输入视频URL，支持.mp4和.m3u8格式
4. **同步观看**: 所有成员自动同步播放状态

## 🛠️ 开发模式

如需修改代码进行开发：

```bash
# 后端开发
docker-compose exec backend sh
go run cmd/server/main.go

# 前端开发  
docker-compose exec frontend sh
npm run dev
```

---
**最后更新**: 2025-12-30  
**维护者**: 稳当 (SRE)
