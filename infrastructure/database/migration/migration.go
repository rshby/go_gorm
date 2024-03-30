package main

import (
	"go_gorm/config"
	"go_gorm/infrastructure/database/connection"
	"go_gorm/model/entity"
	"log"
)

func main() {
	// get config
	cfg, _ := config.LoadConfig()

	db := connection.ConnectToDB(cfg)

	db.AutoMigrate(&entity.User{}, &entity.Wallet{}, &entity.Address{},
		&entity.Product{}, &entity.UserLog{})

	log.Printf("success migration!\n")
}
