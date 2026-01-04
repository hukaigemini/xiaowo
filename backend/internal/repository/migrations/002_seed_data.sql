-- 小窝项目 - 数据初始化脚本
-- 版本: v1.0
-- 创建时间: 2025-12-30
-- 描述: 预置基础数据和演示内容

-- =============================================
-- 1. 系统配置初始化
-- =============================================

-- 应用基础配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_editable) VALUES
('app_name', 'app_name', '小窝同步观影平台', 'string', '应用名称', 1),
('app_version', 'app_version', '1.0.0', 'string', '应用版本', 0),
('max_room_users', 'max_room_users', '7', 'number', '房间最大用户数', 0),
('default_room_timeout', 'default_room_timeout', '600', 'number', '房间默认超时时间(秒)', 1),
('websocket_heartbeat_interval', 'websocket_heartbeat_interval', '30', 'number', 'WebSocket心跳间隔(秒)', 1),
('message_retention_days', 'message_retention_days', '30', 'number', '消息保留天数', 1);

-- 演示片源配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_editable) VALUES
('demo_videos', 'demo_videos', '[
  {
    "title": "Big Buck Bunny",
    "url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
    "description": "开源动画短片，适合测试播放功能",
    "duration": 596,
    "thumbnail": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/BigBuckBunny.jpg"
  },
  {
    "title": "Elephants Dream", 
    "url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
    "description": "另一个经典开源测试视频",
    "duration": 653,
    "thumbnail": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ElephantsDream.jpg"
  },
  {
    "title": "For Bigger Blazes",
    "url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4", 
    "description": "动作短片，适合测试同步功能",
    "duration": 15,
    "thumbnail": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerBlazes.jpg"
  },
  {
    "title": "For Bigger Escapes",
    "url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
    "description": "短动作片，便于快速测试",
    "duration": 15,
    "thumbnail": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerEscapes.jpg"
  },
  {
    "title": "For Bigger Joyrides",
    "url": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4",
    "description": "汽车相关短片，视觉效果丰富",
    "duration": 15,
    "thumbnail": "https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/images/ForBiggerJoyrides.jpg"
  }
]', 'json', '演示视频列表配置', 1);

-- 用户昵称池配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_editable) VALUES
('nickname_pool', 'nickname_pool', '[
  "快乐的考拉", "悠闲的熊猫", "聪明的小狐狸", "活泼的小兔子", "温柔的猫咪",
  "勇敢的小狮子", "可爱的小熊", "聪明的小鸟", "温暖的小狗", "可爱的小猪",
  "神秘的小猫", "聪明的小象", "可爱的小鹿", "聪明的小羊", "勇敢的小虎",
  "温柔的小兔", "活泼的小鸟", "可爱的小鼠", "聪明的小狐", "勇敢的小豹",
  "可爱的小象", "温柔的小鹿", "活泼的小鸟", "聪明的小猪", "勇敢的小牛",
  "可爱的小羊", "温柔的小马", "聪明的小熊", "可爱的小狗", "勇敢的小猫",
  "温柔的小鸟", "活泼的小鼠", "聪明的小鹿", "可爱的小兔", "勇敢的小虎",
  "温柔的小象", "活泼的小豹", "可爱的小猪", "聪明的小狗", "勇敢的小狐",
  "温柔的小鹿", "可爱的小鸟", "活泼的小鼠", "聪明的小熊", "勇敢的小牛",
  "可爱的小羊", "温柔的小马", "聪明的小兔", "勇敢的小猫", "可爱的小虎",
  "活泼的小象", "温柔的小豹", "可爱的小猪", "聪明的小狗", "勇敢的小狐",
  "温柔的小鸟", "活泼的小鼠", "可爱的小鹿", "聪明的小熊", "勇敢的小牛",
  "可爱的小羊", "温柔的小马", "聪明的小兔", "勇敢的小猫", "可爱的小虎",
  "活泼的小象", "温柔的小豹", "可爱的小猪", "聪明的小狗", "勇敢的小狐"
]', 'json', '随机昵称生成池', 1);

-- 头像配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_editable) VALUES
('avatar_pool', 'avatar_pool', '[
  "🐨", "🐼", "🦊", "🐰", "🐱", "🦁", "🐻", "🐦", "🐶", "🐷",
  "🐸", "🐘", "🦒", "🐑", "🐅", "🐹", "🐧", "🐭", "🐨", "🐪",
  "🐮", "🐑", "🐎", "🐻", "🐶", "🐱", "🐦", "🐭", "🦌", "🐰",
  "🐅", "🐘", "🐆", "🐷", "🐶", "🦊", "🐦", "🐭", "🐻", "🐮",
  "🐑", "🐎", "🐱", "🐅", "🐘", "🐆", "🐷", "🐶", "🦊", "🐦",
  "🐭", "🐻", "🦌", "🐰", "🐅", "🐘", "🐆", "🐷", "🐶", "🦊",
  "🐦", "🐭", "🐻", "🐮", "🐑", "🐎", "🐱", "🐅", "🐘", "🐆"
]', 'json', '随机头像表情池', 1);

-- 房间设置默认配置
INSERT INTO system_configs (id, config_key, config_value, config_type, description, is_editable) VALUES
('default_room_settings', 'default_room_settings', '{
  "auto_play": true,
  "allow_control": true,
  "sync_tolerance": 2.0,
  "chat_enabled": true,
  "member_notifications": true,
  "playback_rate_control": true,
  "seek_control": true
}', 'json', '房间默认设置', 1);

-- =============================================
-- 2. 创建示例房间数据 (用于测试)
-- =============================================

-- 创建一个示例房间 (不包含成员，用于UI展示)
INSERT INTO rooms (
    id, name, description, creator_session_id, is_private, password, max_users, 
    status, media_url, media_title, media_duration, playback_state, current_time, 
    playback_rate, settings, version, last_active_at, created_at
) VALUES (
    'DEMO01', '示例观影房间', '这是一个示例房间，用于UI展示和功能测试', 
    'demo-creator-session', 0, NULL, 7, 'active',
    'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4',
    'Big Buck Bunny (示例视频)', 596, 'paused', 0, 1.0,
    '{"auto_play": true, "allow_control": true, "sync_tolerance": 2.0}', 
    0, datetime('now', '-1 hour'), datetime('now', '-1 hour')
);

-- =============================================
-- 3. 数据验证和清理脚本
-- =============================================

-- 清理过期的用户会话
DELETE FROM user_sessions WHERE expires_at < datetime('now');

-- 清理孤立的消息记录 (用户会话已删除但消息还在)
DELETE FROM room_messages 
WHERE session_id NOT IN (SELECT id FROM user_sessions);

-- 清理孤立的房间成员记录
DELETE FROM room_members 
WHERE session_id NOT IN (SELECT id FROM user_sessions)
   OR room_id NOT IN (SELECT id FROM rooms);

-- =============================================
-- 4. 数据完整性检查
-- =============================================

-- 检查数据一致性
SELECT 
    'user_sessions' as table_name,
    COUNT(*) as total_records,
    COUNT(CASE WHEN room_id IS NOT NULL THEN 1 END) as with_room,
    COUNT(CASE WHEN expires_at < datetime(''now'') THEN 1 END) as expired
FROM user_sessions
UNION ALL
SELECT 
    'rooms' as table_name,
    COUNT(*) as total_records,
    COUNT(CASE WHEN status = ''active'' THEN 1 END) as active_rooms,
    COUNT(CASE WHEN last_active_at < datetime(''now'', ''-1 hour'') THEN 1 END) as inactive_1h
FROM rooms
UNION ALL
SELECT 
    'room_members' as table_name,
    COUNT(*) as total_records,
    COUNT(CASE WHEN is_active = 1 THEN 1 END) as active_members,
    COUNT(CASE WHEN left_at IS NOT NULL THEN 1 END) as left_members
FROM room_members;

-- =============================================
-- 数据初始化完成
-- =============================================
-- 此脚本完成了：
-- 1. 系统基础配置设置
-- 2. 演示视频片源配置  
-- 3. 用户昵称和头像池配置
-- 4. 房间默认设置配置
-- 5. 示例房间数据
-- 6. 数据清理和验证脚本