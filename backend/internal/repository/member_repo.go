package repository

import (
	"xiaowo/backend/internal/model"

	"gorm.io/gorm"
)

// RoomMemberRepository interface defines all room member operations
type RoomMemberRepository interface {
	Join(member *model.RoomMember) error
	Leave(roomID, sessionID string) error
	FindMembers(roomID string) ([]*model.RoomMember, error)
	FindBySessionAndRoom(sessionID, roomID string) (*model.RoomMember, error)
	Update(member *model.RoomMember) error
	CountMembers(roomID string) (int64, error)
}

// RoomMemberRepo implements RoomMemberRepository
type RoomMemberRepo struct {
	db *gorm.DB
}

func NewRoomMemberRepo(db *gorm.DB) *RoomMemberRepo {
	return &RoomMemberRepo{db: db}
}

func (r *RoomMemberRepo) Join(member *model.RoomMember) error {
	return r.db.Create(member).Error
}

func (r *RoomMemberRepo) Leave(roomID, sessionID string) error {
	return r.db.Where("room_id = ? AND session_id = ?", roomID, sessionID).Delete(&model.RoomMember{}).Error
}

func (r *RoomMemberRepo) FindMembers(roomID string) ([]*model.RoomMember, error) {
	var members []*model.RoomMember
	if err := r.db.Where("room_id = ?", roomID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *RoomMemberRepo) FindBySessionAndRoom(sessionID, roomID string) (*model.RoomMember, error) {
	var member model.RoomMember
	if err := r.db.Where("session_id = ? AND room_id = ?", sessionID, roomID).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *RoomMemberRepo) Update(member *model.RoomMember) error {
	return r.db.Save(member).Error
}

func (r *RoomMemberRepo) CountMembers(roomID string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.RoomMember{}).Where("room_id = ?", roomID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
