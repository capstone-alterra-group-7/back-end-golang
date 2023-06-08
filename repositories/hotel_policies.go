package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelPoliciesRepository interface {
	GetAllHotelPolicies(page, limit int) ([]models.HotelPolicies, int, error)
	GetHotelPoliciesByIDHotel(id uint) (models.HotelPolicies, error)
	CreateHotelPolicies(HotelPolicies models.HotelPolicies) (models.HotelPolicies, error)
	UpdateHotelPolicies(HotelPolicies models.HotelPolicies) (models.HotelPolicies, error)
	DeleteHotelPolicies(id uint) error
}

type hotelPoliciesRepository struct {
	db *gorm.DB
}

func NewHotelPoliciesRepository(db *gorm.DB) HotelPoliciesRepository {
	return &hotelPoliciesRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelPoliciesRepository) GetAllHotelPolicies(page, limit int) ([]models.HotelPolicies, int, error) {
	var (
		hotelPolicies []models.HotelPolicies
		count         int64
	)
	err := r.db.Find(&hotelPolicies).Count(&count).Error
	if err != nil {
		return hotelPolicies, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&hotelPolicies).Error

	return hotelPolicies, int(count), err
}

func (r *hotelPoliciesRepository) GetHotelPoliciesByIDHotel(id uint) (models.HotelPolicies, error) {
	var HotelPolicies models.HotelPolicies
	err := r.db.Where("hotel_id = ?", id).First(&HotelPolicies).Error
	return HotelPolicies, err
}

func (r *hotelPoliciesRepository) CreateHotelPolicies(HotelPolicies models.HotelPolicies) (models.HotelPolicies, error) {
	err := r.db.Create(&HotelPolicies).Error
	return HotelPolicies, err
}

func (r *hotelPoliciesRepository) UpdateHotelPolicies(HotelPolicies models.HotelPolicies) (models.HotelPolicies, error) {
	err := r.db.Save(&HotelPolicies).Error
	return HotelPolicies, err
}

func (r *hotelPoliciesRepository) DeleteHotelPolicies(id uint) error {
	var HotelPolicies models.HotelPolicies
	err := r.db.Where("hotel_id = ?", id).Delete(&HotelPolicies).Error
	return err
}
