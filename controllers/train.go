package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TrainController interface {
	//admin
	GetAllTrains(c echo.Context) error
	GetAllTrainsByAdmin(c echo.Context) error
	GetTrainByID(c echo.Context) error
	CreateTrain(c echo.Context) error
	UpdateTrain(c echo.Context) error
	DeleteTrain(c echo.Context) error

	//user
	SearchTrainAvailable(c echo.Context) error
}

type trainController struct {
	trainUsecase usecases.TrainUsecase
}

func NewTrainController(trainUsecase usecases.TrainUsecase) TrainController {
	return &trainController{trainUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

// ============================= ADMIN ==================================== \\

func (c *trainController) GetAllTrains(ctx echo.Context) error {
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

	trains, count, err := c.trainUsecase.GetAllTrains(page, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all train",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all trains",
			trains,
			page,
			limit,
			count,
		),
	)
}

func (c *trainController) GetAllTrainsByAdmin(ctx echo.Context) error {
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

	trains, count, err := c.trainUsecase.GetAllTrainsByAdmin(page, limit, searchParam, sortByParam, filterParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all train",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all trains",
			trains,
			page,
			limit,
			count,
		),
	)
}

func (c *trainController) GetTrainByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	train, err := c.trainUsecase.GetTrainByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get train by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get train by id",
			train,
		),
	)

}

func (c *trainController) CreateTrain(ctx echo.Context) error {
	var trainDTO dtos.TrainInput
	if err := ctx.Bind(&trainDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	train, err := c.trainUsecase.CreateTrain(&trainDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a train",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a train",
			train,
		),
	)
}

func (c *trainController) UpdateTrain(ctx echo.Context) error {

	var trainInput dtos.TrainInput
	if err := ctx.Bind(&trainInput); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	id, _ := strconv.Atoi(ctx.Param("id"))

	train, err := c.trainUsecase.GetTrainByID(uint(id))
	if train.TrainID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get train by id",
				helpers.GetErrorData(err),
			),
		)
	}

	trainResp, err := c.trainUsecase.UpdateTrain(uint(id), trainInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to updated a train",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated train",
			trainResp,
		),
	)
}

func (c *trainController) DeleteTrain(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.trainUsecase.DeleteTrain(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete train",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted train",
			nil,
		),
	)
}

// =============================== ADMIN END ================================== \\

// =============================== USER ================================== \\

func (c *trainController) SearchTrainAvailable(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		page = 1
	}

	limitParam := ctx.QueryParam("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 1000
	}

	classParam := ctx.QueryParam("sort_by_class")

	sortByPriceParam := ctx.QueryParam("sort_by_price")
	sortByArriveTimeParam := ctx.QueryParam("sort_by_arrive_time")

	sortByTrainIdParam := ctx.QueryParam("sort_by_train_id")
	sortByTrainId, _ := strconv.Atoi(sortByTrainIdParam)

	stationOriginIdParam := ctx.QueryParam("station_origin_id")
	stationOriginId, _ := strconv.Atoi(stationOriginIdParam)

	stationDestinationIdParam := ctx.QueryParam("station_destination_id")
	stationDestinationId, _ := strconv.Atoi(stationDestinationIdParam)

	trains, count, err := c.trainUsecase.SearchTrainAvailable(page, limit, stationOriginId, stationDestinationId, sortByTrainId, classParam, sortByPriceParam, sortByArriveTimeParam)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all train",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all trains",
			trains,
			page,
			limit,
			count,
		),
	)
}
