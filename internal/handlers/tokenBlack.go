package handlers

import (
	"taskflow-backend/internal/services"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type TokenBlackHandler struct {
	TokenBlacklistService *services.TokenBlacklistService
}

func NewTokenBlackHandler(redisClient *redis.Client, jwtManager *jwt.JWTManager) *TokenBlackHandler {
	return &TokenBlackHandler{
		TokenBlacklistService: services.NewTokenBlacklistService(redisClient, jwtManager),
	}
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出（客户端应删除令牌）
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *TokenBlackHandler) Logout(c *gin.Context) {
	// 在实际应用中，可以将令牌加入黑名单
	// 1. 从上下文获取当前 token
	token, exists := c.Get("token")
	if !exists {
		// 如果没有 token，可能是未登录状态
		Success(c, "登出成功", nil)
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		BadRequest(c, "无效的令牌")
		return
	}

	err := h.TokenBlacklistService.AddToBlacklist(c, tokenStr)
	if err != nil {
		logger.Error("加入黑名单失败",
			zap.String("token", tokenStr),
			zap.Error(err),
		)
		// 即使 Redis 操作失败，也要返回成功让客户端清除本地 token
		Success(c, "登出成功", nil)
		return
	}

	// 6. 返回成功
	Success(c, "登出成功", nil)
}
