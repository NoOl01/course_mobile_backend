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
	"strconv"
	"strings"
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

	access, refresh, err := common.GenerateToken(&user)
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

	access, refresh, err := common.GenerateToken(user)
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

func (dbc *UserController) RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "No Refresh token found",
		})
		return
	}

	tokenData := strings.Split(refreshToken, " ")
	if tokenData[0] != "Refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": nil,
			"error":  "Invalid Refresh token",
		})
		return
	}

	claims, err := common.DecodeToken(tokenData[1])
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
	access, refresh, err := common.Refresh(id)
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
	var user database.User

	if err := c.ShouldBind(&sendOtp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dbc.Db.Where("email = ?", sendOtp.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	otpCode := otp.Generate()
	otp.StoreOTP(sendOtp.Email, otpCode)
	err := otp.SendOtp(sendOtp.Email, otpCode)
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
	var user database.User

	if err := c.ShouldBind(&resetWithOtp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := dbc.Db.Where("email = ?", resetWithOtp.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
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
