package services

import (
	"fmt"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type TokenBlacklistService struct {
	redisClient *redis.Client
	jwtManager  *jwt.JWTManager
}

func NewTokenBlacklistService(redisClient *redis.Client, jwtManager *jwt.JWTManager) *TokenBlacklistService {
	return &TokenBlacklistService{
		redisClient: redisClient,
		jwtManager:  jwtManager,
	}
}

// AddToBlacklist 灏?token 鍔犲叆榛戝悕鍗?func (s *TokenBlacklistService) AddToBlacklist(ctx *gin.Context, token string) error {
	// 瑙ｆ瀽 token 鑾峰彇杩囨湡鏃堕棿
	claims, err := s.jwtManager.ParseToken(token)
	if err != nil {
		return fmt.Errorf("token 瑙ｆ瀽澶辫触: %w", err)
	}

	// 濡傛灉 token 宸茬粡杩囨湡锛屼笉闇€瑕佸姞鍏ラ粦鍚嶅崟
	expTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	if expTime.Before(time.Now()) {
		return nil
	}

	// 璁＄畻鍓╀綑鏃堕棿
	ttl := time.Until(expTime)

	// 浣跨敤 token 鏈韩浣滀负 key
	key := fmt.Sprintf("token:blacklist:%s", token)

	// 鎴栬€呬娇鐢?jti锛堝鏋?JWT 涓湁 JWT ID锛?	// key := fmt.Sprintf("token:blacklist:%s", claims.JTI)

	logger.Info("token宸插姞鍏ラ粦鍚嶅崟",
		zap.Uint("user_id", claims.UserID),
		zap.String("username", claims.Username),
		zap.String("token", token[:20]+"..."), // 鍙褰曢儴鍒?token
		zap.Duration("ttl", ttl),
	)
	return s.redisClient.SetEX(ctx, key, "blacklisted", ttl).Err()
}

// IsBlacklisted 妫€鏌?token 鏄惁鍦ㄩ粦鍚嶅崟涓?func (s *TokenBlacklistService) IsBlacklisted(ctx *gin.Context, token string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", token)

	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// LogoutUser 鐧诲嚭鐢ㄦ埛锛堝彲閫夛細鍙互鐧诲嚭鎵€鏈夎澶囷級
func (s *TokenBlacklistService) LogoutUser(ctx *gin.Context, userID uint) error {
	// 璁板綍鐢ㄦ埛鐧诲嚭鏃堕棿锛屽彲浠ョ敤鏉ヤ娇璇ユ椂闂翠箣鍓嶇殑鎵€鏈?token 澶辨晥
	key := fmt.Sprintf("user:logout:%d", userID)
	return s.redisClient.Set(ctx, key, time.Now().Unix(), 24*time.Hour).Err()
}
