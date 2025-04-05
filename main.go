package main

import (
	"course_mobile/db_connect"
	"course_mobile/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db := db_connect.Connect()
	r := gin.Default()

	r.Use(cors.Default())
	router.ApiRouter(r, db)

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
