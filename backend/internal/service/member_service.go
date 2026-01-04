package service

import (
	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/repository"
	"time"
)

// MemberService 房间成员业务逻辑服务
type MemberService struct {
	memberRepo repository.RoomMemberRepository
	roomRepo   repository.RoomRepository
}

// NewMemberService 创建成员服务
func NewMemberService(memberRepo repository.RoomMemberRepository, roomRepo repository.RoomRepository) *MemberService {
	return &MemberService{
		memberRepo: memberRepo,
		roomRepo:   roomRepo,
	}
}

// AddMember 添加房间成员
func (s *MemberService) AddMember(member *model.RoomMember) error {
	member.JoinedAt = time.Now()
	member.LastSeen = time.Now()

	return s.memberRepo.Join(member)
}

// RemoveMember 移除房间成员
func (s *MemberService) RemoveMember(roomID, sessionID string) error {
	return s.memberRepo.Leave(roomID, sessionID)
}

// GetMember 获取房间成员信息
func (s *MemberService) GetMember(roomID, sessionID string) (*model.RoomMember, error) {
	return s.memberRepo.FindBySessionAndRoom(sessionID, roomID)
}

// GetMemberCount 获取房间成员数量
func (s *MemberService) GetMemberCount(roomID string) (int, error) {
	count, err := s.memberRepo.CountMembers(roomID)
	return int(count), err
}

// GetRoomMembers 获取房间所有成员
func (s *MemberService) GetRoomMembers(roomID string) ([]*model.RoomMember, error) {
	return s.memberRepo.FindMembers(roomID)
}

// UpdateMemberActivity 更新成员活动时间
func (s *MemberService) UpdateMemberActivity(roomID, sessionID string) error {
	member, err := s.GetMember(roomID, sessionID)
	if err != nil {
		return err
	}

	member.LastSeen = time.Now()
	return s.memberRepo.Update(member)
}

// IsMember 检查用户是否为房间成员
func (s *MemberService) IsMember(roomID, sessionID string) bool {
	_, err := s.GetMember(roomID, sessionID)
	return err == nil
}