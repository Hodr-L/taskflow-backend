package handlers

import (
	"taskflow-backend/internal/models"
	"taskflow-backend/pkg/logger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/services"
)

type TeamHandler struct {
	teamService *services.TeamService
	db          *gorm.DB
}

func NewTeamHandler(db *gorm.DB) *TeamHandler {
	return &TeamHandler{
		teamService: services.NewTeamService(db),
		db:          db,
	}
}

// CreateTeam 创建团队
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.CreateTeamParams
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("创建团队参数验证失败", logger.ErrorField(err))
		BadRequest(c, "参数验证失败", err)
		return
	}

	userIDStr, exists := c.Get("user_id")
	if !exists {
		InternalServerError(c, "userID 不存在")
		return
	}

	userID, err := models.ParseUUID(userIDStr.(string))
	if err != nil {
		BadRequest(c, "无效的用户ID格式")
		return
	}

	team, err := h.teamService.CreateTeam(userID, req.Name, req.Description, req.LogoURL)
	if err != nil {
		InternalServerError(c, "创建团队失败", err)
		return
	}

	Success(c, "创建团队成功", team.ToResponse())
}

func (h *TeamHandler) GetListTeams(c *gin.Context) {
	// TODO: 实现获取团队列表逻辑
}

// GetTeams 获取团队列表
func (h *TeamHandler) GetTeams(c *gin.Context) {
	// TODO: 获取用户参与的团队列表
}

// GetTeam 获取团队详情
func (h *TeamHandler) GetTeam(c *gin.Context) {
	// TODO: 实现获取团队详情逻辑
}

// UpdateTeam 更新团队
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	// TODO: 实现更新团队逻辑
}

// DeleteTeam 删除团队
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	// TODO: 实现删除团队逻辑
}
