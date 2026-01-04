package service

import (
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/repository"
)

// 关闭测试数据库
func closeTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// 初始化内存数据库用于测试
func initTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 迁移数据库模式
	if err := db.AutoMigrate(
		&model.Room{},
		&model.RoomMember{},
	); err != nil {
		return nil, err
	}

	return db, nil
}

// 测试MemberService的AddMember功能
func TestMemberService_AddMember(t *testing.T) {
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化服务
	memberRepo := repository.NewRoomMemberRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	memberService := NewMemberService(memberRepo, roomRepo)

	// 先创建一个房间
	room := &model.Room{
		ID:               "test-room-123",
		Name:             "测试房间",
		Description:      "测试房间描述",
		CreatorSessionID: "test-session",
		MediaURL:         "https://example.com/video.mp4",
		MediaType:        "video",
		MediaTitle:       "测试视频",
		MediaDuration:    300,
		Status:           model.RoomStatusActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := roomRepo.Create(room); err != nil {
		t.Fatalf("创建房间失败: %v", err)
	}

	// 测试用例1: 成功添加成员
	t.Run("TestAddMember_Success", func(t *testing.T) {
		member := &model.RoomMember{
			RoomID:    room.ID,
			SessionID: "session-123",
			Nickname:  "测试用户",
			Role:      model.MemberRoleUser,
		}

		err := memberService.AddMember(member)
		if err != nil {
			t.Fatalf("添加成员失败: %v", err)
		}

		// 验证成员已正确添加
		if member.JoinedAt.IsZero() {
			t.Error("加入时间应该被设置")
		}

		if member.LastSeen.IsZero() {
			t.Error("最后可见时间应该被设置")
		}

		// 验证成员确实存在于数据库中
		foundMember, err := memberService.GetMember(room.ID, member.SessionID)
		if err != nil {
			t.Fatalf("获取成员信息失败: %v", err)
		}

		if foundMember.RoomID != room.ID {
			t.Errorf("期望房间ID: %s, 实际: %s", room.ID, foundMember.RoomID)
		}

		if foundMember.SessionID != member.SessionID {
			t.Errorf("期望会话ID: %s, 实际: %s", member.SessionID, foundMember.SessionID)
		}
	})

	// 测试用例2: 添加重复成员
	t.Run("TestAddMember_Duplicate", func(t *testing.T) {
		member := &model.RoomMember{
			RoomID:    room.ID,
			SessionID: "session-456",
			Nickname:  "重复用户",
			Role:      model.MemberRoleUser,
		}

		// 第一次添加应该成功
		err := memberService.AddMember(member)
		if err != nil {
			t.Fatalf("第一次添加成员失败: %v", err)
		}

		// 第二次添加相同成员应该失败（因为有唯一约束）
		err = memberService.AddMember(member)
		if err == nil {
			t.Error("添加重复成员应该返回错误")
		}
	})
}

// 测试MemberService的RemoveMember功能
func TestMemberService_RemoveMember(t *testing.T) {
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化服务
	memberRepo := repository.NewRoomMemberRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	memberService := NewMemberService(memberRepo, roomRepo)

	// 先创建一个房间和成员
	room := &model.Room{
		ID:               "test-room-456",
		Name:             "测试房间2",
		Description:      "测试房间描述2",
		CreatorSessionID: "test-session",
		MediaURL:         "https://example.com/video2.mp4",
		MediaType:        "video",
		MediaTitle:       "测试视频2",
		MediaDuration:    400,
		Status:           model.RoomStatusActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := roomRepo.Create(room); err != nil {
		t.Fatalf("创建房间失败: %v", err)
	}

	member := &model.RoomMember{
		RoomID:    room.ID,
		SessionID: "session-789",
		Nickname:  "待删除用户",
		Role:      model.MemberRoleUser,
	}
	if err := memberService.AddMember(member); err != nil {
		t.Fatalf("添加成员失败: %v", err)
	}

	// 测试用例1: 成功移除成员
	t.Run("TestRemoveMember_Success", func(t *testing.T) {
		err := memberService.RemoveMember(room.ID, member.SessionID)
		if err != nil {
			t.Fatalf("移除成员失败: %v", err)
		}

		// 验证成员已被移除
		_, err = memberService.GetMember(room.ID, member.SessionID)
		if err == nil {
			t.Error("移除成员后应该无法获取到成员信息")
		}
	})

	// 测试用例2: 移除不存在的成员
	t.Run("TestRemoveMember_NonExistent", func(t *testing.T) {
		err := memberService.RemoveMember(room.ID, "non-existent-session")
		if err == nil {
			t.Error("移除不存在的成员应该返回错误")
		}
	})
}
