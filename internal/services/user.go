package services

import (
	"errors"
	"fmt"
	"taskflow-backend/internal/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

// GetUserByID 鏍规嵁ID鑾峰彇鐢ㄦ埛
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// ListUsers 鑾峰彇鐢ㄦ埛鍒楄〃锛堝垎椤碉級
func (s *UserService) ListUsers(req models.GetUsersParams) ([]models.User, int64, int64, error) {
	var users []models.User
	var total int64
	search := req.Search
	page := req.Page
	limit := req.Limit

	query := s.db.Model(&models.User{})
	// todo 缂哄皯鎸夎鑹叉煡鎵?鎸夌姸鎬佹煡鎵?鎸夐偖绠遍獙璇佺姸鎬佹煡鎵?鎸夊垱寤鸿捣濮嬫椂闂村拰缁撴潫鏃堕棿鏌ユ壘
	// 1. 鎼滅储鏉′欢锛堢敤鎴峰悕鎴栭偖绠憋級
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 2. 閭楠岃瘉鐘舵€?
	if req.EmailVerified != nil {
		query = query.Where("email_verified = ?", *req.EmailVerified)
		fmt.Println(query)
	}

	// 3. 鐢ㄦ埛瑙掕壊
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 4. 鐢ㄦ埛鐘舵€?
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 5. 鍒涘缓鏃堕棿鑼冨洿
	if !req.CreatedAtFrom.IsZero() {
		query = query.Where("created_at >= ?", req.CreatedAtFrom)
	}
	if !req.CreatedAtTo.IsZero() {
		query = query.Where("created_at <= ?", req.CreatedAtTo)
	}

	// 6. 鎺掑簭
	query = query.Order("created_at DESC")

	// 璁＄畻鎬绘暟
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 鍒嗛〉鏌ヨ
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, 0, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)
	return users, total, totalPages, nil
}

func (s *UserService) ToUserListResponse(userList []models.User) []interface{} {

	var q []interface{}
	for _, user := range userList {
		q = append(q, user.ToResponse())
	}
	return q
}

// UpdateUserStatus 鏇存柊鐢ㄦ埛鐘舵€?
func (s *UserService) UpdateUserStatus(id uint, status string) error {
	if status != "active" && status != "inactive" && status != "banned" {
		return errors.New("鏃犳晥鐨勭姸鎬佸€?)
	}

	result := s.db.Model(&models.User{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateUserRole 鏇存柊鐢ㄦ埛瑙掕壊
func (s *UserService) UpdateUserRole(id uint, role string) error {
	if role != "user" && role != "admin" && role != "super_admin" {
		return errors.New("鏃犳晥鐨勮鑹插€?)
	}

	result := s.db.Model(&models.User{}).Where("id = ?", id).Update("role", role)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateUser 鏇存柊鐢ㄦ埛淇℃伅
func (s *UserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 鏇存柊鐢ㄦ埛鍚嶏紙濡傛灉鎻愪緵锛?
	if req.Username != nil {
		// 妫€鏌ョ敤鎴峰悕鏄惁宸插瓨鍦?
		var count int64
		s.db.Model(&models.User{}).Where("username = ? AND id != ?", *req.Username, id).Count(&count)
		if count > 0 {
			return nil, errors.New("鐢ㄦ埛鍚嶅凡瀛樺湪")
		}
		user.Username = *req.Username
	}

	// 鏇存柊澶村儚锛堝鏋滄彁渚涳級
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Role != nil {

		user.Role = *req.Role
	}

	if req.Status != nil {
		user.Status = *req.Status
	}

	if req.Email_verified != nil {
		user.EmailVerified = *req.Email_verified
	}

	if req.Fullname != nil {
		user.Fullname = *req.Fullname
	}

	if req.Bio != nil {
		user.Bio = *req.Bio
	}

	// 淇濆瓨鏇存柊
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 淇敼瀵嗙爜
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 楠岃瘉鍘熷瘑鐮?
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidCredentials
	}

	// 楠岃瘉鏂板瘑鐮?
	if len(newPassword) < 6 {
		return ErrInvalidInput
	}

	// 璁剧疆鏂板瘑鐮?
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// 淇濆瓨鍒版暟鎹簱
	if err := s.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) DeleteUserByID(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}
	// 鍒犻櫎鐢ㄦ埛
	if err := s.db.Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func (s *UserService) ResetPassword(id uint, req models.ResetPasswordRequest) error {

	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := user.SetPassword(req.NewPassword); err != nil {
		return err
	}

	return s.db.Save(user).Error

}

func (s *UserService) GetUserStatus() (*models.UserStats, error) {
	var stats models.UserStats
	// 鏋勫缓SQL鏌ヨ
	sql := `
		SELECT
    COUNT(*) as total,
    COUNT(CASE WHEN status = 'active' THEN 1 END) as active,
    COUNT(CASE WHEN status = 'inactive' THEN 1 END) as inactive,
    COUNT(CASE WHEN status = 'banned' THEN 1 END) as banned,
    COUNT(CASE WHEN role = 'admin' THEN 1 END) as admin,
    COUNT(CASE WHEN role = 'super_admin' THEN 1 END) as super_admin,
    COUNT(CASE WHEN email_verified = FALSE THEN 1 END) as unverified
FROM users
WHERE deleted_at IS NULL
	`

	err := s.db.Raw(sql).Scan(&stats).Error
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

func (s *UserService) CreateUser(req models.CreateUserRequest) (*models.User, error) {

	// 楠岃瘉杈撳叆
	if err := validateUserInput(req.Username, req.Email, req.Password); err != nil {
		return nil, ErrInvalidInput
	}

	// 妫€鏌ョ敤鎴锋槸鍚﹀凡瀛樺湪
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	// 鍒涘缓鐢ㄦ埛
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Fullname:  req.Fullname,
		Bio:       req.Bio,
		Role:      req.Role,
		Status:    req.Status,
		AvatarURL: &req.AvatarURL,
	}

	// 璁剧疆瀵嗙爜
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 淇濆瓨鍒版暟鎹簱
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	if req.SendWelcomeEmail == true {
		// todo 鎺ュ叆kafka 鍙戦€佹杩巈mail
	}

	return user, nil
}
