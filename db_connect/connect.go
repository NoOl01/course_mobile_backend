package db_connect

import (
	"course_mobile/db_models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func Connect() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		panic(".env file not found")
	}

	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	migrateError := db.AutoMigrate(&db_models.User{}, &db_models.Product{}, &db_models.ProductImage{},
		&db_models.ProductFilter{}, &db_models.Category{}, &db_models.Filter{})
	if migrateError != nil {
		panic("failed migrate: " + migrateError.Error())
	}

	return db
}
