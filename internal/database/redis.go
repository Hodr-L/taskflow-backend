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

// ConnectRedis 杩炴帴Redis
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

	// 娴嬭瘯杩炴帴
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis杩炴帴澶辫触: %w", err)
	}

	redisClient = client
	logger.Info("鉁?Redis杩炴帴鎴愬姛")

	return redisClient, nil
}

// GetRedis 鑾峰彇Redis瀹㈡埛绔?func GetRedis() *redis.Client {
	return redisClient
}

// CloseRedis 鍏抽棴Redis杩炴帴
func CloseRedis() error {
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			return fmt.Errorf("鍏抽棴Redis杩炴帴澶辫触: %w", err)
		}
		logger.Info("鉁?Redis杩炴帴宸插叧闂?)
	}
	return nil
}

// RedisHealthCheck Redis鍋ュ悍妫€鏌?func RedisHealthCheck() error {
	if redisClient == nil {
		return fmt.Errorf("Redis瀹㈡埛绔湭鍒濆鍖?)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return redisClient.Ping(ctx).Err()
}
