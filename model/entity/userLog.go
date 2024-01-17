package entity

import "time"

type UserLog struct {
	ID        int       `gorm:"column:id;autoIncrement;primaryKey" json:"id,omitempty"`
	UserId    string    `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Action    string    `gorm:"column:action;not null" json:"action,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime;<-:create" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
}

func (u *UserLog) TableName() string {
	return "user_logs"
}
