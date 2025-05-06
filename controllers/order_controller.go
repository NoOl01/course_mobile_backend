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
	"time"
)

type OrderController struct {
	Db *gorm.DB
}

func (dbc *OrderController) GetAllOrders(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  "No Authorization Header found",
		})
		return
	}
	token := strings.Split(authHeader, " ")[1]

	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	var orders []database.Order

	id := common.GetIdFromToken(claims)

	if err := dbc.Db.Where("UserId = ?", id).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": orders,
		"error":  nil,
	})
}

func (dbc *OrderController) BuyProduct(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No Authorization Header found",
		})
		return
	}
	var buyProductDto dto.BuyProductDto

	if err := c.ShouldBind(&buyProductDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token := strings.Split(authHeader, " ")[1]
	claims, err := common.DecodeToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
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

	order := database.Order{
		UserId:    id,
		ProductId: buyProductDto.ProductId,
		Count:     buyProductDto.Count,
		Status:    "",
		Time:      time.Now(),
	}
}
