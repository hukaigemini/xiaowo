package main

import (
	"fmt"
	"log"
	"time"

	"xiaowo/backend/internal/repository"
)

func main() {
	fmt.Println("🚀 小窝数据库连接测试")
	
	// 1. 初始化数据库连接
	fmt.Println("\n📡 正在初始化数据库连接...")
	db, err := repository.InitOptimizedDB()
	if err != nil {
		log.Fatalf("数据库连接初始化失败: %v", err)
	}
	defer repository.Close(db)
	
	fmt.Println("✅ 数据库连接初始化成功")
	
	// 2. 执行数据库健康检查
	fmt.Println("\n🔍 执行数据库健康检查...")
	health := repository.HealthCheck(db)
	
	if health.IsHealthy {
		fmt.Printf("✅ 数据库健康状态: %s\n", health.Message)
		if health.PingLatency != "" {
			fmt.Printf("⏱️  Ping延迟: %s\n", health.PingLatency)
		}
		
		if health.ConnStats != nil {
			stats := health.ConnStats
			fmt.Printf("📊 连接池统计:\n")
			fmt.Printf("   - 活跃连接数: %d\n", stats.InUse)
			fmt.Printf("   - 空闲连接数: %d\n", stats.Idle)
			fmt.Printf("   - 总连接数: %d\n", stats.OpenConnections)
			fmt.Printf("   - 等待连接数: %d\n", stats.WaitCount)
			fmt.Printf("   - 等待时间: %v\n", stats.WaitDuration)
		}
	} else {
		fmt.Printf("❌ 数据库健康检查失败: %s\n", health.Message)
		return
	}
	
	// 3. 验证数据库模式
	fmt.Println("\n🔍 验证数据库模式...")
	if err := repository.ValidateSchema(db); err != nil {
		log.Fatalf("数据库模式验证失败: %v", err)
	}
	fmt.Println("✅ 数据库模式验证通过")
	
	// 4. 执行带重试的ping测试
	fmt.Println("\n🔄 执行带重试的连接验证...")
	if err := repository.PingWithRetry(db, 3, 1*time.Second); err != nil {
		log.Fatalf("连接验证失败: %v", err)
	}
	fmt.Println("✅ 连接验证成功")
	
	// 5. 执行一个简单的查询测试
	fmt.Println("\n🔍 执行简单查询测试...")
	var count int
	result := db.Raw("SELECT COUNT(*) FROM system_configs").Scan(&count)
	if result.Error != nil {
		log.Fatalf("查询测试失败: %v", result.Error)
	}
	fmt.Printf("✅ 系统配置表记录数: %d\n", count)
	
	fmt.Println("\n🎉 数据库连接测试全部通过！")
	fmt.Println("💡 提示: 可以开始使用此连接进行API开发")
}

// 运行命令:
// go run backend/internal/repository/test_db_connection.go
//
// 输出示例:
// 🚀 小窝数据库连接测试
//
// 📡 正在初始化数据库连接...
// ✅ 数据库连接初始化成功
//
// 🔍 执行数据库健康检查...
// ✅ 数据库健康状态: 数据库连接正常
// ⏱️  Ping延迟: 1.234ms
// 📊 连接池统计:
//    - 活跃连接数: 1
//    - 空闲连接数: 2
//    - 总连接数: 3
//    - 等待连接数: 0
//    - 等待时间: 0s
//
// 🔍 验证数据库模式...
// ✅ 数据库模式验证通过
//
// 🔄 执行带重试的连接验证...
// ✅ 连接验证成功
//
// 🔍 执行简单查询测试...
// ✅ 系统配置表记录数: 15
//
// 🎉 数据库连接测试全部通过！
// 💡 提示: 可以开始使用此连接进行API开发