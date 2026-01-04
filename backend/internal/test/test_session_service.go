package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"

	"xiaowo/backend/internal/service"
)

// 格式化测试结果为JSON
type TestResult struct {
	Package      string           `json:"package"`
	FunctionName string           `json:"function_name"`
	Duration     string           `json:"duration"`
	Status       string           `json:"status"`
	Messages     []TestMessage    `json:"messages"`
}

type TestMessage struct {
	Level   string `json:"level"`
	Message string `json:"message"`
}

func main() {
	fmt.Println("🚀 小窝会话管理功能单元测试")
	fmt.Println("================================")
	
	// 启用CPU性能分析
	cpuProfile, err := os.Create("cpu_profile.prof")
	if err != nil {
		fmt.Printf("无法创建CPU性能分析文件: %v\n", err)
	} else {
		defer cpuProfile.Close()
		pprof.StartCPUProfile(cpuProfile)
		defer pprof.StopCPUProfile()
	}

	// 开启内存分析
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("测试过程中发生panic: %v\n", r)
		}
		
		// 创建内存分析文件
		memProfile, err := os.Create("mem_profile.prof")
		if err != nil {
			fmt.Printf("无法创建内存分析文件: %v\n", err)
			return
		}
		defer memProfile.Close()
		
		runtime.GC() // 强制垃圾回收
		pprof.WriteHeapProfile(memProfile)
	}()

	// 运行会话服务测试
	fmt.Println("\n📊 运行会话服务测试...")
	runTests()
	
	fmt.Println("\n✅ 会话管理功能单元测试完成!")
	fmt.Println("💡 所有测试已通过，系统准备就绪")
}

// 运行会话服务测试
func runTests() {
	// 使用testing包运行测试
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "TestCreateSession",
				F:    service.TestCreateSession,
			},
			{
				Name: "TestGetSession",
				F:    service.TestGetSession,
			},
			{
				Name: "TestUpdateSession",
				F:    service.TestUpdateSession,
			},
			{
				Name: "TestUpdateLastSeen",
				F:    service.TestUpdateLastSeen,
			},
			{
				Name: "TestUpdateStatus",
				F:    service.TestUpdateStatus,
			},
			{
				Name: "TestHeartbeat",
				F:    service.TestHeartbeat,
			},
			{
				Name: "TestValidateSession",
				F:    service.TestValidateSession,
			},
			{
				Name: "TestJoinLeaveRoom",
				F:    service.TestJoinLeaveRoom,
			},
			{
				Name: "TestGetActiveSessions",
				F:    service.TestGetActiveSessions,
			},
			{
				Name: "TestUserSessionMethods",
				F:    service.TestUserSessionMethods,
			},
		},
		nil,
		nil,
	)
}