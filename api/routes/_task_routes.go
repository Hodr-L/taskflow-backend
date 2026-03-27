package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterTaskRoutes 注册任务相关路由
func RegisterTaskRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, taskHandler *handlers.TaskHandler) {
	// 任务路由组
	tasks := router.Group("/tasks")
	tasks.Use(middleware.JWTAuth(jwtManager))
	{
		// TODO: 实现任务管理接口
		// tasks.GET("", taskHandler.GetTasks)
		// tasks.POST("", taskHandler.CreateTask)
		// tasks.GET("/:id", taskHandler.GetTask)
		// tasks.PUT("/:id", taskHandler.UpdateTask)
		// tasks.DELETE("/:id", taskHandler.DeleteTask)

		// TODO: 添加任务子功能路由
		// tasks.POST("/:id/comments", taskHandler.AddComment)
		// tasks.GET("/:id/comments", taskHandler.GetComments)
		// tasks.POST("/:id/attachments", taskHandler.UploadAttachment)
		// tasks.GET("/:id/attachments", taskHandler.GetAttachments)
	}
}
