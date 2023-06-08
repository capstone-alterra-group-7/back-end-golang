package repositories

import (
	"back-end-golang/models"

	"gorm.io/gorm"
)

type TrainRepository interface {
	GetAllTrain() ([]models.Train, error)
	GetAllTrain2(search string) ([]models.Train, error)
	GetAllTrains(sortClassName string, sortByTrainId int) ([]models.TrainCarriage, error)
	GetTrainByID(id uint) (models.Train, error)
	GetTrainByID2(id uint) (models.Train, error)
	TrainStationByTrainID(id uint) (models.TrainStation, error)
	GetTrainStationByTrainID(id uint) ([]models.TrainStation, error)
	SearchTrainAvailable(trainId, originId, destinationId uint) ([]models.TrainStation, error)
	GetStationByID(id uint) (models.Station, error)
	GetStationByID2(id uint) (models.Station, error)
	CreateTrain(train models.Train) (models.Train, error)
	UpdateTrain(train models.Train) (models.Train, error)
	DeleteTrain(id uint) error
}

type trainRepository struct {
	db *gorm.DB
}

func NewTrainRepository(db *gorm.DB) TrainRepository {
	return &trainRepository{db}
}

// Implementasi fungsi-fungsi dari interface ItemRepository

func (r *trainRepository) GetAllTrain() ([]models.Train, error) {
	var (
		trains []models.Train
		count  int64
	)

	err := r.db.Find(&trains).Count(&count).Error

	return trains, err
}

func (r *trainRepository) GetAllTrain2(search string) ([]models.Train, error) {
	var (
		trains []models.Train
	)
	err := r.db.Unscoped().Where("code_train LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%").Find(&trains).Error

	return trains, err
}

func (r *trainRepository) GetAllTrains(sortClassName string, sortByTrainId int) ([]models.TrainCarriage, error) {
	var (
		trains []models.TrainCarriage
		err    error
	)

	if sortClassName != "" && sortByTrainId != 0 {
		err = r.db.Raw(`
		SELECT tc.*
		FROM train_carriages tc
		JOIN trains t ON tc.train_id = t.id
		WHERE tc.class = ? AND train_id = ? AND t.status = 'available' AND tc.id IN (
			SELECT MIN(id)
			FROM train_carriages
			GROUP BY class, train_id
		);
	`, sortClassName, sortByTrainId).Scan(&trains).Error
	} else if sortClassName != "" && sortByTrainId == 0 {
		err = r.db.Raw(`
		SELECT tc.*
		FROM train_carriages tc
		JOIN trains t ON tc.train_id = t.id
		WHERE tc.class = ? AND t.status = 'available' AND tc.id IN (
			SELECT MIN(id)
			FROM train_carriages
			GROUP BY class, train_id
		);
	`, sortClassName).Scan(&trains).Error
	} else if sortClassName == "" && sortByTrainId != 0 {
		err = r.db.Raw(`
		SELECT tc.*
		FROM train_carriages tc
		JOIN trains t ON tc.train_id = t.id
		WHERE train_id = ? AND t.status = 'available' AND tc.id IN (
			SELECT MIN(id)
			FROM train_carriages
			GROUP BY class, train_id
		);
	`, sortByTrainId).Scan(&trains).Error
	} else {
		err = r.db.Raw(`
		SELECT tc.*
		FROM train_carriages tc
		JOIN trains t ON tc.train_id = t.id
		WHERE t.status = 'available' AND tc.id IN (
			SELECT MIN(id)
			FROM train_carriages
			GROUP BY class, train_id
		);
		`).Scan(&trains).Error
	}

	if err != nil {
		return trains, err
	}

	return trains, err
}

func (r *trainRepository) GetTrainByID(id uint) (models.Train, error) {
	var train models.Train
	err := r.db.Unscoped().Where("id = ?", id).First(&train).Error
	return train, err
}
func (r *trainRepository) GetTrainByID2(id uint) (models.Train, error) {
	var train models.Train
	err := r.db.Where("id = ?", id).First(&train).Error
	return train, err
}

func (r *trainRepository) GetTrainStationByTrainID(id uint) ([]models.TrainStation, error) {
	var train []models.TrainStation
	err := r.db.Where("train_id = ?", id).Find(&train).Error
	return train, err
}

func (r *trainRepository) TrainStationByTrainID(id uint) (models.TrainStation, error) {
	var train models.TrainStation
	err := r.db.Where("train_id = ?", id).First(&train).Error
	return train, err
}

func (r *trainRepository) SearchTrainAvailable(trainID, originID, destinationID uint) ([]models.TrainStation, error) {
	var train []models.TrainStation
	err := r.db.Where("train_id = ? AND station_id IN (?, ?)", trainID, originID, destinationID).Find(&train).Error
	if err != nil {
		return nil, err
	}

	return train, nil
}

func (r *trainRepository) GetStationByID(id uint) (models.Station, error) {
	var station models.Station
	err := r.db.Where("id = ?", id).Find(&station).Error
	return station, err
}

func (r *trainRepository) GetStationByID2(id uint) (models.Station, error) {
	var station models.Station
	err := r.db.Unscoped().Where("id = ?", id).First(&station).Error
	return station, err
}

func (r *trainRepository) CreateTrain(train models.Train) (models.Train, error) {
	err := r.db.Create(&train).Error
	return train, err
}

func (r *trainRepository) UpdateTrain(train models.Train) (models.Train, error) {
	err := r.db.Save(&train).Error
	return train, err
}

func (r *trainRepository) DeleteTrain(id uint) error {
	var train models.Train
	err := r.db.Where("id = ?", id).Delete(&train).Error
	if err != nil {
		return err
	}

	trainStation := models.TrainStation{
		TrainID: train.ID,
	}
	err = r.db.Where("train_id = ?", trainStation.TrainID).Delete(&trainStation).Error
	if err != nil {
		return err
	}

	trainCarriage := models.TrainCarriage{
		TrainID: train.ID,
	}
	err = r.db.Where("train_id = ?", trainCarriage.TrainID).Delete(&trainCarriage).Error
	if err != nil {
		return err
	}

	return nil
}
