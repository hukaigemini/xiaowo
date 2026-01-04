package model

import (
	"time"

	"gorm.io/gorm"
)

// UserSessionStatus represents the status of a user session
type UserSessionStatus string

const (
	StatusOnline  UserSessionStatus = "online"  // 在线
	StatusOffline UserSessionStatus = "offline" // 离线
)

// UserSession represents a temporary user session without independent user system
type UserSession struct {
	ID        string          `gorm:"primaryKey;type:text" json:"id"` // Session unique ID (UUID)
	Nickname  string          `gorm:"type:text;not null" json:"nickname"`
	Avatar    string          `gorm:"type:text;not null" json:"avatar"`
	RoomID    *string         `gorm:"type:text;index" json:"room_id"` // Current room ID (nullable)
	Status    UserSessionStatus `gorm:"type:text;default:'online'" json:"status"`
	CreatedAt time.Time       `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
	LastSeenAt time.Time      `gorm:"type:datetime;default:CURRENT_TIMESTAMP;index" json:"last_seen_at"`
	ExpiresAt time.Time       `gorm:"type:datetime;default:(datetime('now', '+7 days'))" json:"expires_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName overrides the table name
func (UserSession) TableName() string {
	return "user_sessions"
}

// IsExpired checks if the session has expired
func (s *UserSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsOnline checks if the session is online
func (s *UserSession) IsOnline() bool {
	return s.Status == StatusOnline
}

// UpdateLastSeen updates the last seen timestamp and sets status to online
func (s *UserSession) UpdateLastSeen() {
	s.LastSeenAt = time.Now()
	s.Status = StatusOnline
}

// SetOffline sets the session status to offline
func (s *UserSession) SetOffline() {
	s.Status = StatusOffline
}

// IsActive checks if the session is active (not expired, not deleted and in a room)
func (s *UserSession) IsActive() bool {
	return !s.IsExpired() && s.RoomID != nil && s.Status == StatusOnline
}