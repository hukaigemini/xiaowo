-- 小窝项目 - 数据库初始化脚本
-- 版本: v1.0
-- 创建时间: 2025-12-30
-- 描述: 创建核心数据表和索引

PRAGMA foreign_keys = ON;

-- =============================================
-- 1. 用户会话表 (临时会话，无独立用户系统)
-- =============================================
CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,                    -- 会话唯一ID (UUID)
    nickname TEXT NOT NULL,                 -- 用户昵称 (随机生成的趣味昵称)
    avatar TEXT NOT NULL,                   -- 用户头像URL (随机生成)
    room_id TEXT,                          -- 当前所在房间ID (可为空)
    status TEXT DEFAULT 'online',           -- 会话状态: online/offline
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 最后在线时间
    expires_at DATETIME DEFAULT (datetime('now', '+7 days')), -- 过期时间 (7天后)
    deleted_at DATETIME,                    -- 软删除时间戳
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE SET NULL
);

-- 用户会话表索引
CREATE INDEX idx_user_sessions_room ON user_sessions(room_id);
CREATE INDEX idx_user_sessions_status ON user_sessions(status);
CREATE INDEX idx_user_sessions_expires ON user_sessions(expires_at);
CREATE INDEX idx_user_sessions_last_seen ON user_sessions(last_seen_at);

-- =============================================
-- 2. 房间表 (核心房间信息)
-- =============================================
CREATE TABLE rooms (
    id TEXT PRIMARY KEY,                    -- 房间ID (6位数字+字母混合房间号)
    name TEXT NOT NULL,                     -- 房间名称
    description TEXT,                       -- 房间描述
    creator_session_id TEXT NOT NULL,       -- 创建者会话ID
    is_private BOOLEAN DEFAULT 0,           -- 是否私密房间
    password TEXT,                          -- 房间密码 (如有)
    max_users INTEGER DEFAULT 7,            -- 最大用户数 (固定为7)
    status TEXT DEFAULT 'active',           -- 房间状态: active/inactive
    media_url TEXT NOT NULL,                -- 媒体资源URL
    media_type TEXT DEFAULT 'video',        -- 媒体类型: video/audio/stream
    media_title TEXT,                       -- 媒体标题
    media_duration REAL DEFAULT 0,          -- 媒体总时长 (秒)
    playback_state TEXT DEFAULT 'paused',   -- 播放状态: playing/paused/stopped
    current_time REAL DEFAULT 0,            -- 当前播放时间 (秒)
    playback_rate REAL DEFAULT 1.0,         -- 播放速率 (1.0=正常, 1.5=1.5倍速)
    settings TEXT DEFAULT '{}',             -- 房间设置 (JSON格式)
    version INTEGER DEFAULT 0,              -- 乐观锁版本号
    last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 最后活跃时间
    last_member_left_at DATETIME,           -- 最后一位成员离开时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,     -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,     -- 更新时间
    FOREIGN KEY (creator_session_id) REFERENCES user_sessions(id) ON DELETE CASCADE
);

-- 房间表索引
CREATE INDEX idx_rooms_creator ON rooms(creator_session_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_last_active ON rooms(last_active_at);
CREATE INDEX idx_rooms_is_private ON rooms(is_private);
CREATE INDEX idx_rooms_created_at ON rooms(created_at);

-- =============================================
-- 3. 房间成员表 (房间成员关系)
-- =============================================
CREATE TABLE room_members (
    id TEXT PRIMARY KEY,                    -- 记录唯一ID (UUID)
    room_id TEXT NOT NULL,                  -- 房间ID
    session_id TEXT NOT NULL,               -- 用户会话ID
    nickname TEXT NOT NULL,                 -- 用户昵称 (冗余存储，便于快速查询)
    avatar TEXT NOT NULL,                   -- 用户头像 (冗余存储)
    role TEXT DEFAULT 'guest',              -- 角色: host/guest
    permissions TEXT DEFAULT '{"control": true, "chat": true}', -- 权限设置 (JSON)
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 加入时间
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 最后在线时间
    left_at DATETIME,                      -- 离开时间 (可为空)
    is_active BOOLEAN DEFAULT 1,           -- 是否活跃成员 (用于软删除)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE CASCADE,
    UNIQUE(room_id, session_id) -- 同一用户不能重复加入同一房间
);

-- 房间成员表索引
CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_session ON room_members(session_id);
CREATE INDEX idx_room_members_active ON room_members(is_active, room_id);
CREATE INDEX idx_room_members_role ON room_members(room_id, role);
CREATE INDEX idx_room_members_joined ON room_members(joined_at);

-- =============================================
-- 4. 房间消息表 (聊天和系统消息)
-- =============================================
CREATE TABLE room_messages (
    id TEXT PRIMARY KEY,                    -- 消息唯一ID (UUID)
    room_id TEXT NOT NULL,                  -- 房间ID
    session_id TEXT NOT NULL,               -- 发送者会话ID
    message_type TEXT DEFAULT 'chat',       -- 消息类型: chat/system/notification
    content TEXT NOT NULL,                  -- 消息内容
    metadata TEXT DEFAULT '{}',             -- 消息元数据 (JSON格式)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE CASCADE
);

-- 房间消息表索引
CREATE INDEX idx_messages_room ON room_messages(room_id);
CREATE INDEX idx_messages_session ON room_messages(session_id);
CREATE INDEX idx_messages_type ON room_messages(message_type);
CREATE INDEX idx_messages_created ON room_messages(created_at);

-- =============================================
-- 5. 房间状态历史表 (播放状态变更历史)
-- =============================================
CREATE TABLE room_state_history (
    id TEXT PRIMARY KEY,                    -- 记录唯一ID (UUID)
    room_id TEXT NOT NULL,                  -- 房间ID
    session_id TEXT,                        -- 触发变更的会话ID (可为空，表示系统自动)
    action_type TEXT NOT NULL,              -- 操作类型: play/pause/seek/speed_change/media_change
    old_state TEXT,                         -- 变更前状态 (JSON格式)
    new_state TEXT NOT NULL,                -- 变更后状态 (JSON格式)
    current_time REAL,                      -- 当前播放时间
    playback_rate REAL,                     -- 播放速率
    media_url TEXT,                         -- 媒体URL (如变更)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE SET NULL
);

-- 房间状态历史表索引
CREATE INDEX idx_state_history_room ON room_state_history(room_id);
CREATE INDEX idx_state_history_session ON room_state_history(session_id);
CREATE INDEX idx_state_history_action ON room_state_history(action_type);
CREATE INDEX idx_state_history_created ON room_state_history(created_at);

-- =============================================
-- 6. 系统配置表 (应用配置和字典数据)
-- =============================================
CREATE TABLE system_configs (
    id TEXT PRIMARY KEY,                    -- 配置项ID
    config_key TEXT UNIQUE NOT NULL,        -- 配置键
    config_value TEXT NOT NULL,             -- 配置值
    config_type TEXT DEFAULT 'string',      -- 配置类型: string/number/boolean/json
    description TEXT,                       -- 配置描述
    is_editable BOOLEAN DEFAULT 1,          -- 是否可编辑
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 创建时间
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP  -- 更新时间
);

-- 系统配置表索引
CREATE INDEX idx_configs_key ON system_configs(config_key);

-- =============================================
-- 触发器: 自动更新 updated_at 字段
-- =============================================
CREATE TRIGGER update_rooms_updated_at 
    AFTER UPDATE ON rooms
    BEGIN
        UPDATE rooms SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

CREATE TRIGGER update_system_configs_updated_at 
    AFTER UPDATE ON system_configs
    BEGIN
        UPDATE system_configs SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;

-- =============================================
-- 触发器: 自动更新最后在线时间
-- =============================================
CREATE TRIGGER update_user_sessions_last_seen 
    AFTER UPDATE OF last_seen_at ON user_sessions
    BEGIN
        -- 更新房间最后活跃时间
        UPDATE rooms 
        SET last_active_at = CURRENT_TIMESTAMP 
        WHERE id = NEW.room_id;
    END;

-- =============================================
-- 触发器: 成员加入/离开时更新房间状态
-- =============================================
CREATE TRIGGER on_member_join 
    AFTER INSERT ON room_members
    WHEN NEW.is_active = 1
    BEGIN
        UPDATE rooms 
        SET last_active_at = CURRENT_TIMESTAMP,
            status = 'active'
        WHERE id = NEW.room_id;
    END;

CREATE TRIGGER on_member_leave 
    AFTER UPDATE ON room_members
    WHEN OLD.is_active = 1 AND NEW.is_active = 0
    BEGIN
        -- 检查是否还有活跃成员
        UPDATE rooms 
        SET last_member_left_at = CURRENT_TIMESTAMP,
            status = CASE 
                WHEN (SELECT COUNT(*) FROM room_members WHERE room_id = NEW.room_id AND is_active = 1) = 0 
                THEN 'inactive' 
                ELSE status 
            END
        WHERE id = NEW.room_id;
    END;

-- =============================================
-- 数据库初始化完成
-- =============================================
-- 此脚本创建了小窝项目的核心数据表
-- 包含用户会话、房间管理、成员关系、消息系统和状态历史
-- 支持完整的同步观影功能