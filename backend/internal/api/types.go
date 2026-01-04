package types

import (
	"time"

	"xiaowo/backend/internal/model"
)

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name          string                 `json:"name" binding:"required"`
	Description   string                 `json:"description"`
	IsPrivate     *bool                  `json:"is_private,omitempty"`
	Password      string                 `json:"password,omitempty"`
	MaxUsers      *int                   `json:"max_users,omitempty"`
	MediaURL      string                 `json:"media_url" binding:"required"`
	MediaType     string                 `json:"media_type" binding:"required"`
	MediaTitle    string                 `json:"media_title"`
	MediaDuration int                    `json:"media_duration"`
	Settings      map[string]interface{} `json:"settings"`
}

// UpdateRoomRequest 更新房间请求
type UpdateRoomRequest struct {
	Name       *string                 `json:"name,omitempty"`
	IsPrivate  *bool                   `json:"is_private,omitempty"`
	Password   *string                 `json:"password,omitempty"`
	MaxUsers   *int                    `json:"max_users,omitempty"`
	Settings   map[string]interface{}  `json:"settings"`
}

// JoinRoomRequest 加入房间请求
type JoinRoomRequest struct {
	DisplayName string `json:"display_name"`
	Password    string `json:"password,omitempty"`
}

// PlaybackControlRequest 播放控制请求
type PlaybackControlRequest struct {
	Action     string                 `json:"action" binding:"required"` // play, pause, seek, rate
	Position   *float64               `json:"position,omitempty"`         // 播放位置(秒)
	Rate       *float64               `json:"rate,omitempty"`             // 播放速度
	Data       map[string]interface{} `json:"data,omitempty"`             // 其他参数
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	Nickname string `json:"nickname,omitempty"`
}

// UpdateSessionRequest 更新会话请求
type UpdateSessionRequest struct {
	Nickname *string `json:"nickname,omitempty"`
	Avatar   *string `json:"avatar,omitempty"`
}

// 响应结构

// RoomResponse 房间响应
type RoomResponse struct {
	ID               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	CreatorSessionID string                 `json:"creator_session_id"`
	IsPrivate        bool                   `json:"is_private"`
	MaxUsers         int                    `json:"max_users"`
	CurrentUsers     int                    `json:"current_users"`
	Status           model.RoomStatus       `json:"status"`
	MediaURL         string                 `json:"media_url"`
	MediaType        string                 `json:"media_type"`
	MediaTitle       string                 `json:"media_title"`
	MediaDuration    int                    `json:"media_duration"`
	PlaybackState    string                 `json:"playback_state"`
	CurrentTime      float64                `json:"current_time"`
	PlaybackRate     float64                `json:"playback_rate"`
	Settings         map[string]interface{} `json:"settings"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

// MemberResponse 房间成员响应
type MemberResponse struct {
	SessionID     string    `json:"session_id"`
	DisplayName   string    `json:"display_name"`
	Avatar        string    `json:"avatar"`
	JoinTime      time.Time `json:"join_time"`
	LastActiveAt  time.Time `json:"last_active_at"`
	Role          string    `json:"role"`
}

// SessionResponse 会话响应
type SessionResponse struct {
	ID          string                 `json:"id"`
	Nickname    string                 `json:"nickname"`
	Avatar      string                 `json:"avatar"`
	RoomID      *string                `json:"room_id"`
	Status      model.UserSessionStatus `json:"status"`
	CreatedAt   time.Time              `json:"created_at"`
	LastSeenAt  time.Time              `json:"last_seen_at"`
	ExpiresAt   time.Time              `json:"expires_at"`
}

// MessageResponse 消息响应
type MessageResponse struct {
	ID          string               `json:"id"`
	SessionID   string               `json:"session_id"`
	RoomID      string               `json:"room_id"`
	Type        model.MessageType    `json:"type"`
	Content     string               `json:"content"`
	Data        map[string]interface{} `json:"data,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
}

// PlaybackStateResponse 播放状态响应
type PlaybackStateResponse struct {
	RoomID       string    `json:"room_id"`
	PlaybackState string   `json:"playback_state"`
	CurrentTime  float64   `json:"current_time"`
	PlaybackRate float64   `json:"playback_rate"`
	Version      int       `json:"version"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// PlaybackControlResponse 播放控制响应
type PlaybackControlResponse struct {
	Success    bool                     `json:"success"`
	Message    string                   `json:"message,omitempty"`
	RoomID     string                   `json:"room_id"`
	Action     string                   `json:"action"`
	Position   *float64                 `json:"position,omitempty"`
	Rate       *float64                 `json:"rate,omitempty"`
	Data       map[string]interface{}   `json:"data,omitempty"`
	Version    int                      `json:"version"`
	UpdatedAt  time.Time                `json:"updated_at"`
}

// APIResponse 通用API响应
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError API错误响应
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   APIError  `json:"error"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	Success   bool          `json:"success"`
	Data      interface{}   `json:"data"`
	Total     int           `json:"total"`
	Page      int           `json:"page"`
	PageSize  int           `json:"page_size"`
	TotalPage int           `json:"total_page"`
}