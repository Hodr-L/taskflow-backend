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

// AddToBlacklist 将 token 加入黑名单
func (s *TokenBlacklistService) AddToBlacklist(ctx *gin.Context, token string) error {
	// 解析 token 获取过期时间
	claims, err := s.jwtManager.ParseToken(token)
	if err != nil {
		return fmt.Errorf("token 解析失败: %w", err)
	}

	// 如果 token 已经过期，不需要加入黑名单
	expTime := time.Unix(claims.ExpiresAt.Unix(), 0)
	if expTime.Before(time.Now()) {
		return nil
	}

	// 计算剩余时间
	ttl := time.Until(expTime)

	// 使用 token 本身作为 key
	key := fmt.Sprintf("token:blacklist:%s", token)

	// 或者使用 jti（如果 JWT 中有 JWT ID）
	// key := fmt.Sprintf("token:blacklist:%s", claims.JTI)

	logger.Info("token已加入黑名单",
		zap.Uint("user_id", claims.UserID),
		zap.String("username", claims.Username),
		zap.String("token", token[:20]+"..."), // 只记录部分 token
		zap.Duration("ttl", ttl),
	)
	return s.redisClient.SetEX(ctx, key, "blacklisted", ttl).Err()
}

// IsBlacklisted 检查 token 是否在黑名单中
func (s *TokenBlacklistService) IsBlacklisted(ctx *gin.Context, token string) (bool, error) {
	key := fmt.Sprintf("token:blacklist:%s", token)

	exists, err := s.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// LogoutUser 登出用户（可选：可以登出所有设备）
func (s *TokenBlacklistService) LogoutUser(ctx *gin.Context, userID uint) error {
	// 记录用户登出时间，可以用来使该时间之前的所有 token 失效
	key := fmt.Sprintf("user:logout:%d", userID)
	return s.redisClient.Set(ctx, key, time.Now().Unix(), 24*time.Hour).Err()
}
