package dtos

import "time"

type HotelInput struct {
	Name            string                 `form:"name" json:"name"`
	Class           string                 `form:"class" json:"class"`
	Description     string                 `form:"description" json:"description"`
	PhoneNumber     string                 `form:"phone_number" json:"phone_number"`
	Email           string                 `form:"email" json:"email"`
	Address         string                 `form:"address" json:"address"`
	HotelImage      []HotelImageInput      `form:"hotel_image" json:"hotel_image"`
	HotelFacilities []HotelFacilitiesInput `form:"hotel_facilities" json:"hotel_facilities"`
	HotelPolicy     []HotelPoliciesInput   `form:"hotel_policy" json:"hotel_policy"`
}

type HotelResponse struct {
	HotelID         uint                      `form:"hotel_id" json:"hotel_id"`
	Name            string                    `form:"name" json:"name"`
	Class           string                    `form:"class" json:"class"`
	Description     string                    `form:"description" json:"description"`
	PhoneNumber     string                    `form:"phone_number" json:"phone_number"`
	Email           string                    `form:"email" json:"email"`
	Address         string                    `form:"address" json:"address"`
	HotelImage      []HotelImageResponse      `form:"hotel_image" json:"hotel_image"`
	HotelFacilities []HotelFacilitiesResponse `form:"hotel_facilities" json:"hotel_facilities"`
	HotelPolicy     HotelPoliciesResponse     `form:"hotel_policy" json:"hotel_policy"`
	CreatedAt       time.Time                 `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt       time.Time                 `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
