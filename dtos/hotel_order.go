package dtos

import "time"

type HotelOrderInput struct {
	// HotelID          int                   `form:"hotel_id" json:"hotel_id" example:"1"`
	HotelRoomID      int                   `form:"hotel_room_id" json:"hotel_room_id" example:"1"`
	QuantityAdult    int                   `form:"quantity_adult" json:"quantity_adult" example:"1"`
	QuantityInfant   int                   `form:"quantity_infant" json:"quantity_infant" example:"1"`
	DateStart        string                `form:"date_start" json:"date_start" example:"2023-05-01"`
	DateEnd          string                `form:"date_end" json:"date_end" example:"2023-05-02"`
	PaymentID        int                   `form:"payment_id" json:"payment_id" example:"1"`
	NameOrder        string                `form:"name_order" json:"name_order" example:"Mochammad Hanif"`
	EmailOrder       string                `form:"email_order" json:"email_order" example:"me@hanifz.com"`
	PhoneNumberOrder string                `form:"phone_number_order" json:"phone_number_order" example:"085115151515"`
	SpecialRequest   string                `form:"special_request" json:"special_request" example:"Tambah 1 Bed"`
	TravelerDetail   []TravelerDetailInput `json:"traveler_detail"`
}

type HotelOrderResponse struct {
	HotelOrderID     int                       `json:"hotel_order_id" example:"1"`
	QuantityAdult    int                       `json:"quantity_adult" example:"1"`
	QuantityInfant   int                       `json:"quantity_infant" example:"1"`
	NumberOfNight    int                       `json:"number_of_night" example:"1"`
	DateStart        string                    `json:"check_in_date" example:"2023-05-01"`
	DateEnd          string                    `json:"check_out_date" example:"2023-05-02"`
	Price            int                       `json:"price" example:"50000"`
	TotalAmount      int                       `json:"total_amount" example:"50000"`
	NameOrder        string                    `json:"name_order" example:"Mochammad Hanif"`
	EmailOrder       string                    `json:"email_order" example:"me@hanifz.com"`
	PhoneNumberOrder string                    `json:"phone_number_order" example:"085115151515"`
	SpecialRequest   string                    `json:"special_request" example:"Minta 1 Bed"`
	HotelOrderCode   string                    `json:"ticket_order_code" example:"RANDOMCODE123"`
	IsCheckIn        bool                      `json:"is_check_in" example:"false"`
	IsCheckOut       bool                      `json:"is_check_out" example:"false"`
	Status           string                    `json:"status" example:"unpaid"`
	Hotel            HotelByIDResponses        `json:"hotel"`
	Payment          PaymentResponses          `json:"payment,omitempty"`
	TravelerDetail   []TravelerDetailResponse  `json:"traveler_detail"`
	User             *UserInformationResponses `json:"user,omitempty"`
	CreatedAt        time.Time                 `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt        time.Time                 `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
