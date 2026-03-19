package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Email         string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"size:255;not null" json:"-"`
	AvatarURL     *string        `gorm:"size:255" json:"avatar_url,omitempty"`
	Role          string         `gorm:"size:20;default:'user'" json:"role"`     // super_admin, admin, user
	Status        string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive, banned
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系（暂时注释，等模型创建后再启用）
	// Teams        []TeamMember `gorm:"foreignKey:UserID" json:"teams,omitempty"`
	// OwnedTeams   []Team       `gorm:"foreignKey:OwnerID" json:"owned_teams,omitempty"`
	// OwnedProjects []Project    `gorm:"foreignKey:OwnerID" json:"owned_projects,omitempty"`
	// AssignedTasks []Task       `gorm:"foreignKey:AssigneeID" json:"assigned_tasks,omitempty"`
	// ReportedTasks []Task       `gorm:"foreignKey:ReporterID" json:"reported_tasks,omitempty"`
	// Comments      []TaskComment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	// Attachments   []Attachment  `gorm:"foreignKey:UserID" json:"attachments,omitempty"`
	// Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// SetPassword 设置密码（加密存储）
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin 检查是否是管理员
func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.Role == "super_admin"
}

// IsActive 检查是否活跃
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// UpdateLastLogin 更新最后登录时间
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// ToResponse 转换为API响应格式
func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":             u.ID,
		"username":       u.Username,
		"email":          u.Email,
		"avatar_url":     u.AvatarURL,
		"role":           u.Role,
		"status":         u.Status,
		"email_verified": u.EmailVerified,
		"last_login_at":  u.LastLoginAt,
		"created_at":     u.CreatedAt,
	}
}

// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}

// UserChangePasswordRequest 修改密码请求
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	AvatarURL     *string    `json:"avatar_url,omitempty"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token,omitempty"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	User         interface{} `json:"user"`
}

// PaginatedUsersResponse 分页用户响应
type PaginatedUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
