// database/redis.go
package database

import (
	"context"
	"fmt"
	"time"

	"taskflow-backend/internal/config"
	"taskflow-backend/pkg/logger"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// ConnectRedis 连接Redis
func ConnectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	options := &redis.Options{
		Addr:         cfg.GetAddr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	}

	client := redis.NewClient(options)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis连接失败: %w", err)
	}

	redisClient = client
	logger.Info("✅ Redis连接成功")

	return redisClient, nil
}

// GetRedis 获取Redis客户端
func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() error {
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			return fmt.Errorf("关闭Redis连接失败: %w", err)
		}
		logger.Info("✅ Redis连接已关闭")
	}
	return nil
}

// RedisHealthCheck Redis健康检查
func RedisHealthCheck() error {
	if redisClient == nil {
		return fmt.Errorf("Redis客户端未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}
