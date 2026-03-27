package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterTeamRoutes 娉ㄥ唽鍥㈤槦鐩稿叧璺敱
func RegisterTeamRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, teamHandler *handlers.TeamHandler) {
	// 鍥㈤槦璺敱缁?	teams := router.Group("/teams")
	teams.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
	{
		// 鍥㈤槦鍒楄〃鍜屽垱寤?		teams.GET("", teamHandler.GetTeams)
		teams.POST("", teamHandler.CreateTeam)

		// 鍗曚釜鍥㈤槦鎿嶄綔
		teams.GET("/:id", teamHandler.GetTeam)
		teams.PUT("/:id", teamHandler.UpdateTeam)
		teams.DELETE("/:id", teamHandler.DeleteTeam)

		// TODO: 娣诲姞鍥㈤槦鎴愬憳绠＄悊璺敱
		// teams.POST("/:id/members", teamHandler.AddMember)
		// teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
		// teams.GET("/:id/members", teamHandler.GetMembers)
	}
}
