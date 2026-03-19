package services

import (
	"errors"
	"regexp"

	"gorm.io/gorm"
	"taskflow-backend/internal/models"
)

var (
	ErrUserExists         = errors.New("用户已存在")
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("无效的凭据")
	ErrInvalidInput       = errors.New("无效的输入")
	ErrUserInactive       = errors.New("用户账户未激活")
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

// Register 注册新用户
func (s *AuthService) Register(username, email, password string) (*models.User, error) {
	// 验证输入
	if err := validateUserInput(username, email, password); err != nil {
		return nil, ErrInvalidInput
	}

	// 检查用户是否已存在
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	// 创建用户
	user := &models.User{
		Username: username,
		Email:    email,
		Role:     "user",
		Status:   "active",
	}

	// 设置密码
	if err := user.SetPassword(password); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *AuthService) Login(email, password string) (*models.User, error) {
	// 查找用户
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 检查用户状态
	if !user.IsActive() {
		return nil, ErrUserInactive
	}

	// 验证密码
	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return &user, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *AuthService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (s *AuthService) UpdateUser(id uint, req models.UserUpdateRequest) (*models.User, error) {
	user, err := s.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// 更新用户名（如果提供）
	if req.Username != nil {
		// 检查用户名是否已存在
		var count int64
		s.db.Model(&models.User{}).Where("username = ? AND id != ?", *req.Username, id).Count(&count)
		if count > 0 {
			return nil, errors.New("用户名已存在")
		}
		user.Username = *req.Username
	}

	// 更新头像（如果提供）
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	// 保存更新
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(id uint, oldPassword, newPassword string) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证原密码
	if !user.CheckPassword(oldPassword) {
		return ErrInvalidCredentials
	}

	// 验证新密码
	if len(newPassword) < 6 {
		return ErrInvalidInput
	}

	// 设置新密码
	if err := user.SetPassword(newPassword); err != nil {
		return err
	}

	// 保存到数据库
	if err := s.db.Save(user).Error; err != nil {
		return err
	}

	return nil
}

// validateUserInput 验证用户输入
func validateUserInput(username, email, password string) error {
	// 验证用户名
	if len(username) < 3 || len(username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	// 验证邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("邮箱格式不正确")
	}

	// 验证密码
	if len(password) < 6 {
		return errors.New("密码长度至少6个字符")
	}

	return nil
}

// ListUsers 获取用户列表（分页）
func (s *AuthService) ListUsers(page, limit int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	query := s.db.Model(&models.User{})

	// 搜索条件
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// UpdateUserStatus 更新用户状态
func (s *AuthService) UpdateUserStatus(id uint, status string) error {
	if status != "active" && status != "inactive" && status != "banned" {
		return errors.New("无效的状态值")
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

// UpdateUserRole 更新用户角色
func (s *AuthService) UpdateUserRole(id uint, role string) error {
	if role != "user" && role != "admin" && role != "super_admin" {
		return errors.New("无效的角色值")
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
