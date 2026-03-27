package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, authHandler *handlers.AuthHandler) {
	// 认证路由组
	auth := router.Group("/auth")
	{
		// 公开路由（无需认证）
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)

		// 需要认证的路由
		auth.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
		auth.POST("/logout", tokenBlackHandler.Logout)
	}
}
