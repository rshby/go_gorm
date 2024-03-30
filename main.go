package main

import (
	"fmt"
	"go_gorm/config"
	"go_gorm/infrastructure/database/connection"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Println(err.Error())
		log.Fatal()
	}

	// connect to database
	db := connection.ConnectToDB(cfg)
	fmt.Println(db)
}
