package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理路由（管理员功能）
func RegisterUserRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, userHandler *handlers.UserHandler) {
	users := router.Group("/users")

	// 用户管理路由组 (不需要权限)
	users.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.PUT("/password", userHandler.ChangePassword)
	}

	admin := users.Group("/admin")
	// 用户管理路由组（需要管理员权限）
	admin.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService), middleware.RequireAdmin())
	{
		// TODO: 实现用户管理接口
		admin.GET("", userHandler.GetListUsers)
		admin.POST("", userHandler.CreateUser)
		admin.GET("/:id", userHandler.GetUser)
		admin.PUT("/:id", userHandler.UpdateUser)
		admin.DELETE("/:id", userHandler.DeleteUser)
		admin.POST("/:id/reset-password", userHandler.ResetPassword)
		admin.GET("/stats", userHandler.GetUsersStats)
	}
}
