package model

import (
	"time"
)

// Room represents a viewing room
type Room struct {
	ID                 string     `gorm:"primaryKey;type:text" json:"id"`                       // 房间ID (6位数字+字母混合房间号)
	Name               string     `gorm:"type:text;not null" json:"name"`                       // 房间名称
	Description        string     `gorm:"type:text" json:"description"`                        // 房间描述
	CreatorSessionID   string     `gorm:"type:text;not null" json:"creator_session_id"`        // 创建者会话ID
	IsPrivate          bool       `gorm:"type:integer;default:0" json:"is_private"`            // 是否私密房间
	Password           string     `gorm:"column:room_password;type:text" json:"-"`             // 房间密码 (如有)
	MaxUsers           int        `gorm:"type:integer;default:7" json:"max_users"`             // 最大用户数 (固定为7)
	Status             RoomStatus `gorm:"type:text;default:'active';index" json:"status"`      // 房间状态: active/inactive
	MediaURL           string     `gorm:"type:text;not null" json:"media_url"`                 // 媒体资源URL
	MediaType          string     `gorm:"type:text;default:'video'" json:"media_type"`         // 媒体类型: video/audio/stream
	MediaTitle         string     `gorm:"type:text" json:"media_title"`                        // 媒体标题
	MediaDuration      float64    `gorm:"type:real;default:0" json:"media_duration"`           // 媒体总时长 (秒)
	PlaybackState      string     `gorm:"type:text;default:'paused'" json:"playback_state"`    // 播放状态: playing/paused/stopped
	CurrentTime        float64    `gorm:"type:real;default:0" json:"current_time"`             // 当前播放时间 (秒)
	PlaybackRate       float64    `gorm:"type:real;default:1.0" json:"playback_rate"`          // 播放速率 (1.0=正常, 1.5=1.5倍速)
	Settings           JSON       `gorm:"type:text;default:'{}'" json:"settings"`              // 房间设置 (JSON格式)
	Version            int        `gorm:"type:integer;default:0" json:"version"`               // 乐观锁版本号
	LastActiveAt       time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index" json:"last_active_at"`      // 最后活跃时间
	LastMemberLeftAt   *time.Time `gorm:"type:datetime" json:"last_member_left_at"`            // 最后一位成员离开时间
	CreatedAt          time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`                // 创建时间
	UpdatedAt          time.Time  `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`                // 更新时间

	// Relations
	Members []RoomMember `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
}

// IsFull checks if the room has reached maximum user capacity
func (r *Room) IsFull(currentMemberCount int) bool {
	return currentMemberCount >= r.MaxUsers
}

// IsCreator checks if the given session ID is the room creator
func (r *Room) IsCreator(sessionID string) bool {
	return r.CreatorSessionID == sessionID
}

//State returns the current GetPlayback playback state as JSON
func (r *Room) GetPlaybackState() map[string]interface{} {
	return map[string]interface{}{
		"playback_state": r.PlaybackState,
		"current_time":   r.CurrentTime,
		"playback_rate":  r.PlaybackRate,
		"media_url":      r.MediaURL,
		"media_title":    r.MediaTitle,
		"version":        r.Version,
		"updated_at":     r.UpdatedAt,
	}
}

// UpdatePlaybackState updates the playback state with version checking
func (r *Room) UpdatePlaybackState(newState map[string]interface{}) error {
	// Version check for optimistic locking
	if expectedVersion, ok := newState["expected_version"]; ok {
		if r.Version != expectedVersion.(int) {
			return ErrVersionConflict
		}
	}

	// Update fields
	if state, ok := newState["playback_state"]; ok {
		r.PlaybackState = state.(string)
	}
	if currentTime, ok := newState["current_time"]; ok {
		r.CurrentTime = currentTime.(float64)
	}
	if rate, ok := newState["playback_rate"]; ok {
		r.PlaybackRate = rate.(float64)
	}
	if mediaURL, ok := newState["media_url"]; ok {
		r.MediaURL = mediaURL.(string)
	}
	if mediaTitle, ok := newState["media_title"]; ok {
		r.MediaTitle = mediaTitle.(string)
	}

	// Increment version
	r.Version++

	// Update timestamp
	r.UpdatedAt = time.Now()
	r.LastActiveAt = time.Now()

	return nil
}

// SetLastMemberLeft updates the last member left timestamp
func (r *Room) SetLastMemberLeft() {
	now := time.Now()
	r.LastMemberLeftAt = &now
	r.UpdatedAt = now
}

// TableName overrides the table name
func (Room) TableName() string {
	return "rooms"
}
