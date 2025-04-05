package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApiRouter(r *gin.Engine, db *gorm.DB) {
	test := r.Group("/api")
	{
		test.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}
}
