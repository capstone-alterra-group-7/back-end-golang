package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TemplateMessageRepository interface {
	GetTemplateMessageByID(id uint) (*models.TemplateMessage, error)
}

type templateMessageRepository struct {
	db *gorm.DB
}

func NewTemplateMessageRepository(db *gorm.DB) TemplateMessageRepository {
	return &templateMessageRepository{db}
}

func (r *templateMessageRepository) GetTemplateMessageByID(id uint) (*models.TemplateMessage, error) {
	var templateMessage models.TemplateMessage
	err := r.db.First(&templateMessage, id).Error
	if err != nil {
		return nil, err
	}
	return &templateMessage, nil
}
