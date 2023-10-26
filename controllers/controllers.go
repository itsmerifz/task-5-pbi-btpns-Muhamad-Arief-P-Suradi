package controllers

import (
	"net/http"
	"strconv"

	"rkpbi-go/models"
	"rkpbi-go/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
		query := database.DB.First(&user)

		if query.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
			return
		} else {
			if query.RowsAffected > 0 {
				c.JSON(http.StatusNotAcceptable, gin.H{
					"status":  http.StatusNotAcceptable,
					"message": "Username already taken",
				})
				c.Abort()
				return
			}
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
	var user Credentials
	users := models.User{}
	err := c.Bind(&user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid request",
		})
	}

	// Search user in database
	userResult := database.DB.First(&users, "email = ?", user.Email)
	if userResult != nil {
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
			if !utils.CheckPasswordHash(user.Password, users.Password) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid credentials",
				})
				c.Abort()
				return
			}
		}
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
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
