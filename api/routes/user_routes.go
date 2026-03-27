package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 娉ㄥ唽鐢ㄦ埛绠＄悊璺敱锛堢鐞嗗憳鍔熻兘锛?
func RegisterUserRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, userHandler *handlers.UserHandler) {
	users := router.Group("/users")

	// 鐢ㄦ埛绠＄悊璺敱缁?(涓嶉渶瑕佹潈闄?
	users.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
	{
		users.GET("/profile", userHandler.GetProfile)
		users.PUT("/profile", userHandler.UpdateProfile)
		users.PUT("/password", userHandler.ChangePassword)
	}

	admin := users.Group("/admin")
	// 鐢ㄦ埛绠＄悊璺敱缁勶紙闇€瑕佺鐞嗗憳鏉冮檺锛?
	admin.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService), middleware.RequireAdmin())
	{
		// TODO: 瀹炵幇鐢ㄦ埛绠＄悊鎺ュ彛
		admin.GET("", userHandler.GetListUsers)
		admin.POST("", userHandler.CreateUser)
		admin.GET("/:id", userHandler.GetUser)
		admin.PUT("/:id", userHandler.UpdateUser)
		admin.DELETE("/:id", userHandler.DeleteUser)
		admin.POST("/:id/reset-password", userHandler.ResetPassword)
		admin.GET("/stats", userHandler.GetUsersStats)
	}
}
