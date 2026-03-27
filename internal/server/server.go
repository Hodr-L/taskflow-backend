package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"taskflow-backend/api/routes"
	"taskflow-backend/internal/config"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"
)

type Server struct {
	config      *config.Config
	db          *gorm.DB
	router      *gin.Engine
	server      *http.Server
	redisClient *redis.Client
}

func New(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *Server {
	// 璁剧疆Gin妯″紡
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 鍒涘缓Gin瀹炰緥
	router := gin.New()

	// 娣诲姞涓棿浠?	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())

	// 閰嶇疆CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 鍒涘缓JWT绠＄悊鍣?	jwtManager := jwt.NewJWTManager(cfg.JWT)

	// 娉ㄥ唽璺敱
	s := &Server{
		config:      cfg,
		db:          db,
		router:      router,
		redisClient: redisClient,
	}

	s.setupRoutes(jwtManager)

	// 鍒涘缓HTTP鏈嶅姟鍣?	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return s
}

func (s *Server) setupRoutes(jwtManager *jwt.JWTManager) {
	// 鍒涘缓澶勭悊鍣?	authHandler := handlers.NewAuthHandler(s.db, jwtManager)
	tokenBlackHandler := handlers.NewTokenBlackHandler(s.redisClient, jwtManager)
	teamHandler := handlers.NewTeamHandler(s.db)
	userHandler := handlers.NewUserHandler(s.db)

	// API璺敱缁?	api := s.router.Group("/api/v1")
	{
		// 鍋ュ悍妫€鏌?		api.GET("/health", s.healthCheck)

		// 娉ㄥ唽妯″潡鍖栬矾鐢?		routes.RegisterAuthRoutes(api, jwtManager, tokenBlackHandler, authHandler)
		routes.RegisterUserRoutes(api, jwtManager, tokenBlackHandler, userHandler)
		routes.RegisterTeamRoutes(api, jwtManager, tokenBlackHandler, teamHandler)

		api.GET("/debug/test", func(c *gin.Context) {
			logger.Debug("Debug绾у埆鏃ュ織")
			logger.Info("Info绾у埆鏃ュ織")
			logger.Warn("Warn绾у埆鏃ュ織")
			c.JSON(200, gin.H{"message": "娴嬭瘯鎴愬姛"})
		})
		// TODO: 娉ㄥ唽鍏朵粬妯″潡璺敱锛堝緟瀹炵幇瀵瑰簲handler鍚庯級
		// routes.RegisterProjectRoutes(api, jwtManager, projectHandler)
		// routes.RegisterTaskRoutes(api, jwtManager, taskHandler)
	}

	// 闈欐€佹枃浠舵湇鍔★紙涓婁紶鐨勬枃浠讹級
	if s.config.Upload.StoragePath != "" {
		s.router.Static("/uploads", s.config.Upload.StoragePath)
	}

	// 404澶勭悊
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "鎺ュ彛涓嶅瓨鍦?,
		})
	})
}

// healthCheck 鍋ュ悍妫€鏌ユ帴鍙?func (s *Server) healthCheck(c *gin.Context) {
	// 妫€鏌ユ暟鎹簱杩炴帴锛堝鏋滄暟鎹簱杩炴帴瀛樺湪锛?	if s.db != nil {
		db, err := s.db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "down",
				"message": "鏁版嵁搴撹繛鎺ュけ璐?,
				"error":   err.Error(),
			})
			return
		}

		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "down",
				"message": "鏁版嵁搴損ing澶辫触",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "鏈嶅姟杩愯姝ｅ父",
			"version": s.config.App.Version,
			"env":     s.config.App.Env,
			"time":    time.Now().Format(time.RFC3339),
			"db":      "connected",
		})
	} else {
		// 鏁版嵁搴撴湭杩炴帴锛屼絾鏈嶅姟浠嶅湪杩愯
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "鏈嶅姟杩愯姝ｅ父锛堟暟鎹簱鏈繛鎺ワ級",
			"version": s.config.App.Version,
			"env":     s.config.App.Env,
			"time":    time.Now().Format(time.RFC3339),
			"db":      "disconnected",
			"warning": "鏁版嵁搴撹繛鎺ュけ璐ワ紝閮ㄥ垎鍔熻兘鍙兘涓嶅彲鐢?,
		})
	}
}

// loggingMiddleware 璇锋眰鏃ュ織涓棿浠?func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 澶勭悊璇锋眰
		c.Next()

		// 璁板綍鏃ュ織
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 璺宠繃鍋ュ悍妫€鏌ョ殑璇︾粏鏃ュ織
		if path == "/api/v1/health" {
			return
		}

		fields := []zap.Field{
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", c.Request.UserAgent()),
		}

		// 娣诲姞鐢ㄦ埛淇℃伅锛堝鏋滃凡璁よ瘉锛?		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		// 鏍规嵁鐘舵€佺爜璁板綍涓嶅悓绾у埆鐨勬棩蹇?		switch {
		case status >= 500:
			logger.Error("鏈嶅姟鍣ㄩ敊璇?, fields...)
		case status >= 400:
			logger.Warn("瀹㈡埛绔敊璇?, fields...)
		default:
			logger.Info("璇锋眰瀹屾垚", fields...)
		}
	}
}

// Start 鍚姩鏈嶅姟鍣?func (s *Server) Start() error {
	logger.Info("馃殌 鍚姩HTTP鏈嶅姟鍣?,
		zap.String("address", s.server.Addr),
		zap.String("env", s.config.App.Env),
		zap.String("version", s.config.App.Version),
	)

	return s.server.ListenAndServe()
}

// Shutdown 浼橀泤鍏抽棴鏈嶅姟鍣?func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("馃洃 姝ｅ湪鍏抽棴鏈嶅姟鍣?..")
	return s.server.Shutdown(ctx)
}

// GetRouter 鑾峰彇璺敱鍣紙鐢ㄤ簬娴嬭瘯锛?func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
