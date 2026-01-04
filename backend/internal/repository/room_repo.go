package repository

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"xiaowo/backend/internal/model"
)

// RoomRepository interface defines all room operations
type RoomRepository interface {
	Create(room *model.Room) error
	GetByID(roomID string) (*model.Room, error)
	GetByIDWithMembers(roomID string) (*model.Room, error)
	Update(roomID string, updates map[string]interface{}) (*model.Room, error)
	UpdateWithVersion(roomID string, updates map[string]interface{}, expectedVersion int) error
	Delete(roomID string, sessionID string) error
	GetRooms(filter map[string]interface{}, page, size int) ([]*model.Room, int64, error)
	SearchRooms(keyword string, filter map[string]interface{}, page, size int) ([]*model.Room, int64, error)
	GetActiveRooms(page, size int) ([]*model.Room, int64, error)
	JoinRoom(roomID, sessionID string) error
	LeaveRoom(roomID, sessionID string) error
	UpdatePlaybackState(roomID string, playbackState map[string]interface{}) error
	GetMemberCount(roomID string) (int, error)
	IsMember(roomID, sessionID string) (bool, error)
	CanAccess(roomID, sessionID string) (bool, error)
	GenerateRoomID() string
	ValidateMediaURL(url string) error
	ValidatePlaybackState(state map[string]interface{}) error
	CleanupInactiveRooms() (int64, error)
}

// RoomRepo implements RoomRepository
type RoomRepo struct {
	db *gorm.DB
}

// NewRoomRepo creates a new Room repository
func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

// Create creates a new room
func (r *RoomRepo) Create(room *model.Room) error {
	// Validate room data
	if err := r.validateRoomCreate(room); err != nil {
		return fmt.Errorf("room validation failed: %w", err)
	}

	// Generate room ID if not provided
	if room.ID == "" {
		room.ID = r.GenerateRoomID()
	}

	// Set default values
	if room.Status == "" {
		room.Status = model.RoomStatusActive
	}
	if room.PlaybackState == "" {
		room.PlaybackState = "paused"
	}
	if room.PlaybackRate == 0 {
		room.PlaybackRate = 1.0
	}
	if room.Settings == "" {
		room.Settings = model.JSON(`{
			"auto_sync":    true,
			"allow_control": true,
			"chat_enabled":  true
		}`)
	}

	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	room.LastActiveAt = time.Now()
	room.Version = 0

	if err := r.db.Create(room).Error; err != nil {
		return fmt.Errorf("failed to create room: %w", err)
	}

	return nil
}

// GetByID retrieves a room by ID
func (r *RoomRepo) GetByID(roomID string) (*model.Room, error) {
	var room model.Room
	if err := r.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	return &room, nil
}

// GetByIDWithMembers retrieves a room with its members
func (r *RoomRepo) GetByIDWithMembers(roomID string) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Members", "is_active = ?", true).Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return nil, fmt.Errorf("failed to get room with members: %w", err)
	}

	return &room, nil
}

// Update updates room information
func (r *RoomRepo) Update(roomID string, updates map[string]interface{}) (*model.Room, error) {
	// Add updated timestamp
	updates["updated_at"] = time.Now()
	updates["last_active_at"] = time.Now()

	result := r.db.Model(&model.Room{}).Where("id = ?", roomID).Updates(updates)
	
	if result.Error != nil {
		return nil, fmt.Errorf("failed to update room: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
	}

	// Return updated room
	return r.GetByID(roomID)
}

// UpdateWithVersion updates room with optimistic locking
func (r *RoomRepo) UpdateWithVersion(roomID string, updates map[string]interface{}, expectedVersion int) error {
	// Add version increment and timestamp
	updates["version"] = gorm.Expr("version + 1")
	updates["updated_at"] = time.Now()
	updates["last_active_at"] = time.Now()

	// Perform update with version check
	result := r.db.Model(&model.Room{}).
		Where("id = ? AND version = ?", roomID, expectedVersion).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update room with version: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("%w", model.ErrVersionConflict)
	}

	return nil
}

// Delete deletes a room
func (r *RoomRepo) Delete(roomID, sessionID string) error {
	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get room to check creator
	var room model.Room
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", roomID).First(&room).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return fmt.Errorf("failed to get room for deletion: %w", err)
	}

	// Check if session is the creator
	if !room.IsCreator(sessionID) {
		tx.Rollback()
		return fmt.Errorf("%w", model.ErrNotRoomCreator)
	}

	// Delete room (cascade will handle members)
	if err := tx.Delete(&model.Room{}, "id = ?", roomID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete room: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetRooms retrieves rooms with filters
func (r *RoomRepo) GetRooms(filter map[string]interface{}, page, size int) ([]*model.Room, int64, error) {
	var rooms []*model.Room
	var total int64
	offset := (page - 1) * size

	query := r.db.Model(&model.Room{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Get rooms with pagination
	if err := query.Order("last_active_at DESC").Offset(offset).Limit(size).Find(&rooms).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get rooms: %w", err)
	}

	return rooms, total, nil
}

// SearchRooms searches rooms by keyword and filters
func (r *RoomRepo) SearchRooms(keyword string, filter map[string]interface{}, page, size int) ([]*model.Room, int64, error) {
	var rooms []*model.Room
	var total int64
	offset := (page - 1) * size

	query := r.db.Model(&model.Room{})

	// Apply filters
	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// Search in name and description
	if keyword != "" {
		searchPattern := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count rooms: %w", err)
	}

	// Get rooms with pagination
	if err := query.Order("last_active_at DESC").Offset(offset).Limit(size).Find(&rooms).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search rooms: %w", err)
	}

	return rooms, total, nil
}

// GetActiveRooms retrieves active rooms with pagination
func (r *RoomRepo) GetActiveRooms(page, size int) ([]*model.Room, int64, error) {
	filter := map[string]interface{}{
		"status": model.RoomStatusActive,
	}
	return r.GetRooms(filter, page, size)
}

// JoinRoom adds a session to a room
func (r *RoomRepo) JoinRoom(roomID, sessionID string) error {
	// Check if session exists
	var session model.UserSession
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", model.ErrSessionNotFound, sessionID)
		}
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		return fmt.Errorf("%w: %s", model.ErrSessionExpired, sessionID)
	}

	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get room with lock
	var room model.Room
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", roomID).First(&room).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Check room status
	if room.Status != model.RoomStatusActive {
		tx.Rollback()
		return fmt.Errorf("room is not active: %s", roomID)
	}

	// Check if session is already a member
	var existingMember model.RoomMember
	if err := tx.Where("room_id = ? AND session_id = ? AND is_active = ?", roomID, sessionID, true).First(&existingMember).Error; err == nil {
		tx.Rollback()
		return fmt.Errorf("session already in room: %s", sessionID)
	}

	// Check room capacity
	memberCount, err := r.GetMemberCountWithDB(tx, roomID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get member count: %w", err)
	}

	if room.IsFull(memberCount) {
		tx.Rollback()
		return fmt.Errorf("%w: %d/%d", model.ErrRoomFull, memberCount, room.MaxUsers)
	}

	// Create room member
	role := model.RoleMember
	if room.IsCreator(sessionID) {
		role = model.RoleHost
	}

	member := &model.RoomMember{
		ID:        generateUUID(),
		RoomID:    roomID,
		SessionID: sessionID,
		Nickname:  session.Nickname,
		Avatar:    session.Avatar,
		Role:      role,
		JoinedAt:  time.Now(),
	}

	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create room member: %w", err)
	}

	// Update room's last active time
	if err := tx.Model(&room).Update("last_active_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update room last active: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LeaveRoom removes a session from a room
func (r *RoomRepo) LeaveRoom(roomID, sessionID string) error {
	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get room with lock
	var room model.Room
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", roomID).First(&room).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return fmt.Errorf("failed to get room: %w", err)
	}

	// Update member as inactive
	result := tx.Model(&model.RoomMember{}).
		Where("room_id = ? AND session_id = ? AND is_active = ?", roomID, sessionID, true).
		Updates(map[string]interface{}{
			"is_active": false,
			"left_at":   time.Now(),
		})

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update room member: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("session not found in room: %s", sessionID)
	}

	// Check if room should be set as inactive
	activeMemberCount, err := r.GetMemberCountWithDB(tx, roomID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to get member count: %w", err)
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if activeMemberCount == 0 {
		room.SetLastMemberLeft()
		updates["status"] = model.RoomStatusInactive
		updates["last_member_left_at"] = room.LastMemberLeftAt
	}

	// Update room
	if err := tx.Model(&room).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update room: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdatePlaybackState updates room playback state with version control
func (r *RoomRepo) UpdatePlaybackState(roomID string, playbackState map[string]interface{}) error {
	// Validate playback state
	if err := r.ValidatePlaybackState(playbackState); err != nil {
		return fmt.Errorf("invalid playback state: %w", err)
	}

	// Get expected version
	expectedVersion := 0
	if v, ok := playbackState["expected_version"]; ok {
		expectedVersion = v.(int)
	}

	// Prepare updates
	updates := map[string]interface{}{
		"playback_state": playbackState["playback_state"],
		"current_time":   playbackState["current_time"],
		"playback_rate":  playbackState["playback_rate"],
	}

	if mediaURL, ok := playbackState["media_url"]; ok {
		updates["media_url"] = mediaURL
	}
	if mediaTitle, ok := playbackState["media_title"]; ok {
		updates["media_title"] = mediaTitle
	}
	if mediaType, ok := playbackState["media_type"]; ok {
		updates["media_type"] = mediaType
	}
	if mediaDuration, ok := playbackState["media_duration"]; ok {
		updates["media_duration"] = mediaDuration
	}

	// Update with version control
	return r.UpdateWithVersion(roomID, updates, expectedVersion)
}

// GetMemberCount returns the count of active members in a room
func (r *RoomRepo) GetMemberCount(roomID string) (int, error) {
	return r.GetMemberCountWithDB(r.db, roomID)
}

// IsMember checks if a session is a member of a room
func (r *RoomRepo) IsMember(roomID, sessionID string) (bool, error) {
	var count int64
	err := r.db.Model(&model.RoomMember{}).
		Where("room_id = ? AND session_id = ? AND is_active = ?", roomID, sessionID, true).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}

	return count > 0, nil
}

// CanAccess checks if a session can access a room
func (r *RoomRepo) CanAccess(roomID, sessionID string) (bool, error) {
	// Check if session is member
	isMember, err := r.IsMember(roomID, sessionID)
	if err != nil {
		return false, fmt.Errorf("failed to check membership: %w", err)
	}
	if isMember {
		return true, nil
	}

	// Check if room is public
	var room model.Room
	if err := r.db.Where("id = ?", roomID).First(&room).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%w: %s", model.ErrRoomNotFound, roomID)
		}
		return false, fmt.Errorf("failed to get room: %w", err)
	}

	return !room.IsPrivate, nil
}

// GenerateRoomID generates a unique 6-character room ID
func (r *RoomRepo) GenerateRoomID() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 6

	for {
		// Generate random room ID
		roomID := generateRandomString(charset, length)

		// Check if it already exists
		var count int64
		r.db.Model(&model.Room{}).Where("id = ?", roomID).Count(&count)

		if count == 0 {
			return roomID
		}
	}
}

// ValidateMediaURL validates a media URL
func (r *RoomRepo) ValidateMediaURL(url string) error {
	if url == "" {
		return model.ErrInvalidMediaURL
	}

	// Basic URL validation
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return model.ErrInvalidMediaURL
	}

	// Check for valid characters
	pattern := `^https?://[^\s/$.?#].[^\s]*$`
	matched, err := regexp.MatchString(pattern, url)
	if err != nil || !matched {
		return model.ErrInvalidMediaURL
	}

	return nil
}

// ValidatePlaybackState validates playback state data
func (r *RoomRepo) ValidatePlaybackState(state map[string]interface{}) error {
	if state == nil {
		return model.ErrInvalidPlaybackState
	}

	// Validate required fields
	if state["playback_state"] == nil {
		return model.ErrInvalidPlaybackState
	}

	playbackState, ok := state["playback_state"].(string)
	if !ok {
		return model.ErrInvalidPlaybackState
	}

	// Validate playback state value
	validStates := []string{"playing", "paused", "stopped"}
	if !contains(validStates, playbackState) {
		return model.ErrInvalidPlaybackState
	}

	// Validate current time
	if currentTime, ok := state["current_time"]; ok {
		if _, ok := currentTime.(float64); !ok {
			return model.ErrInvalidPlaybackState
		}
	}

	// Validate playback rate
	if playbackRate, ok := state["playback_rate"]; ok {
		if rate, ok := playbackRate.(float64); !ok || rate < 0.25 || rate > 4.0 {
			return model.ErrInvalidPlaybackState
		}
	}

	// Validate media URL if provided
	if mediaURL, ok := state["media_url"]; ok && mediaURL != nil {
		if err := r.ValidateMediaURL(mediaURL.(string)); err != nil {
			return err
		}
	}

	return nil
}

// CleanupInactiveRooms removes inactive rooms older than specified duration
func (r *RoomRepo) CleanupInactiveRooms() (int64, error) {
	// Delete rooms that have been inactive for more than 7 days
	cutoff := time.Now().Add(-7 * 24 * time.Hour)

	result := r.db.Where("status = ? AND last_active_at < ?", model.RoomStatusInactive, cutoff).Delete(&model.Room{})
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup inactive rooms: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// Helper methods
func (r *RoomRepo) validateRoomCreate(room *model.Room) error {
	if room.Name == "" {
		return fmt.Errorf("room name is required")
	}

	if len(room.Name) > 100 {
		return fmt.Errorf("room name too long (max 100 characters)")
	}

	if room.Description != "" && len(room.Description) > 500 {
		return fmt.Errorf("room description too long (max 500 characters)")
	}

	if room.CreatorSessionID == "" {
		return fmt.Errorf("creator session ID is required")
	}

	if room.MediaURL == "" {
		return fmt.Errorf("media URL is required")
	}

	if err := r.ValidateMediaURL(room.MediaURL); err != nil {
		return fmt.Errorf("invalid media URL: %w", err)
	}

	return nil
}

func (r *RoomRepo) GetMemberCountWithDB(db *gorm.DB, roomID string) (int, error) {
	var count int64
	err := db.Model(&model.RoomMember{}).
		Where("room_id = ? AND is_active = ?", roomID, true).
		Count(&count).Error

	return int(count), err
}

func generateUUID() string {
	// Simple UUID generation for demonstration
	// In production, use github.com/google/uuid
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func generateRandomString(charset string, length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
