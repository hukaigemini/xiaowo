package repository

import (
	"time"

	"gorm.io/gorm"
	"xiaowo/backend/internal/model"
)

// MessageRepository interface defines all message operations
type MessageRepository interface {
	// Basic CRUD operations
	Create(message *model.Message) error
	GetByID(messageID string) (*model.Message, error)
	GetByIDWithRelations(messageID string) (*model.Message, error)
	Update(messageID string, updates map[string]interface{}) (*model.Message, error)
	Delete(messageID string) error

	// Query operations
	GetMessagesByRoom(roomID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error)
	GetMessagesBySession(sessionID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error)
	GetRecentMessages(roomID string, limit int) ([]*model.Message, error)
	GetSystemMessages(roomID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error)

	// Statistics
	GetMessageCountByRoom(roomID string) (int64, error)
	GetMessageCountBySession(sessionID string) (int64, error)
	GetMessageStats(roomID string, startTime, endTime time.Time) (map[string]interface{}, error)

	// Utility operations
	SearchMessages(roomID string, keyword string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error)
	CleanupOldMessages(roomID string, daysOld int) (int64, error)
	ValidateMessageType(messageType model.MessageType) error
	GenerateMessageID() string
}

// messageRepository implements MessageRepository
type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Create creates a new message
func (r *messageRepository) Create(message *model.Message) error {
	// Validate message content
	if err := message.ValidateContent(); err != nil {
		return err
	}

	// Validate message type
	if err := r.ValidateMessageType(message.MessageType); err != nil {
		return err
	}

	// Sanitize content
	message.SanitizeContent()

	// Set message ID if not provided
	if message.ID == "" {
		message.ID = r.GenerateMessageID()
	}

	// Begin transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Create message
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update room's last_active_at timestamp
	if err := tx.Model(&model.Room{}).Where("id = ?", message.RoomID).Update("last_active_at", time.Now()).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	return tx.Commit().Error
}

// GetByID gets a message by ID
func (r *messageRepository) GetByID(messageID string) (*model.Message, error) {
	var message model.Message
	if err := r.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// GetByIDWithRelations gets a message by ID with related data
func (r *messageRepository) GetByIDWithRelations(messageID string) (*model.Message, error) {
	var message model.Message
	if err := r.db.Preload("Session").Preload("Room").Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

// Update updates a message
func (r *messageRepository) Update(messageID string, updates map[string]interface{}) (*model.Message, error) {
	var message model.Message
	if err := r.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}

	// Validate updates if content is being updated
	if content, ok := updates["content"]; ok {
		message.Content = content.(string)
		if err := message.ValidateContent(); err != nil {
			return nil, err
		}
	}

	// Sanitize content if present
	if content, ok := updates["content"]; ok {
		message.Content = content.(string)
		message.SanitizeContent()
		updates["content"] = message.Content
	}

	if err := r.db.Model(&message).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &message, nil
}

// Delete deletes a message
func (r *messageRepository) Delete(messageID string) error {
	result := r.db.Delete(&model.Message{}, "id = ?", messageID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return model.ErrMessageNotFound
	}
	return nil
}

// GetMessagesByRoom gets messages for a specific room with pagination
func (r *messageRepository) GetMessagesByRoom(roomID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("room_id = ?", roomID)

	// Apply filters
	if messageType, ok := filter["message_type"]; ok {
		query = query.Where("message_type = ?", messageType)
	}
	if startTime, ok := filter["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filter["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with pagination
	offset := (page - 1) * size
	if err := query.Preload("Session").Offset(offset).Limit(size).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, total, nil
}

// GetMessagesBySession gets messages sent by a specific session
func (r *messageRepository) GetMessagesBySession(sessionID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("session_id = ?", sessionID)

	// Apply filters
	if messageType, ok := filter["message_type"]; ok {
		query = query.Where("message_type = ?", messageType)
	}
	if roomID, ok := filter["room_id"]; ok {
		query = query.Where("room_id = ?", roomID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with pagination
	offset := (page - 1) * size
	if err := query.Preload("Room").Offset(offset).Limit(size).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, total, nil
}

// GetRecentMessages gets the most recent messages for a room
func (r *messageRepository) GetRecentMessages(roomID string, limit int) ([]*model.Message, error) {
	var messages []*model.Message
	if err := r.db.Preload("Session").Where("room_id = ?", roomID).Order("created_at DESC").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetSystemMessages gets system messages for a room
func (r *messageRepository) GetSystemMessages(roomID string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("room_id = ? AND message_type = ?", roomID, model.MessageTypeSystem)

	// Apply additional filters
	if startTime, ok := filter["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filter["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with pagination
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, total, nil
}

// GetMessageCountByRoom gets the total message count for a room
func (r *messageRepository) GetMessageCountByRoom(roomID string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Message{}).Where("room_id = ?", roomID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetMessageCountBySession gets the total message count sent by a session
func (r *messageRepository) GetMessageCountBySession(sessionID string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Message{}).Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetMessageStats gets message statistics for a room
func (r *messageRepository) GetMessageStats(roomID string, startTime, endTime time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total messages
	var totalMessages int64
	if err := r.db.Model(&model.Message{}).Where("room_id = ? AND created_at BETWEEN ? AND ?", roomID, startTime, endTime).Count(&totalMessages).Error; err != nil {
		return nil, err
	}
	stats["total_messages"] = totalMessages

	// Message count by type
	var messageTypeStats []struct {
		MessageType string `gorm:"column:message_type"`
		Count       int64  `gorm:"column:count"`
	}
	if err := r.db.Model(&model.Message{}).Where("room_id = ? AND created_at BETWEEN ? AND ?", roomID, startTime, endTime).Group("message_type").Pluck("message_type, COUNT(*) as count", &messageTypeStats).Error; err != nil {
		return nil, err
	}
	stats["messages_by_type"] = messageTypeStats

	// Messages per user
	var messagesPerUser []struct {
		SessionID string `gorm:"column:session_id"`
		Nickname  string `gorm:"column:nickname"`
		Count     int64  `gorm:"column:count"`
	}
	if err := r.db.Table("room_messages").Select("room_messages.session_id, user_sessions.nickname, COUNT(*) as count").Joins("JOIN user_sessions ON room_messages.session_id = user_sessions.id").Where("room_messages.room_id = ? AND room_messages.created_at BETWEEN ? AND ?", roomID, startTime, endTime).Group("room_messages.session_id, user_sessions.nickname").Pluck("session_id, nickname, count", &messagesPerUser).Error; err != nil {
		return nil, err
	}
	stats["messages_per_user"] = messagesPerUser

	return stats, nil
}

// SearchMessages searches messages by keyword
func (r *messageRepository) SearchMessages(roomID string, keyword string, filter map[string]interface{}, page, size int) ([]*model.Message, int64, error) {
	var messages []*model.Message
	var total int64

	query := r.db.Model(&model.Message{}).Where("room_id = ? AND content LIKE ?", roomID, "%"+keyword+"%")

	// Apply additional filters
	if messageType, ok := filter["message_type"]; ok {
		query = query.Where("message_type = ?", messageType)
	}
	if startTime, ok := filter["start_time"]; ok {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime, ok := filter["end_time"]; ok {
		query = query.Where("created_at <= ?", endTime)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get messages with pagination
	offset := (page - 1) * size
	if err := query.Preload("Session").Offset(offset).Limit(size).Order("created_at DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// Reverse the slice to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, total, nil
}

// CleanupOldMessages cleans up old messages for a room
func (r *messageRepository) CleanupOldMessages(roomID string, daysOld int) (int64, error) {
	cutoffDate := time.Now().AddDate(0, 0, -daysOld)
	result := r.db.Where("room_id = ? AND created_at < ?", roomID, cutoffDate).Delete(&model.Message{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

// ValidateMessageType validates the message type
func (r *messageRepository) ValidateMessageType(messageType model.MessageType) error {
	switch messageType {
	case model.MessageTypeChat, model.MessageTypeSystem, model.MessageTypeNotification:
		return nil
	default:
		return model.ErrInvalidMessageType
	}
}

// GenerateMessageID generates a unique message ID
func (r *messageRepository) GenerateMessageID() string {
	return generateUUID()
}