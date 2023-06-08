package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelRoomImageRepository interface {
	GetAllHotelRoomImages(page, limit int) ([]models.HotelRoomImage, int, error)
	GetAllHotelRoomImageByID(id uint) ([]models.HotelRoomImage, error)
	GetHotelRoomImageByID(id uint) (models.HotelRoomImage, error)
	CreateHotelRoomImage(roomImage models.HotelRoomImage) (models.HotelRoomImage, error)
	UpdateHotelRoomImage(roomImage models.HotelRoomImage) (models.HotelRoomImage, error)
	DeleteHotelRoomImage(id uint) error
}

type hotelRoomImageRepository struct {
	db *gorm.DB
}

func NewHotelRoomImageRepository(db *gorm.DB) HotelRoomImageRepository {
	return &hotelRoomImageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelRoomImageRepository) GetAllHotelRoomImages(page, limit int) ([]models.HotelRoomImage, int, error) {
	var (
		hotels []models.HotelRoomImage
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

func (r *hotelRoomImageRepository) GetAllHotelRoomImageByID(id uint) ([]models.HotelRoomImage, error) {
	var roomImage []models.HotelRoomImage
	err := r.db.Where("hotel_room_id = ?", id).Find(&roomImage).Error
	return roomImage, err
}

func (r *hotelRoomImageRepository) GetHotelRoomImageByID(id uint) (models.HotelRoomImage, error) {
	var roomImage models.HotelRoomImage
	err := r.db.Where("id = ?", id).First(&roomImage).Error
	return roomImage, err
}

func (r *hotelRoomImageRepository) CreateHotelRoomImage(roomImage models.HotelRoomImage) (models.HotelRoomImage, error) {
	err := r.db.Create(&roomImage).Error
	return roomImage, err
}

func (r *hotelRoomImageRepository) UpdateHotelRoomImage(roomImage models.HotelRoomImage) (models.HotelRoomImage, error) {
	err := r.db.Save(&roomImage).Error
	return roomImage, err
}

func (r *hotelRoomImageRepository) DeleteHotelRoomImage(id uint) error {
	var roomImage models.HotelRoomImage
	err := r.db.Where("hotel_room_id = ?", id).Delete(&roomImage).Error
	return err
}
