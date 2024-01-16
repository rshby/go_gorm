package entity

type Sample struct {
	Id   string `gorm:"column=id" json:"id"`
	Name string `gorm:"column=name" json:"name"`
}
