package main

import (
	"github.com/gin-gonic/gin"
	"github.com/junaid9001/spectr_backend/config"
)

func main() {
	r := gin.Default()

	r.GET("/test", func(ctx *gin.Context) {
		ctx.String(200, "server is running")
	})

	config.ConnectDB()

	r.Run(":8081")
}
