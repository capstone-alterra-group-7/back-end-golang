package models

import "gorm.io/gorm"

type HotelPolicies struct {
	gorm.Model
	HotelID            uint   `form:"hotel_id" json:"hotel_id"`
	Hotel              Hotel  `gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	IsCheckInCheckOut  bool   `form:"is_check_in_check_out" json:"is_check_in_check_out"`
	TimeCheckIn        string `form:"time_check_in" json:"time_check_in"`
	TimeCheckOut       string `form:"time_check_out" json:"time_check_out"`
	IsPolicyCanceled   bool   `form:"is_policy_canceled" json:"is_policy_canceled"`
	PolicyMinimumAge   int    `form:"policy_minimum_age" json:"policy_minimum_age"`
	IsCheckInEarly     bool   `form:"is_check_in_early" json:"is_check_in_early"`
	IsCheckOutOverdue  bool   `form:"is_check_out_overdue" json:"is_check_out_overdue"`
	IsBreakfast        bool   `form:"is_breakfast" json:"is_breakfast"`
	TimeBreakfastStart string `form:"time_breakfast_start" json:"time_breakfast_start"`
	TimeBreakfastEnd   string `form:"time_breakfast_end" json:"time_breakfast_end"`
	IsSmoking          bool   `form:"is_smoking" json:"is_smoking"`
	IsPet              bool   `form:"is_pet" json:"is_pet"`
}
