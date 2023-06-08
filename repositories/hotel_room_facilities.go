package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HotelRoomFacilitiesRepository interface {
	GetAllHotelRoomFacilities(page, limit int) ([]models.HotelRoomFacilities, int, error)
	GetAllHotelRoomFacilitiesByID(id uint) ([]models.HotelRoomFacilities, error)
	GetHotelRoomFacilitiesByID(id uint) (models.HotelRoomFacilities, error)
	CreateHotelRoomFacilities(roomFacilities models.HotelRoomFacilities) (models.HotelRoomFacilities, error)
	UpdateHotelRoomFacilities(roomFacilities models.HotelRoomFacilities) (models.HotelRoomFacilities, error)
	DeleteHotelRoomFacilities(id uint) error
}

type hotelRoomFacilitiesRepository struct {
	db *gorm.DB
}

func NewHotelRoomFacilitiesRepository(db *gorm.DB) HotelRoomFacilitiesRepository {
	return &hotelRoomFacilitiesRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *hotelRoomFacilitiesRepository) GetAllHotelRoomFacilities(page, limit int) ([]models.HotelRoomFacilities, int, error) {
	var (
		rooms []models.HotelRoomFacilities
		count int64
	)
	err := r.db.Find(&rooms).Count(&count).Error
	if err != nil {
		return rooms, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&rooms).Error

	return rooms, int(count), err
}

func (r *hotelRoomFacilitiesRepository) GetAllHotelRoomFacilitiesByID(id uint) ([]models.HotelRoomFacilities, error) {
	var roomFacilities []models.HotelRoomFacilities
	err := r.db.Where("hotel_room_id = ?", id).Find(&roomFacilities).Error
	return roomFacilities, err
}

func (r *hotelRoomFacilitiesRepository) GetHotelRoomFacilitiesByID(id uint) (models.HotelRoomFacilities, error) {
	var roomFacilities models.HotelRoomFacilities
	err := r.db.Where("id = ?", id).First(&roomFacilities).Error
	return roomFacilities, err
}

func (r *hotelRoomFacilitiesRepository) CreateHotelRoomFacilities(roomFacilities models.HotelRoomFacilities) (models.HotelRoomFacilities, error) {
	err := r.db.Create(&roomFacilities).Error
	return roomFacilities, err
}

func (r *hotelRoomFacilitiesRepository) UpdateHotelRoomFacilities(roomFacilities models.HotelRoomFacilities) (models.HotelRoomFacilities, error) {
	err := r.db.Save(&roomFacilities).Error
	return roomFacilities, err
}

func (r *hotelRoomFacilitiesRepository) DeleteHotelRoomFacilities(id uint) error {
	var roomFacilities models.HotelRoomFacilities
	err := r.db.Where("hotel_room_id = ?", id).Delete(&roomFacilities).Error
	return err
}
