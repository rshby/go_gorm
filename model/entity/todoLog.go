package entity

import (
	"gorm.io/gorm"
	"time"
)

type TodoLog struct {
	gorm.Model
	UserId      string `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Title       string `gorm:"column:title;not null" json:"title,omitempty"`
	Description string `gorm:"column:description;null" json:"description,omitempty"`
}

func (t *TodoLog) TableName() string {
	return "todo_logs"
}

var s = TodoLog{
	Model: gorm.Model{
		ID:        0,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
	}}
