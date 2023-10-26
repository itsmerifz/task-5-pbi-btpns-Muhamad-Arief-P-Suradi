package utils

import (
	"golang.org/x/crypto/bcrypt"
	valid "github.com/asaskevich/govalidator"
)

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
