package websocket

import (
	"encoding/json"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketHub WebSocket连接管理中心
type WebSocketHub struct {
	rooms     map[string]*Room
	clients   map[string]*WebSocketConnection
	broadcast chan []byte
	register  chan *WebSocketConnection
	unregister chan *WebSocketConnection
	mu        sync.RWMutex
}

// Room 房间连接管理
type Room struct {
	ID        string
	clients   map[string]*WebSocketConnection
	version   int64 // 乐观锁版本
	state     PlaybackState
	mu        sync.RWMutex
}

// PlaybackState 播放状态
type PlaybackState struct {
	CurrentTime  float64 `json:"current_time"`
	Duration     float64 `json:"duration"`
	IsPlaying    bool    `json:"is_playing"`
	PlaybackRate float64 `json:"playback_rate"`
	VideoURL     string  `json:"video_url"`
	VideoTitle   string  `json:"video_title"`
	LastUpdated  int64   `json:"last_updated"`
}

// Message 基础消息类型
type Message struct {
	Type string `json:"type"`
}

// 消息类型常量
const (
	MsgTypePing    = "ping"
	MsgTypePong    = "pong"
	MsgTypeAuth    = "auth"
	MsgTypeSync    = "sync"
	MsgTypeChat    = "chat"
	MsgTypePlay    = "play"
	MsgTypePause   = "pause"
	MsgTypeSeek    = "seek"
	MsgTypeRate    = "rate"
	MsgTypeError   = "error"
	MsgTypeHeartbeat = "heartbeat"
)

// PingMessage ping 消息
type PingMessage struct {
	Type    string `json:"type"`              // "ping"
	Purpose string `json:"purpose,omitempty"` // "heartbeat" | "calibration"
	ClientSendTime int64 `json:"client_send_time"` // 客户端发送时间戳
}

// PongMessage pong 消息
type PongMessage struct {
	Type           string `json:"type"` // "pong"
	ClientSendTime int64  `json:"client_send_time"`
	ServerRecvTime int64  `json:"server_recv_time"`
	ServerSendTime int64  `json:"server_send_time"`
}

// AuthMessage 认证消息
type AuthMessage struct {
	Type     string `json:"type"`     // "auth"
	Token    string `json:"token"`    // 访问令牌
	RoomID   string `json:"room_id"`  // 房间ID
}

// SyncMessage 同步消息
type SyncMessage struct {
	Type    string `json:"type"` // "sync"
	RoomID  string `json:"room_id"`
	Data    SyncData `json:"data"`
}

// SyncData 同步数据
type SyncData struct {
	CurrentTime  float64 `json:"current_time"`  // 客户端当前播放时间
	Duration     float64 `json:"duration"`      // 视频总时长
	IsPlaying    bool    `json:"is_playing"`    // 是否正在播放
	PlaybackRate float64 `json:"playback_rate"` // 播放倍速
	BaseVersion  int64   `json:"base_version"`  // 基础版本号（乐观锁）
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Type        string    `json:"type"`         // "chat"
	RoomID      string    `json:"room_id"`      // 房间ID
	SessionID   string    `json:"session_id"`   // 发送者会话ID
	DisplayName string    `json:"display_name"` // 发送者显示名称
	Message     string    `json:"message"`      // 消息内容
	Timestamp   int64     `json:"timestamp"`    // 发送时间戳
}

// PlayMessage 播放消息
type PlayMessage struct {
	Type    string `json:"type"`    // "play"
	RoomID  string `json:"room_id"` // 房间ID
	StartTime float64 `json:"start_time"` // 开始播放时间
}

// PauseMessage 暂停消息
type PauseMessage struct {
	Type    string `json:"type"`    // "pause"
	RoomID  string `json:"room_id"` // 房间ID
	PauseTime float64 `json:"pause_time"` // 暂停时间
}

// SeekMessage 跳转消息
type SeekMessage struct {
	Type    string `json:"type"`    // "seek"
	RoomID  string `json:"room_id"` // 房间ID
	TargetTime float64 `json:"target_time"` // 目标播放时间
}

// RateMessage 倍速消息
type RateMessage struct {
	Type    string `json:"type"`    // "rate"
	RoomID  string `json:"room_id"` // 房间ID
	PlaybackRate   float64 `json:"playback_rate"` // 播放倍速
}

// WebSocketConnection WebSocket连接
type WebSocketConnection struct {
	ws        *websocket.Conn
	roomID    string
	sessionID string
	send      chan []byte
	// RTT 相关
	rtts           []int64 // 最近3次RTT测量
	lastCalibrate  time.Time
	timeOffset     int64 // 时钟偏移量
	mu             sync.RWMutex
}

// Register 注册WebSocket连接
func (h *WebSocketHub) Register(conn *websocket.Conn, roomID, sessionID string) {
	// 创建WebSocket连接对象
	wsConn := &WebSocketConnection{
		ws:        conn,
		roomID:    roomID,
		sessionID: sessionID,
		send:      make(chan []byte, 256),
	}
	
	// 注册连接
	h.register <- wsConn
	
	// 启动读写协程
	go wsConn.readPump(h)
	go wsConn.writePump()
}

// readPump 从WebSocket连接读取消息
func (c *WebSocketConnection) readPump(hub *WebSocketHub) {
	defer func() {
		hub.unregister <- c
		c.ws.Close()
	}()
	
	c.ws.SetReadLimit(512)
	c.ws.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(60 * time.Second)); return nil })
	
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 记录错误
			}
			break
		}
		
		// 处理消息
		hub.HandleMessage(c, message)
	}
}

// writePump 向WebSocket连接写入消息
func (c *WebSocketConnection) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	
	for {
		select {
		case message, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// hub 关闭了通道
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			
			w, err := c.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			
			// 添加队列中的其他消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}
			
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SyncAction 同步动作
type SyncAction struct {
	Type          string  `json:"type"`          // "seek", "playback_rate", "play", "pause"
	TargetTime    float64 `json:"target_time,omitempty"`    // 目标时间点
	PlaybackRate  float64 `json:"playback_rate,omitempty"`  // 播放倍速
	Reason        string  `json:"reason"`        // 同步原因
}

// NewWebSocketHub 创建WebSocket Hub
func NewWebSocketHub() *WebSocketHub {
	return &WebSocketHub{
		rooms:       make(map[string]*Room),
		clients:     make(map[string]*WebSocketConnection),
		broadcast:   make(chan []byte),
		register:    make(chan *WebSocketConnection),
		unregister:  make(chan *WebSocketConnection),
	}
}

// Run 运行WebSocket Hub主循环
func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToAll(message)
		}
	}
}

// RegisterClient 注册客户端连接
func (h *WebSocketHub) RegisterClient(conn *WebSocketConnection) {
	h.register <- conn
}

// UnregisterClient 取消注册客户端连接
func (h *WebSocketHub) UnregisterClient(conn *WebSocketConnection) {
	h.unregister <- conn
}

// registerClient 注册客户端连接（内部方法）
func (h *WebSocketHub) registerClient(conn *WebSocketConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[conn.sessionID] = conn

	// 加入房间
	room := h.getOrCreateRoom(conn.roomID)
	room.mu.Lock()
	room.clients[conn.sessionID] = conn
	room.mu.Unlock()

	// 发送房间当前状态
	h.sendRoomState(conn.roomID, conn)
	
	// 广播成员加入通知给房间内其他成员
	h.broadcastToRoom(conn.roomID, map[string]interface{}{
		"type":        "member_join",
		"session_id":  conn.sessionID,
		"room_id":     conn.roomID,
		"timestamp":   time.Now().Unix(),
		"member_count": len(room.clients),
	})
}

// unregisterClient 取消注册客户端连接（内部方法）
func (h *WebSocketHub) unregisterClient(conn *WebSocketConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.clients[conn.sessionID]; ok {
		delete(h.clients, conn.sessionID)
		close(conn.send)

		// 从房间移除
		if room, ok := h.rooms[conn.roomID]; ok {
			room.mu.Lock()
			memberCount := len(room.clients)
			delete(room.clients, conn.sessionID)
			room.mu.Unlock()

			// 广播成员退出通知给房间内其他成员
			h.broadcastToRoom(conn.roomID, map[string]interface{}{
				"type":        "member_leave",
				"session_id":  conn.sessionID,
				"room_id":     conn.roomID,
				"timestamp":   time.Now().Unix(),
				"member_count": memberCount - 1,
			})

			// 如果房间为空，清理房间
			if len(room.clients) == 0 {
				delete(h.rooms, conn.roomID)
			}
		}
	}
}

// getOrCreateRoom 获取或创建房间
func (h *WebSocketHub) getOrCreateRoom(roomID string) *Room {
	room, exists := h.rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			clients: make(map[string]*WebSocketConnection),
			version: 0,
			state: PlaybackState{
				CurrentTime:  0,
				Duration:     0,
				IsPlaying:    false,
				PlaybackRate: 1.0,
				LastUpdated:  time.Now().Unix(),
			},
		}
		h.rooms[roomID] = room
	}
	return room
}

// HandleMessage 处理WebSocket消息
func (h *WebSocketHub) HandleMessage(conn *WebSocketConnection, message []byte) {
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		h.sendError(conn, "invalid_message", "消息格式错误")
		return
	}

	switch msg.Type {
	case MsgTypePing:
		h.handlePing(conn, message)
	case MsgTypePong:
		h.handlePong(conn, message)
	case MsgTypeAuth:
		h.handleAuth(conn, message)
	case MsgTypeSync:
		h.handleSync(conn, message)
	case MsgTypeChat:
		h.handleChat(conn, message)
	case MsgTypePlay:
		h.handlePlay(conn, message)
	case MsgTypePause:
		h.handlePause(conn, message)
	case MsgTypeSeek:
		h.handleSeek(conn, message)
	case MsgTypeRate:
		h.handleRate(conn, message)
	default:
		h.sendError(conn, "unknown_message_type", "未知消息类型")
	}
}

// handlePing 处理ping消息
func (h *WebSocketHub) handlePing(conn *WebSocketConnection, message []byte) {
	var pingMsg PingMessage
	if err := json.Unmarshal(message, &pingMsg); err != nil {
		return
	}

	// 根据purpose决定是心跳还是对时
	if pingMsg.Purpose == "calibration" {
		// 对时消息，计算RTT和时钟偏移
		h.handleCalibration(conn, pingMsg)
	} else {
		// 心跳消息，只返回pong
		h.sendPong(conn, pingMsg.ClientSendTime)
	}
}

// handleCalibration 处理对时消息
func (h *WebSocketHub) handleCalibration(conn *WebSocketConnection, pingMsg PingMessage) {
	serverRecvTime := time.Now().UnixMilli()
	
	// 发送pong响应
	pongMsg := PongMessage{
		Type:           MsgTypePong,
		ClientSendTime: pingMsg.ClientSendTime,
		ServerRecvTime: serverRecvTime,
		ServerSendTime: time.Now().UnixMilli(),
	}
	
	// 计算RTT
	rawRTT := pongMsg.ServerSendTime - pingMsg.ClientSendTime
	
	// RTT平滑处理 (中位数算法)
	smoothedRTT := conn.calculateSmoothedRTT(rawRTT)
	
	// 计算时钟偏移
	offset := pongMsg.ServerRecvTime - pingMsg.ClientSendTime - (smoothedRTT / 2)
	conn.SetTimeOffset(offset)
	
	// 记录校准时间
	conn.lastCalibrate = time.Now()
	
	// 发送pong响应
	h.sendJSON(conn, pongMsg)
}

// calculateSmoothedRTT 计算平滑RTT
func (c *WebSocketConnection) calculateSmoothedRTT(newRTT int64) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	
	c.rtts = append(c.rtts, newRTT)
	if len(c.rtts) > 3 {
		c.rtts = c.rtts[1:] // 保持最近3次
	}
	if len(c.rtts) < 3 {
		return newRTT // 不足3次直接返回
	}
	
	// 使用中位数平滑处理
	sortedRTTs := make([]int64, len(c.rtts))
	copy(sortedRTTs, c.rtts)
	sort.Slice(sortedRTTs, func(i, j int) bool { return sortedRTTs[i] < sortedRTTs[j] })
	return sortedRTTs[len(sortedRTTs)/2] // 中位数
}

// SetTimeOffset 设置时钟偏移
func (c *WebSocketConnection) SetTimeOffset(offset int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.timeOffset = offset
}

// GetAdjustedTime 获取修正时间
func (c *WebSocketConnection) GetAdjustedTime(clientTime int64) int64 {
	return clientTime + c.timeOffset
}

// handlePong 处理pong消息
func (h *WebSocketHub) handlePong(conn *WebSocketConnection, message []byte) {
	// 这里主要用于心跳响应，不需要特殊处理
}

// handleAuth 处理认证消息
func (h *WebSocketHub) handleAuth(conn *WebSocketConnection, message []byte) {
	var authMsg AuthMessage
	if err := json.Unmarshal(message, &authMsg); err != nil {
		h.sendError(conn, "auth_failed", "认证消息格式错误")
		return
	}
	
	// TODO: 验证JWT令牌
	// 暂时简化实现
	
	// 认证成功，发送确认
	successMsg := map[string]interface{}{
		"type":    "auth_success",
		"room_id": authMsg.RoomID,
		"status":  "authenticated",
	}
	h.sendJSON(conn, successMsg)
}

// handleSync 处理同步消息
func (h *WebSocketHub) handleSync(conn *WebSocketConnection, message []byte) {
	var syncMsg SyncMessage
	if err := json.Unmarshal(message, &syncMsg); err != nil {
		h.sendError(conn, "sync_failed", "同步消息格式错误")
		return
	}

	room := h.getOrCreateRoom(syncMsg.RoomID)
	room.mu.Lock()
	defer room.mu.Unlock()

	currentTime := syncMsg.Data.CurrentTime
	targetTime := h.calculateTargetTime(room)
	timeDiff := targetTime - currentTime

	var action SyncAction
	switch {
	case timeDiff > 2.0: // Level 2: 误差 > 2s，触发 Seek
		action = SyncAction{
			Type:       "seek",
			TargetTime: targetTime,
			Reason:     "large_difference",
		}
	case timeDiff > 0.5: // Level 3: 误差 0.5s - 2s，动态调整倍速
		speed := 1.0 + math.Min(timeDiff/10.0, 0.5) // 最大1.5倍速
		action = SyncAction{
			Type:          "playback_rate",
			PlaybackRate:  speed,
			TargetTime:    targetTime,
			Reason:        "medium_difference",
		}
	case timeDiff < -0.5: // 客户端超前，减速
		speed := 1.0 + math.Max(timeDiff/10.0, -0.3) // 最小0.7倍速
		action = SyncAction{
			Type:          "playback_rate",
			PlaybackRate:  speed,
			TargetTime:    targetTime,
			Reason:        "client_ahead",
		}
	default: // Level 4: 误差 < 0.5s，忽略
		return
	}

	// 广播同步指令给房间内其他用户
	h.broadcastToRoom(room.ID, action)
}

// calculateTargetTime 计算目标时间
func (h *WebSocketHub) calculateTargetTime(room *Room) float64 {
	// 简化实现：返回房间当前播放时间
	return room.state.CurrentTime
}

// handleChat 处理聊天消息
func (h *WebSocketHub) handleChat(conn *WebSocketConnection, message []byte) {
	var chatMsg ChatMessage
	if err := json.Unmarshal(message, &chatMsg); err != nil {
		h.sendError(conn, "chat_failed", "聊天消息格式错误")
		return
	}

	// 广播聊天消息给房间内所有用户
	h.broadcastToRoom(conn.roomID, chatMsg)
}

// handlePlay 处理播放消息
func (h *WebSocketHub) handlePlay(conn *WebSocketConnection, message []byte) {
	room := h.getOrCreateRoom(conn.roomID)
	room.mu.Lock()
	defer room.mu.Unlock()

	room.state.IsPlaying = true
	room.state.LastUpdated = time.Now().Unix()
	room.version++

	// 广播播放状态给房间内其他用户
	h.broadcastToRoom(conn.roomID, map[string]interface{}{
		"type":         "play",
		"current_time": room.state.CurrentTime,
		"version":      room.version,
	})
}

// handlePause 处理暂停消息
func (h *WebSocketHub) handlePause(conn *WebSocketConnection, message []byte) {
	room := h.getOrCreateRoom(conn.roomID)
	room.mu.Lock()
	defer room.mu.Unlock()

	room.state.IsPlaying = false
	room.state.LastUpdated = time.Now().Unix()
	room.version++

	// 广播暂停状态给房间内其他用户
	h.broadcastToRoom(conn.roomID, map[string]interface{}{
		"type":         "pause",
		"current_time": room.state.CurrentTime,
		"version":      room.version,
	})
}

// handleSeek 处理拖拽消息
func (h *WebSocketHub) handleSeek(conn *WebSocketConnection, message []byte) {
	var seekMsg SeekMessage
	if err := json.Unmarshal(message, &seekMsg); err != nil {
		h.sendError(conn, "seek_failed", "拖拽消息格式错误")
		return
	}

	room := h.getOrCreateRoom(conn.roomID)
	room.mu.Lock()
	defer room.mu.Unlock()

	room.state.CurrentTime = seekMsg.TargetTime
	room.state.LastUpdated = time.Now().Unix()
	room.version++

	// 广播拖拽状态给房间内其他用户
	h.broadcastToRoom(conn.roomID, map[string]interface{}{
		"type":         "seek",
		"target_time":  seekMsg.TargetTime,
		"version":      room.version,
	})
}

// handleRate 处理倍速消息
func (h *WebSocketHub) handleRate(conn *WebSocketConnection, message []byte) {
	var rateMsg RateMessage
	if err := json.Unmarshal(message, &rateMsg); err != nil {
		h.sendError(conn, "rate_failed", "倍速消息格式错误")
		return
	}

	room := h.getOrCreateRoom(conn.roomID)
	room.mu.Lock()
	defer room.mu.Unlock()

	room.state.PlaybackRate = rateMsg.PlaybackRate
	room.state.LastUpdated = time.Now().Unix()
	room.version++

	// 广播倍速状态给房间内其他用户
	h.broadcastToRoom(conn.roomID, map[string]interface{}{
		"type":          "rate",
		"playback_rate": rateMsg.PlaybackRate,
		"version":       room.version,
	})
}

// sendRoomState 发送房间状态
func (h *WebSocketHub) sendRoomState(roomID string, conn *WebSocketConnection) {
	room := h.getOrCreateRoom(roomID)
	room.mu.RLock()
	defer room.mu.RUnlock()

	stateMsg := map[string]interface{}{
		"type":    "room_state",
		"state":   room.state,
		"version": room.version,
		"members": len(room.clients),
	}
	h.sendJSON(conn, stateMsg)
}

// broadcastToRoom 向房间内广播消息
func (h *WebSocketHub) broadcastToRoom(roomID string, data interface{}) {
	room := h.getOrCreateRoom(roomID)
	room.mu.RLock()
	defer room.mu.RUnlock()

	message, _ := json.Marshal(data)
	for _, conn := range room.clients {
		select {
		case conn.send <- message:
		default:
			close(conn.send)
			h.unregister <- conn
		}
	}
}

// broadcastToAll 向所有连接广播消息
func (h *WebSocketHub) broadcastToAll(message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, conn := range h.clients {
		select {
		case conn.send <- message:
		default:
			close(conn.send)
			h.unregister <- conn
		}
	}
}

// sendJSON 发送JSON消息
func (h *WebSocketHub) sendJSON(conn *WebSocketConnection, data interface{}) {
	message, _ := json.Marshal(data)
	select {
	case conn.send <- message:
	default:
		close(conn.send)
		h.unregister <- conn
	}
}

// sendPong 发送pong响应
func (h *WebSocketHub) sendPong(conn *WebSocketConnection, clientSendTime int64) {
	pongMsg := PongMessage{
		Type:           MsgTypePong,
		ClientSendTime: clientSendTime,
		ServerRecvTime: time.Now().UnixMilli(),
		ServerSendTime: time.Now().UnixMilli(),
	}
	h.sendJSON(conn, pongMsg)
}

// sendError 发送错误消息
func (h *WebSocketHub) sendError(conn *WebSocketConnection, code, message string) {
	errorMsg := map[string]interface{}{
		"type":    MsgTypeError,
		"code":    code,
		"message": message,
	}
	h.sendJSON(conn, errorMsg)
}

// StartPeriodicTasks 启动周期性任务
func (h *WebSocketHub) StartPeriodicTasks() {
	// 心跳任务 - 每30秒
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			h.broadcastHeartbeat()
		}
	}()

	// 对时任务 - 每5分钟
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			h.triggerCalibration()
		}
	}()
}

// broadcastHeartbeat 广播心跳
func (h *WebSocketHub) broadcastHeartbeat() {
	heartbeatMsg := map[string]interface{}{
		"type": MsgTypeHeartbeat,
	}
	
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	message, _ := json.Marshal(heartbeatMsg)
	for _, conn := range h.clients {
		select {
		case conn.send <- message:
		default:
			close(conn.send)
			h.unregister <- conn
		}
	}
}

// triggerCalibration 触发对时
func (h *WebSocketHub) triggerCalibration() {
	h.mu.RLock()
	defer h.mu.RUnlock()
	
	calibrationMsg := map[string]interface{}{
		"type":    MsgTypePing,
		"purpose": "calibration",
	}
	
	message, _ := json.Marshal(calibrationMsg)
	for _, conn := range h.clients {
		select {
		case conn.send <- message:
		default:
			close(conn.send)
			h.unregister <- conn
		}
	}
}