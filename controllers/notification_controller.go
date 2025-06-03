package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type NotificationController struct {
	Db *gorm.DB
}

func (dbc *NotificationController) GetAllNotifications(c *gin.Context) {
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

	id, err := common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "invalid token",
		})
		return
	}

	var notifications []database.Notification
	if err := dbc.Db.Where("user_id = ?", id).Find(&notifications).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": notifications,
		"error":  nil,
	})
}
