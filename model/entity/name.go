package entity

type Name struct {
	FirstName  string `gorm:"column=first_name;not null" json:"first_name,omitempty"`
	MiddleName string `gorm:"column=middle_name;null" json:"middle_name,omitempty"`
	LastName   string `gorm:"column=last_name;null" json:"last_name,omitempty"`
}
