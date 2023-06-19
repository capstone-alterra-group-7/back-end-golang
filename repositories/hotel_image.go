package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelImageRepository interface {
	GetAllHotelImages(page, limit int) ([]models.HotelImage, int, error)
	GetAllHotelImageByID(id uint) ([]models.HotelImage, error)
	GetHotelImageByID(id uint) (models.HotelImage, error)
	CreateHotelImage(hotelImage models.HotelImage) (models.HotelImage, error)
	UpdateHotelImage(hotelImage models.HotelImage) (models.HotelImage, error)
	DeleteHotelImage(id uint) error
}

type hotelImageRepository struct {
	db *gorm.DB
}

func NewHotelImageRepository(db *gorm.DB) HotelImageRepository {
	return &hotelImageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelImageRepository) GetAllHotelImages(page, limit int) ([]models.HotelImage, int, error) {
	var (
		hotels []models.HotelImage
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

func (r *hotelImageRepository) GetAllHotelImageByID(id uint) ([]models.HotelImage, error) {
	var hotelImage []models.HotelImage
	err := r.db.Where("hotel_id = ?", id).Unscoped().Find(&hotelImage).Error
	return hotelImage, err
}

func (r *hotelImageRepository) GetHotelImageByID(id uint) (models.HotelImage, error) {
	var hotelImage models.HotelImage
	err := r.db.Where("id = ?", id).First(&hotelImage).Error
	return hotelImage, err
}

func (r *hotelImageRepository) CreateHotelImage(hotelImage models.HotelImage) (models.HotelImage, error) {
	err := r.db.Create(&hotelImage).Error
	return hotelImage, err
}

func (r *hotelImageRepository) UpdateHotelImage(hotelImage models.HotelImage) (models.HotelImage, error) {
	err := r.db.Save(&hotelImage).Error
	return hotelImage, err
}

func (r *hotelImageRepository) DeleteHotelImage(id uint) error {
	var hotelImage models.HotelImage
	err := r.db.Where("hotel_id = ?", id).Delete(&hotelImage).Error
	return err
}
