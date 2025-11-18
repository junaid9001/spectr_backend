package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	AuthRoutes(r)
	UserRoutes(r)
	ViewRoutes(r)
}
