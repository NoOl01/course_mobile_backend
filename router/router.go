package router

import (
	"backend_course/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AppRouter(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api/v1")
	{
		user := api.Group("/user")
		{
			userController := controllers.UserController{Db: db}

			user.POST("/register", userController.Register)
			user.POST("/login", userController.Login)
			user.POST("/sendOtp", userController.SendPasswordResetCode)
			user.POST("/checkOtp", userController.OtpCheck)
			user.POST("/resetPassword", userController.ResetPassword)
		}
		product := api.Group("/product")
		{
			productController := controllers.ProductController{Db: db}

			product.GET("/getAll", productController.GetAllProducts)
			product.GET("/getByCategoryId", productController.GetProductByCategoryId)
			product.GET("/getById", productController.GetProductInfoById)
		}
		order := api.Group("/order")
		{
			orderController := controllers.OrderController{Db: db}

			order.GET("/getAll", orderController.GetAllOrders)
		}
		cart := api.Group("/cart")
		{
			cartController := controllers.CartController{Db: db}

			cart.GET("/getAll", cartController.GetAllCart)
			cart.POST("/add", cartController.AddToCart)
			cart.POST("/delete", cartController.DeleteFromCart)
		}
	}
}
