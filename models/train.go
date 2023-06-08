package models

import "gorm.io/gorm"

type Train struct {
	gorm.Model
	CodeTrain string
	Name      string
	Status    string `gorm:"type:ENUM('available','unavailable'); default:'unavailable'"`
}
