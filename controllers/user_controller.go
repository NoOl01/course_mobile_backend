package controllers

import (
	"backend_course/common"
	"backend_course/database"
	"backend_course/dto"
	"backend_course/otp"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type UserController struct {
	Db *gorm.DB
}

func (dbc *UserController) Register(c *gin.Context) {
	var newUser dto.RegisterDto

	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	hash := common.Encrypt(newUser.Password)
	user := database.User{

		FirstName: newUser.FirstName,
		Email:     newUser.Email,
		Password:  hash,
		Balance:   0,
	}

	access, refresh, err := common.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	if err := dbc.Db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"access_token":  access,
			"refresh_token": refresh,
		},
		"error": nil,
	})
}

func (dbc *UserController) Login(c *gin.Context) {
	var login dto.LoginDto

	if err := c.ShouldBind(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}
	user, err := common.CheckPass(login, dbc.Db)
	if err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
			return
		}
		if errors.Is(err, common.ErrWrongPassword) {
			c.JSON(http.StatusBadRequest, gin.H{
				"result": nil,
				"error":  err.Error(),
			})
			return
		}
	}

	access, refresh, err := common.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": nil,
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": gin.H{
			"access_token":  access,
			"refresh_token": refresh,
		},
		"error": nil,
	})
}

func (dbc *UserController) SendPasswordResetCode(c *gin.Context) {
	var sendOtp dto.SendOtpDto
	var user bool

	if err := c.ShouldBind(&sendOtp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := dbc.Db.Model(database.User{}).Where("email = ?", sendOtp.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !user {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User %s not found", sendOtp.Email),
		})
		return
	}

	otpCode := otp.OtpGenerate()
	otp.StoreOTP(sendOtp.Email, otpCode)
	err = otp.SendOtp(sendOtp.Email, otpCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *UserController) OtpCheck(c *gin.Context) {
	var resetWithOtp dto.ResetPasswordWithCodeDto
	var user bool

	if err := c.ShouldBind(&resetWithOtp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := dbc.Db.Model(database.User{}).Where("email = ?", resetWithOtp.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if !user {
		c.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("User %s not found", resetWithOtp.Email),
		})
		return
	}

	isVerified := otp.VerifyOTP(resetWithOtp.Email, resetWithOtp.Code)
	if !isVerified {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong otp code",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *UserController) ResetPassword(c *gin.Context) {
	var resetPassword dto.ResetPasswordDto
	var user database.User

	if err := c.ShouldBind(&resetPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dbc.Db.Where("email = ?", resetPassword.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("User %s not found", resetPassword.Email),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	hash := common.Encrypt(resetPassword.Password)

	if err := dbc.Db.Model(&user).Update("password", hash).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func (dbc *UserController) UpdateProfile(c *gin.Context) {

}
