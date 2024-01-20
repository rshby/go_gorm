package entity

import "time"

type Wallet struct {
	Id        string    `gorm:"column:id;primaryKey;not null" json:"id,omitempty"`
	UserId    string    `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Balance   int64     `gorm:"column:balance;not null" json:"balance,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	User      *User     `gorm:"foreignKey:user_id;references:id" json:"user,omitempty"`
}

func (w *Wallet) TableName() string {
	return "wallets"
}
