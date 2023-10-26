package main

import (
	"rkpbi-go/router"

	"github.com/gin-gonic/gin"
)

func main() {
	routers := gin.Default()
	router.UserRouter(routers)
	router.PhotoRouter(routers)

	routers.Run(":8172")
}