package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `json:"description,omitempty"`
	TeamID      uint           `gorm:"not null" json:"team_id"`
	OwnerID     uint           `gorm:"not null" json:"owner_id"`
	Status      string         `gorm:"size:20;default:'active'" json:"status"` // active, archived, deleted
	Color       string         `gorm:"size:7;default:'#3B82F6'" json:"color"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴
	Team  *Team `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Owner *User `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	// Tasks   []Task `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// TableName 鎸囧畾琛ㄥ悕
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate 鍒涘缓鍓嶇殑閽╁瓙
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 鏇存柊鍓嶇殑閽╁瓙
func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
