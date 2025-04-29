package controllers

import (
	"backend_course/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ProductController struct {
	Db *gorm.DB
}

func (dbc *ProductController) GetAllProducts(c *gin.Context) {
	var products []database.Product

	if err := dbc.Db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": products,
		"error":  nil,
	})
}

func (dbc *ProductController) GetProductByCategoryId(c *gin.Context) {
	var product []database.Product

	categoryId := c.Query("category_id")
	if categoryId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  "Query \"category_id\" is required",
		})
		return
	}

	if err := dbc.Db.Where("category_id = ?", categoryId).Find(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": product,
		"error":  nil,
	})
}
