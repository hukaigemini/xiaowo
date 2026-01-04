package v1

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xiaowo/backend/internal/model"
	"xiaowo/backend/internal/service"
)

// RoomHandler 房间相关API处理器
type RoomHandler struct {
	roomService   *service.RoomService
	memberService *service.MemberService
}

// NewRoomHandler 创建房间处理器
func NewRoomHandler(roomService *service.RoomService, memberService *service.MemberService) *RoomHandler {
	return &RoomHandler{
		roomService:   roomService,
		memberService: memberService,
	}
}

// CreateRoom 创建房间
// @Summary 创建房间
// @Description 创建一个新的同步播放房间
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body CreateRoomRequest true "创建房间请求"
// @Success 200 {object} RoomResponse
// @Router /api/v1/rooms [post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
	var req CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}

	// 验证请求参数
	if err := h.validateCreateRoomRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 生成匿名用户会话ID
	sessionID := generateSessionID()

	// 转换为服务层的CreateRoomRequest
	serviceReq := &service.CreateRoomRequest{
		Name:          req.Name,
		Description:   req.Description,
		IsPrivate:     &req.IsPrivate,
		Password:      req.Password,
		MaxUsers:      &req.MaxUsers,
		MediaURL:      req.MediaURL,
		MediaType:     req.MediaType,
		MediaTitle:    req.MediaTitle,
		MediaDuration: int(req.MediaDuration),
		Settings:      nil,
	}

	// 创建房间
	room, err := h.roomService.CreateRoom(serviceReq, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建房间失败",
			"detail": err.Error(),
		})
		return
	}

	// 自动创建者为房间成员
	member := &model.RoomMember{
		RoomID:      room.ID,
		SessionID:   sessionID,
		Nickname:    generateDisplayName(),
		JoinedAt:    time.Now(),
		LastSeen:    time.Now(),
	}

	if err := h.memberService.AddMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "添加房间成员失败",
			"detail": err.Error(),
		})
		return
	}

	// 生成访问令牌
	// token, err := generateRoomToken(room.ID, sessionID)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "生成访问令牌失败",
	// 		"detail": err.Error(),
	// 	})
	// 	return
	// }

	// 返回创建结果
	resp := &RoomResponse{
		Room:        room,
		MemberCount: 1, // 创建者自动加入
		CreatedBy:   room.CreatorSessionID,
		CreatedAt:   room.CreatedAt,
		IsCreator:   true, // 当前用户是创建者
	}

	c.JSON(http.StatusCreated, resp)
}

// GetRoom 获取房间信息
// @Summary 获取房间信息
// @Description 根据房间ID获取房间详细信息
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Success 200 {object} RoomDetailResponse
// @Router /api/v1/rooms/{room_id} [get]
func (h *RoomHandler) GetRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	
	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在",
			"detail": err.Error(),
		})
		return
	}

	// 获取房间成员数量
	memberCount, err := h.memberService.GetMemberCount(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取房间成员数量失败",
			"detail": err.Error(),
		})
		return
	}

	resp := &RoomDetailResponse{
		Room:        room,
		MemberCount: memberCount,
	}

	c.JSON(http.StatusOK, resp)
}

// JoinRoom 加入房间
// @Summary 加入房间
// @Description 加入一个现有的同步播放房间
// @Tags rooms
// @Accept json
// @Produce json
// @Param room_id path string true "房间ID"
// @Param request body JoinRoomRequest true "加入房间请求"
// @Success 200 {object} JoinRoomResponse
// @Router /api/v1/rooms/{room_id}/join [post]
func (h *RoomHandler) JoinRoom(c *gin.Context) {
	// 从URL参数获取roomID
	roomID := c.Param("room_id")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "房间ID不能为空",
		})
		return
	}

	var req JoinRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}

	// 验证房间是否存在
	room, err := h.roomService.GetRoom(roomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "房间不存在",
			"detail": err.Error(),
		})
		return
	}

	// 验证房间密码（如果需要）
	if room.IsPrivate && req.Password != room.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "房间密码错误",
		})
		return
	}

	// 检查房间人数限制
	memberCount, err := h.memberService.GetMemberCount(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "检查房间人数失败",
			"detail": err.Error(),
		})
		return
	}

	if memberCount >= room.MaxUsers {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "房间已满",
		})
		return
	}

	// 生成会话ID和显示名称
	sessionID := generateSessionID()
	displayName := req.DisplayName
	if displayName == "" {
		displayName = generateDisplayName()
	}

	// 添加房间成员
	member := &model.RoomMember{
		RoomID:    roomID,
		SessionID: sessionID,
		Nickname:  displayName,
		JoinedAt:  time.Now(),
		LastSeen:  time.Now(),
	}

	if err := h.memberService.AddMember(member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "加入房间失败",
			"detail": err.Error(),
		})
		return
	}

	// 生成访问令牌
	token, err := generateRoomToken(roomID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成访问令牌失败",
			"detail": err.Error(),
		})
		return
	}

	resp := &JoinRoomResponse{
		Room:      room,
		SessionID: sessionID,
		Token:     token,
		JoinURL:   generateJoinURL(roomID, token),
	}

	c.JSON(http.StatusOK, resp)
}

// CloseRoom 关闭房间
// @Summary 关闭房间
// @Description 房间创建者关闭房间
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Param session_id query string true "会话ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/rooms/{room_id} [delete]
func (h *RoomHandler) CloseRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	sessionID := c.Query("session_id")
	
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	// 验证权限（只有房间创建者可以关闭房间）
	if !h.roomService.IsCreator(roomID, sessionID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "只有房间创建者可以关闭房间",
		})
		return
	}

	// 关闭房间
	if err := h.roomService.CloseRoom(roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "关闭房间失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "成功关闭房间",
	})
}

// LeaveRoom 离开房间
// @Summary 离开房间
// @Description 用户离开房间
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Param session_id query string true "会话ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/rooms/{room_id}/leave [post]
func (h *RoomHandler) LeaveRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	sessionID := c.Query("session_id")
	
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "会话ID不能为空",
		})
		return
	}

	// 移除房间成员
	if err := h.memberService.RemoveMember(roomID, sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "离开房间失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "成功离开房间",
	})
}

// GetRoomMembers 获取房间成员列表
// @Summary 获取房间成员列表
// @Description 获取房间的所有成员信息
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Success 200 {array} model.RoomMember
// @Router /api/v1/rooms/{room_id}/members [get]
func (h *RoomHandler) GetRoomMembers(c *gin.Context) {
	roomID := c.Param("room_id")
	
	// 获取房间成员列表
	members, err := h.memberService.GetRoomMembers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取房间成员列表失败",
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, members)
}

// UpdateRoom 更新房间信息
// @Summary 更新房间信息
// @Description 房间创建者可以更新房间信息
// @Tags rooms
// @Accept json
// @Produce json
// @Param room_id path string true "房间ID"
// @Param session_id query string true "会话ID"
// @Param request body UpdateRoomRequest true "更新房间请求"
// @Success 200 {object} RoomResponse
// @Router /api/v1/rooms/{room_id} [put]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomID := c.Param("room_id")
	sessionID := c.Query("session_id")
	
	var req UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}

	// 验证权限（只有房间创建者可以更新）
	if !h.roomService.IsCreator(roomID, sessionID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "只有房间创建者可以更新房间信息",
		})
		return
	}

	// 转换为服务层的UpdateRoomRequest
	serviceReq := &service.UpdateRoomRequest{
		Settings: nil,
	}
	if req.Name != "" {
		serviceReq.Name = &req.Name
	}
	if req.IsPrivate != nil {
		serviceReq.IsPrivate = req.IsPrivate
	}
	if req.Password != "" {
		serviceReq.Password = &req.Password
	}
	if req.MaxUsers != nil {
		serviceReq.MaxUsers = req.MaxUsers
	}

	// 更新房间
	room, err := h.roomService.UpdateRoom(roomID, serviceReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "更新房间失败",
			"detail": err.Error(),
		})
		return
	}

	resp := &RoomResponse{
		Room:        room,
		MemberCount: 0, // 后续可以从memberService获取
		CreatedBy:   room.CreatorSessionID,
		CreatedAt:   room.CreatedAt,
		IsCreator:   true, // 当前用户是创建者
	}

	c.JSON(http.StatusOK, resp)
}

// PlayVideo 播放视频
// @Summary 播放视频
// @Description 播放房间内的视频
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/rooms/{room_id}/play [post]
func (h *RoomHandler) PlayVideo(c *gin.Context) {
	roomID := c.Param("room_id")
	
	// 播放视频
	if err := h.roomService.PlayVideo(roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "播放视频失败",
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "视频播放成功",
	})
}

// PauseVideo 暂停视频
// @Summary 暂停视频
// @Description 暂停房间内的视频
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/rooms/{room_id}/pause [post]
func (h *RoomHandler) PauseVideo(c *gin.Context) {
	roomID := c.Param("room_id")
	
	// 暂停视频
	if err := h.roomService.PauseVideo(roomID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "暂停视频失败",
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "视频暂停成功",
	})
}

// SeekVideo 跳转视频
// @Summary 跳转视频
// @Description 跳转房间内的视频到指定时间点
// @Tags rooms
// @Accept json
// @Produce json
// @Param room_id path string true "房间ID"
// @Param request body struct{CurrentTime float64 `json:"current_time"`} true "跳转请求"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/rooms/{room_id}/seek [post]
func (h *RoomHandler) SeekVideo(c *gin.Context) {
	roomID := c.Param("room_id")
	
	var req struct {
		CurrentTime float64 `json:"current_time" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}
	
	// 跳转视频
	if err := h.roomService.SeekVideo(roomID, req.CurrentTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "跳转视频失败",
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "视频跳转成功",
	})
}

// GetPlaybackStatus 获取播放状态
// @Summary 获取播放状态
// @Description 获取房间内视频的当前播放状态
// @Tags rooms
// @Produce json
// @Param room_id path string true "房间ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/rooms/{room_id}/status [get]
func (h *RoomHandler) GetPlaybackStatus(c *gin.Context) {
	roomID := c.Param("room_id")
	
	// 获取播放状态
	status, err := h.roomService.GetPlaybackStatus(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取播放状态失败",
			"detail": err.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, status)
}

// ListRooms 获取房间列表
// @Summary 获取房间列表
// @Description 获取公开房间列表
// @Tags rooms
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Success 200 {object} RoomListResponse
// @Router /api/v1/rooms [get]
func (h *RoomHandler) ListRooms(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	
	// 限制每页最大数量
	if size > 100 {
		size = 100
	}

	rooms, total, err := h.roomService.ListRooms(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取房间列表失败",
			"detail": err.Error(),
		})
		return
	}

	resp := &RoomListResponse{
		Rooms: rooms,
		Total: total,
		Page:  page,
		Size:  size,
	}

	c.JSON(http.StatusOK, resp)
}

// validateCreateRoomRequest 验证创建房间请求
func (h *RoomHandler) validateCreateRoomRequest(req *CreateRoomRequest) error {
	if req.Name == "" {
		return errors.New("房间名称不能为空")
	}
	
	if len(req.Name) > 100 {
		return errors.New("房间名称过长")
	}
	
	if req.MaxUsers <= 0 || req.MaxUsers > 1000 {
		return errors.New("房间人数限制必须在1-1000之间")
	}
	
	if req.IsPrivate && req.Password == "" {
		return errors.New("私密房间必须设置密码")
	}
	
	if len(req.Password) > 50 {
		return errors.New("密码过长")
	}

	return nil
}