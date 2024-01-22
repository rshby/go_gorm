package connection

import (
	"fmt"
	"go_gorm/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

// function connection menggunakan gorm
func ConnectToDB() *gorm.DB {
	config, _ := config.LoadConfig()

	user := config.GetString("database.user")
	password := config.GetString("database.password")
	host := config.GetString("database.host")
	port := config.GetInt("database.port")
	dbName := config.GetString("database.name")

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal()
	}

	sqlDb, err := db.DB()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal()
	}

	sqlDb.SetMaxOpenConns(100)
	sqlDb.SetMaxIdleConns(50)
	sqlDb.SetConnMaxLifetime(1 * time.Hour)
	sqlDb.SetConnMaxIdleTime(40 * time.Minute)

	// success
	return db
}
