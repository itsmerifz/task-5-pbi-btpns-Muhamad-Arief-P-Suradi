package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint		  `gorm:"primaryKey,not null" json:"id,omitempty"`
	Username  string 		`gorm:"not null" json:"username,omitempty"`
	Email     string    `gorm:"unique,not null" json:"email,omitempty"`
	Password  string    `gorm:"not null,size:>6" json:"password,omitempty"`
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
