package routes

import (
	"github.com/gin-gonic/gin"
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"
)

// RegisterProjectRoutes 娉ㄥ唽椤圭洰鐩稿叧璺敱
func RegisterProjectRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, projectHandler *handlers.ProjectHandler) {
	// 椤圭洰璺敱缁?	projects := router.Group("/projects")
	projects.Use(middleware.JWTAuth(jwtManager))
	{
		// TODO: 瀹炵幇椤圭洰绠＄悊鎺ュ彛
		// projects.GET("", projectHandler.GetProjects)
		// projects.POST("", projectHandler.CreateProject)
		// projects.GET("/:id", projectHandler.GetProject)
		// projects.PUT("/:id", projectHandler.UpdateProject)
		// projects.DELETE("/:id", projectHandler.DeleteProject)

		// TODO: 娣诲姞椤圭洰浠诲姟绠＄悊璺敱
		// projects.GET("/:id/tasks", projectHandler.GetProjectTasks)
		// projects.GET("/:id/members", projectHandler.GetProjectMembers)
	}
}
