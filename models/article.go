package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title       string `json:"title" form:"title"`
	Image       string `json:"image" form:"image"`
	Description string `json:"description" form:"description"`
	Label       string `json:"label" form:"label"`
}
