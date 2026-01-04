# 小窝同步观影平台 - API接口契约文档

**版本**: v1.0  
**创建时间**: 2025-12-30  
**描述**: 基于数据库Schema设计的完整API接口契约，为前端开发和API实现提供统一的规范

---

## 📋 目录

- [1. API基础信息](#1-api基础信息)
- [2. 通用响应格式](#2-通用响应格式)
- [3. 数据模型定义](#3-数据模型定义)
- [4. RESTful API接口](#4-restful-api接口)
- [5. WebSocket事件契约](#5-websocket事件契约)
- [6. 错误码定义](#6-错误码定义)
- [7. 请求示例](#7-请求示例)

---

## 1. API基础信息

### 1.1 API版本控制
- **Base URL**: `http://localhost:8080/api/v1`
- **版本策略**: URL路径版本控制
- **Content-Type**: `application/json; charset=utf-8`

### 1.2 认证方式
- **用户认证**: 基于Session Token（UUID格式）
- **房间访问**: 通过房间ID和密码（如需要）
- **权限控制**: 基于角色（host/guest）

### 1.3 通用HTTP状态码
- `200 OK` - 请求成功
- `201 Created` - 创建成功
- `400 Bad Request` - 请求参数错误
- `401 Unauthorized` - 未认证
- `403 Forbidden` - 权限不足
- `404 Not Found` - 资源不存在
- `409 Conflict` - 资源冲突
- `429 Too Many Requests` - 请求频率限制
- `500 Internal Server Error` - 服务器内部错误

---

## 2. 通用响应格式

### 2.1 成功响应
```json
{
    "code": 0,
    "message": "success",
    "data": {
        // 实际数据
    },
    "timestamp": "2025-12-30T10:30:00Z"
}
```

### 2.2 错误响应
```json
{
    "code": 40001,
    "message": "Invalid parameter",
    "error": "详细错误信息",
    "timestamp": "2025-12-30T10:30:00Z"
}
```

### 2.3 分页响应
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "items": [],
        "pagination": {
            "page": 1,
            "size": 20,
            "total": 100,
            "pages": 5
        }
    }
}
```

---

## 3. 数据模型定义

### 3.1 UserSession（用户会话）
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "快乐的小熊猫",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=panda",
    "room_id": "ABC123",
    "created_at": "2025-12-30T10:30:00Z",
    "last_seen_at": "2025-12-30T10:35:00Z",
    "expires_at": "2026-01-06T10:30:00Z"
}
```

### 3.2 Room（房间）
```json
{
    "id": "ABC123",
    "name": "周末电影时光",
    "description": "一起观看经典电影",
    "creator_session_id": "550e8400-e29b-41d4-a716-446655440000",
    "is_private": false,
    "password": null,
    "max_users": 7,
    "status": "active",
    "media_url": "https://example.com/movie.mp4",
    "media_type": "video",
    "media_title": "肖申克的救赎",
    "media_duration": 8520,
    "playback_state": "paused",
    "current_time": 0,
    "playback_rate": 1.0,
    "settings": {
        "auto_sync": true,
        "allow_control": true,
        "chat_enabled": true
    },
    "version": 1,
    "last_active_at": "2025-12-30T10:30:00Z",
    "created_at": "2025-12-30T10:30:00Z",
    "updated_at": "2025-12-30T10:35:00Z"
}
```

### 3.3 RoomMember（房间成员）
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "room_id": "ABC123",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "nickname": "快乐的小熊猫",
    "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=panda",
    "role": "host",
    "permissions": {
        "control": true,
        "chat": true,
        "invite": false
    },
    "joined_at": "2025-12-30T10:30:00Z",
    "last_seen_at": "2025-12-30T10:35:00Z",
    "left_at": null,
    "is_active": true
}
```

### 3.4 RoomMessage（房间消息）
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "room_id": "ABC123",
    "session_id": "550e8400-e29b-41d4-a716-446655440000",
    "message_type": "chat",
    "content": "这部电影真的太棒了！",
    "metadata": {
        "emojis": ["👍", "❤️"],
        "reply_to": null
    },
    "created_at": "2025-12-30T10:35:00Z"
}
```

### 3.5 SystemConfig（系统配置）
```json
{
    "id": "app_name",
    "config_key": "app_name",
    "config_value": "小窝同步观影平台",
    "config_type": "string",
    "description": "应用名称",
    "is_editable": true,
    "created_at": "2025-12-30T10:30:00Z",
    "updated_at": "2025-12-30T10:30:00Z"
}
```

---

## 4. RESTful API接口

### 4.1 用户会话管理

#### 4.1.1 创建会话
**POST** `/sessions`

**请求体**:
```json
{
    "nickname": "快乐的小熊猫"
}
```

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "session": {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "nickname": "快乐的小熊猫",
            "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=panda",
            "room_id": null,
            "expires_at": "2026-01-06T10:30:00Z"
        },
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

#### 4.1.2 更新会话
**PUT** `/sessions/{session_id}`

**请求体**:
```json
{
    "nickname": "新的昵称",
    "last_seen_at": "2025-12-30T10:35:00Z"
}
```

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "session": {
            // 更新的会话信息
        }
    }
}
```

#### 4.1.3 获取会话信息
**GET** `/sessions/{session_id}`

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "session": {
            // 会话详细信息
        }
    }
}
```

### 4.2 房间管理

#### 4.2.1 创建房间
**POST** `/rooms`

**请求头**: `Authorization: Bearer {session_token}`

**请求体**:
```json
{
    "name": "周末电影时光",
    "description": "一起观看经典电影",
    "is_private": false,
    "password": null,
    "media_url": "https://example.com/movie.mp4",
    "media_title": "肖申克的救赎",
    "media_type": "video",
    "settings": {
        "auto_sync": true,
        "allow_control": true
    }
}
```

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "room": {
            "id": "ABC123",
            "name": "周末电影时光",
            // 完整房间信息
        }
    }
}
```

#### 4.2.2 获取房间信息
**GET** `/rooms/{room_id}`

**查询参数**:
- `include_members`: boolean, 是否包含成员信息

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "room": {
            // 房间详细信息
        },
        "members": [
            // 房间成员列表（如果include_members=true）
        ]
    }
}
```

#### 4.2.3 更新房间信息
**PUT** `/rooms/{room_id}`

**请求头**: `Authorization: Bearer {session_token}`

**请求体**:
```json
{
    "name": "新的房间名称",
    "description": "新的描述",
    "settings": {
        "auto_sync": false,
        "allow_control": false
    }
}
```

#### 4.2.4 获取房间列表
**GET** `/rooms`

**查询参数**:
- `status`: string, 房间状态筛选
- `is_private`: boolean, 私密房间筛选
- `page`: number, 页码
- `size`: number, 每页数量

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "rooms": [
            // 房间列表
        ],
        "pagination": {
            "page": 1,
            "size": 20,
            "total": 5,
            "pages": 1
        }
    }
}
```

#### 4.2.5 删除房间
**DELETE** `/rooms/{room_id}`

**请求头**: `Authorization: Bearer {session_token}`

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": null
}
```

### 4.3 房间成员管理

#### 4.3.1 加入房间
**POST** `/rooms/{room_id}/join`

**请求体**:
```json
{
    "password": "房间密码（如需要）"
}
```

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "member": {
            // 成员信息
        }
    }
}
```

#### 4.3.2 离开房间
**POST** `/rooms/{room_id}/leave`

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": null
}
```

#### 4.3.3 获取房间成员
**GET** `/rooms/{room_id}/members`

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "members": [
            // 房间成员列表
        ]
    }
}
```

#### 4.3.4 更新成员权限
**PUT** `/rooms/{room_id}/members/{session_id}/permissions`

**请求体**:
```json
{
    "role": "guest",
    "permissions": {
        "control": true,
        "chat": true,
        "invite": false
    }
}
```

### 4.4 消息管理

#### 4.4.1 发送消息
**POST** `/rooms/{room_id}/messages`

**请求体**:
```json
{
    "message_type": "chat",
    "content": "这部电影真的太棒了！",
    "metadata": {
        "reply_to": null
    }
}
```

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "message": {
            "id": "550e8400-e29b-41d4-a716-446655440002",
            "room_id": "ABC123",
            "session_id": "550e8400-e29b-41d4-a716-446655440000",
            "message_type": "chat",
            "content": "这部电影真的太棒了！",
            "created_at": "2025-12-30T10:35:00Z"
        }
    }
}
```

#### 4.4.2 获取消息历史
**GET** `/rooms/{room_id}/messages`

**查询参数**:
- `message_type`: string, 消息类型筛选
- `since`: string, 获取指定时间之后的消息
- `limit`: number, 消息数量限制

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "messages": [
            // 消息列表
        ]
    }
}
```

### 4.5 播放状态控制

#### 4.5.1 获取播放状态
**GET** `/rooms/{room_id}/playback`

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "playback_state": {
            "playback_state": "paused",
            "current_time": 0,
            "playback_rate": 1.0,
            "media_url": "https://example.com/movie.mp4",
            "media_title": "肖申克的救赎",
            "last_updated": "2025-12-30T10:35:00Z",
            "version": 1
        }
    }
}
```

#### 4.5.2 播放控制
**POST** `/rooms/{room_id}/playback/play`

**请求体**:
```json
{
    "current_time": 0,
    "playback_rate": 1.0
}
```

#### 4.5.3 暂停控制
**POST** `/rooms/{room_id}/playback/pause`

#### 4.5.4 跳转控制
**POST** `/rooms/{room_id}/playback/seek`

**请求体**:
```json
{
    "current_time": 120.5
}
```

#### 4.5.5 播放速度控制
**POST** `/rooms/{room_id}/playback/speed`

**请求体**:
```json
{
    "playback_rate": 1.5
}
```

#### 4.5.6 媒体切换
**POST** `/rooms/{room_id}/playback/media`

**请求体**:
```json
{
    "media_url": "https://example.com/new-movie.mp4",
    "media_title": "新电影",
    "media_type": "video"
}
```

### 4.6 系统配置

#### 4.6.1 获取配置
**GET** `/configs`

**查询参数**:
- `key`: string, 指定配置键（可选）

**响应体**:
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "configs": [
            {
                "id": "app_name",
                "config_key": "app_name",
                "config_value": "小窝同步观影平台",
                "config_type": "string",
                "description": "应用名称",
                "is_editable": true
            }
        ]
    }
}
```

#### 4.6.2 更新配置
**PUT** `/configs/{config_key}`

**请求体**:
```json
{
    "config_value": "新的配置值",
    "description": "新的描述"
}
```

---

## 5. WebSocket事件契约

### 5.1 连接建立
**客户端发送**:
```json
{
    "type": "auth",
    "session_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "room_id": "ABC123"
}
```

**服务器响应**:
```json
{
    "type": "auth_success",
    "data": {
        "session": {},
        "room": {},
        "members": []
    }
}
```

### 5.2 播放状态同步事件

#### 5.2.1 播放状态变更通知
**服务器广播**:
```json
{
    "type": "playback_update",
    "data": {
        "room_id": "ABC123",
        "session_id": "550e8400-e29b-41d4-a716-446655440000",
        "action": "play",
        "current_time": 125.5,
        "playback_rate": 1.0,
        "version": 2,
        "timestamp": "2025-12-30T10:35:00Z"
    }
}
```

#### 5.2.2 播放状态请求
**客户端发送**:
```json
{
    "type": "get_playback_state",
    "room_id": "ABC123"
}
```

**服务器响应**:
```json
{
    "type": "playback_state",
    "data": {
        "playback_state": "playing",
        "current_time": 125.5,
        "playback_rate": 1.0,
        "version": 2
    }
}
```

### 5.3 消息事件

#### 5.3.1 新消息通知
**服务器广播**:
```json
{
    "type": "new_message",
    "data": {
        "message": {
            "id": "550e8400-e29b-41d4-a716-446655440002",
            "room_id": "ABC123",
            "session_id": "550e8400-e29b-41d4-a716-446655440000",
            "nickname": "快乐的小熊猫",
            "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=panda",
            "message_type": "chat",
            "content": "这部电影真的太棒了！",
            "created_at": "2025-12-30T10:35:00Z"
        }
    }
}
```

### 5.4 成员状态事件

#### 5.4.1 成员加入通知
**服务器广播**:
```json
{
    "type": "member_joined",
    "data": {
        "member": {
            "session_id": "550e8400-e29b-41d4-a716-446655440000",
            "nickname": "快乐的小熊猫",
            "avatar": "https://api.dicebear.com/7.x/avataaars/svg?seed=panda",
            "role": "guest",
            "joined_at": "2025-12-30T10:35:00Z"
        }
    }
}
```

#### 5.4.2 成员离开通知
**服务器广播**:
```json
{
    "type": "member_left",
    "data": {
        "session_id": "550e8400-e29b-41d4-a716-446655440000",
        "left_at": "2025-12-30T10:40:00Z"
    }
}
```

### 5.5 心跳检测
**客户端发送**:
```json
{
    "type": "ping",
    "timestamp": "2025-12-30T10:35:00Z"
}
```

**服务器响应**:
```json
{
    "type": "pong",
    "timestamp": "2025-12-30T10:35:00Z"
}
```

---

## 6. 错误码定义

### 6.1 系统错误码（0-999）
- `0` - 成功
- `1` - 未知错误
- `100` - 请求参数错误
- `101` - 请求体格式错误
- `102` - 请求参数缺失

### 6.2 认证错误码（1000-1999）
- `1001` - 未提供认证信息
- `1002` - 认证信息无效
- `1003` - 会话已过期
- `1004` - 用户不存在

### 6.3 权限错误码（2000-2999）
- `2001` - 权限不足
- `2002` - 不是房间创建者
- `2003` - 房间访问被拒绝
- `2004` - 房间密码错误

### 6.5 资源错误码（3000-3999）
- `3001` - 房间不存在
- `3002` - 会话不存在
- `3003` - 消息不存在
- `3004` - 房间已满员
- `3005` - 已在房间中

### 6.6 业务错误码（4000-4999）
- `4001` - 房间名称已存在
- `4002` - 媒体URL无效
- `4003` - 播放状态冲突
- `4004` - 版本冲突

### 6.7 系统限制错误码（5000-5999）
- `5001` - 请求频率过高
- `5002` - 房间数量超限
- `5003` - 消息过长
- `5004` - 文件大小超限

---

## 7. 请求示例

### 7.1 创建房间完整流程

```bash
# 1. 创建用户会话
curl -X POST http://localhost:8080/api/v1/sessions \
  -H "Content-Type: application/json" \
  -d '{"nickname": "快乐的小熊猫"}'

# 响应获取 session_id 和 token
# session_id: "550e8400-e29b-41d4-a716-446655440000"
# token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# 2. 创建房间
curl -X POST http://localhost:8080/api/v1/rooms \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -d '{
    "name": "周末电影时光",
    "description": "一起观看经典电影",
    "is_private": false,
    "media_url": "https://example.com/movie.mp4",
    "media_title": "肖申克的救赎"
  }'

# 响应获取 room_id: "ABC123"
```

### 7.2 加入房间流程

```bash
# 1. 加入房间
curl -X POST http://localhost:8080/api/v1/rooms/ABC123/join \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}"

# 2. 获取房间信息
curl -X GET "http://localhost:8080/api/v1/rooms/ABC123?include_members=true" \
  -H "Authorization: Bearer {session_token}"

# 3. 发送消息
curl -X POST http://localhost:8080/api/v1/rooms/ABC123/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}" \
  -d '{
    "message_type": "chat",
    "content": "这部电影真的太棒了！"
  }'
```

### 7.3 播放控制流程

```bash
# 1. 获取当前播放状态
curl -X GET http://localhost:8080/api/v1/rooms/ABC123/playback \
  -H "Authorization: Bearer {session_token}"

# 2. 开始播放
curl -X POST http://localhost:8080/api/v1/rooms/ABC123/playback/play \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}" \
  -d '{"current_time": 0, "playback_rate": 1.0}'

# 3. 暂停播放
curl -X POST http://localhost:8080/api/v1/rooms/ABC123/playback/pause \
  -H "Authorization: Bearer {session_token}"

# 4. 跳转播放位置
curl -X POST http://localhost:8080/api/v1/rooms/ABC123/playback/seek \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {session_token}" \
  -d '{"current_time": 120.5}'
```

---

**文档版本**: v1.0  
**最后更新**: 2025-12-30  
**负责人**: 后盾（后端架构师）

此API契约文档为前端开发和后端API实现提供了完整的规范，确保双方开发的一致性和可靠性。