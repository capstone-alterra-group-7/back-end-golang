package repositories

import "gorm.io/gorm"

type NotificationRepository interface {
	CreateNotificationRepository(userID uint, templateMessageID uint, content string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) CreateNotificationRepository(userID uint, templateMessageID uint, content string) error {
	notification := map[string]interface{}{
		"user_id":             userID,
		"template_message_id": templateMessageID,
		"content":             content,
	}
	err := r.db.Create(&notification).Error
	return err
}
