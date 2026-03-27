package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理路由（管理员功能）
func RegisterUserRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, userHandler *handlers.UserHandler) {
	// 用户管理路由组（需要管理员权限）
	users := router.Group("/users")

	{
		users.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.PUT("/password", userHandler.ChangePassword)
	}

	users.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService), middleware.RequireAdmin())
	{
		// TODO: 实现用户管理接口
		users.GET("", userHandler.GetListUsers)
		users.GET("/:id", userHandler.GetUser)
		users.PUT("/:id", userHandler.UpdateUser)
		users.DELETE("/:id", userHandler.DeleteUser)
		users.POST("/:id/reset-password", userHandler.ResetPassword)
		users.GET("/stats", userHandler.GetUsersStats)
	}
}
