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

// Create 鍒涘缓鍥㈤槦
func (r *TeamRepository) Create(team *models.Team) error {
	return r.db.Create(team).Error
}

// FindByID 鏍规嵁ID鏌ユ壘鍥㈤槦
func (r *TeamRepository) FindByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.First(&team, id).Error
	return &team, err
}

// FindByUserID 鏌ユ壘鐢ㄦ埛鍙備笌鐨勫洟闃?func (r *TeamRepository) FindByUserID(userID uint, page, limit int) ([]models.Team, int64, error) {
	var teams []models.Team
	var total int64

	// 鏌ヨ鐢ㄦ埛鍙備笌鐨勫洟闃?	query := r.db.Model(&models.Team{}).
		Joins("JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ?", userID)

	// 鑾峰彇鎬绘暟
	query.Count(&total)

	// 鍒嗛〉鏌ヨ
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&teams).Error

	return teams, total, err
}

// Update 鏇存柊鍥㈤槦
func (r *TeamRepository) Update(team *models.Team) error {
	return r.db.Save(team).Error
}

// Delete 鍒犻櫎鍥㈤槦
func (r *TeamRepository) Delete(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}
