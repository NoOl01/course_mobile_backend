package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type SettingsController struct {
	Db *gorm.DB
}

func (dbc *SettingsController) MoneyAction(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "No Authorization Header found",
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

	id, err := common.GetIdFromToken(claims)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	var body dto.SettingMoneyAction
	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var user database.User
	if err := dbc.Db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch body.Action {
	case "add":
		user.Balance += body.Money
	case "remove":
		user.Balance -= body.Money
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid action",
		})
		return
	}

	if err := dbc.Db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
