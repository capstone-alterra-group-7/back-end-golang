package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName       string
	Email          string `gorm:"unique"`
	Password       string
	PhoneNumber    string
	BirthDate      *time.Time `gorm:"type:DATE"`
	ProfilePicture string
	Citizen        string
	Role           string `gorm:"type:ENUM('user','admin')"`
}
