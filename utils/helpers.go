package utils

import (
	"errors"
	"time"
	valid "github.com/asaskevich/govalidator"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type CustomJWTClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

var jwtKey = []byte("secret")

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ValidUser(email string, password string) bool {
	if valid.IsEmail(email) && len(password) >= 6 {
		return true
	}
	return false
}

func ValidPhoto(title string, caption string, photoUrl string) bool {
	if valid.IsAlphanumeric(title) && valid.IsAlphanumeric(caption) && valid.IsURL(photoUrl) {
		return true
	}
	return false
}

func GenerateJWT(id string) (token string, err error) {
	expTime := time.Now().Add(time.Hour * 24).Unix()
	claims := &CustomJWTClaims{
		UserID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,
		},
	}
	sign := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
	token, err = sign.SignedString(jwtKey)
	return
}

func ValidateToken(tokenString string) (err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*CustomJWTClaims)
	if !ok {
		err = errors.New("token is invalid")
		return
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = errors.New("token is expired")
		return
	}
	return
}

func GetUserIDFromToken(tokenString string) (id string, err error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomJWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*CustomJWTClaims)
	if !ok {
		err = errors.New("token is invalid")
		return
	}
	id = claims.UserID
	return
}
