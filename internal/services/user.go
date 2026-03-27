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

// GetUserByID 根据ID获取用户
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

// ListUsers 获取用户列表（分页）
func (s *UserService) ListUsers(req models.GetUsersParams) ([]models.User, int64, int64, error) {
	var users []models.User
	var total int64
	search := req.Search
	page := req.Page
	limit := req.Limit

	query := s.db.Model(&models.User{})
	// todo 缺少按角色查找 按状态查找 按邮箱验证状态查找 按创建起始时间和结束时间查找
	// 1. 搜索条件（用户名或邮箱）
	if search != "" {
		query = query.Where("username LIKE ? OR email LIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// 2. 邮箱验证状态
	if req.EmailVerified != nil {
		query = query.Where("email_verified = ?", *req.EmailVerified)
		fmt.Println(query)
	}

	// 3. 用户角色
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 4. 用户状态
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 5. 创建时间范围
	if !req.CreatedAtFrom.IsZero() {
		query = query.Where("created_at >= ?", req.CreatedAtFrom)
	}
	if !req.CreatedAtTo.IsZero() {
		query = query.Where("created_at <= ?", req.CreatedAtTo)
	}

	// 6. 排序
	query = query.Order("created_at DESC")

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	// 分页查询
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

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(id uint, status string) error {
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
func (s *UserService) UpdateUserRole(id uint, role string) error {
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

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id uint, req models.UpdateUserRequest) (*models.User, error) {
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

	// 保存更新
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
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

func (s *UserService) DeleteUserByID(id uint) error {
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}
	// 删除用户
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
	// 构建SQL查询
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

	// 验证输入
	if err := validateUserInput(req.Username, req.Email, req.Password); err != nil {
		return nil, ErrInvalidInput
	}

	// 检查用户是否已存在
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, ErrUserExists
	}

	// 创建用户
	user := &models.User{
		Username:  req.Username,
		Email:     req.Email,
		Fullname:  req.Fullname,
		Bio:       req.Bio,
		Role:      req.Role,
		Status:    req.Status,
		AvatarURL: &req.AvatarURL,
	}

	// 设置密码
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	// 保存到数据库
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	if req.SendWelcomeEmail == true {
		// todo 接入kafka 发送欢迎email
	}

	return user, nil
}
