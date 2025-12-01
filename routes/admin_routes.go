package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/controllers"
	"github.com/junaid9001/spectr_backend/middlewares"
)

func AdminRoutes(r *gin.Engine) {
	db := config.DB
	admin := r.Group("/admin")

	admin.Use(middlewares.AdminAuthMiddleware())

	//user manage (done) postman
	{
		admin.GET("/users", controllers.AllUsers(db))
		admin.PUT("/users/:id/role", controllers.UpdateUserRole(db))
		admin.PUT("/users/:id/status", controllers.UpdateUserStatus(db))

	}

	//product related done postman
	{
		//done postman
		admin.POST("/product", controllers.CreateProduct(db))

		admin.PUT("/product/:id", controllers.UpdateProductByID(db))

		admin.DELETE("/product/:id", controllers.DeleteProductByID(db))

		//product public

		//done Postman
		r.GET("/products", controllers.GetAllProducts(db))

		r.GET("/product/:id", controllers.GetProductByID(db))
	}

	//orders related

	{ //done postman
		admin.GET("/orders", controllers.GetAllOrders(db))
		admin.PATCH("/order/:id", controllers.UpdateOrderStatus(db))
	}

	//category related
	{
		admin.POST("/categories", controllers.AddCategory(db))
		admin.DELETE("/categories/:id", controllers.DeleteCategoryByID(db))
		admin.GET("/categories/tree", controllers.CategoryTree(db))
		admin.GET("/categories", controllers.AllCategories(db))
		admin.PUT("/edit/category/:id", controllers.EditCategoryByID(db))
	}

	//add filters
	{
		admin.POST("/add/filter", controllers.AddFilter(db))
		//add option like male female
		admin.POST("/add/filter_option", controllers.AddFilterOption(db))
		//link product with filteroption
		admin.POST("/link_filter", controllers.AddFilterOptionToProduct(db))

		admin.GET("/filtered_products/:id", controllers.ViewProductsByFilterOptionId(db))
	}

}
