package router

import (
	"course_mobile/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApiRouter(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", func(c *gin.Context) {
				controllers.UserRegister(c, db)
			})
			users.POST("/login", func(c *gin.Context) {
				controllers.UserLogin(c, db)
			})
			users.POST("/forgot_password", func(c *gin.Context) {
				controllers.ForgotPassword(c, db)
			})
			users.POST("/verify_otp", func(c *gin.Context) {
				controllers.VerifyOtp(c)
			})
		}
	}
}
