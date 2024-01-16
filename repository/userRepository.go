package repository

import (
	"context"
	"go_gorm/model/entity"
	_interface "go_gorm/repository/interface"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// function provider
func NewUserRepository(db *gorm.DB) _interface.IUserRepository {
	return &UserRepository{db}
}

func (u UserRepository) Insert(ctx context.Context, entity *entity.User) (*entity.User, error) {
	//TODO implement me
	panic("implement me")
}
