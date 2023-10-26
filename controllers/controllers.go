package controllers

import (
	"net/http"
	"strconv"

	"rkpbi-go/models"
	"rkpbi-go/utils"

	jwt "github.com/dgrijalva/jwt-go"
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

	// Validate input
	parsedId, _ := strconv.ParseUint(id, 10, 32)
	if utils.ValidUser(email, password) {
		user.ID = uint(parsedId)
		user.Username = username
		user.Email = email
		user.Password = utils.HashPassword(password)

		// Check avability username
		query := database.DB.First(&user, "username = ?", username)

		if query.RowsAffected > 0 {
			c.JSON(http.StatusNotAcceptable, gin.H{
				"status":  http.StatusNotAcceptable,
				"message": "Username already taken",
			})
			c.Abort()
			return
		}

		database.DB.Create(&user)
		result = gin.H{
			"result": user,
		}

		c.JSON(http.StatusOK, result)
		c.Abort()
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Unvalid input",
		})
		c.Abort()
		return
	}
}

func (database *Database) DeleteUser(c *gin.Context) {
	var(
		user models.User
		result gin.H
	)

	id := c.Param("userId")
	parsedId, _ := strconv.ParseUint(id, 10, 32)

	query := database.DB.First(&user, "id = ?", parsedId)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message": "Internal server error",
		})
		c.Abort()
		return
	} else {
		if query.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"message": "User not found",
			})
			c.Abort()
			return
		} else {
			database.DB.Where("id = ?", parsedId).Unscoped().Delete(&user)
			result = gin.H{
				"status": http.StatusOK,
				"message": "User deleted",
			}
			c.JSON(http.StatusOK, result)
			c.Abort()
			return
		}
	}
}

func (database *Database) CheckUser(c *gin.Context) {
	var user []models.User
	query := database.DB.Find(&user)

	if query.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": user})
}

func (database *Database) LoginHandler(c *gin.Context) {
	users := models.User{}
	email := c.PostForm("email")
	password := c.PostForm("password")
	

	if !utils.ValidUser(email, password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid credentials",
		})
	}

	// Search user in database
	userResult := database.DB.Select("username", "password", "email").First(&users, "email = ?", email)
	if userResult.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User not found",
		})
		c.Abort()
		return
	} else {
		if users.Username == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid credentials",
			})
			c.Abort()
			return
		} else {
			// Compare password
			if !utils.CheckPasswordHash(password, users.Password) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid credentials",
				})
				c.Abort()
				return
			}
		}
	}

	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), jwt.MapClaims{
		"user": users,
	})
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal server error",
		})
		c.Abort()
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Login success",
			"token":   token,
		})
		return
	}

}
