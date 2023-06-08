package models

import "gorm.io/gorm"

type TrainCarriage struct {
	gorm.Model
	TrainID uint
	Train   Train  `gorm:"foreignKey:TrainID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Class   string `gorm:"type:varchar(255)"`
	Name    string
	Price   int
}
