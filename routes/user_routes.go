package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/controllers"
	"github.com/junaid9001/spectr_backend/middlewares"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")

	user.Use(middlewares.UserAuthMiddleware())

	db := config.DB

	//done psotman
	user.GET("/profile", controllers.GetUserProfile(db))
	user.PUT("/profile", controllers.UpdateUserProfile(db))

	{ //user cart related (done) postman
		user.POST("/cart", controllers.AddProductToCart(db))
		user.GET("/cart", controllers.GetUserCart(db))
		user.PATCH("/cart/:id", controllers.UpdateQuantityInCartByID(db))
		user.DELETE("/cart/:id", controllers.DeleteCartItemByID(db))
	}

	{ //wishlist related (done) postman
		user.POST("/wishlist", controllers.AddToWishlist(db))
		user.GET("/wishlist", controllers.GetWishList(db))
		user.DELETE("/wishlist/:product_id", controllers.DeleteFromWishList(db))
	}

	{ //order related (done) postman
		user.POST("/order", controllers.PlaceOrder(db))
		user.GET("/orders", controllers.GetOrderHistory(db))
		user.GET("/order/:id", controllers.GetDetailsOfOrder(db))
		user.DELETE("/order/:id", controllers.DeleteOrderById(db))
		user.PATCH("/order/:id/cancel", controllers.CancelOrderAndRestock(db))

	}

	{ //payment

		user.POST("/order/:id/payments", controllers.CreatePayment(db))
		user.POST("/payment/:payment_id/confirm", controllers.ConfirmPayment(db))

	}
}
