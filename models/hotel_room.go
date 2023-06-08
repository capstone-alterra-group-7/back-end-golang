package models

import "gorm.io/gorm"

type HotelRoom struct {
	gorm.Model
	HotelID          uint   `form:"hotel_id" json:"hotel_id"`
	Hotel            Hotel  `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name             string `form:"name" json:"name"`
	SizeOfRoom       int    `form:"size_of_room" json:"size_of_room"`
	QuantityOfRoom   int    `form:"quantity_of_room" json:"quantity_of_room"`
	Description      string `form:"description" json:"description"`
	NormalPrice      int    `form:"normal_price" json:"normal_price"`
	Discount         int    `form:"discount" json:"discount"`
	DiscountPrice    int    `form:"discount_price" json:"discount_price"`
	NumberOfGuest    int    `form:"number_of_guest" json:"number_of_guest"`
	MattressSize     string `form:"mattress_size" json:"mattress_size"`
	NumberOfMattress int    `form:"number_of_mattress" json:"number_of_mattress"`
}
