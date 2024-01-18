package entity

type UserLog struct {
	ID        int    `gorm:"column:id;autoIncrement;primaryKey" json:"id,omitempty"`
	UserId    string `gorm:"column:user_id;not null" json:"user_id,omitempty"`
	Action    string `gorm:"column:action;not null" json:"action,omitempty"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli;<-:create" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli" json:"updated_at,omitempty"`
}

func (u *UserLog) TableName() string {
	return "user_logs"
}
