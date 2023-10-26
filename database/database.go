package database

import (
	"rkpbi-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dbUrl := "root:@tcp(127.0.0.1:3306)/rkpbi?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(models.User{}, models.Photo{})
	return db
}
