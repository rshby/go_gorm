package entity

import "time"

type UserLog struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id,omitempty"`
	UserId    string    `gorm:"column:user_id" json:"user_id,omitempty"`
	Action    string    `gorm:"column:action" json:"action,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
}

func (u *UserLog) TableName() string {
	return "user_logs"
}
