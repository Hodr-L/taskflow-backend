package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"size:50;uniqueIndex;not null" json:"username"`
	Fullname      string         `gorm:"size:50;" json:"fullname"`
	Email         string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"size:255;not null" json:"-"`
	AvatarURL     *string        `gorm:"size:255" json:"avatar_url,omitempty"`
	Role          string         `gorm:"size:20;default:'user'" json:"role"`     // super_admin, admin, user
	Status        string         `gorm:"size:20;default:'active'" json:"status"` // active, inactive, banned
	Bio           string         `gorm:"size:255'" json:"bio"`                   // active, inactive, banned
	EmailVerified bool           `gorm:"default:false" json:"email_verified"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴锛堟殏鏃舵敞閲婏紝绛夋ā鍨嬪垱寤哄悗鍐嶅惎鐢級
	// Teams        []TeamMember `gorm:"foreignKey:UserID" json:"teams,omitempty"`
	// OwnedTeams   []Team       `gorm:"foreignKey:OwnerID" json:"owned_teams,omitempty"`
	// OwnedProjects []Project    `gorm:"foreignKey:OwnerID" json:"owned_projects,omitempty"`
	// AssignedTasks []Task       `gorm:"foreignKey:AssigneeID" json:"assigned_tasks,omitempty"`
	// ReportedTasks []Task       `gorm:"foreignKey:ReporterID" json:"reported_tasks,omitempty"`
	// Comments      []TaskComment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
	// Attachments   []Attachment  `gorm:"foreignKey:UserID" json:"attachments,omitempty"`
	// Notifications []Notification `gorm:"foreignKey:UserID" json:"notifications,omitempty"`
}

// TableName 鎸囧畾琛ㄥ悕
func (User) TableName() string {
	return "users"
}

// BeforeCreate 鍒涘缓鍓嶇殑閽╁瓙
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 鏇存柊鍓嶇殑閽╁瓙
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// SetPassword 璁剧疆瀵嗙爜锛堝姞瀵嗗瓨鍌級
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword 楠岃瘉瀵嗙爜
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// IsAdmin 妫€鏌ユ槸鍚︽槸绠＄悊鍛?
func (u *User) IsAdmin() bool {
	return u.Role == "admin" || u.Role == "super_admin"
}

// IsActive 妫€鏌ユ槸鍚︽椿璺?
func (u *User) IsActive() bool {
	return u.Status == "active"
}

// UpdateLastLogin 鏇存柊鏈€鍚庣櫥褰曟椂闂?
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLoginAt = &now
}

// ToResponse 杞崲涓篈PI鍝嶅簲鏍煎紡
func (u *User) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":             u.ID,
		"username":       u.Username,
		"email":          u.Email,
		"avatar_url":     u.AvatarURL,
		"role":           u.Role,
		"fullName":       u.Fullname,
		"bio":            u.Bio,
		"status":         u.Status,
		"email_verified": u.EmailVerified,
		"last_login_at":  u.LastLoginAt,
		"created_at":     u.CreatedAt,
		"updated_at":     u.UpdatedAt,
	}
}

// UserCreateRequest 鍒涘缓鐢ㄦ埛璇锋眰
type UserCreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=100"`
}

// UserLoginRequest 鐢ㄦ埛鐧诲綍璇锋眰
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// UpdateProfileRequest 鐢ㄦ埛鏇存柊鑷繁璇锋眰
type UpdateProfileRequest struct {
	Username  *string `json:"username,omitempty" binding:"omitempty,min=3,max=50"`
	AvatarURL *string `json:"avatar_url,omitempty" binding:"omitempty,url"`
}

// UpdateProfileRequest 鏇存柊鐢ㄦ埛璇锋眰
type UpdateUserRequest struct {
	UpdateProfileRequest
	Email          *string `json:"email,omitempty" binding:"omitempty,email"`
	Role           *string `json:"role,omitempty" binding:"omitempty,oneof=user admin super_admin"`
	Status         *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive banned"`
	Email_verified *bool   `json:"email_verified,omitempty" binding:"omitempty"`
	Fullname       *string `json:"fullname,omitempty" binding:"omitempty,min=3,max=50"`
	Bio            *string `json:"bio,omitempty" binding:"omitempty,min=3,max=100"`
}

// UserChangePasswordRequest 淇敼瀵嗙爜璇锋眰
type UserChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=100"`
}

// GetUsersParams 鑾峰彇鐢ㄦ埛鍒楄〃鏌ヨ鍙傛暟
type GetUsersParams struct {
	Page          int       `form:"page,default=1" binding:"min=1"`
	Limit         int       `form:"limit,default=20" binding:"min=1,max=100"`
	Search        string    `form:"search" binding:"omitempty,max=100"`
	Role          string    `form:"role" binding:"omitempty,oneof=user admin super_admin"`
	Status        string    `form:"status" binding:"omitempty,oneof=active inactive banned"`
	EmailVerified *bool     `form:"email_verified"`
	CreatedAtFrom time.Time `form:"created_at_from" binding:"omitempty" time_format:"2006-01-02T15:04:05Z"`
	CreatedAtTo   time.Time `form:"created_at_to" binding:"omitempty" time_format:"2006-01-02T15:04:05Z"`
}

// Pagination 鍒嗛〉淇℃伅缁撴瀯浣?
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

// UserListResponse 鐢ㄦ埛鍒楄〃鍝嶅簲缁撴瀯浣擄紙浣跨敤 data 瀛楁锛?
type UserListResponse struct {
	User       interface{} `json:"user"`
	Pagination Pagination  `json:"Pagination"`
}

type ResetPasswordRequest struct {
	NewPassword      string `json:"new_password" binding:"required,min=6,max=100"`
	SendNotification bool   `json:"send_notification" binding:"omitempty"`
}

type UserStats struct {
	Total      int64 `json:"total"`
	Active     int64 `json:"active"`
	Inactive   int64 `json:"inactive"`
	Banned     int64 `json:"banned"`
	Admin      int64 `json:"admin"`
	SuperAdmin int64 `json:"super_admin"`
	Unverified int64 `json:"unverified"`
}

// CreateUserRequest 鍒涘缓鐢ㄦ埛璇锋眰缁撴瀯浣?
type CreateUserRequest struct {
	// 蹇呭～瀛楁
	Username string `json:"username" binding:"required,min=3,max=50" validate:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" binding:"required,min=6,max=100" validate:"required,min=6,max=100"`

	// 鍙€夊瓧娈?
	Fullname string `json:"fullname,omitempty" binding:"omitempty,min=1,max=100" validate:"omitempty,min=1,max=100"`
	Bio      string `json:"bio,omitempty" binding:"omitempty,max=500" validate:"omitempty,max=500"`

	// 鏋氫妇瀛楁
	Role   string `json:"role,omitempty" binding:"omitempty,oneof=user admin super_admin" validate:"omitempty,oneof=user admin super_admin"`
	Status string `json:"status,omitempty" binding:"omitempty,oneof=active inactive banned" validate:"omitempty,oneof=active inactive banned"`

	// URL瀛楁
	AvatarURL string `json:"avatar_url,omitempty" binding:"omitempty,url" validate:"omitempty,url" `

	// 甯冨皵瀛楁
	SendWelcomeEmail bool `json:"send_welcome_email,omitempty" binding:"omitempty"`
}
