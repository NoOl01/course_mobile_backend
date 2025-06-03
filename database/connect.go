package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL") + "?parseTime=true"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("Error connecting to database. %s\n", err.Error())
	}

	autoMigrateErr := db.AutoMigrate(&User{}, &Category{}, &Brand{}, &Product{}, &ProductImage{}, &Cart{}, &Favorite{}, &Order{}, &Notification{})
	if autoMigrateErr != nil {
		log.Panicf("Error running migrations. %s\n", autoMigrateErr.Error())
	}

	return db
}
