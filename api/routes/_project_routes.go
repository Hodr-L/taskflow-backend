package routes

import (
	"github.com/gin-gonic/gin"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"
)

// RegisterProjectRoutes 注册项目相关路由
func RegisterProjectRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, projectHandler *handlers.ProjectHandler) {
	// 项目路由组
	projects := router.Group("/projects")
	projects.Use(middleware.JWTAuth(jwtManager))
	{
		// TODO: 实现项目管理接口
		// projects.GET("", projectHandler.GetProjects)
		// projects.POST("", projectHandler.CreateProject)
		// projects.GET("/:id", projectHandler.GetProject)
		// projects.PUT("/:id", projectHandler.UpdateProject)
		// projects.DELETE("/:id", projectHandler.DeleteProject)
		
		// TODO: 添加项目任务管理路由
		// projects.GET("/:id/tasks", projectHandler.GetProjectTasks)
		// projects.GET("/:id/members", projectHandler.GetProjectMembers)
	}
}