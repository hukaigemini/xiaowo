package v1

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ==================== 中间件 ====================

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "缺少访问令牌",
			})
			c.Abort()
			return
		}

		// TODO: 验证JWT令牌
		// 暂时简化实现，检查令牌格式
		if len(token) < 10 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的访问令牌",
			})
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

// CORSMiddleware CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// ==================== 错误处理 ====================

// HandleError 全局错误处理
func HandleError() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "内部服务器错误",
					"detail": "系统发生异常",
				})
			}
		}()
		c.Next()
	}
}

// ==================== 健康检查 ====================

// HealthCheck 健康检查处理器
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// @Summary 健康检查
// @Description 检查服务健康状态
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now().Unix(),
		"version":   "1.0.0",
	})
}

// @Summary 就绪检查
// @Description 检查服务就绪状态
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /ready [get]
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// TODO: 检查数据库连接等依赖服务
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"timestamp": time.Now().Unix(),
	})
}

// ==================== API版本信息 ====================

// VersionHandler 版本信息处理器
type VersionHandler struct{}

func NewVersionHandler() *VersionHandler {
	return &VersionHandler{}
}

// @Summary API版本信息
// @Description 获取API版本信息
// @Tags version
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /version [get]
func (h *VersionHandler) GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"api_version": "v1",
		"build_time":  "2024-01-01T00:00:00Z",
		"git_commit":  "abc123",
		"go_version":  "1.21",
	})
}