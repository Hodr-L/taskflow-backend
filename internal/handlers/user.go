package handlers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/services"
)

type UserHandler struct {
	userService *services.UserService
	db          *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: services.NewUserService(db),
		db:          db,
	}
}

// GetTeams 获取用户列表
func (h *UserHandler) GetListUsers(c *gin.Context) {
	// TODO: 实现获取用户列表逻辑

}

// GetTeam 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	// TODO: 实现获取用户详情逻辑
}

// UpdateTeam 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// TODO: 实现更新用户逻辑
}

// DeleteTeam 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// TODO: 实现删除用户逻辑
}
