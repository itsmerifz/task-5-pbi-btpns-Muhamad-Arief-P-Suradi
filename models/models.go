package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint `gorm:"primaryKey,not null"`
	Username  string
	Email     string    `gorm:"unique,not null"`
	Password  string    `gorm:"not null,size:>6"`
	Photo     Photo     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Photo struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Title    string
	Caption  string
	PhotoUrl string
	UserID   uint
}
