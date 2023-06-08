package usecases

import (
	"back-end-golang/dtos"
	"back-end-golang/models"
	"back-end-golang/repositories"
	"errors"
	"sort"
	"strings"
)

type StationUsecase interface {
	GetAllStations(page, limit int) ([]dtos.StationResponse, int, error)
	GetAllStationsByAdmin(page, limit int, search, sortBy, filter string) ([]dtos.StationResponse, int, error)
	GetStationByID(id uint) (dtos.StationResponse, error)
	CreateStation(station *dtos.StationInput) (dtos.StationResponse, error)
	UpdateStation(id uint, station dtos.StationInput) (dtos.StationResponse, error)
	DeleteStation(id uint) error
}

type stationUsecase struct {
	stationRepo repositories.StationRepository
}

func NewStationUsecase(StationRepo repositories.StationRepository) StationUsecase {
	return &stationUsecase{StationRepo}
}

// GetAllStations godoc
// @Summary      Get all stations
// @Description  Get all stations
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Success      200 {object} dtos.GetAllStationStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/station [get]
func (u *stationUsecase) GetAllStations(page, limit int) ([]dtos.StationResponse, int, error) {
	stations, count, err := u.stationRepo.GetAllStations(page, limit)
	if err != nil {
		return nil, 0, err
	}

	var stationResponses []dtos.StationResponse
	for _, station := range stations {
		stationResponse := dtos.StationResponse{
			StationID: station.ID,
			Origin:    station.Origin,
			Name:      station.Name,
			Initial:   station.Initial,
			CreatedAt: station.CreatedAt,
			UpdatedAt: station.UpdatedAt,
		}
		stationResponses = append(stationResponses, stationResponse)
	}

	return stationResponses, count, nil
}

// GetAllStations godoc
// @Summary      Get all stations
// @Description  Get all stations
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search data"
// @Param sort_by query string false "Sort by name"
// @Param filter query string false "Filter data"
// @Success      200 {object} dtos.GetAllStationStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/station [get]
// @Security BearerAuth
func (u *stationUsecase) GetAllStationsByAdmin(page, limit int, search, sortBy, filter string) ([]dtos.StationResponse, int, error) {
	stations, count, err := u.stationRepo.GetAllStationsByAdmin(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	filter = strings.ToLower(filter)

	var stationResponses []dtos.StationResponse
	for _, station := range stations {
		deletedStation := ""

		if filter == "inactive" && station.DeletedAt.Time.IsZero() {
			continue
		} else if filter == "active" && !station.DeletedAt.Time.IsZero() {
			continue
		}

		if !station.DeletedAt.Time.IsZero() {
			deletedStation = station.DeletedAt.Time.Format("2006-01-02T15:04:05.000-07:00")
		}
		stationResponse := dtos.StationResponse{
			StationID: station.ID,
			Origin:    station.Origin,
			Name:      station.Name,
			Initial:   station.Initial,
			CreatedAt: station.CreatedAt,
			UpdatedAt: station.UpdatedAt,
			DeletedAt: &deletedStation,
		}
		stationResponses = append(stationResponses, stationResponse)
	}

	// Sort trainResponses based on price
	if strings.ToLower(sortBy) == "asc" {
		sort.SliceStable(stationResponses, func(i, j int) bool {
			return stationResponses[i].Name < stationResponses[j].Name
		})
	} else if strings.ToLower(sortBy) == "desc" {
		sort.SliceStable(stationResponses, func(i, j int) bool {
			return stationResponses[i].Name > stationResponses[j].Name
		})
	}
	return stationResponses, count, nil
}

// GetStationByID godoc
// @Summary      Get station by ID
// @Description  Get station by ID
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param id path integer true "ID station"
// @Success      200 {object} dtos.StationStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /public/station/{id} [get]
func (u *stationUsecase) GetStationByID(id uint) (dtos.StationResponse, error) {
	var stationResponses dtos.StationResponse
	station, err := u.stationRepo.GetStationByID2(id)
	if err != nil {
		return stationResponses, err
	}
	stationResponse := dtos.StationResponse{
		StationID: station.ID,
		Origin:    station.Origin,
		Name:      station.Name,
		Initial:   station.Initial,
		CreatedAt: station.CreatedAt,
		UpdatedAt: station.UpdatedAt,
	}
	return stationResponse, nil
}

// CreateStation godoc
// @Summary      Create a new station
// @Description  Create a new station
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param        request body dtos.StationInput true "Payload Body [RAW]"
// @Success      201 {object} dtos.StationCreeatedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/station [post]
// @Security BearerAuth
func (u *stationUsecase) CreateStation(station *dtos.StationInput) (dtos.StationResponse, error) {
	var stationResponses dtos.StationResponse
	if station.Initial == "" || station.Name == "" || station.Origin == "" {
		return stationResponses, errors.New("Failed to create station")
	}
	createStation := models.Station{
		Origin:  station.Origin,
		Name:    station.Name,
		Initial: station.Initial,
	}

	createdStation, err := u.stationRepo.CreateStation(createStation)
	if err != nil {
		return stationResponses, err
	}

	stationResponse := dtos.StationResponse{
		StationID: createdStation.ID,
		Origin:    createdStation.Origin,
		Name:      createdStation.Name,
		Initial:   createdStation.Initial,
		CreatedAt: createdStation.CreatedAt,
		UpdatedAt: createdStation.UpdatedAt,
	}
	return stationResponse, nil
}

// UpdateStation godoc
// @Summary      Update station
// @Description  Update station
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param id path integer true "ID station"
// @Param        request body dtos.StationInput true "Payload Body [RAW]"
// @Success      200 {object} dtos.StationStatusOKResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/station/{id} [put]
// @Security BearerAuth
func (u *stationUsecase) UpdateStation(id uint, stationInput dtos.StationInput) (dtos.StationResponse, error) {
	var station models.Station
	var stationResponse dtos.StationResponse
	if stationInput.Initial == "" || stationInput.Name == "" || stationInput.Origin == "" {
		return stationResponse, errors.New("Failed to update station")
	}

	station, err := u.stationRepo.GetStationByID(id)
	if err != nil {
		return stationResponse, err
	}

	station.Origin = stationInput.Origin
	station.Name = stationInput.Name
	station.Initial = stationInput.Initial

	station, err = u.stationRepo.UpdateStation(station)

	if err != nil {
		return stationResponse, err
	}

	stationResponse.StationID = station.ID
	stationResponse.Origin = station.Origin
	stationResponse.Name = station.Name
	stationResponse.Initial = station.Initial
	stationResponse.CreatedAt = station.CreatedAt
	stationResponse.UpdatedAt = station.UpdatedAt

	return stationResponse, nil

}

// DeleteStation godoc
// @Summary      Delete a station
// @Description  Delete a station
// @Tags         Admin - Station
// @Accept       json
// @Produce      json
// @Param id path integer true "ID station"
// @Success      200 {object} dtos.StatusOKDeletedResponse
// @Failure      400 {object} dtos.BadRequestResponse
// @Failure      401 {object} dtos.UnauthorizedResponse
// @Failure      403 {object} dtos.ForbiddenResponse
// @Failure      404 {object} dtos.NotFoundResponse
// @Failure      500 {object} dtos.InternalServerErrorResponse
// @Router       /admin/station/{id} [delete]
// @Security BearerAuth
func (u *stationUsecase) DeleteStation(id uint) error {
	return u.stationRepo.DeleteStation(id)
}
