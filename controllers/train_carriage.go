package controllers

import (
	"back-end-golang/dtos"
	"back-end-golang/helpers"
	"back-end-golang/usecases"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TrainCarriageController interface {
	GetAllTrainCarriages(c echo.Context) error
	GetTrainCarriageByID(c echo.Context) error
	CreateTrainCarriage(c echo.Context) error
	UpdateTrainCarriage(c echo.Context) error
	DeleteTrainCarriage(c echo.Context) error
}

type trainCarriageController struct {
	trainCarriageUsecase usecases.TrainCarriageUsecase
}

func NewTrainCarriageController(trainCarriageUsecase usecases.TrainCarriageUsecase) TrainCarriageController {
	return &trainCarriageController{trainCarriageUsecase}
}

// Implementasi fungsi-fungsi dari interface ItemController

func (c *trainCarriageController) GetAllTrainCarriages(ctx echo.Context) error {
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
	trainCarriages, count, err := c.trainCarriageUsecase.GetAllTrainCarriages(page, limit)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get all train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewPaginationResponse(
			http.StatusOK,
			"Successfully get all train carriage",
			trainCarriages,
			page,
			limit,
			count,
		),
	)
}

func (c *trainCarriageController) GetTrainCarriageByID(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))
	trainCarriage, err := c.trainCarriageUsecase.GetTrainCarriageByID(uint(id))

	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get train carriage by id",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully to get train carriage by id",
			trainCarriage,
		),
	)

}

func (c *trainCarriageController) CreateTrainCarriage(ctx echo.Context) error {
	var trainCarriageDTO []dtos.TrainCarriageInput
	if err := ctx.Bind(&trainCarriageDTO); err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed binding train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	trainCarriage, err := c.trainCarriageUsecase.CreateTrainCarriage(trainCarriageDTO)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to created a train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusCreated,
		helpers.NewResponse(
			http.StatusCreated,
			"Successfully to created a train carriage",
			trainCarriage,
		),
	)
}

func (c *trainCarriageController) UpdateTrainCarriage(ctx echo.Context) error {

	var trainCarriageInput dtos.TrainCarriageInput
	if err := ctx.Bind(&trainCarriageInput); err != nil {
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

	trainCarriage, err := c.trainCarriageUsecase.GetTrainCarriageByID(uint(id))
	if trainCarriage.TrainCarriageID == 0 {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to get train carriage by id",
				helpers.GetErrorData(err),
			),
		)
	}

	trainCarriageResp, err := c.trainCarriageUsecase.UpdateTrainCarriage(uint(id), trainCarriageInput)
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to updated a train carriage",
				helpers.GetErrorData(err),
			),
		)
	}

	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully updated train carriage",
			trainCarriageResp,
		),
	)
}

func (c *trainCarriageController) DeleteTrainCarriage(ctx echo.Context) error {
	id, _ := strconv.Atoi(ctx.Param("id"))

	err := c.trainCarriageUsecase.DeleteTrainCarriage(uint(id))
	if err != nil {
		return ctx.JSON(
			http.StatusBadRequest,
			helpers.NewErrorResponse(
				http.StatusBadRequest,
				"Failed to delete train carriage",
				helpers.GetErrorData(err),
			),
		)
	}
	return ctx.JSON(
		http.StatusOK,
		helpers.NewResponse(
			http.StatusOK,
			"Successfully deleted train carriage",
			nil,
		),
	)
}
