package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TicketOrderRepository interface {
	GetTicketOrders(page, limit int, status string) ([]models.TicketOrder, int, error)
	GetTicketOrderByStatusAndID(id, userID uint, status string) (models.TicketOrder, error)
	GetTicketOrderByID(id, userID uint) (models.TicketOrder, error)
	CreateTicketOrder(ticketOrder models.TicketOrder) (models.TicketOrder, error)
	UpdateTicketOrder(ticketOrder models.TicketOrder) (models.TicketOrder, error)
}

type ticketOrderRepository struct {
	db *gorm.DB
}

func NewTicketOrderRepository(db *gorm.DB) TicketOrderRepository {
	return &ticketOrderRepository{db}
}

func (r *ticketOrderRepository) GetTicketOrders(page, limit int, status string) ([]models.TicketOrder, int, error) {
	var (
		ticketOrders []models.TicketOrder
		count        int64
	)
	err := r.db.Where("status = ?", status).Find(&ticketOrders).Count(&count).Error
	if err != nil {
		return ticketOrders, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&ticketOrders).Error

	return ticketOrders, int(count), err
}

func (r *ticketOrderRepository) GetTicketOrderByStatusAndID(id, userID uint, status string) (models.TicketOrder, error) {
	var ticketOrder models.TicketOrder
	if userID == 1 {
		err := r.db.Where("id = ? AND status = ?", id, status).First(&ticketOrder).Error
		return ticketOrder, err
	}
	err := r.db.Where("id = ? AND user_id = ? AND status = ?", id, userID, status).First(&ticketOrder).Error
	return ticketOrder, err
}

func (r *ticketOrderRepository) GetTicketOrderByID(id, userID uint) (models.TicketOrder, error) {
	var ticketOrder models.TicketOrder
	if userID == 1 {
		err := r.db.Where("id = ?", id).First(&ticketOrder).Error
		return ticketOrder, err
	}
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&ticketOrder).Error
	return ticketOrder, err
}

func (r *ticketOrderRepository) CreateTicketOrder(ticketOrder models.TicketOrder) (models.TicketOrder, error) {
	err := r.db.Create(&ticketOrder).Error
	return ticketOrder, err
}

func (r *ticketOrderRepository) UpdateTicketOrder(ticketOrder models.TicketOrder) (models.TicketOrder, error) {
	err := r.db.Save(ticketOrder).Error
	return ticketOrder, err
}
