package repositories

import (
	"gorm.io/gorm"

	"taskflow-backend/internal/models"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

// Create 创建团队
func (r *TeamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

// FindByID 根据ID查找团队
func (r *TeamRepository) FindByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	return &team, err
}

// FindByUserID 查找用户参与的团队
func (r *TeamRepository) FindByUserID(userID uint, page, limit int) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	// 查询用户参与的团队
	query := r.db.Model(&models.Team{}).
		Joins("JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userID)

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&teams).Error

	return teams, total, err
}

// Update 更新团队
func (r *TeamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

// Delete 删除团队
func (r *TeamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}
