package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID        uint   `gorm:"primaryKey,not null, autoIncrement" json:"id,omitempty"`
	Username  string `gorm:"not null" json:"username,omitempty"`
	Email     string `gorm:"unique,not null" json:"email,omitempty"`
	Password  string `gorm:"not null,size:>6" json:"password,omitempty"`
	Photo     Photo  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Photo struct {
	gorm.Model
	ID       uint `gorm:"primaryKey;not null"`
	Title    string
	Caption  string
	PhotoUrl string
	UserID   uint
	CreatedAt time.Time
	UpdatedAt time.Time
}
