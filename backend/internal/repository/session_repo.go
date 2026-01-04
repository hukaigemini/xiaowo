package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"xiaowo/backend/internal/model"
)

// SessionRepository interface defines all session operations
type SessionRepository interface {
	Create(nickname string) (*model.UserSession, error)
	GetByID(sessionID string) (*model.UserSession, error)
	Update(sessionID string, updates map[string]interface{}) (*model.UserSession, error)
	UpdateLastSeen(sessionID string) error
	UpdateStatus(sessionID, status string) error
	JoinRoom(sessionID, roomID string) error
	LeaveRoom(sessionID string) error
	GetActiveSessions() ([]*model.UserSession, error)
	GetExpiredSessions() ([]*model.UserSession, error)
	GetByStatus(status string) ([]*model.UserSession, error)
	GetOnlineSessions() ([]*model.UserSession, error)
	Delete(sessionID string) error
	SoftDelete(sessionID string) error
	CleanupExpired() (int64, error)
	GenerateNickname() string
	GenerateAvatar() string
}

// SessionRepo implements SessionRepository
type SessionRepo struct {
	db *gorm.DB
}

// NewSessionRepo creates a new Session repository
func NewSessionRepo(db *gorm.DB) *SessionRepo {
	return &SessionRepo{db: db}
}

// Create creates a new user session
func (r *SessionRepo) Create(nickname string) (*model.UserSession, error) {
	// Validate nickname
	if nickname == "" {
		nickname = r.GenerateNickname()
	}

	// Generate avatar if not provided
	avatar := r.GenerateAvatar()

	// Create new session
	session := &model.UserSession{
		ID:        uuid.New().String(),
		Nickname:  nickname,
		Avatar:    avatar,
		CreatedAt: time.Now(),
		LastSeenAt: time.Now(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days from now
	}

	if err := r.db.Create(session).Error; err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetByID retrieves a session by ID
func (r *SessionRepo) GetByID(sessionID string) (*model.UserSession, error) {
	var session model.UserSession
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		return nil, fmt.Errorf("session expired: %s", sessionID)
	}

	return &session, nil
}

// Update updates session information
func (r *SessionRepo) Update(sessionID string, updates map[string]interface{}) (*model.UserSession, error) {
	var session model.UserSession
	
	// Start transaction for atomic update
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Check if session exists and is not expired
	if err := tx.Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("session not found: %s", sessionID)
		}
		return nil, fmt.Errorf("failed to get session for update: %w", err)
	}

	// Check expiration
	if session.IsExpired() {
		tx.Rollback()
		return nil, fmt.Errorf("session expired: %s", sessionID)
	}

	// Validate nickname if updating
	if nickname, ok := updates["nickname"]; ok && nickname != "" {
		if len(nickname.(string)) > 50 {
			tx.Rollback()
			return nil, fmt.Errorf("nickname too long (max 50 characters)")
		}
	}

	// Add updated timestamp
	updates["last_seen_at"] = time.Now()

	// Perform update
	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// Refresh session data
	if err := r.db.Where("id = ?", sessionID).First(&session).Error; err != nil {
		return nil, fmt.Errorf("failed to refresh session: %w", err)
	}

	return &session, nil
}

// UpdateLastSeen updates the last seen timestamp
func (r *SessionRepo) UpdateLastSeen(sessionID string) error {
	result := r.db.Model(&model.UserSession{}).
		Where("id = ?", sessionID).
		Update("last_seen_at", time.Now())

	if result.Error != nil {
		return fmt.Errorf("failed to update last seen: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// UpdateStatus updates the session status
func (r *SessionRepo) UpdateStatus(sessionID string, status string) error {
	// Validate status
	if status != string(model.StatusOnline) && status != string(model.StatusOffline) {
		return fmt.Errorf("invalid status: %v (must be 'online' or 'offline')", status)
	}

	// Start transaction for atomic update
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get session with lock
	var session model.UserSession
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("session not found: %s", sessionID)
		}
		return fmt.Errorf("failed to get session for status update: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		tx.Rollback()
		return fmt.Errorf("session expired: %s", sessionID)
	}

	// Update status and last_seen_at
	updates := map[string]interface{}{
		"status":       status,
		"last_seen_at": time.Now(),
	}

	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// JoinRoom adds a session to a room
func (r *SessionRepo) JoinRoom(sessionID, roomID string) error {
	var session model.UserSession
	
	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get session with lock
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("session not found: %s", sessionID)
		}
		return fmt.Errorf("failed to get session for join room: %w", err)
	}

	// Check if session is expired
	if session.IsExpired() {
		tx.Rollback()
		return fmt.Errorf("session expired: %s", sessionID)
	}

	// Update room_id and last_seen_at
	updates := map[string]interface{}{
		"room_id":     roomID,
		"last_seen_at": time.Now(),
	}

	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to join room: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// LeaveRoom removes a session from a room
func (r *SessionRepo) LeaveRoom(sessionID string) error {
	var session model.UserSession
	
	// Start transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Get session with lock
	if err := tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", sessionID).First(&session).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("session not found: %s", sessionID)
		}
		return fmt.Errorf("failed to get session for leave room: %w", err)
	}

	// Update room_id to NULL and last_seen_at
	updates := map[string]interface{}{
		"room_id":     nil,
		"last_seen_at": time.Now(),
	}

	if err := tx.Model(&session).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to leave room: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetActiveSessions retrieves all active (not expired) sessions
func (r *SessionRepo) GetActiveSessions() ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	
	err := r.db.Where("expires_at > ?", time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	return sessions, nil
}

// GetExpiredSessions retrieves all expired sessions
func (r *SessionRepo) GetExpiredSessions() ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	
	err := r.db.Where("expires_at <= ?", time.Now()).Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get expired sessions: %w", err)
	}

	return sessions, nil
}

// GetByStatus retrieves sessions by status
func (r *SessionRepo) GetByStatus(status string) ([]*model.UserSession, error) {
	var sessions []*model.UserSession
	
	// Validate status
	if status != string(model.StatusOnline) && status != string(model.StatusOffline) {
		return nil, fmt.Errorf("invalid status: %s (must be 'online' or 'offline')", status)
	}
	
	err := r.db.Where("status = ?", status).Find(&sessions).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions by status: %w", err)
	}

	return sessions, nil
}

// GetOnlineSessions retrieves all online sessions
func (r *SessionRepo) GetOnlineSessions() ([]*model.UserSession, error) {
	return r.GetByStatus(string(model.StatusOnline))
}

// Delete removes a session
func (r *SessionRepo) Delete(sessionID string) error {
	result := r.db.Delete(&model.UserSession{}, "id = ?", sessionID)
	
	if result.Error != nil {
		return fmt.Errorf("failed to delete session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// SoftDelete performs a soft delete on a session
func (r *SessionRepo) SoftDelete(sessionID string) error {
	result := r.db.Model(&model.UserSession{}).
		Where("id = ?", sessionID).
		Updates(map[string]interface{}{
			"deleted_at": time.Now(),
			"status":     model.StatusOffline,
		})
	
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete session: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	return nil
}

// CleanupExpired removes all expired sessions and returns the count of removed sessions
func (r *SessionRepo) CleanupExpired() (int64, error) {
	result := r.db.Where("expires_at <= ?", time.Now()).Delete(&model.UserSession{})
	
	if result.Error != nil {
		return 0, fmt.Errorf("failed to cleanup expired sessions: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// GenerateNickname generates a random fun nickname
func (r *SessionRepo) GenerateNickname() string {
	adjectives := []string{
		"快乐", "聪明", "勇敢", "温柔", "活泼", "安静", "幽默", "善良",
		"可爱", "酷炫", "热情", "冷静", "机智", "幽默", "开朗", "贴心",
		"细心", "果断", "体贴", "浪漫", "梦幻", "清新", "阳光", "星辰",
		"海洋", "森林", "雪山", "花朵", "蝴蝶", "彩虹", "星星", "月亮",
	}

	animals := []string{
		"小熊猫", "小猫咪", "小狗子", "小兔子", "小松鼠", "小鸟",
		"小海豚", "小企鹅", "小熊猫", "小狐狸", "小鹿", "小熊",
		"小鲸鱼", "小海马", "小螃蟹", "小章鱼", "小蜜蜂", "小蝴蝶",
		"小蜻蜓", "小蚂蚁", "小蜗牛", "小青蛙", "小乌龟", "小鸽子",
	}

	// Simple hash-based selection for consistency
	now := time.Now()
	seed := now.UnixNano()
	adjIndex := int(seed) % len(adjectives)
	animalIndex := int(seed/1000) % len(animals)

	return adjectives[adjIndex] + animals[animalIndex]
}

// GenerateAvatar generates a random avatar URL using DiceBear
func (r *SessionRepo) GenerateAvatar() string {
	now := time.Now()
	seed := now.UnixNano()
	avatarSeed := fmt.Sprintf("xiaowo_%d", seed)
	
	return fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s&backgroundColor=ffd5dc,5199e4,64c8ff,fa709a,fee140&radius=50", avatarSeed)
}