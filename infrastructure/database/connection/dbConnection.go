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
func ConnectToDB(config config.IConfig) *gorm.DB {
	config = config.GetConfig()

	user := config.GetConfig().Database.User
	password := config.GetConfig().Database.Password
	host := config.GetConfig().Database.Host
	port := config.GetConfig().Database.Port
	dbName := config.GetConfig().Database.Name

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

	log.Println("success connection to database")
	// success
	return db
}
