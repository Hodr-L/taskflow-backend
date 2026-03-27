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

// Logout 鐢ㄦ埛鐧诲嚭
// @Summary 鐢ㄦ埛鐧诲嚭
// @Description 鐢ㄦ埛鐧诲嚭锛堝鎴风搴斿垹闄や护鐗岋級
// @Tags 璁よ瘉
// @Security BearerAuth
// @Produce json
// @Success 200 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /auth/logout [post]
func (h *TokenBlackHandler) Logout(c *gin.Context) {
	// 鍦ㄥ疄闄呭簲鐢ㄤ腑锛屽彲浠ュ皢浠ょ墝鍔犲叆榛戝悕鍗?	// 1. 浠庝笂涓嬫枃鑾峰彇褰撳墠 token
	token, exists := c.Get("token")
	if !exists {
		// 濡傛灉娌℃湁 token锛屽彲鑳芥槸鏈櫥褰曠姸鎬?		Success(c, "鐧诲嚭鎴愬姛", nil)
		return
	}

	tokenStr, ok := token.(string)
	if !ok {
		BadRequest(c, "鏃犳晥鐨勪护鐗?)
		return
	}

	err := h.TokenBlacklistService.AddToBlacklist(c, tokenStr)
	if err != nil {
		logger.Error("鍔犲叆榛戝悕鍗曞け璐?,
			zap.String("token", tokenStr),
			zap.Error(err),
		)
		// 鍗充娇 Redis 鎿嶄綔澶辫触锛屼篃瑕佽繑鍥炴垚鍔熻瀹㈡埛绔竻闄ゆ湰鍦?token
		Success(c, "鐧诲嚭鎴愬姛", nil)
		return
	}

	// 6. 杩斿洖鎴愬姛
	Success(c, "鐧诲嚭鎴愬姛", nil)
}
