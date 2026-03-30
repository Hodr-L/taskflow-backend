package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `json:"description,omitempty"`
	OwnerID     UUID           `gorm:"type:char(36);not null" json:"owner_id"`
	LogoURL     *string        `gorm:"size:255" json:"logo_url,omitempty"`
	InviteCode  *string        `gorm:"size:50;uniqueIndex" json:"invite_code,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Owner    *User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members  []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Projects []Project    `gorm:"foreignKey:TeamID" json:"projects,omitempty"`
}

type TeamMember struct {
	ID        UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	TeamID    UUID           `gorm:"type:char(36);not null" json:"team_id"`
	UserID    UUID           `gorm:"type:char(36);not null" json:"user_id"`
	Role      string         `gorm:"size:20;default:'member'" json:"role"` // owner, admin, member
	JoinedAt  time.Time      `json:"joined_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Team *Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Team) TableName() string {
	return "teams"
}

func (TeamMember) TableName() string {
	return "team_members"
}

// BeforeCreate 创建前的钩子
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	if t.ID.IsZero() {
		t.ID = NewUUID()
	}
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	if tm.ID.IsZero() {
		tm.ID = NewUUID()
	}
	tm.JoinedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}

type CreateTeamParams struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
}

// ToResponse 转换为API响应格式
func (t *Team) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"name":        t.Name,
		"description": t.Description,
		"owner_id":    t.OwnerID,
		"logo_url":    t.LogoURL,
		"created_at":  t.CreatedAt,
		"updated_at":  t.UpdatedAt,
	}
}
