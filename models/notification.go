package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID            uint `json:"user_id" from:"user_id"`
	TemplateMessageID uint `json:"template_message_id" from:"template_message_id" gorm:"foreignKey:TemplateMessageID"`
}
