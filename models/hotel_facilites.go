package models

import "gorm.io/gorm"

type HotelFacilities struct {
	gorm.Model
	HotelID uint   `form:"hotel_id" json:"hotel_id"`
	Hotel   Hotel  `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name    string `form:"name" json:"name"`
}
