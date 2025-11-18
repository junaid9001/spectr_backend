package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/middlewares"
)

func AdminRoutes(r *gin.Engine) {
	admin := r.Group("/admin")

	admin.Use(middlewares.AdminAuthMiddleware())

	//admin/user
	{

	}
}
