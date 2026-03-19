package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskComment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TaskID    uint           `gorm:"not null" json:"task_id"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Task *Task `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (TaskComment) TableName() string {
	return "task_comments"
}

func (tc *TaskComment) BeforeCreate(tx *gorm.DB) error {
	tc.CreatedAt = time.Now()
	tc.UpdatedAt = time.Now()
	return nil
}

func (tc *TaskComment) BeforeUpdate(tx *gorm.DB) error {
	tc.UpdatedAt = time.Now()
	return nil
}
