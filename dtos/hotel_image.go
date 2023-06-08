package dtos

type HotelImageInput struct {
	HotelID  uint   `form:"hotel_id" json:"hotel_id"`
	ImageUrl string `form:"image_url" json:"image_url"`
}

type HotelImageResponse struct {
	HotelID  uint   `form:"hotel_id" json:"hotel_id"`
	ImageUrl string `form:"image_url" json:"image_url"`
}
