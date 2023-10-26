package router

import (
	"rkpbi-go/controllers"
	"rkpbi-go/database"
	"rkpbi-go/utils"

	"github.com/gin-gonic/gin"
)

func UserRouter(route *gin.Engine) {
	router := route.Group("/users")
	db := database.InitDB()
	userController := &controllers.Database{DB: db}

	router.POST("/register", userController.CreateUser)
	router.PUT("/:userId", userController.CreateUser)
	router.DELETE("/:userId", userController.CreateUser)
	router.POST("/login", userController.LoginHandler)
	router.GET("/", userController.CheckUser)
}

func PhotoRouter(route *gin.Engine) {
	router := route.Group("/photos")
	db := database.InitDB()
	photoController := &controllers.Database{DB: db}

	router.GET("/", photoController.CreateUser)
	router.POST("/", utils.AuthMiddleware, photoController.CreateUser)
	router.PUT("/:photoId", utils.AuthMiddleware, photoController.CreateUser)
	router.DELETE("/:photoId", utils.AuthMiddleware, photoController.CreateUser)
}


