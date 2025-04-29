package controllers

import (
	"backend_course/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type CategoryController struct {
	Db *gorm.DB
}

func (dbc *CategoryController) GetAllCategories(c *gin.Context) {
	var categories []database.Category

	if err := dbc.Db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": categories,
		"error":  nil,
	})
}
