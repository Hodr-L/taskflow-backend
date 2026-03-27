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
	// 1. 鍔犺浇閰嶇疆
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("鉂?鍔犺浇閰嶇疆澶辫触: %v", err)
	}

	// 2. 鍒濆鍖栨棩蹇?	logger.Init(cfg.Log.Level, cfg.Log.Format)
	logger.Info("馃殌 鍚姩 TaskFlow 鍚庣鏈嶅姟",
		zap.String("env", cfg.App.Env),
		zap.String("version", cfg.App.Version),
		zap.Int("port", cfg.Server.Port),
	)

	// 3. 杩炴帴鏁版嵁搴?	db, err := database.Connect(cfg.Database)
	if err != nil {
		logger.Error("鈿狅笍 鏁版嵁搴撹繛鎺ュけ璐ワ紝浣跨敤妯℃嫙鏁版嵁搴?, logger.ErrorField(err))
		// 鍒涘缓妯℃嫙鐨勬暟鎹簱杩炴帴锛岃鏈嶅姟鑷冲皯鑳藉惎鍔?		db = nil
	} else {

		defer func() {
			if err := database.Close(); err != nil {
				logger.Error("鍏抽棴鏁版嵁搴撹繛鎺ュけ璐?, logger.ErrorField(err))
			} else {
				logger.Info("鉁?鏁版嵁搴撹繛鎺ュ凡鍏抽棴")
			}
		}()

		// 4. 鑷姩杩佺Щ鏁版嵁搴撹〃锛堜粎寮€鍙戠幆澧冿級
		if cfg.App.Env == "development" {
			logger.Info("馃攧 姝ｅ湪鑷姩杩佺Щ鏁版嵁搴撹〃...")
			if err := database.AutoMigrate(db); err != nil {
				logger.Error("鉂?鏁版嵁搴撹縼绉诲け璐?, logger.ErrorField(err))
			} else {
				logger.Info("鉁?鏁版嵁搴撹縼绉诲畬鎴?)
			}
		}
	}

	// 5. 杩炴帴Redis
	redisClient, err := database.ConnectRedis(cfg.Redis)
	if err != nil {
		logger.Error("鈿狅笍 Redis杩炴帴澶辫触", logger.ErrorField(err))
		// 鍙互閫夋嫨缁х画鍚姩锛屼絾榛戝悕鍗曞姛鑳藉皢涓嶅彲鐢?		logger.Warn("榛戝悕鍗曞姛鑳藉皢涓嶅彲鐢?)
	} else {
		defer func() {
			if err := database.CloseRedis(); err != nil {
				logger.Error("鍏抽棴Redis杩炴帴澶辫触", logger.ErrorField(err))
			}
		}()

		// Redis鍋ュ悍妫€鏌?		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()

			for range ticker.C {
				if err := database.RedisHealthCheck(); err != nil {
					logger.Error("Redis鍋ュ悍妫€鏌ュけ璐?, logger.ErrorField(err))
				}
			}
		}()
	}

	// 5. 鍒涘缓HTTP鏈嶅姟鍣?	srv := server.New(cfg, db, redisClient)

	// 6. 鍚姩鏈嶅姟鍣紙鍦╣oroutine涓級
	go func() {
		logger.Info("馃寪 鍚姩HTTP鏈嶅姟鍣?, zap.Int("port", cfg.Server.Port))
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("鉂?鏈嶅姟鍣ㄥ惎鍔ㄥけ璐?, logger.ErrorField(err))
		}
	}()

	// 7. 绛夊緟涓柇淇″彿
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("馃洃 鏀跺埌鍋滄淇″彿锛屾鍦ㄥ叧闂湇鍔?..")

	// 8. 浼橀泤鍏抽棴
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("鉂?鏈嶅姟鍣ㄥ叧闂け璐?, logger.ErrorField(err))
	}

	logger.Info("馃憢 鏈嶅姟宸插畨鍏ㄩ€€鍑?)
}

