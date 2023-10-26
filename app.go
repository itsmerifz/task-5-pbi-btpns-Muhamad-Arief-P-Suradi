package main

import (
	"rkpbi-go/router"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func main() {
	routers := gin.Default()
	router.UserRouter(routers)
	router.PhotoRouter(routers)

	routers.Run(":8172")
}
