package api

import (
	"net/http"

	"devops/internal/auth"
	"devops/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService *service.UserService
	jwtManager  *auth.JWTManager
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(db *gorm.DB, rdb *redis.Client, jwtSecret string) *AuthHandler {
	userService := service.NewUserService(db, rdb)
	jwtManager := auth.NewJWTManager(jwtSecret, 3600, 604800)
	
	return &AuthHandler{
		userService: userService,
		jwtManager:  jwtManager,
	}
}

// Login 用户登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证用户凭据
	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: err.Error(),
		})
		return
	}

	// 生成令牌
	accessToken, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "生成访问令牌失败",
		})
		return
	}

	refreshToken, err := h.jwtManager.GenerateRefreshToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "生成刷新令牌失败",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "登录成功",
		Data: LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User: User{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
				Role:     user.Role,
				Status:   user.Status,
			},
		},
	})
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 刷新令牌
	accessToken, newRefreshToken, err := h.jwtManager.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "刷新令牌失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "刷新成功",
		Data: gin.H{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		},
	})
}

// GetProfile 获取用户信息
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "用户未认证",
		})
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "用户不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data: User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	})
}

// Logout 用户登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 在Redis中可以维护一个黑名单来失效token
	// 这里简单返回成功
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "登出成功",
	})
}