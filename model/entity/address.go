package entity

import "time"

type Address struct {
	ID        int64     `gorm:"column:id;not null;primary_key;autoIncrement" json:"id,omitempty"`
	UserId    string    `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Address   string    `gorm:"column:address;not null" json:"address,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`

	User *User `gorm:"foreignKey:user_id;references:id"`
}

func (a *Address) TableName() string {
	return "addresses"
}
