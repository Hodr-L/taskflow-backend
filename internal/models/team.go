package models

import (
	"time"

	"gorm.io/gorm"
)

type Team struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `json:"description,omitempty"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	AvatarURL   *string        `gorm:"size:255" json:"avatar_url,omitempty"`
	InviteCode  *string        `gorm:"size:50;uniqueIndex" json:"invite_code,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴
	Owner    *User        `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members  []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
	Projects []Project    `gorm:"foreignKey:TeamID" json:"projects,omitempty"`
}

type TeamMember struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TeamID    uint           `gorm:"not null" json:"team_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Role      string         `gorm:"size:20;default:'member'" json:"role"` // owner, admin, member
	JoinedAt  time.Time      `json:"joined_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴
	Team *Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 鎸囧畾琛ㄥ悕
func (Team) TableName() string {
	return "teams"
}

func (TeamMember) TableName() string {
	return "team_members"
}

// BeforeCreate 鍒涘缓鍓嶇殑閽╁瓙
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	tm.JoinedAt = time.Now()
	return nil
}

// BeforeUpdate 鏇存柊鍓嶇殑閽╁瓙
func (t *Team) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
