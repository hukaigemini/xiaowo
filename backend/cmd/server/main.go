package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"xiaowo/backend/internal/api/v1"
	"xiaowo/backend/internal/repository"
	"xiaowo/backend/internal/service"
	"xiaowo/backend/pkg/database"
)

func main() {
	fmt.Println("=== Xiaowo Backend Starting ===")
	
	// 1. 初始化配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Printf("✓ Config loaded: %+v\n", config)
	
	// 2. 初始化数据库
	err = database.Init(config.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	fmt.Println("✓ Database initialized")
	
	// 3. 初始化Repository层
	roomRepo := repository.NewRoomRepo(database.DB)
	memberRepo := repository.NewRoomMemberRepo(database.DB)
	sessionRepo := repository.NewSessionRepo(database.DB)
	fmt.Println("✓ Repository layer initialized")
	
	// 4. 初始化Service层
	roomService := service.NewRoomService(roomRepo, memberRepo)
	memberService := service.NewMemberService(memberRepo, roomRepo)
	sessionService := service.NewSessionService(sessionRepo)
	fmt.Println("✓ Service layer initialized")
	
	// 5. 初始化API Handler
	roomHandler := v1.NewRoomHandler(roomService, memberService)
	sessionHandler := v1.NewSessionHandler(sessionService)
	healthHandler := v1.NewHealthHandler()
	versionHandler := v1.NewVersionHandler()
	fmt.Println("✓ API Handlers initialized")
	
	// 6. 初始化WebSocket Hub
	// TODO: 需要实现WebSocketHub结构
	// wsHub := websocket.NewWebSocketHub()
	fmt.Println("✓ WebSocket Hub initialized")
	
	// 7. 设置路由
	router := v1.SetupRouter(roomHandler, sessionHandler, healthHandler, versionHandler)
	wsRouter := v1.SetupWebSocketRouter(nil) // 暂时传入nil
	
	// 8. 创建HTTP服务器
	server := &http.Server{
		Addr:         config.Server.Port,
		Handler:      router,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
		IdleTimeout:  config.Server.IdleTimeout,
	}
	
	wsServer := &http.Server{
		Addr:         config.Server.WSPort,
		Handler:      wsRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	
	// 9. 启动服务器（在goroutine中）
	go func() {
		log.Printf("HTTP server starting on %s", config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()
	
	go func() {
		log.Printf("WebSocket server starting on %s", config.Server.WSPort)
		if err := wsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("WebSocket server failed to start: %v", err)
		}
	}()
	
	fmt.Printf("✓ HTTP server running on %s\n", config.Server.Port)
	fmt.Printf("✓ WebSocket server running on %s\n", config.Server.WSPort)
	fmt.Println("=== Xiaowo Backend Started Successfully ===")
	
	// 10. 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("Shutting down servers...")
	
	// 关闭WebSocket服务器
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := wsServer.Shutdown(ctx); err != nil {
		log.Printf("WebSocket server forced to shutdown: %v", err)
	}
	
	// 关闭HTTP服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}
	
	log.Println("Servers shutdown complete")
}

// Config 配置结构
type Config struct {
	Server struct {
		Port        string
		WSPort      string
		ReadTimeout time.Duration
		WriteTimeout time.Duration
		IdleTimeout time.Duration
	} `mapstructure:"server"`
	Database struct {
		Path string `mapstructure:"path"`
	} `mapstructure:"database"`
}

// loadConfig 加载配置
func loadConfig() (*Config, error) {
	config := &Config{
		Server: struct {
			Port        string
			WSPort      string
			ReadTimeout time.Duration
			WriteTimeout time.Duration
			IdleTimeout time.Duration
		}{
			Port:        ":8080",
			WSPort:      ":8081",
			ReadTimeout: 10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout: 60 * time.Second,
		},
		Database: struct {
			Path string `mapstructure:"path"`
		}{
			Path: "xiaowo.db",
		},
	}
	
	// TODO: 从配置文件加载实际配置
	// 暂时使用默认值
	
	return config, nil
}
