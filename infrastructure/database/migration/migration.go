package main

import (
	"fmt"
	"go_gorm/config"
	"go_gorm/infrastructure/database/connection"
	"go_gorm/model/entity"
)

func main() {
	// get config
	cfg, _ := config.LoadConfig()
	appConfig := config.ConvertToObject(cfg)

	db := connection.ConnectToDB()

	db.AutoMigrate(&entity.User{}, &entity.Wallet{}, &entity.Address{},
		&entity.Product{}, &entity.UserLog{})

	fmt.Println(appConfig.Database.Name)
}
