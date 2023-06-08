package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type StationController interface {
	GetAllStations(c echo.Context) error
	GetAllStationsByAdmin(c echo.Context) error
	GetStationByID(c echo.Context) error
	CreateStation(c echo.Context) error
	UpdateStation(c echo.Context) error
	DeleteStation(c echo.Context) error
}

type stationController struct {
	stationUsecase usecases.StationUsecase
}

func NewStationController(stationUsecase usecases.StationUsecase) StationController {
	return &stationController{stationUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *stationController) GetAllStations(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	stations, count, err := c.stationUsecase.GetAllStations(page, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching station",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all stations",
			stations,
			page,
			limit,
			count,
		),
	)
}

func (c *stationController) GetAllStationsByAdmin(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}

	searchParam := ctx.QueryParam("search")
	sortByParam := ctx.QueryParam("sort_by")
	filterParam := ctx.QueryParam("filter")

	stations, count, err := c.stationUsecase.GetAllStationsByAdmin(page, limit, searchParam, sortByParam, filterParam)
	if err != nil {
		return ctx.JSON(
			http.StatusInternalServerError,
			helpers.NewErrorResponse(
				http.StatusInternalServerError,
				"Failed fetching station",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all stations",
			stations,
			page,
			limit,
			count,
		),
	)
}

func (c *stationController) GetStationByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	station, err := c.stationUsecase.GetStationByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get station by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get station by id",
			station,
		),
	)

}

func (c *stationController) CreateStation(ctx echo.Context) error {
	var stationDTO dtos.StationInput
	if err := ctx.Bind(&stationDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding station",
				helpers.GetErrorData(err),
			),
		)
	}

	station, err := c.stationUsecase.CreateStation(&stationDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a station",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a station",
			station,
		),
	)
}

func (c *stationController) UpdateStation(ctx echo.Context) error {

	var stationInput dtos.StationInput
	if err := ctx.Bind(&stationInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed fetching station",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	station, err := c.stationUsecase.GetStationByID(uint(id))
	if station.StationID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get station by id",
				helpers.GetErrorData(err),
			),
		)
	}

	stationResp, err := c.stationUsecase.UpdateStation(uint(id), stationInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed update station",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated station",
			stationResp,
		),
	)
}

func (c *stationController) DeleteStation(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.stationUsecase.DeleteStation(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete station",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted station",
			nil,
		),
	)
}
