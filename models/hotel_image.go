package models

import "gorm.io/gorm"

type HotelImage struct {
	gorm.Model
	HotelID  uint   `form:"hotel_id" json:"hotel_id"`
	Hotel    Hotel  `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageUrl string `form:"image_url" json:"image_url"`
}
