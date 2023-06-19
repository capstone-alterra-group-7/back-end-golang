package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelRoomRepository interface {
	GetAllHotelRooms(page, limit int) ([]models.HotelRoom, int, error)
	GetAllHotelRoomByHotelID(id uint) ([]models.HotelRoom, error)
	GetHotelRoomByID(id uint) (models.HotelRoom, error)
	GetMinimumPriceHotelRoomByHotelID(id uint) (models.HotelRoom, error)
	CreateHotelRoom(hotelRoom models.HotelRoom) (models.HotelRoom, error)
	UpdateHotelRoom(hotelRoom models.HotelRoom) (models.HotelRoom, error)
	DeleteHotelRoom(id uint) error
}

type hotelRoomRepository struct {
	db *gorm.DB
}

func NewHotelRoomRepository(db *gorm.DB) HotelRoomRepository {
	return &hotelRoomRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelRoomRepository) GetAllHotelRooms(page, limit int) ([]models.HotelRoom, int, error) {
	var (
		hotelrooms []models.HotelRoom
		count      int64
	)
	err := r.db.Find(&hotelrooms).Count(&count).Error
	if err != nil {
		return hotelrooms, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&hotelrooms).Error

	return hotelrooms, int(count), err
}

func (r *hotelRoomRepository) GetAllHotelRoomByHotelID(id uint) ([]models.HotelRoom, error) {
	var hotelRoom []models.HotelRoom
	err := r.db.Where("hotel_id = ?", id).Unscoped().Find(&hotelRoom).Error
	return hotelRoom, err
}

func (r *hotelRoomRepository) GetMinimumPriceHotelRoomByHotelID(id uint) (models.HotelRoom, error) {
	var hotelRoom models.HotelRoom
	err := r.db.Where("hotel_id = ?", id).Order("discount_price ASC").First(&hotelRoom).Error
	return hotelRoom, err
}

func (r *hotelRoomRepository) GetHotelRoomByID(id uint) (models.HotelRoom, error) {
	var hotelRoom models.HotelRoom
	err := r.db.Where("id = ?", id).First(&hotelRoom).Error
	return hotelRoom, err
}

func (r *hotelRoomRepository) CreateHotelRoom(hotelRoom models.HotelRoom) (models.HotelRoom, error) {
	err := r.db.Create(&hotelRoom).Error
	return hotelRoom, err
}

func (r *hotelRoomRepository) UpdateHotelRoom(hotelRoom models.HotelRoom) (models.HotelRoom, error) {
	err := r.db.Save(&hotelRoom).Error
	return hotelRoom, err
}

func (r *hotelRoomRepository) DeleteHotelRoom(id uint) error {
	var hotelRoom models.HotelRoom
	err := r.db.Where("id = ?", id).Delete(&hotelRoom).Error
	return err
}
