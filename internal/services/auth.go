package services

import (
	"errors"
	"regexp"

	"taskflow-backend/internal/models"

	"gorm.io/gorm"
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
