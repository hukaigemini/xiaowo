-- 房间表
CREATE TABLE rooms (
    id TEXT PRIMARY KEY, -- SQLite使用TEXT存储UUID
    name TEXT NOT NULL,
    description TEXT,
    creator_session_id TEXT NOT NULL, -- 仅记录会话ID
    is_private INTEGER DEFAULT 0, -- SQLite使用INTEGER存储BOOLEAN (0/1)
    room_password TEXT,
    max_users INTEGER DEFAULT 50,
    settings TEXT, -- JSON存储为TEXT
    playback_state TEXT, -- JSON存储为TEXT
    version INTEGER DEFAULT 0, -- 乐观锁版本号
    status TEXT DEFAULT 'active',
    last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- 用于自动销毁
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 房间成员表
CREATE TABLE room_members (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    session_id TEXT NOT NULL, -- 对应临时会话
    role TEXT DEFAULT 'member',
    nickname TEXT,
    is_muted INTEGER DEFAULT 0,
    last_seen DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE
);

-- 复合唯一索引
CREATE UNIQUE INDEX idx_room_session ON room_members(room_id, session_id);

-- 索引优化
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_last_active ON rooms(last_active_at);
