package dto

type UserResponse struct {
	Id        string `gorm:"column=id" json:"id,omitempty"`
	FirstName string `gorm:"column=first_name" json:"first_name,omitempty"`
	LastName  string `gorm:"last_name" json:"last_name,omitempty"`
}
