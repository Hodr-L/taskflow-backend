package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes 娉ㄥ唽璁よ瘉鐩稿叧璺敱
func RegisterAuthRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, authHandler *handlers.AuthHandler) {
	// 璁よ瘉璺敱缁?
	auth := router.Group("/auth")
	{
		// 鍏紑璺敱锛堟棤闇€璁よ瘉锛?
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
	}

	// 闇€瑕佽璇佺殑璺敱
	token := auth.Group("/token")
	token.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
	{
		token.POST("/logout", tokenBlackHandler.Logout)
	}
}
