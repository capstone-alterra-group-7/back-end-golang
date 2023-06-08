package dtos

type HotelRoomImageInput struct {
	HotelID     uint   `form:"hotel_id" json:"hotel_id"`
	HotelRoomID uint   `form:"hotel_room_id" json:"hotel_room_id"`
	ImageUrl    string `form:"image_url" json:"image_url"`
}

type HotelRoomImageResponse struct {
	HotelID     uint   `form:"hotel_id" json:"hotel_id"`
	HotelRoomID uint   `form:"hotel_room_id" json:"hotel_room_id"`
	ImageUrl    string `form:"image_url" json:"image_url"`
}
