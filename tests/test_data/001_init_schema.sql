-- 小沃项目测试数据库初始化脚本
-- 版本: v1.0
-- 创建时间: 2025-12-30

-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255) UNIQUE NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    avatar TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建房间表
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_private BOOLEAN DEFAULT FALSE,
    password VARCHAR(255),
    max_users INTEGER DEFAULT 10,
    media_url TEXT NOT NULL,
    media_title VARCHAR(255),
    media_type VARCHAR(50) DEFAULT 'video',
    media_duration INTEGER DEFAULT 0,
    settings JSONB DEFAULT '{}',
    creator_session_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建房间成员表
CREATE TABLE IF NOT EXISTS room_members (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) REFERENCES rooms(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,
    nickname VARCHAR(100) NOT NULL,
    avatar TEXT,
    role VARCHAR(50) DEFAULT 'member',
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMP,
    UNIQUE(room_id, session_id)
);

-- 创建播放状态表
CREATE TABLE IF NOT EXISTS playback_states (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) REFERENCES rooms(id) ON DELETE CASCADE,
    current_time DECIMAL DEFAULT 0,
    is_playing BOOLEAN DEFAULT FALSE,
    volume DECIMAL DEFAULT 1.0,
    quality VARCHAR(20) DEFAULT 'auto',
    updated_by VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS operation_logs (
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) REFERENCES rooms(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,
    operation_type VARCHAR(50) NOT NULL,
    operation_data JSONB,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_rooms_creator ON rooms(creator_session_id);
CREATE INDEX IF NOT EXISTS idx_room_members_room_id ON room_members(room_id);
CREATE INDEX IF NOT EXISTS idx_room_members_session_id ON room_members(session_id);
CREATE INDEX IF NOT EXISTS idx_playback_states_room_id ON playback_states(room_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_room_id ON operation_logs(room_id);
CREATE INDEX IF NOT EXISTS idx_operation_logs_timestamp ON operation_logs(timestamp);

-- 插入测试数据
INSERT INTO users (session_id, nickname, avatar) VALUES
    ('sess_test_001', '测试用户001', 'https://example.com/avatar1.jpg'),
    ('sess_test_002', '测试用户002', 'https://example.com/avatar2.jpg'),
    ('sess_test_003', '测试用户003', 'https://example.com/avatar3.jpg')
ON CONFLICT (session_id) DO NOTHING;
