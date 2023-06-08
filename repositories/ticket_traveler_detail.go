package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TicketTravelerDetailRepository interface {
	GetAllTicketTravelerDetails() ([]models.TicketTravelerDetail, int, error)
	GetTicketTravelerDetailByID(id uint) (models.TicketTravelerDetail, error)
	GetTicketTravelerDetailByTicketOrderID(id uint) ([]models.TicketTravelerDetail, error)
	GetTicketTravelerDetailByTicketOrderIDAndTrainID(ticketOrderId, trainId uint) (models.TicketTravelerDetail, error)
	CreateTicketTravelerDetail(ticketTravelerDetail models.TicketTravelerDetail) (models.TicketTravelerDetail, error)
	UpdateTicketTravelerDetail(ticketTravelerDetail models.TicketTravelerDetail) (models.TicketTravelerDetail, error)
}

type ticketTravelerDetailRepository struct {
	db *gorm.DB
}

func NewTicketTravelerDetailRepository(db *gorm.DB) TicketTravelerDetailRepository {
	return &ticketTravelerDetailRepository{db}
}

func (r *ticketTravelerDetailRepository) GetAllTicketTravelerDetails() ([]models.TicketTravelerDetail, int, error) {
	var (
		ticketTravelerDetails []models.TicketTravelerDetail
		count                 int64
	)

	err := r.db.Find(&ticketTravelerDetails).Count(&count).Error
	if err != nil {
		return ticketTravelerDetails, int(count), err
	}

	return ticketTravelerDetails, int(count), err
}

func (r *ticketTravelerDetailRepository) GetTicketTravelerDetailByID(id uint) (models.TicketTravelerDetail, error) {
	var ticketTravelerDetail models.TicketTravelerDetail
	err := r.db.Where("id = ?", id).First(&ticketTravelerDetail).Error
	return ticketTravelerDetail, err
}

func (r *ticketTravelerDetailRepository) GetTicketTravelerDetailByTicketOrderID(id uint) ([]models.TicketTravelerDetail, error) {
	var ticketTravelerDetail []models.TicketTravelerDetail
	err := r.db.Where("ticket_order_id = ?", id).Find(&ticketTravelerDetail).Error
	return ticketTravelerDetail, err
}

func (r *ticketTravelerDetailRepository) GetTicketTravelerDetailByTicketOrderIDAndTrainID(ticketOrderId, trainId uint) (models.TicketTravelerDetail, error) {
	var ticketTravelerDetail models.TicketTravelerDetail
	err := r.db.Where("ticket_order_id = ? AND train_id = ?", ticketOrderId, trainId).First(&ticketTravelerDetail).Error
	return ticketTravelerDetail, err
}

func (r *ticketTravelerDetailRepository) CreateTicketTravelerDetail(ticketTravelerDetail models.TicketTravelerDetail) (models.TicketTravelerDetail, error) {
	err := r.db.Create(&ticketTravelerDetail).Error
	return ticketTravelerDetail, err
}

func (r *ticketTravelerDetailRepository) UpdateTicketTravelerDetail(ticketTravelerDetail models.TicketTravelerDetail) (models.TicketTravelerDetail, error) {
	err := r.db.Save(ticketTravelerDetail).Error
	return ticketTravelerDetail, err
}
