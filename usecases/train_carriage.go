package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
)

type TrainCarriageUsecase interface {
	GetAllTrainCarriages(page, limit int) ([]dtos.TrainCarriageResponse, int, error)
	GetTrainCarriageByID(id uint) (dtos.TrainCarriageResponse, error)
	CreateTrainCarriage(trainCarriage []dtos.TrainCarriageInput) ([]dtos.TrainCarriageResponse, error)
	UpdateTrainCarriage(id uint, trainCarriageInput dtos.TrainCarriageInput) (dtos.TrainCarriageResponse, error)
	DeleteTrainCarriage(id uint) error
}

type trainCarriageUsecase struct {
	trainCarriageRepo repositories.TrainCarriageRepository
	trainRepo         repositories.TrainRepository
}

func NewTrainCarriageUsecase(TrainCarriageRepo repositories.TrainCarriageRepository, TrainRepo repositories.TrainRepository) TrainCarriageUsecase {
	return &trainCarriageUsecase{TrainCarriageRepo, TrainRepo}
}

// GetAllTrainCarriages godoc
// @Summary      Get all train carriage
// @Description  Get all train carriage
// @Tags         Admin - Train Carriage
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllTrainCarriageStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/train-carriage [get]
func (u *trainCarriageUsecase) GetAllTrainCarriages(page, limit int) ([]dtos.TrainCarriageResponse, int, error) {
	trainCarriages, count, err := u.trainCarriageRepo.GetAllTrainCarriages(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var trainCarriageResponses []dtos.TrainCarriageResponse
	visitedIDs := make(map[uint]map[uint]bool)
	for _, trainCarriage := range trainCarriages {
		trainStationResponses := make([]dtos.TrainStationResponse, 0) // Reset trainStationResponses for each trainCarriage

		train, err := u.trainCarriageRepo.GetTrainByID2(trainCarriage.TrainID)
		if err != nil {
			return trainCarriageResponses, 0, err
		}

		trainSeat, err := u.trainCarriageRepo.GetTrainSeatsByClass(trainCarriage.Class)
		if err != nil {
			return trainCarriageResponses, 0, err
		}

		var trainSeatResponses []dtos.TrainSeatResponse
		for _, trainSeat := range trainSeat {
			trainSeatRespon := dtos.TrainSeatResponse{
				ID:   int(trainSeat.ID),
				Name: trainSeat.Name,
			}
			trainSeatResponses = append(trainSeatResponses, trainSeatRespon)
		}

		getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(train.ID)
		if err != nil {
			return trainCarriageResponses, count, err
		}

		for _, train := range getTrainStation {
			getStation, err := u.trainRepo.GetStationByID(train.StationID)
			if err != nil {
				return trainCarriageResponses, count, err
			}

			if visitedIDs[trainCarriage.ID] == nil {
				visitedIDs[trainCarriage.ID] = make(map[uint]bool)
			}
			if visitedIDs[trainCarriage.ID][getStation.ID] {
				continue
			}
			visitedIDs[trainCarriage.ID][getStation.ID] = true

			trainStationResponse := dtos.TrainStationResponse{
				StationID: train.StationID,
				Station: dtos.StationInput{
					Origin:  getStation.Origin,
					Name:    getStation.Name,
					Initial: getStation.Initial,
				},
				ArriveTime: train.ArriveTime,
			}

			trainStationResponses = append(trainStationResponses, trainStationResponse)
		}

		trainCarriageResponse := dtos.TrainCarriageResponse{
			TrainCarriageID: trainCarriage.ID,
			Train: dtos.TrainResponse{
				TrainID:   trainCarriage.TrainID,
				CodeTrain: train.CodeTrain,
				Name:      train.Name,
				Class:     trainCarriage.Class,
				Price:     trainCarriage.Price,
				Route:     trainStationResponses,
				Status:    train.Status,
			},
			Name:      trainCarriage.Name,
			Seat:      trainSeatResponses,
			CreatedAt: trainCarriage.CreatedAt,
			UpdatedAt: trainCarriage.UpdatedAt,
		}
		trainCarriageResponses = append(trainCarriageResponses, trainCarriageResponse)
	}
	return trainCarriageResponses, count, nil
}

// GetTrainCarriageByID godoc
// @Summary      Get train carriage by ID
// @Description  Get train carriage by ID
// @Tags         Admin - Train Carriage
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train carriage"
// @Success      200 {object} dtos.TrainCarriageStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/train-carriage/{id} [get]
func (u *trainCarriageUsecase) GetTrainCarriageByID(id uint) (dtos.TrainCarriageResponse, error) {
	var trainCarriageResponses dtos.TrainCarriageResponse
	trainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID2(id)
	if err != nil {
		return trainCarriageResponses, err
	}

	train, err := u.trainCarriageRepo.GetTrainByID2(trainCarriage.TrainID)
	if err != nil {
		return trainCarriageResponses, err
	}

	trainSeat, err := u.trainCarriageRepo.GetTrainSeatsByClass(trainCarriage.Class)
	if err != nil {
		return trainCarriageResponses, err
	}

	visitedIDs := make(map[uint]map[uint]bool)

	var trainSeatResponses []dtos.TrainSeatResponse
	var trainStationResponses []dtos.TrainStationResponse
	for _, trainSeat := range trainSeat {
		trainSeatRespon := dtos.TrainSeatResponse{
			ID:   int(trainSeat.ID),
			Name: trainSeat.Name,
		}
		trainSeatResponses = append(trainSeatResponses, trainSeatRespon)
	}

	getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(train.ID)
	if err != nil {
		return trainCarriageResponses, err
	}

	for _, train := range getTrainStation {
		getStation, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainCarriageResponses, err
		}

		if visitedIDs[trainCarriage.ID] == nil {
			visitedIDs[trainCarriage.ID] = make(map[uint]bool)
		}
		if visitedIDs[trainCarriage.ID][getStation.ID] {
			continue
		}
		visitedIDs[trainCarriage.ID][getStation.ID] = true

		trainStationResponse := dtos.TrainStationResponse{
			StationID: train.StationID,
			Station: dtos.StationInput{
				Origin:  getStation.Origin,
				Name:    getStation.Name,
				Initial: getStation.Initial,
			},
			ArriveTime: train.ArriveTime,
		}

		trainStationResponses = append(trainStationResponses, trainStationResponse)
	}
	trainCarriageResponse := dtos.TrainCarriageResponse{
		TrainCarriageID: trainCarriage.ID,
		Train: dtos.TrainResponse{
			TrainID:   trainCarriage.TrainID,
			CodeTrain: train.CodeTrain,
			Name:      train.Name,
			Class:     trainCarriage.Class,
			Price:     trainCarriage.Price,
			Route:     trainStationResponses,
			Status:    train.Status,
		},
		Name:      trainCarriage.Name,
		Seat:      trainSeatResponses,
		CreatedAt: trainCarriage.CreatedAt,
		UpdatedAt: trainCarriage.UpdatedAt,
	}
	return trainCarriageResponse, nil
}

// CreateTrainCarriage godoc
// @Summary      Create a new train carriage
// @Description  Create a new train carriage
// @Tags         Admin - Train Carriage
// @Accept       json
// @Produce      json
// @Param        request body []dtos.TrainCarriageInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.TrainCarriageCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train-carriage [post]
// @Security BearerAuth
func (u *trainCarriageUsecase) CreateTrainCarriage(trainCarriages []dtos.TrainCarriageInput) ([]dtos.TrainCarriageResponse, error) {
	var trainCarriageResponses []dtos.TrainCarriageResponse

	for _, trainCarriageInput := range trainCarriages {
		if trainCarriageInput.TrainID < 1 || trainCarriageInput.Name == "" || trainCarriageInput.Class == "" || trainCarriageInput.Price < 1 {
			return trainCarriageResponses, errors.New("Failed to create train carriage")
		}
		createTrainCarriage := models.TrainCarriage{
			TrainID: trainCarriageInput.TrainID,
			Class:   trainCarriageInput.Class,
			Name:    trainCarriageInput.Name,
			Price:   trainCarriageInput.Price,
		}

		createdTrainCarriage, err := u.trainCarriageRepo.CreateTrainCarriage(createTrainCarriage)
		if err != nil {
			return trainCarriageResponses, err
		}

		train, err := u.trainCarriageRepo.GetTrainByID2(createdTrainCarriage.TrainID)
		if err != nil {
			return trainCarriageResponses, err
		}

		trainSeats, err := u.trainCarriageRepo.GetTrainSeatsByClass(trainCarriageInput.Class)
		if err != nil {
			return trainCarriageResponses, err
		}

		var trainSeatResponses []dtos.TrainSeatResponse
		var trainStationResponses []dtos.TrainStationResponse
		for _, trainSeat := range trainSeats {
			trainSeatResponse := dtos.TrainSeatResponse{
				ID:   int(trainSeat.ID),
				Name: trainSeat.Name,
			}
			trainSeatResponses = append(trainSeatResponses, trainSeatResponse)
		}

		getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(train.ID)
		if err != nil {
			return trainCarriageResponses, err
		}

		for _, train := range getTrainStation {
			getStation, err := u.trainRepo.GetStationByID2(train.StationID)
			if err != nil {
				return trainCarriageResponses, err
			}

			trainStationResponse := dtos.TrainStationResponse{
				StationID: train.StationID,
				Station: dtos.StationInput{
					Origin:  getStation.Origin,
					Name:    getStation.Name,
					Initial: getStation.Initial,
				},
				ArriveTime: train.ArriveTime,
			}

			trainStationResponses = append(trainStationResponses, trainStationResponse)
		}

		trainCarriageResponse := dtos.TrainCarriageResponse{
			TrainCarriageID: createdTrainCarriage.ID,
			Train: dtos.TrainResponse{
				TrainID:   createdTrainCarriage.TrainID,
				CodeTrain: train.CodeTrain,
				Name:      train.Name,
				Class:     createdTrainCarriage.Class,
				Price:     createdTrainCarriage.Price,
				Route:     trainStationResponses,
				Status:    train.Status,
			},
			Name:      createTrainCarriage.Name,
			Seat:      trainSeatResponses,
			CreatedAt: createdTrainCarriage.CreatedAt,
			UpdatedAt: createdTrainCarriage.UpdatedAt,
		}
		trainCarriageResponses = append(trainCarriageResponses, trainCarriageResponse)
	}

	return trainCarriageResponses, nil
}

// UpdateTrainCarriage godoc
// @Summary      Update train carriage
// @Description  Update train carriage
// @Tags         Admin - Train Carriage
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train carriage"
// @Param        request body dtos.TrainCarriageInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.TrainCarriageStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train-carriage/{id} [put]
// @Security BearerAuth
func (u *trainCarriageUsecase) UpdateTrainCarriage(id uint, trainCarriageInput dtos.TrainCarriageInput) (dtos.TrainCarriageResponse, error) {
	var trainCarriage models.TrainCarriage
	var trainCarriageResponse dtos.TrainCarriageResponse
	var trainCarriageResponsee dtos.TrainCarriageResponse

	if trainCarriageInput.TrainID < 1 || trainCarriageInput.Name == "" || trainCarriageInput.Class == "" || trainCarriageInput.Price < 1 {
		return trainCarriageResponse, errors.New("Failed to create train carriage")
	}

	trainCarriage, err := u.trainCarriageRepo.GetTrainCarriageByID(id)

	if err != nil {
		return trainCarriageResponse, err
	}

	trainCarriage.TrainID = trainCarriageInput.TrainID
	trainCarriage.Class = trainCarriageInput.Class
	trainCarriage.Name = trainCarriageInput.Name
	trainCarriage.Price = trainCarriageInput.Price

	trainCarriage, err = u.trainCarriageRepo.UpdateTrainCarriage(trainCarriage)

	if err != nil {
		return trainCarriageResponsee, err
	}

	train, err := u.trainCarriageRepo.GetTrainByID2(trainCarriage.TrainID)
	if err != nil {
		return trainCarriageResponsee, err
	}

	trainSeat, err := u.trainCarriageRepo.GetTrainSeatsByClass(trainCarriage.Class)
	if err != nil {
		return trainCarriageResponsee, err
	}

	var trainStationResponses []dtos.TrainStationResponse
	var trainSeatResponses []dtos.TrainSeatResponse
	for _, trainSeat := range trainSeat {
		trainSeatRespon := dtos.TrainSeatResponse{
			ID:   int(trainSeat.ID),
			Name: trainSeat.Name,
		}
		trainSeatResponses = append(trainSeatResponses, trainSeatRespon)
	}

	getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(train.ID)
	if err != nil {
		return trainCarriageResponsee, err
	}

	for _, train := range getTrainStation {
		getStation, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainCarriageResponsee, err
		}

		trainStationResponse := dtos.TrainStationResponse{
			StationID: train.StationID,
			Station: dtos.StationInput{
				Origin:  getStation.Origin,
				Name:    getStation.Name,
				Initial: getStation.Initial,
			},
			ArriveTime: train.ArriveTime,
		}

		trainStationResponses = append(trainStationResponses, trainStationResponse)
	}
	trainCarriageResponse = dtos.TrainCarriageResponse{
		TrainCarriageID: trainCarriage.ID,
		Train: dtos.TrainResponse{
			TrainID:   trainCarriage.TrainID,
			CodeTrain: train.CodeTrain,
			Name:      train.Name,
			Class:     trainCarriage.Class,
			Price:     trainCarriage.Price,
			Route:     trainStationResponses,
			Status:    train.Status,
		},
		Name:      trainCarriage.Name,
		Seat:      trainSeatResponses,
		CreatedAt: trainCarriage.CreatedAt,
		UpdatedAt: trainCarriage.UpdatedAt,
	}

	return trainCarriageResponse, nil

}

// DeleteTrainCarriage godoc
// @Summary      Delete a train carriage
// @Description  Delete a train carriage
// @Tags         Admin - Train Carriage
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train carriage"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train-carriage/{id} [delete]
// @Security BearerAuth
func (u *trainCarriageUsecase) DeleteTrainCarriage(id uint) error {
	return u.trainCarriageRepo.DeleteTrainCarriage(id)
}
