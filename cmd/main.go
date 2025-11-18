package main

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
	"github.com/junaid9001/spectr_backend/routes"
)

func main() {

	//database setups
	config.LoadEnv()
	config.ConnectDB()
	config.MigrateAll()

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	//register all routes
	routes.RegisterRoutes(r)

	r.GET("/signup", func(ctx *gin.Context) {
		ctx.HTML(200, "signup.html", "")
	})

	r.Run(":8080")
}
