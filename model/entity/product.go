package entity

import "time"

type Product struct {
	ID          string    `gorm:"column:id;not null;primaryKey" json:"id,omitempty"`
	Name        string    `gorm:"column:name;not null" json:"name,omitempty"`
	Price       int64     `gorm:"column:price;not null" json:"price,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoCreateTime;autoUpdateTime" json:"updated_at,omitempty"`
	LikeByUsers []User    `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:product_id;references:id;joinReferences:user_id" json:"like_by_users,omitempty"`
}

func (p *Product) TableName() string {
	return "products"
}
