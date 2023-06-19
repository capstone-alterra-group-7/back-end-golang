package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type HistorySearchHotelRepository interface {
	GetHistorySearchHotelByID(userID, id uint) (models.HistorySearchHotel, error)
	GetHistorySearchHotelsByUserID(userID uint, page, limit int) ([]models.HistorySearchHotel, int, error)
	CreateHistorySearchHotel(historySearchHotel models.HistorySearchHotel) (models.HistorySearchHotel, error)
	DeleteHistorySearchHotel(userID, id uint) error
}

type historySearchHotelRepository struct {
	db *gorm.DB
}

func NewHistorySearchHotelRepository(db *gorm.DB) HistorySearchHotelRepository {
	return &historySearchHotelRepository{db}
}

func (r *historySearchHotelRepository) GetHistorySearchHotelByID(userID, id uint) (models.HistorySearchHotel, error) {
	var historySearchHotel models.HistorySearchHotel
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&historySearchHotel).Error
	return historySearchHotel, err
}

func (r *historySearchHotelRepository) GetHistorySearchHotelsByUserID(userID uint, page, limit int) ([]models.HistorySearchHotel, int, error) {
	var historySearchHotels []models.HistorySearchHotel
	var count int64
	err := r.db.Where("user_id = ?", userID).Find(&historySearchHotels).Count(&count).Error
	if err != nil {
		return historySearchHotels, 0, err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&historySearchHotels).Error

	return historySearchHotels, int(count), err
}

func (r *historySearchHotelRepository) CreateHistorySearchHotel(historySearchHotel models.HistorySearchHotel) (models.HistorySearchHotel, error) {
	err := r.db.Create(&historySearchHotel).Error
	return historySearchHotel, err
}

func (r *historySearchHotelRepository) DeleteHistorySearchHotel(userID, id uint) error {
	var historySearchHotel models.HistorySearchHotel
	err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&historySearchHotel).Error
	return err
}
