package models

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	Type          string
	ImageUrl      string
	Name          string
	AccountName   string
	AccountNumber string
}
