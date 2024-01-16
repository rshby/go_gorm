package entity

import "time"

// table users
type User struct {
	ID          string    `gorm:"column=id;primaryKey;<-:create'" json:"id"`
	Password    string    `gorm:"column=password" json:"password"`
	Name        string    `gorm:"name" json:"name"`
	CreatedAt   time.Time `gorm:"column=created_at;autoCreateTime;<-:create" json:"created_at"`
	UpdatedAt   time.Time `gorm:"updated_at;autoUpdateTime" json:"updated_at"`
	Information string    `gorm:"-"`
}

// set table name
func (u *User) TableName() string {
	return "users"
}
