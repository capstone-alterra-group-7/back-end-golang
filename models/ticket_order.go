package models

import "gorm.io/gorm"

type TicketOrder struct {
	gorm.Model
	UserID           uint `form:"user_id" json:"user_id"`
	User             User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	QuantityAdult    int
	QuantityInfant   int
	Price            int
	WithReturn       bool
	PaymentID        int
	TotalAmount      int
	NameOrder        string
	EmailOrder       string
	PhoneNumberOrder string
	TicketOrderCode  string
	Status           string `gorm:"type:ENUM('unpaid', 'paid', 'done', 'canceled', 'refund')"`
}
