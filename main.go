package main

import (
	"backend_course/database"
	"backend_course/images"
	"backend_course/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panicf("Error loading .env file. %s\n", err.Error())
	}

	images.CheckFolder("./upload/products")
	images.CheckFolder("./upload/avatars")

	r := gin.Default()
	db := database.Connect()

	r.Use(cors.Default())
	router.AppRouter(r, db)

	if err := r.Run(":8080"); err != nil {
		log.Panicf("Error starting server. %s\n", err.Error())
	}
}
