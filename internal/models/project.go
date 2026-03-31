package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `json:"description,omitempty"`
	TeamID      UUID           `gorm:"type:char(36);not null" json:"team_id"`
	OwnerID     UUID           `gorm:"type:char(36);not null" json:"owner_id"`
	Status      string         `gorm:"size:20;default:'active'" json:"status"` // active, archived, deleted
	Color       string         `gorm:"size:7;default:'#3B82F6'" json:"color"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Team  *Team  `gorm:"foreignKey:TeamID" json:"team,omitempty"`
	Owner *User  `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Tasks []Task `gorm:"foreignKey:ProjectID" json:"tasks,omitempty"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate 创建前的钩子
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.ID.IsZero() {
		p.ID = NewUUID()
	}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 更新前的钩子
func (p *Project) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}
