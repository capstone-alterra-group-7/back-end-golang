package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelFacilitiesRepository interface {
	GetAllHotelFacilities(page, limit int) ([]models.HotelFacilities, int, error)
	GetAllHotelFacilitiesByID(id uint) ([]models.HotelFacilities, error)
	GetHotelFacilitiesByID(id uint) (models.HotelFacilities, error)
	CreateHotelFacilities(HotelFacilities models.HotelFacilities) (models.HotelFacilities, error)
	UpdateHotelFacilities(HotelFacilities models.HotelFacilities) (models.HotelFacilities, error)
	DeleteHotelFacilities(id uint) error
}

type hotelFacilitiesRepository struct {
	db *gorm.DB
}

func NewHotelFacilitiesRepository(db *gorm.DB) HotelFacilitiesRepository {
	return &hotelFacilitiesRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelFacilitiesRepository) GetAllHotelFacilities(page, limit int) ([]models.HotelFacilities, int, error) {
	var (
		hotels []models.HotelFacilities
		count  int64
	)
	err := r.db.Find(&hotels).Count(&count).Error
	if err != nil {
		return hotels, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&hotels).Error

	return hotels, int(count), err
}

func (r *hotelFacilitiesRepository) GetAllHotelFacilitiesByID(id uint) ([]models.HotelFacilities, error) {
	var HotelFacilities []models.HotelFacilities
	err := r.db.Where("hotel_id = ?", id).Unscoped().Find(&HotelFacilities).Error
	return HotelFacilities, err
}

func (r *hotelFacilitiesRepository) GetHotelFacilitiesByID(id uint) (models.HotelFacilities, error) {
	var HotelFacilities models.HotelFacilities
	err := r.db.Where("id = ?", id).First(&HotelFacilities).Error
	return HotelFacilities, err
}

func (r *hotelFacilitiesRepository) CreateHotelFacilities(HotelFacilities models.HotelFacilities) (models.HotelFacilities, error) {
	err := r.db.Create(&HotelFacilities).Error
	return HotelFacilities, err
}

func (r *hotelFacilitiesRepository) UpdateHotelFacilities(HotelFacilities models.HotelFacilities) (models.HotelFacilities, error) {
	err := r.db.Save(&HotelFacilities).Error
	return HotelFacilities, err
}

func (r *hotelFacilitiesRepository) DeleteHotelFacilities(id uint) error {
	var HotelFacilities models.HotelFacilities
	err := r.db.Where("hotel_id = ?", id).Delete(&HotelFacilities).Error
	return err
}
