package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

type CartController struct {
	Db *gorm.DB
}

func (dbc *CartController) GetAllCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "authorization header is required",
		})
		return
	}
	token := strings.Split(authHeader, " ")[1]

	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "invalid token",
		})
		return
	}

	var cart database.Cart
	id := common.GetIdFromToken(claims)

	if err := dbc.Db.Where("id = ?", id).Find(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"result": nil,
				"error":  "cart not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"result": cart,
		"error":  nil,
	})
}

func (dbc *CartController) AddToCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authorization header is required",
		})
		return
	}

	var cartDto dto.CartDto

	if err := c.ShouldBind(&cartDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	id, err := strconv.ParseInt(common.GetIdFromToken(claims), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	var cart = database.Cart{
		UserId:    id,
		ProductId: cartDto.ProductId,
		Count:     cartDto.Count,
	}

	if err := dbc.Db.Create(&cart).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *CartController) DeleteFromCart(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "authorization header is required",
		})
		return
	}
	token := strings.Split(authHeader, " ")[1]

	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}
	id := common.GetIdFromToken(claims)

	if err := dbc.Db.Where("id = ?", id).Delete(&database.Cart{}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "cart not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
