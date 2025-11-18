package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ViewRoutes(r *gin.Engine) {
	r.GET("/login", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login page",
		})
	})
}
