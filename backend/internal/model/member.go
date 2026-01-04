package model

import (
	"time"
)

// RoomMember represents a user in a room
type RoomMember struct {
	ID        string    `gorm:"primaryKey;type:text" json:"id"`
	RoomID    string    `gorm:"type:text;not null;index:idx_room_session,unique" json:"room_id"`
	SessionID string    `gorm:"type:text;not null;index:idx_room_session,unique" json:"session_id"`
	Role      RoomRole  `gorm:"type:text;default:'member'" json:"role"`
	Nickname  string    `gorm:"type:text" json:"nickname"`
	Avatar    string    `gorm:"type:text" json:"avatar"`
	IsMuted   bool      `gorm:"type:integer;default:0" json:"is_muted"`
	JoinedAt  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"joined_at"`
	LastSeen  time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"last_seen"`
	CreatedAt time.Time `gorm:"type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName overrides the table name
func (RoomMember) TableName() string {
	return "room_members"
}


