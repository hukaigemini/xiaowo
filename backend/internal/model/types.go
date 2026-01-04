package model

import (
	"database/sql/driver"
	"errors"
)

// JSON is a helper type for storing JSON in SQLite TEXT fields
type JSON string

// Value implements driver.Valuer interface
func (j *JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(*j), nil
}

// Scan implements sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = ""
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*j = JSON(bytes)
	return nil
}

// Error types
var (
	ErrVersionConflict    = errors.New("version conflict: optimistic locking failed")
	ErrRoomNotFound       = errors.New("room not found")
	ErrRoomFull           = errors.New("room is full")
	ErrRoomPasswordInvalid = errors.New("invalid room password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrNotRoomCreator     = errors.New("not room creator")
	ErrInvalidMediaURL    = errors.New("invalid media URL")
	ErrInvalidPlaybackState = errors.New("invalid playback state")
	
	// Message errors
	ErrMessageNotFound    = errors.New("message not found")
	ErrMessageEmpty       = errors.New("message content is empty")
	ErrMessageTooLong     = errors.New("message content is too long")
	ErrInvalidMessageType = errors.New("invalid message type")
)

// RoomStatus represents the status of a room
type RoomStatus string

const (
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
	RoomStatusDeleted  RoomStatus = "deleted"
)

// RoomRole represents the role of a member in a room
type RoomRole string

const (
	RoleHost   RoomRole = "host"
	RoleMember RoomRole = "member"
)


