package service

import (
	"fmt"
	"time"

	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/repository"
)

// SessionService 会话业务逻辑服务
type SessionService struct {
	sessionRepo repository.SessionRepository
}

// NewSessionService 创建会话服务
func NewSessionService(sessionRepo repository.SessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

// CreateSession 创建新会话
func (s *SessionService) CreateSession(nickname string) (*model.UserSession, error) {
	if nickname == "" {
		nickname = s.sessionRepo.GenerateNickname()
	}

	session, err := s.sessionRepo.Create(nickname)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession 获取会话信息
func (s *SessionService) GetSession(sessionID string) (*model.UserSession, error) {
	return s.sessionRepo.GetByID(sessionID)
}

// UpdateSession 更新会话信息
func (s *SessionService) UpdateSession(sessionID string, updates map[string]interface{}) (*model.UserSession, error) {
	return s.sessionRepo.Update(sessionID, updates)
}

// UpdateLastSeen 更新最后在线时间
func (s *SessionService) UpdateLastSeen(sessionID string) error {
	return s.sessionRepo.UpdateLastSeen(sessionID)
}

// UpdateStatus 更新会话状态
func (s *SessionService) UpdateStatus(sessionID, status string) error {
	return s.sessionRepo.UpdateStatus(sessionID, status)
}

// Heartbeat 心跳处理，更新最后在线时间和状态
func (s *SessionService) Heartbeat(sessionID string) error {
	var updates = map[string]interface{}{
		"last_seen_at": time.Now(),
		"status":       model.StatusOnline,
	}
	
	_, err := s.sessionRepo.Update(sessionID, updates)
	return err
}

// ValidateSession 验证会话是否有效
func (s *SessionService) ValidateSession(sessionID string) (bool, error) {
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return false, err
	}
	
	return session.IsActive(), nil
}

// JoinRoom 加入房间
func (s *SessionService) JoinRoom(sessionID, roomID string) error {
	return s.sessionRepo.JoinRoom(sessionID, roomID)
}

// LeaveRoom 离开房间
func (s *SessionService) LeaveRoom(sessionID string) error {
	return s.sessionRepo.LeaveRoom(sessionID)
}

// GetActiveSessions 获取活跃会话
func (s *SessionService) GetActiveSessions() ([]*model.UserSession, error) {
	return s.sessionRepo.GetActiveSessions()
}

// GetExpiredSessions 获取过期会话
func (s *SessionService) GetExpiredSessions() ([]*model.UserSession, error) {
	return s.sessionRepo.GetExpiredSessions()
}

// GetSessionsByStatus 根据状态获取会话
func (s *SessionService) GetSessionsByStatus(status string) ([]*model.UserSession, error) {
	return s.sessionRepo.GetByStatus(status)
}

// GetOnlineSessions 获取在线会话
func (s *SessionService) GetOnlineSessions() ([]*model.UserSession, error) {
	return s.sessionRepo.GetOnlineSessions()
}

// DeleteSession 删除会话
func (s *SessionService) DeleteSession(sessionID string) error {
	return s.sessionRepo.Delete(sessionID)
}

// SoftDeleteSession 软删除会话
func (s *SessionService) SoftDeleteSession(sessionID string) error {
	return s.sessionRepo.SoftDelete(sessionID)
}

// CleanupExpired 清理过期会话
func (s *SessionService) CleanupExpired() (int64, error) {
	return s.sessionRepo.CleanupExpired()
}

// IsSessionValid 检查会话是否有效
func (s *SessionService) IsSessionValid(sessionID string) bool {
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return false
	}
	return !session.IsExpired()
}

// GenerateSession 创建并生成会话信息
func (s *SessionService) GenerateSession(nickname string) (sessionID string, displayName string, avatar string) {
	session, err := s.CreateSession(nickname)
	if err != nil {
		// 如果创建失败，生成临时会话信息
		sessionID = generateSessionID()
		displayName = s.sessionRepo.GenerateNickname()
		avatar = s.sessionRepo.GenerateAvatar()
		return
	}

	return session.ID, session.Nickname, session.Avatar
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	return "sess_" + generateUUID()[:8]
}

// generateUUID 生成UUID（复制自repository）
func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}