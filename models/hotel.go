package models

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	Name        string `form:"name" json:"name"`
	Class       string `form:"class" json:"class"`
	Description string `form:"description" json:"description"`
	PhoneNumber string `form:"phone_number" json:"phone_number"`
	Email       string `form:"email" json:"email"`
	Address     string `form:"address" json:"address"`
}
