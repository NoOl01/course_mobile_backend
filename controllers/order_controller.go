package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"backend_course/statuses"
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
			"result": nil,
			"error":  "No Authorization Header found",
		})
		return
	}
	var buyProductDto dto.BuyProductDto

	if err := c.ShouldBind(&buyProductDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  err.Error(),
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

	id, err := strconv.ParseInt(common.GetIdFromToken(claims), 10, 64)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "invalid token",
		})
		return
	}

	var product database.Product
	var user database.User
	status, err := common.TryTransaction(dbc.Db, id, buyProductDto.ProductId, &product, &user)
	if err != nil {
		if errors.Is(err, statuses.UserNotFound) || errors.Is(err, statuses.ProductNotFound) {
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

	order := database.Order{
		UserId:    id,
		ProductId: buyProductDto.ProductId,
		Count:     buyProductDto.Count,
		Time:      time.Now(),
	}

	switch status {
	case statuses.ResultAwaitingPayment:
		order.Status = statuses.ResultAwaitingPayment

		if err := dbc.createOrder(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
		}
	case statuses.ResultOk:
		order.Status = statuses.ResultOk

		if err := dbc.Db.Model(&user).Where("id = ?", id).Update("Balance", user.Balance-product.Price).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
			return
		}

		if err := dbc.createOrder(&order); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"result": order,
		"error":  nil,
	})
}

func (dbc *OrderController) createOrder(order *database.Order) error {
	if err := dbc.Db.Create(order).Error; err != nil {
		return err
	}
	return nil
}
