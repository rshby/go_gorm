package main

import (
	"fmt"
	"go_gorm/config"
	"go_gorm/infrastructure/database/connection"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Println(err.Error())
		log.Fatal()
	}

	// test config
	fmt.Println(config.GetString("app.name"))

	// connect to database
	db := connection.ConnectToDB()
	fmt.Println(db)
}
