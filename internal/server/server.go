package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"taskflow-backend/internal/config"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"
	"taskflow-backend/pkg/logger"
)

type Server struct {
	config *config.Config
	db     *gorm.DB
	router *gin.Engine
	server *http.Server
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	// 设置Gin模式
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建Gin实例
	router := gin.New()

	// 添加中间件
	router.Use(gin.Recovery())
	router.Use(loggingMiddleware())

	// 配置CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Server.CORSAllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 创建JWT管理器
	jwtManager := jwt.NewJWTManager(cfg.JWT)

	// 注册路由
	s := &Server{
		config: cfg,
		db:     db,
		router: router,
	}

	s.setupRoutes(jwtManager)

	// 创建HTTP服务器
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout) * time.Second,
	}

	return s
}

func (s *Server) setupRoutes(jwtManager *jwt.JWTManager) {
	// 创建处理器
	authHandler := handlers.NewAuthHandler(s.db, jwtManager)
	teamHandler := handlers.NewTeamHandler(s.db)

	// API路由组
	api := s.router.Group("/api/v1")
	{
		// 健康检查
		api.GET("/health", s.healthCheck)

		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// 需要认证的路由
			auth.Use(middleware.JWTAuth(jwtManager))
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/profile", authHandler.GetProfile)
			auth.PUT("/profile", authHandler.UpdateProfile)
			auth.PUT("/password", authHandler.ChangePassword)
		}

		// 用户管理路由（需要管理员权限）
		users := api.Group("/users")
		users.Use(middleware.JWTAuth(jwtManager), middleware.RequireAdmin())
		{
			// TODO: 实现用户管理接口
			// users.GET("", userHandler.ListUsers)
			// users.GET("/:id", userHandler.GetUser)
			// users.PUT("/:id", userHandler.UpdateUser)
			// users.DELETE("/:id", userHandler.DeleteUser)
		}

		// 团队路由
		teams := api.Group("/teams")
		teams.Use(middleware.JWTAuth(jwtManager))
		{
			// TODO: 实现团队管理接口
			teams.GET("", teamHandler.GetTeams)
			teams.POST("", teamHandler.CreateTeam)
			teams.GET("/:id", teamHandler.GetTeam)
			teams.PUT("/:id", teamHandler.UpdateTeam)
			teams.DELETE("/:id", teamHandler.DeleteTeam)
		}

		// 项目路由
		projects := api.Group("/projects")
		projects.Use(middleware.JWTAuth(jwtManager))
		{
			// TODO: 实现项目管理接口
		}

		// 任务路由
		tasks := api.Group("/tasks")
		tasks.Use(middleware.JWTAuth(jwtManager))
		{
			// TODO: 实现任务管理接口
		}
	}

	// 静态文件服务（上传的文件）
	if s.config.Upload.StoragePath != "" {
		s.router.Static("/uploads", s.config.Upload.StoragePath)
	}

	// 404处理
	s.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "接口不存在",
		})
	})
}

// healthCheck 健康检查接口
func (s *Server) healthCheck(c *gin.Context) {
	// 检查数据库连接（如果数据库连接存在）
	if s.db != nil {
		db, err := s.db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "down",
				"message": "数据库连接失败",
				"error":   err.Error(),
			})
			return
		}

		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "down",
				"message": "数据库ping失败",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "服务运行正常",
			"version": s.config.App.Version,
			"env":     s.config.App.Env,
			"time":    time.Now().Format(time.RFC3339),
			"db":      "connected",
		})
	} else {
		// 数据库未连接，但服务仍在运行
		c.JSON(http.StatusOK, gin.H{
			"status":  "up",
			"message": "服务运行正常（数据库未连接）",
			"version": s.config.App.Version,
			"env":     s.config.App.Env,
			"time":    time.Now().Format(time.RFC3339),
			"db":      "disconnected",
			"warning": "数据库连接失败，部分功能可能不可用",
		})
	}
}

// loggingMiddleware 请求日志中间件
func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		// 跳过健康检查的详细日志
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

		// 添加用户信息（如果已认证）
		if userID, exists := c.Get("user_id"); exists {
			fields = append(fields, zap.Any("user_id", userID))
		}

		// 根据状态码记录不同级别的日志
		switch {
		case status >= 500:
			logger.Error("服务器错误", fields...)
		case status >= 400:
			logger.Warn("客户端错误", fields...)
		default:
			logger.Info("请求完成", fields...)
		}
	}
}

// Start 启动服务器
func (s *Server) Start() error {
	logger.Info("🚀 启动HTTP服务器",
		zap.String("address", s.server.Addr),
		zap.String("env", s.config.App.Env),
		zap.String("version", s.config.App.Version),
	)

	return s.server.ListenAndServe()
}

// Shutdown 优雅关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("🛑 正在关闭服务器...")
	return s.server.Shutdown(ctx)
}

// GetRouter 获取路由器（用于测试）
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}
