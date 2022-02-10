package main

import (
	"cloudDisk/src/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/user/login", controller.LoginHandler)
	router.POST("/user/register", controller.RegisterHandler)
	router.Run(":8080")
}
