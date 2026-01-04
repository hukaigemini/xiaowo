package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/repository"
)

// 初始化内存数据库用于测试
func initTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open test database: %w", err)
	}

	// 迁移数据库模式
	if err := db.AutoMigrate(
		&model.UserSession{},
		&model.Room{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

// 关闭测试数据库
func closeTestDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}
}

// 测试CreateSession功能
func TestCreateSession(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 测试用例1: 使用指定昵称创建会话
	t.Run("TestCreateSessionWithNickname", func(t *testing.T) {
		nickname := "测试用户1"
		session, err := sessionService.CreateSession(nickname)
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}

		// 验证返回的会话信息
		if session.Nickname != nickname {
			t.Errorf("期望昵称: %s, 实际: %s", nickname, session.Nickname)
		}

		// 检查会话ID不为空
		if session.ID == "" {
			t.Error("会话ID不应为空")
		}

		// 检查头像URL不为空
		if session.Avatar == "" {
			t.Error("头像URL不应为空")
		}

		// 检查状态为在线
		if session.Status != model.StatusOnline {
			t.Errorf("期望状态: %s, 实际: %s", model.StatusOnline, session.Status)
		}

		// 检查过期时间是否正确（7天后）
		expectedExpires := time.Now().Add(7 * 24 * time.Hour)
		if session.ExpiresAt.Before(expectedExpires.Add(-time.Minute)) || 
		   session.ExpiresAt.After(expectedExpires.Add(time.Minute)) {
			t.Errorf("期望过期时间: %v, 实际: %v", expectedExpires, session.ExpiresAt)
		}
	})

	// 测试用例2: 使用空昵称创建会话（自动生成）
	t.Run("TestCreateSessionWithoutNickname", func(t *testing.T) {
		session, err := sessionService.CreateSession("")
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}

		// 检查是否自动生成了昵称
		if session.Nickname == "" {
			t.Error("自动生成的昵称不应为空")
		}

		// 检查头像URL不为空
		if session.Avatar == "" {
			t.Error("头像URL不应为空")
		}
	})
}

// 测试GetSession功能
func TestGetSession(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户2"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 测试获取会话
	t.Run("TestGetSession", func(t *testing.T) {
		foundSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证返回的会话信息
		if foundSession.ID != session.ID {
			t.Errorf("期望ID: %s, 实际: %s", session.ID, foundSession.ID)
		}

		if foundSession.Nickname != session.Nickname {
			t.Errorf("期望昵称: %s, 实际: %s", session.Nickname, foundSession.Nickname)
		}
	})

	// 测试获取不存在的会话
	t.Run("TestGetNonExistentSession", func(t *testing.T) {
		nonExistentID := "non-existent-id"
		_, err := sessionService.GetSession(nonExistentID)
		if err == nil {
			t.Error("获取不存在的会话应该返回错误")
		}
	})
}

// 测试UpdateSession功能
func TestUpdateSession(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户3"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 测试更新会话昵称
	t.Run("TestUpdateNickname", func(t *testing.T) {
		newNickname := "更新后的昵称"
		updates := map[string]interface{}{
			"nickname": newNickname,
		}

		updatedSession, err := sessionService.UpdateSession(session.ID, updates)
		if err != nil {
			t.Fatalf("更新会话失败: %v", err)
		}

		// 验证昵称已更新
		if updatedSession.Nickname != newNickname {
			t.Errorf("期望昵称: %s, 实际: %s", newNickname, updatedSession.Nickname)
		}
	})

	// 测试更新不存在的会话
	t.Run("TestUpdateNonExistentSession", func(t *testing.T) {
		updates := map[string]interface{}{
			"nickname": "不存在的用户",
		}

		_, err := sessionService.UpdateSession("non-existent-id", updates)
		if err == nil {
			t.Error("更新不存在的会话应该返回错误")
		}
	})
}

// 测试UpdateLastSeen功能
func TestUpdateLastSeen(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户4"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 记录原始时间
	originalLastSeen := session.LastSeenAt

	// 等待一点时间
	time.Sleep(1 * time.Second)

	// 测试更新最后在线时间
	t.Run("TestUpdateLastSeen", func(t *testing.T) {
		err := sessionService.UpdateLastSeen(session.ID)
		if err != nil {
			t.Fatalf("更新最后在线时间失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证最后在线时间已更新
		if !updatedSession.LastSeenAt.After(originalLastSeen) {
			t.Errorf("最后在线时间应该更新，当前: %v, 原: %v", updatedSession.LastSeenAt, originalLastSeen)
		}
	})

	// 测试更新不存在的会话
	t.Run("TestUpdateNonExistentSessionLastSeen", func(t *testing.T) {
		err := sessionService.UpdateLastSeen("non-existent-id")
		if err == nil {
			t.Error("更新不存在的会话应该返回错误")
		}
	})
}

// 测试UpdateStatus功能
func TestUpdateStatus(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户5"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 测试更新状态为离线
	t.Run("TestUpdateStatusToOffline", func(t *testing.T) {
		err := sessionService.UpdateStatus(session.ID, string(model.StatusOffline))
		if err != nil {
			t.Fatalf("更新会话状态失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证状态已更新为离线
		if updatedSession.Status != model.StatusOffline {
			t.Errorf("期望状态: %s, 实际: %s", model.StatusOffline, updatedSession.Status)
		}
	})

	// 测试更新状态为在线
	t.Run("TestUpdateStatusToOnline", func(t *testing.T) {
		err := sessionService.UpdateStatus(session.ID, string(model.StatusOnline))
		if err != nil {
			t.Fatalf("更新会话状态失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证状态已更新为在线
		if updatedSession.Status != model.StatusOnline {
			t.Errorf("期望状态: %s, 实际: %s", model.StatusOnline, updatedSession.Status)
		}
	})

	// 测试更新无效状态
	t.Run("TestUpdateInvalidStatus", func(t *testing.T) {
		err := sessionService.UpdateStatus(session.ID, "invalid-status")
		if err == nil {
			t.Error("更新无效状态应该返回错误")
		}
	})
}

// 测试Heartbeat功能
func TestHeartbeat(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户6"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 记录原始时间和状态
	originalLastSeen := session.LastSeenAt
	originalStatus := session.Status

	// 等待一点时间
	time.Sleep(1 * time.Second)

	// 测试心跳功能
	t.Run("TestHeartbeat", func(t *testing.T) {
		err := sessionService.Heartbeat(session.ID)
		if err != nil {
			t.Fatalf("心跳失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证最后在线时间已更新
		if !updatedSession.LastSeenAt.After(originalLastSeen) {
			t.Errorf("最后在线时间应该更新，当前: %v, 原: %v", updatedSession.LastSeenAt, originalLastSeen)
		}

		// 验证状态已更新为在线
		if updatedSession.Status != model.StatusOnline {
			t.Errorf("期望状态: %s, 实际: %s", model.StatusOnline, updatedSession.Status)
		}
	})

	// 验证原始状态和更新后的状态不同
	if originalStatus == model.StatusOnline {
		t.Log("原始状态已经是在线，心跳测试仍然会更新最后在线时间")
	}
}

// 测试ValidateSession功能
func TestValidateSession(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个会话
	nickname := "测试用户7"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 测试验证有效会话
	t.Run("TestValidateValidSession", func(t *testing.T) {
		isValid, err := sessionService.ValidateSession(session.ID)
		if err != nil {
			t.Fatalf("验证会话失败: %v", err)
		}
		if !isValid {
			t.Error("会话应该是有效的")
		}
	})

	// 测试验证不存在的会话
	t.Run("TestValidateNonExistentSession", func(t *testing.T) {
		isValid, err := sessionService.ValidateSession("non-existent-id")
		if err == nil {
			t.Error("验证不存在的会话应该返回错误")
		}
		if isValid {
			t.Error("不存在的会话应该不是有效的")
		}
	})
}

// 测试JoinRoom和LeaveRoom功能
func TestJoinLeaveRoom(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个房间
	roomName := "测试房间"
	room, err := createTestRoom(db, roomName)
	if err != nil {
		t.Fatalf("创建房间失败: %v", err)
	}

	// 先创建一个会话
	nickname := "测试用户8"
	session, err := sessionService.CreateSession(nickname)
	if err != nil {
		t.Fatalf("创建会话失败: %v", err)
	}

	// 测试加入房间
	t.Run("TestJoinRoom", func(t *testing.T) {
		err := sessionService.JoinRoom(session.ID, room.ID)
		if err != nil {
			t.Fatalf("加入房间失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证会话已加入房间
		if updatedSession.RoomID == nil || *updatedSession.RoomID != room.ID {
			t.Errorf("期望房间ID: %s, 实际: %s", room.ID, *updatedSession.RoomID)
		}
	})

	// 测试离开房间
	t.Run("TestLeaveRoom", func(t *testing.T) {
		err := sessionService.LeaveRoom(session.ID)
		if err != nil {
			t.Fatalf("离开房间失败: %v", err)
		}

		// 获取更新后的会话
		updatedSession, err := sessionService.GetSession(session.ID)
		if err != nil {
			t.Fatalf("获取会话失败: %v", err)
		}

		// 验证会话已离开房间
		if updatedSession.RoomID != nil {
			t.Errorf("期望房间ID: nil, 实际: %s", *updatedSession.RoomID)
		}
	})

	// 测试加入不存在的房间
	t.Run("TestJoinNonExistentRoom", func(t *testing.T) {
		err := sessionService.JoinRoom(session.ID, "non-existent-room")
		if err == nil {
			t.Error("加入不存在的房间应该返回错误")
		}
	})
}

// 测试GetActiveSessions功能
func TestGetActiveSessions(t *testing.T) {
	// 初始化测试数据库
	db, err := initTestDB()
	if err != nil {
		t.Fatalf("无法初始化测试数据库: %v", err)
	}
	defer closeTestDB(db)

	// 初始化仓库和服务
	sessionRepo := repository.NewSessionRepo(db)
	sessionService := NewSessionService(sessionRepo)

	// 先创建一个房间
	roomName := "测试房间2"
	room, err := createTestRoom(db, roomName)
	if err != nil {
		t.Fatalf("创建房间失败: %v", err)
	}

	// 创建多个会话
	sessions := make([]*model.UserSession, 3)
	for i := 0; i < 3; i++ {
		nickname := fmt.Sprintf("测试用户%d", i+9)
		session, err := sessionService.CreateSession(nickname)
		if err != nil {
			t.Fatalf("创建会话失败: %v", err)
		}

		// 第一个和第二个会话加入房间，第三个不加入
		if i < 2 {
			err := sessionService.JoinRoom(session.ID, room.ID)
			if err != nil {
				t.Fatalf("加入房间失败: %v", err)
			}
		}

		sessions[i] = session
	}

	// 测试获取活跃会话
	t.Run("TestGetActiveSessions", func(t *testing.T) {
		activeSessions, err := sessionService.GetActiveSessions()
		if err != nil {
			t.Fatalf("获取活跃会话失败: %v", err)
		}

		// 验证活跃会话数量
		if len(activeSessions) != 2 {
			t.Errorf("期望活跃会话数: 2, 实际: %d", len(activeSessions))
		}

		// 验证返回的会话ID
		var foundIDs []string
		for _, session := range activeSessions {
			foundIDs = append(foundIDs, session.ID)
		}

		if len(foundIDs) != 2 {
			t.Errorf("返回的会话ID数量不正确: %d", len(foundIDs))
		}

		// 验证返回的会话确实是加入房间的
		for _, session := range activeSessions {
			if session.RoomID == nil || *session.RoomID != room.ID {
				t.Errorf("会话 %s 应该属于房间 %s", session.ID, room.ID)
			}
		}
	})

	// 测试获取过期会话
	t.Run("TestGetExpiredSessions", func(t *testing.T) {
		// 创建一个过期的会话
		expiredSession, err := createExpiredSession(db, "过期用户")
		if err != nil {
			t.Fatalf("创建过期会话失败: %v", err)
		}

		expiredSessions, err := sessionService.GetExpiredSessions()
		if err != nil {
			t.Fatalf("获取过期会话失败: %v", err)
		}

		// 验证过期会话列表中包含我们的过期会话
		found := false
		for _, session := range expiredSessions {
			if session.ID == expiredSession.ID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("过期会话列表中应该包含会话 %s", expiredSession.ID)
		}
	})

	// 测试根据状态获取会话
	t.Run("TestGetSessionsByStatus", func(t *testing.T) {
		// 设置一个会话为离线状态
		err := sessionService.UpdateStatus(sessions[0].ID, string(model.StatusOffline))
		if err != nil {
			t.Fatalf("更新会话状态失败: %v", err)
		}

		offlineSessions, err := sessionService.GetSessionsByStatus(string(model.StatusOffline))
		if err != nil {
			t.Fatalf("获取离线会话失败: %v", err)
		}

		// 验证离线会话列表中包含我们的离线会话
		found := false
		for _, session := range offlineSessions {
			if session.ID == sessions[0].ID {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("离线会话列表中应该包含会话 %s", sessions[0].ID)
		}
	})
}

// 创建测试房间
func createTestRoom(db *gorm.DB, name string) (*model.Room, error) {
	room := &model.Room{
		ID:          fmt.Sprintf("room-%d", time.Now().Unix()),
		Name:        name,
		Description: "测试房间描述",
		CreatedAt:   time.Now(),
	}

	if err := db.Create(room).Error; err != nil {
		return nil, err
	}

	return room, nil
}

// 创建过期会话
func createExpiredSession(db *gorm.DB, nickname string) (*model.UserSession, error) {
	repo := repository.NewSessionRepo(db)
	session, err := repo.Create(nickname)
	if err != nil {
		return nil, err
	}

	// 设置过期时间为过去
	expiredSession := &model.UserSession{
		ID:          session.ID,
		Nickname:    session.Nickname,
		Avatar:      session.Avatar,
		RoomID:      session.RoomID,
		Status:      session.Status,
		CreatedAt:   session.CreatedAt,
		LastSeenAt:  session.LastSeenAt,
		ExpiresAt:   time.Now().Add(-time.Hour), // 1小时前过期
		DeletedAt:   session.DeletedAt,
	}

	if err := db.Model(&model.UserSession{}).Where("id = ?", session.ID).Updates(expiredSession).Error; err != nil {
		return nil, err
	}

	return expiredSession, nil
}

// 格式化测试结果为JSON
func formatTestResult(testName string, passed bool, message string) string {
	result := map[string]interface{}{
		"test_name": testName,
		"passed":    passed,
		"message":   message,
	}

	jsonResult, _ := json.Marshal(result)
	return string(jsonResult)
}

// 测试UserSession模型方法
func TestUserSessionMethods(t *testing.T) {
	now := time.Now()
	session := &model.UserSession{
		ID:          "test-session-id",
		Nickname:    "测试用户",
		Avatar:      "https://example.com/avatar.png",
		Status:      model.StatusOnline,
		CreatedAt:   now,
		LastSeenAt:  now,
		ExpiresAt:   now.Add(7 * 24 * time.Hour), // 7天后过期
	}

	// 测试IsExpired方法
	t.Run("TestIsExpired", func(t *testing.T) {
		// 创建未过期的会话
		if session.IsExpired() {
			t.Error("未过期的会话应该返回false")
		}

		// 创建已过期的会话
		expiredSession := *session
		expiredSession.ExpiresAt = now.Add(-time.Hour) // 1小时前过期
		if !expiredSession.IsExpired() {
			t.Error("已过期的会话应该返回true")
		}
	})

	// 测试IsOnline方法
	t.Run("TestIsOnline", func(t *testing.T) {
		if !session.IsOnline() {
			t.Error("在线状态的会话应该返回true")
		}

		// 创建离线会话
		offlineSession := *session
		offlineSession.Status = model.StatusOffline
		if offlineSession.IsOnline() {
			t.Error("离线状态的会话应该返回false")
		}
	})

	// 测试UpdateLastSeen方法
	t.Run("TestUpdateLastSeen", func(t *testing.T) {
		originalLastSeen := session.LastSeenAt

		// 等待一点时间
		time.Sleep(1 * time.Second)

		session.UpdateLastSeen()

		if !session.LastSeenAt.After(originalLastSeen) {
			t.Errorf("最后在线时间应该更新，当前: %v, 原: %v", session.LastSeenAt, originalLastSeen)
		}

		if session.Status != model.StatusOnline {
			t.Errorf("更新最后在线时间后状态应该为在线，实际: %s", session.Status)
		}
	})

	// 测试SetOffline方法
	t.Run("TestSetOffline", func(t *testing.T) {
		session.SetOffline()

		if session.Status != model.StatusOffline {
			t.Errorf("设置离线后状态应该为离线，实际: %s", session.Status)
		}
	})

	// 测试IsActive方法
	t.Run("TestIsActive", func(t *testing.T) {
		// 创建完整活动的会话（在线且在房间中）
		roomID := "room-123"
		session.RoomID = &roomID
		session.Status = model.StatusOnline

		if !session.IsActive() {
			t.Error("完整活动的会话应该返回true")
		}

		// 测试不在线的会话
		session.Status = model.StatusOffline
		if session.IsActive() {
			t.Error("不在线的会话应该返回false")
		}

		// 恢复在线状态
		session.Status = model.StatusOnline

		// 测试不在房间中的会话
		session.RoomID = nil
		if session.IsActive() {
			t.Error("不在房间中的会话应该返回false")
		}

		// 测试过期的会话
		expiredSession := *session
		expiredSession.ExpiresAt = time.Now().Add(-time.Hour) // 1小时前过期
		if expiredSession.IsActive() {
			t.Error("过期的会话应该返回false")
		}
	})
}