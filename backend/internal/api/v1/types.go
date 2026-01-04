package v1

import (
	"errors"
	"strconv"
	"time"

	"xiaowo/backend/internal/model"
)

// ==================== 请求结构体 ====================

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,max=100" example:"电影之夜"`          // 房间名称
	Description string `json:"description" example:"一起看《阿凡达》"`                       // 房间描述
	IsPrivate   bool   `json:"is_private" example:"false"`                               // 是否私密房间
	Password    string `json:"password" example:""`                                      // 房间密码（私密房间必需）
	MaxUsers    int    `json:"max_users" binding:"min=1,max=1000" example:"10"`          // 最大用户数
	
	// 媒体信息
	MediaURL    string  `json:"media_url" binding:"required" example:"https://example.com/video.mp4"` // 媒体资源URL
	MediaType   string  `json:"media_type" example:"video"`                                         // 媒体类型: video/audio/stream
	MediaTitle  string  `json:"media_title" example:"阿凡达"`                                       // 媒体标题
	MediaDuration float64 `json:"media_duration" example:"7200"`                                    // 媒体总时长(秒)
	
	Settings    struct {
		AutoPlay       bool    `json:"auto_play" example:"true"`           // 自动播放
		AllowSeek      bool    `json:"allow_seek" example:"true"`          // 允许拖拽
		DefaultVolume  float64 `json:"default_volume" example:"1.0"`       // 默认音量
		Quality        string  `json:"quality" example:"auto"`             // 默认画质
		SubtitleOn     bool    `json:"subtitle_on" example:"false"`       // 字幕开关
		SubtitleLang   string  `json:"subtitle_lang" example:"zh"`        // 字幕语言
	} `json:"settings"` // 房间设置
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	RoomID      string `json:"room_id" binding:"required" example:"room_123"` // 房间ID
	DisplayName string `json:"display_name" example:"电影爱好者"` // 显示名称（可选，不提供则自动生成）
	Password    string `json:"password" example:""`             // 房间密码（私密房间必需）
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	Name        string `json:"name" binding:"max=100" example:"新的电影之夜"` // 房间名称
	Description string `json:"description" example:"今晚看什么电影？"`         // 房间描述
	IsPrivate   *bool  `json:"is_private"`                              // 是否私密房间
	Password    string `json:"password" example:""`                    // 房间密码
	MaxUsers    *int   `json:"max_users" binding:"min=1,max=1000"`     // 最大用户数
	
	// 媒体信息（可选更新）
	MediaURL    *string  `json:"media_url" example:"https://example.com/video2.mp4"` // 媒体资源URL
	MediaType   *string  `json:"media_type" example:"video"`                         // 媒体类型
	MediaTitle  *string  `json:"media_title" example:"新的视频标题"`                 // 媒体标题
	MediaDuration *float64 `json:"media_duration" example:"5400"`                    // 媒体总时长(秒)
	
	Settings    struct {
		AutoPlay       *bool    `json:"auto_play"`           // 自动播放
		AllowSeek      *bool    `json:"allow_seek"`          // 允许拖拽
		DefaultVolume  *float64 `json:"default_volume"`      // 默认音量
		Quality        *string  `json:"quality"`             // 默认画质
		SubtitleOn     *bool    `json:"subtitle_on"`         // 字幕开关
		SubtitleLang   *string  `json:"subtitle_lang"`       // 字幕语言
	} `json:"settings"` // 房间设置
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Nickname string `json:"nickname" example:"电影爱好者"` // 昵称（可选，不提供则自动生成）
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	Nickname string `json:"nickname" example:"新的昵称"` // 新昵称
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"` // 头像URL
}

// ==================== 响应结构体 ====================

// RoomResponse 房间响应
type RoomResponse struct {
	Room        *model.Room `json:"room"`                  // 房间信息
	MemberCount int         `json:"member_count"`          // 当前成员数量
	CreatedBy   string      `json:"created_by"`            // 创建者显示名称
	CreatedAt   time.Time   `json:"created_at"`            // 创建时间
	IsCreator   bool        `json:"is_creator"`            // 当前用户是否为创建者
}

// SessionResponse 会话响应
type SessionResponse struct {
	SessionID  string    `json:"session_id"`  // 会话ID
	Nickname   string    `json:"nickname"`    // 用户昵称
	Avatar     string    `json:"avatar"`      // 头像URL
	RoomID     string    `json:"room_id"`     // 所在房间ID
	Status     string    `json:"status"`      // 会话状态: online/offline
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	LastSeenAt time.Time `json:"last_seen_at"` // 最后在线时间
	ExpiresAt  time.Time `json:"expires_at"`  // 过期时间
	IsExpired  bool      `json:"is_expired"`  // 是否已过期
	IsOnline   bool      `json:"is_online"`   // 是否在线
	IsActive   bool      `json:"is_active"`   // 是否活跃（未过期且未软删除）
}

// SessionValidationResponse 会话验证响应
type SessionValidationResponse struct {
	SessionID string    `json:"session_id"`  // 会话ID
	IsValid   bool      `json:"is_valid"`    // 是否有效
	Message   string    `json:"message"`     // 验证消息
	ExpiresAt time.Time `json:"expires_at,omitempty"` // 过期时间
}

// RoomDetailResponse 房间详情响应
type RoomDetailResponse struct {
	Room        *model.Room `json:"room"`                  // 房间信息
	MemberCount int         `json:"member_count"`          // 当前成员数量
	CreatedBy   string      `json:"created_by"`            // 创建者显示名称
	CreatedAt   time.Time   `json:"created_at"`            // 创建时间
	IsCreator   bool        `json:"is_creator"`            // 当前用户是否为创建者
}

// JoinRoomResponse 加入房间响应
type JoinRoomResponse struct {
	Room      *model.Room `json:"room"`          // 房间信息
	SessionID string      `json:"session_id"`    // 会话ID
	Token     string      `json:"token"`         // 访问令牌
	JoinURL   string      `json:"join_url"`      // 加入链接
}

// RoomListResponse 房间列表响应
type RoomListResponse struct {
	Rooms []*model.Room `json:"rooms"` // 房间列表
	Total int64         `json:"total"` // 总数量
	Page  int           `json:"page"`  // 当前页码
	Size  int           `json:"size"`  // 每页数量
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Message string    `json:"message"` // 成功消息
	Data    interface{} `json:"data,omitempty"` // 附加数据
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`   // 错误信息
	Detail  string `json:"detail,omitempty"` // 错误详情
	Code    string `json:"code,omitempty"`   // 错误码
}

// RoomsListResponse 房间列表响应
type RoomsListResponse struct {
	Rooms  []*RoomResponse `json:"rooms"`  // 房间列表
	Total  int64           `json:"total"`  // 总数
	Page   int             `json:"page"`   // 当前页
	Size   int             `json:"size"`   // 每页数量
	HasMore bool           `json:"has_more"` // 是否有更多
}

// ==================== 工具函数 ====================

// generateSessionID 生成匿名会话ID
func generateSessionID() string {
	return "sess_" + time.Now().Format("20060102150405") + "_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

// generateDisplayName 生成随机显示名称
func generateDisplayName() string {
	prefixes := []string{
		"电影爱好者", "追剧达人", "观影客", "影迷", "视频党",
		"剧迷", "放映员", "观众", "影评人", "观影者",
	}
	
	suffixes := []string{
		"001", "007", "666", "888", "999",
		"2024", "2025", "Alpha", "Beta", "Gamma",
		"星空", "海洋", "森林", "山峰", "沙漠",
	}
	
	prefix := prefixes[time.Now().Unix()%int64(len(prefixes))]
	suffix := suffixes[time.Now().Unix()%int64(len(suffixes))]
	
	return prefix + suffix
}

// generateRoomToken 生成房间访问令牌
func generateRoomToken(roomID, sessionID string) (string, error) {
	// 这里应该实现 JWT 令牌生成逻辑
	// 为了简化 MVP，暂时使用简单的字符串拼接
	// 实际生产环境应该使用 proper JWT library
	
	// 简化实现：使用简单的字符串拼接
	token := "xiaowo_" + roomID + "_" + sessionID
	
	return token, nil
}

// generateJoinURL 生成加入链接
func generateJoinURL(roomID, token string) string {
	// 这里应该使用实际的前端域名
	baseURL := "http://localhost:3000" // 开发环境
	return baseURL + "/room/" + roomID + "?token=" + token
}

// ValidatePassword 验证房间密码
func ValidatePassword(inputPassword, roomPassword string) bool {
	if inputPassword == roomPassword {
		return true
	}
	
	// TODO: 实现密码哈希验证（生产环境需要）
	// 这里暂时使用明文比较，MVP 简化实现
	
	return false
}

// HashPassword 密码哈希（预留接口）
func HashPassword(password string) (string, error) {
	// TODO: 实现密码哈希
	// 生产环境应该使用 bcrypt 或 argon2
	
	return password, nil
}

// CheckRoomPermission 检查房间权限
func CheckRoomPermission(room *model.Room, sessionID string, action string) error {
	switch action {
	case "update":
		if room.CreatorSessionID != sessionID {
			return errors.New("只有房间创建者可以更新房间信息")
		}
	case "delete":
		if room.CreatorSessionID != sessionID {
			return errors.New("只有房间创建者可以删除房间")
		}
	}
	
	return nil
}
