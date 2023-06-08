package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TrainStationRepository interface {
	GetAllTrainStations(page, limit int) ([]models.TrainStation, int, error)
	GetTrainStationByID(id uint) (models.TrainStation, error)
	GetTrainStationByTrainIDStationID(trainID, stationID uint) (models.TrainStation, error)
	CreateTrainStation(trainStation models.TrainStation) (models.TrainStation, error)
	UpdateTrainStation(trainStation models.TrainStation) (models.TrainStation, error)
	DeleteTrainStation(trainStation models.TrainStation) error
	DeleteTrainStationById(id uint) error
}

type trainStationRepository struct {
	db *gorm.DB
}

func NewTrainStationRepository(db *gorm.DB) TrainStationRepository {
	return &trainStationRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *trainStationRepository) GetAllTrainStations(page, limit int) ([]models.TrainStation, int, error) {
	var (
		trainStations []models.TrainStation
		count         int64
	)
	err := r.db.Find(&trainStations).Count(&count).Error
	if err != nil {
		return trainStations, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&trainStations).Error

	return trainStations, int(count), err
}

func (r *trainStationRepository) GetTrainStationByID(id uint) (models.TrainStation, error) {
	var trainStation models.TrainStation
	err := r.db.Where("train_id = ?", id).First(&trainStation).Error
	return trainStation, err
}

func (r *trainStationRepository) GetTrainStationByTrainIDStationID(trainID, stationID uint) (models.TrainStation, error) {
	var trainStation models.TrainStation
	err := r.db.Unscoped().Where("train_id = ? AND station_id = ?", trainID, stationID).First(&trainStation).Error
	return trainStation, err
}

func (r *trainStationRepository) CreateTrainStation(trainStation models.TrainStation) (models.TrainStation, error) {
	err := r.db.Create(&trainStation).Error
	return trainStation, err
}

func (r *trainStationRepository) UpdateTrainStation(trainStation models.TrainStation) (models.TrainStation, error) {
	err := r.db.Save(&trainStation).Error
	return trainStation, err
}

func (r *trainStationRepository) DeleteTrainStation(trainStation models.TrainStation) error {
	err := r.db.Delete(&trainStation).Error
	return err
}

func (r *trainStationRepository) DeleteTrainStationById(id uint) error {
	err := r.db.Where("train_id = ?", id).Delete(&models.TrainStation{}).Error
	return err
}
