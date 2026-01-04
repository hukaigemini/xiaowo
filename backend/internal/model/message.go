package model

import (
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	MessageTypeChat         MessageType = "chat"
	MessageTypeSystem       MessageType = "system"
	MessageTypeNotification MessageType = "notification"
)

// Message represents a room message
type Message struct {
	ID          string      `gorm:"primaryKey;type:text" json:"id"`           // 消息唯一ID (UUID)
	RoomID      string      `gorm:"type:text;not null;index" json:"room_id"` // 房间ID
	SessionID   string      `gorm:"type:text;not null;index" json:"session_id"` // 发送者会话ID
	MessageType MessageType `gorm:"type:text;default:'chat';index" json:"message_type"` // 消息类型: chat/system/notification
	Content     string      `gorm:"type:text;not null" json:"content"`       // 消息内容
	Metadata    JSON        `gorm:"type:text;default:'{}'" json:"metadata"`  // 消息元数据 (JSON格式)
	CreatedAt   time.Time   `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index" json:"created_at"` // 创建时间

	// Relations
	Room    *Room        `gorm:"foreignKey:RoomID;constraint:OnDelete:CASCADE" json:"room,omitempty"`
	Session *UserSession `gorm:"foreignKey:SessionID;constraint:OnDelete:CASCADE;references:ID" json:"session,omitempty"`
}

// TableName overrides the table name
func (Message) TableName() string {
	return "room_messages"
}

// IsSystem checks if the message is a system message
func (m *Message) IsSystem() bool {
	return m.MessageType == MessageTypeSystem
}

// IsChat checks if the message is a chat message
func (m *Message) IsChat() bool {
	return m.MessageType == MessageTypeChat
}

// IsNotification checks if the message is a notification
func (m *Message) IsNotification() bool {
	return m.MessageType == MessageTypeNotification
}

// GetSenderNickname returns the sender's nickname from the session if available
func (m *Message) GetSenderNickname() string {
	if m.Session != nil {
		return m.Session.Nickname
	}
	return "未知用户"
}

// GetSenderAvatar returns the sender's avatar from the session if available
func (m *Message) GetSenderAvatar() string {
	if m.Session != nil {
		return m.Session.Avatar
	}
	return ""
}

// ValidateContent validates the message content
func (m *Message) ValidateContent() error {
	if len(m.Content) == 0 {
		return ErrMessageEmpty
	}
	if len(m.Content) > 2000 {
		return ErrMessageTooLong
	}
	return nil
}

// SanitizeContent removes potentially harmful content
func (m *Message) SanitizeContent() {
	// Basic sanitization - in production, you might want more sophisticated filtering
	// This is a placeholder for content moderation logic
}

// TruncateContent truncates content to a maximum length
func (m *Message) TruncateContent(maxLen int) {
	if len(m.Content) > maxLen {
		m.Content = m.Content[:maxLen] + "..."
	}
}

// GetFormattedTime returns a formatted time string
func (m *Message) GetFormattedTime() string {
	return m.CreatedAt.Format("15:04:05")
}

// GetTimestamp returns Unix timestamp
func (m *Message) GetTimestamp() int64 {
	return m.CreatedAt.Unix()
}