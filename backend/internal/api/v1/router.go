package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	gorillaWs "github.com/gorilla/websocket"
	"xiaowo/backend/internal/websocket"
)

// SetupRouter 设置路由
func SetupRouter(roomHandler *RoomHandler, sessionHandler *SessionHandler, healthHandler *HealthHandler, versionHandler *VersionHandler) *gin.Engine {
	// 设置为发布模式（生产环境）
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.New()
	
	// 全局中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())
	router.Use(HandleError())
	
	// 健康检查路由（不需要认证）
	healthGroup := router.Group("/")
	{
		healthGroup.GET("/health", healthHandler.HealthCheck)
		healthGroup.GET("/ready", healthHandler.ReadinessCheck)
		healthGroup.GET("/version", versionHandler.GetVersion)
	}
	
	// API v1 路由
	v1 := router.Group("/api/v1")
	{
		// 房间相关路由
		roomGroup := v1.Group("/rooms")
		{
			roomGroup.POST("", roomHandler.CreateRoom)
			roomGroup.GET("", roomHandler.ListRooms)
			roomGroup.GET("/:room_id", roomHandler.GetRoom)
			roomGroup.PUT("/:room_id", roomHandler.UpdateRoom)
			roomGroup.DELETE("/:room_id", roomHandler.CloseRoom)
			roomGroup.GET("/:room_id/members", roomHandler.GetRoomMembers)
			roomGroup.POST("/:room_id/join", roomHandler.JoinRoom)
			roomGroup.POST("/:room_id/leave", roomHandler.LeaveRoom)
			roomGroup.POST("/:room_id/play", roomHandler.PlayVideo)
			roomGroup.POST("/:room_id/pause", roomHandler.PauseVideo)
			roomGroup.POST("/:room_id/seek", roomHandler.SeekVideo)
			roomGroup.GET("/:room_id/status", roomHandler.GetPlaybackStatus)
		}
		
		// 会话相关路由
		sessionGroup := v1.Group("/sessions")
		{
			sessionGroup.POST("", sessionHandler.CreateSession)
			sessionGroup.GET("/:session_id", sessionHandler.GetSession)
			sessionGroup.PUT("/:session_id", sessionHandler.UpdateSession)
			sessionGroup.POST("/:session_id/heartbeat", sessionHandler.Heartbeat)
			sessionGroup.GET("/:session_id/validate", sessionHandler.ValidateSession)
			sessionGroup.DELETE("/:session_id", sessionHandler.DeleteSession)
		}
	}
	
	return router
}

// SetupWebSocketRouter 设置 WebSocket 路由
func SetupWebSocketRouter(hub *websocket.WebSocketHub) *gin.Engine {
	router := gin.New()
	
	// WebSocket 连接路由
	router.GET("/ws/room/:room_id", func(c *gin.Context) {
		roomID := c.Param("room_id")
		token := c.Query("token")
		
		if roomID == "" {
			c.JSON(400, gin.H{"error": "房间ID不能为空"})
			c.Abort()
			return
		}
		
		if token == "" {
			c.JSON(400, gin.H{"error": "访问令牌不能为空"})
			c.Abort()
			return
		}
		
		// TODO: 验证令牌和房间权限
		
		// 升级为 WebSocket 连接
		WebSocketHandler(c.Writer, c.Request, hub, roomID, token)
	})
	
	return router
}

// WebSocketHandler WebSocket 连接处理器
func WebSocketHandler(w http.ResponseWriter, r *http.Request, hub *websocket.WebSocketHub, roomID, token string) {
	// 解析 token 获取 session_id
	sessionID := parseTokenSessionID(token)
	if sessionID == "" {
		http.Error(w, "无效的访问令牌", http.StatusUnauthorized)
		return
	}
	
	// 配置 WebSocket 升级器
	upgrader := gorillaWs.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 允许所有来源，实际生产环境应该配置允许的域名
			return true
		},
	}
	
	// 升级 HTTP 连接为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket 连接升级失败", http.StatusInternalServerError)
		return
	}
	
	// 注册连接到 hub
	hub.Register(conn, roomID, sessionID)
}

// parseTokenSessionID 从令牌中解析会话ID
func parseTokenSessionID(token string) string {
	// 简化实现：从令牌中提取 session_id
	// 实际应该使用 JWT 解析
	
	if len(token) < 20 {
		return ""
	}
	
	// 临时实现：假设 token 格式为 "xiaowo_roomid_sessionid_uuid"
	parts := splitToken(token)
	if len(parts) >= 3 {
		return parts[2]
	}
	
	return ""
}

// splitToken 分割令牌
func splitToken(token string) []string {
	// 简化实现：按下划线分割
	var parts []string
	current := ""
	
	for _, char := range token {
		if char == '_' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	return parts
}