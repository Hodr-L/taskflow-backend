package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Type      string         `gorm:"size:50;not null" json:"type"`
	Title     string         `gorm:"size:200;not null" json:"title"`
	Content   string         `json:"content,omitempty"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time      `json:"created_at"`
	ReadAt    *time.Time     `json:"read_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Notification) TableName() string {
	return "notifications"
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	n.CreatedAt = time.Now()
	return nil
}
