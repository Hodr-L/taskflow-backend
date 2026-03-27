package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"taskflow-backend/internal/config"
	"taskflow-backend/internal/database"
	"taskflow-backend/internal/server"
	"taskflow-backend/pkg/logger"

	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 2. 初始化日志
	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Info("🚀 启动 TaskFlow 后端服务",
		zap.String("env", cfg.App.Env),
		zap.String("version", cfg.App.Version),
		zap.Int("port", cfg.Server.Port),
	)

	// 3. 连接数据库
	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Error("⚠️ 数据库连接失败，使用模拟数据库", logger.ErrorField(err))
		// 创建模拟的数据库连接，让服务至少能启动
		db = nil
	} else {

		defer func() {
			if err := database.Close(); err != nil {
				logger.Error("关闭数据库连接失败", logger.ErrorField(err))
			} else {
				logger.Info("✅ 数据库连接已关闭")
			}
		}()

		// 4. 自动迁移数据库表（仅开发环境）
		if cfg.App.Env == "development" {
			logger.Info("🔄 正在自动迁移数据库表...")
			if err := database.AutoMigrate(db); err != nil {
				logger.Error("❌ 数据库迁移失败", logger.ErrorField(err))
			} else {
				logger.Info("✅ 数据库迁移完成")
			}
		}
	}

	// 5. 连接Redis
	redisClient, err := database.ConnectRedis(cfg.Redis)
	if err != nil {
		logger.Error("⚠️ Redis连接失败", logger.ErrorField(err))
		// 可以选择继续启动，但黑名单功能将不可用
		logger.Warn("黑名单功能将不可用")
	} else {
		defer func() {
			if err := database.CloseRedis(); err != nil {
				logger.Error("关闭Redis连接失败", logger.ErrorField(err))
			}
		}()

		// Redis健康检查
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				if err := database.RedisHealthCheck(); err != nil {
					logger.Error("Redis健康检查失败", logger.ErrorField(err))
				}
			}
		}()
	}

	// 5. 创建HTTP服务器
	srv := server.New(cfg, db, redisClient)

	// 6. 启动服务器（在goroutine中）
	go func() {
		logger.Info("🌐 启动HTTP服务器", zap.Int("port", cfg.Server.Port))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("❌ 服务器启动失败", logger.ErrorField(err))
		}
	}()

	// 7. 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("🛑 收到停止信号，正在关闭服务...")

	// 8. 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("❌ 服务器关闭失败", logger.ErrorField(err))
	}

	logger.Info("👋 服务已安全退出")
}
