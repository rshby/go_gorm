package entity

import "time"

// User adalah struct representasi tabel users
type User struct {
	ID          string    `gorm:"column:id;primaryKey;<-:create" json:"id,omitempty"`
	Password    string    `gorm:"column:password" json:"password,omitempty"`
	Name        Name      `gorm:"embedded"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime;<-:create" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	Information string    `gorm:"-" json:"information,omitempty"` // ignore column
}

// TableName digunakan untuk override nama table entity apabila convention
// berbeda dengan default konfigurasi gorm
func (u *User) TableName() string {
	return "users"
}
