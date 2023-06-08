package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TravelerDetailRepository interface {
	GetAllTravelerDetails(page, limit int) ([]models.TravelerDetail, int, error)
	GetTravelerDetailByID(id uint) (models.TravelerDetail, error)
	GetTravelerDetailByID2(ticketOrderID uint) ([]models.TravelerDetail, error)
	GetTravelerDetailByTicketOrderID(id uint) ([]models.TravelerDetail, error)
	CreateTravelerDetail(travelerDetail models.TravelerDetail) (models.TravelerDetail, error)
}

type travelerDetailRepository struct {
	db *gorm.DB
}

func NewTravelerDetailRepository(db *gorm.DB) TravelerDetailRepository {
	return &travelerDetailRepository{db}
}

func (r *travelerDetailRepository) GetAllTravelerDetails(page, limit int) ([]models.TravelerDetail, int, error) {
	var (
		travelerDetails []models.TravelerDetail
		count           int64
	)
	err := r.db.Find(&travelerDetails).Count(&count).Error
	if err != nil {
		return travelerDetails, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&travelerDetails).Error

	return travelerDetails, int(count), err
}

func (r *travelerDetailRepository) GetTravelerDetailByID(id uint) (models.TravelerDetail, error) {
	var travelerDetail models.TravelerDetail
	err := r.db.Where("id = ?", id).First(&travelerDetail).Error
	return travelerDetail, err
}

func (r *travelerDetailRepository) GetTravelerDetailByID2(ticketOrderID uint) ([]models.TravelerDetail, error) {
	var travelerDetail []models.TravelerDetail
	err := r.db.Where("ticket_order_id = ?", ticketOrderID).Find(&travelerDetail).Error
	return travelerDetail, err
}

func (r *travelerDetailRepository) GetTravelerDetailByTicketOrderID(id uint) ([]models.TravelerDetail, error) {
	var travelerDetail []models.TravelerDetail
	err := r.db.Where("ticket_order_id = ?", id).Find(&travelerDetail).Error
	return travelerDetail, err
}

func (r *travelerDetailRepository) CreateTravelerDetail(travelerDetail models.TravelerDetail) (models.TravelerDetail, error) {
	err := r.db.Create(&travelerDetail).Error
	return travelerDetail, err
}
