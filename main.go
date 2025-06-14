package main

import (
	"backend_course/database"
	"backend_course/images"
	"backend_course/router"
	"fmt"
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

	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	db := database.Connect()

	r.Static("/upload/products", "./upload/products")
	r.Static("/upload/avatars", "./upload/avatars")

	r.Use(cors.Default())
	router.AppRouter(r, db)

	fmt.Println("Server started")
	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Panicf("Error starting server. %s\n", err.Error())
	}
}
