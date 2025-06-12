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
			user.POST("/refresh", userController.RefreshToken)
			user.POST("/sendOtp", userController.SendPasswordResetCode)
			user.POST("/checkOtp", userController.OtpCheck)
			user.POST("/resetPassword", userController.ResetPassword)
			user.POST("/updateProfile", userController.UpdateProfile)
			user.POST("/updateProfileAvatar", userController.UpdateProfileAvatar)
			user.GET("/getInfo", userController.GetProfileInfo)
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
			order.POST("/buy", orderController.BuyProduct)
		}
		category := api.Group("category")
		{
			categoryController := controllers.CategoryController{Db: db}

			category.GET("/getAll", categoryController.GetAllCategories)
		}
		cart := api.Group("/cart")
		{
			cartController := controllers.CartController{Db: db}

			cart.GET("/getAll", cartController.GetAllCart)
			cart.POST("/add", cartController.AddToCart)
			cart.POST("/delete", cartController.DeleteFromCart)
			cart.POST("/updateCount", cartController.UpdateProductsCount)
		}
		notification := api.Group("/notification")
		{
			notificationController := controllers.NotificationController{Db: db}
			notification.GET("/getAll", notificationController.GetAllNotifications)
		}
		favourite := api.Group("/favourite")
		{
			favouriteController := controllers.FavouriteController{Db: db}

			favourite.POST("/add", favouriteController.AddToFavourite)
			favourite.POST("/delete", favouriteController.DeleteFromFavourite)
			favourite.GET("/getAll", favouriteController.GetAllFavouriteProducts)
		}
		settings := api.Group("/settings")
		{
			settingsController := controllers.SettingsController{Db: db}

			settings.POST("/moneyAction", settingsController.MoneyAction)
		}
	}
}
