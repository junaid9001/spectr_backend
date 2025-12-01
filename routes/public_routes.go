package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/controllers"
)

func PublicRoutes(r *gin.Engine) {
	db := config.DB
	r.GET("/search", controllers.SearchProduct(db))
	r.GET("/products/filter-by-category_id", controllers.FilterProductByCategoryID(db))
	r.GET("/products/filter-by-brand", controllers.FilterProductByBrand(db))
	r.GET("/products/filter-by-price", controllers.FilterProductByPrice(db))
}
