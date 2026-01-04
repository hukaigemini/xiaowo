package service

import (
	"fmt"
	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/repository"
	"time"
)

// CreateRoomRequest 创建房间请求
type CreateRoomRequest struct {
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	IsPrivate     *bool                  `json:"is_private,omitempty"`
	Password      string                 `json:"password,omitempty"`
	MaxUsers      *int                   `json:"max_users,omitempty"`
	MediaURL      string                 `json:"media_url"`
	MediaType     string                 `json:"media_type"`
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

// RoomService 房间业务逻辑服务
type RoomService struct {
	roomRepo   repository.RoomRepository
	memberRepo repository.RoomMemberRepository
}

// NewRoomService 创建房间服务
func NewRoomService(roomRepo repository.RoomRepository, memberRepo repository.RoomMemberRepository) *RoomService {
	return &RoomService{
		roomRepo:   roomRepo,
		memberRepo: memberRepo,
	}
}

// CreateRoom 创建房间
func (s *RoomService) CreateRoom(req *CreateRoomRequest, creatorSessionID string) (*model.Room, error) {
	// 验证媒体URL
	if err := s.roomRepo.ValidateMediaURL(req.MediaURL); err != nil {
		return nil, err
	}

	// 设置默认值
	isPrivate := false
	if req.IsPrivate != nil {
		isPrivate = *req.IsPrivate
	}
	
	maxUsers := 10
	if req.MaxUsers != nil {
		maxUsers = *req.MaxUsers
	}
	
	mediaDuration := float64(req.MediaDuration)

	room := &model.Room{
		ID:               s.roomRepo.GenerateRoomID(),
		Name:             req.Name,
		Description:      req.Description,
		CreatorSessionID: creatorSessionID,
		IsPrivate:        isPrivate,
		Password:         req.Password,
		MaxUsers:         maxUsers,
		Status:           model.RoomStatusActive,
		MediaURL:         req.MediaURL,
		MediaType:        req.MediaType,
		MediaTitle:       req.MediaTitle,
		MediaDuration:    mediaDuration,
		PlaybackState:    "paused",
		CurrentTime:      0,
		PlaybackRate:     1.0,
		Settings:         s.convertSettings(req.Settings),
		Version:          0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.roomRepo.Create(room); err != nil {
		return nil, err
	}

	return room, nil
}

// GetRoom 获取房间信息
func (s *RoomService) GetRoom(roomID string) (*model.Room, error) {
	return s.roomRepo.GetByID(roomID)
}

// UpdateRoom 更新房间信息
func (s *RoomService) UpdateRoom(roomID string, req *UpdateRoomRequest) (*model.Room, error) {
	// 构建更新字段映射
	updates := make(map[string]interface{})
	
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.IsPrivate != nil {
		updates["is_private"] = *req.IsPrivate
	}
	if req.Password != nil {
		updates["password"] = *req.Password
	}
	if req.MaxUsers != nil {
		updates["max_users"] = *req.MaxUsers
	}
	if req.Settings != nil {
		updates["settings"] = s.convertSettings(req.Settings)
	}

	updatedRoom, err := s.roomRepo.Update(roomID, updates)
	if err != nil {
		return nil, err
	}

	return updatedRoom, nil
}

// ListRooms 获取房间列表
func (s *RoomService) ListRooms(page, size int) ([]*model.Room, int64, error) {
	return s.roomRepo.GetActiveRooms(page, size)
}

// IsCreator 检查是否为房间创建者
func (s *RoomService) IsCreator(roomID, sessionID string) bool {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return false
	}
	return room.CreatorSessionID == sessionID
}

// DeleteRoom 删除房间（软删除）
func (s *RoomService) DeleteRoom(roomID, sessionID string) error {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}

	// 只能由创建者删除
	if room.CreatorSessionID != sessionID {
		return nil
	}

	updates := map[string]interface{}{
		"status":     model.RoomStatusDeleted,
		"updated_at": time.Now(),
	}

	_, err = s.roomRepo.Update(roomID, updates)
	return err
}

// CloseRoom 关闭房间
func (s *RoomService) CloseRoom(roomID string) error {
	updates := map[string]interface{}{
		"status":     model.RoomStatusDeleted,
		"updated_at": time.Now(),
	}

	_, err := s.roomRepo.Update(roomID, updates)
	return err
}

// PlayVideo 播放视频
func (s *RoomService) PlayVideo(roomID string) error {
	// 获取当前房间
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}

	// 更新播放状态
	updates := map[string]interface{}{
		"playback_state": "playing",
		"version":        room.Version + 1,
		"updated_at":     time.Now(),
	}

	err = s.roomRepo.UpdatePlaybackState(roomID, updates)
	return err
}

// PauseVideo 暂停视频
func (s *RoomService) PauseVideo(roomID string) error {
	// 获取当前房间
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}

	// 更新播放状态
	updates := map[string]interface{}{
		"playback_state": "paused",
		"version":        room.Version + 1,
		"updated_at":     time.Now(),
	}

	err = s.roomRepo.UpdatePlaybackState(roomID, updates)
	return err
}

// SeekVideo 视频跳转
func (s *RoomService) SeekVideo(roomID string, currentTime float64) error {
	// 获取当前房间
	room, err := s.GetRoom(roomID)
	if err != nil {
		return err
	}

	// 更新播放状态
	updates := map[string]interface{}{
		"current_time": currentTime,
		"version":      room.Version + 1,
		"updated_at":   time.Now(),
	}

	err = s.roomRepo.UpdatePlaybackState(roomID, updates)
	return err
}

// GetPlaybackStatus 获取播放状态
func (s *RoomService) GetPlaybackStatus(roomID string) (map[string]interface{}, error) {
	room, err := s.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	return room.GetPlaybackState(), nil
}

// convertSettings 转换设置格式
func (s *RoomService) convertSettings(settings interface{}) model.JSON {
	// 使用默认设置
	return model.JSON(`{
		"auto_play":      true,
		"allow_seek":     true,
		"default_volume": 1.0,
		"quality":        "auto",
		"subtitle_on":    false,
		"subtitle_lang":  "zh"
	}`)
}

// getDefaultPlaybackState 获取默认播放状态
func (s *RoomService) getDefaultPlaybackState() model.JSON {
	jsonStr := `{
		"current_time":  0.0,
		"duration":      0.0,
		"is_playing":    false,
		"playback_rate": 1.0,
		"video_url":     "",
		"video_title":   "",
		"last_updated":  %d
	}`
	
	playbackJSON := fmt.Sprintf(jsonStr, time.Now().Unix())
	return model.JSON(playbackJSON)
}