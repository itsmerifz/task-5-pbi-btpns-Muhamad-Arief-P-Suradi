package controllers

import (
	"net/http"
	"strconv"

	"rkpbi-go/models"
	"rkpbi-go/utils"

	"github.com/gin-gonic/gin"
)

func (database *Database) CreateUser(c *gin.Context) {
	var (
		user   models.User
		result gin.H
	)

	id := c.PostForm("id")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	parsedId, _ := strconv.ParseUint(id, 10, 32)

	user.ID = uint(parsedId)
	user.Username = username
	user.Email = email
	user.Password = utils.HashPassword(password)

	database.DB.Create(&user)
	result = gin.H{
		"result": user,
	}

	c.JSON(http.StatusOK, result)
}

func (database *Database) CheckUser(c *gin.Context) {
	var user []models.User
	query := database.DB.Find(&user)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": query.Rows})
}
