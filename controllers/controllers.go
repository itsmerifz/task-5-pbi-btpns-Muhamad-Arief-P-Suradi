package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"rkpbi-go/models"
	"rkpbi-go/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
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
		query := database.DB.First(&user, "username = ?", username)

		if query.RowsAffected > 0 {
			c.AbortWithStatusJSON(http.StatusNotAcceptable, gin.H{
				"status":  http.StatusNotAcceptable,
				"message": "Username already taken",
			})
			return
		}

		database.DB.Create(&user)
		result = gin.H{
			"status":  http.StatusCreated,
			"message": "User created",
		}

		c.JSON(http.StatusCreated, result)
		return
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Unvalid input",
		})
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"message": "Internal server error",
		})
		return
	} else {
		if query.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"message": "User not found",
			})
			return
		} else {
			database.DB.Where("id = ?", parsedId).Unscoped().Delete(&user)
			result = gin.H{
				"status": http.StatusOK,
				"message": "User deleted",
			}
			c.JSON(http.StatusOK, result)
			return
		}
	}
}

func (database *Database) CheckUser(c *gin.Context) {
	var user []models.User
	query := database.DB.Find(&user)

	if query.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": user})
}

func (database *Database) LoginHandler(c *gin.Context) {
	users := models.User{}
	email := c.PostForm("email")
	password := c.PostForm("password")
	

	if !utils.ValidUser(email, password) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Invalid credentials",
		})
		return
	}

	// Search user in database
	userResult := database.DB.Select("id", "username", "password", "email").First(&users, "email = ?", email)
	if userResult.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "User not found",
		})
		return
	} else {
		if users.Username == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": "Invalid credentials",
			})
			return
		} else {
			// Compare password
			if !utils.CheckPasswordHash(password, users.Password) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"status":  http.StatusUnauthorized,
					"message": "Invalid credentials",
				})
				return
			}
		}
	}
	// Create token
	token, err := utils.GenerateJWT(strconv.Itoa(int(users.ID)))
	
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Internal server error",
		})
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

func (database *Database) UpdateUser(c *gin.Context) {
	var (
		user   models.User
		result gin.H
	)

	id := c.Param("userId")
	parsedId, _ := strconv.ParseUint(id, 10, 32)

	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	// Validate input
	if utils.ValidUser(email, password) {
		user.ID = uint(parsedId)
		user.Username = username
		user.Email = email
		user.Password = utils.HashPassword(password)

		// Check avability username
		err := database.DB.First(&user, "id = ?", parsedId).Error

		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "User not found",
			})
			return
		}

		// Update user
		database.DB.Model(&user).Updates(models.User{Username: username, Email: email, Password: utils.HashPassword(password)})

		result = gin.H{
			"status":  http.StatusOK,
			"message": "User updated",
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

func (database *Database) GetPhoto(c *gin.Context) {
	var photo []models.Photo
	query := database.DB.Find(&photo)

	if query.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": query.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": photo})
}

func (database *Database) CreatePhoto(c *gin.Context) {
	// Check Authorization
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// Get ID
	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	id, err := utils.GetUserIDFromToken(auth)
	parsedId, _ := strconv.ParseUint(id, 10, 32)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}else{
		var (
			photo  models.Photo
			result gin.H
		)

		// Get form value
		title := c.PostForm("title")
		caption := c.PostForm("caption")
		photoUrl := c.PostForm("photoUrl")

		// Validate input
		if !utils.ValidPhoto(title, caption, photoUrl) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unvalid input",
			})
			return
		}

		// Create photo
		photo.Title = title
		photo.Caption = caption
		photo.PhotoUrl = photoUrl
		photo.UserID = uint(parsedId)

		database.DB.Create(&photo)
		result = gin.H{
			"status":  http.StatusOK,
			"message": "Photo created",
		}

		c.JSON(http.StatusOK, result)
		return
	}
}

func (database *Database) UpdatePhoto(c *gin.Context) {
	// Check Authorization
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// Get ID
	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	uid, err := utils.GetUserIDFromToken(auth)
	parsedUid, _ := strconv.ParseUint(uid, 10, 32)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}else{
		var (
			photo  models.Photo
			result gin.H
		)

		id := c.Param("photoId")
		parsedId, _ := strconv.ParseUint(id, 10, 32)

		err := database.DB.Having("user_id", parsedUid).First(&photo, "id = ?", parsedId).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  http.StatusNotFound,
				"message": "Photo not found",
			})
			return
		}

		// Get form value
		title := c.PostForm("title")
		caption := c.PostForm("caption")
		photoUrl := c.PostForm("photoUrl")

		// Validate input
		if !utils.ValidPhoto(title, caption, photoUrl) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Unvalid input",
			})
			return
		}

		// Update photo
		err = database.DB.Model(&photo).Updates(models.Photo{Title: title, Caption: caption, PhotoUrl: photoUrl, UserID: uint(parsedId)}).Error
		if err != nil {
			result = gin.H{
				"status":  http.StatusInternalServerError,
				"message": "Internal server error",
			}
		} else {
			result = gin.H{
				"status":  http.StatusOK,
				"message": "Photo updated",
			}
			c.JSON(http.StatusOK, result)
			return
		}
	}
}

func (database *Database) DeletePhoto(c *gin.Context) {
	// Check Authorization
	auth := c.Request.Header.Get("Authorization")

	if auth == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// Get ID
	splitToken := strings.Split(auth, "Bearer ")
	auth = splitToken[1]

	uid, err := utils.GetUserIDFromToken(auth)

	if err != nil && uid == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}else{
		var (
			photo  models.Photo
			result gin.H
		)

		id := c.Param("photoId")
		parsedId, _ := strconv.ParseUint(id, 10, 32)

		query := database.DB.First(&photo, "id = ?", parsedId)
		if query.Error != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"message": "Internal server error",
			})
			return
		} else {
			if query.RowsAffected == 0 {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"status": http.StatusNotFound,
					"message": "Photo not found",
				})
				return
			} else {
				database.DB.Where("id = ?", parsedId).Unscoped().Delete(&photo)
				result = gin.H{
					"status": http.StatusOK,
					"message": "Photo deleted",
				}
				c.JSON(http.StatusOK, result)
				return
			}
		}
	}
}