package dtos

type HotelFacilitiesInput struct {
	Name string `form:"name" json:"name"`
}

type HotelFacilitiesResponse struct {
	HotelID uint   `form:"hotel_id" json:"hotel_id"`
	Name    string `form:"name" json:"name"`
}
