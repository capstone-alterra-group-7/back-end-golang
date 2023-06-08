package dtos

import "time"

type TicketOrderInput struct {
	QuantityAdult                 int                         `form:"quantity_adult" json:"quantity_adult" example:"1"`
	QuantityInfant                int                         `form:"quantity_infant" json:"quantity_infant" example:"1"`
	WithReturn                    bool                        `form:"with_return" json:"with_return" example:"true"`
	PaymentID                     int                         `form:"payment_id" json:"payment_id" example:"1"`
	NameOrder                     string                      `form:"name_order" json:"name_order" example:"Mochammad Hanif"`
	EmailOrder                    string                      `form:"email_order" json:"email_order" example:"me@hanifz.com"`
	PhoneNumberOrder              string                      `form:"phone_number_order" json:"phone_number_order" example:"085115151515"`
	TravelerDetail                []TravelerDetailInput       `json:"traveler_detail"`
	TicketTravelerDetailDeparture []TicketTravelerDetailInput `json:"ticket_traveler_detail_departure"`
	TicketTravelerDetailReturn    []TicketTravelerDetailInput `json:"ticket_traveler_detail_return"`
}

type TravelerDetailInput struct {
	Title        string `form:"title" json:"title" example:"Saudara"`
	FullName     string `form:"full_name" json:"full_name" example:"Mochammad Hanif"`
	IDCardNumber string `form:"id_card_number" json:"id_card_number" example:"1902389012801211"`
}

type TravelerDetailResponse struct {
	ID           int    `form:"traveler_detail_id" json:"traveler_detail_id" example:"1"`
	Title        string `form:"title" json:"title" example:"Saudara"`
	FullName     string `form:"full_name" json:"full_name" example:"Mochammad Hanif"`
	IDCardNumber string `form:"id_card_number" json:"id_card_number" example:"1902389012801211"`
}

type TicketTravelerDetailInput struct {
	TrainID              int    `form:"train_id" json:"train_id" example:"1"`
	TrainCarriageID      int    `form:"train_carriage_id" json:"train_carriage_id" example:"1"`
	TrainSeatID          int    `form:"train_seat_id" json:"train_seat_id" example:"2"`
	StationOriginID      int    `form:"station_origin_id" json:"station_origin_id" example:"1"`
	StationDestinationID int    `form:"station_destination_id" json:"station_destination_id" example:"2"`
	Date                 string `form:"date" json:"date" example:"2023-05-31"`
}

type TicketOrderResponse struct {
	TicketOrderID        int                            `json:"ticket_order_id" example:"1"`
	QuantityAdult        int                            `json:"quantity_adult" example:"1"`
	QuantityInfant       int                            `json:"quantity_infant" example:"1"`
	Price                int                            `json:"price" example:"50000"`
	TotalAmount          int                            `json:"total_amount" example:"50000"`
	NameOrder            string                         `json:"name_order" example:"Mochammad Hanif"`
	EmailOrder           string                         `json:"email_order" example:"me@hanifz.com"`
	PhoneNumberOrder     string                         `json:"phone_number_order" example:"085115151515"`
	TicketOrderCode      string                         `json:"ticket_order_code" example:"RANDOMCODE123"`
	Status               string                         `json:"status" example:"unpaid"`
	Payment              PaymentResponses               `json:"payment,omitempty"`
	TicketTravelerDetail []TicketTravelerDetailResponse `json:"ticket_traveler_detail"`
	CreatedAt            time.Time                      `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt            time.Time                      `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}

type TicketTravelerDetailResponse struct {
	TicketTravelerDetailID int                    `json:"ticket_traveler_detail_id" example:"1"`
	TravelerDetail         TravelerDetailResponse `json:"traveler_detail"`
	Train                  TrainResponsesSimply   `json:"train"`
	StationOrigin          StationResponseSimply  `json:"station_origin"`
	StationDestination     StationResponseSimply  `json:"station_destination"`
	Date                   time.Time              `json:"date" example:"2023-05-31"`
	BoardingTicketCode     string                 `json:"boarding_ticket_code" example:"RANDOMBOARDINGTICKETCODE123"`
}

type TicketTravelerDetailOrderResponse struct {
	TicketOrderID      int                       `json:"ticket_order_id" example:"1"`
	QuantityAdult      int                       `json:"quantity_adult" example:"1"`
	QuantityInfant     int                       `json:"quantity_infant" example:"1"`
	NameOrder          string                    `json:"name_order" example:"Mochammad Hanif"`
	EmailOrder         string                    `json:"email_order" example:"me@hanifz.com"`
	PhoneNumberOrder   string                    `json:"phone_number_order" example:"085115151515"`
	TicketOrderCode    string                    `json:"ticket_order_code" example:"RANDOMCODE123"`
	User               *UserInformationResponses `json:"user,omitempty"`
	Payment            PaymentResponses          `json:"payment"`
	TravelerDetail     []TravelerDetailResponse  `json:"traveler_detail"`
	Train              TrainResponsesSimply      `json:"train"`
	StationOrigin      StationResponseSimply     `json:"station_origin"`
	StationDestination StationResponseSimply     `json:"station_destination"`
	Date               time.Time                 `json:"date" example:"2023-05-31"`
	BoardingTicketCode string                    `json:"boarding_ticket_code" example:"RANDOMBOARDINGTICKETCODE123"`
	Status             string                    `json:"status" example:"2023-05-31"`
	CreatedAt          time.Time                 `json:"created_at" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt          time.Time                 `json:"updated_at" example:"2023-05-17T15:07:16.504+07:00"`
}
