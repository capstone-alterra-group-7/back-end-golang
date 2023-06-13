package models

import "gorm.io/gorm"

type TemplateMessage struct {
	gorm.Model
	Name    string `json:"name"`
	Content string `json:"content"`
}
