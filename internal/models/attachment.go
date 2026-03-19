package models

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TaskID      uint           `gorm:"not null" json:"task_id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Filename    string         `gorm:"size:255;not null" json:"filename"`
	FileSize    int64          `json:"file_size"`
	MimeType    string         `gorm:"size:100" json:"mime_type"`
	StoragePath string         `gorm:"size:500;not null" json:"storage_path"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Attachment) TableName() string {
	return "attachments"
}

func (a *Attachment) BeforeCreate(tx *gorm.DB) error {
	a.CreatedAt = time.Now()
	return nil
}
