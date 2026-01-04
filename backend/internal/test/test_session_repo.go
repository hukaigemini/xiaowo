package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"

	"xiaowo/backend/internal/model"
)

// TestSessionRepo tests the Session Repository functionality
func TestSessionRepo(db *gorm.DB) {
	fmt.Println("🧪 开始测试Session Repository")
	
	repo := NewSessionRepo(db)

	// Test 1: Create session
	fmt.Println("\n📝 测试1: 创建会话")
	nickname := "测试用户"
	session, err := repo.Create(nickname)
	if err != nil {
		log.Fatalf("创建会话失败: %v", err)
	}
	fmt.Printf("✅ 创建会话成功: ID=%s, Nickname=%s\n", session.ID, session.Nickname)

	// Test 2: Get session by ID
	fmt.Println("\n🔍 测试2: 获取会话")
	foundSession, err := repo.GetByID(session.ID)
	if err != nil {
		log.Fatalf("获取会话失败: %v", err)
	}
	fmt.Printf("✅ 获取会话成功: %+v\n", foundSession)

	// Test 3: Update session
	fmt.Println("\n✏️ 测试3: 更新会话")
	updates := map[string]interface{}{
		"nickname": "更新后的用户",
	}
	updatedSession, err := repo.Update(session.ID, updates)
	if err != nil {
		log.Fatalf("更新会话失败: %v", err)
	}
	fmt.Printf("✅ 更新会话成功: Nickname=%s\n", updatedSession.Nickname)

	// Test 4: Update status
	fmt.Println("\n🔄 测试4: 更新会话状态")
	err = repo.UpdateStatus(session.ID, model.StatusOffline)
	if err != nil {
		log.Fatalf("更新会话状态失败: %v", err)
	}
	
	// Verify status change
	statusSession, err := repo.GetByID(session.ID)
	if err != nil {
		log.Fatalf("获取会话失败: %v", err)
	}
	
	if statusSession.Status != model.StatusOffline {
		log.Fatalf("状态更新失败，期望: %s, 实际: %s", model.StatusOffline, statusSession.Status)
	}
	fmt.Printf("✅ 更新会话状态成功: Status=%s\n", statusSession.Status)

	// Test 5: Get sessions by status
	fmt.Println("\n📊 测试5: 按状态获取会话")
	offlineSessions, err := repo.GetByStatus(model.StatusOffline)
	if err != nil {
		log.Fatalf("获取离线会话失败: %v", err)
	}
	fmt.Printf("✅ 离线会话数量: %d\n", len(offlineSessions))
	
	// Set session back to online
	err = repo.UpdateStatus(session.ID, model.StatusOnline)
	if err != nil {
		log.Fatalf("恢复会话状态失败: %v", err)
	}

	// Test 6: Generate nickname and avatar
	fmt.Println("\n🎨 测试6: 生成昵称和头像")
	for i := 0; i < 5; i++ {
		nickname := repo.GenerateNickname()
		avatar := repo.GenerateAvatar()
		fmt.Printf("   生成 %d: %s -> %s\n", i+1, nickname, avatar)
	}

	// Test 7: Get active sessions
	fmt.Println("\n📊 测试7: 获取活跃会话")
	activeSessions, err := repo.GetActiveSessions()
	if err != nil {
		log.Fatalf("获取活跃会话失败: %v", err)
	}
	fmt.Printf("✅ 活跃会话数量: %d\n", len(activeSessions))

	// Test 8: Update last seen
	fmt.Println("\n👁️ 测试8: 更新最后在线时间")
	err = repo.UpdateLastSeen(session.ID)
	if err != nil {
		log.Fatalf("更新最后在线时间失败: %v", err)
	}
	fmt.Println("✅ 更新最后在线时间成功")

	// Test 9: Soft delete
	fmt.Println("\n🗑️ 测试9: 软删除会话")
	err = repo.SoftDelete(session.ID)
	if err != nil {
		log.Fatalf("软删除会话失败: %v", err)
	}
	
	// Verify soft delete
	deletedSession, err := repo.GetByID(session.ID)
	if err != nil {
		fmt.Printf("✅ 软删除检查正常: %v\n", err)
	} else {
		fmt.Printf("❌ 软删除检查异常: 会话仍然存在: %+v\n", deletedSession)
	}

	// Test 10: Test session expiration check
	fmt.Println("\n⏰ 测试10: 会话过期检查")
	expiredSession, err := repo.Create("过期测试用户")
	if err != nil {
		log.Fatalf("创建过期测试会话失败: %v", err)
	}

	// Manually set session to expired (for testing)
	expiredSession.ExpiresAt = time.Now().Add(-1 * time.Hour)
	err = db.Save(expiredSession).Error
	if err != nil {
		log.Fatalf("设置会话过期失败: %v", err)
	}

	_, err = repo.GetByID(expiredSession.ID)
	if err != nil {
		fmt.Printf("✅ 过期检查正常: %v\n", err)
	} else {
		fmt.Println("❌ 过期检查异常: 应该返回错误但没有")
	}

	// Clean up test data
	err = repo.Delete(session.ID)
	if err != nil {
		log.Printf("清理测试数据失败: %v", err)
	}
	err = repo.Delete(expiredSession.ID)
	if err != nil {
		log.Printf("清理测试数据失败: %v", err)
	}

	fmt.Println("\n🎉 Session Repository测试完成!")
}

// DemoSessionRepo demonstrates the Session Repository usage
func DemoSessionRepo(db *gorm.DB) {
	fmt.Println("🎬 小窝Session Repository演示")
	fmt.Println(strings.Repeat("=", 50))

	repo := NewSessionRepo(db)

	// Create multiple sessions
	fmt.Println("\n1. 创建多个用户会话")
	sessions := make([]*model.UserSession, 5)
	for i := 0; i < 5; i++ {
		session, err := repo.Create("")
		if err != nil {
			log.Printf("创建会话 %d 失败: %v", i+1, err)
			continue
		}
		sessions[i] = session
		fmt.Printf("   会话 %d: %s (%s)\n", i+1, session.Nickname, session.Avatar)
	}

	// Simulate room joining
	fmt.Println("\n2. 模拟加入房间")
	roomID := "DEMO123"
	for i := 0; i < 3 && i < len(sessions) && sessions[i] != nil; i++ {
		err := repo.JoinRoom(sessions[i].ID, roomID)
		if err != nil {
			log.Printf("会话 %s 加入房间失败: %v", sessions[i].ID, err)
			continue
		}
		fmt.Printf("   %s 加入房间 %s\n", sessions[i].Nickname, roomID)
	}

	// Update last seen for some sessions
	fmt.Println("\n3. 更新在线状态")
	for i := 0; i < 3 && i < len(sessions) && sessions[i] != nil; i++ {
		err := repo.UpdateLastSeen(sessions[i].ID)
		if err != nil {
			log.Printf("更新 %s 在线状态失败: %v", sessions[i].Nickname, err)
			continue
		}
		fmt.Printf("   %s 更新在线状态\n", sessions[i].Nickname)
	}

	// Get active sessions
	fmt.Println("\n4. 获取活跃会话")
	activeSessions, err := repo.GetActiveSessions()
	if err != nil {
		log.Printf("获取活跃会话失败: %v", err)
	} else {
		fmt.Printf("   活跃会话总数: %d\n", len(activeSessions))
		for _, session := range activeSessions {
			roomStatus := "未加入房间"
			if session.RoomID != nil {
				roomStatus = fmt.Sprintf("在房间 %s", *session.RoomID)
			}
			fmt.Printf("   - %s: %s\n", session.Nickname, roomStatus)
		}
	}

	// Cleanup
	fmt.Println("\n5. 清理测试数据")
	for i, session := range sessions {
		if session != nil {
			err := repo.Delete(session.ID)
			if err != nil {
				log.Printf("删除会话 %d 失败: %v", i+1, err)
			} else {
				fmt.Printf("   删除会话: %s\n", session.Nickname)
			}
		}
	}

	fmt.Println("\n✨ Session Repository演示完成!")
}