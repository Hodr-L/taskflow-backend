package routes

import (
	"github.com/gin-gonic/gin"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, authHandler *handlers.AuthHandler) {
	// 认证路由组
	auth := router.Group("/auth")
	{
		// 公开路由（无需认证）
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
}
