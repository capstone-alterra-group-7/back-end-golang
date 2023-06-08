package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	GetAllPayments(page, limit int) ([]models.Payment, int, error)
	GetPaymentByID(id uint) (models.Payment, error)
	GetPaymentByID2(id uint) (models.Payment, error)
	CreatePayment(payment models.Payment) (models.Payment, error)
	UpdatePayment(payment models.Payment) (models.Payment, error)
	DeletePayment(id uint) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *paymentRepository) GetAllPayments(page, limit int) ([]models.Payment, int, error) {
	var (
		payments []models.Payment
		count    int64
	)
	err := r.db.Find(&payments).Count(&count).Error
	if err != nil {
		return payments, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&payments).Error

	return payments, int(count), err
}

func (r *paymentRepository) GetPaymentByID(id uint) (models.Payment, error) {
	var payment models.Payment
	err := r.db.Unscoped().Where("id = ?", id).First(&payment).Error
	return payment, err
}

func (r *paymentRepository) GetPaymentByID2(id uint) (models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("id = ?", id).First(&payment).Error
	return payment, err
}

func (r *paymentRepository) CreatePayment(payment models.Payment) (models.Payment, error) {
	err := r.db.Create(&payment).Error
	return payment, err
}

func (r *paymentRepository) UpdatePayment(payment models.Payment) (models.Payment, error) {
	err := r.db.Save(&payment).Error
	return payment, err
}

func (r *paymentRepository) DeletePayment(id uint) error {
	var payment models.Payment
	err := r.db.Where("id = ?", id).Delete(&payment).Error
	return err
}
