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

	user.GET("/profile", controllers.GetUserProfile(db))
	user.PUT("/profile", controllers.UpdateUserProfile(db))
}
