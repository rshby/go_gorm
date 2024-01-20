package entity

import "time"

// table users
type User struct {
	ID          string    `gorm:"column=id;primaryKey;<-:create'" json:"id,omitempty"`
	Password    string    `gorm:"column=password" json:"password,omitempty"`
	Name        Name      `gorm:"embedded" json:"name"`
	CreatedAt   time.Time `gorm:"column=created_at;autoCreateTime;<-:create" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"updated_at;autoUpdateTime" json:"updated_at,omitempty"`
	Information string    `gorm:"-" json:"information,omitempty"`

	Wallet    *Wallet   `gorm:"foreignKey:UserId;references:ID"`
	Addresses []Address `gorm:"foreignKey:UserId;references:ID"`
}

// set table name
func (u *User) TableName() string {
	return "users"
}
