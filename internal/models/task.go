package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Title          string         `gorm:"size:200;not null" json:"title"`
	Description    string         `json:"description,omitempty"`
	ProjectID      uint           `gorm:"not null" json:"project_id"`
	AssigneeID     *uint          `json:"assignee_id,omitempty"`
	ReporterID     uint           `gorm:"not null" json:"reporter_id"`
	Priority       string         `gorm:"size:10;default:'medium'" json:"priority"` // low, medium, high, urgent
	Status         string         `gorm:"size:20;default:'todo'" json:"status"`     // todo, in_progress, review, done, cancelled
	DueDate        *time.Time     `json:"due_date,omitempty"`
	EstimatedHours float64        `json:"estimated_hours,omitempty"`
	ActualHours    float64        `gorm:"default:0" json:"actual_hours"`
	Tags           string         `gorm:"type:json" json:"tags,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	CompletedAt    *time.Time     `json:"completed_at,omitempty"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 鍏宠仈鍏崇郴
	Project  *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Assignee *User    `gorm:"foreignKey:AssigneeID" json:"assignee,omitempty"`
	Reporter *User    `gorm:"foreignKey:ReporterID" json:"reporter,omitempty"`
}

// TableName 鎸囧畾琛ㄥ悕
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate 鍒涘缓鍓嶇殑閽╁瓙
func (t *Task) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate 鏇存柊鍓嶇殑閽╁瓙
func (t *Task) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()

	// 濡傛灉鐘舵€佸彉涓哄畬鎴愶紝璁剧疆瀹屾垚鏃堕棿
	if t.Status == "done" && t.CompletedAt == nil {
		now := time.Now()
		t.CompletedAt = &now
	}

	return nil
}
