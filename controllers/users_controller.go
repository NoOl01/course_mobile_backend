package controllers

import (
	"course_mobile/common"
	"course_mobile/db_models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UserRegister(c *gin.Context, db *gorm.DB) {
	var newUser common.NewUser
	var user db_models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": nil,
			"error":   "Bad Request",
		})
		return
	}

	hash := common.Encrypt(newUser.Password)
	user.Email = newUser.Email
	user.Password = hash

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": nil,
			"error":   "Internal Server Error",
		})
		return
	}

	token, err := common.GenerateToken(user, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": nil,
			"error":   "Internal Server Error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": token,
		"error":   nil,
	})
}

func UserLogin(c *gin.Context, db *gorm.DB) {
	var newUser common.NewUser
	var user db_models.User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": nil,
			"error":   "Bad Request",
		})
		return
	}

	checkPass := common.CheckPass(newUser, db)
	if checkPass != "Ok" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": nil,
			"error":   "Wrong password",
		})
		return
	}

	token, err := common.GenerateToken(user, db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": nil,
			"error":   "Internal Server Error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": token,
		"error":   nil,
	})
}

func ForgotPassword(c *gin.Context, db *gorm.DB) {
	var email common.Email
	var user db_models.User

	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	err := db.Find(&user, "email = ?", email.Email).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request: " + err.Error(),
		})
		return
	}

	otpCode := common.OtpGenerate()
	common.StoreOTP(email.Email, otpCode)
	err = common.SendOtp(email.Email, otpCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func VerifyOtp(c *gin.Context) {
	var otp common.Otp

	if err := c.BindJSON(&otp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	verified := common.VerifyOTP(otp.Email, otp.Otp)
	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Wrong otp code",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
