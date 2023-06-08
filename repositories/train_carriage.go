package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TrainCarriageRepository interface {
	GetAllTrainCarriages(page, limit int) ([]models.TrainCarriage, int, error)
	GetTrainCarriageByID(id uint) (models.TrainCarriage, error)
	GetTrainCarriageByID2(id uint) (models.TrainCarriage, error)
	GetTrainCarriageByID3(id uint, class string) ([]models.TrainCarriage, error)
	GetTrainByID2(id uint) (models.Train, error)
	GetStationByID2(id uint) (models.Station, error)
	GetTrainSeatsByClass(class string) ([]models.TrainSeat, error)
	CreateTrainCarriage(trainCarriage models.TrainCarriage) (models.TrainCarriage, error)
	UpdateTrainCarriage(trainCarriage models.TrainCarriage) (models.TrainCarriage, error)
	DeleteTrainCarriage(id uint) error
}

type trainCarriageRepository struct {
	db *gorm.DB
}

func NewTrainCarriageRepository(db *gorm.DB) TrainCarriageRepository {
	return &trainCarriageRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *trainCarriageRepository) GetAllTrainCarriages(page, limit int) ([]models.TrainCarriage, int, error) {
	var (
		trainCarriages []models.TrainCarriage
		count          int64
	)
	err := r.db.Find(&trainCarriages).Count(&count).Error
	if err != nil {
		return trainCarriages, int(count), err
	}

	offset := (page - 1) * limit

	err = r.db.Limit(limit).Offset(offset).Find(&trainCarriages).Error

	return trainCarriages, int(count), err
}

func (r *trainCarriageRepository) GetTrainCarriageByID(id uint) (models.TrainCarriage, error) {
	var trainCarriage models.TrainCarriage
	err := r.db.Unscoped().Where("id = ?", id).First(&trainCarriage).Error
	return trainCarriage, err
}

func (r *trainCarriageRepository) GetTrainCarriageByID2(id uint) (models.TrainCarriage, error) {
	var trainCarriage models.TrainCarriage
	err := r.db.Where("id = ?", id).First(&trainCarriage).Error
	return trainCarriage, err
}

func (r *trainCarriageRepository) GetTrainCarriageByID3(id uint, class string) ([]models.TrainCarriage, error) {
	var trainCarriage []models.TrainCarriage
	err := r.db.Where("train_id = ? AND class = ?", id, class).Find(&trainCarriage).Error
	return trainCarriage, err
}

func (r *trainCarriageRepository) GetTrainByID2(id uint) (models.Train, error) {
	var train models.Train
	err := r.db.Unscoped().Where("id = ?", id).First(&train).Error
	return train, err
}

func (r *trainCarriageRepository) GetStationByID2(id uint) (models.Station, error) {
	var station models.Station
	err := r.db.Where("id = ?", id).First(&station).Error
	return station, err
}

func (r *trainCarriageRepository) GetTrainSeatsByClass(class string) ([]models.TrainSeat, error) {
	var trainSeats []models.TrainSeat
	err := r.db.Where("class = ?", class).Find(&trainSeats).Error
	return trainSeats, err
}

func (r *trainCarriageRepository) CreateTrainCarriage(trainCarriage models.TrainCarriage) (models.TrainCarriage, error) {
	err := r.db.Create(&trainCarriage).Error
	return trainCarriage, err
}

func (r *trainCarriageRepository) UpdateTrainCarriage(trainCarriage models.TrainCarriage) (models.TrainCarriage, error) {
	err := r.db.Save(trainCarriage).Error
	return trainCarriage, err
}

func (r *trainCarriageRepository) DeleteTrainCarriage(id uint) error {
	var trainCarriage models.TrainCarriage
	err := r.db.Where("id = ?", id).Delete(&trainCarriage).Error
	return err
}
