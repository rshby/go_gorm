package _interface

import (
	"context"
	"go_gorm/model/entity"
)

type IUserRepository interface {
	Insert(ctx context.Context, entity *entity.User) (*entity.User, error)
}
