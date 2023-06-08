package dtos

import "time"

type PaymentInput struct {
	Type          string `form:"type" json:"type" example:"Virtual Account"`
	ImageUrl      string `form:"image_url" json:"image_url" example:"bca.png"`
	Name          string `form:"name" json:"name" example:"Bank Central Asia"`
	AccountName   string `form:"account_name" json:"account_name" example:"PT Tripease"`
	AccountNumber string `form:"account_number" json:"account_number" example:"00832982893221132"`
}

type PaymentResponse struct {
	ID            int        `json:"id" example:"1"`
	Type          string     `json:"type" example:"Virtual Account"`
	ImageUrl      string     `json:"image_url" example:"bca.png"`
	Name          string     `json:"name" example:"Bank Central Asia"`
	AccountName   string     `json:"account_name" example:"PT Tripease"`
	AccountNumber string     `json:"account_number" example:"00832982893221132"`
	CreatedAt     *time.Time `json:"created_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}

type PaymentResponses struct {
	ID            int        `json:"payment_id" example:"1"`
	Type          string     `json:"type" example:"Virtual Account"`
	ImageUrl      string     `json:"image_url" example:"bca.png"`
	Name          string     `json:"name" example:"Bank Central Asia"`
	AccountName   string     `json:"account_name" example:"PT Tripease"`
	AccountNumber string     `json:"account_number" example:"00832982893221132"`
	CreatedAt     *time.Time `json:"created_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" example:"2023-05-17T15:07:16.504+07:00"`
}
