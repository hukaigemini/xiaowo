package main

import (
	"fmt"
	"log"
	"time"

	"xiaowo/backend/internal/repository"
)

func main() {
	fmt.Println("🚀 小窝Session Repository完整功能测试")
	fmt.Println("=" * 60)
	
	// 1. 初始化数据库连接
	fmt.Println("\n📡 正在初始化数据库连接...")
	db, err := repository.InitOptimizedDB()
	if err != nil {
		log.Fatalf("数据库连接初始化失败: %v", err)
	}
	defer repository.Close(db)
	
	fmt.Println("✅ 数据库连接初始化成功")
	
	// 2. 健康检查
	fmt.Println("\n🔍 执行数据库健康检查...")
	health := repository.HealthCheck(db)
	if !health.IsHealthy {
		log.Fatalf("数据库健康检查失败: %s", health.Message)
	}
	fmt.Println("✅ 数据库健康检查通过")
	
	// 3. 执行Session Repository测试
	fmt.Println("\n🧪 开始Session Repository功能测试...")
	repository.TestSessionRepo(db)
	
	// 4. 执行Session Repository演示
	fmt.Println("\n🎬 Session Repository演示...")
	repository.DemoSessionRepo(db)
	
	// 5. 清理测试数据
	fmt.Println("\n🧹 清理测试数据...")
	cleanedCount, err := repository.CleanupExpired(db)
	if err != nil {
		log.Printf("清理过期数据失败: %v", err)
	} else {
		fmt.Printf("✅ 清理过期数据: %d 条记录\n", cleanedCount)
	}
	
	fmt.Println("\n🎉 Session Repository完整测试完成!")
	fmt.Println("💡 Session Repository已经可以使用，可以开始API层开发")
}

// 运行命令:
// go run backend/internal/repository/test_session_integration.go
//
// 预期输出:
// 🚀 小窝Session Repository完整功能测试
// ============================================================
//
// 📡 正在初始化数据库连接...
// ✅ 数据库连接初始化成功
//
// 🔍 执行数据库健康检查...
// ✅ 数据库健康检查通过
//
// 🧪 开始Session Repository功能测试...
// 🧪 开始测试Session Repository
//
// 📝 测试1: 创建会话
// ✅ 创建会话成功: ID=xxx, Nickname=测试用户
//
// 🔍 测试2: 获取会话
// ✅ 获取会话成功: ...
//
// ✏️ 测试3: 更新会话
// ✅ 更新会话成功: Nickname=更新后的用户
//
// 🎨 测试4: 生成昵称和头像
//    生成 1: 快乐的小熊猫 -> https://api.dicebear.com/...
//    生成 2: 聪明的小猫咪 -> https://api.dicebear.com/...
//    ...
//
// 📊 测试5: 获取活跃会话
// ✅ 活跃会话数量: 3
//
// 👁️ 测试6: 更新最后在线时间
// ✅ 更新最后在线时间成功
//
// ⏰ 测试7: 会话过期检查
// ✅ 过期检查正常: session expired: xxx
//
// 🎉 Session Repository测试完成!
//
// 🎬 Session Repository演示...
// 🎬 小窝Session Repository演示
// ===================================================
//
// 1. 创建多个用户会话
//    会话 1: 快乐的小熊猫 (https://api.dicebear.com/...)
//    会话 2: 聪明的小猫咪 (https://api.dicebear.com/...)
//    ...
//
// 2. 模拟加入房间
//    快乐的小熊猫 加入房间 DEMO123
//    聪明的小猫咪 加入房间 DEMO123
//    ...
//
// 3. 更新在线状态
//    快乐的小熊猫 更新在线状态
//    ...
//
// 4. 获取活跃会话
//    活跃会话总数: 3
//    - 快乐的小熊猫: 在房间 DEMO123
//    - 聪明的小猫咪: 在房间 DEMO123
//    - ...
//
// 5. 清理测试数据
//    删除会话: 快乐的小熊猫
//    ...
//
// ✨ Session Repository演示完成!
//
// 🧹 清理测试数据...
// ✅ 清理过期数据: 2 条记录
//
// 🎉 Session Repository完整测试完成!
// 💡 Session Repository已经可以使用，可以开始API层开发