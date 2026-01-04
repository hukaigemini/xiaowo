# 小窝项目详细设计文档 (Phase 2)

## 执行摘要

本文档是小窝同步观影项目的详细设计文档，基于Phase 1的技术架构方案，输出完整的数据库设计、API契约和核心业务逻辑设计。

**核心设计原则**：
- **MVP优先**：快速迭代，最小可行产品先行
- **技术成熟**：采用经过验证的技术栈
- **可扩展性**：预留扩展空间，支持业务增长
- **开发效率**：平衡开发速度与代码质量

---

## 1. Schema设计：数据库架构与ER图

### 1.1 核心业务ER图描述

```
┌─────────────────┐     ┌──────────────────┐     ┌─────────────────┐
│   UserSession   │     │      Room        │     │   RoomMember    │
│  (用户会话)      │     │    (房间)        │     │   (房间成员)     │
├─────────────────┤     ├──────────────────┤     ├─────────────────┤
│ *id: UUID       │     │ *id: RoomCode    │     │ *id: UUID       │
│  nickname       │     │  name            │     │  room_id        │
│  avatar         │     │  description     │     │  session_id     │
│  room_id (FK)   │────▶│  creator_id      │     │  nickname       │
│  created_at     │     │  is_private      │     │  avatar         │
└─────────────────┘     │  password        │     │  role           │
                        │  max_users       │     │  joined_at      │
┌─────────────────┐     │  status          │     │  last_seen_at   │
│  PlaybackState  │     │  media_url       │◀────┤  created_at     │
│  (播放状态)      │     │  playback_state  │     └─────────────────┘
├─────────────────┤     │  current_time    │              ▲
│ *room_id        │     │  duration        │              │
│  is_playing     │     │  version         │              │
│  current_time   │     │  last_active_at  │              │
│  duration       │     │  created_at      │              │
│  playback_rate  │     │  updated_at      │              │
│  updated_at     │     └──────────────────┘              │
└─────────────────┘              │                         │
                                │                         │
┌─────────────────┐             │                         │
│  SyncMessage    │             │                         │
│  (同步消息)      │             │                         │
├─────────────────┤             │                         │
│ *id: UUID       │             │                         │
│  room_id        │─────────────┘                         │
│  session_id     │                                       │
│  message_type   │                                       │
│  payload        │                                       │
│  created_at     │                                       │
└─────────────────┘                                       │
                                                          │
┌─────────────────┐                                       │
│  RoomActivity   │                                       │
│  (房间活动日志)  │                                       │
├─────────────────┤                                       │
│ *id: UUID       │                                       │
│  room_id        │───────────────────────────────────────┘
│  activity_type  │
│  actor_session  │
│  details        │
│  created_at     │
└─────────────────┘
```

**关系说明**：
- **UserSession ↔ Room**: 一对多关系（一个会话只能在一个房间，一个房间有多个会话）
- **Room ↔ RoomMember**: 一对多关系（一个房间有多个成员）
- **UserSession ↔ RoomMember**: 一对一关系（每个会话对应一个成员身份）
- **Room ↔ PlaybackState**: 一对一关系（每个房间有唯一的播放状态）
- **Room ↔ SyncMessage**: 一对多关系（房间接收多个同步消息）
- **Room ↔ RoomActivity**: 一对多关系（房间记录多次活动日志）

### 1.2 SQLite DDL建表语句（完整版）

```sql
-- 启用外键约束检查
PRAGMA foreign_keys = ON;

-- 用户会话表（临时会话）
CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,                    -- UUID主键
    nickname TEXT NOT NULL,                 -- 用户昵称
    avatar TEXT NOT NULL,                   -- 头像URL
    room_id TEXT,                          -- 当前所在房间ID（可空）
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    
    -- 字段注释：
    -- nickname: 趣味性昵称，如"快乐的考拉"
    -- avatar: 随机头像URL
    -- room_id: 关联到房间，支持会话迁移
    CONSTRAINT fk_user_session_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON DELETE SET NULL ON UPDATE CASCADE
);

-- 房间表（核心实体）
CREATE TABLE rooms (
    id TEXT PRIMARY KEY,                    -- 6位房间号主键
    name TEXT NOT NULL,                     -- 房间名称
    description TEXT,                       -- 房间描述
    creator_session_id TEXT NOT NULL,       -- 创建者会话ID
    is_private BOOLEAN DEFAULT 0,          -- 是否私密房间
    password TEXT,                          -- 房间密码
    max_users INTEGER DEFAULT 7,           -- 最大用户数（固定7人）
    status TEXT DEFAULT 'active',          -- 房间状态：active/inactive/archived
    media_url TEXT NOT NULL,               -- 媒体资源URL
    playback_state TEXT DEFAULT 'paused',  -- 播放状态：playing/paused/buffering
    current_time REAL DEFAULT 0,           -- 当前播放时间（秒）
    duration REAL DEFAULT 0,               -- 媒体总时长（秒）
    playback_rate REAL DEFAULT 1.0,        -- 播放倍速
    settings TEXT DEFAULT '{}',            -- JSON格式房间设置
    version INTEGER DEFAULT 0,             -- 乐观锁版本号
    last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 最后活跃时间
    last_member_left_at DATETIME,          -- 最后成员离开时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,      -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,      -- 更新时间
    
    -- 字段注释：
    -- id: 6位数字+字母混合房间号，如"8XA92B"
    -- version: 用于并发控制，防止脏写
    -- settings: 包含音量、画质等用户偏好设置
    -- last_member_left_at: 用于房间自动清理机制
    
    CONSTRAINT fk_room_creator 
        FOREIGN KEY (creator_session_id) REFERENCES user_sessions(id) 
        ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT chk_room_max_users CHECK (max_users > 0 AND max_users <= 20),
    CONSTRAINT chk_room_playback_rate CHECK (playback_rate >= 0.25 AND playback_rate <= 4.0),
    CONSTRAINT chk_room_current_time CHECK (current_time >= 0),
    CONSTRAINT chk_room_duration CHECK (duration >= 0)
);

-- 房间成员表（关联表）
CREATE TABLE room_members (
    id TEXT PRIMARY KEY,                    -- UUID主键
    room_id TEXT NOT NULL,                  -- 房间ID
    session_id TEXT NOT NULL,              -- 用户会话ID
    nickname TEXT NOT NULL,                 -- 成员昵称（冗余存储，便于快速查询）
    avatar TEXT NOT NULL,                   -- 成员头像（冗余存储）
    role TEXT DEFAULT 'guest',             -- 角色：host/guest/admin
    permissions TEXT DEFAULT '{}',          -- JSON格式权限配置
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 加入时间
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 最后在线时间
    left_at DATETIME,                       -- 离开时间（可空，表示当前在线）
    is_active BOOLEAN DEFAULT 1,           -- 活跃成员标记
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,   -- 创建时间
    
    -- 字段注释：
    -- nickname/avatar: 冗余存储，避免频繁JOIN查询
    -- role: host（房主）、guest（普通成员）、admin（管理员）
    -- permissions: 详细的权限控制，如播放控制、成员管理等
    -- is_active: 软删除标记，支持成员历史记录
    
    CONSTRAINT fk_member_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_member_session 
        FOREIGN KEY (session_id) REFERENCES user_sessions(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT uk_member_room_session UNIQUE (room_id, session_id),
    CONSTRAINT chk_member_role CHECK (role IN ('host', 'guest', 'admin'))
);

-- 播放状态表（房间播放状态快照）
CREATE TABLE playback_states (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),  -- UUID主键
    room_id TEXT NOT NULL UNIQUE,          -- 房间ID（唯一约束）
    is_playing BOOLEAN DEFAULT 0,         -- 是否正在播放
    current_time REAL DEFAULT 0,          -- 当前播放时间
    duration REAL DEFAULT 0,              -- 媒体总时长
    playback_rate REAL DEFAULT 1.0,       -- 播放倍速
    volume REAL DEFAULT 1.0,              -- 音量（0.0-1.0）
    quality TEXT DEFAULT 'auto',          -- 播放画质：auto/1080p/720p/480p
    buffering_progress REAL DEFAULT 0,    -- 缓冲进度（0.0-1.0）
    error_message TEXT,                   -- 播放错误信息
    updated_by_session TEXT,              -- 最后更新者会话ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 更新时间
    
    -- 字段注释：
    -- room_id: 与rooms表一对一关系，确保房间播放状态一致性
    -- buffering_progress: 支持流媒体播放的缓冲状态
    -- error_message: 播放错误信息，便于问题排查
    
    CONSTRAINT fk_playback_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_playback_updater 
        FOREIGN KEY (updated_by_session) REFERENCES user_sessions(id) 
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT chk_playback_rate CHECK (playback_rate >= 0.25 AND playback_rate <= 4.0),
    CONSTRAINT chk_playback_volume CHECK (volume >= 0.0 AND volume <= 1.0),
    CONSTRAINT chk_playback_buffering CHECK (buffering_progress >= 0.0 AND buffering_progress <= 1.0)
);

-- 同步消息表（WebSocket消息持久化，用于消息重放）
CREATE TABLE sync_messages (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),  -- UUID主键
    room_id TEXT NOT NULL,                 -- 房间ID
    session_id TEXT NOT NULL,             -- 发送者会话ID
    message_type TEXT NOT NULL,           -- 消息类型：play/pause/seek/sync/chat
    payload TEXT NOT NULL,                -- JSON格式消息载荷
    client_timestamp DATETIME NOT NULL,   -- 客户端时间戳
    server_timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 服务器时间戳
    processed BOOLEAN DEFAULT 0,          -- 是否已处理
    processing_error TEXT,                -- 处理错误信息
    broadcast_to TEXT DEFAULT 'all',      -- 广播目标：all/specific_users
    priority INTEGER DEFAULT 5,           -- 消息优先级（1-10，10最高）
    expires_at DATETIME,                  -- 消息过期时间
    
    -- 字段注释：
    -- message_type: 支持多种消息类型，便于扩展
    -- payload: 灵活的JSON载荷，支持不同消息格式
    -- client_timestamp: 用于时间同步和延迟计算
    -- processed: 消息处理状态，支持重试机制
    -- expires_at: 自动清理过期消息
    
    CONSTRAINT fk_message_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_message_session 
        FOREIGN KEY (session_id) REFERENCES user_sessions(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT chk_message_type CHECK (message_type IN ('play', 'pause', 'seek', 'sync', 'chat', 'system')),
    CONSTRAINT chk_message_priority CHECK (priority >= 1 AND priority <= 10)
);

-- 房间活动日志表（审计和统计）
CREATE TABLE room_activities (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),  -- UUID主键
    room_id TEXT NOT NULL,                 -- 房间ID
    activity_type TEXT NOT NULL,           -- 活动类型
    actor_session TEXT NOT NULL,           -- 活动执行者会话ID
    target_session TEXT,                   -- 目标会话ID（可空）
    details TEXT,                          -- JSON格式详细信息
    metadata TEXT,                         -- JSON格式元数据
    ip_address TEXT,                       -- IP地址（用于安全审计）
    user_agent TEXT,                       -- 用户代理
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,  -- 创建时间
    
    -- 字段注释：
    -- activity_type: room_created/member_joined/member_left/play_started 等
    -- details: 详细的操作信息，如播放位置、错误信息等
    -- metadata: 扩展信息，如设备信息、网络状态等
    -- ip_address/user_agent: 安全审计需要
    
    CONSTRAINT fk_activity_room 
        FOREIGN KEY (room_id) REFERENCES rooms(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_activity_actor 
        FOREIGN KEY (actor_session) REFERENCES user_sessions(id) 
        ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_activity_target 
        FOREIGN KEY (target_session) REFERENCES user_sessions(id) 
        ON DELETE SET NULL ON UPDATE CASCADE,
    CONSTRAINT chk_activity_type CHECK (activity_type IN (
        'room_created', 'room_updated', 'room_archived',
        'member_joined', 'member_left', 'member_kicked',
        'play_started', 'play_paused', 'play_seeked',
        'sync_error', 'connection_lost', 'connection_restored'
    ))
);

-- 索引优化（查询性能）
CREATE INDEX idx_user_sessions_room ON user_sessions(room_id);
CREATE INDEX idx_user_sessions_created ON user_sessions(created_at);

CREATE INDEX idx_rooms_creator ON rooms(creator_session_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_last_active ON rooms(last_active_at);
CREATE INDEX idx_rooms_created ON rooms(created_at);
CREATE INDEX idx_rooms_media_url ON rooms(media_url);

CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_session ON room_members(session_id);
CREATE INDEX idx_room_members_role ON room_members(role);
CREATE INDEX idx_room_members_active ON room_members(is_active, left_at);
CREATE INDEX idx_room_members_joined ON room_members(joined_at);

CREATE INDEX idx_playback_room ON playback_states(room_id);
CREATE INDEX idx_playback_updated ON playback_states(updated_at);

CREATE INDEX idx_sync_messages_room ON sync_messages(room_id);
CREATE INDEX idx_sync_messages_session ON sync_messages(session_id);
CREATE INDEX idx_sync_messages_type ON sync_messages(message_type);
CREATE INDEX idx_sync_messages_timestamp ON sync_messages(server_timestamp);
CREATE INDEX idx_sync_messages_processed ON sync_messages(processed, server_timestamp);
CREATE INDEX idx_sync_messages_expires ON sync_messages(expires_at) WHERE expires_at IS NOT NULL;

CREATE INDEX idx_room_activities_room ON room_activities(room_id);
CREATE INDEX idx_room_activities_type ON room_activities(activity_type);
CREATE INDEX idx_room_activities_actor ON room_activities(actor_session);
CREATE INDEX idx_room_activities_created ON room_activities(created_at);

-- 复合索引（多字段查询优化）
CREATE INDEX idx_room_members_room_active ON room_members(room_id, is_active, left_at);
CREATE INDEX idx_sync_messages_room_type_time ON sync_messages(room_id, message_type, server_timestamp);
CREATE INDEX idx_room_activities_room_type_time ON room_activities(room_id, activity_type, created_at);

-- 触发器：自动更新 updated_at 字段
CREATE TRIGGER user_sessions_update_updated_at 
    AFTER UPDATE ON user_sessions
    BEGIN
        UPDATE user_sessions SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER rooms_update_updated_at 
    AFTER UPDATE ON rooms
    BEGIN
        UPDATE rooms SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER playback_states_update_updated_at 
    AFTER UPDATE ON playback_states
    BEGIN
        UPDATE playback_states SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 触发器：房间成员离开时自动更新 last_seen_at
CREATE TRIGGER room_members_update_last_seen 
    AFTER UPDATE OF left_at ON room_members
    WHEN NEW.left_at IS NOT NULL
    BEGIN
        UPDATE room_members SET last_seen_at = NEW.left_at WHERE id = NEW.id;
    END;

-- 触发器：房间状态变化时更新房间活跃时间
CREATE TRIGGER rooms_update_last_active_on_playback 
    AFTER UPDATE OF playback_state ON rooms
    WHEN NEW.playback_state = 'playing' AND OLD.playback_state != 'playing'
    BEGIN
        UPDATE rooms SET last_active_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- 触发器：自动创建播放状态记录
CREATE TRIGGER rooms_create_playback_state 
    AFTER INSERT ON rooms
    BEGIN
        INSERT INTO playback_states (room_id, created_at, updated_at) 
        VALUES (NEW.id, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
    END;

-- 触发器：房间成员变化时更新房间成员离开时间
CREATE TRIGGER rooms_update_member_left_time 
    AFTER DELETE ON room_members
    WHEN (SELECT COUNT(*) FROM room_members WHERE room_id = OLD.room_id AND is_active = 1) = 0
    BEGIN
        UPDATE rooms SET last_member_left_at = CURRENT_TIMESTAMP WHERE id = OLD.room_id;
    END;

-- 房间自动清理机制视图
CREATE VIEW room_cleanup_candidates AS
SELECT 
    r.id,
    r.name,
    r.status,
    r.last_member_left_at,
    r.created_at,
    COUNT(rm.id) as active_member_count,
    CASE 
        WHEN r.last_member_left_at IS NOT NULL 
        THEN (julianday('now') - julianday(r.last_member_left_at)) * 24 * 60  -- 分钟
        ELSE NULL 
    END as minutes_since_last_member_left
FROM rooms r
LEFT JOIN room_members rm ON r.id = rm.room_id AND rm.is_active = 1
WHERE r.status = 'active'
GROUP BY r.id, r.name, r.status, r.last_member_left_at, r.created_at
HAVING active_member_count = 0;

-- 性能统计视图
CREATE VIEW room_performance_stats AS
SELECT 
    r.id as room_id,
    r.name,
    COUNT(DISTINCT rm.id) as total_members,
    COUNT(sm.id) as total_sync_messages,
    AVG(
        CASE 
            WHEN sm.message_type = 'sync' 
            THEN CAST(json_extract(sm.payload, '$.delay_ms') AS REAL)
            ELSE NULL 
        END
    ) as avg_sync_delay_ms,
    MAX(sm.server_timestamp) as last_activity_time
FROM rooms r
LEFT JOIN room_members rm ON r.id = rm.room_id
LEFT JOIN sync_messages sm ON r.id = sm.room_id
WHERE r.status = 'active'
GROUP BY r.id, r.name;
```

---

## 2. API契约：RESTful接口规范与OpenAPI定义

### 2.1 RESTful接口规范

**基础URL**: `https://api.xiaowo.com/v1` (生产环境) / `http://localhost:8080/v1` (开发环境)

**统一响应格式**:
```json
{
  "code": 200,
  "message": "success", 
  "data": {},
  "request_id": "uuid-string",
  "timestamp": "2025-12-30T10:30:00Z"
}
```

**错误码体系**:
- **系统级错误 (1000-1999)**: 1001-内部错误, 1002-数据库错误, 1003-网络错误
- **认证级错误 (2000-2999)**: 2001-无效Token, 2002-Token过期, 2003-权限不足
- **业务级错误 (3000-3999)**: 3001-房间不存在, 3002-房间已满, 3003-用户已在房间中
- **参数级错误 (4000-4999)**: 4001-参数缺失, 4002-参数无效, 4003-参数过长

**核心接口列表**:

**A. 房间管理接口**
```
POST   /rooms              # 创建房间
GET    /rooms/{room_id}    # 获取房间详情
PUT    /rooms/{room_id}    # 更新房间信息
DELETE /rooms/{room_id}    # 删除/归档房间
GET    /rooms              # 房间列表（支持分页和筛选）
```

**B. 用户会话接口**
```
POST   /sessions           # 创建用户会话
GET    /sessions/{session_id}  # 获取会话信息
PUT    /sessions/{session_id}  # 更新会话信息
DELETE /sessions/{session_id}  # 销毁会话
```

**C. 成员管理接口**
```
POST   /rooms/{room_id}/members     # 加入房间
DELETE /rooms/{room_id}/members/{session_id}  # 离开房间
GET    /rooms/{room_id}/members     # 获取房间成员列表
PUT    /rooms/{room_id}/members/{session_id}  # 更新成员信息
```

**D. 播放控制接口**
```
GET    /rooms/{room_id}/playback    # 获取播放状态
PUT    /rooms/{room_id}/playback    # 更新播放状态
POST   /rooms/{room_id}/play        # 播放
POST   /rooms/{room_id}/pause       # 暂停
POST   /rooms/{room_id}/seek        # 跳转
POST   /rooms/{room_id}/sync        # 强制同步
```

**E. 系统接口**
```
GET    /health                    # 健康检查
GET    /stats                     # 系统统计
GET    /rooms/{room_id}/stats     # 房间统计
```

### 2.2 OpenAPI/Swagger接口定义

完整的OpenAPI 3.0.3规范已保存在 `docs/openapi.yaml` 文件中，包含：
- 20+个核心接口的详细定义
- 完整的请求/响应模型
- 错误处理机制
- 认证和权限控制
- 接口示例和说明

---

## 3. 核心逻辑：业务流程时序图与异常处理设计

### 3.1 核心业务流程时序图

**A. 房间创建与初始化流程**
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │     │   Backend    │     │  SQLite DB  │     │  WebSocket  │     │  Logging    │
│   (用户)     │     │   (API)      │     │  (数据库)    │     │    Hub      │     │   System    │
└──────┬──────┘     └──────┬───────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │                  │                  │
       │ POST /rooms       │                   │                  │                  │
       │ {name, media_url} │                   │                  │                  │
       │------------------>│                   │                  │                  │
       │                   │ 1. 验证参数       │                  │                  │
       │                   │ 2. 生成房间ID     │                  │                  │
       │                   │ 3. 创建房间记录   │                  │                  │
       │                   │----------------->│ INSERT rooms     │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 4. 创建播放状态  │                  │                  │
       │                   │----------------->│ INSERT playback  │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 5. 记录活动日志  │                  │                  │
       │                   │----------------->│ INSERT activity  │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 6. 广播房间创建  │                  │ 广播房间创建消息  │
       │                   │----------------->│----------------->│----------------->│
       │                   │ 7. 记录操作日志  │                  │              记录成功
       │                   │----------------->│----------------->│----------------->│
       │                   │              ←--│-- 201 Created --│-- 房间ID: 8XA92B --│
       │<------------------│              ←--│-- 房间对象数据 --│-- 播放状态 --│
       │ 房间创建成功       │                   │                  │                  │
       │ 房间ID: 8XA92B     │                   │                  │                  │
       │                   │                   │                  │                  │
```

**B. 用户加入房间与实时同步流程**
```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │     │   Backend    │     │  SQLite DB  │     │  WebSocket  │     │  Media      │
│   (用户)     │     │   (API)      │     │  (数据库)    │     │    Hub      │     │  Player     │
└──────┬──────┘     └──────┬───────┘     └──────┬──────┘     └──────┬──────┘     └──────┬──────┘
       │                   │                   │                  │                  │
       │ 1. 加入房间        │                   │                  │                  │
       │ POST /rooms/xx/members│                │                  │                  │
       │------------------>│                   │                  │                  │
       │                   │ 2. 权限验证       │                  │                  │
       │                   │ 3. 检查房间状态   │                  │                  │
       │                   │----------------->│ SELECT room      │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-room data-------|                  │
       │                   │ 4. 检查容量      │                  │                  │
       │                   │ 5. 创建成员记录  │                  │                  │
       │                   │----------------->│ INSERT member    │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 6. 更新会话信息  │                  │                  │
       │                   │----------------->│ UPDATE session   │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 7. 获取播放状态  │                  │                  │
       │                   │----------------->│ SELECT playback  │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-playback-------|                  │
       │                   │ 8. 广播新成员    │                  │                  │
       │                   │----------------->│----------------->│ 成员加入广播    │
       │                   │ 9. 记录活动      │                  │              ←--│--加入通知
       │                   │----------------->│ INSERT activity  │                  │
       │                   │                  │----------------->│                  │
       │                   │              ←--│-- 201 Created --│-- 成员信息 --│
       │<------------------│              ←--│-- 播放状态数据 --│-- 同步数据 --│
       │ 加入成功           │                   │                  │                  │
       │ 播放状态: paused   │                   │                  │                  │
       │                   │                   │                  │                  │
       │ 10. 建立WebSocket │                   │                  │                  │
       │ 连接到实时同步    │                   │                  │                  │
       │<==============WebSocket Connect=============>│                  │                  │
       │                   │ 11. WebSocket注册│                  │                  │
       │                   │----------------->│ REGISTER session │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-Registered------│                  │
       │                   │ 12. 发送同步消息 │                  │                  │
       │------------------>│----------------->│----------------->│ SYNC_REQUEST   │
       │ 播放:play()       │                   │                  │              ←--│--SYNC_RESPONSE
       │                   │ 13. 播放控制验证 │                  │                  │
       │                   │----------------->│ SELECT permissions│                 │
       │                   │                  │----------------->│                  │
       │                   │                  │<-permissions-----|                  │
       │                   │ 14. 更新播放状态 │                  │                  │
       │                   │----------------->│ UPDATE playback  │                  │
       │                   │                  │----------------->│                  │
       │                   │                  │<-SQLite OK-------|                  │
       │                   │ 15. 广播播放状态 │                  │                  │
       │                   │----------------->│----------------->│ PLAY_BROADCAST│
       │              ←--播放状态更新--│              ←--所有成员--│              ←--│--更新播放
       │                   │                   │                  │              同步播放
       │                   │                   │                  │                  │
```

### 3.2 异常处理与幂等性设计

**A. 异常处理策略**

**1. 分层异常处理架构**
```go
// 异常处理层次结构
type ErrorHandler struct {
    dbHandler    *DatabaseErrorHandler
    wsHandler    *WebSocketErrorHandler
    businessHandler *BusinessErrorHandler
    retryManager *RetryManager
}

// 数据库异常处理
type DatabaseErrorHandler struct{}

func (h *DatabaseErrorHandler) Handle(err error, ctx *RequestContext) *ErrorResponse {
    switch err {
    case sql.ErrNoRows:
        return &ErrorResponse{
            Code:    3001, // 资源不存在
            Message: "数据不存在",
            Details: fmt.Sprintf("查询的资源不存在: %s", ctx.ResourceID),
        }
    case context.DeadlineExceeded:
        return &ErrorResponse{
            Code:    1002, // 数据库超时
            Message: "数据库操作超时",
            Details: "请稍后重试",
        }
    default:
        // 记录详细错误信息
        log.WithFields(log.Fields{
            "error":     err.Error(),
            "resource":  ctx.ResourceID,
            "operation": ctx.Operation,
            "user_id":   ctx.UserID,
        }).Error("database_error")
        
        return &ErrorResponse{
            Code:    1002, // 数据库错误
            Message: "数据库操作失败",
            Details: "系统内部错误，请联系管理员",
        }
    }
}
```

**B. 幂等性设计**

**1. 幂等性保证机制**
```go
// 幂等性管理器
type IdempotencyManager struct {
    redis *redis.Client
    ttl   time.Duration
}

func NewIdempotencyManager(redisClient *redis.Client) *IdempotencyManager {
    return &IdempotencyManager{
        redis: redisClient,
        ttl:   24 * time.Hour, // 24小时TTL
    }
}

// 幂等性检查
func (m *IdempotencyManager) Check(idempotencyKey string) (*IdempotencyResult, error) {
    key := fmt.Sprintf("idempotency:%s", idempotencyKey)
    
    // 检查是否已存在
    result, err := m.redis.Get(key)
    if err != nil && err != redis.Nil {
        return nil, err
    }
    
    if err == nil {
        // 已存在，返回之前的结果
        var savedResult IdempotencyResult
        if err := json.Unmarshal([]byte(result), &savedResult); err != nil {
            return nil, err
        }
        return &savedResult, nil
    }
    
    // 不存在，创建新的幂等性记录
    return &IdempotencyResult{
        Key:        idempotencyKey,
        Status:     "pending",
        CreatedAt:  time.Now(),
        ExpiresAt:  time.Now().Add(m.ttl),
    }, nil
}
```

**C. 重试策略与降级方案**

**1. 重试管理器**
```go
// 重试策略
type RetryStrategy int

const (
    Immediate RetryStrategy = iota
    ExponentialBackoff
    LinearBackoff
)

// 重试配置
type RetryConfig struct {
    MaxAttempts    int           // 最大重试次数
    InitialDelay   time.Duration // 初始延迟
    MaxDelay       time.Duration // 最大延迟
    Strategy       RetryStrategy // 重试策略
    Jitter         bool          // 是否添加随机抖动
}

// 重试管理器
type RetryManager struct {
    configs map[string]*RetryConfig
}

func NewRetryManager() *RetryManager {
    return &RetryManager{
        configs: map[string]*RetryConfig{
            "database": {
                MaxAttempts:  3,
                InitialDelay: 100 * time.Millisecond,
                MaxDelay:     5 * time.Second,
                Strategy:     ExponentialBackoff,
                Jitter:       true,
            },
            "websocket": {
                MaxAttempts:  5,
                InitialDelay: 1 * time.Second,
                MaxDelay:     30 * time.Second,
                Strategy:     ExponentialBackoff,
                Jitter:       true,
            },
            "external_api": {
                MaxAttempts:  3,
                InitialDelay: 500 * time.Millisecond,
                MaxDelay:     10 * time.Second,
                Strategy:     ExponentialBackoff,
                Jitter:       true,
            },
        },
    }
}
```

**2. 降级策略实现**
```go
// 降级管理器
type DegradationManager struct {
    circuitBreakers map[string]*CircuitBreaker
    fallbacks       map[string]FallbackHandler
    metrics         *MetricsCollector
}

type CircuitState int

const (
    Closed CircuitState = iota
    Open
    HalfOpen
)

// 熔断器
type CircuitBreaker struct {
    name            string
    state           CircuitState
    failureCount    int
    successCount    int
    lastFailureTime time.Time
    config          CircuitBreakerConfig
}

type CircuitBreakerConfig struct {
    FailureThreshold int           // 失败阈值
    RecoveryTimeout  time.Duration // 恢复超时
    SuccessThreshold int           // 半开状态成功阈值
}

// 降级处理器
type FallbackHandler func(ctx context.Context, req interface{}) (interface{}, error)
```

---

## 4. 核心设计亮点总结

### 4.1 数据一致性保障
- **乐观锁版本控制**：防止并发冲突
- **外键约束**：确保数据完整性
- **事务边界**：明确划分操作范围

### 4.2 高可用架构
- **熔断器模式**：防止雪崩效应
- **多级降级策略**：保障核心功能
- **幂等性设计**：消除重复操作影响

### 4.3 性能优化
- **针对性索引策略**：优化查询性能
- **连接池和查询优化**：提升数据库性能
- **内存缓存与数据库WAL模式**：减少I/O操作

### 4.4 开发友好
- **完整的API文档**：便于前后端并行开发
- **清晰的代码组织结构**：便于维护和扩展
- **统一的错误处理和日志记录**：便于调试和监控

---

## 5. 实施建议

### 5.1 开发阶段划分
1. **第一阶段**：数据库初始化和基础API实现
2. **第二阶段**：WebSocket实时同步功能
3. **第三阶段**：异常处理和性能优化
4. **第四阶段**：测试和部署优化

### 5.2 技术要点
1. **数据库迁移**：使用提供的DDL脚本初始化数据库
2. **API开发**：严格按照OpenAPI文档实现接口
3. **WebSocket管理**：实现连接池和消息队列
4. **监控告警**：接入日志系统和性能监控

### 5.3 质量保证
1. **单元测试**：核心业务逻辑的全面测试
2. **集成测试**：端到端业务流程验证
3. **性能测试**：并发和负载测试
4. **安全测试**：权限和数据安全验证

---

**文档版本**: v1.0  
**创建日期**: 2025-12-30  
**最后更新**: 2025-12-30  
**作者**: 后盾（Principal Backend Architect）