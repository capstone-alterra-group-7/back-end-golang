package models

import "gorm.io/gorm"

type HistorySearchHotel struct {
	gorm.Model
	UserID  uint `json:"user_id" form:"user_id"`
	HotelID uint `json:"hotel_id" form:"hotel_id"`
}
