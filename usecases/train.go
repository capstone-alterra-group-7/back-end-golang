package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
	"sort"
	"strings"
)

type TrainUsecase interface {
	// admin
	GetAllTrains(page, limit int) ([]dtos.TrainResponses, int, error)
	GetAllTrainsByAdmin(page, limit int, search, sortBy, filter string) ([]dtos.TrainResponses, int, error)
	GetTrainByID(id uint) (dtos.TrainResponses, error)
	CreateTrain(train *dtos.TrainInput) (dtos.TrainResponses, error)
	UpdateTrain(id uint, trainInput dtos.TrainInput) (dtos.TrainResponses, error)
	DeleteTrain(id uint) error

	// user
	SearchTrainAvailable(page, limit, stationOrigin, stationDestination, sortByTrainId int, sortClassName, sortByPrice, sortByArriveTime string) ([]dtos.TrainResponse, int, error)
}

type trainUsecase struct {
	trainRepo        repositories.TrainRepository
	trainStationRepo repositories.TrainStationRepository
}

func NewTrainUsecase(TrainRepo repositories.TrainRepository, TrainStationRepo repositories.TrainStationRepository) TrainUsecase {
	return &trainUsecase{TrainRepo, TrainStationRepo}
}

// =============================== ADMIN ================================== \\

// GetAllTrains godoc
// @Summary      Get all train
// @Description  Get all train
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllTrainStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/train [get]
func (u *trainUsecase) GetAllTrains(page, limit int) ([]dtos.TrainResponses, int, error) {

	trains, err := u.trainRepo.GetAllTrain()
	if err != nil {
		return nil, 0, err
	}

	var trainResponses []dtos.TrainResponses
	for _, train := range trains {
		getTrain, err := u.trainRepo.GetTrainByID(train.ID)
		if err != nil {
			return trainResponses, 0, err
		}

		// if getTrain.Status != "available" {
		// 	continue
		// }

		getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(getTrain.ID)
		if err != nil {
			return trainResponses, 0, err
		}

		var trainStationResponses []dtos.TrainStationResponse

		for _, train := range getTrainStation {
			getStation, err := u.trainRepo.GetStationByID2(train.StationID)
			if err != nil {
				return trainResponses, 0, err
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

		trainResponse := dtos.TrainResponses{
			TrainID:   getTrain.ID,
			CodeTrain: getTrain.CodeTrain,
			Name:      getTrain.Name,
			Route:     trainStationResponses,
			Status:    getTrain.Status,
			CreatedAt: getTrain.CreatedAt,
			UpdatedAt: getTrain.UpdatedAt,
		}
		trainResponses = append(trainResponses, trainResponse)
	}
	// Apply offset and limit to trainResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of trainResponses
	if start >= len(trainResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of trainResponses
	if end > len(trainResponses) {
		end = len(trainResponses)
	}

	subsetTrainResponses := trainResponses[start:end]

	return subsetTrainResponses, len(trainResponses), nil
}

// GetAllTrainsByAdmin godoc
// @Summary      Get all train
// @Description  Get all train
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search data"
// @Param sort_by query string false "Sort by name"
// @Param filter query string false "Filter data"
// @Success      200 {object} dtos.GetAllTrainStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train [get]
// @Security BearerAuth
func (u *trainUsecase) GetAllTrainsByAdmin(page, limit int, search, sortBy, filter string) ([]dtos.TrainResponses, int, error) {

	trains, err := u.trainRepo.GetAllTrain2(search)
	if err != nil {
		return nil, 0, err
	}

	filter = strings.ToLower(filter)

	var trainResponses []dtos.TrainResponses
	for _, train := range trains {
		getTrain, err := u.trainRepo.GetTrainByID(train.ID)
		if err != nil {
			return trainResponses, 0, err
		}
		deletedTrain := ""

		if filter != "" {
			isAvailable := getTrain.Status == "available" && getTrain.DeletedAt.Time.IsZero()

			if filter == "active" && !isAvailable {
				continue
			} else if filter != "active" && isAvailable {
				continue
			}
		}

		if !getTrain.DeletedAt.Time.IsZero() {
			deletedTrain = getTrain.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
		}

		getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(getTrain.ID)
		if err != nil {
			return trainResponses, 0, err
		}

		var trainStationResponses []dtos.TrainStationResponse

		for _, train := range getTrainStation {
			getStation, err := u.trainRepo.GetStationByID2(train.StationID)
			if err != nil {
				return trainResponses, 0, err
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

		trainResponse := dtos.TrainResponses{
			TrainID:   getTrain.ID,
			CodeTrain: strings.ToUpper(getTrain.CodeTrain),
			Name:      strings.ToUpper(getTrain.Name),
			Route:     trainStationResponses,
			Status:    getTrain.Status,
			CreatedAt: getTrain.CreatedAt,
			UpdatedAt: getTrain.UpdatedAt,
			DeletedAt: &deletedTrain,
		}
		trainResponses = append(trainResponses, trainResponse)
	}
	// Apply offset and limit to trainResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of trainResponses
	if start >= len(trainResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of trainResponses
	if end > len(trainResponses) {
		end = len(trainResponses)
	}

	subsetTrainResponses := trainResponses[start:end]

	// Sort trainResponses based on price
	if strings.ToLower(sortBy) == "asc" {
		sort.SliceStable(subsetTrainResponses, func(i, j int) bool {
			return subsetTrainResponses[i].Name < subsetTrainResponses[j].Name
		})
	} else if strings.ToLower(sortBy) == "desc" {
		sort.SliceStable(subsetTrainResponses, func(i, j int) bool {
			return subsetTrainResponses[i].Name > subsetTrainResponses[j].Name
		})
	}
	return subsetTrainResponses, len(trainResponses), nil
}

// GetTrainByID godoc
// @Summary      Get train by ID
// @Description  Get train by ID
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train"
// @Success      200 {object} dtos.TrainStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/train/{id} [get]
func (u *trainUsecase) GetTrainByID(id uint) (dtos.TrainResponses, error) {
	var trainResponses dtos.TrainResponses
	train, err := u.trainRepo.GetTrainByID2(id)
	if err != nil {
		return trainResponses, err
	}

	getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(train.ID)
	if err != nil {
		return trainResponses, err
	}

	var trainStationResponses []dtos.TrainStationResponse

	for _, train := range getTrainStation {
		getStation, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainResponses, err
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

	trainResponse := dtos.TrainResponses{
		TrainID:   train.ID,
		CodeTrain: train.CodeTrain,
		Name:      train.Name,
		Route:     trainStationResponses,
		Status:    train.Status,
		CreatedAt: train.CreatedAt,
		UpdatedAt: train.UpdatedAt,
	}
	return trainResponse, nil
}

// CreateTrain godoc
// @Summary      Create a new train
// @Description  Create a new train
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param        request body dtos.TrainInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.TrainCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train [post]
// @Security BearerAuth
func (u *trainUsecase) CreateTrain(train *dtos.TrainInput) (dtos.TrainResponses, error) {
	var trainResponse dtos.TrainResponses
	if train.CodeTrain == "" || train.Name == "" || train.Route == nil || train.Status == "" {
		return trainResponse, errors.New("Failed to create train")
	}

	createTrain := models.Train{
		CodeTrain: train.CodeTrain,
		Name:      train.Name,
		Status:    train.Status,
	}

	createdTrain, err := u.trainRepo.CreateTrain(createTrain)
	if err != nil {
		return trainResponse, err
	}

	for _, train := range train.Route {
		if train.ArriveTime == "" || train.StationID < 1 {
			return trainResponse, errors.New("Failed to create train")
		}
		station, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainResponse, err
		}
		trainStation := models.TrainStation{
			TrainID:    createdTrain.ID,
			StationID:  station.ID,
			ArriveTime: train.ArriveTime,
		}
		_, err = u.trainStationRepo.CreateTrainStation(trainStation)
		if err != nil {
			return trainResponse, err
		}
	}

	getTrain, err := u.trainRepo.GetTrainByID(createdTrain.ID)
	if err != nil {
		return trainResponse, err
	}

	getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(getTrain.ID)
	if err != nil {
		return trainResponse, err
	}

	var trainStationResponses []dtos.TrainStationResponse
	for _, train := range getTrainStation {
		getStation, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainResponse, err
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

	trainResponse = dtos.TrainResponses{
		TrainID:   getTrain.ID,
		CodeTrain: getTrain.CodeTrain,
		Name:      getTrain.Name,
		Route:     trainStationResponses,
		Status:    getTrain.Status,
		CreatedAt: getTrain.CreatedAt,
		UpdatedAt: getTrain.UpdatedAt,
	}
	return trainResponse, nil
}

// UpdateTrain godoc
// @Summary      Update train
// @Description  Update train
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train"
// @Param        request body dtos.TrainInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.TrainStatusOKResponses
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train/{id} [put]
// @Security BearerAuth
func (u *trainUsecase) UpdateTrain(id uint, train dtos.TrainInput) (dtos.TrainResponses, error) {
	var trains models.Train
	var trainResponse dtos.TrainResponses
	if train.CodeTrain == "" || train.Name == "" || train.Route == nil || train.Status == "" {
		return trainResponse, errors.New("Failed to update train")
	}

	trains, err := u.trainRepo.GetTrainByID(id)
	if err != nil {
		return trainResponse, err
	}

	train.Name = strings.ToUpper(train.Name)
	train.CodeTrain = strings.ToUpper(train.CodeTrain)

	trains.CodeTrain = train.CodeTrain
	trains.Name = train.Name
	trains.Status = train.Status

	createdTrain, err := u.trainRepo.UpdateTrain(trains)
	if err != nil {
		return trainResponse, err
	}

	// trainStation, err := u.trainStationRepo.GetTrainStationByID(trains.ID)

	err = u.trainStationRepo.DeleteTrainStationById(trains.ID)
	if err != nil {
		return trainResponse, err
	}

	for _, train := range train.Route {
		station, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainResponse, err
		}
		trainStation := models.TrainStation{
			TrainID:    createdTrain.ID,
			StationID:  station.ID,
			ArriveTime: train.ArriveTime,
		}
		_, err = u.trainStationRepo.CreateTrainStation(trainStation)
		if err != nil {
			return trainResponse, err
		}
	}

	getTrain, err := u.trainRepo.GetTrainByID(createdTrain.ID)
	if err != nil {
		return trainResponse, err
	}

	getTrainStation, err := u.trainRepo.GetTrainStationByTrainID(getTrain.ID)
	if err != nil {
		return trainResponse, err
	}

	var trainStationResponses []dtos.TrainStationResponse
	for _, train := range getTrainStation {
		getStation, err := u.trainRepo.GetStationByID2(train.StationID)
		if err != nil {
			return trainResponse, err
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

	trainResponse = dtos.TrainResponses{
		TrainID:   getTrain.ID,
		CodeTrain: getTrain.CodeTrain,
		Name:      getTrain.Name,
		Route:     trainStationResponses,
		Status:    getTrain.Status,
		CreatedAt: getTrain.CreatedAt,
		UpdatedAt: getTrain.UpdatedAt,
	}
	return trainResponse, nil
}

// DeleteTrain godoc
// @Summary      Delete a train
// @Description  Delete a train
// @Tags         Admin - Train
// @Accept       json
// @Produce      json
// @Param id path integer true "ID train"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/train/{id} [delete]
// @Security BearerAuth
func (u *trainUsecase) DeleteTrain(id uint) error {
	return u.trainRepo.DeleteTrain(id)
}

// =============================== ADMIN END ================================== \\

// =============================== USER ================================== \\

// SearchTrainAvailable godoc
// @Summary      Search Train Available
// @Description  Search Train
// @Tags         User - Train
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param station_origin_id query int true "Station origin id"
// @Param station_destination_id query int true "Station destination id"
// @Param sort_by_train_id query int false "Filter by train id"
// @Param sort_by_class query string false "Filter by class name"
// @Param sort_by_price query string false "Filter by price"
// @Param sort_by_arrive_time query string false "Filter by arrive time"
// @Success      200 {object} dtos.GetAllTrainStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/train/search [get]
func (u *trainUsecase) SearchTrainAvailable(page, limit, stationOriginId, stationDestinationId, sortByTrainId int, sortClassName, sortByPrice, sortByArriveTime string) ([]dtos.TrainResponse, int, error) {
	trains, err := u.trainRepo.GetAllTrains(sortClassName, sortByTrainId)
	if err != nil {
		return nil, 0, err
	}

	var trainResponses []dtos.TrainResponse

	for _, train := range trains {
		getTrain, err := u.trainRepo.GetTrainByID(train.TrainID)
		if err != nil {
			return trainResponses, 0, err
		}

		getTrainStation, err := u.trainRepo.SearchTrainAvailable(getTrain.ID, uint(stationOriginId), uint(stationDestinationId))
		if err != nil {
			return trainResponses, 0, err
		}

		if getTrain.Status != "available" {
			continue
		}

		// Check if route[0] matches stationOriginId and route[1] matches stationDestinationId
		if len(getTrainStation) < 2 || getTrainStation[0].StationID != uint(stationOriginId) || getTrainStation[1].StationID != uint(stationDestinationId) {
			continue
		}

		// Check if class has already been added or doesn't match the filterClass
		if getTrain.ID != uint(sortByTrainId) && sortByTrainId != 0 {
			continue
		}

		// Check if class name matches the desired class filter
		if strings.ToLower(sortClassName) != "" && strings.ToLower(train.Class) != strings.ToLower(sortClassName) {
			continue
		}

		var trainStationResponses []dtos.TrainStationResponse

		for _, trainStation := range getTrainStation {
			getStation, err := u.trainRepo.GetStationByID2(trainStation.StationID)
			if err != nil {
				return trainResponses, 0, err
			}

			trainStationResponse := dtos.TrainStationResponse{
				StationID: trainStation.StationID,
				Station: dtos.StationInput{
					Origin:  getStation.Origin,
					Name:    getStation.Name,
					Initial: getStation.Initial,
				},
				ArriveTime: trainStation.ArriveTime,
			}

			trainStationResponses = append(trainStationResponses, trainStationResponse)
		}

		trainResponse := dtos.TrainResponse{
			TrainID:   getTrain.ID,
			CodeTrain: getTrain.CodeTrain,
			Name:      getTrain.Name,
			Class:     train.Class,
			Price:     train.Price,
			Route:     trainStationResponses,
			Status:    getTrain.Status,
			CreatedAt: getTrain.CreatedAt,
			UpdatedAt: getTrain.UpdatedAt,
		}
		trainResponses = append(trainResponses, trainResponse)

		// Sort trainResponses based on price
		if strings.ToLower(sortByPrice) == "asc" {
			sort.SliceStable(trainResponses, func(i, j int) bool {
				return trainResponses[i].Price < trainResponses[j].Price
			})
		} else if strings.ToLower(sortByPrice) == "desc" {
			sort.SliceStable(trainResponses, func(i, j int) bool {
				return trainResponses[i].Price > trainResponses[j].Price
			})
		}

		// Sort trainResponses based on arrive time
		if strings.ToLower(sortByArriveTime) == "asc" {
			sort.SliceStable(trainResponses, func(i, j int) bool {
				return trainResponses[i].Route[0].ArriveTime < trainResponses[j].Route[0].ArriveTime
			})
		} else if strings.ToLower(sortByArriveTime) == "desc" {
			sort.SliceStable(trainResponses, func(i, j int) bool {
				return trainResponses[i].Route[0].ArriveTime > trainResponses[j].Route[0].ArriveTime
			})
		}
	}
	// Apply offset and limit to trainResponses
	start := (page - 1) * limit
	end := start + limit

	// Ensure that `start` is within the range of trainResponses
	if start >= len(trainResponses) {
		return nil, 0, nil
	}

	// Ensure that `end` does not exceed the length of trainResponses
	if end > len(trainResponses) {
		end = len(trainResponses)
	}

	subsetTrainResponses := trainResponses[start:end]

	return subsetTrainResponses, len(trainResponses), nil
}
