package dtos

type HotelRoomFacilitiesInput struct {
	HotelID     uint   `form:"hotel_id" json:"hotel_id"`
	HotelRoomID uint   `form:"hotel_room_id" json:"hotel_room_id"`
	Name        string `form:"name" json:"name"`
}

type HotelRoomFacilitiesResponse struct {
	HotelID     uint   `form:"hotel_id" json:"hotel_id"`
	HotelRoomID uint   `form:"hotel_room_id" json:"hotel_room_id"`
	Name        string `form:"name" json:"name"`
}
