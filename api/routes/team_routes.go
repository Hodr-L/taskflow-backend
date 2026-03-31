package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterTeamRoutes 注册团队相关路由
func RegisterTeamRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, tokenBlackHandler *handlers.TokenBlackHandler, teamHandler *handlers.TeamHandler) {
	// 团队路由组
	teams := router.Group("/teams")
	teams.Use(middleware.JWTAuth(jwtManager, tokenBlackHandler.TokenBlacklistService))
	{
		// 团队列表和创建
		teams.GET("", teamHandler.GetListTeams)
		teams.POST("", teamHandler.CreateTeam)

		// 单个团队操作
		teams.GET("/:id", teamHandler.GetTeam)
		teams.PUT("/:id", teamHandler.UpdateTeam)
		teams.DELETE("/:id", teamHandler.DeleteTeam)

		// TODO: 添加团队成员管理路由
		// teams.POST("/:id/members", teamHandler.AddMember)
		// teams.DELETE("/:id/members/:userId", teamHandler.RemoveMember)
		// teams.GET("/:id/members", teamHandler.GetMembers)
	}
}
