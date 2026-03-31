package services

import (
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"taskflow-backend/internal/models"
)

var ErrTeamNotFound = errors.New("团队不存在")

type TeamService struct {
	db *gorm.DB
}

func NewTeamService(db *gorm.DB) *TeamService {
	return &TeamService{db: db}
}

// CreateTeam 创建团队
func (s *TeamService) CreateTeam(userID models.UUID, name, description, logoUrl string) (*models.Team, error) {
	userService := NewUserService(s.db)
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	team := models.Team{
		Name:        name,
		Description: description,
		OwnerID:     userID,
		LogoURL:     &logoUrl,
		Owner:       user,
		InviteCode:  nil,
	}

	err = team.BeforeCreate(s.db)
	if err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := s.db.Create(&team).Error; err != nil {
		return nil, err
	}

	teamMember := models.TeamMember{
		TeamID: team.ID,
		UserID: userID,
		User:   user,
		Role:   "owner",
	}

	err = teamMember.BeforeCreate(s.db)
	if err != nil {
		return nil, err
	}

	// 创建团队所有者成员记录
	if err := s.db.Create(&teamMember).Error; err != nil {
		return nil, err
	}

	team.Members = append(team.Members, teamMember)

	// 更新团队所有者成员记录
	if err := s.db.Save(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// GetListTeams 获取团队列表逻辑
func (s *TeamService) GetListTeams(c *gin.Context, req models.GetTeamsParams) ([]models.Team, int64, int64, error) {
	var teams []models.Team
	var total int64
	search := req.Search
	page := req.Page
	limit := req.Limit

	query := s.db.Model(&models.Team{})

	query = query.Where("deleted_at IS NULL")

	if search != "" {
		query = query.Where("name LIKE ?",
			"%"+search+"%")
	}

	if req.OwnerId != "" {
		query = query.Where("owner_id = ?", req.OwnerId)
	}

	query = query.Order("created_at DESC")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Preload("Members.User").Preload("Projects").Find(&teams).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)
	return teams, total, totalPages, nil
}

// GetTeamById 获取团队详情
func (s *TeamService) GetTeamById(teamID models.UUID) (*models.Team, error) {
	var team models.Team
	if err := s.db.Preload("Owner").Preload("Members.User").Preload("Projects").First(&team, "id = ?", teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}
	return &team, nil
}

// UpdateTeam 更新团队信息
func (s *TeamService) UpdateTeam(teamID models.UUID, updates models.UpdateTeamRequest) (*models.Team, error) {
	var team models.Team

	if err := s.db.First(&team, "id = ?", teamID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTeamNotFound
		}
		return nil, err
	}

	team.Name = updates.Name
	team.Description = updates.Description
	team.LogoURL = updates.LogoURL

	if err := s.db.Save(&team).Error; err != nil {
		return nil, err
	}

	return &team, nil
}

// DeleteTeam 删除团队
func (s *TeamService) DeleteTeam(teamID models.UUID) error {
	team, err := s.GetTeamById(teamID)
	if err != nil {
		return err
	}

	if err := s.db.Delete(&team).Error; err != nil {
		return err
	}

	return nil
}

// AddTeamMember 添加团队成员
func (s *TeamService) AddTeamMember(teamID models.UUID, ownerID models.UUID, userID models.UUID, role string) error {
	// TODO: 实现添加成员逻辑
	return nil
}

// RemoveTeamMember 移除团队成员
func (s *TeamService) RemoveTeamMember(teamID models.UUID, ownerID models.UUID, userID models.UUID) error {
	// TODO: 实现移除成员逻辑
	return nil
}

func (t *TeamService) ToTeamResponse(team []models.Team) []interface{} {

	var q []interface{}
	for _, team := range team {
		q = append(q, team.ToResponse())
	}

	return q
}

func (t *TeamService) ToTeamMembersResponse(teamMembers []models.TeamMember) []interface{} {

	var q []interface{}
	for _, teamMember := range teamMembers {
		q = append(q, teamMember.ToTeamMemberResponse())
	}

	return q
}

func (t *TeamService) ToTeamProjectsResponse(teamProjects []models.Project) []interface{} {

	var q []interface{}
	for _, teamProject := range teamProjects {
		q = append(q, teamProject.ToTeamProjectResponse())
	}

	return q
}
