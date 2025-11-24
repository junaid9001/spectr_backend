package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/models"
)

func ViewRoutes(r *gin.Engine) {

	db := config.DB

	view := r.Group("/view")

	view.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	view.GET("/admin/dashboard", func(c *gin.Context) {

		var stats models.AppStats
		db.First(&stats, 1)

		c.HTML(http.StatusOK, "dashboard.html", gin.H{"data": stats})
	})

	//users
	view.GET("/admin/users", func(c *gin.Context) {

		var users []models.User

		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "users.html", gin.H{"users": users})
	})

	//products
	view.GET("/admin/products", func(c *gin.Context) {

		products := []models.Product{}

		if err := db.Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"error":  err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "products.html", gin.H{"data": products})
	})

	//orders

	view.GET("/admin/orders", func(c *gin.Context) {
		orders := []models.Order{}
		db.Preload("OrderItems").Find(&orders)

		c.HTML(http.StatusOK, "orders.html", gin.H{"data": orders})
	})

}
