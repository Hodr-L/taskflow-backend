package routes

import (
	"taskflow-backend/internal/handlers"
	"taskflow-backend/internal/middleware"
	"taskflow-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// RegisterTaskRoutes 娉ㄥ唽浠诲姟鐩稿叧璺敱
func RegisterTaskRoutes(router *gin.RouterGroup, jwtManager *jwt.JWTManager, taskHandler *handlers.TaskHandler) {
	// 浠诲姟璺敱缁?	tasks := router.Group("/tasks")
	tasks.Use(middleware.JWTAuth(jwtManager))
	{
		// TODO: 瀹炵幇浠诲姟绠＄悊鎺ュ彛
		// tasks.GET("", taskHandler.GetTasks)
		// tasks.POST("", taskHandler.CreateTask)
		// tasks.GET("/:id", taskHandler.GetTask)
		// tasks.PUT("/:id", taskHandler.UpdateTask)
		// tasks.DELETE("/:id", taskHandler.DeleteTask)

		// TODO: 娣诲姞浠诲姟瀛愬姛鑳借矾鐢?		// tasks.POST("/:id/comments", taskHandler.AddComment)
		// tasks.GET("/:id/comments", taskHandler.GetComments)
		// tasks.POST("/:id/attachments", taskHandler.UploadAttachment)
		// tasks.GET("/:id/attachments", taskHandler.GetAttachments)
	}
}
