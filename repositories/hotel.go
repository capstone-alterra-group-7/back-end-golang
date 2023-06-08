package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelRepository interface {
	GetAllHotels(page, limit int) ([]models.Hotel, int, error)
	GetHotelByID(id uint) (models.Hotel, error)
	CreateHotel(hotel models.Hotel) (models.Hotel, error)
	UpdateHotel(hotel models.Hotel) (models.Hotel, error)
	DeleteHotel(id uint) error
}

type hotelRepository struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelRepository) GetAllHotels(page, limit int) ([]models.Hotel, int, error) {
	var (
		hotels []models.Hotel
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

func (r *hotelRepository) GetHotelByID(id uint) (models.Hotel, error) {
	var hotel models.Hotel
	err := r.db.Where("id = ?", id).First(&hotel).Error
	return hotel, err
}

func (r *hotelRepository) CreateHotel(hotel models.Hotel) (models.Hotel, error) {
	err := r.db.Create(&hotel).Error
	return hotel, err
}

func (r *hotelRepository) UpdateHotel(hotel models.Hotel) (models.Hotel, error) {
	err := r.db.Save(&hotel).Error
	return hotel, err
}

func (r *hotelRepository) DeleteHotel(id uint) error {
	var hotel models.Hotel
	err := r.db.Where("id = ?", id).Delete(&hotel).Error
	return err
}
