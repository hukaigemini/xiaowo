package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xiaowo/backend/internal/service"
)

// SessionHandler 会话相关API处理器
type SessionHandler struct {
	sessionService *service.SessionService
}

// NewSessionHandler 创建会话处理器
func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// CreateSession 创建新会话
// @Summary 创建新会话
// @Description 创建匿名用户会话
// @Tags sessions
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "创建会话请求"
// @Success 201 {object} SessionResponse
// @Router /api/v1/sessions [post]
func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}

	// 创建会话
	session, err := h.sessionService.CreateSession(req.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建会话失败",
			"detail": err.Error(),
		})
		return
	}

	// 处理RoomID为*string的情况
	var roomID string
	if session.RoomID != nil {
		roomID = *session.RoomID
	}
	
	resp := &SessionResponse{
		SessionID:  session.ID,
		Nickname:   session.Nickname,
		Avatar:     session.Avatar,
		RoomID:     roomID,
		Status:     string(session.Status),
		CreatedAt:  session.CreatedAt,
		LastSeenAt: session.LastSeenAt,
		ExpiresAt:  session.ExpiresAt,
		IsExpired:  session.IsExpired(),
		IsOnline:   session.IsOnline(),
		IsActive:   session.IsActive(),
	}

	c.JSON(http.StatusCreated, resp)
}

// GetSession 获取会话信息
// @Summary 获取会话信息
// @Description 根据会话ID获取会话信息
// @Tags sessions
// @Produce json
// @Param session_id path string true "会话ID"
// @Success 200 {object} SessionResponse
// @Router /api/v1/sessions/{session_id} [get]
func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "会话不存在",
			"detail": err.Error(),
		})
		return
	}

	// 处理RoomID为*string的情况
	var roomID string
	if session.RoomID != nil {
		roomID = *session.RoomID
	}
	
	resp := &SessionResponse{
		SessionID:  session.ID,
		Nickname:   session.Nickname,
		Avatar:     session.Avatar,
		RoomID:     roomID,
		Status:     string(session.Status),
		CreatedAt:  session.CreatedAt,
		LastSeenAt: session.LastSeenAt,
		ExpiresAt:  session.ExpiresAt,
		IsExpired:  session.IsExpired(),
		IsOnline:   session.IsOnline(),
		IsActive:   session.IsActive(),
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateSession 更新会话信息
// @Summary 更新会话信息
// @Description 更新用户昵称等信息
// @Tags sessions
// @Accept json
// @Produce json
// @Param session_id path string true "会话ID"
// @Param request body UpdateSessionRequest true "更新会话请求"
// @Success 200 {object} SessionResponse
// @Router /api/v1/sessions/{session_id} [put]
func (h *SessionHandler) UpdateSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	var req UpdateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
			"detail": err.Error(),
		})
		return
	}

	// 更新会话
	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	session, err := h.sessionService.UpdateSession(sessionID, updates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "更新会话失败",
			"detail": err.Error(),
		})
		return
	}

	// 处理RoomID为*string的情况
	var roomID string
	if session.RoomID != nil {
		roomID = *session.RoomID
	}
	
	resp := &SessionResponse{
		SessionID:  session.ID,
		Nickname:   session.Nickname,
		Avatar:     session.Avatar,
		RoomID:     roomID,
		Status:     string(session.Status),
		CreatedAt:  session.CreatedAt,
		LastSeenAt: session.LastSeenAt,
		ExpiresAt:  session.ExpiresAt,
		IsExpired:  session.IsExpired(),
		IsOnline:   session.IsOnline(),
		IsActive:   session.IsActive(),
	}

	c.JSON(http.StatusOK, resp)
}

// Heartbeat 心跳保活
// @Summary 心跳保活
// @Description 更新会话最后在线时间并设置状态为在线
// @Tags sessions
// @Produce json
// @Param session_id path string true "会话ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/sessions/{session_id}/heartbeat [post]
func (h *SessionHandler) Heartbeat(c *gin.Context) {
	sessionID := c.Param("session_id")

	if err := h.sessionService.Heartbeat(sessionID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "心跳更新失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "心跳更新成功",
	})
}

// ValidateSession 验证会话
// @Summary 验证会话
// @Description 验证会话是否有效
// @Tags sessions
// @Produce json
// @Param session_id path string true "会话ID"
// @Success 200 {object} SessionValidationResponse
// @Router /api/v1/sessions/{session_id}/validate [get]
func (h *SessionHandler) ValidateSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := h.sessionService.GetSession(sessionID)
	if err != nil {
		c.JSON(http.StatusOK, SessionValidationResponse{
			SessionID: sessionID,
			IsValid:   false,
			Message:   "会话不存在",
		})
		return
	}

	isValid := !session.IsExpired()
	message := "会话有效"
	if !isValid {
		message = "会话已过期"
	}

	c.JSON(http.StatusOK, SessionValidationResponse{
		SessionID: sessionID,
		IsValid:   isValid,
		Message:   message,
		ExpiresAt: session.ExpiresAt,
	})
}

// DeleteSession 删除会话
// @Summary 删除会话
// @Description 删除用户会话
// @Tags sessions
// @Produce json
// @Param session_id path string true "会话ID"
// @Success 200 {object} SuccessResponse
// @Router /api/v1/sessions/{session_id} [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID := c.Param("session_id")

	if err := h.sessionService.DeleteSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "删除会话失败",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "会话删除成功",
	})
}