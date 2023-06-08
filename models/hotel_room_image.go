package models

import "gorm.io/gorm"

type HotelRoomImage struct {
	gorm.Model
	HotelID     uint      `form:"hotel_id" json:"hotel_id"`
	Hotel       Hotel     `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	HotelRoomID uint      `form:"hotel_room_id" json:"hotel_room_id"`
	HotelRoom   HotelRoom `gorm:"foreignKey:HotelRoomID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ImageUrl    string    `form:"image_url" json:"image_url"`
}
