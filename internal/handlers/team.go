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
	var req models.GetTeamsParams
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.Warn("查询参数验证失败", logger.ErrorField(err))
		BadRequest(c, "参数验证失败", err)
		return
	}

	teamList, total, totalPages, err := h.teamService.GetListTeams(c, req)
	if err != nil {
		InternalServerError(c, "获取团队失败！", err)
		return
	}

	teamJson := h.teamService.ToTeamResponse(teamList)

	teamReq := models.TeamListResponse{
		Teams: teamJson,
		Pagination: models.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	Success(c, "团队列表获取成功", teamReq)
}

// GetTeams 获取团队列表
func (h *TeamHandler) GetTeams(c *gin.Context) {
	// TODO: 获取用户参与的团队列表
}

// GetTeam 获取团队详情
func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamID, err := ParseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的团队ID格式")
		return
	}

	team, err := h.teamService.GetTeamById(teamID)
	if err != nil {
		if err == services.ErrUserNotFound {
			NotFound(c, "团队不存在")
		} else {
			InternalServerError(c, "获取团队信息失败", err)
		}
		return
	}

	req := models.GetTeamByIDResponse{
		Team: models.Team{ID: teamID,
			Name:        team.Name,
			Description: team.Description,
			LogoURL:     team.LogoURL,
			OwnerID:     team.OwnerID,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt},

		Members:  h.teamService.ToTeamMembersResponse(team.Members),
		Projects: h.teamService.ToTeamProjectsResponse(team.Projects),
	}

	Success(c, "获取团队信息成功", req)
}

// UpdateTeam 更新团队
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	teamID, err := ParseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的团队ID格式")
		return
	}

	var request models.UpdateTeamRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		BadRequest(c, "参数验证失败", err)
		return
	}

	team, err := h.teamService.UpdateTeam(teamID, request)
	if err != nil {
		InternalServerError(c, "更细失败", err)
		return
	}

	var req = map[string]map[string]interface{}{
		"team": {
			"Id":          team.ID,
			"Name":        team.Name,
			"Description": team.Description,
			"LogoURL":     team.LogoURL,
			"UpdatedAt":   team.UpdatedAt,
		},
	}

	Success(c, "团队更新成功", req)
}

// DeleteTeam 删除团队
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	teamID, err := ParseUUIDParam(c, "id")
	if err != nil {
		BadRequest(c, "无效的团队ID格式", err)
		return
	}

	err = h.teamService.DeleteTeam(teamID)
	if err != nil {
		InternalServerError(c, "删除团队失败", err)
		return
	}
}
