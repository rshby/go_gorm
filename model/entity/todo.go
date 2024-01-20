package entity

import (
	"gorm.io/gorm"
	"time"
)

type Todo struct {
	ID          int64          `gorm:"column:id;autoIncrement;primaryKey" json:"id,omitempty"`
	UserId      string         `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Title       string         `gorm:"column:title;not null" json:"title,omitempty"`
	Description string         `gorm:"column:description;null" json:"description,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime;not null" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime;autoCreateTime;not null" json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}

func (t *Todo) TableName() string {
	return "todos"
}
