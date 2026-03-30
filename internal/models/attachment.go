package models

import (
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	ID          UUID           `gorm:"type:char(36);primaryKey" json:"id"`
	TaskID      UUID           `gorm:"type:char(36);not null" json:"task_id"`
	UserID      UUID           `gorm:"type:char(36);not null" json:"user_id"`
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
	if a.ID.IsZero() {
		a.ID = NewUUID()
	}
	a.CreatedAt = time.Now()
	return nil
}
