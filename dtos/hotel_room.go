package dtos

import "time"

type HotelRoomInput struct {
	HotelID           uint                          `form:"hotel_id" json:"hotel_id"`
	Name              string                        `form:"name" json:"name"`
	SizeOfRoom        int                           `form:"size_of_room" json:"size_of_room"`
	QuantityOfRoom    int                           `form:"quantity_of_room" json:"quantity_of_room"`
	Description       string                        `form:"description" json:"description"`
	NormalPrice       int                           `form:"normal_price" json:"normal_price"`
	Discount          int                           `form:"discount" json:"discount"`
	DiscountPrice     int                           `form:"discount_price" json:"discount_price"`
	NumberOfGuest     int                           `form:"number_of_guest" json:"number_of_guest"`
	MattressSize      string                        `form:"mattress_size" json:"mattress_size"`
	NumberOfMattress  int                           `form:"number_of_mattress" json:"number_of_mattress"`
	HotelRoomImage    []HotelRoomImageResponse      `form:"hotel_room_image" json:"hotel_room_image"`
	HotelRoomFacility []HotelRoomFacilitiesResponse `form:"hotel_room_facility" json:"hotel_room_facility"`
}

type HotelRoomResponse struct {
	HotelRoomID       uint                          `form:"hotel_room_id" json:"hotel_room_id"`
	HotelID           uint                          `form:"hotel_id" json:"hotel_id"`
	Name              string                        `form:"name" json:"name"`
	SizeOfRoom        int                           `form:"size_of_room" json:"size_of_room"`
	QuantityOfRoom    int                           `form:"quantity_of_room" json:"quantity_of_room"`
	Description       string                        `form:"description" json:"description"`
	NormalPrice       int                           `form:"normal_price" json:"normal_price"`
	Discount          int                           `form:"discount" json:"discount"`
	DiscountPrice     int                           `form:"discount_price" json:"discount_price"`
	NumberOfGuest     int                           `form:"number_of_guest" json:"number_of_guest"`
	MattressSize      string                        `form:"mattress_size" json:"mattress_size"`
	NumberOfMattress  int                           `form:"number_of_mattress" json:"number_of_mattress"`
	HotelRoomImage    []HotelRoomImageResponse      `form:"hotel_room_image" json:"hotel_room_image"`
	HotelRoomFacility []HotelRoomFacilitiesResponse `form:"hotel_room_facility" json:"hotel_room_facility"`
	CreatedAt         time.Time                     `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt         time.Time                     `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
