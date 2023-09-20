package initializer

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var Db *gorm.DB

func ConnectToDB() {
	dsn := os.Getenv("DB_URL")
	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error loading db" + err.Error())
	}
}
