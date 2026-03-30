package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	UserID    UUID           `gorm:"type:char(36);not null" json:"user_id"`
	Type      string         `gorm:"size:50;not null" json:"type"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `json:"content,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID.IsZero() {
		n.ID = NewUUID()
	}
	n.CreatedAt = time.Now()
	return nil
}
